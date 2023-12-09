package graph

import (
	"context"
	"cubawheeler.io/pkg/cubawheeler"
	"fmt"
)

type messageResolver struct{ *Resolver }

// Order is the resolver for the order field.
func (r *messageResolver) Order(ctx context.Context, obj *cubawheeler.Message) (string, error) {
	panic(fmt.Errorf("not implemented: Order - order"))
}
