package graph

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"

	"cubawheeler.io/pkg/cubawheeler"
	"cubawheeler.io/pkg/graph/model"
)

var _ MutationResolver = &mutationResolver{}

type mutationResolver struct{ *Resolver }

// ConfirmTransaction implements MutationResolver.
func (r *mutationResolver) ConfirmTransaction(ctx context.Context, id string, pin string) (*cubawheeler.Response, error) {
	var response = cubawheeler.Response{
		Success: true,
		Code:    http.StatusNoContent,
	}

	value := url.Values{
		"pin": []string{pin},
		"id":  []string{id},
	}
	resp, err := makeRequest(ctx, http.MethodPost, fmt.Sprintf("%s/v1/wallet/transfer/confirm", r.WalletService), value)
	if err != nil {
		response.Success = false
		response.Message = err.Error()
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNoContent {
		response.Success = false
		response.Message = fmt.Sprintf("error confirming transaction: %s", resp.Status)
	}
	return &response, nil
}

// Transfer implements MutationResolver.
func (r *mutationResolver) Transfer(ctx context.Context, to string, amount int, currency string, typeArg model.TransferType) (*model.Transaction, error) {
	value := url.Values{
		"to":       []string{to},
		"amount":   []string{fmt.Sprintf("%d", amount)},
		"currency": []string{currency},
		"type":     []string{string(typeArg)},
	}
	resp, err := makeRequest(ctx, http.MethodPost, fmt.Sprintf("%s/v1/wallet/transfer", r.WalletService), value)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v: %w", err, cubawheeler.ErrInternal)
	}
	defer resp.Body.Close()
	var transaction model.Transaction
	if err := json.NewDecoder(resp.Body).Decode(&transaction); err != nil {
		return nil, fmt.Errorf("error decoding response: %v: %w", err, cubawheeler.ErrInternal)
	}
	return &transaction, nil
}

// Redeem implements MutationResolver.
func (r *mutationResolver) Redeem(ctx context.Context, input string) (*cubawheeler.Response, error) {
	var rsp = cubawheeler.Response{
		Success: true,
		Code:    http.StatusOK,
	}
	_, err := r.coupon.Redeem(ctx, input)
	if err != nil {
		rsp.Success = false
		rsp.Code = http.StatusBadRequest
		rsp.Message = err.Error()
	}
	return &rsp, nil
}

