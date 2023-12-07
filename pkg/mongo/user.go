package mongo

import (
	"bytes"
	"context"
	"errors"
	"fmt"

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
}

func NewUserService(db *DB) *UserService {
	return &UserService{
		db:         db,
		collection: db.client.Database(database).Collection(UsersCollection.String()),
	}
}

func (s *UserService) Login(ctx context.Context, input cubawheeler.LoginRequest) (*cubawheeler.User, error) {
	user, err := s.FindByEmail(ctx, input.Email)
	if err != nil && errors.Is(err, e.ErrNotFound) {
		user = &cubawheeler.User{
			ID:     cubawheeler.NewID().String(),
			Email:  input.Email,
			Status: cubawheeler.UserStatusInactive,
			Otp:    []byte("000000"),
		}
		// TODO: generate a new OTP for the user
		err = s.CreateUser(ctx, user)
		if err != nil {
			return nil, err
		}
		user, err = s.FindByID(ctx, user.ID)
		if err != nil {
			return nil, err
		}
		return user, nil
	}
	if user != nil && user.Status == cubawheeler.UserStatusInactive && input.Otp != nil {
		if !bytes.Equal(user.Otp, []byte(*input.Otp)) {
			return nil, e.ErrAccessDenied
		}
		user.Status = cubawheeler.UserStatusOnReview
		user.Otp = nil
		if err := updateUser(ctx, s.db.client, user, bson.D{{Key: "status", Value: user.Status}, {Key: "otp", Value: nil}}); err != nil {
			return nil, err
		}
		return user, nil
	}

	return user, nil
}

func (s *UserService) CreateUser(ctx context.Context, user *cubawheeler.User) error {
	_, err := s.collection.InsertOne(ctx, user)
	if err != nil {
		return fmt.Errorf("unable to store the user: %w", err)
	}
	if _, err := createProfile(ctx, s.db, &cubawheeler.ProfileRequest{}, user); err != nil {
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
	users, _, err := findAllUsers(ctx, s.collection, &cubawheeler.UserFilter{
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

func (s *UserService) FindByEmail(ctx context.Context, email string) (*cubawheeler.User, error) {
	return findUserByEmail(ctx, s.db.client, email)
}

func (s *UserService) FindAll(ctx context.Context, filter *cubawheeler.UserFilter) (*cubawheeler.UserList, error) {
	user := cubawheeler.UserFromContext(ctx)
	if user == nil {
		return nil, errors.New("invalid token provided")
	}
	if user.Role != cubawheeler.RoleAdmin {
		return nil, errors.New("access denied")
	}
	users, token, err := findAllUsers(ctx, s.collection, filter)
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
		location.Lat = input.Location.Lat
		location.Long = input.Location.Long
	}
	usr.Locations = append(usr.Locations, location)
	if err := updateAddFavoritesPlaces(ctx, s.collection, usr); err != nil {
		return nil, fmt.Errorf("unable to store the favorite palces: %w", err)
	}
	return &location, nil
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

func (s *UserService) Trips(ctx context.Context, filter *cubawheeler.TripFilter) (*cubawheeler.TripList, error) {
	usr := cubawheeler.UserFromContext(ctx)
	if usr == nil {
		return nil, e.ErrNotFound
	}

	tripsCollection := s.db.client.Database(database).Collection("trips")
	trips, token, err := findAllTrips(ctx, tripsCollection, filter)
	if err != nil {
		return nil, err
	}
	return &cubawheeler.TripList{Data: trips, Token: token}, nil
}

func (s *UserService) LastNAddress(ctx context.Context, number int) ([]*cubawheeler.Location, error) {
	usr := cubawheeler.UserFromContext(ctx)
	if usr != nil {
		return nil, errors.New("invalid token profived")
	}
	panic("implement me")
}

func (s *UserService) UpdateOTP(ctx context.Context, otp string, u2 uint64) error {
	//TODO implement me
	panic("implement me")
}

func findUserByEmail(ctx context.Context, client *mongo.Client, email string) (*cubawheeler.User, error) {
	collection := client.Database(database).Collection(UsersCollection.String())
	users, _, err := findAllUsers(ctx, collection, &cubawheeler.UserFilter{Email: email, Limit: 1})
	if err != nil {
		return nil, err
	}
	if len(users) == 0 {
		return nil, e.ErrNotFound
	}
	return users[0], nil
}

func findAllUsers(ctx context.Context, collection *mongo.Collection, filter *cubawheeler.UserFilter) ([]*cubawheeler.User, string, error) {
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

func updateAddFavoritesPlaces(ctx context.Context, collection *mongo.Collection, usr *cubawheeler.User) error {
	f := bson.D{{"$set", bson.E{"locations", usr.Locations}}}
	_, err := collection.UpdateOne(ctx, bson.D{{"_id", usr.ID}}, f)
	if err != nil {
		return fmt.Errorf("unable to update the location: %w", err)
	}
	return nil
}

func updateUser(ctx context.Context, client *mongo.Client, user *cubawheeler.User, data bson.D) error {
	collection := client.Database(database).Collection(UsersCollection.String())
	_, err := collection.UpdateOne(ctx, bson.D{{"email", user.Email}}, bson.D{{"$set", data}})
	if err != nil {
		return fmt.Errorf("unable to update the user: %v: %w", err, e.ErrInternal)
	}
	return nil
}
