package graph

import (
	"context"

	"cubawheeler.io/pkg/cubawheeler"
)

var _ UserResolver = &userResolver{}

type userResolver struct{ *Resolver }

func (r *userResolver) Orders(ctx context.Context, obj *cubawheeler.User) ([]*cubawheeler.Order, error) {
	return obj.Orders, nil
}

// ID is the resolver for the id field.
func (r *userResolver) ID(ctx context.Context, obj *cubawheeler.User) (string, error) {
	return obj.ID, nil
}

// Password is the resolver for the password field.
func (r *userResolver) Password(ctx context.Context, obj *cubawheeler.User) (*string, error) {
	pwd := ""
	return &pwd, nil
}

// Pin is the resolver for the pin field.
func (r *userResolver) Pin(ctx context.Context, obj *cubawheeler.User) (string, error) {
	return "", nil
}

// Otp is the resolver for the pin field.
func (r *userResolver) Otp(ctx context.Context, obj *cubawheeler.User) (string, error) {
	return "", nil
}

// ActiveVehicle is the resolver for the active_vehicle field.
func (r *userResolver) ActiveVehicle(ctx context.Context, obj *cubawheeler.User) (*cubawheeler.Vehicle, error) {
	return r.vehicle.FindByID(ctx, obj.ActiveVehicle)
}

// Plan is the resolver for the plan field.
func (r *userResolver) Plan(ctx context.Context, obj *cubawheeler.User) (*cubawheeler.Plan, error) {
	return r.plan.FindByID(ctx, obj.Plan)
}

// Reviews is the resolver for the reviews field.
func (r *userResolver) Reviews(ctx context.Context, obj *cubawheeler.User) ([]*cubawheeler.Review, error) {
	reviews, _, err := r.review.FindAll(ctx, cubawheeler.ReviewFilter{
		From: obj.ID,
	})
	return reviews, err
}