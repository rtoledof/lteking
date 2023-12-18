package graph

import (
	"context"
	"fmt"

	"cubawheeler.io/pkg/cubawheeler"
)

type planResolver struct{ *Resolver }

// Orders is the resolver for the orders field.
func (r *planResolver) Orders(ctx context.Context, obj *cubawheeler.Plan) (int, error) {
	panic(fmt.Errorf("not implemented: Orders - orders"))
}
