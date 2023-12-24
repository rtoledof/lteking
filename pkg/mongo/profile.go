package mongo

import (
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"cubawheeler.io/pkg/cubawheeler"
)

var (
	_ cubawheeler.ProfileService = &ProfileService{}
)

const ProfileCollection Collections = "profiles"

type ProfileService struct {
	db         *DB
	collection *mongo.Collection
}

func NewProfileService(db *DB) *ProfileService {
	return &ProfileService{
		db:         db,
		collection: db.client.Database(database).Collection(ProfileCollection.String()),
	}
}

func (s *ProfileService) Create(ctx context.Context, request *cubawheeler.UpdateProfile) (*cubawheeler.Profile, error) {
	usr := cubawheeler.UserFromContext(ctx)
	if usr == nil {
		return nil, errors.New("invalid token provided")
	}
	return createProfile(ctx, s.db, request, usr)
}

func (s *ProfileService) Update(ctx context.Context, request *cubawheeler.UpdateProfile) (*cubawheeler.Profile, error) {
	usr := cubawheeler.UserFromContext(ctx)
	if usr == nil {
		return nil, errors.New("invalid token provided")
	}
	profile, err := s.FindByUser(ctx)
	if err != nil {
		return nil, err
	}
	params := bson.D{}
	if request.Name != nil {
		params = append(params, primitive.E{Key: "name", Value: *request.Name})
		profile.Name = *request.Name
	}
	if request.LastName != nil {
		params = append(params, primitive.E{Key: "last_name", Value: *request.LastName})
		profile.LastName = *request.LastName
	}
	if request.Dob != nil {
		params = append(params, primitive.E{Key: "dob", Value: *request.Dob})
		profile.DOB = *request.Dob
	}
	if request.Phone != nil {
		params = append(params, primitive.E{Key: "phone", Value: *request.Phone})
		profile.Phone = *request.Phone
	}
	if request.Photo != nil {
		params = append(params, primitive.E{Key: "photo", Value: *request.Photo})
		profile.Phone = *request.Phone
	}
	if request.Gender != nil {
		params = append(params, primitive.E{Key: "gender", Value: *request.Gender})
		profile.Gender = *request.Gender
	}
	if request.Licence != nil {
		params = append(params, primitive.E{Key: "licence", Value: *request.Licence})
		profile.Licence = *request.Licence
	}
	if request.Dni != nil {
		params = append(params, primitive.E{Key: "dni", Value: *request.Dni})
		profile.Dni = *request.Dni
	}

	if profile.IsCompleted(usr.Role) {
		profile.Status = cubawheeler.ProfileStatusCompleted
		params = append(params, primitive.E{Key: "status", Value: profile.Status})
	}

	_, err = s.collection.UpdateOne(ctx,
		bson.D{
			{Key: "_id", Value: profile.ID},
		},
		bson.D{
			{Key: "$set", Value: params},
		})
	if err != nil {
		return nil, fmt.Errorf("unable to update user profile: %w", err)
	}
	profile, _ = s.FindByUser(ctx)
	if profile.IsCompleted(usr.Role) {
		usr.Status = cubawheeler.UserStatusActive
		if err := updateUser(ctx, s.db, usr); err != nil {
			return nil, err
		}
	}

	return profile, nil
}

func (s *ProfileService) ChangePin(ctx context.Context, old *string, pin string) (*cubawheeler.Profile, error) {
	usr := cubawheeler.UserFromContext(ctx)
	if usr == nil {
		return nil, errors.New("invalid token provided")
	}
	if err := usr.ComparePin(*old); err != nil || usr.Pin != nil {
		return nil, errors.New("user/pin incorrect")
	}
	if err := usr.EncryptPin(pin); err != nil {
		return nil, err
	}
	return s.Update(ctx, &cubawheeler.UpdateProfile{})
}

func (s *ProfileService) FindAll(ctx context.Context, filter *cubawheeler.ProfileFilter) ([]*cubawheeler.Profile, string, error) {
	return findAllProfiles(ctx, s.collection, filter)
}

func (s *ProfileService) FindByUser(ctx context.Context) (*cubawheeler.Profile, error) {
	usr := cubawheeler.UserFromContext(ctx)
	if usr == nil {
		return nil, errors.New("nil user in context")
	}
	profiles, _, err := findAllProfiles(ctx, s.collection, &cubawheeler.ProfileFilter{
		User:  usr.ID,
		Limit: 1,
	})
	if err != nil {
		return nil, err
	}
	if len(profiles) == 0 {
		return nil, cubawheeler.ErrNotFound
	}
	return profiles[0], nil
}

func findAllProfiles(ctx context.Context, collection *mongo.Collection, filter *cubawheeler.ProfileFilter) ([]*cubawheeler.Profile, string, error) {
	var profiles []*cubawheeler.Profile
	var token string
	f := bson.D{}

	if len(filter.Token) > 0 {
		f = append(f, primitive.E{Key: "_id", Value: primitive.A{"$gt", filter.Token}})
	}
	if len(filter.Dni) > 0 {
		f = append(f, primitive.E{Key: "dni", Value: filter.Dni})
	}
	if len(filter.User) > 0 {
		f = append(f, primitive.E{Key: "user_id", Value: filter.User})
	}
	cur, err := collection.Find(ctx, f)
	if err != nil {
		return nil, "", err
	}
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var profile cubawheeler.Profile
		err := cur.Decode(&profile)
		if err != nil {
			return nil, "", err
		}
		profiles = append(profiles, &profile)
		if len(profiles) == filter.Limit+1 {
			break
		}
	}
	if len(profiles) > filter.Limit {
		token = profiles[filter.Limit].ID
		profiles = profiles[:filter.Limit]
	}
	if err := cur.Err(); err != nil {
		return nil, "", err
	}
	return profiles, token, err
}

func createProfile(ctx context.Context, db *DB, request *cubawheeler.UpdateProfile, usr *cubawheeler.User) (*cubawheeler.Profile, error) {
	user := cubawheeler.UserFromContext(ctx)
	if user == nil {
		return nil, fmt.Errorf("nil user in context")
	}
	profile := &cubawheeler.Profile{
		ID:     cubawheeler.NewID().String(),
		UserId: usr.ID,
		Status: cubawheeler.ProfileStatusIncompleted,
	}
	user.Profile = cubawheeler.Profile{}
	collection := db.client.Database(database).Collection(string(UsersCollection))
	_, err := collection.InsertOne(ctx, profile)
	if err != nil {
		return nil, fmt.Errorf("unable to crete the profile: %w", err)
	}
	return profile, nil
}
