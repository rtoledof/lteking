package graph

import (
	"context"

	"cubawheeler.io/pkg/cubawheeler"
)

var _ OrderResolver = &orderResolver{}

type orderResolver struct{ *Resolver }

func (r *orderResolver) Price(ctx context.Context, obj *cubawheeler.Order) (int, error) {
	return int(obj.Price), nil
}

// / Rider is the resolver for the rider field.
func (r *orderResolver) Rider(ctx context.Context, obj *cubawheeler.Order) (*cubawheeler.User, error) {
	return r.user.FindByID(ctx, obj.Rider)
}

// Driver is the resolver for the driver field.
func (r *orderResolver) Driver(ctx context.Context, obj *cubawheeler.Order) (*cubawheeler.User, error) {
	return r.user.FindByID(ctx, obj.Driver)
}

// History is the resolver for the history field.
func (r *orderResolver) History(ctx context.Context, obj *cubawheeler.Order) ([]*cubawheeler.Point, error) {
	return obj.History, nil
}

// Coupon is the resolver for the coupon field.
func (r *orderResolver) Coupon(ctx context.Context, obj *cubawheeler.Order) (*cubawheeler.Coupon, error) {
	return r.coupon.FindByID(ctx, obj.Coupon)
}

// Review is the resolver for the review field.
func (r *orderResolver) Review(ctx context.Context, obj *cubawheeler.Order) ([]*cubawheeler.Review, error) {
	reviews, _, err := r.review.FindAll(ctx, cubawheeler.ReviewFilter{
		From: obj.Rider,
	})
	return reviews, err
}

// Items is the resolver for the items field.
func (r *orderResolver) Items(ctx context.Context, obj *cubawheeler.Order) ([]*cubawheeler.OrderItem, error) {
	return obj.Items, nil
}

type orderItemResolver struct{ *Resolver }

// Seconds is the resolver for the seconds field.
func (r *orderItemResolver) Seconds(ctx context.Context, obj *cubawheeler.OrderItem) (int, error) {
	return int(obj.Seconds), nil
}

// M is the resolver for the m field.
func (r *orderItemResolver) M(ctx context.Context, obj *cubawheeler.OrderItem) (float64, error) {
	return float64(obj.Meters), nil
}

type createOrderRequestResolver struct{ *Resolver }

// Riders is the resolver for the riders field.
func (r *createOrderRequestResolver) Riders(ctx context.Context, obj *cubawheeler.CreateOrderRequest, data *int) error {
	obj.Riders = data
	return nil
}

// Baggages is the resolver for the baggages field.
func (r *createOrderRequestResolver) Baggages(ctx context.Context, obj *cubawheeler.CreateOrderRequest, data *int) error {
	obj.Baggages = data
	return nil
}
