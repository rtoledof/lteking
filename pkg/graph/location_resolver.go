package graph

import (
	"context"

	"cubawheeler.io/pkg/cubawheeler"
	"github.com/99designs/gqlgen/graphql"
)

type updatePlaceResolver struct{ *Resolver }

// Palte implements UpdateVehicleResolver.
func (*updatePlaceResolver) Palte(ctx context.Context, obj *cubawheeler.UpdateVehicle, data *string) error {
	panic("unimplemented")
}

// Pictures implements UpdateVehicleResolver.
func (*updatePlaceResolver) Pictures(ctx context.Context, obj *cubawheeler.UpdateVehicle, data []*graphql.Upload) error {
	panic("unimplemented")
}

// VehicleType implements UpdateVehicleResolver.
func (*updatePlaceResolver) VehicleType(ctx context.Context, obj *cubawheeler.UpdateVehicle, data *cubawheeler.VehicleType) error {
	panic("unimplemented")
}

// Location is the resolver for the location field.
func (r *updatePlaceResolver) Location(ctx context.Context, obj *cubawheeler.UpdatePlace, data *cubawheeler.LocationInput) error {
	obj.Location = data
	return nil
}
