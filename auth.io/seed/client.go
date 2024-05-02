package seed

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"auth.io/models"
	"auth.io/mongo"
	"github.com/go-chi/jwtauth"
	"github.com/lestrrat-go/jwx/jwt"
)

var _ Seeder = &Client{}

type Client struct {
	service  models.ClientService
	features []models.Client
}

func NewClient(db *mongo.DB) *Client {
	return &Client{
		service: mongo.NewClientService(db),
		features: []models.Client{
			{
				Name: "Rider",
				Type: models.ClientTypeRider,
			},
			{
				Name: "Driver",
				Type: models.ClientTypeDriver,
			},
		},
	}
}

func (s *Client) Up() error {
	ctx := prepateContext()
	for _, v := range s.features {
		result, _, err := s.service.FindClients(ctx, models.ClientFilter{
			Name: []string{v.Name},
		})
		if err != nil && errors.Is(err, models.ErrNotFound) || result == nil {
			err := s.service.Create(ctx, &v)
			if err != nil {
				return nil
			}
			fmt.Println(v)
		}
	}
	return nil
}

func (s *Client) Down() error {
	//TODO implement me
	panic("implement me")
}

func prepateContext(roles ...models.Role) context.Context {

	ctx := context.Background()

	token := jwt.New()
	token.Set("id", models.NewID().String())
	user := models.User{
		ID:   models.NewID().String(),
		Role: models.RoleAdmin,
	}
	if roles != nil {
		user.Role = roles[0]
	}
	ctx = models.NewContextWithUser(ctx, &user)
	userData, _ := json.Marshal(user)
	token.Set("user", userData)

	return jwtauth.NewContext(ctx, token, nil)
}
