package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.43

import (
	"context"
	"fmt"

	"order.io/graph/model"
	"order.io/pkg/order"
)

// CreateRide is the resolver for the createRide field.
func (r *mutationResolver) CreateRide(ctx context.Context, input model.RideInput) (*model.Order, error) {
	order, err := r.order.Create(ctx, assembleOrderItem(&input))
	if err != nil {
		return nil, err
	}
	return assembleModelOrder(order)
}

// UpdateRide is the resolver for the updateRide field.
func (r *mutationResolver) UpdateRide(ctx context.Context, id string, input model.RideInput) (*model.Order, error) {
	order, err := r.order.Update(ctx, id, assembleOrderItem(&input))
	if err != nil {
		return nil, err
	}
	return assembleModelOrder(order)
}

// ConfirmRide is the resolver for the confirmRide field.
func (r *mutationResolver) ConfirmRide(ctx context.Context, input model.ConfirmRideInput) (*model.Response, error) {
	rsp := &model.Response{
		Success: true,
	}
	if !input.Category.IsValid() {
		rsp.Success = false
		rsp.Errors = append(rsp.Errors, &model.Error{
			Field:   "category",
			Message: "invalid category",
		})
		return rsp, nil
	}
	if !input.Method.IsValid() {
		rsp.Success = false
		rsp.Errors = append(rsp.Errors, &model.Error{
			Field:   "method",
			Message: "invalid method",
		})
		return rsp, nil
	}

	if err := r.order.ConfirmOrder(ctx, order.ConfirmOrder{
		OrderID:  input.ID,
		Category: order.VehicleCategory(input.Category),
		Method:   order.ChargeMethod(input.Method),
	}); err != nil {
		rsp.Success = false
		rsp.Errors = append(rsp.Errors, &model.Error{
			Field:   "order",
			Message: err.Error(),
		})
	}
	return rsp, nil
}

// CancelRide is the resolver for the cancelRide field.
func (r *mutationResolver) CancelRide(ctx context.Context, id string) (*model.Response, error) {
	rsp := &model.Response{
		Success: true,
	}
	if err := r.order.CancelOrder(ctx, id); err != nil {
		rsp.Success = false
		rsp.Errors = append(rsp.Errors, &model.Error{
			Field:   "order",
			Message: err.Error(),
		})
	}
	return rsp, nil
}

// AcceptRide is the resolver for the acceptRide field.
func (r *mutationResolver) AcceptRide(ctx context.Context, id string) (*model.Response, error) {
	rsp := &model.Response{
		Success: true,
	}
	if err := r.order.AcceptOrder(ctx, id); err != nil {
		rsp.Success = false
		rsp.Errors = append(rsp.Errors, &model.Error{
			Field:   "order",
			Message: err.Error(),
		})
	}
	return rsp, nil
}

// StartRide is the resolver for the startRide field.
func (r *mutationResolver) StartRide(ctx context.Context, id string) (*model.Response, error) {
	rsp := &model.Response{
		Success: true,
	}
	if err := r.order.StartOrder(ctx, id); err != nil {
		rsp.Success = false
		rsp.Errors = append(rsp.Errors, &model.Error{
			Field:   "order",
			Message: err.Error(),
		})
	}
	return rsp, nil
}

// FinishRide is the resolver for the finishRide field.
func (r *mutationResolver) FinishRide(ctx context.Context, id string) (*model.Response, error) {
	rsp := &model.Response{
		Success: true,
	}
	if err := r.order.FinishOrder(ctx, id); err != nil {
		rsp.Success = false
		rsp.Errors = append(rsp.Errors, &model.Error{
			Field:   "order",
			Message: err.Error(),
		})
	}
	return rsp, nil
}

// RateRide is the resolver for the rateRide field.
func (r *mutationResolver) RateRide(ctx context.Context, id string, rate float64, comment *string) (*model.Response, error) {
	rsp := &model.Response{
		Success: true,
	}
	if err := r.order.RateOrder(ctx, id, rate, *comment); err != nil {
		rsp.Success = false
		rsp.Errors = append(rsp.Errors, &model.Error{
			Field:   "order",
			Message: err.Error(),
		})
	}
	return rsp, nil
}

// PayRide is the resolver for the payRide field.
// TODO: move this to payment service
func (r *mutationResolver) PayRide(ctx context.Context, id string, method model.PaymentMethod) (*model.Response, error) {
	panic(fmt.Errorf("not implemented: PayRide - payRide"))
}

// TODO: move this to identity service
func (r *mutationResolver) RateRider(ctx context.Context, id string, rate float64, comment *string) (*model.Response, error) {
	panic(fmt.Errorf("not implemented: RateRider - rateRider"))
}

// Orders is the resolver for the orders field.
func (r *queryResolver) Orders(ctx context.Context, filter model.OrderListFilter) (*model.OrdersResponse, error) {
	orders, err := r.order.FindAll(ctx, assembleOrderFilter(filter))
	if err != nil {
		return nil, err
	}
	items := make([]*model.Order, len(orders.Data))
	for i, o := range orders.Data {
		item, err := assembleModelOrder(o)
		if err != nil {
			return nil, err
		}
		items[i] = item
	}
	return &model.OrdersResponse{
		Items: items,
		Token: orders.Token,
	}, nil
}

// Order is the resolver for the order field.
func (r *queryResolver) Order(ctx context.Context, id string) (*model.Order, error) {
	order, err := r.order.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return assembleModelOrder(order)
}

// Categories is the resolver for the categories field.
func (r *queryResolver) Categories(ctx context.Context, order string) ([]*model.CategoryPrice, error) {
	categories, err := r.order.Categories(ctx, order)
	if err != nil {
		return nil, err
	}
	return assembleCategoryPrices(categories), nil
}

// PaymentMethods is the resolver for the paymentMethods field.
func (r *queryResolver) PaymentMethods(ctx context.Context) ([]model.PaymentMethod, error) {
	return []model.PaymentMethod{
		model.PaymentMethodCash,
		model.PaymentMethodBalance,
		model.PaymentMethodMLCTransaction,
		model.PaymentMethodCUPTransaction,
	}, nil
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
