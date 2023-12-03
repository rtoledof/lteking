package graph

import (
	"context"

	"cubawheeler.io/pkg/cubawheeler"
)

type addPlaceResolver struct{ *Resolver }

// Lat is the resolver for the lat field.
func (r *addPlaceResolver) Lat(ctx context.Context, obj *cubawheeler.AddPlace, data float64) error {
	return nil
}

// Long is the resolver for the long field.
func (r *addPlaceResolver) Long(ctx context.Context, obj *cubawheeler.AddPlace, data float64) error {
	return nil
}
