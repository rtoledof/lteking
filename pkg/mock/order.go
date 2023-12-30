package mock

import (
	"context"

	"cubawheeler.io/pkg/cubawheeler"
)

var _ cubawheeler.OrderService = &OrderService{}

type OrderService struct {
	CreateFunc   func(ctx context.Context, req *cubawheeler.DirectionRequest) (*cubawheeler.Order, error)
	UpdateFunc   func(ctx context.Context, req *cubawheeler.DirectionRequest) (*cubawheeler.Order, error)
	FindByIDFunc func(ctx context.Context, id string) (*cubawheeler.Order, error)
	FindAllFunc  func(ctx context.Context, filter *cubawheeler.OrderFilter) (*cubawheeler.OrderList, error)

	AcceptOrderFunc  func(ctx context.Context, id string) (*cubawheeler.Order, error)
	CancelOrderFunc  func(ctx context.Context, id string) (*cubawheeler.Order, error)
	ConfirmOrderFunc func(ctx context.Context, req cubawheeler.ConfirmOrder) error
	FinishOrderFunc  func(ctx context.Context, id string) (*cubawheeler.Order, error)
	StartOrderFunc   func(ctx context.Context, id string) (*cubawheeler.Order, error)
}

// AcceptOrder implements cubawheeler.OrderService.
func (s *OrderService) AcceptOrder(ct context.Context, order string) (*cubawheeler.Order, error) {
	return s.AcceptOrderFunc(ct, order)
}

// CancelOrder implements cubawheeler.OrderService.
func (s *OrderService) CancelOrder(ctx context.Context, order string) (*cubawheeler.Order, error) {
	return s.CancelOrderFunc(ctx, order)
}

// ConfirmOrder implements cubawheeler.OrderService.
func (s *OrderService) ConfirmOrder(ctx context.Context, req cubawheeler.ConfirmOrder) error {
	return s.ConfirmOrderFunc(ctx, req)
}

// Create implements cubawheeler.OrderService.
func (s *OrderService) Create(ctx context.Context, req *cubawheeler.DirectionRequest) (*cubawheeler.Order, error) {
	return s.CreateFunc(ctx, req)
}

// FindAll implements cubawheeler.OrderService.
func (s *OrderService) FindAll(ctx context.Context, filter *cubawheeler.OrderFilter) (*cubawheeler.OrderList, error) {
	return s.FindAllFunc(ctx, filter)
}

// FindByID implements cubawheeler.OrderService.
func (s *OrderService) FindByID(ctx context.Context, order string) (*cubawheeler.Order, error) {
	return s.FindByIDFunc(ctx, order)
}

// FinishOrder implements cubawheeler.OrderService.
func (s *OrderService) FinishOrder(ctx context.Context, order string) (*cubawheeler.Order, error) {
	return s.FinishOrderFunc(ctx, order)
}

// StartOrder implements cubawheeler.OrderService.
func (s *OrderService) StartOrder(ctx context.Context, order string) (*cubawheeler.Order, error) {
	return s.StartOrderFunc(ctx, order)
}

// Update implements cubawheeler.OrderService.
func (s *OrderService) Update(ctx context.Context, req *cubawheeler.DirectionRequest) (*cubawheeler.Order, error) {
	return s.UpdateFunc(ctx, req)
}
