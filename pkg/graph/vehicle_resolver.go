package graph

import (
	"context"

	"cubawheeler.io/pkg/cubawheeler"
	"github.com/99designs/gqlgen/graphql"
)

type vehicleResolver struct{ *Resolver }

// ID is the resolver for the id field.
func (r *vehicleResolver) ID(ctx context.Context, obj *cubawheeler.Vehicle) (string, error) {
	return obj.ID, nil
}

// Model is the resolver for the model field.
func (r *vehicleResolver) Model(ctx context.Context, obj *cubawheeler.Vehicle) (string, error) {
	return obj.CarModel, nil
}

// Status is the resolver for the status field.
func (r *vehicleResolver) Status(ctx context.Context, obj *cubawheeler.Vehicle) (cubawheeler.VehicleStatus, error) {
	return obj.Status, nil
}

type updateVehicleResolver struct{ *Resolver }

// ID is the resolver for the id field.
func (r *updateVehicleResolver) ID(ctx context.Context, obj *cubawheeler.UpdateVehicle, data string) error {
	return nil
}

// Palte is the resolver for the palte field.
func (r *updateVehicleResolver) Palte(ctx context.Context, obj *cubawheeler.UpdateVehicle, data *string) error {
	return nil
}

// VehicleType is the resolver for the vehicleType field.
func (r *updateVehicleResolver) VehicleType(ctx context.Context, obj *cubawheeler.UpdateVehicle, data *cubawheeler.VehicleType) error {
	return nil
}

func (r *updateVehicleResolver) Pictures(ctx context.Context, obj *cubawheeler.UpdateVehicle, data []*graphql.Upload) error {
	panic("implement me")
}
