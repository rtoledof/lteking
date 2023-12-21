package mongo

import (
	"context"
	"fmt"

	"cubawheeler.io/pkg/cubawheeler"
	"go.mongodb.org/mongo-driver/bson"
)

var _ cubawheeler.VehicleCategoryRateService = &VehicleCategoryRateService{}

const VehicleCategoryRateCollection Collections = "vehicle_category_rates"

type VehicleCategoryRateService struct {
	db *DB
}

func NewVehicleCategoryRateService(db *DB) *VehicleCategoryRateService {
	return &VehicleCategoryRateService{
		db: db,
	}
}

// Create implements cubawheeler.VehicleCategoryRateService.
func (s *VehicleCategoryRateService) Create(ctx context.Context, req *cubawheeler.VehicleCategoryRateRequest) (*cubawheeler.VehicleCategoryRate, error) {
	usr := cubawheeler.UserFromContext(ctx)
	if usr == nil || usr.Role != cubawheeler.RoleAdmin {
		return nil, cubawheeler.ErrAccessDenied
	}
	var rate cubawheeler.VehicleCategoryRate
	assembleVehicleCategoryRate(&rate, *req)
	if err := rate.Validate(); err != nil {
		return nil, err
	}
	if err := insertVehicleCategoryRate(ctx, s.db, &rate); err != nil {
		return nil, err
	}
	return &rate, nil
}

func assembleVehicleCategoryRate(rate *cubawheeler.VehicleCategoryRate, req cubawheeler.VehicleCategoryRateRequest) {
	rate.ID = req.ID
	if req.ID == "" {
		rate.ID = cubawheeler.NewID().String()
	}
	rate.Category = req.Category
	rate.Factor = req.Factor
}

// Update implements cubawheeler.VehicleCategoryRateService.
func (s *VehicleCategoryRateService) Update(ctx context.Context, req *cubawheeler.VehicleCategoryRateRequest) (*cubawheeler.VehicleCategoryRate, error) {
	usr := cubawheeler.UserFromContext(ctx)
	if usr == nil || usr.Role != cubawheeler.RoleAdmin {
		return nil, cubawheeler.ErrAccessDenied
	}
	rate, err := findVehicleCategoryRateByID(ctx, s.db, req.ID)
	if err != nil {
		return nil, err
	}
	assembleVehicleCategoryRate(rate, *req)
	if err := rate.Validate(); err != nil {
		return nil, err
	}
	if err := updateVehicleCategoryRate(ctx, s.db, rate.ID, rate); err != nil {
		return nil, err
	}
	return rate, nil
}

// FindByCategory implements cubawheeler.VehicleCategoryRateService.
func (s *VehicleCategoryRateService) FindByCategory(ctx context.Context, category cubawheeler.VehicleCategory) (*cubawheeler.VehicleCategoryRate, error) {
	usr := cubawheeler.UserFromContext(ctx)
	if usr == nil || usr.Role != cubawheeler.RoleAdmin {
		return nil, cubawheeler.ErrAccessDenied
	}
	return findVehicleCategoryRateByCategory(ctx, s.db, category)
}

// FindByID implements cubawheeler.VehicleCategoryRateService.
func (s *VehicleCategoryRateService) FindByID(ctx context.Context, id string) (*cubawheeler.VehicleCategoryRate, error) {
	usr := cubawheeler.UserFromContext(ctx)
	if usr == nil || usr.Role != cubawheeler.RoleAdmin {
		return nil, cubawheeler.ErrAccessDenied
	}

	return findVehicleCategoryRateByID(ctx, s.db, id)
}

// FindAll implements cubawheeler.VehicleCategoryRateService.
func (s *VehicleCategoryRateService) FindAll(ctx context.Context, filter cubawheeler.VehicleCategoryRateFilter) ([]*cubawheeler.VehicleCategoryRate, string, error) {
	usr := cubawheeler.UserFromContext(ctx)
	if usr == nil || usr.Role != cubawheeler.RoleAdmin {
		return nil, "", cubawheeler.ErrAccessDenied
	}
	return findVehicleCategoriesRate(ctx, s.db, filter)
}

