package mongo

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/jwtauth"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"identity.io/pkg/cannon"
	"identity.io/pkg/derrors"
	"identity.io/pkg/identity"
	"identity.io/pkg/redis"
)

var _ identity.UserService = &UserService{}

const RiderCollection Collections = "riders"
const DriverCollection Collections = "drivers"

type UserService struct {
	db        *DB
	redis     *redis.Redis
	tokenAuth *jwtauth.JWTAuth
}

func NewUserService(
	db *DB,
	wallet string,
	done chan struct{},
	redis *redis.Redis,
	tokenAuth *jwtauth.JWTAuth,
) *UserService {

	indexes := []mongo.IndexModel{
		{
			Keys: bson.D{{Key: "email", Value: 1}},
			Options: &options.IndexOptions{
				Unique: &[]bool{true}[0],
			},
		},
	}

	_, err := db.Collection(RiderCollection).Indexes().CreateMany(
		context.Background(),
		indexes,
	)
	if err != nil {
		panic("unable to create user email index")
	}
	_, err = db.Collection(DriverCollection).Indexes().CreateMany(
		context.Background(),
		indexes,
	)
	if err != nil {
		panic("unable to create user email index")
	}

	s := &UserService{
		db:        db,
		redis:     redis,
		tokenAuth: tokenAuth,
	}

	return s
}

func (s *UserService) Login(ctx context.Context, email, otp string, refer ...string) (_ *identity.User, err error) {
	defer derrors.Wrap(&err, "mongo.UserService.Login")
	app := identity.ClientFromContext(ctx)
	if app == nil {
		return nil, fmt.Errorf("no application provided: %w", identity.ErrAccessDenied)
	}
	user, err := s.FindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, identity.ErrNotFound) {
			if app == nil {
				return nil, fmt.Errorf("no application provided: %w", identity.ErrAccessDenied)
			}
			user = &identity.User{
				ID:     identity.NewID().String(),
				Email:  email,
				Status: identity.UserStatusOnReview,
			}
			if refer != nil {
				user.Referer = refer[0]
				user.Referal = user.ID[:6]
			}

			switch app.Type {
			case identity.ClientTypeDriver:
				user.Role = identity.RoleDriver
			default:
				user.Role = identity.RoleRider
			}

			err = s.CreateUser(ctx, user)
			if err != nil {
				return nil, err
			}
			user, err = findUserByEmail(ctx, s.db, user.Email)
			if err != nil {
				return nil, err
			}
			return user, nil
		}
		return nil, err
	}

	user.Otp = ""
	if err := updateUser(ctx, s.db, user); err != nil {
		return nil, err
	}

	if !user.IsActive() {
		return nil, identity.ErrAccessDenied
	}

	return user, nil
}

// Token implements identity.UserService.
func (s *UserService) Token(ctx context.Context, user *identity.User) (_ string, err error) {
	defer derrors.Wrap(&err, "mongo.UserService.Token")
	app := identity.ClientFromContext(ctx)
	if app == nil {
		return "", identity.NewUnauthorizedError(identity.ErrAccessDenied)
	}
	_, token, err := s.tokenAuth.Encode(user.Claim())
	if err != nil {
		return "", identity.NewInternalError(err)
	}
	if err := s.redis.Set(ctx, user.ID, token, identity.ExpireIn); err != nil {
		return "", identity.NewInternalError(err)
	}
	return token, nil
}

func (s *UserService) CreateUser(ctx context.Context, user *identity.User) (err error) {
	defer derrors.Wrap(&err, "mongo.UserService.CreateUser")
	app := identity.ClientFromContext(ctx)
	if app == nil {
		return identity.NewUnauthorizedError(identity.ErrAccessDenied)
	}
	logger := cannon.LoggerFromContext(ctx)
	logger.Info("start otp handler")
	if user.Referer != "" {
		user.Referer = identity.NewID().String()[:8]
	}

	tx, err := s.db.client.StartSession()
	if err != nil {
		return fmt.Errorf("unable to start a new session: %w", err)
	}
	defer tx.EndSession(ctx)
	if err := tx.StartTransaction(); err != nil {
		return fmt.Errorf("unable to start a new transaction: %v: %w", err, identity.ErrInternal)
	}
	collection := RiderCollection
	if user.Role == identity.RoleDriver {
		collection = DriverCollection
	}

	_, err = s.db.Collection(collection).InsertOne(ctx, user)
	if err != nil {
		tx.AbortTransaction(ctx)
		return fmt.Errorf("unable to store the user: %w", err)
	}
	return tx.CommitTransaction(ctx)
}

