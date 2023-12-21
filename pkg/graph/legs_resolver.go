package graph

import (
	"context"

	"cubawheeler.io/pkg/cubawheeler"
)

type legsResolver struct{ *Resolver }

// Distance is the resolver for the distance field.
func (r *legsResolver) Distance(ctx context.Context, obj *cubawheeler.Legs) (int, error) {
	return int(obj.Distance), nil
}

// Duration is the resolver for the duration field.
func (r *legsResolver) Duration(ctx context.Context, obj *cubawheeler.Legs) (int, error) {
	return int(obj.Duration), nil
}
