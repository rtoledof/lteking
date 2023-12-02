package mongo

import (
	"context"
	"cubawheeler.io/pkg/cubawheeler"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var _ cubawheeler.UserService = &UserService{}

type UserService struct {
	db         *DB
	collection *mongo.Collection
}

func NewUserService(db *DB) *UserService {
	return &UserService{
		db:         db,
		collection: db.client.Database("cubawheeler").Collection("users"),
	}
}

func (s *UserService) CreateUser(ctx context.Context, user *cubawheeler.User) error {
	_, err := s.collection.InsertOne(ctx, user)
	if err != nil {
		return fmt.Errorf("unable to store the user: %w", err)
	}
	return nil
}

func (s *UserService) FindByID(ctx context.Context, id string) (*cubawheeler.User, error) {
	users, _, err := findAll(ctx, s.collection, &cubawheeler.UserFilter{
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
	users, _, err := findAll(ctx, s.collection, &cubawheeler.UserFilter{
		Email: email,
		Limit: 1,
	})
	if err != nil || len(users) == 0 {
		return nil, errors.New("user not found")
	}
	return users[0], nil
}

func (s *UserService) FindAll(ctx context.Context, filter *cubawheeler.UserFilter) ([]*cubawheeler.User, string, error) {
	return findAll(ctx, s.collection, filter)
}

func (s *UserService) UpdateOTP(ctx context.Context, otp string, u2 uint64) error {
	//TODO implement me
	panic("implement me")
}

func findAll(ctx context.Context, collection *mongo.Collection, filter *cubawheeler.UserFilter) ([]*cubawheeler.User, string, error) {
	var users []*cubawheeler.User
	var token string
	f := bson.D{}
	if len(filter.Ids) > 0 {
		ids := make([]cubawheeler.ID, len(filter.Ids))
		for i, v := range filter.Ids {
			var id cubawheeler.ID
			if err := id.UnmarshalText([]byte(v)); err != nil {
				return nil, "", err
			}
			ids[i] = id
		}

		f = append(f, primitive.E{"_id", primitive.A{"$in", filter.Ids}})
	}
	if filter.Email != "" {
		f = append(f, primitive.E{"email", filter.Email})
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
		users = append(users, &user)
		if len(users) == filter.Limit+1 {
			break
		}
	}
	if len(users) > filter.Limit+1 {
		token = users[filter.Limit].ID.String()
		users = users[:filter.Limit]
	}
	if err := cur.Err(); err != nil {
		return nil, "", err
	}
	cur.Close(ctx)
	return users, token, nil
}