// CreateOrder implements MutationResolver.
func (r *mutationResolver) CreateOrder(ctx context.Context, input *cubawheeler.DirectionRequest) (*cubawheeler.Order, error) {
	value := url.Values{
		"coupon":   {input.Coupon},
		"riders":   {strconv.Itoa(input.Riders)},
		"baggages": {strconv.FormatBool(input.Baggages)},
		"currency": {input.Currency},
	}
	for _, v := range input.Points {
		value.Add("points", fmt.Sprintf("%f,%f", v.Lat, v.Lng))
	}

	resp, err := makeRequest(ctx, http.MethodPost, fmt.Sprintf("%s/v1/orders", r.OrderService), value)
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

// UpdateOrder implements MutationResolver.
func (r *mutationResolver) UpdateOrder(ctx context.Context, input *cubawheeler.DirectionRequest) (*cubawheeler.Order, error) {
	value := url.Values{
		"coupon":   {input.Coupon},
		"riders":   {strconv.Itoa(input.Riders)},
		"baggages": {strconv.FormatBool(input.Baggages)},
		"currency": {input.Currency},
	}
	for _, v := range input.Points {
		value.Add("points", fmt.Sprintf("%f,%f", v.Lat, v.Lng))
	}

	resp, err := makeRequest(ctx, http.MethodPut, fmt.Sprintf("%s/%s", r.OrderService, input.ID), value)
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

// CancelOrder implements MutationResolver.
func (r *mutationResolver) CancelOrder(ctx context.Context, order string) (*cubawheeler.Response, error) {
	response := cubawheeler.Response{
		Success: true,
		Code:    http.StatusOK,
	}
	resp, err := makeRequest(ctx, http.MethodPost, fmt.Sprintf("%s/%s/cancel", r.OrderService, order), nil)
	if err != nil {
		response.Success = false
		response.Code = http.StatusBadRequest
		response.Message = err.Error()
	}
	defer resp.Body.Close()
	return &response, nil
}

// StartOrder implements MutationResolver.
func (r *mutationResolver) StartOrder(ctx context.Context, order string) (*cubawheeler.Response, error) {
	response := cubawheeler.Response{
		Success: true,
		Code:    http.StatusOK,
	}
	resp, err := makeRequest(ctx, http.MethodPost, fmt.Sprintf("%s/%s/start", r.OrderService, order), nil)
	if err != nil {
		response.Success = false
		response.Code = http.StatusBadRequest
		response.Message = err.Error()
	}
	defer resp.Body.Close()
	return &response, nil
}

func (r *mutationResolver) AcceptOrder(ctx context.Context, order string) (*cubawheeler.Response, error) {
	panic("implement me")
}

func (r *mutationResolver) ConfirmOrder(ctx context.Context, req cubawheeler.ConfirmOrder) (*cubawheeler.Response, error) {
	response := cubawheeler.Response{
		Success: true,
		Code:    http.StatusOK,
		Message: "order confirmed",
	}
	value := url.Values{
		"category": {string(req.Category)},
		"method":   {string(req.Method)},
		"currency": {req.Currency},
	}
	resp, err := makeRequest(ctx, http.MethodPost, fmt.Sprintf("%s/v1/orders/%s/confirm", r.OrderService, req.OrderID), value)
	if err != nil {
		response.Success = false
		response.Code = http.StatusBadRequest
		response.Message = err.Error()
	}
	defer resp.Body.Close()
	return &response, nil
}

// Login is the resolver for the login field.
func (r *mutationResolver) Login(ctx context.Context, email, otp string) (*cubawheeler.Token, error) {
	value := url.Values{
		"grant_type": {"password"},
		"username":   {email},
		"password":   {otp},
	}

	var token cubawheeler.Token

	resp, err := makeRequest(ctx, http.MethodPost, fmt.Sprintf("%s/login", r.AuthService), value)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v: %w", err, cubawheeler.ErrInternal)
	}
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(&token); err != nil {
		return nil, fmt.Errorf("error decoding response: %v: %w", err, cubawheeler.ErrInternal)
	}

	return &token, nil
}

func (r *mutationResolver) Authorize(ctx context.Context, clientID string, clientSecret string, refreshToken *string) (*cubawheeler.Token, error) {
	value := url.Values{
		"grant_type":    {"client_credentials"},
		"client_id":     {clientID},
		"client_secret": {clientSecret},
	}
	var token cubawheeler.Token
	authURL := fmt.Sprintf("%s/authorize", r.AuthService)
	slog.Info(fmt.Sprintf("authURL: %s", authURL))

	req, err := http.NewRequest(http.MethodPost, authURL, bytes.NewBufferString(value.Encode()))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v: %w", err, cubawheeler.ErrInternal)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v: %w", err, cubawheeler.ErrInternal)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		slog.Info("error decoding response: %v: %w", err, cubawheeler.ErrInternal)
		return nil, fmt.Errorf("error decoding response: %v: %w", err, cubawheeler.ErrAccessDenied)
	}
	responseData, _ := io.ReadAll(resp.Body)
	slog.Info(fmt.Sprintf("responseData: %s", responseData))
	if err := json.NewDecoder(bytes.NewBuffer(responseData)).Decode(&token); err != nil {
		return nil, fmt.Errorf("error decoding response: %v: %w", err, cubawheeler.ErrInternal)
	}
	return &token, nil
}

// Otp is the resolver for the otp field.
func (r *mutationResolver) Otp(ctx context.Context, email string) (*cubawheeler.Response, error) {
	var rsp = cubawheeler.Response{
		Success: true,
		Code:    http.StatusOK,
	}
	value := url.Values{
		"email": {email},
	}

	resp, err := makeRequest(ctx, http.MethodPost, fmt.Sprintf("%s/otp", r.AuthService), value)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var code string
	if err := json.NewDecoder(resp.Body).Decode(&code); err != nil {
		slog.Info("error decoding response: %v: %w", err, cubawheeler.ErrInternal)
	}
	rsp.Message = code
	return &rsp, nil
}

// UpdateProfile is the resolver for the updateProfile field.
func (r *mutationResolver) UpdateProfile(ctx context.Context, profile cubawheeler.UpdateProfile) (*cubawheeler.Response, error) {
	var rsp = cubawheeler.Response{
		Success: true,
		Code:    http.StatusOK,
	}
	value := url.Values{}

	if profile.Name != nil {
		value.Add("name", *profile.Name)
	}
	if profile.LastName != nil {
		value.Add("last_name", *profile.LastName)
	}
	if profile.Dob != nil {
		value.Add("dob", *profile.Dob)
	}
	if profile.Phone != nil {
		value.Add("phone", *profile.Phone)
	}
	if profile.Photo != nil {
		value.Add("photo", *profile.Photo)
	}
	if profile.Gender != nil {
		value.Add("gender", profile.Gender.String())
	}
	if profile.Dni != nil {
		value.Add("dni", *profile.Dni)
	}

	_, err := makeRequest(ctx, http.MethodPut, fmt.Sprintf("%s/profile", r.AuthService), value)
	if err != nil {
		rsp.Success = false
		rsp.Message = err.Error()
		rsp.Code = http.StatusBadRequest
		return nil, err
	}
	return &rsp, nil
}

