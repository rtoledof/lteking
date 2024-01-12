package graph

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"cubawheeler.io/pkg/cannon"
	"cubawheeler.io/pkg/cubawheeler"
	"cubawheeler.io/pkg/graph/model"
)

var _ QueryResolver = &queryResolver{}

type queryResolver struct{ *Resolver }

// Balance implements QueryResolver.
func (r *queryResolver) Balance(ctx context.Context) ([]*model.Balance, error) {
	resp, err := makeRequest(ctx, http.MethodGet, fmt.Sprintf("%s/v1/wallet", r.WalletService), nil)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v: %w", err, cubawheeler.ErrInternal)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error making request: %v: %w", err, cubawheeler.ErrInternal)
	}

	var balance cubawheeler.Balance
	if err := json.NewDecoder(resp.Body).Decode(&balance); err != nil {
		return nil, fmt.Errorf("error decoding response: %v: %w", err, cubawheeler.ErrInternal)
	}
	var rsp []*model.Balance
	for currency, value := range balance.Amount {
		rsp = append(rsp, &model.Balance{
			Currency: currency,
			Amount:   int(value),
		})
	}
	return rsp, nil
}

// Transactions implements QueryResolver.
func (r *queryResolver) Transactions(ctx context.Context) ([]*model.Transaction, error) {
	resp, err := makeRequest(ctx, http.MethodGet, fmt.Sprintf("%s/v1/wallet/transactions", r.WalletService), nil)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v: %w", err, cubawheeler.ErrInternal)
	}
	defer resp.Body.Close()
	var transactions []*model.Transaction
	if err := json.NewDecoder(resp.Body).Decode(&transactions); err != nil {
		return nil, fmt.Errorf("error decoding response: %v: %w", err, cubawheeler.ErrInternal)
	}
	return transactions, nil
}

// PaymentMethods implements QueryResolver.
func (*queryResolver) PaymentMethods(ctx context.Context) ([]*model.ChargeMethod, error) {

	return []*model.ChargeMethod{
		{
			Name:        cubawheeler.ChargeMethodCash.String(),
			Description: "Pago en Efectivo",
		},
		{
			Name:        cubawheeler.ChargeMethodCUPTransaction.String(),
			Description: "Transferencia Bancaria (CUP)",
		},
		{
			Name:        cubawheeler.ChargeMethodMLCTransaction.String(),
			Description: "Transferencia Bancaria (MLC)",
		},
		{
			Name:        cubawheeler.ChargeMethodBalance.String(),
			Description: "Balance",
		},
	}, nil
}

// Users is the resolver for the users field.
func (r *queryResolver) Users(ctx context.Context, filter *cubawheeler.UserFilter) (*cubawheeler.UserList, error) {
	return r.user.FindAll(ctx, filter)
}

// Trips is the resolver for the trips field.
func (r *queryResolver) Orders(ctx context.Context, filter *cubawheeler.OrderFilter) (*cubawheeler.OrderList, error) {
	if filter == nil {
		filter = &cubawheeler.OrderFilter{}
	}
	if filter.Limit == 0 {
		filter.Limit = 10
	}
	value := url.Values{
		"limit": []string{fmt.Sprintf("%d", filter.Limit)},
		"ids":   filter.IDs,
	}
	if filter.Token != nil {
		value.Add("token", *filter.Token)
	}
	if filter.Rider != nil {
		value.Add("rider", *filter.Rider)
	}
	if filter.Driver != nil {
		value.Add("driver", *filter.Driver)
	}
	if filter.Status != nil {
		value.Add("status", *filter.Status)
	}

	resp, err := makeRequest(ctx, http.MethodGet, fmt.Sprintf("%s/v1/orders?%s", r.OrderService, value.Encode()), nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	orderList := new(cubawheeler.OrderList)
	if err := json.NewDecoder(resp.Body).Decode(orderList); err != nil {
		return nil, fmt.Errorf("error decoding response: %v: %w", err, cubawheeler.ErrInternal)
	}

	return orderList, nil
}

// Charges is the resolver for the charges field.
func (r *queryResolver) Charges(ctx context.Context, filter cubawheeler.ChargeRequest) (*cubawheeler.ChargeList, error) {
	return r.charge.FindAll(ctx, filter)
}

// Profile is the resolver for the profile field.
func (r *queryResolver) Me(ctx context.Context) (*model.ProfileOutput, error) {
	resp, err := makeRequest(ctx, http.MethodGet, fmt.Sprintf("%s/me", r.AuthService), nil)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v: %w", err, cubawheeler.ErrInternal)
	}
	defer resp.Body.Close()
	var user cubawheeler.User
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, fmt.Errorf("error decoding response: %v: %w", err, cubawheeler.ErrInternal)
	}
	return assambleUser(&user), nil
}

func assambleUser(user *cubawheeler.User) *model.ProfileOutput {
	return &model.ProfileOutput{
		ID:       user.ID,
		Email:    user.Email,
		Name:     user.Profile.Name,
		LastName: user.Profile.LastName,
		Phone:    user.Profile.Phone,
		Dob:      user.Profile.DOB,
		Photo:    user.Profile.Photo,
		Rate:     user.Rate,
		Status:   user.Profile.Status.String(),
		Gender:   user.Profile.Gender.String(),
	}
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
	logger := cannon.LoggerFromContext(ctx)
	logger.Info("Making request to %s", fmt.Sprintf("%s/v1/orders/%s", r.OrderService, id), nil)
	resp, err := makeRequest(ctx, http.MethodGet, fmt.Sprintf("%s/v1/orders/%s", r.OrderService, id), nil)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v: %w", err, cubawheeler.ErrInternal)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v: %w", err, cubawheeler.ErrInternal)
	}
	logger.Info("Response body: %s", string(data), resp.StatusCode)
	defer resp.Body.Close()

	var order cubawheeler.Order
	if err := json.NewDecoder(bytes.NewBuffer(data)).Decode(&order); err != nil {
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

	if len(users) == 0 {
		return nil, fmt.Errorf("no drivers found")
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
