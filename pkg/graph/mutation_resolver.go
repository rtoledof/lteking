package graph

import (
	"context"
	"fmt"

	"cubawheeler.io/pkg/cubawheeler"
)

type mutationResolver struct{ *Resolver }

// Login is the resolver for the login field.
func (r *mutationResolver) Login(ctx context.Context, input cubawheeler.LoginRequest) (*cubawheeler.Token, error) {
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
	return "000000", nil
}

// RequestTrip is the resolver for the requestTrip field.
func (r *mutationResolver) RequestTrip(ctx context.Context, input cubawheeler.RequestTrip) (*cubawheeler.Trip, error) {
	return r.trip.Create(ctx, &input)
}

// UpdateProfile is the resolver for the updateProfile field.
func (r *mutationResolver) UpdateProfile(ctx context.Context, profile cubawheeler.UpdateProfile) (*cubawheeler.Profile, error) {
	return r.profile.Update(ctx, &cubawheeler.ProfileRequest{
		Name:     profile.Name,
		LastName: profile.LastName,
		Dob:      profile.DOB,
		Phone:    profile.Phone,
		Photo:    profile.Photo,
		Gender:   profile.Gender,
		// Licence:  profile.Licence,
		Dni: profile.Dni,
	})
}

// UpdateTrip is the resolver for the updateTrip field.
func (r *mutationResolver) UpdateTrip(ctx context.Context, update *cubawheeler.UpdateTrip) (*cubawheeler.Trip, error) {
	return r.trip.Update(ctx, update)
}

// AddFavoritePlace is the resolver for the addFavoritePlace field.
func (r *mutationResolver) AddFavoritePlace(ctx context.Context, input cubawheeler.AddPlace) (*cubawheeler.Location, error) {
	return r.user.AddFavoritePlace(ctx, input)
}

// FavoritePlaces is the resolver for the favoritePlaces field.
func (r *mutationResolver) FavoritePlaces(ctx context.Context) ([]*cubawheeler.Location, error) {
	panic(fmt.Errorf("not implemented: FavoritePlaces - favoritePlaces"))
}

// AddFavoriteVehicle is the resolver for the addFavoriteVehicle field.
func (r *mutationResolver) AddFavoriteVehicle(ctx context.Context, plate *string) (*cubawheeler.Vehicle, error) {
	panic(fmt.Errorf("not implemented: AddFavoriteVehicle - addFavoriteVehicle"))
}

// FavoriteVehicles is the resolver for the favoriteVehicles field.
func (r *mutationResolver) FavoriteVehicles(ctx context.Context) ([]*cubawheeler.Vehicle, error) {
	panic(fmt.Errorf("not implemented: FavoriteVehicles - favoriteVehicles"))
}

// UpdatePlace is the resolver for the updatePlace field.
func (r *mutationResolver) UpdatePlace(ctx context.Context, input *cubawheeler.UpdatePlace) (*cubawheeler.Location, error) {
	panic(fmt.Errorf("not implemented: UpdatePlace - updatePlace"))
}

// FindVehicle is the resolver for the findVehicle field.
func (r *mutationResolver) FindVehicle(ctx context.Context, vehicle string) (*cubawheeler.Vehicle, error) {
	panic(fmt.Errorf("not implemented: FindVehicle - findVehicle"))
}

// UpdateVehicle is the resolver for the updateVehicle field.
func (r *mutationResolver) UpdateVehicle(ctx context.Context, input *cubawheeler.UpdateVehicle) (*cubawheeler.Vehicle, error) {
	panic(fmt.Errorf("not implemented: UpdateVehicle - updateVehicle"))
}

// CancelTrip is the resolver for the cancelTrip field.
func (r *mutationResolver) CancelTrip(ctx context.Context, trip string) (*cubawheeler.Trip, error) {
	panic(fmt.Errorf("not implemented: CancelTrip - cancelTrip"))
}

// AddRate is the resolver for the addRate field.
func (r *mutationResolver) AddRate(ctx context.Context, input cubawheeler.RateRequest) (*cubawheeler.Rate, error) {
	panic(fmt.Errorf("not implemented: AddRate - addRate"))
}

// UpdateRate is the resolver for the updateRate field.
func (r *mutationResolver) UpdateRate(ctx context.Context, input cubawheeler.RateRequest) (*cubawheeler.Rate, error) {
	panic(fmt.Errorf("not implemented: UpdateRate - updateRate"))
}

func (r *mutationResolver) ChangePin(ctx context.Context, old *string, pin string) (*cubawheeler.Profile, error) {
	return r.profile.ChangePin(ctx, old, pin)
}

// AcceptTrip is the resolver for the acceptTrip field.
func (r *mutationResolver) AcceptTrip(ctx context.Context, trip string) (*cubawheeler.Trip, error) {
	panic(fmt.Errorf("not implemented: AcceptTrip - acceptTrip"))
}

// CreateApplication is the resolver for the createApplication field.
func (r *mutationResolver) CreateApplication(ctx context.Context, input cubawheeler.ApplicationRequest) (*cubawheeler.Application, error) {
	panic(fmt.Errorf("not implemented: CreateApplication - createApplication"))
}

// UpdateApplicationCredentials is the resolver for the updateApplicationCredentials field.
func (r *mutationResolver) UpdateApplicationCredentials(ctx context.Context, application string) (*cubawheeler.Application, error) {
	panic(fmt.Errorf("not implemented: UpdateApplicationCredentials - updateApplicationCredentials"))
}
