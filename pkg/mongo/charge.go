package mongo

import (
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"

	"cubawheeler.io/pkg/cubawheeler"
	"cubawheeler.io/pkg/currency"
	"cubawheeler.io/pkg/derrors"
)

var _ cubawheeler.ChargeService = &ChargeService{}

var ChargesCollection Collections = "charges"

type ChargeService struct {
	db *DB
}

func NewChargeService(db *DB) *ChargeService {
	return &ChargeService{
		db: db,
	}
}

func (s *ChargeService) Create(ctx context.Context, request *cubawheeler.ChargeRequest) (_ *cubawheeler.Charge, err error) {
	defer derrors.Wrap(&err, "mongo.ChargeService.Create")
	usr := cubawheeler.UserFromContext(ctx)
	if usr == nil {
		return nil, errors.New("invalid token provided")
	}
	if usr.Role != cubawheeler.RoleRider && usr.Role != cubawheeler.RoleDriver {
		return nil, errors.New("access denied")
	}
	if usr.Role == cubawheeler.RoleRider {
		request.ReceiptEmail = usr.Email
	}
	id := cubawheeler.NewID().String()
	cur, err := currency.Parse(request.Currency)
	if err != nil {
		return nil, fmt.Errorf("invalid currency: %w", err)
	}
	charge := &cubawheeler.Charge{
		ID: id,
		Amount: currency.Amount{
			Currency: cur,
			Amount:   int64(request.Amount),
		},
		Description:       request.Description,
		Order:             request.Order,
		Disputed:          *request.Disputed,
		ReceiptEmail:      request.ReceiptEmail,
		Status:            cubawheeler.ChargeStatusPending,
		Method:            request.Method,
		ExternalReference: *request.Reference,
	}
	// TODO: calculate the fees and apply it to the charge
	collection := s.db.Collection(ChargesCollection)
	_, err = collection.InsertOne(ctx, charge)
	if err != nil {
		return nil, fmt.Errorf("unable to store the charge: %w", err)
	}
	return charge, nil
}

func (s *ChargeService) Update(ctx context.Context, request *cubawheeler.ChargeRequest) (_ *cubawheeler.Charge, err error) {
	defer derrors.Wrap(&err, "mongo.ChargeService.Update")
	//TODO implement me
	panic("implement me")
}

func (s *ChargeService) FindByID(ctx context.Context, id string) (_ *cubawheeler.Charge, err error) {
	defer derrors.Wrap(&err, "mongo.ChargeService.FindByID")
	charges, _, err := findAllCharges(ctx, s.db, cubawheeler.ChargeRequest{
		Ids:   []string{id},
		Limit: 1,
	})
	if err != nil {
		return nil, err
	}
	return charges[0], nil
}

func (s *ChargeService) FindAll(ctx context.Context, request cubawheeler.ChargeRequest) (_ *cubawheeler.ChargeList, err error) {
	defer derrors.Wrap(&err, "mongo.ChargeService.FindAll")
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
	charges, token, err := findAllCharges(ctx, s.db, request)
	if err != nil {
		return nil, err
	}
	return &cubawheeler.ChargeList{Data: charges, Token: token}, nil
}

func findAllCharges(ctx context.Context, db *DB, filter cubawheeler.ChargeRequest) ([]*cubawheeler.Charge, string, error) {
	var charges []*cubawheeler.Charge
	var token string
	collection := db.Collection(ChargesCollection)
	f := bson.D{}
	if len(filter.Ids) > 0 {
		f = append(f, bson.E{Key: "_id", Value: bson.A{"$in", filter.Ids}})
	}
	if filter.Rider != nil {
		f = append(f, bson.E{Key: "rider", Value: filter.Rider})
	}
	if filter.Driver != nil {
		f = append(f, bson.E{Key: "driver", Value: filter.Driver})
	}
	if len(filter.Token) > 0 {
		f = append(f, bson.E{Key: "_id", Value: bson.A{"$gt", filter.Token}})
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
