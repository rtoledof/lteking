package seed

import (
	"errors"

	"order.io/pkg/mongo"
	"order.io/pkg/order"
)

type VehicleCategoryRate struct {
	service  order.VehicleCategoryRateService
	features []order.VehicleCategoryRateRequest
}

func NewVehicleCategoryRate(db *mongo.DB) *VehicleCategoryRate {
	return &VehicleCategoryRate{
		service: mongo.NewVehicleCategoryRateService(db),
		features: []order.VehicleCategoryRateRequest{
			{
				ID:       order.NewID().String(),
				Category: order.VehicleCategoryX,
				Factor:   1.0,
			},
			{
				ID:       order.NewID().String(),
				Category: order.VehicleCategoryXl,
				Factor:   1.2,
			},
			{
				ID:       order.NewID().String(),
				Category: order.VehicleCategoryConfort,
				Factor:   1.5,
			},
			{
				ID:       order.NewID().String(),
				Category: order.VehicleCategoryGreen,
				Factor:   1.0,
			},
			{
				ID:       order.NewID().String(),
				Category: order.VehicleCategoryPets,
				Factor:   1.8,
			},
			{
				ID:       order.NewID().String(),
				Category: order.VehicleCategoryPackage,
				Factor:   1.0,
			},
			{
				ID:       order.NewID().String(),
				Category: order.VehicleCategoryPriority,
				Factor:   2.0,
			},
		},
	}
}

func (r *VehicleCategoryRate) Up() error {
	ctx := prepateContext()
	for _, feature := range r.features {
		_, err := r.service.FindByCategory(ctx, feature.Category)
		if err != nil && errors.Is(err, order.ErrNotFound) {
			if _, err := r.service.Create(ctx, &feature); err != nil {
				return err
			}
		}
	}
	return nil
}

func (r *VehicleCategoryRate) Down() error {
	return nil
}
