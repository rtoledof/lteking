package seed

import "cubawheeler.io/pkg/mongo"

type Seed interface {
	Up() error
	Down() error
}

type seed struct {
	app Seed
}

func NewSeed(db *mongo.DB) Seed {
	return &seed{
		app: NewApplication(db),
	}
}

func (s *seed) Up() error {
	if err := s.app.Up(); err != nil {
		return err
	}
	return nil
}

func (s *seed) Down() error {
	panic("implement me")
}
