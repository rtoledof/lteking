package graph

import (
	"context"

	"cubawheeler.io/pkg/cubawheeler"
)

type tripResolver struct{ *Resolver }

// ID is the resolver for the id field.
func (r *tripResolver) ID(ctx context.Context, obj *cubawheeler.Trip) (string, error) {
	return obj.ID, nil
}

// Driver is the resolver for the driver field.
func (r *tripResolver) Driver(ctx context.Context, obj *cubawheeler.Trip) (*cubawheeler.User, error) {
	return r.user.FindByID(ctx, obj.Driver)
}

// Rider is the resolver for the rider field.
func (r *tripResolver) Rider(ctx context.Context, obj *cubawheeler.Trip) (*cubawheeler.User, error) {
	return r.user.FindByID(ctx, obj.Rider)
}

// StatusHistory is the resolver for the status_history field.
func (r *tripResolver) StatusHistory(ctx context.Context, obj *cubawheeler.Trip) ([]*cubawheeler.TripStatusHistory, error) {
	return obj.StatusHistory, nil
}

// Coupon is the resolver for the coupon field.
func (r *tripResolver) Coupon(ctx context.Context, obj *cubawheeler.Trip) (*cubawheeler.Coupon, error) {
	return r.coupon.FindByID(ctx, obj.Coupon)
}

// Review is the resolver for the review field.
func (r *tripResolver) Review(ctx context.Context, obj *cubawheeler.Trip) (*cubawheeler.Review, error) {
	return r.review.FindById(ctx, obj.Review)
}

type requestTripResolver struct{ *Resolver }

// PickUpLat is the resolver for the pick_up_lat field.
func (r *requestTripResolver) PickUpLat(ctx context.Context, obj *cubawheeler.RequestTrip, data float64) error {
	return nil
}

// PickUpLong is the resolver for the pick_up_long field.
func (r *requestTripResolver) PickUpLong(ctx context.Context, obj *cubawheeler.RequestTrip, data float64) error {
	return nil
}

// DropOffLat is the resolver for the drop_off_lat field.
func (r *requestTripResolver) DropOffLat(ctx context.Context, obj *cubawheeler.RequestTrip, data float64) error {
	return nil
}

// DropOffLong is the resolver for the drop_off_long field.
func (r *requestTripResolver) DropOffLong(ctx context.Context, obj *cubawheeler.RequestTrip, data float64) error {
	return nil
}

// Route is the resolver for the route field.
func (r *requestTripResolver) Route(ctx context.Context, obj *cubawheeler.RequestTrip, data []*cubawheeler.LocationInput) error {
	return nil
}