// AddFavoritePlace is the resolver for the addFavoritePlace field.
func (r *mutationResolver) AddFavoritePlace(ctx context.Context, input cubawheeler.AddPlace) (*cubawheeler.Location, error) {
	return r.user.AddFavoritePlace(ctx, input)
}

// FavoritePlaces is the resolver for the favoritePlaces field.
func (r *mutationResolver) FavoritePlaces(ctx context.Context) ([]*cubawheeler.Location, error) {
	return r.user.FavoritePlaces(ctx)
}

// AddFavoriteVehicle is the resolver for the addFavoriteVehicle field.
func (r *mutationResolver) AddFavoriteVehicle(ctx context.Context, plate *string) (*cubawheeler.Vehicle, error) {
	return r.user.AddFavoriteVehicle(ctx, plate)
}

// FavoriteVehicles is the resolver for the favoriteVehicles field.
func (r *mutationResolver) FavoriteVehicles(ctx context.Context) ([]*cubawheeler.Vehicle, error) {
	return r.user.FavoriteVehicles(ctx)
}

// UpdatePlace is the resolver for the updatePlace field.
func (r *mutationResolver) UpdatePlace(ctx context.Context, input *cubawheeler.UpdatePlace) (*cubawheeler.Location, error) {
	return r.user.UpdatePlace(ctx, input)
}

// FindVehicle is the resolver for the findVehicle field.
func (r *mutationResolver) FindVehicle(ctx context.Context, plate string) (*cubawheeler.Vehicle, error) {
	return r.vehicle.FindByPlate(ctx, plate)
}

// UpdateVehicle is the resolver for the updateVehicle field.
func (r *mutationResolver) UpdateVehicle(ctx context.Context, input *cubawheeler.UpdateVehicle) (*cubawheeler.Vehicle, error) {
	return r.vehicle.Update(ctx, *input)
}

// CancelTrip is the resolver for the cancelTrip field.
func (r *mutationResolver) CancelTrip(ctx context.Context, trip string) (*cubawheeler.Order, error) {
	return r.order.CancelOrder(ctx, trip)
}

// AddRate is the resolver for the addRate field.
func (r *mutationResolver) AddRate(ctx context.Context, input cubawheeler.RateRequest) (*cubawheeler.Rate, error) {
	return r.rate.Create(ctx, input)
}

// UpdateRate is the resolver for the updateRate field.
func (r *mutationResolver) UpdateRate(ctx context.Context, input cubawheeler.RateRequest) (*cubawheeler.Rate, error) {
	return r.rate.Update(ctx, &input)
}

func (r *mutationResolver) ChangePin(ctx context.Context, old *string, pin string) (*cubawheeler.Response, error) {
	var rsp = cubawheeler.Response{
		Success: true,
		Code:    http.StatusOK,
	}
	_, err := r.profile.ChangePin(ctx, old, pin)
	if err != nil {
		rsp.Success = false
		rsp.Message = err.Error()
		rsp.Code = http.StatusBadRequest
	}
	return &rsp, nil
}

// CreateApplication is the resolver for the createApplication field.
func (r *mutationResolver) CreateApplication(ctx context.Context, input cubawheeler.ApplicationRequest) (*cubawheeler.Application, error) {
	return r.app.CreateApplication(ctx, input)
}

// UpdateApplicationCredentials is the resolver for the updateApplicationCredentials field.
func (r *mutationResolver) UpdateApplicationCredentials(ctx context.Context, application string) (*cubawheeler.Application, error) {
	return r.app.UpdateApplicationCredentials(ctx, application)
}

// AddDevice is the resolver for the addDevice field.
func (r *mutationResolver) AddDevice(ctx context.Context, device string) (*cubawheeler.Response, error) {
	var rsp = cubawheeler.Response{
		Success: true,
		Code:    http.StatusOK,
		Message: "device added",
	}
	value := url.Values{
		"device_id": {device},
	}
	_, err := makeRequest(ctx, http.MethodPost, fmt.Sprintf("%s/profile/devices", r.AuthService), value)
	if err != nil {
		rsp.Message = err.Error()
		rsp.Code = http.StatusBadRequest
		rsp.Success = false
	}
	return &rsp, nil
}
