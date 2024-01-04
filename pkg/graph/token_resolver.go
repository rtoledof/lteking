package graph

import (
	"context"

	"cubawheeler.io/pkg/cubawheeler"
)

var _ TokenResolver = &tokenResolver{}

type tokenResolver struct{ *Resolver }

// RefreshExpireIn implements TokenResolver.
func (*tokenResolver) RefreshExpireIn(ctx context.Context, obj *cubawheeler.Token) (int, error) {
	return int(obj.RefreshTokenExpiresIn), nil
}

// ExpiryAt is the resolver for the expiry_at field.
func (r *tokenResolver) ExpiryAt(ctx context.Context, obj *cubawheeler.Token) (int, error) {
	return int(obj.AccessTokenExpiresIn), nil
}
