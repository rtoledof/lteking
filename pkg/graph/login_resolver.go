package graph

import (
	"context"

	"cubawheeler.io/pkg/cubawheeler"
)

type loginRequestResolver struct{ *Resolver }

// GrantType is the resolver for the grant_type field.
func (r *loginRequestResolver) GrantType(ctx context.Context, obj *cubawheeler.LoginRequest, data string) error {
	var grantType cubawheeler.GrantType
	if err := grantType.UnmarshalGQL(data); err != nil {
		return err
	}
	obj.GrantType = grantType
	return nil
}

// ClientID is the resolver for the client_id field.
func (r *loginRequestResolver) ClientID(ctx context.Context, obj *cubawheeler.LoginRequest, data *string) error {
	obj.Client = *data
	return nil
}

// ClientSecret is the resolver for the client_secret field.
func (r *loginRequestResolver) ClientSecret(ctx context.Context, obj *cubawheeler.LoginRequest, data *string) error {
	obj.Secret = *data
	return nil
}
