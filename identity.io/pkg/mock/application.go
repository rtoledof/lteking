package mock

import (
	"context"

	"identity.io/pkg/identity"
)

var _ identity.ClientService = &ClientService{}

type ClientService struct {
	CreateFn      func(context.Context, *identity.Client) error
	DeleteByIDFn  func(context.Context, identity.ID) error
	FindByIDFn    func(context.Context, identity.ID) (*identity.Client, error)
	FindByKeyFn   func(context.Context, string) (*identity.Client, error)
	FindClientsFn func(context.Context, identity.ClientFilter) ([]*identity.Client, string, error)
	UpdateFn      func(context.Context, *identity.Client) error
	UpdateKeyFn   func(context.Context, identity.ID, bool) error
}

// Create implements identity.ClientService.
func (s *ClientService) Create(ctx context.Context, client *identity.Client) error {
	return s.CreateFn(ctx, client)
}

// DeleteByID implements identity.ClientService.
func (s *ClientService) DeleteByID(ctx context.Context, id identity.ID) error {
	return s.DeleteByIDFn(ctx, id)
}

// FindByID implements identity.ClientService.
func (s *ClientService) FindByID(ctx context.Context, id identity.ID) (*identity.Client, error) {
	return s.FindByIDFn(ctx, id)
}

// FindByKey implements identity.ClientService.
func (s *ClientService) FindByKey(ctx context.Context, key string) (*identity.Client, error) {
	return s.FindByKeyFn(ctx, key)
}

// FindClients implements identity.ClientService.
func (s *ClientService) FindClients(ctx context.Context, filter identity.ClientFilter) ([]*identity.Client, string, error) {
	return s.FindClientsFn(ctx, filter)
}

// Update implements identity.ClientService.
func (s *ClientService) Update(ctx context.Context, client *identity.Client) error {
	return s.UpdateFn(ctx, client)
}

// UpdateKey implements identity.ClientService.
func (s *ClientService) UpdateKey(ctx context.Context, id identity.ID, b bool) error {
	return s.UpdateKeyFn(ctx, id, b)
}
