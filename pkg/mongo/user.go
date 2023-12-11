package mongo

import (
	"context"
	"cubawheeler.io/pkg/pusher"
	"cubawheeler.io/pkg/redis"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"cubawheeler.io/pkg/cubawheeler"
	e "cubawheeler.io/pkg/errors"
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
	s := &UserService{
		db:         db,
		collection: db.client.Database(database).Collection(UsersCollection.String()),
		beansToken: beansToken,
		beans:      beans,
	}
	ticker := time.NewTicker(time.Hour * 1)
	go func() {

		for {
			select {
			case <-ticker.C:
				if err := s.generateTokens(context.Background()); err != nil {
					slog.Info(err.Error())
				}
			case <-done:
				return
			}
		}

	}()
	return s
}

func (s *UserService) Login(ctx context.Context, input cubawheeler.LoginRequest) (*cubawheeler.User, error) {
	app := cubawheeler.ClientFromContext(ctx)
	user, err := s.FindByEmail(ctx, input.Email)
	if err != nil && errors.Is(err, e.ErrNotFound) {
		if app == nil {
			return nil, fmt.Errorf("no application provided: %w", e.ErrAccessDenied)
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
		// TODO: generate a new OTP for the user

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

	if !user.IsActive() {
		return nil, e.ErrAccessDenied
	}

	return user, nil
}

func (s *UserService) CreateUser(ctx context.Context, user *cubawheeler.User) error {
	_, err := s.collection.InsertOne(ctx, user)
	if err != nil {
		return fmt.Errorf("unable to store the user: %w", err)
	}
	if _, err := createProfile(ctx, s.db, &cubawheeler.UpdateProfile{}, user); err != nil {
		return err
	}
	return nil
}

func (s *UserService) FindByID(ctx context.Context, id string) (*cubawheeler.User, error) {
	usr := cubawheeler.UserFromContext(ctx)
	if usr == nil {
		return nil, errors.New("invalid token provided")
	}
	if usr.Role != cubawheeler.RoleAdmin {
		return usr, nil
	}
	return findUserByID(ctx, s.db, id)
}

func (s *UserService) FindByEmail(ctx context.Context, email string) (*cubawheeler.User, error) {
	return findUserByEmail(ctx, s.db, email)
}

func (s *UserService) FindAll(ctx context.Context, filter *cubawheeler.UserFilter) (*cubawheeler.UserList, error) {
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

func (s *UserService) AddFavoritePlace(ctx context.Context, input cubawheeler.AddPlace) (*cubawheeler.Location, error) {
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

func (s *UserService) FavoritePlaces(ctx context.Context) ([]*cubawheeler.Location, error) {
	usr := cubawheeler.UserFromContext(ctx)
	if usr == nil {
		return nil, fmt.Errorf("nil user in context: %w", e.ErrAccessDenied)
	}
	return usr.Locations, nil
}

func (s *UserService) AddFavoriteVehicle(ctx context.Context, plate *string) (*cubawheeler.Vehicle, error) {
	usr := cubawheeler.UserFromContext(ctx)
	if usr == nil {
		return nil, fmt.Errorf("nil user in context: %w", e.ErrAccessDenied)
	}
	vehicle, err := findVehicleByPlate(ctx, s.db, *plate)
	if err != nil {
		return nil, err
	}
	usr.FavoriteVehicles = append(usr.FavoriteVehicles, vehicle)
	if err := updateUser(ctx, s.db, usr, bson.D{{Key: "favorite_vehicles", Value: usr.FavoriteVehicles}}); err != nil {
		return nil, err
	}
	return vehicle, nil
}

func (s *UserService) UpdatePlace(ctx context.Context, input *cubawheeler.UpdatePlace) (*cubawheeler.Location, error) {
	user := cubawheeler.UserFromContext(ctx)
	if user == nil {
		return nil, fmt.Errorf("unable to update the place: %w", e.ErrAccessDenied)
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

func (s *UserService) FavoriteVehicles(ctx context.Context) ([]*cubawheeler.Vehicle, error) {
	usr := cubawheeler.UserFromContext(ctx)
	if usr == nil {
		return nil, fmt.Errorf("nil user in context: %w", e.ErrAccessDenied)
	}
	return usr.FavoriteVehicles, nil
}

func (s *UserService) Me(ctx context.Context) (*cubawheeler.Profile, error) {
	usr := cubawheeler.UserFromContext(ctx)
	if usr == nil {
		return nil, errors.New("invalid token provided")
	}

	profiles, _, err := findAllProfiles(ctx, s.db.client.Database(database).Collection(ProfileCollection.String()), &cubawheeler.ProfileFilter{
		User:  usr.ID,
		Limit: 1,
	})
	if err != nil {
		return nil, err
	}
	return profiles[0], nil
}

func (s *UserService) Orders(ctx context.Context, filter *cubawheeler.OrderFilter) (*cubawheeler.OrderList, error) {
	usr := cubawheeler.UserFromContext(ctx)
	if usr == nil {
		return nil, e.ErrNotFound
	}

	ordersCollection := s.db.client.Database(database).Collection(OrderCollection.String())
	orders, token, err := findOrders(ctx, ordersCollection, filter)
	if err != nil {
		return nil, err
	}
	return &cubawheeler.OrderList{Data: orders, Token: token}, nil
}

func (s *UserService) LastNAddress(ctx context.Context, number int) ([]*cubawheeler.Location, error) {
	usr := cubawheeler.UserFromContext(ctx)
	if usr != nil {
		return nil, errors.New("invalid token profived")
	}
	panic("implement me")
}

func (s *UserService) generateTokens(ctx context.Context) error {
	users, _, err := findAllUsers(ctx, s.db, &cubawheeler.UserFilter{
		Status: []cubawheeler.UserStatus{cubawheeler.UserStatusActive},
	})
	if err != nil {
		return err
	}
	for _, v := range users {
		if _, err := s.beansToken.GetBeansToken(ctx, v.ID); err != nil && errors.Is(e.ErrNotFound, err) {
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
		return nil, e.ErrNotFound
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
	f := bson.D{{"$set", bson.E{"locations", usr.Locations}}}
	_, err := collection.UpdateOne(ctx, bson.D{{"_id", usr.ID}}, f)
	if err != nil {
		return fmt.Errorf("unable to update the location: %w", err)
	}
	return nil
}

func updateUser(ctx context.Context, db *DB, user *cubawheeler.User, data bson.D) error {
	collection := db.client.Database(database).Collection(UsersCollection.String())
	_, err := collection.UpdateOne(ctx, bson.D{{"email", user.Email}}, bson.D{{"$set", data}})
	if err != nil {
		return fmt.Errorf("unable to update the user: %v: %w", err, e.ErrInternal)
	}
	return nil
}
