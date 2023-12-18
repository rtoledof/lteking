package seed

import (
	"context"
	"math"

	"cubawheeler.io/pkg/cubawheeler"
	"cubawheeler.io/pkg/mongo"
)

type Plan struct {
	service  cubawheeler.PlanService
	features []cubawheeler.PlanRequest
}

func NewPlan(db *mongo.DB) *Plan {
	return &Plan{
		service: mongo.NewPlanService(db),
		features: []cubawheeler.PlanRequest{
			{
				Name:       func() *string { s := "Basic"; return &s }(),
				TotalTrips: func() *int { t := 80; return &t }(),
				Price:      func() *int { p := 1500; return &p }(),
				Interval:   func() *cubawheeler.Interval { i := cubawheeler.IntervalDay; return &i }(),
				Code:       func() *string { c := "B1500"; return &c }(),
			},
			{
				Name:       func() *string { s := "Platinum"; return &s }(),
				TotalTrips: func() *int { t := 160; return &t }(),
				Price:      func() *int { p := 2500; return &p }(),
				Interval:   func() *cubawheeler.Interval { i := cubawheeler.IntervalWeek; return &i }(),
				Code:       func() *string { c := "P2500"; return &c }(),
			},
			{
				Name:       func() *string { s := "Gold"; return &s }(),
				TotalTrips: func() *int { t := 300; return &t }(),
				Price:      func() *int { p := 3500; return &p }(),
				Interval:   func() *cubawheeler.Interval { i := cubawheeler.IntervalMonth; return &i }(),
				Code:       func() *string { c := "G3500"; return &c }(),
			},
			{
				Name:       func() *string { s := "Plus"; return &s }(),
				TotalTrips: func() *int { t := 1000; return &t }(),
				Price:      func() *int { p := 7000; return &p }(),
				Interval:   func() *cubawheeler.Interval { i := cubawheeler.IntervalMonth; return &i }(),
				Code:       func() *string { c := "PL7000"; return &c }(),
			},
			{
				Name:       func() *string { s := "Vip"; return &s }(),
				TotalTrips: func() *int { t := math.MaxInt; return &t }(),
				Price:      func() *int { p := 10000; return &p }(),
				Interval:   func() *cubawheeler.Interval { i := cubawheeler.IntervalMonth; return &i }(),
				Code:       func() *string { c := "V10000"; return &c }(),
			},
		},
	}
}

func (s *Plan) Up() error {
	usr := cubawheeler.User{
		Role: cubawheeler.RoleAdmin,
	}
	ctx := cubawheeler.NewContextWithUser(context.TODO(), &usr)
	for _, v := range s.features {
		plans, _, err := s.service.FindAll(ctx, &cubawheeler.PlanFilter{Name: *v.Name, Limit: 1})
		if err != nil || len(plans) == 0 {
			_, err := s.service.Create(ctx, &v)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *Plan) Down() error {
	//TODO implement me
	panic("implement me")
}
