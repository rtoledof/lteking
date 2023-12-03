package graph

import (
	"context"

	"cubawheeler.io/pkg/cubawheeler"
)

type profileResolver struct{ *Resolver }

// Dob is the resolver for the dob field.
func (r *profileResolver) Dob(ctx context.Context, obj *cubawheeler.Profile) (*string, error) {
	return &obj.DOB, nil
}

// User is the resolver for the user field.
func (r *profileResolver) User(ctx context.Context, obj *cubawheeler.Profile) (*cubawheeler.User, error) {
	user, err := r.user.FindByID(ctx, obj.UserId)
	if err != nil {
		return nil, err
	}
	return user, nil
}

type updateProfileResolver struct{ *Resolver }

// LastName is the resolver for the last_name field.
func (r *updateProfileResolver) LastName(ctx context.Context, obj *cubawheeler.UpdateProfile, data *string) error {
	obj.LastName = data
	return nil
}

// Phone is the resolver for the phone field.
func (r *updateProfileResolver) Phone(ctx context.Context, obj *cubawheeler.UpdateProfile, data *string) error {
	obj.Phone = data
	return nil
}

// License is the resolver for the license field.
func (r *updateProfileResolver) License(ctx context.Context, obj *cubawheeler.UpdateProfile, data *string) error {
	obj.Licence = data
	return nil
}

// Dni is the resolver for the dni field.
func (r *updateProfileResolver) Dni(ctx context.Context, obj *cubawheeler.UpdateProfile, data *string) error {
	obj.Dni = data
	return nil
}
