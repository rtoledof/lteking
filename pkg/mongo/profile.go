package mongo

import (
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"cubawheeler.io/pkg/cubawheeler"
	e "cubawheeler.io/pkg/errors"
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

func (s *ProfileService) Create(ctx context.Context, request *cubawheeler.ProfileRequest) (*cubawheeler.Profile, error) {
	usr := cubawheeler.UserFromContext(ctx)
	if usr == nil {
		return nil, errors.New("invalid token provided")
	}
	return createProfile(ctx, s.db, request, usr)
}

func (s *ProfileService) Update(ctx context.Context, request *cubawheeler.ProfileRequest) (*cubawheeler.Profile, error) {
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
	}
	if request.LastName != nil {
		params = append(params, primitive.E{Key: "last_name", Value: *request.LastName})
	}
	if request.DOB != nil {
		params = append(params, primitive.E{Key: "dob", Value: *request.DOB})
	}
	if request.Phone != nil {
		params = append(params, primitive.E{Key: "phone", Value: *request.Phone})
	}
	if request.Photo != nil {
		params = append(params, primitive.E{Key: "photo", Value: *request.Photo})
	}
	if request.Gender != nil {
		params = append(params, primitive.E{Key: "gender", Value: *request.Gender})
	}
	if request.Licence != nil {
		params = append(params, primitive.E{Key: "licence", Value: *request.Licence})
	}
	if request.Dni != nil {
		params = append(params, primitive.E{Key: "dni", Value: *request.Dni})
	}
	if request.Pin != nil {
		params = append(params, primitive.E{Key: "dni", Value: *request.Dni})
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
	return s.Update(ctx, &cubawheeler.ProfileRequest{})
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
		return nil, e.ErrNotFound
	}
	return profiles[0], nil
}

func addLastLocations(ctx context.Context, db *mongo.Client, location cubawheeler.Location) error {
	collection := db.Database(database).Collection(ProfileCollection.String())
	if collection == nil {
		return errors.New("unable to retrieve the collectoion")
	}
	panic("implement me")
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

func createProfile(ctx context.Context, db *DB, request *cubawheeler.ProfileRequest, usr *cubawheeler.User) (*cubawheeler.Profile, error) {
	profile := &cubawheeler.Profile{
		ID:     cubawheeler.NewID().String(),
		UserId: usr.ID,
	}
	collection := db.client.Database(database).Collection(string(ProfileCollection))
	_, err := collection.InsertOne(ctx, profile)
	if err != nil {
		return nil, fmt.Errorf("unable to crete the profile: %w", err)
	}
	return profile, nil
}
