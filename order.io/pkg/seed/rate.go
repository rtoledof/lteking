package seed

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-chi/jwtauth"
	"github.com/lestrrat-go/jwx/jwt"
	"order.io/pkg/mongo"
	"order.io/pkg/order"
)

var _ Seeder = &Rate{}

type Rate struct {
	service  order.RateService
	features []order.RateRequest
}

func NewRate(db *mongo.DB) *Rate {
	return &Rate{
		service: mongo.NewRateService(db),
		features: []order.RateRequest{
			{
				Code:              "R1",
				BasePrice:         35000,
				PricePerKm:        17500,
				PricePerPassenger: 10000,
				PricePerBaggage:   5000,
				StartTime:         "05:30",
				EndTime:           "08:30",
				MaxKm:             12000,
			},
			{
				Code:              "R2",
				BasePrice:         10000,
				PricePerKm:        12500,
				PricePerPassenger: 10000,
				PricePerBaggage:   5000,
				StartTime:         "08:30",
				EndTime:           "15:59",
				MaxKm:             12000,
			},
			{
				Code:              "R3",
				BasePrice:         1000,
				PricePerKm:        17500,
				PricePerPassenger: 10000,
				PricePerBaggage:   5000,
				StartTime:         "16:00",
				EndTime:           "19:59",
				MaxKm:             12000,
			},
			{
				Code:              "R4",
				BasePrice:         1000,
				PricePerKm:        15000,
				PricePerPassenger: 10000,
				PricePerBaggage:   5000,
				StartTime:         "20:00",
				EndTime:           "23:30",
				MaxKm:             12000,
			},
			{
				Code:              "R5",
				BasePrice:         1000,
				PricePerKm:        20000,
				PricePerPassenger: 10000,
				PricePerBaggage:   5000,
				StartTime:         "23:30",
				EndTime:           "05:30",
				MaxKm:             12000,
			},
		},
	}
}

func (s *Rate) Up() error {
	ctx := prepateContext()
	for _, v := range s.features {
		result, err := s.service.FindByCode(ctx, v.Code)
		if err != nil && errors.Is(err, order.ErrNotFound) || result == nil {
			_, err := s.service.Create(ctx, v)
			if err != nil {
				return nil
			}
			fmt.Println(v)
		}
	}
	return nil
}

func (s *Rate) Down() error {
	//TODO implement me
	panic("implement me")
}

func prepateContext(roles ...order.Role) context.Context {

	ctx := context.Background()

	token := jwt.New()
	token.Set("id", order.NewID().String())
	user := order.User{
		ID:   order.NewID().String(),
		Role: order.RoleAdmin,
	}
	if roles != nil {
		user.Role = roles[0]
	}
	ctx = order.NewContextWithUser(ctx, user)
	token.Set("user", user.Claim())

	return jwtauth.NewContext(ctx, token, nil)
}
