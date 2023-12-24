package graph

import (
	"context"

	"cubawheeler.io/pkg/cubawheeler"
	"cubawheeler.io/pkg/graph/model"
)

var _ CouponResolver = &couponResolver{}

type couponResolver struct{ *Resolver }

// Amount is the resolver for the amount field.
func (r *couponResolver) Amount(ctx context.Context, obj *cubawheeler.Coupon) (*model.Amount, error) {
	return &model.Amount{
		Amount:   int(obj.Amount.Amount),
		Currency: obj.Amount.Currency.String(),
	}, nil
}
