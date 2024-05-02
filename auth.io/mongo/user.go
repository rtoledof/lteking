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

	"auth.io/cannon"
	"auth.io/derrors"
	"auth.io/models"
	"auth.io/redis"
)

var _ models.UserService = &UserService{}

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

	s := &UserService{
		db:        db,
		redis:     redis,
		tokenAuth: tokenAuth,
	}

	return s
}

func (s *UserService) Login(ctx context.Context, email, otp string, refer ...string) (_ *models.User, err error) {
	defer derrors.Wrap(&err, "mongo.UserService.Login")
	app := models.ClientFromContext(ctx)
	if app == nil {
		return nil, fmt.Errorf("no application provided: %w", models.ErrAccessDenied)
	}
	user, err := s.FindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			if app == nil {
				return nil, fmt.Errorf("no application provided: %w", models.ErrAccessDenied)
			}
			user = &models.User{
				ID:     models.NewID().String(),
				Email:  email,
				Status: models.UserStatusOnReview,
			}
			if refer != nil {
				user.Referer = refer[0]
				user.Referal = user.ID[:6]
			}

			switch app.Type {
			case models.ClientTypeDriver:
				user.Role = models.RoleDriver
			default:
				user.Role = models.RoleRider
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
		return nil, models.ErrAccessDenied
	}

	return user, nil
}

// Token implements models.UserService.
func (s *UserService) Token(ctx context.Context, user *models.User) (_ string, err error) {
	defer derrors.Wrap(&err, "mongo.UserService.Token")
	app := models.ClientFromContext(ctx)
	if app == nil {
		return "", models.NewUnauthorizedError(models.ErrAccessDenied)
	}
	_, token, err := s.tokenAuth.Encode(user.Claim())
	if err != nil {
		return "", models.NewInternalError(err)
	}
	if err := s.redis.Set(ctx, user.ID, token, models.ExpireIn); err != nil {
		return "", models.NewInternalError(err)
	}
	return token, nil
}

func (s *UserService) CreateUser(ctx context.Context, user *models.User) (err error) {
	defer derrors.Wrap(&err, "mongo.UserService.CreateUser")
	app := models.ClientFromContext(ctx)
	if app == nil {
		return models.NewUnauthorizedError(models.ErrAccessDenied)
	}
	logger := cannon.LoggerFromContext(ctx)
	logger.Info("start otp handler")
	if user.Referer != "" {
		user.Referer = models.NewID().String()[:8]
	}

	tx, err := s.db.client.StartSession()
	if err != nil {
		return fmt.Errorf("unable to start a new session: %w", err)
	}
	defer tx.EndSession(ctx)
	if err := tx.StartTransaction(); err != nil {
		return fmt.Errorf("unable to start a new transaction: %v: %w", err, models.ErrInternal)
	}
	collection := RiderCollection
	if user.Role == models.RoleDriver {
		collection = DriverCollection
	}

	_, err = s.db.Collection(collection).InsertOne(ctx, user)
	if err != nil {
		tx.AbortTransaction(ctx)
		return fmt.Errorf("unable to store the user: %w", err)
	}
	return tx.CommitTransaction(ctx)
}

// Logout implements models.UserService.
func (s *UserService) Logout(ctx context.Context) error {
	token := models.GetTokenTypeFromContext(ctx)
	if token == "" {
		return models.NewUnauthorizedError(models.ErrAccessDenied)
	}
	user := models.UserFromContext(ctx)
	if user == nil {
		return models.NewUnauthorizedError(models.ErrAccessDenied)
	}
	return s.redis.Del(ctx, user.ID)
}

func (s *UserService) FindByID(ctx context.Context, id string) (_ *models.User, err error) {
	defer derrors.Wrap(&err, "mongo.UserService.FindByID")
	if _, err := checkRole(ctx, s.db, models.RoleAdmin); err != nil {
		return nil, err
	}
	return findUserByID(ctx, s.db, id)
}

