package graph

import (
	"context"
	"fmt"
	"net/http"

	"cubawheeler.io/pkg/cubawheeler"
)

var _ MutationResolver = &mutationResolver{}

type mutationResolver struct{ *Resolver }

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
	return r.order.Create(ctx, input)
}

// UpdateOrder implements MutationResolver.
func (r *mutationResolver) UpdateOrder(ctx context.Context, input *cubawheeler.DirectionRequest) (*cubawheeler.Order, error) {
	return r.order.Update(ctx, input)
}

// CancelOrder implements MutationResolver.
func (r *mutationResolver) CancelOrder(ctx context.Context, order string) (*cubawheeler.Response, error) {
	response := cubawheeler.Response{
		Success: true,
		Code:    http.StatusOK,
	}
	_, err := r.order.CancelOrder(ctx, order)
	if err != nil {
		response.Success = false
		response.Code = http.StatusBadRequest
		response.Message = err.Error()
	}
	return &response, nil
}

// StartOrder implements MutationResolver.
func (r *mutationResolver) StartOrder(ctx context.Context, order string) (*cubawheeler.Response, error) {
	response := cubawheeler.Response{
		Success: true,
		Code:    http.StatusOK,
	}
	if _, err := r.order.StartOrder(ctx, order); err != nil {
		response.Success = false
		response.Code = http.StatusBadRequest
		response.Message = err.Error()
	}
	return &response, nil
}

func (r *mutationResolver) AcceptOrder(ctx context.Context, order string) (*cubawheeler.Response, error) {
	response := cubawheeler.Response{
		Success: true,
		Code:    http.StatusOK,
	}
	_, err := r.order.AcceptOrder(ctx, order)
	if err != nil {
		response.Success = false
		response.Code = http.StatusBadRequest
	}
	return &response, nil
}

func (r *mutationResolver) ConfirmOrder(ctx context.Context, order string, cost string) (*cubawheeler.Response, error) {
	// response := cubawheeler.Response{
	// 	Success: true,
	// 	Code:    http.StatusOK,
	// }
	// _, err := r.order.ConfirmOrder(ctx, order, cost)
	// if err != nil {
	// 	response.Success = false
	// 	response.Code = http.StatusBadRequest
	// }
	// return &response, nil
	panic(fmt.Errorf("not implemented: ConfirmOrder - confirmOrder"))
}

// Login is the resolver for the login field.
func (r *mutationResolver) Login(ctx context.Context, input cubawheeler.LoginRequest) (*cubawheeler.Token, error) {
	if err := r.otp.Otp(ctx, input.Otp, input.Email); err != nil {
		return nil, err
	}
	user, err := r.user.Login(ctx, input)
	if err != nil {
		return nil, err
	}
	token, err := user.GenToken()
	if err != nil {
		return nil, err
	}
	// TODO: add tokens to redis cache to avoid inecesaries queries if the user is login
	return token, nil
}

// Otp is the resolver for the otp field.
func (r *mutationResolver) Otp(ctx context.Context, email string) (*cubawheeler.Response, error) {
	// TODO: generate a new otp and send the email
	var rsp = cubawheeler.Response{
		Success: true,
		Code:    http.StatusOK,
	}
	otp, err := r.otp.Create(ctx, email)
	if err != nil {
		rsp.Success = false
		rsp.Code = http.StatusBadRequest
		rsp.Message = err.Error()
	} else {
		rsp.Message = otp
	}
	return &rsp, nil
}

// UpdateProfile is the resolver for the updateProfile field.
func (r *mutationResolver) UpdateProfile(ctx context.Context, profile cubawheeler.UpdateProfile) (*cubawheeler.Response, error) {
	var rsp = cubawheeler.Response{
		Success: true,
		Code:    http.StatusOK,
	}
	if err := r.user.UpdateProfile(ctx, &cubawheeler.UpdateProfile{
		Name:     profile.Name,
		LastName: profile.LastName,
		Dob:      profile.Dob,
		Phone:    profile.Phone,
		Photo:    profile.Photo,
		Gender:   profile.Gender,
		Dni:      profile.Dni,
	}); err != nil {
		rsp.Success = false
		rsp.Message = err.Error()
		rsp.Code = http.StatusBadRequest
		return nil, err
	}
	return &rsp, nil
}

// UpdateOrder is the resolver for the updateTrip field.
// func (r *mutationResolver) UpdateOrder(ctx context.Context, update *cubawheeler.UpdateOrder) (*cubawheeler.Order, error) {
// 	return r.order.Update(ctx, update)
// }

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
	}
	if err := r.user.AddDevice(ctx, device); err != nil {
		rsp.Message = err.Error()
		rsp.Code = http.StatusBadRequest
		rsp.Success = false
	}
	return &rsp, nil
}
