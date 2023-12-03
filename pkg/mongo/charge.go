package mongo

import (
	"context"
	"cubawheeler.io/pkg/cubawheeler"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var _ cubawheeler.ChargeService = &ChargeService{}

type ChargeService struct {
	db         *DB
	collection *mongo.Collection
}

func NewChargeService(db *DB) *ChargeService {
	return &ChargeService{
		db:         db,
		collection: db.client.Database(database).Collection("charges"),
	}
}

func (s *ChargeService) Create(ctx context.Context, request *cubawheeler.ChargeRequest) (*cubawheeler.Charge, error) {
	usr := cubawheeler.UserFromContext(ctx)
	if usr == nil {
		return nil, errors.New("invalid token provided")
	}
	if usr.Role != cubawheeler.RoleRider || usr.Role != cubawheeler.RoleAdmin {
		return nil, errors.New("access denied")
	}
	if usr.Role == cubawheeler.RoleRider {
		request.ReceiptEmail = usr.Email
	}
	id := cubawheeler.NewID().String()
	charge := &cubawheeler.Charge{
		ID:                id,
		Amount:            request.Amount,
		Currency:          request.Currency,
		Description:       request.Description,
		Trip:              request.Trip,
		Disputed:          request.Disputed,
		ReceiptEmail:      request.ReceiptEmail,
		Status:            cubawheeler.ChargeStatusPending,
		Method:            request.Method,
		ExternalReference: request.Reference,
	}
	// TODO: calculate the fees and apply it to the charge
	_, err := s.collection.InsertOne(ctx, charge)
	if err != nil {
		return nil, fmt.Errorf("unable to store the charge: %w", err)
	}
	return charge, nil
}

func (s *ChargeService) Update(ctx context.Context, request *cubawheeler.ChargeRequest) (*cubawheeler.Charge, error) {
	//TODO implement me
	panic("implement me")
}

func (s *ChargeService) FindByID(ctx context.Context, id string) (*cubawheeler.Charge, error) {
	charges, _, err := findAllCharges(ctx, s.collection, cubawheeler.ChargeRequest{
		Ids:   []string{id},
		Limit: 1,
	})
	if err != nil {
		return nil, err
	}
	return charges[0], nil
}

func (s *ChargeService) FindAll(ctx context.Context, request cubawheeler.ChargeRequest) (*cubawheeler.ChargeList, error) {
	usr := cubawheeler.UserFromContext(ctx)
	if usr == nil {
		return nil, errors.New("invalid token provided")
	}
	switch usr.Role {
	case cubawheeler.RoleRider:
		request.Rider = &usr.ID
	case cubawheeler.RoleDriver:
		request.Driver = &usr.ID
	}
	charges, token, err := findAllCharges(ctx, s.collection, request)
	if err != nil {
		return nil, err
	}
	return &cubawheeler.ChargeList{Data: charges, Token: token}, nil
}

func findAllCharges(ctx context.Context, collection *mongo.Collection, filter cubawheeler.ChargeRequest) ([]*cubawheeler.Charge, string, error) {
	var charges []*cubawheeler.Charge
	var token string
	f := bson.D{}
	if len(filter.Ids) > 0 {
		f = append(f, bson.E{"_id", bson.A{"$in", filter.Ids}})
	}
	if filter.Rider != nil {
		f = append(f, bson.E{"rider", filter.Rider})
	}
	if filter.Driver != nil {
		f = append(f, bson.E{"driver", filter.Driver})
	}
	if len(filter.Token) > 0 {
		f = append(f, bson.E{"_id", bson.A{"$gt", filter.Token}})
	}

	cur, err := collection.Find(ctx, f)
	if err != nil {
		return nil, "", err
	}
	for cur.Next(ctx) {
		var charge cubawheeler.Charge
		err := cur.Decode(&charge)
		if err != nil {
			return nil, "", err
		}
		charges = append(charges, &charge)
		if len(charges) == filter.Limit+1 {
			token = charges[filter.Limit].ID
			charges = charges[:filter.Limit]
			break
		}
	}
	if err := cur.Err(); err != nil {
		return nil, "", err
	}
	cur.Close(ctx)
	return charges, token, nil
}
