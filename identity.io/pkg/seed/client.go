package seed

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/go-chi/jwtauth"
	"github.com/lestrrat-go/jwx/jwt"
	"identity.io/pkg/identity"
	"identity.io/pkg/mongo"
)

var _ Seeder = &Client{}

type Client struct {
	service  identity.ClientService
	features []identity.Client
}

func NewClient(db *mongo.DB) *Client {
	return &Client{
		service: mongo.NewClientService(db),
		features: []identity.Client{
			{
				Name: "Rider",
				Type: identity.ClientTypeRider,
			},
			{
				Name: "Driver",
				Type: identity.ClientTypeDriver,
			},
		},
	}
}

func (s *Client) Up() error {
	ctx := prepateContext()
	for _, v := range s.features {
		result, _, err := s.service.FindClients(ctx, identity.ClientFilter{
			Name: []string{v.Name},
		})
		if err != nil && errors.Is(err, identity.ErrNotFound) || result == nil {
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

func prepateContext(roles ...identity.Role) context.Context {

	ctx := context.Background()

	token := jwt.New()
	token.Set("id", identity.NewID().String())
	user := identity.User{
		ID:   identity.NewID().String(),
		Role: identity.RoleAdmin,
	}
	if roles != nil {
		user.Role = roles[0]
	}
	ctx = identity.NewContextWithUser(ctx, &user)
	userData, _ := json.Marshal(user)
	token.Set("user", userData)

	return jwtauth.NewContext(ctx, token, nil)
}
