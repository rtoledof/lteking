package mongo

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/ably/ably-go/ably"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"cubawheeler.io/pkg/cubawheeler"
)

//db.restaurants.find({ location:
//	{ $geoWithin:
//		{ $centerSphere: [ [ -73.93414657, 40.82302903 ], 5 / 3963.2 ] } } })

const DriversLocationCollection Collections = "drivers_location"

var DriverLocations = make(chan cubawheeler.Location, 10000)

type RealTimeService struct {
	db        *DB
	exit      chan struct{}
	orderChan chan *cubawheeler.Order
	rest      *ably.REST
}

func NewRealTimeService(
	db *DB,
	exit chan struct{},
	orderChan chan *cubawheeler.Order,
	rest *ably.REST,
) *RealTimeService {
	index := mongo.IndexModel{
		Keys: bson.D{{Key: "location.geo", Value: "2dsphere"}},
	}
	_, err := db.client.Database(database).Collection(DriversLocationCollection.String()).Indexes().CreateOne(context.Background(), index)
	if err != nil {
		panic("unable to create geo index")
	}

	service := &RealTimeService{
		db:        db,
		exit:      exit,
		orderChan: orderChan,
		rest:      rest,
	}

	return service
}

func (s *RealTimeService) FindNearByDrivers(ctx context.Context, location cubawheeler.GeoLocation) ([]*cubawheeler.Location, error) {
	collection := s.db.client.Database(database).Collection(DriversLocationCollection.String())

	var locations []*cubawheeler.Location

	mongoDBHQ := bson.D{
		{Key: "type", Value: "Point"},
		{Key: "coordinates", Value: []float64{location.Long, location.Lat}},
	}
	filter := bson.D{
		{Key: "location.geo",
			Value: bson.D{
				{Key: "$nearSphere", Value: bson.D{
					{Key: "$geometry", Value: mongoDBHQ},
					{Key: "$maxDistance", Value: 5000},
				}},
			},
		},
	}
	cur, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("unable to find the data: %v: %w", err, cubawheeler.ErrNotFound)
	}
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var location cubawheeler.Location
		if err := cur.Decode(&location); err != nil {
			return nil, fmt.Errorf("unable to decode the location: %v: %w", err, cubawheeler.ErrInternal)
		}
		locations = append(locations, &location)
	}
	return locations, nil
}

func (s *RealTimeService) UpdateLocation(context.Context, string, cubawheeler.GeoLocation) error {
	collection := s.db.client.Database(database).Collection(DriversLocationCollection.String())
	ctx := context.Background()
	for v := range DriverLocations {
		v.ID = cubawheeler.NewID().String()
		v.UpdatedAt = time.Now().Unix()
		result := collection.FindOneAndUpdate(ctx, bson.D{{Key: "user_id", Value: v.User}}, v)
		if result == nil {
			v.CreatedAt = time.Now().Unix()
			if _, err := collection.InsertOne(ctx, v); err != nil {
				slog.Info("unable to insert location data: %v: %w", err, cubawheeler.ErrInternal)
			}
		}
	}
	return nil
}
