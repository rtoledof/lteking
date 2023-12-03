package mongo

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/mongo"

	"cubawheeler.io/pkg/cubawheeler"
)

var _ cubawheeler.RateService = &RateService{}

type RateService struct {
	db         *DB
	collection *mongo.Collection
}

func NewRateService(db *DB) *RateService {
	return &RateService{
		db:         db,
		collection: db.client.Database(database).Collection("rates"),
	}
}

func (s *RateService) Create(ctx context.Context, request *cubawheeler.RateRequest) (*cubawheeler.Rate, error) {
	usr := cubawheeler.UserFromContext(ctx)
	if usr == nil {
		return nil, errors.New("invalid token provided")
	}
	if usr.Role != cubawheeler.RoleAdmin {
		return nil, errors.New("access denied")
	}
	rate := &cubawheeler.Rate{
		ID:                cubawheeler.NewID().String(),
		Code:              request.Code,
		BasePrice:         request.BasePrice,
		PricePerMin:       *request.PricePerMin,
		PricePerKm:        *request.PricePerKm,
		PricePerPassenger: request.PricePerPassenger,
		PricePerBaggage:   *request.PricePerBaggage,
		StartDate:         request.StartDate,
		EndDate:           request.EndDate,
		StartTime:         request.StartTime,
		EndTime:           request.EndTime,
		MinKm:             request.MinKm,
		MaxKm:             request.MaxKm,
	}
	if _, err := s.collection.InsertOne(ctx, rate); err != nil {
		return nil, err
	}
	return rate, nil
}

func (s *RateService) Update(ctx context.Context, request *cubawheeler.RateRequest) (*cubawheeler.Rate, error) {
	//TODO implement me
	panic("implement me")
}

func (s *RateService) FindID(ctx context.Context, id string) (*cubawheeler.Rate, error) {
	//TODO implement me
	panic("implement me")
}

func (s *RateService) FindAll(ctx context.Context, request *cubawheeler.RateRequest) ([]cubawheeler.Rate, string, error) {
	//TODO implement me
	panic("implement me")
}
