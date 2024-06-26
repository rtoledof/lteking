package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.43

import (
	"context"

	"wallet.io/graph/model"
)

// SetPin is the resolver for the setPin field.
func (r *mutationResolver) SetPin(ctx context.Context, pin string, old *string) (*model.Response, error) {
	rsp := &model.Response{
		Success: true,
	}
	if err := r.wallet.SetPin(ctx, *old, pin); err != nil {
		rsp.Success = false
		rsp.Errors = append(rsp.Errors, &model.Error{
			Message: err.Error(),
		})
	}
	return rsp, nil
}

// Withdraw is the resolver for the withdraw field.
func (r *mutationResolver) Withdraw(ctx context.Context, amount int, currency string) (*model.Response, error) {
	rsp := &model.Response{
		Success: true,
	}
	if err := r.wallet.Withdraw(ctx, int64(amount), currency); err != nil {
		rsp.Success = false
		rsp.Errors = append(rsp.Errors, &model.Error{
			Message: err.Error(),
		})
	}
	return rsp, nil
}

// Transfer is the resolver for the transfer field.
func (r *mutationResolver) Transfer(ctx context.Context, amount int, currency string, to string) (*model.Transfer, error) {
	transfer, err := r.wallet.Transfer(ctx, to, int64(amount), currency)
	if err != nil {
		return nil, err
	}
	return assembleModelTransfer(transfer), nil
}

// ConfirmTransfer is the resolver for the confirmTransfer field.
func (r *mutationResolver) ConfirmTransfer(ctx context.Context, id string, pin string) (*model.Response, error) {
	rsp := &model.Response{
		Success: true,
	}
	if err := r.wallet.ConfirmTransfer(ctx, id, pin); err != nil {
		rsp.Success = false
		rsp.Errors = append(rsp.Errors, &model.Error{
			Message: err.Error(),
		})
	}
	return rsp, nil
}

// Balance is the resolver for the balance field.
func (r *queryResolver) Balance(ctx context.Context, currency string) (int, error) {
	balance, err := r.wallet.Balance(ctx)
	if err != nil {
		return 0, err
	}
	return int(balance.Amount[currency]), nil
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