func (s *UserService) FindByEmail(ctx context.Context, email string) (_ *models.User, err error) {
	defer derrors.Wrap(&err, "mongo.UserService.FindByEmail")
	return findUserByEmail(ctx, s.db, email)
}

func (s *UserService) FindAll(ctx context.Context, filter *models.UserFilter) (_ *models.UserList, err error) {
	defer derrors.Wrap(&err, "mongo.UserService.FindAll")
	if _, err := checkRole(ctx, s.db, models.RoleAdmin); err != nil {
		return nil, err
	}
	users, token, err := findAllUsers(ctx, s.db, filter)
	if err != nil {
		return nil, err
	}
	return &models.UserList{Data: users, Token: token}, nil
}

func (s *UserService) AddFavoritePlace(ctx context.Context, name string, point models.Point) (_ *models.Location, err error) {
	defer derrors.Wrap(&err, "mongo.UserService.AddFavoritePlace")
	user, err := checkRole(ctx, s.db, models.RoleRider)
	if err != nil {
		return nil, err
	}

	location := models.Location{
		ID:   models.NewID().String(),
		Name: name,
	}

	location.Geolocation = models.GeoLocation{
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

func (s *UserService) FavoritePlaces(ctx context.Context) (_ []*models.Location, err error) {
	defer derrors.Wrap(&err, "mongo.UserService.FavoritePlaces")
	usr, err := checkRole(ctx, s.db, models.RoleRider)
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
	usr, err := checkRole(ctx, s.db, models.RoleRider)
	if err != nil {
		return err
	}
	usr.AddFavoriteVehicle(plate, *name)
	return nil
}

func (s *UserService) UpdatePlace(ctx context.Context, input *models.UpdatePlace) (_ *models.Location, err error) {
	defer derrors.Wrap(&err, "mongo.UserService.UpdatePlace")
	user := models.UserFromContext(ctx)
	if user == nil {
		return nil, fmt.Errorf("unable to update the place: %w", models.ErrAccessDenied)
	}
	var location models.Location
	for i := range user.Locations {
		place := user.Locations[i]
		if place.Name == input.Name {
			place.Geolocation = models.GeoLocation{
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
	usr, err := checkRole(ctx, s.db, models.RoleRider)
	if err != nil {
		return nil, err
	}
	return usr.GetFavoriteVehicles(), nil
}

func (s *UserService) Me(ctx context.Context) (_ *models.User, err error) {
	defer derrors.Wrap(&err, "mongo.UserService.Me")
	usr, err := checkRole(ctx, s.db, models.Role(""))
	if err != nil {
		return nil, err
	}

	return usr, nil
}

func (s *UserService) LastNAddress(ctx context.Context, number int) (_ []*models.Location, err error) {
	defer derrors.Wrap(&err, "mongo.UserService.LastNAddress")
	usr, err := checkRole(ctx, s.db, models.RoleRider)
	if err != nil {
		return nil, err
	}
	return usr.LastNAddress(number), nil
}

func (s *UserService) UpdateProfile(ctx context.Context, request *models.UpdateProfile) (err error) {
	defer derrors.Wrap(&err, "mongo.UserService.UpdateProfile")
	usr, err := checkRole(ctx, s.db, models.Role(""))
	if err != nil {
		return err
	}

	if usr.Profile == nil {
		usr.Profile = &models.Profile{}
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

	if request.PreferedCurrency != "" {
		usr.Profile.PreferedCurrency = request.PreferedCurrency
	}

	if usr.Profile.IsCompleted(usr.Role) {
		usr.Profile.Status = models.ProfileStatusCompleted
		usr.Status = models.UserStatusActive
	}

	if err := updateUser(ctx, s.db, usr); err != nil {
		return err
	}

	return nil
}

func (s *UserService) AddDevice(ctx context.Context, device string) (err error) {
	defer derrors.Wrap(&err, "mongo.UserService.AddDevice")
	usr, err := checkRole(ctx, s.db, models.Role(""))
	if err != nil {
		return err
	}
	if usr.HasDevice(device) {
		return models.NewError(nil, http.StatusBadRequest, "device already added")
	}
	usr.Devices = append(usr.Devices, models.Device{Token: device, Active: true})
	return updateUser(ctx, s.db, usr)
}

func (s *UserService) GetUserDevices(ctx context.Context, filter models.UserFilter) (_ []string, err error) {
	defer derrors.Wrap(&err, "mongo.UserService.GetUserDevices")
	usr, err := checkRole(ctx, s.db, models.Role(""))
	if err != nil {
		return nil, err
	}
	return usr.GetDevices(), nil
}

func (s *UserService) SetAvailability(ctx context.Context, available bool) (err error) {
	defer derrors.Wrap(&err, "mongo.UserService.SetAvailability")
	usr, err := checkRole(ctx, s.db, models.RoleDriver)
	if err != nil {
		return err
	}
	usr.Available = available
	return updateUser(ctx, s.db, usr)
}

func (s *UserService) Update(ctx context.Context, user *models.User) error {
	usr, err := checkRole(ctx, s.db, models.Role(""))
	if err != nil {
		return err
	}
	if usr.ID != user.ID && usr.Role != models.RoleAdmin {
		return models.NewError(nil, http.StatusUnauthorized, "invalid token provided")
	}

	if user.Profile.IsCompleted(user.Role) {
		user.Profile.Status = models.ProfileStatusCompleted
		user.Status = models.UserStatusActive
	}
	return updateUser(ctx, s.db, user)
}

// AddVehicle implements models.UserService.
func (s *UserService) AddVehicle(ctx context.Context, vehicle *models.Vehicle) (err error) {
	defer derrors.Wrap(&err, "mongo.UserService.AddVehicle")
	use, err := checkRole(ctx, nil, models.RoleDriver)
	if err != nil {
		return err
	}
	if vehicle.ID == "" {
		vehicle.ID = models.NewID().String()
	}
	use.Vehicles = append(use.Vehicles, vehicle)
	return updateUser(ctx, s.db, use)
}

// UpdateVehicle implements models.UserService.
func (s *UserService) UpdateVehicle(ctx context.Context, v *models.Vehicle) (err error) {
	defer derrors.Wrap(&err, "mongo.UserService.UpdateVehicle")
	use, err := checkRole(ctx, nil, models.RoleDriver)
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

// UpdateFavoritePlace implements models.UserService.
func (s *UserService) UpdateFavoritePlace(ctx context.Context, id string, p models.UpdatePlace) (err error) {
	defer derrors.Wrap(&err, "mongo.UserService.UpdateFavoritePlace")
	use, err := checkRole(ctx, nil, models.RoleRider)
	if err != nil {
		return err
	}
	for i, location := range use.Locations {
		if location.ID == id {
			location.Name = p.Name
			location.Geolocation = models.GeoLocation{
				Type:        "Point",
				Coordinates: []float64{p.Location.Long, p.Location.Lat},
			}
			use.Locations[i] = location
			break
		}
	}
	return updateUser(ctx, s.db, use)
}

// DeleteFavoritePlace implements models.UserService.
func (s *UserService) DeleteFavoritePlace(ctx context.Context, id string) (err error) {
	defer derrors.Wrap(&err, "mongo.UserService.DeleteFavoritePlace")
	use, err := checkRole(ctx, nil, models.RoleRider)
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

// DeleteVehicle implements models.UserService.
func (s *UserService) DeleteVehicle(ctx context.Context, id string) (err error) {
	defer derrors.Wrap(&err, "mongo.UserService.DeleteVehicle")
	use, err := checkRole(ctx, nil, models.RoleDriver)
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

// DeleteFavoriteVehicle implements models.UserService.
func (s *UserService) DeleteFavoriteVehicle(ctx context.Context, plate string) error {
	user, err := checkRole(ctx, nil, models.RoleRider)
	if err != nil {
		return err
	}
	user.DeleteFavoriteVehicle(plate)
	return updateUser(ctx, s.db, user)
}

// SetActiveVehicle implements models.UserService.
func (s *UserService) SetActiveVehicle(ctx context.Context, plate string) error {
	user, err := checkRole(ctx, nil, models.RoleRider)
	if err != nil {
		return err
	}

	if !user.SetActiveVehicle(plate) {
		return models.NewError(nil, http.StatusBadRequest, "vehicle not found")
	}
	return updateUser(ctx, s.db, user)
}

// SetPreferedCurrency implements models.UserService.
func (s *UserService) SetPreferedCurrency(ctx context.Context, currency string) (err error) {
	defer derrors.Wrap(&err, "mongo.UserService.SetPreferedCurrency")
	user, err := checkRole(ctx, nil, models.RoleRider)
	if err != nil {
		return err
	}
	if user.Profile == nil {
		user.Profile = &models.Profile{}
	}
	user.Profile.PreferedCurrency = currency
	return updateUser(ctx, s.db, user)
}

// AddDeviceToken implements models.UserService.
func (s *UserService) AddDeviceToken(ctx context.Context, token, name string) error {
	user, err := checkRole(ctx, nil, models.RoleRider)
	if err != nil {
		return err
	}
	if user.HasDevice(token) {
		return models.NewError(nil, http.StatusBadRequest, "device already added")
	}
	if err := user.AddDevice(token, name); err != nil {
		return err
	}
	return updateUser(ctx, s.db, user)
}

// DeviceTokens implements models.UserService.
func (s *UserService) DeviceTokens(ctx context.Context) ([]string, error) {
	user, err := checkRole(ctx, nil, models.Role(""))
	if err != nil {
		return nil, err
	}
	return user.GetDevices(), nil
}

// RemoveDeviceToken implements models.UserService.
func (s *UserService) RemoveDeviceToken(ctx context.Context, token string) error {
	user, err := checkRole(ctx, nil, models.Role(""))
	if err != nil {
		return err
	}

	if err := user.RemoveDevice(token); err != nil {
		return err
	}
	return updateUser(ctx, s.db, user)
}

// Vehicle implements models.UserService.
func (s *UserService) Vehicle(ctx context.Context, vehicle string) (_ *models.Vehicle, err error) {
	defer derrors.Wrap(&err, "mongo.UserService.Vehicle")
	user, err := checkRole(ctx, nil, models.RoleDriver)
	if err != nil {
		return nil, err
	}
	for _, v := range user.Vehicles {
		if v.ID == vehicle || v.Plate == vehicle {
			return v, nil
		}
	}
	return nil, models.NewNotFound("vehicle not found")
}

// Vehicles implements models.UserService.
func (s *UserService) Vehicles(ctx context.Context) (_ []*models.Vehicle, err error) {
	defer derrors.Wrap(&err, "mongo.UserService.Vehicles")
	user, err := checkRole(ctx, nil, models.RoleDriver)
	if err != nil {
		return nil, err
	}
	return user.GetVehicles(), nil
}

// FavoritePlace implements models.UserService.
func (s *UserService) FavoritePlace(ctx context.Context, name string) (*models.Location, error) {
	user, err := checkRole(ctx, nil, models.RoleRider)
	if err != nil {
		return nil, err
	}
	return user.GetFavoritePlace(name), nil
}

func findUserByEmail(ctx context.Context, db *DB, email string) (*models.User, error) {
	users, _, err := findAllUsers(ctx, db, &models.UserFilter{Email: email, Limit: 1})
	if err != nil {
		return nil, err
	}
	if len(users) == 0 {
		return nil, models.ErrNotFound
	}
	return users[0], nil
}

func findUserByID(ctx context.Context, db *DB, id string) (*models.User, error) {
	users, _, err := findAllUsers(ctx, db, &models.UserFilter{
		Ids:   []string{id},
		Limit: 1,
	})
	if err != nil {
		return nil, err
	}
	if len(users) == 0 {
		return nil, models.NewNotFound("user not found")
	}
	return users[0], nil
}

func findAllUsers(ctx context.Context, db *DB, filter *models.UserFilter) ([]*models.User, string, error) {
	var users []*models.User
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
		return nil, "", models.NewInternalError(fmt.Errorf("mongo error: %v: %w", err, models.ErrInternal))
	}

	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var user models.User
		err := cur.Decode(&user)
		if err != nil {
			return nil, "", err
		}
		if user.Role == models.Role("") {
			user.Role = models.RoleRider
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
		var user models.User
		err := cur1.Decode(&user)
		if err != nil {
			return nil, "", err
		}
		if user.Role == models.Role("") {
			user.Role = models.RoleDriver
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
	return nil, "", models.NewNotFound("users not found")
}

func updateAddFavoritesPlaces(ctx context.Context, db *DB, usr *models.User) error {
	collection := RiderCollection
	if usr.Role == models.RoleDriver {
		collection = DriverCollection
	}
	f := bson.D{{Key: "$set", Value: bson.E{Key: "locations", Value: usr.Locations}}}
	_, err := db.Collection(collection).UpdateOne(ctx, bson.D{{Key: "_id", Value: usr.ID}}, f)
	if err != nil {
		return fmt.Errorf("unable to update the location: %w", err)
	}
	return nil
}

func updateUser(ctx context.Context, db *DB, user *models.User) error {
	collection := RiderCollection
	if user.Role == models.RoleDriver {
		collection = DriverCollection
	}
	_, err := db.Collection(collection).UpdateOne(ctx, bson.D{{Key: "email", Value: user.Email}}, bson.D{{Key: "$set", Value: user}})
	if err != nil {
		return models.NewError(models.ErrInternal, http.StatusInternalServerError, "unable to update the user")
	}
	return nil
}

func getDevices(ctx context.Context, db *DB, filter models.UserFilter) ([]string, error) {
	user := models.UserFromContext(ctx)
	if user == nil {
		return nil, models.NewError(models.ErrAccessDenied, http.StatusUnauthorized, "invalid token provided")
	}

	var users []models.User

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
		return nil, fmt.Errorf("unable to get devices: %v: %w", err, models.ErrInternal)
	}
	defer cur.Close(ctx)
	var devices []string
	var addedDevices = map[string]bool{}
	for cur.Next(ctx) {
		var token models.User
		err := cur.Decode(&token)
		if err != nil {
			return nil, fmt.Errorf("unable to decode user: %v: %w", err, models.ErrInternal)
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
		return nil, fmt.Errorf("unable to get devices: %v: %w", err, models.ErrInternal)
	}
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var token models.User
		err := cur.Decode(&token)
		if err != nil {
			return nil, fmt.Errorf("unable to decode user: %v: %w", err, models.ErrInternal)
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

func CreateUser(ctx context.Context, db *DB, user *models.User) error {
	collection := RiderCollection
	if user.Role == models.RoleDriver {
		collection = DriverCollection
	}
	if user.Referer != "" {
		user.Referer = models.NewReferalCode()
	}
	_, err := db.Collection(collection).InsertOne(ctx, user)
	if err != nil {
		return fmt.Errorf("unable to store the user: %w", err)
	}
	return nil
}

func checkRole(ctx context.Context, db *DB, role models.Role) (*models.User, error) {
	_, claims, err := jwtauth.FromContext(ctx)
	if err != nil {
		return nil, models.ErrNilUserInContext
	}
	if claims == nil {
		return nil, models.ErrNilUserInContext
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

	return nil, models.ErrAccessDenied
}
