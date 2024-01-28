package mongo

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"

	"order.io/pkg/order"
)

var _ order.RateService = &RateService{}

var RatesCollection Collections = "rates"

type RateService struct {
	db *DB
}

func NewRateService(db *DB) *RateService {
	return &RateService{
		db: db,
	}
}

func (s *RateService) Create(ctx context.Context, request order.RateRequest) (*order.Rate, error) {
	usr := order.UserFromContext(ctx)
	if usr == nil || usr.Role != order.RoleAdmin {
		return nil, order.ErrAccessDenied
	}

	var rate order.Rate
	assembleRate(&rate, request)
	if err := rate.Validate(); err != nil {
		return nil, err
	}
	if err := storeRate(ctx, s.db, &rate); err != nil {
		return nil, err
	}

	return &rate, nil
}

func (s *RateService) Update(ctx context.Context, request *order.RateRequest) (*order.Rate, error) {
	usr := order.UserFromContext(ctx)
	if usr == nil || usr.Role != order.RoleAdmin {
		return nil, order.ErrAccessDenied
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

func (s *RateService) FindByID(ctx context.Context, id string) (*order.Rate, error) {
	return findRateByID(ctx, s.db, id)
}

// FindByCode implements order.RateService.
func (s *RateService) FindByCode(ctx context.Context, code string) (*order.Rate, error) {
	usr := order.UserFromContext(ctx)
	if usr == nil || usr.Role != order.RoleAdmin {
		return nil, order.ErrAccessDenied
	}
	rates, _, err := findRates(ctx, s.db, &order.RateFilter{Code: []string{code}, Limit: 1})
	if err != nil {
		return nil, err
	}
	if len(rates) == 0 {
		return nil, fmt.Errorf("unable to find rate: %v: %w", code, order.ErrNotFound)
	}
	return rates[0], nil
}

func (s *RateService) FindAll(ctx context.Context, request order.RateFilter) ([]*order.Rate, string, error) {
	return findRates(ctx, s.db, &request)
}

func updateRate(ctx context.Context, db *DB, rate *order.Rate) error {
	collection := db.client.Database(database).Collection(RatesCollection.String())
	if _, err := collection.UpdateOne(ctx, bson.D{{Key: "_id", Value: rate.ID}}, bson.D{{Key: "$set", Value: rate}}); err != nil {
		return fmt.Errorf("unable to update rate: %v: %w", err, order.ErrInternal)
	}
	return nil
}

func findRates(ctx context.Context, db *DB, filter *order.RateFilter) ([]*order.Rate, string, error) {
	collection := db.client.Database(database).Collection(RatesCollection.String())
	var rates []*order.Rate
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
		return nil, "", fmt.Errorf("unable to find rates: %v: %w", err, order.ErrInternal)
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var rate order.Rate
		if err := cursor.Decode(&rate); err != nil {
			return nil, "", fmt.Errorf("unable to decode rate: %v: %w", err, order.ErrInternal)
		}
		rates = append(rates, &rate)
		if len(rates) == filter.Limit+1 && filter.Limit > 0 {
			token = rates[filter.Limit].ID
			rates = rates[:filter.Limit]
		}
	}
	if err := cursor.Err(); err != nil {
		return nil, "", fmt.Errorf("unable to iterate over rates: %v: %w", err, order.ErrInternal)
	}
	return rates, token, nil
}

func findRateByID(ctx context.Context, db *DB, id string) (*order.Rate, error) {
	rates, _, err := findRates(ctx, db, &order.RateFilter{Ids: []string{id}, Limit: 1})
	if err != nil {
		return nil, err
	}
	if len(rates) == 0 {
		return nil, fmt.Errorf("unable to find rate: %v: %w", id, order.ErrNotFound)
	}
	return rates[0], nil
}

func storeRate(ctx context.Context, db *DB, rate *order.Rate) error {
	collection := db.client.Database(database).Collection(RatesCollection.String())
	if _, err := collection.InsertOne(ctx, rate); err != nil {
		return fmt.Errorf("unable to store the rate: %v: %w", err, order.ErrInternal)
	}
	return nil
}

func assembleRate(rate *order.Rate, req order.RateRequest) {
	id := req.ID
	if id == "" {
		id = order.NewID().String()
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
		startDate := time.Unix(*req.StartDate, 0)
		rate.StartDate = startDate.Format("2006-01-02")
	}
	if req.EndDate != nil {
		endDate := time.Unix(*req.EndDate, 0)
		rate.EndDate = endDate.Format("2006-01-02")
	}
	if req.MinKm != nil {
		rate.MinKm = *req.MinKm
	}
	if req.MaxKm != nil {
		rate.MaxKm = *req.MaxKm
	}
}
