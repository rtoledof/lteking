package mongo

import (
	"context"
	"errors"
	"fmt"
	"time"

	"cubawheeler.io/pkg/cubawheeler"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var _ cubawheeler.VehicleService = &VehicleService{}

type VehicleService struct {
	db         *DB
	collection *mongo.Collection
}

func NewVehicleService(db *DB) *VehicleService {
	return &VehicleService{
		db:         db,
		collection: db.client.Database(database).Collection("vehicles"),
	}
}

func (s *VehicleService) Store(ctx context.Context, vehicle *cubawheeler.Vehicle) error {
	vehicle.ID = cubawheeler.NewID().String()
	vehicle.CreatedAt = time.Now().UTC().Unix()
	_, err := s.collection.InsertOne(ctx, vehicle)
	if err != nil {
		return fmt.Errorf("unable to store the vehicle: %w", err)
	}
	return nil
}

func (s *VehicleService) Update(ctx context.Context, input cubawheeler.UpdateVehicle) error {
	vehicles, _, err := findAllVehicles(ctx, s.collection, &cubawheeler.VehicleFilter{
		Ids:   []string{input.ID},
		Limit: 1,
	})
	if err != nil {
		return err
	}
	if len(vehicles) > 0 {
		return errors.New("vehicle not found")
	}
	f := bson.D{}
	params := bson.D{}
	vehicle := vehicles[0]
	if len(input.Plate) > 0 {
		vehicle.Plate = &input.Plate
		params = append(params, bson.E{"plate", vehicle.Plate})
	}
	if input.Category.IsValid() {
		vehicle.Category = input.Category
		params = append(params, bson.E{"category", vehicle.Category})
	}
	if input.Type.IsValid() {
		vehicle.Type = input.Type
		params = append(params, bson.E{"type", vehicle.Type})
	}
	if input.Year > 0 {
		vehicle.Year = input.Year
		params = append(params, bson.E{"year", vehicle.Year})
	}
	if len(input.Facilities) > 0 {
		vehicle.Facilities = input.Facilities
		params = append(params, bson.E{"facilities", vehicle.Facilities})
	}
	f = append(bson.D{{"$set", params}})
	_, err = s.collection.UpdateOne(ctx, bson.D{{"_id", input.ID}}, f)
	if err != nil {
		return fmt.Errorf("unable to update the vehicle")
	}
	return nil
}

func (s *VehicleService) FindByID(ctx context.Context, id string) (*cubawheeler.Vehicle, error) {
	vehicles, _, err := findAllVehicles(ctx, s.collection, &cubawheeler.VehicleFilter{
		Ids: []string{id},
	})
	if err != nil {
		return nil, fmt.Errorf("vehicle not found")
	}
	return vehicles[0], nil
}

func (s *VehicleService) FindAll(ctx context.Context, filter *cubawheeler.VehicleFilter) ([]*cubawheeler.Vehicle, string, error) {
	return findAllVehicles(ctx, s.collection, filter)
}

func findAllVehicles(ctx context.Context, collection *mongo.Collection, filter *cubawheeler.VehicleFilter) ([]*cubawheeler.Vehicle, string, error) {
	var vehicles []*cubawheeler.Vehicle
	var token string
	usr := cubawheeler.UserFromContext(ctx)
	if usr == nil {
		return nil, "", errors.New("invalid token provided")
	}
	if usr.Role != cubawheeler.RoleAdmin {
		filter.User = usr.ID
	}
	f := bson.D{}
	if len(filter.Ids) > 0 {
		f = append(f, primitive.E{Key: "_id", Value: primitive.A{"$in", filter.Ids}})
	}
	if len(filter.Token) > 0 {
		f = append(f, bson.E{"_id", primitive.E{"$gt", filter.Token}})
	}
	if len(filter.Plate) > 0 {
		f = append(f, bson.E{"plate", filter.Plate})
	}
	if len(filter.Color) > 0 {
		f = append(f, bson.E{"color", filter.Color})
	}
	if len(filter.Model) > 0 {
		f = append(f, bson.E{"model", filter.Model})
	}
	if len(filter.User) > 0 {
		f = append(f, bson.E{"user_id", filter.User})
	}

	cur, err := collection.Find(ctx, f)
	if err != nil {
		return nil, "", err
	}
	for cur.Next(ctx) {
		var vehicle cubawheeler.Vehicle
		err := cur.Decode(&vehicle)
		if err != nil {
			return nil, "", err
		}
		vehicles = append(vehicles, &vehicle)
		if len(vehicles) == filter.Limit+1 {
			token = vehicles[filter.Limit].ID
			vehicles = vehicles[:filter.Limit]
			break
		}
	}
	if err := cur.Err(); err != nil {
		return nil, "", err
	}
	cur.Close(ctx)
	return vehicles, token, nil
}
