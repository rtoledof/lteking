package seed

import (
	"context"

	"cubawheeler.io/pkg/cubawheeler"
	"cubawheeler.io/pkg/mongo"
)

type VehicleCategoryRate struct {
	service  cubawheeler.VehicleCategoryRateService
	features []cubawheeler.VehicleCategoryRateRequest
}

func NewVehicleCategoryRate(db *mongo.DB) *VehicleCategoryRate {
	return &VehicleCategoryRate{
		service: mongo.NewVehicleCategoryRateService(db),
		features: []cubawheeler.VehicleCategoryRateRequest{
			{
				ID:       cubawheeler.NewID().String(),
				Category: cubawheeler.VehicleCategoryX,
				Factor:   1.0,
			},
			{
				ID:       cubawheeler.NewID().String(),
				Category: cubawheeler.VehicleCategoryXl,
				Factor:   1.2,
			},
			{
				ID:       cubawheeler.NewID().String(),
				Category: cubawheeler.VehicleCategoryConfort,
				Factor:   1.5,
			},
			{
				ID:       cubawheeler.NewID().String(),
				Category: cubawheeler.VehicleCategoryGreen,
				Factor:   1.0,
			},
			{
				ID:       cubawheeler.NewID().String(),
				Category: cubawheeler.VehicleCategoryPets,
				Factor:   1.8,
			},
			{
				ID:       cubawheeler.NewID().String(),
				Category: cubawheeler.VehicleCategoryPackage,
				Factor:   1.0,
			},
			{
				ID:       cubawheeler.NewID().String(),
				Category: cubawheeler.VehicleCategoryPriority,
				Factor:   2.0,
			},
		},
	}
}

func (r *VehicleCategoryRate) Up() error {
	for _, feature := range r.features {
		if _, err := r.service.Create(context.Background(), &feature); err != nil {
			return err
		}
	}
	return nil
}

func (r *VehicleCategoryRate) Down() error {
	return nil
}
