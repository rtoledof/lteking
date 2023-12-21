package graph

import (
	"context"

	"cubawheeler.io/pkg/cubawheeler"
)

type routeResolver struct{ *Resolver }

// Distance is the resolver for the distance field.
func (r *routeResolver) Distance(ctx context.Context, obj *cubawheeler.Route) (int, error) {
	return int(obj.Distance), nil
}

// Duration is the resolver for the duration field.
func (r *routeResolver) Duration(ctx context.Context, obj *cubawheeler.Route) (int, error) {
	return int(obj.Duration), nil
}
