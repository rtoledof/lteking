package mongo

import (
	"context"
	"cubawheeler.io/pkg/cubawheeler"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var _ cubawheeler.PlanService = &PlanService{}

type PlanService struct {
	db         *DB
	collection *mongo.Collection
}

func NewPlanService(db *DB) *PlanService {
	return &PlanService{
		db:         db,
		collection: db.client.Database(database).Collection("plans"),
	}
}

func (s *PlanService) Create(ctx context.Context, request *cubawheeler.PlanRequest) (*cubawheeler.Plan, error) {
	plan := &cubawheeler.Plan{
		ID:         cubawheeler.NewID().String(),
		Name:       request.Name,
		Recurrintg: request.Recurring,
		Trips:      request.TotalTrips,
		Price:      request.Price,
		Interval:   request.Interval,
		Code:       request.Code,
	}
	_, err := s.collection.InsertOne(ctx, plan)
	if err != nil {
		return nil, fmt.Errorf("unable to create the plan: %w", err)
	}
	return plan, nil
}

func (s *PlanService) Update(ctx context.Context, request *cubawheeler.PlanRequest) (*cubawheeler.Plan, error) {
	//TODO implement me
	panic("implement me")
}

func (s *PlanService) FindByID(ctx context.Context, id string) (*cubawheeler.Plan, error) {
	plans, _, err := findAllPlans(ctx, s.collection, &cubawheeler.PlanFilter{Ids: []string{id}, Limit: 1})
	if err != nil {
		return nil, err
	}
	if len(plans) == 0 {
		return nil, errors.New("plan not found")
	}
	return plans[0], nil
}

func (s *PlanService) FindAll(ctx context.Context, filter *cubawheeler.PlanFilter) ([]*cubawheeler.Plan, string, error) {
	return findAllPlans(ctx, s.collection, filter)
}

func findAllPlans(ctx context.Context, collection *mongo.Collection, filter *cubawheeler.PlanFilter) ([]*cubawheeler.Plan, string, error) {
	var plans []*cubawheeler.Plan
	var token string
	f := bson.D{}
	// TODO: add filters here
	cur, err := collection.Find(ctx, f)
	if err != nil {
		return nil, "", err
	}
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var plan cubawheeler.Plan
		err := cur.Decode(&plan)
		if err != nil {
			return nil, "", err
		}
		plans = append(plans, &plan)
		if len(plans) > filter.Limit+1 {
			token = plans[filter.Limit].ID
			plans = plans[:filter.Limit]
			break
		}
	}
	if err := cur.Err(); err != nil {
		return nil, "", err
	}
	return plans, token, nil
}
