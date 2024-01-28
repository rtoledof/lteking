package order

import "context"

type DirectionService interface {
	GetRoute(context.Context, DirectionRequest) (_ *DirectionResponse, _ string, err error)
}
