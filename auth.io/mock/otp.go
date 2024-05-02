package mock

import (
	"context"

	"auth.io/models"
)

var _ models.OtpService = &OtpService{}

type OtpService struct {
	CreateFn func(context.Context, string) (string, error)
	OtpFn    func(context.Context, string, string) error
}

// Create implements cubawheeler.OtpService.
func (s *OtpService) Create(ctx context.Context, email string) (string, error) {
	return s.CreateFn(ctx, email)
}

// Otp implements cubawheeler.OtpService.
func (s *OtpService) Otp(ctx context.Context, email, otp string) error {
	return s.OtpFn(ctx, email, otp)
}