// Logout implements identity.UserService.
func (s *UserService) Logout(ctx context.Context) error {
	token := identity.GetTokenTypeFromContext(ctx)
	if token == "" {
		return identity.NewUnauthorizedError(identity.ErrAccessDenied)
	}
	user := identity.UserFromContext(ctx)
	if user == nil {
		return identity.NewUnauthorizedError(identity.ErrAccessDenied)
	}
	return s.redis.Del(ctx, user.ID)
}

func (s *UserService) FindByID(ctx context.Context, id string) (_ *identity.User, err error) {
	defer derrors.Wrap(&err, "mongo.UserService.FindByID")
	if _, err := checkRole(ctx, s.db, identity.RoleAdmin); err != nil {
		return nil, err
	}
	return findUserByID(ctx, s.db, id)
}

func (s *UserService) FindByEmail(ctx context.Context, email string) (_ *identity.User, err error) {
	defer derrors.Wrap(&err, "mongo.UserService.FindByEmail")
	return findUserByEmail(ctx, s.db, email)
}

func (s *UserService) FindAll(ctx context.Context, filter *identity.UserFilter) (_ *identity.UserList, err error) {
	defer derrors.Wrap(&err, "mongo.UserService.FindAll")
	if _, err := checkRole(ctx, s.db, identity.RoleAdmin); err != nil {
		return nil, err
	}
	users, token, err := findAllUsers(ctx, s.db, filter)
	if err != nil {
		return nil, err
	}
	return &identity.UserList{Data: users, Token: token}, nil
}

func (s *UserService) AddFavoritePlace(ctx context.Context, name string, point identity.Point) (_ *identity.Location, err error) {
	defer derrors.Wrap(&err, "mongo.UserService.AddFavoritePlace")
	user, err := checkRole(ctx, s.db, identity.RoleRider)
	if err != nil {
		return nil, err
	}

	location := identity.Location{
		ID:   identity.NewID().String(),
		Name: name,
	}

	location.Geolocation = identity.GeoLocation{
		Type:        "Point",
		Coordinates: []float64{point.Lng, point.Lat},
		Lat:         point.Lat,
		Long:        point.Lng,
	}
	user.Locations = append(user.Locations, &location)
	if err := updateAddFavoritesPlaces(ctx, s.db, user); err != nil {
		return nil, fmt.Errorf("unable to store the favorite palces: %w", err)
	}
	return &location, nil
}

func (s *UserService) FavoritePlaces(ctx context.Context) (_ []*identity.Location, err error) {
	defer derrors.Wrap(&err, "mongo.UserService.FavoritePlaces")
	usr, err := checkRole(ctx, s.db, identity.RoleRider)
	if err != nil {
		return nil, err
	}
	return usr.Locations, nil
}

func (s *UserService) AddFavoriteVehicle(ctx context.Context, plate string, name *string) (err error) {
	defer derrors.Wrap(&err, "mongo.UserService.AddFavoriteVehicle")
	if plate == "" {
		return fmt.Errorf("plate is required")
	}
	usr, err := checkRole(ctx, s.db, identity.RoleRider)
	if err != nil {
		return err
	}
	usr.AddFavoriteVehicle(plate, *name)
	return nil
}

func (s *UserService) UpdatePlace(ctx context.Context, input *identity.UpdatePlace) (_ *identity.Location, err error) {
	defer derrors.Wrap(&err, "mongo.UserService.UpdatePlace")
	user := identity.UserFromContext(ctx)
	if user == nil {
		return nil, fmt.Errorf("unable to update the place: %w", identity.ErrAccessDenied)
	}
	var location identity.Location
	for i := range user.Locations {
		place := user.Locations[i]
		if place.Name == input.Name {
			place.Geolocation = identity.GeoLocation{
				Type:        "Point",
				Coordinates: []float64{input.Location.Long, input.Location.Lat},
			}
			user.Locations[i] = place
			break
		}
	}
	return &location, nil
}

func (s *UserService) FavoriteVehicles(ctx context.Context) (_ []string, err error) {
	defer derrors.Wrap(&err, "mongo.UserService.FavoriteVehicles")
	usr, err := checkRole(ctx, s.db, identity.RoleRider)
	if err != nil {
		return nil, err
	}
	return usr.GetFavoriteVehicles(), nil
}

