package mock

import (
	"context"

	"cubawheeler.io/pkg/cubawheeler"
)

var _ cubawheeler.ApplicationService = &ApplicationService{}

type ApplicationService struct {
	CreateApplicationFn            func(ctx context.Context, input cubawheeler.ApplicationRequest) (*cubawheeler.Application, error)
	FindApplicationsFn             func(ctx context.Context, input *cubawheeler.ApplicationFilter) (*cubawheeler.ApplicationList, error)
	FindByClientFn                 func(ctx context.Context, client string) (*cubawheeler.Application, error)
	FindByIDFn                     func(ctx context.Context, input string) (*cubawheeler.Application, error)
	UpdateApplicationCredentialsFn func(ctx context.Context, application string) (*cubawheeler.Application, error)
}

// CreateApplication implements cubawheeler.ApplicationService.
func (s *ApplicationService) CreateApplication(ctx context.Context, input cubawheeler.ApplicationRequest) (*cubawheeler.Application, error) {
	return s.CreateApplicationFn(ctx, input)
}

// FindApplications implements cubawheeler.ApplicationService.
func (s *ApplicationService) FindApplications(ctx context.Context, input *cubawheeler.ApplicationFilter) (*cubawheeler.ApplicationList, error) {
	return s.FindApplicationsFn(ctx, input)
}

// FindByClient implements cubawheeler.ApplicationService.
func (s *ApplicationService) FindByClient(ctx context.Context, client string) (*cubawheeler.Application, error) {
	return s.FindByClientFn(ctx, client)
}

// FindByID implements cubawheeler.ApplicationService.
func (s *ApplicationService) FindByID(ctx context.Context, input string) (*cubawheeler.Application, error) {
	return s.FindByIDFn(ctx, input)
}

// UpdateApplicationCredentials implements cubawheeler.ApplicationService.
func (s *ApplicationService) UpdateApplicationCredentials(ctx context.Context, application string) (*cubawheeler.Application, error) {
	return s.UpdateApplicationCredentialsFn(ctx, application)
}
