package seed

import (
	"context"

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
				PricePerKm:        func() *int { p := 175; return &p }(),
				PricePerPassenger: func() *int { p := 100; return &p }(),
				PricePerBaggage:   func() *int { p := 50; return &p }(),
				StartTime:         "05:30",
				EndTime:           "08:30",
				MaxKm:             func() *int { p := 120; return &p }(),
			},
			{
				ID:                cubawheeler.NewID().String(),
				Code:              "R2",
				BasePrice:         1000,
				PricePerKm:        func() *int { p := 125; return &p }(),
				PricePerPassenger: func() *int { p := 100; return &p }(),
				PricePerBaggage:   func() *int { p := 50; return &p }(),
				StartTime:         "08:30",
				EndTime:           "15:59",
				MaxKm:             func() *int { p := 120; return &p }(),
			},
			{
				ID:                cubawheeler.NewID().String(),
				Code:              "R3",
				BasePrice:         1000,
				PricePerKm:        func() *int { p := 175; return &p }(),
				PricePerPassenger: func() *int { p := 100; return &p }(),
				PricePerBaggage:   func() *int { p := 50; return &p }(),
				StartTime:         "16:00",
				EndTime:           "19:59",
				MaxKm:             func() *int { p := 120; return &p }(),
			},
			{
				ID:                cubawheeler.NewID().String(),
				Code:              "R4",
				BasePrice:         1000,
				PricePerKm:        func() *int { p := 150; return &p }(),
				PricePerPassenger: func() *int { p := 100; return &p }(),
				PricePerBaggage:   func() *int { p := 50; return &p }(),
				StartTime:         "20:00",
				EndTime:           "23:30",
				MaxKm:             func() *int { p := 120; return &p }(),
			},
			{
				ID:                cubawheeler.NewID().String(),
				Code:              "R4",
				BasePrice:         1000,
				PricePerKm:        func() *int { p := 200; return &p }(),
				PricePerPassenger: func() *int { p := 100; return &p }(),
				PricePerBaggage:   func() *int { p := 50; return &p }(),
				StartTime:         "23:30",
				EndTime:           "05:30",
				MaxKm:             func() *int { p := 120; return &p }(),
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
		if _, err := s.service.Create(ctx, r); err != nil {
			return err
		}
	}
	return nil
}

func (s *Rate) Down() error {
	panic("implement me")
}
