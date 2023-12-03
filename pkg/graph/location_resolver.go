package graph

import (
	"context"

	"cubawheeler.io/pkg/cubawheeler"
)

type updatePlaceResolver struct{ *Resolver }

// Location is the resolver for the location field.
func (r *updatePlaceResolver) Location(ctx context.Context, obj *cubawheeler.UpdatePlace, data *cubawheeler.LocationInput) error {
	obj.Location = data
	return nil
}
