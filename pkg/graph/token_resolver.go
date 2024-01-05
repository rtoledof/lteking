package graph

import (
	"context"

	"cubawheeler.io/pkg/cubawheeler"
)

var _ TokenResolver = &tokenResolver{}

type tokenResolver struct{ *Resolver }

// ExpiresAt implements TokenResolver.
func (*tokenResolver) ExpiresAt(ctx context.Context, obj *cubawheeler.Token) (int, error) {
	return int(obj.ExpiresAt.Unix()), nil
}

// ExpiresIn implements TokenResolver.
func (*tokenResolver) ExpiresIn(ctx context.Context, obj *cubawheeler.Token) (int, error) {
	return int(obj.ExpiresIn), nil
}

// RefreshExpireIn implements TokenResolver.
func (*tokenResolver) RefreshExpireIn(ctx context.Context, obj *cubawheeler.Token) (int, error) {
	return int(obj.RefreshTokenExpiresIn), nil
}
