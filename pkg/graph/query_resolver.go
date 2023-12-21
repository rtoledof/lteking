package graph

import (
	"context"
	"fmt"

	"cubawheeler.io/pkg/cubawheeler"
)

var _ QueryResolver = &queryResolver{}

type queryResolver struct{ *Resolver }

// Users is the resolver for the users field.
func (r *queryResolver) Users(ctx context.Context, filter *cubawheeler.UserFilter) (*cubawheeler.UserList, error) {
	return r.user.FindAll(ctx, filter)
}

// Trips is the resolver for the trips field.
func (r *queryResolver) Orders(ctx context.Context, filter *cubawheeler.OrderFilter) (*cubawheeler.OrderList, error) {
	if filter == nil {
		filter = &cubawheeler.OrderFilter{}
	}
	return r.user.Orders(ctx, filter)
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
func (r *queryResolver) Order(ctx context.Context, id string) (*cubawheeler.Order, error) {
	return r.order.FindByID(ctx, id)
}

// FindVehicle is the resolver for the findVehicle field.
func (r *queryResolver) FindVehicle(ctx context.Context, vehicle string) (*cubawheeler.Vehicle, error) {
	return r.vehicle.FindByID(ctx, vehicle)
}

// FindApplications is the resolver for the findApplications field.
func (r *queryResolver) FindApplications(ctx context.Context, input *cubawheeler.ApplicationFilter) (*cubawheeler.ApplicationList, error) {
	panic(fmt.Errorf("not implemented: FindApplications - findApplications"))
}

// NearByDrivers is the resolver for the nearByDrivers field.
func (r *queryResolver) NearByDrivers(ctx context.Context, input *cubawheeler.PointInput) ([]*cubawheeler.NearByResponse, error) {
	locations, err := r.realTimeLocation.FindNearByDrivers(ctx, cubawheeler.GeoLocation{
		Type: "Point",
		Long: input.Lng,
		Lat:  input.Lat,
	})
	if err != nil {
		return nil, err
	}
	var users []string
	var userLocations = make(map[string]*cubawheeler.Location)
	for _, l := range locations {
		users = append(users, l.User)
		userLocations[l.User] = l
	}

	rsp, err := r.user.FindAll(ctx, &cubawheeler.UserFilter{
		Ids:    users,
		Status: []cubawheeler.UserStatus{cubawheeler.UserStatusActive},
	})
	if err != nil {
		return nil, err
	}
	var nearByResponses []*cubawheeler.NearByResponse
	for _, v := range rsp.Data {
		rsp := cubawheeler.NearByResponse{
			Driver: v,
		}
		if l, ok := userLocations[v.ID]; ok {
			rsp.Location = l
			nearByResponses = append(nearByResponses, &rsp)
		}
	}

	return nearByResponses, nil
}
