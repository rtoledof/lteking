package mock

import (
	"context"

	"order.io/pkg/order"
)

var _ order.OrderService = &OrderService{}

type OrderService struct {
	CreateFunc   func(context.Context, order.Item) (*order.Order, error)
	UpdateFunc   func(context.Context, string, order.Item) (*order.Order, error)
	FindByIDFunc func(context.Context, string) (*order.Order, error)
	FindAllFunc  func(context.Context, *order.OrderFilter) (*order.OrderList, error)

	AcceptOrderFunc  func(context.Context, string) (*order.Order, error)
	CancelOrderFunc  func(context.Context, string) (*order.Order, error)
	ConfirmOrderFunc func(context.Context, order.ConfirmOrder) error
	FinishOrderFunc  func(context.Context, string) (*order.Order, error)
	StartOrderFunc   func(context.Context, string) (*order.Order, error)
}

// AcceptOrder implements order.OrderService.
func (*OrderService) AcceptOrder(context.Context, string) error {
	panic("unimplemented")
}

// CancelOrder implements order.OrderService.
func (*OrderService) CancelOrder(context.Context, string) error {
	panic("unimplemented")
}

// Categories implements order.OrderService.
func (*OrderService) Categories(context.Context, string) ([]*order.CategoryPrice, error) {
	panic("unimplemented")
}

// ConfirmOrder implements order.OrderService.
func (*OrderService) ConfirmOrder(context.Context, order.ConfirmOrder) error {
	panic("unimplemented")
}

// Create implements order.OrderService.
func (*OrderService) Create(context.Context, order.Item) (*order.Order, error) {
	panic("unimplemented")
}

// FindAll implements order.OrderService.
func (*OrderService) FindAll(context.Context, order.OrderFilter) (*order.OrderList, error) {
	panic("unimplemented")
}

// FindByID implements order.OrderService.
func (*OrderService) FindByID(context.Context, string) (*order.Order, error) {
	panic("unimplemented")
}

// FinishOrder implements order.OrderService.
func (*OrderService) FinishOrder(context.Context, string) error {
	panic("unimplemented")
}

// RateOrder implements order.OrderService.
func (*OrderService) RateOrder(context.Context, string, float64, string) error {
	panic("unimplemented")
}

// StartOrder implements order.OrderService.
func (*OrderService) StartOrder(context.Context, string) error {
	panic("unimplemented")
}

// Update implements order.OrderService.
func (*OrderService) Update(context.Context, string, order.Item) (*order.Order, error) {
	panic("unimplemented")
}