func (s *UserService) Me(ctx context.Context) (_ *identity.Profile, err error) {
	defer derrors.Wrap(&err, "mongo.UserService.Me")
	usr, err := checkRole(ctx, s.db, identity.RoleRider)
	if err != nil {
		return nil, err
	}

	return usr.Profile, nil
}

func (s *UserService) LastNAddress(ctx context.Context, number int) (_ []*identity.Location, err error) {
	defer derrors.Wrap(&err, "mongo.UserService.LastNAddress")
	usr, err := checkRole(ctx, s.db, identity.RoleRider)
	if err != nil {
		return nil, err
	}
	return usr.LastNAddress(number), nil
}

func (s *UserService) UpdateProfile(ctx context.Context, request *identity.UpdateProfile) (err error) {
	defer derrors.Wrap(&err, "mongo.UserService.UpdateProfile")
	usr, err := checkRole(ctx, s.db, identity.Role(""))
	if err != nil {
		return err
	}

	if request.Name != "" {
		usr.Profile.Name = request.Name
	}
	if request.LastName != "" {
		usr.Profile.LastName = request.LastName
	}
	if request.Dob != "" {
		usr.Profile.DOB = request.Dob
	}
	if request.Phone != "" {
		usr.Profile.Phone = request.Phone
	}
	if request.Photo != "" {
		usr.Profile.Photo = request.Photo
	}
	if request.Gender != "" {
		usr.Profile.Gender = request.Gender
	}
	if request.Licence != "" {
		usr.Profile.Licence = request.Licence
	}
	if request.Dni != "" {
		usr.Profile.Dni = request.Dni
	}

	if usr.Profile.IsCompleted(usr.Role) {
		usr.Profile.Status = identity.ProfileStatusCompleted
		usr.Status = identity.UserStatusActive
	}

	if err := updateUser(ctx, s.db, usr); err != nil {
		return err
	}

	return nil
}

func (s *UserService) AddDevice(ctx context.Context, device string) (err error) {
	defer derrors.Wrap(&err, "mongo.UserService.AddDevice")
	usr, err := checkRole(ctx, s.db, identity.Role(""))
	if err != nil {
		return err
	}
	if usr.HasDevice(device) {
		return identity.NewError(nil, http.StatusBadRequest, "device already added")
	}
	usr.Devices = append(usr.Devices, identity.Device{Token: device, Active: true})
	return updateUser(ctx, s.db, usr)
}

func (s *UserService) GetUserDevices(ctx context.Context, filter identity.UserFilter) (_ []string, err error) {
	defer derrors.Wrap(&err, "mongo.UserService.GetUserDevices")
	usr, err := checkRole(ctx, s.db, identity.Role(""))
	if err != nil {
		return nil, err
	}
	return usr.GetDevices(), nil
}

func (s *UserService) SetAvailability(ctx context.Context, available bool) (err error) {
	defer derrors.Wrap(&err, "mongo.UserService.SetAvailability")
	usr, err := checkRole(ctx, s.db, identity.RoleDriver)
	if err != nil {
		return err
	}
	usr.Available = available
	return updateUser(ctx, s.db, usr)
}

func (s *UserService) Update(ctx context.Context, user *identity.User) error {
	usr, err := checkRole(ctx, s.db, identity.Role(""))
	if err != nil {
		return err
	}
	if usr.ID != user.ID && usr.Role != identity.RoleAdmin {
		return identity.NewError(nil, http.StatusUnauthorized, "invalid token provided")
	}

	if user.Profile.IsCompleted(user.Role) {
		user.Profile.Status = identity.ProfileStatusCompleted
		user.Status = identity.UserStatusActive
	}
	return updateUser(ctx, s.db, user)
}

// AddVehicle implements identity.UserService.
func (s *UserService) AddVehicle(ctx context.Context, vehicle *identity.Vehicle) (err error) {
	defer derrors.Wrap(&err, "mongo.UserService.AddVehicle")
	use, err := checkRole(ctx, nil, identity.RoleDriver)
	if err != nil {
		return err
	}
	if vehicle.ID == "" {
		vehicle.ID = identity.NewID().String()
	}
	use.Vehicles = append(use.Vehicles, vehicle)
	return updateUser(ctx, s.db, use)
}

