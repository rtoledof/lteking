package mongo

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"cubawheeler.io/pkg/cubawheeler"
	"cubawheeler.io/pkg/derrors"
	"cubawheeler.io/pkg/pusher"
	"cubawheeler.io/pkg/redis"
)

var _ cubawheeler.UserService = &UserService{}

const UsersCollection Collections = "users"

type UserService struct {
	db         *DB
	collection *mongo.Collection
	beansToken *redis.BeansToken
	beans      *pusher.PushNotification
}

func NewUserService(
	db *DB,
	beansToken *redis.BeansToken,
	beans *pusher.PushNotification,
	done chan struct{},
) *UserService {

	_, err := db.client.Database(database).Collection(UsersCollection.String()).Indexes().CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys: bson.D{{Key: "email", Value: "text"}},
		})
	if err != nil {
		panic("unable to create user email index")
	}

	s := &UserService{
		db:         db,
		collection: db.client.Database(database).Collection(UsersCollection.String()),
		beansToken: beansToken,
		beans:      beans,
	}

	//	ticker := time.NewTicker(time.Hour * 1)
	//	go func() {
	//
	//		for {
	//			select {
	//			case <-ticker.C:
	//				if err := s.generateTokens(context.Background()); err != nil {
	//					slog.Info(err.Error())
	//				}
	//			case <-done:
	//				return
	//			}
	//		}
	//
	//	}()
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
				ID:     cubawheeler.NewID().String(),
				Email:  input.Email,
				Status: cubawheeler.UserStatusOnReview,
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

	if !user.IsActive() {
		return nil, cubawheeler.ErrAccessDenied
	}

	return user, nil
}

func (s *UserService) CreateUser(ctx context.Context, user *cubawheeler.User) (err error) {
	defer derrors.Wrap(&err, "mongo.UserService.CreateUser")
	_, err = s.collection.InsertOne(ctx, user)
	if err != nil {
		return fmt.Errorf("unable to store the user: %w", err)
	}
	if _, err := createProfile(ctx, s.db, &cubawheeler.UpdateProfile{}, user); err != nil {
		return err
	}
	return nil
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
	usr := cubawheeler.UserFromContext(ctx)
	if usr == nil {
		return nil, errors.New("invalid token provided")
	}

	return &usr.Profile, nil
}

func (s *UserService) Orders(ctx context.Context, filter *cubawheeler.OrderFilter) (*cubawheeler.OrderList, error) {
	usr := cubawheeler.UserFromContext(ctx)
	if usr == nil {
		return nil, cubawheeler.ErrNotFound
	}

	ordersCollection := s.db.client.Database(database).Collection(OrderCollection.String())
	orders, token, err := findOrders(ctx, ordersCollection, filter)
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
	params := bson.D{}
	if request.Name != nil {
		params = append(params, primitive.E{Key: "name", Value: *request.Name})
		usr.Profile.Name = *request.Name
	}
	if request.LastName != nil {
		params = append(params, primitive.E{Key: "profile.last_name", Value: *request.LastName})
		usr.Profile.LastName = *request.LastName
	}
	if request.Dob != nil {
		params = append(params, primitive.E{Key: "profile.dob", Value: *request.Dob})
		usr.Profile.DOB = *request.Dob
	}
	if request.Phone != nil {
		params = append(params, primitive.E{Key: "profile.phone", Value: *request.Phone})
		usr.Profile.Phone = *request.Phone
	}
	if request.Photo != nil {
		params = append(params, primitive.E{Key: "profile.photo", Value: *request.Photo})
		usr.Profile.Photo = *request.Photo
	}
	if request.Gender != nil {
		params = append(params, primitive.E{Key: "profile.gender", Value: *request.Gender})
		usr.Profile.Gender = *request.Gender
	}
	if request.Licence != nil {
		params = append(params, primitive.E{Key: "profile.licence", Value: *request.Licence})
		usr.Profile.Licence = *request.Licence
	}
	if request.Dni != nil {
		params = append(params, primitive.E{Key: "profile.dni", Value: *request.Dni})
		usr.Profile.Dni = *request.Dni
	}

	if usr.Profile.IsCompleted(usr.Role) {
		usr.Profile.Status = cubawheeler.ProfileStatusCompleted
		usr.Status = cubawheeler.UserStatusActive
		params = append(params, primitive.E{Key: "profile.status", Value: usr.Profile.Status})
		params = append(params, primitive.E{Key: "status", Value: usr.Status})
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
	usr.Devices = append(usr.Devices, cubawheeler.Device{Token: device, Active: true})
	return updateUser(ctx, s.db, usr)
}

func (s *UserService) GetUserDevices(ctx context.Context, users []string) (_ []*cubawheeler.User, err error) {
	defer derrors.Wrap(&err, "mongo.UserService.GetUserDevices")
	panic("implement me")
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

func (s *UserService) generateTokens(ctx context.Context) error {
	users, _, err := findAllUsers(ctx, s.db, &cubawheeler.UserFilter{
		Status: []cubawheeler.UserStatus{cubawheeler.UserStatusActive},
	})
	if err != nil {
		return err
	}
	for _, v := range users {
		if _, err := s.beansToken.GetBeansToken(ctx, v.ID); err != nil && errors.Is(cubawheeler.ErrNotFound, err) {
			token, err := s.beans.GenerateToken(v.ID)
			if err != nil {
				slog.Info("unabe to generate a new beans token: %v", err)
				continue
			}
			if err := s.beansToken.StoreBeansToken(ctx, v.ID, token); err != nil {
				slog.Info("unabe to store a beans token: %v", err)
			}
		}
	}
	return nil
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
	collection := db.client.Database(database).Collection(UsersCollection.String())
	var users []*cubawheeler.User
	var token string
	f := bson.D{}
	if len(filter.Ids) > 0 {
		f = append(f, primitive.E{Key: "_id", Value: primitive.A{"$in", filter.Ids}})
	}
	if filter.Email != "" {
		f = append(f, primitive.E{Key: "email", Value: filter.Email})
	}
	if len(filter.Otp) > 0 {
		f = append(f, primitive.E{Key: "otp", Value: filter.Otp})
	}
	if len(filter.Pin) > 0 {
		f = append(f, primitive.E{Key: "pin", Value: filter.Email})
	}
	if len(filter.Status) > 0 {
		f = append(f, primitive.E{Key: "status", Value: primitive.A{"$in", filter.Status}})
	}
	if len(filter.Role) > 0 {
		f = append(f, primitive.E{Key: "role", Value: filter.Role})
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

		if len(users) == filter.Limit+1 {
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
