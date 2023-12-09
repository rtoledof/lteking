package graph

import (
	"context"
	"fmt"

	"cubawheeler.io/pkg/cubawheeler"
)

type pointResolver struct{ *Resolver }

// Long is the resolver for the long field.
func (r *pointResolver) Long(ctx context.Context, obj *cubawheeler.Point) (float64, error) {
	panic(fmt.Errorf("not implemented: Long - long"))
}
