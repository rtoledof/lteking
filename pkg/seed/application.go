package seed

import (
	"context"
	"errors"
	"fmt"

	"cubawheeler.io/pkg/cubawheeler"
	e "cubawheeler.io/pkg/errors"
	"cubawheeler.io/pkg/mongo"
)

var _ Seed = &Application{}

type Application struct {
	service  cubawheeler.ApplicationService
	features []cubawheeler.ApplicationRequest
}

func NewApplication(db *mongo.DB) *Application {
	return &Application{
		service: mongo.NewApplicationService(db),
		features: []cubawheeler.ApplicationRequest{
			{
				Name:   "Rider",
				Type:   cubawheeler.ApplicationTypeRider,
				Client: "rider",
				Secret: "secret",
			},
			{
				Name:   "Driver",
				Type:   cubawheeler.ApplicationTypeDriver,
				Client: "driver",
				Secret: "secret",
			},
		},
	}
}

func (s *Application) Up() error {
	usr := cubawheeler.User{
		Role: cubawheeler.RoleAdmin,
	}
	ctx := cubawheeler.NewContextWithUser(context.TODO(), &usr)
	for _, v := range s.features {
		_, err := s.service.FindByClient(ctx, v.Client)
		if err != nil && errors.Is(err, e.ErrNotFound) {
			app, err := s.service.CreateApplication(ctx, v)
			if err != nil {
				return nil
			}
			fmt.Println(app)
		}
	}
	return nil
}

func (s *Application) Down() error {
	//TODO implement me
	panic("implement me")
}
