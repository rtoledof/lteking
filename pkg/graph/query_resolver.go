package graph

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

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
	value := url.Values{
		"limit":  []string{fmt.Sprintf("%d", filter.Limit)},
		"token":  []string{*filter.Token},
		"ids":    filter.IDs,
		"rider":  []string{*filter.Rider},
		"driver": []string{*filter.Driver},
		"status": []string{*filter.Status},
	}
	jwtToken := cubawheeler.JWTFromContext(ctx)

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s?%s", r.OrderService, value.Encode()), nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v: %w", err, cubawheeler.ErrInternal)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", jwtToken))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v: %w", err, cubawheeler.ErrInternal)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusUnauthorized {
			return nil, cubawheeler.ErrAccessDenied
		}
	}
	var orderList cubawheeler.OrderList
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %v: %w", err, cubawheeler.ErrInternal)
	}
	if err := json.Unmarshal(data, &orderList); err != nil {
		return nil, fmt.Errorf("error decoding response: %v: %w", err, cubawheeler.ErrInternal)
	}

	return &orderList, nil
}

// Charges is the resolver for the charges field.
func (r *queryResolver) Charges(ctx context.Context, filter cubawheeler.ChargeRequest) (*cubawheeler.ChargeList, error) {
	return r.charge.FindAll(ctx, filter)
}

// Profile is the resolver for the profile field.
func (r *queryResolver) Me(ctx context.Context) (*cubawheeler.Profile, error) {
	resp, err := makeRequest(ctx, http.MethodGet, r.AuthService, nil)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v: %w", err, cubawheeler.ErrInternal)
	}
	defer resp.Body.Close()
	var profile cubawheeler.Profile
	if err := json.NewDecoder(resp.Body).Decode(&profile); err != nil {
		return nil, fmt.Errorf("error decoding response: %v: %w", err, cubawheeler.ErrInternal)
	}
	return &profile, nil
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
	resp, err := makeRequest(ctx, http.MethodGet, fmt.Sprintf("%s/%s", r.OrderService, id), nil)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v: %w", err, cubawheeler.ErrInternal)
	}
	defer resp.Body.Close()

	var order cubawheeler.Order
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %v: %w", err, cubawheeler.ErrInternal)
	}
	if err := json.Unmarshal(data, &order); err != nil {
		return nil, fmt.Errorf("error decoding response: %v: %w", err, cubawheeler.ErrInternal)
	}

	return &order, nil
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
