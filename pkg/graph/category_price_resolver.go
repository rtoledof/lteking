package graph

import (
	"context"

	"cubawheeler.io/pkg/cubawheeler"
)

type categoryPriceResolver struct{ *Resolver }

// Price is the resolver for the price field.
func (r *categoryPriceResolver) Price(ctx context.Context, obj *cubawheeler.CategoryPrice) (int, error) {
	return int(obj.Price.Amount), nil
}
