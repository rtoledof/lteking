package seed

import (
	"context"
	"errors"

	"cubawheeler.io/pkg/cubawheeler"
	"cubawheeler.io/pkg/mongo"
)

type Rate struct {
	service  cubawheeler.RateService
	features []cubawheeler.RateRequest
}

func NewRate(db *mongo.DB) *Rate {
	return &Rate{
		service: mongo.NewRateService(db),
		features: []cubawheeler.RateRequest{
			{
				ID:                cubawheeler.NewID().String(),
				Code:              "R1",
				BasePrice:         350,
				PricePerKm:        func() *int { p := 17500; return &p }(),
				PricePerPassenger: func() *int { p := 10000; return &p }(),
				PricePerBaggage:   func() *int { p := 5000; return &p }(),
				StartTime:         "05:30",
				EndTime:           "08:30",
				MaxKm:             func() *int { p := 12000; return &p }(),
			},
			{
				ID:                cubawheeler.NewID().String(),
				Code:              "R2",
				BasePrice:         1000,
				PricePerKm:        func() *int { p := 12500; return &p }(),
				PricePerPassenger: func() *int { p := 10000; return &p }(),
				PricePerBaggage:   func() *int { p := 5000; return &p }(),
				StartTime:         "08:30",
				EndTime:           "15:59",
				MaxKm:             func() *int { p := 12000; return &p }(),
			},
			{
				ID:                cubawheeler.NewID().String(),
				Code:              "R3",
				BasePrice:         1000,
				PricePerKm:        func() *int { p := 17500; return &p }(),
				PricePerPassenger: func() *int { p := 10000; return &p }(),
				PricePerBaggage:   func() *int { p := 5000; return &p }(),
				StartTime:         "16:00",
				EndTime:           "19:59",
				MaxKm:             func() *int { p := 12000; return &p }(),
			},
			{
				ID:                cubawheeler.NewID().String(),
				Code:              "R4",
				BasePrice:         1000,
				PricePerKm:        func() *int { p := 15000; return &p }(),
				PricePerPassenger: func() *int { p := 10000; return &p }(),
				PricePerBaggage:   func() *int { p := 5000; return &p }(),
				StartTime:         "20:00",
				EndTime:           "23:30",
				MaxKm:             func() *int { p := 120; return &p }(),
			},
			{
				ID:                cubawheeler.NewID().String(),
				Code:              "R5",
				BasePrice:         1000,
				PricePerKm:        func() *int { p := 20000; return &p }(),
				PricePerPassenger: func() *int { p := 10000; return &p }(),
				PricePerBaggage:   func() *int { p := 5000; return &p }(),
				StartTime:         "23:30",
				EndTime:           "05:30",
				MaxKm:             func() *int { p := 12000; return &p }(),
			},
		},
	}
}

func (s *Rate) Up() error {
	usr := cubawheeler.User{
		Role: cubawheeler.RoleAdmin,
	}
	ctx := cubawheeler.NewContextWithUser(context.TODO(), &usr)
	for _, r := range s.features {
		_, err := s.service.FindByCode(ctx, r.Code)
		if err != nil && errors.Is(err, cubawheeler.ErrNotFound) {
			if _, err := s.service.Create(ctx, r); err != nil {
				return err
			}
		}

	}
	return nil
}

func (s *Rate) Down() error {
	panic("implement me")
}
