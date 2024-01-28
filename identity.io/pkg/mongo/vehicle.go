package mongo

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"identity.io/pkg/derrors"
	"identity.io/pkg/identity"
)

var _ identity.VehicleService = &VehicleService{}

var VehiclesCollection Collections = "vehicles"

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

func (s *VehicleService) Store(ctx context.Context, vehicle *identity.Vehicle) (err error) {
	defer derrors.Wrap(&err, "mongo.VehicleService.Store")
	vehicle.ID = identity.NewID().String()
	vehicle.CreatedAt = time.Now().UTC().Unix()
	return storeVehicle(ctx, s.db, vehicle)
}

func (s *VehicleService) Update(ctx context.Context, input identity.UpdateVehicle) (_ *identity.Vehicle, err error) {
	defer derrors.Wrap(&err, "mongo.VehicleService.Update")
	user := identity.UserFromContext(ctx)
	if user == nil {
		return nil, fmt.Errorf("unable to update the vehicle: %w", identity.ErrAccessDenied)
	}
	vehicles, _, err := findAllVehicles(ctx, s.collection, &identity.VehicleFilter{
		Ids:   []string{input.ID},
		Limit: 1,
	})
	if err != nil {
		return nil, err
	}
	if len(vehicles) > 0 {
		return nil, errors.New("vehicle not found")
	}
	vehicle := vehicles[0]
	if vehicle.User != user.ID && user.Role != identity.RoleAdmin {
		return nil, fmt.Errorf("access denied: %v: %w", err, identity.ErrAccessDenied)
	}
	f := bson.D{}
	params := bson.D{}

	if len(input.Plate) > 0 {
		vehicle.Plate = input.Plate
		params = append(params, bson.E{Key: "plate", Value: vehicle.Plate})
	}
	if input.Category.IsValid() {
		vehicle.Category = input.Category
		params = append(params, bson.E{Key: "category", Value: vehicle.Category})
	}
	if input.Type.IsValid() {
		vehicle.Type = input.Type
		params = append(params, bson.E{Key: "type", Value: vehicle.Type})
	}
	if input.Year > 0 {
		vehicle.Year = input.Year
		params = append(params, bson.E{Key: "year", Value: vehicle.Year})
	}
	if len(input.Facilities) > 0 {
		vehicle.Facilities = input.Facilities
		params = append(params, bson.E{Key: "facilities", Value: vehicle.Facilities})
	}
	f = append(f, bson.E{Key: "$set", Value: params})
	_, err = s.collection.UpdateOne(ctx, bson.D{{Key: "_id", Value: input.ID}}, f)
	if err != nil {
		return nil, fmt.Errorf("unable to update the vehicle")
	}
	return vehicle, nil
}

func (s *VehicleService) FindByID(ctx context.Context, id string) (*identity.Vehicle, error) {
	vehicles, _, err := findAllVehicles(ctx, s.collection, &identity.VehicleFilter{
		Ids: []string{id},
	})
	if err != nil {
		return nil, fmt.Errorf("vehicle not found")
	}
	return vehicles[0], nil
}

func (s *VehicleService) FindByPlate(ctx context.Context, plate string) (*identity.Vehicle, error) {
	return findVehicleByPlate(ctx, s.db, plate)
}

func (s *VehicleService) FindAll(ctx context.Context, filter *identity.VehicleFilter) ([]*identity.Vehicle, string, error) {
	return findAllVehicles(ctx, s.collection, filter)
}

func findVehicleByPlate(ctx context.Context, db *DB, plate string) (*identity.Vehicle, error) {
	collection := db.client.Database(database).Collection(VehiclesCollection.String())
	vehicles, _, err := findAllVehicles(ctx, collection, &identity.VehicleFilter{Plate: plate, Limit: 1})
	if err != nil {
		return nil, err
	}
	if len(vehicles) == 0 {
		return nil, identity.ErrNotFound
	}
	return vehicles[0], nil
}

func findAllVehicles(ctx context.Context, collection *mongo.Collection, filter *identity.VehicleFilter) ([]*identity.Vehicle, string, error) {
	var vehicles []*identity.Vehicle
	var token string
	usr := identity.UserFromContext(ctx)
	if usr == nil {
		return nil, "", errors.New("invalid token provided")
	}
	if usr.Role != identity.RoleAdmin {
		filter.User = usr.ID
	}
	f := bson.D{}
	if len(filter.Ids) > 0 {
		f = append(f, primitive.E{Key: "_id", Value: primitive.A{"$in", filter.Ids}})
	}
	if len(filter.Token) > 0 {
		f = append(f, bson.E{Key: "_id", Value: primitive.E{Key: "$gt", Value: filter.Token}})
	}
	if len(filter.Plate) > 0 {
		f = append(f, bson.E{Key: "plate", Value: filter.Plate})
	}
	if len(filter.Color) > 0 {
		f = append(f, bson.E{Key: "color", Value: filter.Color})
	}
	if len(filter.Model) > 0 {
		f = append(f, bson.E{Key: "model", Value: filter.Model})
	}
	if len(filter.User) > 0 {
		f = append(f, bson.E{Key: "user_id", Value: filter.User})
	}

	cur, err := collection.Find(ctx, f)
	if err != nil {
		return nil, "", err
	}
	for cur.Next(ctx) {
		var vehicle identity.Vehicle
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

func storeVehicle(ctx context.Context, db *DB, vehicle *identity.Vehicle) error {
	_, err := db.Collection(VehiclesCollection).InsertOne(ctx, vehicle)
	if err != nil {
		return fmt.Errorf("unable to store the vehicle: %w", err)
	}
	return nil
}

func setupVethicles(t *testing.T, db *DB) []*identity.Vehicle {
	t.Helper()
	// Setup vehicles categories
	var vehicles = []*identity.Vehicle{
		{
			ID:       "1",
			Category: identity.VehicleCategoryX,
			Brand:    "Kia",
			Year:     2019,
			Seats:    4,
			CarModel: "Picanto",
			Colors:   []string{"red", "blue"},
			Status:   identity.VehicleStatusActive,
			Type:     identity.TypeAuto,
		},
	}

	for _, c := range vehicles {
		if err := storeVehicle(context.Background(), db, c); err != nil {
			t.Fatalf("unable to store vehicle: %v", err)
		}
	}
	return vehicles
}
