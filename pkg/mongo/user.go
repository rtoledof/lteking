package mongo

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"cubawheeler.io/pkg/cannon"
	"cubawheeler.io/pkg/client/wallet"
	"cubawheeler.io/pkg/cubawheeler"
	"cubawheeler.io/pkg/derrors"
)

var _ cubawheeler.UserService = &UserService{}

const UsersCollection Collections = "users"

type UserService struct {
	db     *DB
	wallet string
}

func NewUserService(
	db *DB,
	wallet string,
	done chan struct{},
) *UserService {

	indexes := []mongo.IndexModel{
		{
			Keys: bson.D{{Key: "email", Value: 1}},
			Options: &options.IndexOptions{
				Unique: &[]bool{true}[0],
			},
		},
	}
	db.client.Database(database).Collection(UsersCollection.String()).Indexes().DropAll(context.Background())

	_, err := db.Collection(UsersCollection).Indexes().CreateMany(
		context.Background(),
		indexes,
	)
	if err != nil {
		panic("unable to create user email index")
	}

	s := &UserService{
		db:     db,
		wallet: wallet,
	}

	return s
}

func (s *UserService) Login(ctx context.Context, input cubawheeler.LoginRequest) (_ *cubawheeler.User, err error) {
	defer derrors.Wrap(&err, "mongo.UserService.Login")
	app := cubawheeler.ClientFromContext(ctx)
	user, err := s.FindByEmail(ctx, input.Email)
	if err != nil {
		if errors.Is(err, cubawheeler.ErrNotFound) {
			if app == nil {
				return nil, fmt.Errorf("no application provided: %w", cubawheeler.ErrAccessDenied)
			}
			user = &cubawheeler.User{
				ID:      cubawheeler.NewID().String(),
				Email:   input.Email,
				Status:  cubawheeler.UserStatusOnReview,
				Referal: input.Referer,
			}

			switch app.Type {
			case cubawheeler.ApplicationTypeDriver:
				user.Role = cubawheeler.RoleDriver
			default:
				user.Role = cubawheeler.RoleRider
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
		return nil, cubawheeler.ErrAccessDenied
	}

	return user, nil
}

func (s *UserService) CreateUser(ctx context.Context, user *cubawheeler.User) (err error) {
	defer derrors.Wrap(&err, "mongo.UserService.CreateUser")
	logger := cannon.LoggerFromContext(ctx)
	logger.Info("start otp handler")
	if user.Referer != "" {
		user.Referer = cubawheeler.NewID().String()[:8]
	}

	tx, err := s.db.client.StartSession()
	if err != nil {
		return fmt.Errorf("unable to start a new session: %w", err)
	}
	defer tx.EndSession(ctx)
	if err := tx.StartTransaction(); err != nil {
		return fmt.Errorf("unable to start a new transaction: %v: %w", err, cubawheeler.ErrInternal)
	}
	if err := createWallet(ctx, s.wallet, user.ID); err != nil {
		tx.AbortTransaction(ctx)
		return err
	}
	_, err = s.db.Collection(UsersCollection).InsertOne(ctx, user)
	if err != nil {
		tx.AbortTransaction(ctx)
		return fmt.Errorf("unable to store the user: %w", err)
	}
	return tx.CommitTransaction(ctx)
}

func (s *UserService) FindByID(ctx context.Context, id string) (_ *cubawheeler.User, err error) {
	defer derrors.Wrap(&err, "mongo.UserService.FindByID")
	usr := cubawheeler.UserFromContext(ctx)
	if usr == nil {
		return nil, errors.New("invalid token provided")
	}
	if usr.Role != cubawheeler.RoleAdmin {
		return usr, nil
	}
	return findUserByID(ctx, s.db, id)
}

func (s *UserService) FindByEmail(ctx context.Context, email string) (_ *cubawheeler.User, err error) {
	defer derrors.Wrap(&err, "mongo.UserService.FindByEmail")
	return findUserByEmail(ctx, s.db, email)
}

func (s *UserService) FindAll(ctx context.Context, filter *cubawheeler.UserFilter) (_ *cubawheeler.UserList, err error) {
	defer derrors.Wrap(&err, "mongo.UserService.FindAll")
	user := cubawheeler.UserFromContext(ctx)
	if user == nil {
		return nil, errors.New("invalid token provided")
	}
	if user.Role != cubawheeler.RoleAdmin {
		return nil, errors.New("access denied")
	}
	users, token, err := findAllUsers(ctx, s.db, filter)
	if err != nil {
		return nil, err
	}
	return &cubawheeler.UserList{Data: users, Token: token}, nil
}

func (s *UserService) AddFavoritePlace(ctx context.Context, input cubawheeler.AddPlace) (_ *cubawheeler.Location, err error) {
	defer derrors.Wrap(&err, "mongo.UserService.AddFavoritePlace")
	usr := cubawheeler.UserFromContext(ctx)
	if usr == nil {
		return nil, errors.New("invalid token provided")
	}
	if usr.Role != cubawheeler.RoleRider {
		return nil, errors.New("access denied")
	}

	location := cubawheeler.Location{
		ID:   cubawheeler.NewID().String(),
		Name: input.Name,
	}
	if input.Location != nil {
		location.Geolocation = cubawheeler.GeoLocation{
			Type:        "Point",
			Coordinates: []float64{input.Location.Long, input.Location.Lat},
		}
	}
	usr.Locations = append(usr.Locations, &location)
	if err := updateAddFavoritesPlaces(ctx, s.db, usr); err != nil {
		return nil, fmt.Errorf("unable to store the favorite palces: %w", err)
	}
	return &location, nil
}

func (s *UserService) FavoritePlaces(ctx context.Context) (_ []*cubawheeler.Location, err error) {
	defer derrors.Wrap(&err, "mongo.UserService.FavoritePlaces")
	usr := cubawheeler.UserFromContext(ctx)
	if usr == nil {
		return nil, fmt.Errorf("nil user in context: %w", cubawheeler.ErrAccessDenied)
	}
	return usr.Locations, nil
}

func (s *UserService) AddFavoriteVehicle(ctx context.Context, plate *string) (_ *cubawheeler.Vehicle, err error) {
	defer derrors.Wrap(&err, "mongo.UserService.AddFavoriteVehicle")
	usr := cubawheeler.UserFromContext(ctx)
	if usr == nil {
		return nil, fmt.Errorf("nil user in context: %w", cubawheeler.ErrAccessDenied)
	}
	vehicle, err := findVehicleByPlate(ctx, s.db, *plate)
	if err != nil {
		return nil, err
	}
	usr.FavoriteVehicles = append(usr.FavoriteVehicles, vehicle)
	if err := updateUser(ctx, s.db, usr); err != nil {
		return nil, err
	}
	return vehicle, nil
}

func (s *UserService) UpdatePlace(ctx context.Context, input *cubawheeler.UpdatePlace) (_ *cubawheeler.Location, err error) {
	defer derrors.Wrap(&err, "mongo.UserService.UpdatePlace")
	user := cubawheeler.UserFromContext(ctx)
	if user == nil {
		return nil, fmt.Errorf("unable to update the place: %w", cubawheeler.ErrAccessDenied)
	}
	var location cubawheeler.Location
	for i := range user.Locations {
		place := user.Locations[i]
		if place.Name == input.Name {
			place.Geolocation = cubawheeler.GeoLocation{
				Type:        "Point",
				Coordinates: []float64{input.Location.Long, input.Location.Lat},
			}
			user.Locations[i] = place
			break
		}
	}
	return &location, nil
}

func (s *UserService) FavoriteVehicles(ctx context.Context) (_ []*cubawheeler.Vehicle, err error) {
	defer derrors.Wrap(&err, "mongo.UserService.FavoriteVehicles")
	usr := cubawheeler.UserFromContext(ctx)
	if usr == nil {
		return nil, fmt.Errorf("nil user in context: %w", cubawheeler.ErrAccessDenied)
	}
	return usr.FavoriteVehicles, nil
}

func (s *UserService) Me(ctx context.Context) (*cubawheeler.Profile, error) {
	user := cubawheeler.UserFromContext(ctx)
	if user == nil {
		return nil, fmt.Errorf("nil user in context: %w", cubawheeler.ErrAccessDenied)
	}

	return &user.Profile, nil
}

func (s *UserService) Orders(ctx context.Context, filter *cubawheeler.OrderFilter) (*cubawheeler.OrderList, error) {
	usr := cubawheeler.UserFromContext(ctx)
	if usr == nil {
		return nil, cubawheeler.ErrNotFound
	}

	orders, token, err := findOrders(ctx, s.db, filter)
	if err != nil {
		return nil, err
	}
	return &cubawheeler.OrderList{Data: orders, Token: token}, nil
}

func (s *UserService) LastNAddress(ctx context.Context, number int) (_ []*cubawheeler.Location, err error) {
	defer derrors.Wrap(&err, "mongo.UserService.LastNAddress")
	usr := cubawheeler.UserFromContext(ctx)
	if usr != nil {
		return nil, errors.New("invalid token profived")
	}
	panic("implement me")
}

func (s *UserService) UpdateProfile(ctx context.Context, request *cubawheeler.UpdateProfile) (err error) {
	defer derrors.Wrap(&err, "mongo.UserService.UpdateProfile")
	usr := cubawheeler.UserFromContext(ctx)
	if usr == nil {
		return errors.New("invalid token provided")
	}
	usr, err = findUserByEmail(ctx, s.db, usr.Email)
	if err != nil {
		return err
	}

	if request.Name != nil {
		usr.Profile.Name = *request.Name
	}
	if request.LastName != nil {
		usr.Profile.LastName = *request.LastName
	}
	if request.Dob != nil {
		usr.Profile.DOB = *request.Dob
	}
	if request.Phone != nil {
		usr.Profile.Phone = *request.Phone
	}
	if request.Photo != nil {
		usr.Profile.Photo = *request.Photo
	}
	if request.Gender != nil {
		usr.Profile.Gender = *request.Gender
	}
	if request.Licence != nil {
		usr.Profile.Licence = *request.Licence
	}
	if request.Dni != nil {
		usr.Profile.Dni = *request.Dni
	}

	if usr.Profile.IsCompleted(usr.Role) {
		usr.Profile.Status = cubawheeler.ProfileStatusCompleted
		usr.Status = cubawheeler.UserStatusActive
	}

	if err := updateUser(ctx, s.db, usr); err != nil {
		return err
	}

	return nil
}

func (s *UserService) AddDevice(ctx context.Context, device string) (err error) {
	defer derrors.Wrap(&err, "mongo.UserService.AddDevice")
	usr := cubawheeler.UserFromContext(ctx)
	if usr == nil {
		return cubawheeler.ErrNilUserInContext
	}
	usr, err = findUserByEmail(ctx, s.db, usr.Email)
	if err != nil {
		return err
	}
	usr.Devices = append(usr.Devices, cubawheeler.Device{Token: device, Active: true})
	return updateUser(ctx, s.db, usr)
}

func (s *UserService) GetUserDevices(ctx context.Context, users []string) (_ []string, err error) {
	defer derrors.Wrap(&err, "mongo.UserService.GetUserDevices")
	usr := cubawheeler.UserFromContext(ctx)
	if usr == nil {
		return nil, cubawheeler.ErrNilUserInContext
	}
	if usr.Role != cubawheeler.RoleAdmin {
		return nil, cubawheeler.ErrAccessDenied
	}
	return getDevices(ctx, s.db, users)
}

func (s *UserService) SetAvailability(ctx context.Context, user string, available bool) (err error) {
	defer derrors.Wrap(&err, "mongo.UserService.SetAvailability")
	usr, err := findUserByID(ctx, s.db, user)
	if err != nil {
		return err
	}
	usr.Available = available
	return updateUser(ctx, s.db, usr)
}

func (s *UserService) Update(ctx context.Context, user *cubawheeler.User) error {
	if user.Profile.IsCompleted(user.Role) {
		user.Profile.Status = cubawheeler.ProfileStatusCompleted
		user.Status = cubawheeler.UserStatusActive
	}
	return updateUser(ctx, s.db, user)
}

func findUserByEmail(ctx context.Context, db *DB, email string) (*cubawheeler.User, error) {
	users, _, err := findAllUsers(ctx, db, &cubawheeler.UserFilter{Email: email, Limit: 1})
	if err != nil {
		return nil, err
	}
	if len(users) == 0 {
		return nil, cubawheeler.ErrNotFound
	}
	return users[0], nil
}

func findUserByID(ctx context.Context, db *DB, id string) (*cubawheeler.User, error) {
	users, _, err := findAllUsers(ctx, db, &cubawheeler.UserFilter{
		Ids:   []string{id},
		Limit: 1,
	})
	if err != nil {
		return nil, err
	}
	if len(users) == 0 {
		return nil, errors.New("user not found")
	}
	return users[0], nil
}

func findAllUsers(ctx context.Context, db *DB, filter *cubawheeler.UserFilter) ([]*cubawheeler.User, string, error) {
	collection := db.Collection(UsersCollection)
	var users []*cubawheeler.User
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
	cur, err := collection.Find(ctx, f)
	if err != nil {
		return nil, "", err
	}
	for cur.Next(ctx) {
		var user cubawheeler.User
		err := cur.Decode(&user)
		if err != nil {
			return nil, "", err
		}
		if user.Role == cubawheeler.Role("") {
			user.Role = cubawheeler.RoleRider
		}
		users = append(users, &user)

		if len(users) == filter.Limit+1 && filter.Limit > 0 {
			token = users[filter.Limit].ID
			users = users[:filter.Limit]
			break
		}
	}

	if err := cur.Err(); err != nil {
		return nil, "", err
	}
	cur.Close(ctx)
	return users, token, nil
}

func updateAddFavoritesPlaces(ctx context.Context, db *DB, usr *cubawheeler.User) error {
	collection := db.client.Database(database).Collection(UsersCollection.String())
	f := bson.D{{Key: "$set", Value: bson.E{Key: "locations", Value: usr.Locations}}}
	_, err := collection.UpdateOne(ctx, bson.D{{Key: "_id", Value: usr.ID}}, f)
	if err != nil {
		return fmt.Errorf("unable to update the location: %w", err)
	}
	return nil
}

func updateUser(ctx context.Context, db *DB, user *cubawheeler.User) error {
	collection := db.client.Database(database).Collection(UsersCollection.String())
	_, err := collection.UpdateOne(ctx, bson.D{{Key: "email", Value: user.Email}}, bson.D{{Key: "$set", Value: user}})
	if err != nil {
		return fmt.Errorf("unable to update the user: %v: %w", err, cubawheeler.ErrInternal)
	}
	return nil
}

func getDevices(ctx context.Context, db *DB, users []string) ([]string, error) {
	collection := db.Collection(UsersCollection)
	var devices []string

	f := bson.D{
		{Key: "_id", Value: bson.A{"$in", users}},
		{Key: "devices.active", Value: true},
	}
	projection := bson.D{{Key: "devices.token", Value: 1}}
	cur, err := collection.Find(ctx, f, &options.FindOptions{Projection: projection})
	if err != nil {
		return nil, fmt.Errorf("unable to get devices: %v: %w", err, cubawheeler.ErrInternal)
	}
	for cur.Next(ctx) {
		var token string
		err := cur.Decode(&token)
		if err != nil {
			return nil, fmt.Errorf("unable to decode user: %v: %w", err, cubawheeler.ErrInternal)
		}
		devices = append(devices, token)
	}
	return devices, nil
}

func createWallet(ctx context.Context, walletURL, owner string) error {
	logger := cannon.LoggerFromContext(ctx)
	logger.Info("start create wallet")
	token := cubawheeler.JWTFromContext(ctx)
	if token == "" {
		logger.Info("nil token in context")
		return cubawheeler.NewError(cubawheeler.ErrAccessDenied, http.StatusUnauthorized, "token is nil")
	}
	walletTransport := wallet.AuthTransport{
		Token: token,
	}
	walletClient, err := wallet.NewClient(walletTransport.Client(), walletURL)
	if err != nil {
		logger.Info(fmt.Sprintf("create waller: %v", err))
		return fmt.Errorf("error creating wallet client: %v: %w", err, cubawheeler.ErrInternal)
	}
	if _, err = walletClient.Service.Create(ctx, wallet.CreateRequest{
		Owner: owner,
	}); err != nil {
		logger.Info(fmt.Sprintf("create wallet: %v", err))
		return fmt.Errorf("error creating wallet: %v: %w", err, cubawheeler.ErrInternal)
	}
	return nil
}
