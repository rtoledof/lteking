package seed

import "cubawheeler.io/pkg/mongo"

type Seed interface {
	Up() error
	Down() error
}

type seed struct {
	seeders []Seed
}

func NewSeed(db *mongo.DB) Seed {
	return &seed{
		seeders: []Seed{
			NewApplication(db),
			NewPlan(db),
			NewRate(db),
			NewVehicleCategoryRate(db),
		},
	}
}

func (s *seed) Up() error {
	for _, v := range s.seeders {
		if err := v.Up(); err != nil {
			return err
		}
	}
	return nil
}

func (s *seed) Down() error {
	panic("implement me")
}
