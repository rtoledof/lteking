package graph

import (
	"context"
	"fmt"

	"cubawheeler.io/pkg/cubawheeler"
)

var _ OrderResolver = &orderResolver{}

type orderResolver struct{ *Resolver }

func (r *orderResolver) Price(ctx context.Context, obj *cubawheeler.Order) (int, error) {
	//TODO implement me
	panic("implement me")
}

/// Rider is the resolver for the rider field.
func (r *orderResolver) Rider(ctx context.Context, obj *cubawheeler.Order) (*cubawheeler.User, error) {
	panic(fmt.Errorf("not implemented: Rider - rider"))
}

// Driver is the resolver for the driver field.
func (r *orderResolver) Driver(ctx context.Context, obj *cubawheeler.Order) (*cubawheeler.User, error) {
	panic(fmt.Errorf("not implemented: Driver - driver"))
}

// History is the resolver for the history field.
func (r *orderResolver) History(ctx context.Context, obj *cubawheeler.Order) ([]*cubawheeler.Point, error) {
	panic(fmt.Errorf("not implemented: History - history"))
}

// Coupon is the resolver for the coupon field.
func (r *orderResolver) Coupon(ctx context.Context, obj *cubawheeler.Order) (*cubawheeler.Coupon, error) {
	panic(fmt.Errorf("not implemented: Coupon - coupon"))
}

// Review is the resolver for the review field.
func (r *orderResolver) Review(ctx context.Context, obj *cubawheeler.Order) ([]*cubawheeler.Review, error) {
	panic(fmt.Errorf("not implemented: Review - review"))
}

// Items is the resolver for the items field.
func (r *orderResolver) Items(ctx context.Context, obj *cubawheeler.Order) ([]*cubawheeler.OrderItem, error) {
	panic(fmt.Errorf("not implemented: Items - items"))
}

type orderItemResolver struct{ *Resolver }

// Seconds is the resolver for the seconds field.
func (r *orderItemResolver) Seconds(ctx context.Context, obj *cubawheeler.OrderItem) (int, error) {
	panic(fmt.Errorf("not implemented: Seconds - seconds"))
}

// M is the resolver for the m field.
func (r *orderItemResolver) M(ctx context.Context, obj *cubawheeler.OrderItem) (float64, error) {
	panic(fmt.Errorf("not implemented: M - m"))
}
