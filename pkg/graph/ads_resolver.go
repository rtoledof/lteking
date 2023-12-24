package graph

import (
	"context"

	"cubawheeler.io/pkg/cubawheeler"
)

var _ AdsResolver = &adsResolver{}

type adsResolver struct{ *Resolver }

// ID is the resolver for the id field.
func (r *adsResolver) ID(ctx context.Context, obj *cubawheeler.Ads) (string, error) {
	return obj.ID, nil
}

// Owner is the resolver for the owner field.
func (r *adsResolver) Owner(ctx context.Context, obj *cubawheeler.Ads) (*cubawheeler.Client, error) {
	client, err := r.client.FindById(ctx, obj.Client)
	if err != nil {
		return nil, err
	}
	return client, nil
}
