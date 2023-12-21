package graph

import (
	"context"

	"cubawheeler.io/pkg/cubawheeler"
)

type directionResponseResolver struct{ *Resolver }

// Distance is the resolver for the distance field.
func (r *directionResponseResolver) Distance(ctx context.Context, obj *cubawheeler.DirectionResponse) (int, error) {
	return int(obj.Distance), nil
}

// Duration is the resolver for the duration field.
func (r *directionResponseResolver) Duration(ctx context.Context, obj *cubawheeler.DirectionResponse) (int, error) {
	return int(obj.Duration), nil
}

// Waitpoints is the resolver for the waitpoints field.
func (r *directionResponseResolver) Waitpoints(ctx context.Context, obj *cubawheeler.DirectionResponse) ([]*cubawheeler.WaitPoint, error) {
	return obj.WayPoints, nil
}
