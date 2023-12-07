package graph

import (
	"context"
	"fmt"

	"cubawheeler.io/pkg/cubawheeler"
)

type queryResolver struct{ *Resolver }

// Users is the resolver for the users field.
func (r *queryResolver) Users(ctx context.Context, filter *cubawheeler.UserFilter) (*cubawheeler.UserList, error) {
	return r.user.FindAll(ctx, filter)
}

// Trips is the resolver for the trips field.
func (r *queryResolver) Trips(ctx context.Context, filter *cubawheeler.TripFilter) (*cubawheeler.TripList, error) {
	if filter == nil {
		filter = &cubawheeler.TripFilter{}
	}
	return r.user.Trips(ctx, filter)
}

// Charges is the resolver for the charges field.
func (r *queryResolver) Charges(ctx context.Context, filter cubawheeler.ChargeRequest) (*cubawheeler.ChargeList, error) {
	return r.charge.FindAll(ctx, filter)
}

// Profile is the resolver for the profile field.
func (r *queryResolver) Me(ctx context.Context) (*cubawheeler.Profile, error) {
	return r.user.Me(ctx)
}

// LastNAddress is the resolver for the lastNAddress field.
func (r *queryResolver) LastNAddress(ctx context.Context, number int) ([]*cubawheeler.Location, error) {
	panic(fmt.Errorf("not implemented: LastNAddress - lastNAddress"))
}

// Charge is the resolver for the charge field.
func (r *queryResolver) Charge(ctx context.Context, id *string) (*cubawheeler.Charge, error) {
	return r.charge.FindByID(ctx, *id)
}

// Trip is the resolver for the trip field.
func (r *queryResolver) Trip(ctx context.Context, id string) (*cubawheeler.Trip, error) {
	return r.trip.FindByID(ctx, id)
}

// FindVehicle is the resolver for the findVehicle field.
func (r *queryResolver) FindVehicle(ctx context.Context, vehicle string) (*cubawheeler.Vehicle, error) {
	return r.vehicle.FindByID(ctx, vehicle)
}

// FindApplications is the resolver for the findApplications field.
func (r *queryResolver) FindApplications(ctx context.Context, input *cubawheeler.ApplicationFilter) (*cubawheeler.ApplicationList, error) {
	panic(fmt.Errorf("not implemented: FindApplications - findApplications"))
}
