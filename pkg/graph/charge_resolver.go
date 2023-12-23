package graph

import (
	"context"

	"cubawheeler.io/pkg/cubawheeler"
	"cubawheeler.io/pkg/graph/model"
)

var _ ChargeResolver = &chargeResolver{}

type chargeResolver struct{ *Resolver }

// Amount implements ChargeResolver.
func (*chargeResolver) Amount(ctx context.Context, obj *cubawheeler.Charge) (*model.Amount, error) {
	return &model.Amount{
		Currency: obj.Amount.Currency.String(),
		Amount:   int(obj.Amount.Amount),
	}, nil
}

// Method implements ChargeResolver.
func (*chargeResolver) Method(ctx context.Context, obj *cubawheeler.Charge) (*string, error) {
	method := obj.Method.String()
	return &method, nil
}

// ID is the resolver for the id field.
func (r *chargeResolver) ID(ctx context.Context, obj *cubawheeler.Charge) (string, error) {
	return obj.ID, nil
}

// Rider is the resolver for the rider field.
func (r *chargeResolver) Rider(ctx context.Context, obj *cubawheeler.Charge) (*cubawheeler.User, error) {
	return r.user.FindByID(ctx, obj.Rider)
}

// Trip is the resolver for the trip field.
func (r *chargeResolver) Order(ctx context.Context, obj *cubawheeler.Charge) (*cubawheeler.Order, error) {
	return r.order.FindByID(ctx, obj.Order)
}

var _ ChargeRequestResolver = &chargeRequestResolver{}

type chargeRequestResolver struct{ *Resolver }

// Method implements ChargeRequestResolver.
func (*chargeRequestResolver) Method(ctx context.Context, obj *cubawheeler.ChargeRequest, data *string) error {
	obj.Method = cubawheeler.ChargeMethod(*data)
	return nil
}

// Dispute is the resolver for the dispute field.
func (r *chargeRequestResolver) Dispute(ctx context.Context, obj *cubawheeler.ChargeRequest, data *bool) error {
	obj.Disputed = data
	return nil
}

func (r *chargeRequestResolver) Order(ctx context.Context, obj *cubawheeler.ChargeRequest, data *string) error {
	obj.Order = *data
	return nil
}