// UpdateVehicle implements identity.UserService.
func (s *UserService) UpdateVehicle(ctx context.Context, v *identity.Vehicle) (err error) {
	defer derrors.Wrap(&err, "mongo.UserService.UpdateVehicle")
	use, err := checkRole(ctx, nil, identity.RoleDriver)
	if err != nil {
		return err
	}
	for i, vehicle := range use.Vehicles {
		if vehicle.ID == v.ID {
			use.Vehicles[i] = v
			break
		}
	}
	return updateUser(ctx, s.db, use)
}

// UpdateFavoritePlace implements identity.UserService.
func (s *UserService) UpdateFavoritePlace(ctx context.Context, id string, p identity.UpdatePlace) (err error) {
	defer derrors.Wrap(&err, "mongo.UserService.UpdateFavoritePlace")
	use, err := checkRole(ctx, nil, identity.RoleRider)
	if err != nil {
		return err
	}
	for i, location := range use.Locations {
		if location.ID == id {
			location.Name = p.Name
			location.Geolocation = identity.GeoLocation{
				Type:        "Point",
				Coordinates: []float64{p.Location.Long, p.Location.Lat},
			}
			use.Locations[i] = location
			break
		}
	}
	return updateUser(ctx, s.db, use)
}

// DeleteFavoritePlace implements identity.UserService.
func (s *UserService) DeleteFavoritePlace(ctx context.Context, id string) (err error) {
	defer derrors.Wrap(&err, "mongo.UserService.DeleteFavoritePlace")
	use, err := checkRole(ctx, nil, identity.RoleRider)
	if err != nil {
		return err
	}
	for i, location := range use.Locations {
		if location.ID == id {
			use.Locations = append(use.Locations[:i], use.Locations[i+1:]...)
			break
		}
	}
	return updateUser(ctx, s.db, use)
}

// DeleteVehicle implements identity.UserService.
func (s *UserService) DeleteVehicle(ctx context.Context, id string) (err error) {
	defer derrors.Wrap(&err, "mongo.UserService.DeleteVehicle")
	use, err := checkRole(ctx, nil, identity.RoleDriver)
	if err != nil {
		return err
	}
	for i, vehicle := range use.Vehicles {
		if vehicle.ID == id {
			use.Vehicles = append(use.Vehicles[:i], use.Vehicles[i+1:]...)
			break
		}
	}
	return updateUser(ctx, s.db, use)
}

// DeleteFavoriteVehicle implements identity.UserService.
func (s *UserService) DeleteFavoriteVehicle(ctx context.Context, plate string) error {
	user, err := checkRole(ctx, nil, identity.RoleRider)
	if err != nil {
		return err
	}
	user.DeleteFavoriteVehicle(plate)
	return updateUser(ctx, s.db, user)
}

// SetActiveVehicle implements identity.UserService.
func (s *UserService) SetActiveVehicle(ctx context.Context, plate string) error {
	user, err := checkRole(ctx, nil, identity.RoleRider)
	if err != nil {
		return err
	}

	if !user.SetActiveVehicle(plate) {
		return identity.NewError(nil, http.StatusBadRequest, "vehicle not found")
	}
	return updateUser(ctx, s.db, user)
}

// SetPreferedCurrency implements identity.UserService.
func (s *UserService) SetPreferedCurrency(ctx context.Context, currency string) (err error) {
	defer derrors.Wrap(&err, "mongo.UserService.SetPreferedCurrency")
	user, err := checkRole(ctx, nil, identity.RoleRider)
	if err != nil {
		return err
	}
	if user.Profile == nil {
		user.Profile = &identity.Profile{}
	}
	user.Profile.PreferedCurrency = currency
	return updateUser(ctx, s.db, user)
}

// AddDeviceToken implements identity.UserService.
func (s *UserService) AddDeviceToken(ctx context.Context, token, name string) error {
	user, err := checkRole(ctx, nil, identity.RoleRider)
	if err != nil {
		return err
	}
	if user.HasDevice(token) {
		return identity.NewError(nil, http.StatusBadRequest, "device already added")
	}
	if err := user.AddDevice(token, name); err != nil {
		return err
	}
	return updateUser(ctx, s.db, user)
}

// DeviceTokens implements identity.UserService.
func (s *UserService) DeviceTokens(ctx context.Context) ([]string, error) {
	user, err := checkRole(ctx, nil, identity.Role(""))
	if err != nil {
		return nil, err
	}
	return user.GetDevices(), nil
}

