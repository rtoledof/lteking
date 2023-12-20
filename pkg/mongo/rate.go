package mongo

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"

	"cubawheeler.io/pkg/cubawheeler"
)

var _ cubawheeler.RateService = &RateService{}

var RatesCollection Collections = "rates"

type RateService struct {
	db *DB
}

func NewRateService(db *DB) *RateService {
	return &RateService{
		db: db,
	}
}

func (s *RateService) Create(ctx context.Context, request cubawheeler.RateRequest) (*cubawheeler.Rate, error) {
	usr := cubawheeler.UserFromContext(ctx)
	if usr == nil || usr.Role != cubawheeler.RoleAdmin {
		return nil, cubawheeler.ErrAccessDenied
	}

	var rate cubawheeler.Rate
	assembleRate(&rate, request)
	if err := rate.Validate(); err != nil {
		return nil, err
	}
	if err := insertRate(ctx, s.db, &rate); err != nil {
		return nil, err
	}

	return &rate, nil
}

func (s *RateService) Update(ctx context.Context, request *cubawheeler.RateRequest) (*cubawheeler.Rate, error) {
	usr := cubawheeler.UserFromContext(ctx)
	if usr == nil || usr.Role != cubawheeler.RoleAdmin {
		return nil, cubawheeler.ErrAccessDenied
	}

	rate, err := findRateByID(ctx, s.db, request.ID)
	if err != nil {
		return nil, err
	}
	assembleRate(rate, *request)
	if err := rate.Validate(); err != nil {
		return nil, err
	}
	if err := updateRate(ctx, s.db, rate); err != nil {
		return nil, err
	}
	return rate, nil
}

func (s *RateService) FindByID(ctx context.Context, id string) (*cubawheeler.Rate, error) {
	return findRateByID(ctx, s.db, id)
}

func (s *RateService) FindAll(ctx context.Context, request cubawheeler.RateFilter) ([]*cubawheeler.Rate, string, error) {
	return findRates(ctx, s.db, &request)
}

func updateRate(ctx context.Context, db *DB, rate *cubawheeler.Rate) error {
	collection := db.client.Database(database).Collection(RatesCollection.String())
	if _, err := collection.UpdateOne(ctx, bson.D{{Key: "_id", Value: rate.ID}}, bson.D{{Key: "$set", Value: rate}}); err != nil {
		return fmt.Errorf("unable to update rate: %v: %w", err, cubawheeler.ErrInternal)
	}
	return nil
}

func findRates(ctx context.Context, db *DB, filter *cubawheeler.RateFilter) ([]*cubawheeler.Rate, string, error) {
	collection := db.client.Database(database).Collection(RatesCollection.String())
	var rates []*cubawheeler.Rate
	var token string
	f := bson.D{}
	if len(filter.Ids) > 0 {
		f = append(f, bson.E{Key: "_id", Value: bson.D{{Key: "$in", Value: filter.Ids}}})
	}
	if len(filter.Code) > 0 {
		f = append(f, bson.E{Key: "code", Value: bson.D{{Key: "$in", Value: filter.Code}}})
	}
	if filter.MinPrice > 0 {
		f = append(f, bson.E{Key: "base_price", Value: bson.D{{Key: "$gte", Value: filter.MinPrice}}})
	}
	if filter.MaxPrice > 0 {
		f = append(f, bson.E{Key: "base_price", Value: bson.D{{Key: "$lte", Value: filter.MaxPrice}}})
	}
	if filter.StartDate > 0 {
		f = append(f, bson.E{Key: "start_date", Value: bson.D{{Key: "$gte", Value: filter.StartDate}}})
	}
	if filter.EndDate > 0 {
		f = append(f, bson.E{Key: "end_date", Value: bson.D{{Key: "$lte", Value: filter.EndDate}}})
	}
	if filter.StartTime > 0 {
		f = append(f, bson.E{Key: "start_time", Value: bson.D{{Key: "$gte", Value: filter.StartTime}}})
	}
	if filter.EndTime > 0 {
		f = append(f, bson.E{Key: "end_time", Value: bson.D{{Key: "$lte", Value: filter.EndTime}}})
	}
	cursor, err := collection.Find(ctx, f)
	if err != nil {
		return nil, "", fmt.Errorf("unable to find rates: %v: %w", err, cubawheeler.ErrInternal)
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var rate cubawheeler.Rate
		if err := cursor.Decode(&rate); err != nil {
			return nil, "", fmt.Errorf("unable to decode rate: %v: %w", err, cubawheeler.ErrInternal)
		}
		rates = append(rates, &rate)
		if len(rates) == filter.Limit+1 && filter.Limit > 0 {
			token = rates[filter.Limit].ID
			rates = rates[:filter.Limit]
		}
	}
	if err := cursor.Err(); err != nil {
		return nil, "", fmt.Errorf("unable to iterate over rates: %v: %w", err, cubawheeler.ErrInternal)
	}
	return rates, token, nil
}

func findRateByID(ctx context.Context, db *DB, id string) (*cubawheeler.Rate, error) {
	rates, _, err := findRates(ctx, db, &cubawheeler.RateFilter{Ids: []string{id}, Limit: 1})
	if err != nil {
		return nil, err
	}
	if len(rates) == 0 {
		return nil, fmt.Errorf("unable to find rate: %v: %w", id, cubawheeler.ErrNotFound)
	}
	return rates[0], nil
}

func insertRate(ctx context.Context, db *DB, rate *cubawheeler.Rate) error {
	collection := db.client.Database(database).Collection(RatesCollection.String())
	if _, err := collection.InsertOne(ctx, rate); err != nil {
		return fmt.Errorf("unable to store the rate: %v: %w", err, cubawheeler.ErrInternal)
	}
	return nil
}

func assembleRate(rate *cubawheeler.Rate, req cubawheeler.RateRequest) {
	id := req.ID
	if id == "" {
		id = cubawheeler.NewID().String()
	}
	rate.ID = id
	if len(req.Code) > 0 {
		rate.Code = req.Code
	}
	if req.BasePrice > 0 {
		rate.BasePrice = req.BasePrice
	}
	if req.PricePerMin != nil {
		rate.PricePerMin = *req.PricePerMin
	}
	if req.PricePerKm != nil {
		rate.PricePerKm = *req.PricePerKm
	}
	if req.PricePerPassenger != nil {
		rate.PricePerPassenger = *req.PricePerPassenger
	}
	if req.PricePerBaggage != nil {
		rate.PricePerBaggage = *req.PricePerBaggage
	}
	if len(req.StartTime) > 0 {
		rate.StartTime = req.StartTime
	}
	if len(req.EndTime) > 0 {
		rate.EndTime = req.EndTime
	}
	if req.StartDate != nil {
		rate.StartDate = *req.StartDate
	}
	if req.EndDate != nil {
		rate.EndDate = *req.EndDate
	}
	if req.MinKm != nil {
		rate.MinKm = *req.MinKm
	}
	if req.MaxKm != nil {
		rate.MaxKm = *req.MaxKm
	}
}
