package graph

import (
	"context"
	"fmt"

	"cubawheeler.io/pkg/cubawheeler"
)

type messageResolver struct{ *Resolver }

// Order is the resolver for the order field.
func (r *messageResolver) Order(ctx context.Context, obj *cubawheeler.Message) (string, error) {
	panic(fmt.Errorf("not implemented: Order - order"))
}