// RemoveDeviceToken implements identity.UserService.
func (s *UserService) RemoveDeviceToken(ctx context.Context, token string) error {
	user, err := checkRole(ctx, nil, identity.Role(""))
	if err != nil {
		return err
	}

	if err := user.RemoveDevice(token); err != nil {
		return err
	}
	return updateUser(ctx, s.db, user)
}

// Vehicle implements identity.UserService.
func (s *UserService) Vehicle(ctx context.Context, vehicle string) (_ *identity.Vehicle, err error) {
	defer derrors.Wrap(&err, "mongo.UserService.Vehicle")
	user, err := checkRole(ctx, nil, identity.RoleDriver)
	if err != nil {
		return nil, err
	}
	for _, v := range user.Vehicles {
		if v.ID == vehicle || v.Plate == vehicle {
			return v, nil
		}
	}
	return nil, identity.NewNotFound("vehicle not found")
}

// Vehicles implements identity.UserService.
func (s *UserService) Vehicles(ctx context.Context) (_ []*identity.Vehicle, err error) {
	defer derrors.Wrap(&err, "mongo.UserService.Vehicles")
	user, err := checkRole(ctx, nil, identity.RoleDriver)
	if err != nil {
		return nil, err
	}
	return user.GetVehicles(), nil
}

// FavoritePlace implements identity.UserService.
func (s *UserService) FavoritePlace(ctx context.Context, name string) (*identity.Location, error) {
	user, err := checkRole(ctx, nil, identity.RoleRider)
	if err != nil {
		return nil, err
	}
	return user.GetFavoritePlace(name), nil
}

func findUserByEmail(ctx context.Context, db *DB, email string) (*identity.User, error) {
	users, _, err := findAllUsers(ctx, db, &identity.UserFilter{Email: email, Limit: 1})
	if err != nil {
		return nil, err
	}
	if len(users) == 0 {
		return nil, identity.ErrNotFound
	}
	return users[0], nil
}

func findUserByID(ctx context.Context, db *DB, id string) (*identity.User, error) {
	users, _, err := findAllUsers(ctx, db, &identity.UserFilter{
		Ids:   []string{id},
		Limit: 1,
	})
	if err != nil {
		return nil, err
	}
	if len(users) == 0 {
		return nil, identity.NewNotFound("user not found")
	}
	return users[0], nil
}

func findAllUsers(ctx context.Context, db *DB, filter *identity.UserFilter) ([]*identity.User, string, error) {
	var users []*identity.User
	var token string
	f := bson.D{}
	if len(filter.Ids) > 0 {
		f = append(f, bson.E{Key: "_id", Value: bson.A{"$in", filter.Ids}})
	}
	if filter.Email != "" {
		f = append(f, bson.E{Key: "email", Value: filter.Email})
	}
	if len(filter.Otp) > 0 {
		f = append(f, bson.E{Key: "otp", Value: filter.Otp})
	}
	if len(filter.Pin) > 0 {
		f = append(f, bson.E{Key: "pin", Value: filter.Email})
	}
	if len(filter.Status) > 0 {
		f = append(f, bson.E{Key: "status", Value: primitive.A{"$in", filter.Status}})
	}
	if len(filter.Role) > 0 {
		f = append(f, bson.E{Key: "role", Value: filter.Role})
	}
	cur, err := db.Collection(RiderCollection).Find(ctx, f)
	if err != nil {
		return nil, "", identity.NewInternalError(fmt.Errorf("mongo error: %v: %w", err, identity.ErrInternal))
	}

	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var user identity.User
		err := cur.Decode(&user)
		if err != nil {
			return nil, "", err
		}
		if user.Role == identity.Role("") {
			user.Role = identity.RoleRider
		}
		users = append(users, &user)

		if len(users) == filter.Limit+1 && filter.Limit > 0 {
			token = users[filter.Limit].ID
			users = users[:filter.Limit]
			return users, token, nil
		}
	}
	cur1, err := db.Collection(DriverCollection).Find(ctx, f)
	if err != nil {
		return nil, "", err
	}
	defer cur1.Close(ctx)
	for cur1.Next(ctx) {
		var user identity.User
		err := cur1.Decode(&user)
		if err != nil {
			return nil, "", err
		}
		if user.Role == identity.Role("") {
			user.Role = identity.RoleDriver
		}
		users = append(users, &user)

		if len(users) == filter.Limit+1 && filter.Limit > 0 {
			token = users[filter.Limit].ID
			users = users[:filter.Limit]
			return users, token, nil
		}
	}
	if len(users) > 0 {
		return users, "", nil
	}
	return nil, "", identity.NewNotFound("users not found")
}

