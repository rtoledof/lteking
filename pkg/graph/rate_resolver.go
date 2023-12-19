package graph

import (
	"context"

	"cubawheeler.io/pkg/cubawheeler"
)

type rateRequestResolver struct{ *Resolver }

// HighDemand is the resolver for the high_demand field.
func (r *rateRequestResolver) HighDemand(ctx context.Context, obj *cubawheeler.RateRequest, data *bool) error {
	obj.HiDemand = data
	return nil
}
