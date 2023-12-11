package graph

import (
	"context"

	"cubawheeler.io/pkg/cubawheeler"
	"github.com/99designs/gqlgen/graphql"
)

type updatePlaceResolver struct{ *Resolver }

func (r *updatePlaceResolver) Location(ctx context.Context, obj *cubawheeler.UpdatePlace, data *cubawheeler.PointInput) error {
	panic("implemet me")
}

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

type locationResolver struct{ *Resolver }

// Address is the resolver for the address field.
func (r *locationResolver) Address(ctx context.Context, obj *cubawheeler.Location) (*cubawheeler.Address, error) {
	return &obj.Address, nil
}

// GeoLocation is the resolver for the geo_location field.
func (r *locationResolver) GeoLocation(ctx context.Context, obj *cubawheeler.Location) (*cubawheeler.GeoLocation, error) {
	return &obj.Geolocation, nil
}

type geoLocationResolver struct{ *Resolver }

// Type is the resolver for the type field.
func (r *geoLocationResolver) Type(ctx context.Context, obj *cubawheeler.GeoLocation) (cubawheeler.ShapeType, error) {
	return obj.Type, nil
}
