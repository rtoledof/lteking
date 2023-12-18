package mongo

import (
	"context"
	"errors"
	"fmt"

	"cubawheeler.io/pkg/cubawheeler"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var _ cubawheeler.PlanService = &PlanService{}

var PlanCollection Collections = "plans"

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
	id := request.ID
	if id == "" {
		id = cubawheeler.NewID().String()
	}
	plan := &cubawheeler.Plan{
		ID: id,
	}
	if err := assamblePlan(plan, request); err != nil {
		return nil, err
	}
	_, err := s.collection.InsertOne(ctx, plan)
	if err != nil {
		return nil, fmt.Errorf("unable to create the plan: %w", err)
	}
	return plan, nil
}

func (s *PlanService) Update(ctx context.Context, request *cubawheeler.PlanRequest) (*cubawheeler.Plan, error) {
	plan, err := findPlanById(ctx, s.db, request.ID)
	if err != nil {
		return nil, err
	}
	if err := assamblePlan(plan, request); err != nil {
		return nil, err
	}
	err = updatePlan(ctx, s.db, plan)
	if err != nil {
		return nil, err
	}
	return plan, nil
}

func (s *PlanService) FindByID(ctx context.Context, id string) (*cubawheeler.Plan, error) {
	return findPlanById(ctx, s.db, id)
}

func (s *PlanService) FindAll(ctx context.Context, filter *cubawheeler.PlanFilter) ([]*cubawheeler.Plan, string, error) {
	return findAllPlans(ctx, s.db, filter)
}

func findAllPlans(ctx context.Context, db *DB, filter *cubawheeler.PlanFilter) ([]*cubawheeler.Plan, string, error) {
	collection := db.client.Database(database).Collection(PlanCollection.String())
	var plans []*cubawheeler.Plan
	var token string
	f := bson.D{}
	if len(filter.Ids) > 0 {
		f = append(f, primitive.E{Key: "_id", Value: bson.D{{Key: "$in", Value: filter.Ids}}})
	}
	if filter.Name != "" {
		f = append(f, primitive.E{Key: "name", Value: filter.Name})
	}
	if filter.Price > 0 {
		f = append(f, primitive.E{Key: "price", Value: filter.Price})
	}
	if len(filter.Code) > 0 {
		f = append(f, primitive.E{Key: "code", Value: filter.Code})
	}
	if filter.TotalTrips > 0 {
		f = append(f, primitive.E{Key: "trips", Value: filter.TotalTrips})
	}
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
		if filter.Limit > 0 && len(plans) >= filter.Limit+1 {
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

func findPlanById(ctx context.Context, db *DB, id string) (*cubawheeler.Plan, error) {
	plans, _, err := findAllPlans(ctx, db, &cubawheeler.PlanFilter{Ids: []string{id}, Limit: 1})
	if err != nil {
		return nil, err
	}
	if len(plans) == 0 {
		return nil, errors.New("plan not found")
	}
	return plans[0], nil
}

func updatePlan(ctx context.Context, db *DB, plan *cubawheeler.Plan) error {
	collection := db.client.Database(database).Collection(PlanCollection.String())
	_, err := collection.UpdateOne(ctx, bson.D{{Key: "_id", Value: plan.ID}}, bson.D{{Key: "$set", Value: plan}})
	if err != nil {
		return fmt.Errorf("unable to update the plan: %w", err)
	}
	return nil
}

func assamblePlan(plan *cubawheeler.Plan, request *cubawheeler.PlanRequest) error {
	if request.Name != nil {
		plan.Name = *request.Name
	} else {
		return errors.New("name is required")
	}
	if request.Recurring != nil {
		plan.Recurrintg = *request.Recurring
	}
	if request.TotalTrips != nil {
		plan.Trips = *request.TotalTrips
	} else {
		return errors.New("total trips is required")
	}
	if request.Price != nil {
		plan.Price = *request.Price
	} else {
		return errors.New("price is required")
	}
	interval := cubawheeler.IntervalDay
	plan.Interval = interval
	if request.Interval != nil {
		switch *request.Interval {
		case cubawheeler.IntervalDay,
			cubawheeler.IntervalWeek,
			cubawheeler.IntervalMonth,
			cubawheeler.IntervalYear:
			plan.Interval = *request.Interval
		default:
			return errors.New("invalid interval")
		}
	}
	if request.Code != nil {
		plan.Code = *request.Code
	} else {
		return errors.New("code is required")
	}
	return nil
}
