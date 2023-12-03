package graph

import (
	"context"
	"errors"
	"fmt"

	"cubawheeler.io/pkg/cubawheeler"
)

type mutationResolver struct{ *Resolver }

// Login is the resolver for the login field.
func (r *mutationResolver) Login(ctx context.Context, input cubawheeler.LoginRequest) (*cubawheeler.Token, error) {
	user, err := r.user.FindByEmail(ctx, input.Email)
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

// Register is the resolver for the register field.
func (r *mutationResolver) Register(ctx context.Context, email string, otp string) (*cubawheeler.Token, error) {
	_, err := r.user.FindByEmail(ctx, email)
	if err == nil {
		return nil, errors.New("email aready in use")
	}
	user := &cubawheeler.User{
		ID:    cubawheeler.NewID().String(),
		Email: email,
		Code:  cubawheeler.NewReferalCode(),
	}
	if err := r.user.CreateUser(ctx, user); err != nil {
		return nil, err
	}
	token, err := user.GenToken()
	if err != nil {
		return nil, err
	}
	return token, nil
}

// Otp is the resolver for the otp field.
func (r *mutationResolver) Otp(ctx context.Context, email string) (string, error) {
	return "000000", nil
}

// RequestTrip is the resolver for the requestTrip field.
func (r *mutationResolver) RequestTrip(ctx context.Context, input cubawheeler.RequestTrip) (*cubawheeler.Trip, error) {
	usr := cubawheeler.UserFromContext(ctx)
	if usr == nil {
		return nil, errors.New("invalid token provided")
	}
	trip := cubawheeler.Trip{
		PickUp: &cubawheeler.Location{
			Lat:  input.PickUp.Lat,
			Long: input.PickUp.Long,
		},
		DropOff: &cubawheeler.Location{
			Lat:  input.DropOff.Lat,
			Long: input.DropOff.Long,
		},
		Rider:  usr.ID,
		Status: cubawheeler.TripStatusNew,
	}

	for _, l := range input.Route {
		trip.Route = append(trip.History, cubawheeler.Location{
			Lat:  l.Lat,
			Long: l.Long,
		})
	}

	panic(fmt.Errorf("not implemented: RequestTrip - requestTrip"))
}

// UpdateProfile is the resolver for the updateProfile field.
func (r *mutationResolver) UpdateProfile(ctx context.Context, profile cubawheeler.UpdateProfile) (*cubawheeler.Profile, error) {
	return r.profile.Update(ctx, &cubawheeler.ProfileRequest{
		Name:     profile.Name,
		LastName: profile.LastName,
		DOB:      profile.DOB,
		Phone:    profile.Phone,
		Photo:    profile.Photo,
		Gender:   profile.Gender,
		Licence:  profile.Licence,
		Dni:      profile.Dni,
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
