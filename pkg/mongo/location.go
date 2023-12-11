package mongo

import (
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"cubawheeler.io/pkg/cubawheeler"
)

var _ cubawheeler.LocationService = &LocationService{}
var _ cubawheeler.LastLocations = &LocationService{}

const LocationsCollections Collections = "locations"

type LocationService struct {
	db         *DB
	collection *mongo.Collection
}

func NewLocationService(db *DB) *LocationService {
	return &LocationService{
		db:         db,
		collection: db.client.Database(database).Collection(LocationsCollections.String()),
	}
}

func (s *LocationService) Create(ctx context.Context, request *cubawheeler.LocationRequest) (*cubawheeler.Location, error) {
	usr := cubawheeler.UserFromContext(ctx)
	if usr == nil {
		return nil, errors.New("invalid token provided")
	}
	if usr.Role != cubawheeler.RoleAdmin {
		request.User = &usr.ID
	}
	location := &cubawheeler.Location{
		ID:   cubawheeler.NewID().String(),
		Name: request.Name,
		User: *request.User,
		Geolocation: cubawheeler.GeoLocation{
			Type:        "Point",
			Coordinates: []float64{request.Long, request.Lat},
		},
	}
	_, err := s.collection.InsertOne(ctx, location)
	if err != nil {
		return nil, fmt.Errorf("unable to store the location: %w", err)
	}
	return location, nil

}

func (s *LocationService) Update(ctx context.Context, request *cubawheeler.LocationRequest) (*cubawheeler.Location, error) {
	//TODO implement me
	panic("implement me")
}

func (s *LocationService) FindByID(ctx context.Context, id string) (*cubawheeler.Location, error) {
	locations, _, err := findAllLocations(ctx, s.collection, &cubawheeler.LocationRequest{
		Ids: []string{id},
	})
	if err != nil {
		return nil, errors.New("location not found")
	}
	return locations[0], nil
}

func (s *LocationService) FindAll(ctx context.Context, request *cubawheeler.LocationRequest) ([]*cubawheeler.Location, string, error) {
	return findAllLocations(ctx, s.collection, request)
}

func (s *LocationService) Locations(ctx context.Context, n int) ([]*cubawheeler.Location, error) {
	//TODO implement me
	panic("implement me")
}

func findAllLocations(ctx context.Context, collection *mongo.Collection, filter *cubawheeler.LocationRequest) ([]*cubawheeler.Location, string, error) {
	var locations []*cubawheeler.Location
	var token string
	f := bson.D{}
	// TODO: add missing filters here
	cur, err := collection.Find(ctx, f)
	if err != nil {
		return nil, "", err
	}
	for cur.Next(ctx) {
		var location cubawheeler.Location
		err := cur.Decode(&location)
		if err != nil {
			return nil, "", err
		}
		locations = append(locations, &location)
		if len(locations) > filter.Limit+1 {
			token = locations[filter.Limit].ID
			locations = locations[:filter.Limit]
			break
		}
	}
	if err := cur.Err(); err != nil {
		return nil, "", err
	}
	cur.Close(ctx)
	return locations, token, nil
}

//
//func findLastNLocations(ctx context.Context, client *mongo.Client, n int) ([]*cubawheeler.Location, error) {
//	collection := client.Database(database).Collection(LastLocationsCollections.String())
//
//}
