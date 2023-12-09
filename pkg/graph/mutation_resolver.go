package graph

import (
	"context"
	"fmt"

	"cubawheeler.io/pkg/cubawheeler"
)

var _ MutationResolver = &mutationResolver{}

type mutationResolver struct{ *Resolver }

func (r *mutationResolver) CreateOrder(ctx context.Context, input []*cubawheeler.Item) (*cubawheeler.Order, error) {
	//TODO implement me
	panic("implement me")
}

func (r *mutationResolver) UpdatOrder(ctx context.Context, update *cubawheeler.UpdateOrder) (*cubawheeler.Order, error) {
	//TODO implement me
	panic("implement me")
}

func (r *mutationResolver) CancelOrder(ctx context.Context, order string) (*cubawheeler.Order, error) {
	//TODO implement me
	panic("implement me")
}

func (r *mutationResolver) AcceptOrder(ctx context.Context, order string) (*cubawheeler.Order, error) {
	//TODO implement me
	panic("implement me")
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
func (r *mutationResolver) Otp(ctx context.Context, email string) (string, error) {
	// TODO: generate a new otp and send the email
	return "", r.otp.Create(ctx, email)
}

// UpdateProfile is the resolver for the updateProfile field.
func (r *mutationResolver) UpdateProfile(ctx context.Context, profile cubawheeler.UpdateProfile) (*cubawheeler.Profile, error) {
	return r.profile.Update(ctx, &cubawheeler.UpdateProfile{
		Name:     profile.Name,
		LastName: profile.LastName,
		Dob:      profile.Dob,
		Phone:    profile.Phone,
		Photo:    profile.Photo,
		Gender:   profile.Gender,
		// Licence:  profile.Licence,
		Dni: profile.Dni,
	})
}

// UpdateTrip is the resolver for the updateTrip field.
func (r *mutationResolver) UpdateOrder(ctx context.Context, update *cubawheeler.UpdateOrder) (*cubawheeler.Order, error) {
	return r.order.Update(ctx, update)
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
	panic(fmt.Errorf("not implemented: CancelTrip - cancelTrip"))
}

// AddRate is the resolver for the addRate field.
func (r *mutationResolver) AddRate(ctx context.Context, input cubawheeler.RateRequest) (*cubawheeler.Rate, error) {
	return r.rate.Create(ctx, input)
}

// UpdateRate is the resolver for the updateRate field.
func (r *mutationResolver) UpdateRate(ctx context.Context, input cubawheeler.RateRequest) (*cubawheeler.Rate, error) {
	panic(fmt.Errorf("not implemented: UpdateRate - updateRate"))
}

func (r *mutationResolver) ChangePin(ctx context.Context, old *string, pin string) (*cubawheeler.Profile, error) {
	return r.profile.ChangePin(ctx, old, pin)
}

// AcceptTrip is the resolver for the acceptTrip field.
func (r *mutationResolver) AcceptTrip(ctx context.Context, trip string) (*cubawheeler.Order, error) {
	panic(fmt.Errorf("not implemented: AcceptTrip - acceptTrip"))
}

// CreateApplication is the resolver for the createApplication field.
func (r *mutationResolver) CreateApplication(ctx context.Context, input cubawheeler.ApplicationRequest) (*cubawheeler.Application, error) {
	return r.app.CreateApplication(ctx, input)
}

// UpdateApplicationCredentials is the resolver for the updateApplicationCredentials field.
func (r *mutationResolver) UpdateApplicationCredentials(ctx context.Context, application string) (*cubawheeler.Application, error) {
	return r.app.UpdateApplicationCredentials(ctx, application)
}