func findVehicleCategoriesRate(ctx context.Context, db *DB, filter cubawheeler.VehicleCategoryRateFilter) ([]*cubawheeler.VehicleCategoryRate, string, error) {
	var rates []*cubawheeler.VehicleCategoryRate
	var token string
	f := bson.D{}
	if len(filter.Ids) > 0 {
		f = append(f, bson.E{Key: "_id", Value: bson.D{{Key: "$in", Value: filter.Ids}}})
	}
	if len(filter.Category) > 0 {
		f = append(f, bson.E{Key: "category", Value: bson.D{{Key: "$in", Value: filter.Category}}})
	}
	if filter.Token != "" {
		f = append(f, bson.E{Key: "_id", Value: bson.D{{Key: "$gt", Value: filter.Token}}})
	}
	collection := db.Collection(VehicleCategoryRateCollection)
	cur, err := collection.Find(ctx, f)
	if err != nil {
		return nil, "", fmt.Errorf("find error: %v: %w", err, cubawheeler.ErrInternal)
	}
	defer cur.Close(ctx)

	for cur.Next(ctx) {
		var rate cubawheeler.VehicleCategoryRate
		err := cur.Decode(&rate)
		if err != nil {
			return nil, "", err
		}
		rates = append(rates, &rate)
		if len(rates) == filter.Limit+1 && filter.Limit > 0 {
			token = rates[filter.Limit].ID
			rates = rates[:filter.Limit]
			break
		}
	}
	if err := cur.Err(); err != nil {
		return nil, "", fmt.Errorf("cursor error: %v: %w", err, cubawheeler.ErrInternal)
	}

	return rates, token, nil
}

func findVehicleCategoryRateByID(ctx context.Context, db *DB, id string) (*cubawheeler.VehicleCategoryRate, error) {
	rates, _, err := findVehicleCategoriesRate(ctx, db, cubawheeler.VehicleCategoryRateFilter{Ids: []string{id}, Limit: 1})
	if err != nil {
		return nil, err
	}
	if len(rates) == 0 {
		return nil, fmt.Errorf("unable to find rate: %v: %w", id, cubawheeler.ErrNotFound)
	}
	return rates[0], nil
}

func findVehicleCategoryRateByCategory(ctx context.Context, db *DB, category cubawheeler.VehicleCategory) (*cubawheeler.VehicleCategoryRate, error) {
	rates, _, err := findVehicleCategoriesRate(ctx, db, cubawheeler.VehicleCategoryRateFilter{Category: []cubawheeler.VehicleCategory{category}, Limit: 1})
	if err != nil {
		return nil, err
	}
	if len(rates) == 0 {
		return nil, fmt.Errorf("unable to find rate: %v: %w", category, cubawheeler.ErrNotFound)
	}
	return rates[0], nil
}

func insertVehicleCategoryRate(ctx context.Context, db *DB, rate *cubawheeler.VehicleCategoryRate) error {
	collection := db.Collection(VehicleCategoryRateCollection)
	if _, err := collection.InsertOne(ctx, rate); err != nil {
		return fmt.Errorf("unable to insert rate: %v: %w", err, cubawheeler.ErrInternal)
	}
	return nil
}

func updateVehicleCategoryRate(ctx context.Context, db *DB, id string, rate *cubawheeler.VehicleCategoryRate) error {
	collection := db.Collection(VehicleCategoryRateCollection)
	if _, err := collection.UpdateOne(ctx, bson.D{{Key: "_id", Value: id}}, bson.D{{Key: "$set", Value: rate}}); err != nil {
		return fmt.Errorf("unable to update rate: %v: %w", err, cubawheeler.ErrInternal)
	}
	return nil
}
