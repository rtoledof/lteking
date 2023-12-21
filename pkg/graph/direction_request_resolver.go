package graph

import (
	"context"

	"cubawheeler.io/pkg/cubawheeler"
)

type directionRequestResolver struct{ *Resolver }

// Points is the resolver for the points field.
func (r *directionRequestResolver) Points(ctx context.Context, obj *cubawheeler.DirectionRequest, data []*cubawheeler.PointInput) error {
	for _, v := range data {
		obj.Points = append(obj.Points, &cubawheeler.Point{
			Lat: v.Lat,
			Lng: v.Lng,
		})
	}
	return nil
}
