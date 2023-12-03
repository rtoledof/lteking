package graph

import (
	"context"

	"cubawheeler.io/pkg/cubawheeler"
)

type chargeResolver struct{ *Resolver }

// ID is the resolver for the id field.
func (r *chargeResolver) ID(ctx context.Context, obj *cubawheeler.Charge) (string, error) {
	return obj.ID, nil
}

// Rider is the resolver for the rider field.
func (r *chargeResolver) Rider(ctx context.Context, obj *cubawheeler.Charge) (*cubawheeler.User, error) {
	return r.user.FindByID(ctx, obj.Rider)
}

// Trip is the resolver for the trip field.
func (r *chargeResolver) Trip(ctx context.Context, obj *cubawheeler.Charge) (*cubawheeler.Trip, error) {
	return r.trip.FindByID(ctx, obj.Trip)
}

type chargeRequestResolver struct{ *Resolver }

// Dispute is the resolver for the dispute field.
func (r *chargeRequestResolver) Dispute(ctx context.Context, obj *cubawheeler.ChargeRequest, data *bool) error {
	obj.Disputed = data
	return nil
}
