package graph

import (
	"context"

	"cubawheeler.io/pkg/cubawheeler"
)

type stepResolver struct{ *Resolver }

// Distance is the resolver for the distance field.
func (r *stepResolver) Distance(ctx context.Context, obj *cubawheeler.Step) (int, error) {
	return int(obj.Distance), nil
}

// Duration is the resolver for the duration field.
func (r *stepResolver) Duration(ctx context.Context, obj *cubawheeler.Step) (int, error) {
	return int(obj.Duration), nil
}
