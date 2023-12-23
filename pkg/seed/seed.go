package seed

import "cubawheeler.io/pkg/mongo"

type Seeder interface {
	Up() error
	Down() error
}

type seed struct {
	seeders []Seeder
}

func NewSeed(db *mongo.DB) Seeder {
	return &seed{
		seeders: []Seeder{
			NewApplication(db),
			NewPlan(db),
			NewRate(db),
			NewVehicleCategoryRate(db),
			NewWallet(db),
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
