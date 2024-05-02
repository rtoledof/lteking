package mock

import (
	"context"

	"auth.io/models"
)

var _ models.ClientService = &ClientService{}

type ClientService struct {
	CreateFn      func(context.Context, *models.Client) error
	DeleteByIDFn  func(context.Context, models.ID) error
	FindByIDFn    func(context.Context, models.ID) (*models.Client, error)
	FindByKeyFn   func(context.Context, string) (*models.Client, error)
	FindClientsFn func(context.Context, models.ClientFilter) ([]*models.Client, string, error)
	UpdateFn      func(context.Context, *models.Client) error
	UpdateKeyFn   func(context.Context, models.ID, bool) error
}

// Create implements models.ClientService.
func (s *ClientService) Create(ctx context.Context, client *models.Client) error {
	return s.CreateFn(ctx, client)
}

// DeleteByID implements models.ClientService.
func (s *ClientService) DeleteByID(ctx context.Context, id models.ID) error {
	return s.DeleteByIDFn(ctx, id)
}

// FindByID implements models.ClientService.
func (s *ClientService) FindByID(ctx context.Context, id models.ID) (*models.Client, error) {
	return s.FindByIDFn(ctx, id)
}

// FindByKey implements models.ClientService.
func (s *ClientService) FindByKey(ctx context.Context, key string) (*models.Client, error) {
	return s.FindByKeyFn(ctx, key)
}

// FindClients implements models.ClientService.
func (s *ClientService) FindClients(ctx context.Context, filter models.ClientFilter) ([]*models.Client, string, error) {
	return s.FindClientsFn(ctx, filter)
}

// Update implements models.ClientService.
func (s *ClientService) Update(ctx context.Context, client *models.Client) error {
	return s.UpdateFn(ctx, client)
}

// UpdateKey implements models.ClientService.
func (s *ClientService) UpdateKey(ctx context.Context, id models.ID, b bool) error {
	return s.UpdateKeyFn(ctx, id, b)
}
