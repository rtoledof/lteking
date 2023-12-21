package graph

import (
	"context"

	"cubawheeler.io/pkg/cubawheeler"
)

var _ OrderResolver = &orderResolver{}

type orderResolver struct{ *Resolver }

// Items implements OrderResolver.
func (*orderResolver) Items(ctx context.Context, obj *cubawheeler.Order) ([]*cubawheeler.OrderItem, error) {
	panic("unimplemented")
}

// Cost implements OrderResolver.
func (*orderResolver) Cost(ctx context.Context, obj *cubawheeler.Order) ([]*cubawheeler.CategoryPrice, error) {
	return obj.CategoryPrice, nil
}

// SelectedCost implements OrderResolver.
func (*orderResolver) SelectedCost(ctx context.Context, obj *cubawheeler.Order) (*cubawheeler.CategoryPrice, error) {
	return &obj.SelectedCategory, nil
}

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