func updateAddFavoritesPlaces(ctx context.Context, db *DB, usr *identity.User) error {
	collection := RiderCollection
	if usr.Role == identity.RoleDriver {
		collection = DriverCollection
	}
	f := bson.D{{Key: "$set", Value: bson.E{Key: "locations", Value: usr.Locations}}}
	_, err := db.Collection(collection).UpdateOne(ctx, bson.D{{Key: "_id", Value: usr.ID}}, f)
	if err != nil {
		return fmt.Errorf("unable to update the location: %w", err)
	}
	return nil
}

func updateUser(ctx context.Context, db *DB, user *identity.User) error {
	collection := RiderCollection
	if user.Role == identity.RoleDriver {
		collection = DriverCollection
	}
	_, err := db.Collection(collection).UpdateOne(ctx, bson.D{{Key: "email", Value: user.Email}}, bson.D{{Key: "$set", Value: user}})
	if err != nil {
		return identity.NewError(identity.ErrInternal, http.StatusInternalServerError, "unable to update the user")
	}
	return nil
}

func getDevices(ctx context.Context, db *DB, filter identity.UserFilter) ([]string, error) {
	user := identity.UserFromContext(ctx)
	if user == nil {
		return nil, identity.NewError(identity.ErrAccessDenied, http.StatusUnauthorized, "invalid token provided")
	}

	var users []identity.User

	f := bson.D{
		{Key: "_id",
			Value: bson.D{
				{Key: "$in",
					Value: bson.A{
						filter.Ids,
					},
				},
			},
		},
	}
	// projection := bson.D{{Key: "devices.id", Value: 1}}
	cur, err := db.Collection(RiderCollection).Find(ctx, f)
	if err != nil {
		return nil, fmt.Errorf("unable to get devices: %v: %w", err, identity.ErrInternal)
	}
	defer cur.Close(ctx)
	var devices []string
	var addedDevices = map[string]bool{}
	for cur.Next(ctx) {
		var token identity.User
		err := cur.Decode(&token)
		if err != nil {
			return nil, fmt.Errorf("unable to decode user: %v: %w", err, identity.ErrInternal)
		}
		users = append(users, token)
		for _, v := range users {
			for _, d := range v.Devices {
				if _, ok := addedDevices[d.Token]; ok {
					continue
				}
				addedDevices[d.Token] = true
				devices = append(devices, d.Token)
			}
		}
	}
	// TODO: simplify this code
	cur, err = db.Collection(DriverCollection).Find(ctx, f)
	if err != nil {
		return nil, fmt.Errorf("unable to get devices: %v: %w", err, identity.ErrInternal)
	}
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var token identity.User
		err := cur.Decode(&token)
		if err != nil {
			return nil, fmt.Errorf("unable to decode user: %v: %w", err, identity.ErrInternal)
		}
		users = append(users, token)
		for _, v := range users {
			for _, d := range v.Devices {
				if _, ok := addedDevices[d.Token]; ok {
					continue
				}
				addedDevices[d.Token] = true
				devices = append(devices, d.Token)
			}
		}
	}

	return devices, nil
}

func CreateUser(ctx context.Context, db *DB, user *identity.User) error {
	collection := RiderCollection
	if user.Role == identity.RoleDriver {
		collection = DriverCollection
	}
	_, err := db.Collection(collection).InsertOne(ctx, user)
	if err != nil {
		return fmt.Errorf("unable to store the user: %w", err)
	}
	return nil
}

func checkRole(ctx context.Context, db *DB, role identity.Role) (*identity.User, error) {
	_, claims, err := jwtauth.FromContext(ctx)
	if err != nil {
		return nil, identity.ErrNilUserInContext
	}
	if claims == nil {
		return nil, identity.ErrNilUserInContext
	}
	if role == identity.Role("") {
		return nil, identity.ErrAccessDenied
	}

	if user, ok := claims["user"]; ok {
		for key, r := range user.(map[string]interface{}) {
			if key == "email" {
				usr, err := findUserByEmail(ctx, db, r.(string))
				if err != nil {
					return nil, err
				}
				return usr, nil
			}

		}
	}

	return nil, identity.ErrAccessDenied
}
