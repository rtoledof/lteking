package mongo

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"cubawheeler.io/pkg/cubawheeler"
)

var _ cubawheeler.TripService = &TripService{}

type TripService struct {
	db         *DB
	collection *mongo.Collection
	tripChan   chan *cubawheeler.Trip
}

func NewTripService(db *DB) *TripService {
	return &TripService{
		db:         db,
		tripChan:   make(chan *cubawheeler.Trip, 10000),
		collection: db.client.Database(database).Collection("trips"),
	}
}

func (s *TripService) Create(ctx context.Context, trip *cubawheeler.Trip) error {
	user := cubawheeler.UserFromContext(ctx)
	if user == nil {
		return errors.New("invalid token provided")
	}
	trip.ID = cubawheeler.NewID().String()
	trip.CreatedAt = time.Now().UTC().UnixNano()
	_, err := s.collection.InsertOne(ctx, trip)
	if err != nil {
		return fmt.Errorf("unable to store the trip: %w", err)
	}
	return nil
}

func (s *TripService) Update(ctx context.Context, trip *cubawheeler.UpdateTrip) (*cubawheeler.Trip, error) {
	//TODO implement me
	panic("implement me")
}

func (s *TripService) FindByID(ctx context.Context, id string) (*cubawheeler.Trip, error) {
	limit := 1
	trips, _, err := findAllTrips(ctx, s.collection, &cubawheeler.TripFilter{
		Ids:   []*string{&id},
		Limit: &limit,
	})
	if err != nil && len(trips) == 0 {
		return nil, errors.New("trip not found")
	}
	return trips[0], nil
}

func (s *TripService) FindAll(ctx context.Context, filter *cubawheeler.TripFilter) (*cubawheeler.TripList, error) {
	trips, token, err := findAllTrips(ctx, s.collection, filter)
	if err != nil {
		return nil, err
	}
	return &cubawheeler.TripList{Data: trips, Token: token}, nil
}

func findAllTrips(ctx context.Context, collection *mongo.Collection, filter *cubawheeler.TripFilter) ([]*cubawheeler.Trip, string, error) {
	user := cubawheeler.UserFromContext(ctx)
	if user == nil {
		return nil, "", errors.New("invalid token provided")
	}
	switch user.Role {
	case cubawheeler.RoleRider:
		filter.Rider = &user.ID
	case cubawheeler.RoleDriver:
		filter.Driver = &user.ID
	}
	var trips []*cubawheeler.Trip
	var token string
	f := bson.D{}
	if filter.Rider != nil {
		f = append(f, bson.E{Key: "rider", Value: filter.Rider})
	}
	if filter.Driver != nil {
		f = append(f, bson.E{Key: "driver", Value: filter.Driver})
	}
	if filter.Token != nil {
		f = append(f, bson.E{Key: "_id", Value: primitive.E{"$gt", filter.Token}})
	}

	cur, err := collection.Find(ctx, f)
	if err != nil {
		return nil, "", err
	}
	for cur.Next(ctx) {
		var trip cubawheeler.Trip
		err := cur.Decode(&trip)
		if err != nil {
			return nil, "", err
		}
		trips = append(trips, &trip)
		if len(trips) == *filter.Limit+1 {
			token = trips[*filter.Limit].ID
			trips = trips[:*filter.Limit]
			break
		}
	}

	if err := cur.Err(); err != nil {
		return nil, "", err
	}
	cur.Close(ctx)
	return trips, token, nil
}
