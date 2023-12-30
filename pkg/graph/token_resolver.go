package graph

import (
	"context"

	"cubawheeler.io/pkg/cubawheeler"
)

type tokenResolver struct{ *Resolver }

// ExpiryAt is the resolver for the expiry_at field.
func (r *tokenResolver) ExpiryAt(ctx context.Context, obj *cubawheeler.Token) (int, error) {
	return int(obj.AccessTokenExpiresIn), nil
}
