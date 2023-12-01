package otp

import (
	"context"
	"cubawheeler.io/pkg/cubawheeler"
	"math/rand"
	"time"
)

type OtpService struct {
	user cubawheeler.UserService
}

func NewOtpService(user cubawheeler.UserService) *OtpService {
	return &OtpService{user: user}
}

func (s *OtpService) New(ctx context.Context, email string) (int, error) {
	rand.Seed(time.Now().UnixNano())
	otp := uint64(rand.Intn(999999 - 100000))
	if err := s.user.UpdateOTP(ctx, email, otp); err != nil {
		return 0, err
	}
	return int(otp), nil
}
