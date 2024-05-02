package mock

import (
	"context"

	"auth.io/models"
)

var _ models.TenantService = &TenantService{}

type TenantService struct {
	CreateFn     func(context.Context, models.Tenant) (*models.Tenant, error)
	FindAllFn    func(context.Context, models.TenantFilter) ([]models.Tenant, string, error)
	FindByIDFn   func(context.Context, models.ID) (*models.Tenant, error)
	UpdateFn     func(context.Context, models.Tenant) (*models.Tenant, error)
	DeleteByIDFn func(context.Context, models.ID) error
}

// Create implements models.TenantService.
func (s *TenantService) Create(ctx context.Context, tenant models.Tenant) (*models.Tenant, error) {
	return s.CreateFn(ctx, tenant)
}

// DeleteByID implements models.TenantService.
func (s *TenantService) DeleteByID(ctx context.Context, id models.ID) error {
	return s.DeleteByIDFn(ctx, id)
}

// FindAll implements models.TenantService.
func (s *TenantService) FindAll(ctx context.Context, filter models.TenantFilter) ([]models.Tenant, string, error) {
	return s.FindAllFn(ctx, filter)
}

// FindByID implements models.TenantService.
func (s *TenantService) FindByID(ctx context.Context, id models.ID) (*models.Tenant, error) {
	return s.FindByIDFn(ctx, id)
}

// Update implements models.TenantService.
func (s *TenantService) Update(ctx context.Context, tenant models.Tenant) (*models.Tenant, error) {
	return s.UpdateFn(ctx, tenant)
}
