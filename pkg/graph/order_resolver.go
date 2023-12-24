package graph

import (
	"context"
	"fmt"

	"cubawheeler.io/pkg/cubawheeler"
	"cubawheeler.io/pkg/graph/model"
)

var _ OrderResolver = &orderResolver{}

type orderResolver struct{ *Resolver }

// Items implements OrderResolver.
func (*orderResolver) Items(ctx context.Context, obj *cubawheeler.Order) ([]*cubawheeler.OrderItem, error) {
	panic("unimplemented")
}

// Price implements OrderResolver.
func (*orderResolver) Price(ctx context.Context, obj *cubawheeler.Order) (*model.Amount, error) {
	return &model.Amount{
		Amount:   int(obj.Price.Amount),
		Currency: obj.Price.Currency.String(),
	}, nil
}

// Cost implements OrderResolver.
func (*orderResolver) Cost(ctx context.Context, obj *cubawheeler.Order) ([]*cubawheeler.CategoryPrice, error) {
	return obj.CategoryPrice, nil
}

// SelectedCost implements OrderResolver.
func (*orderResolver) SelectedCost(ctx context.Context, obj *cubawheeler.Order) (*cubawheeler.CategoryPrice, error) {
	return &obj.SelectedCategory, nil
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
	if obj.Riders == nil {
		riders := 1
		obj.Riders = &riders
	}
	return nil
}

// Baggages is the resolver for the baggages field.
func (r *createOrderRequestResolver) Baggages(ctx context.Context, obj *cubawheeler.CreateOrderRequest, data *int) error {
	obj.Baggages = data
	return nil
}

type confirmOrderResolver struct{ *Resolver }

// Order is the resolver for the order field.
func (r *confirmOrderResolver) Order(ctx context.Context, obj *cubawheeler.ConfirmOrder, data string) error {
	obj.OrderID = data
	return nil
}

var _ OrderItemResolver = &orderItemResolver{}

type orderItemResolver struct{ *Resolver }

// DropOff implements OrderItemResolver.
func (*orderItemResolver) DropOff(ctx context.Context, obj *cubawheeler.OrderItem) (*cubawheeler.Point, error) {
	if len(obj.Points) == 0 {
		return nil, fmt.Errorf("no points")
	}
	return obj.Points[0], nil
}

// PickUp is the resolver for the pick_up field.
func (r *orderItemResolver) PickUp(ctx context.Context, obj *cubawheeler.OrderItem) (*cubawheeler.Point, error) {
	if len(obj.Points) < 2 {
		return nil, fmt.Errorf("no points")
	}
	return obj.Points[len(obj.Points)-1], nil
}
