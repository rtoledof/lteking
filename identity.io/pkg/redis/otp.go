package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"identity.io/pkg/identity"

	"identity.io/pkg/mailer"
)

var _ identity.OtpService = &OtpService{}

type token struct {
	Otp       string        `json:"otp"`
	Email     string        `json:"email"`
	ExpireIn  time.Duration `json:"expire_in"`
	CreatedAt time.Time     `json:"created_at"`
}

type OtpService struct {
	redis *Redis
}

func NewOtpService(client *Redis) *OtpService {
	return &OtpService{redis: client}
}

func (s *OtpService) Create(ctx context.Context, email string) (string, error) {
	client := identity.ClientFromContext(ctx)
	if client == nil {
		return "", identity.ErrUnauthorized
	}

	otp := token{
		Email:     email,
		Otp:       identity.NewOtp(),
		ExpireIn:  time.Minute * 3,
		CreatedAt: time.Now(),
	}
	data, err := json.Marshal(otp)
	if err != nil {
		return "", fmt.Errorf("unable to marshal otp: %v: %w", err, identity.ErrInternal)
	}
	if err = s.redis.client.Set(ctx, otp.Otp, data, otp.ExpireIn).Err(); err != nil {
		return "", fmt.Errorf("unable to store otp token: %v: %w", err, identity.ErrInternal)
	}
	go func() {
		textTemplate := fmt.Sprintf("Your otp is: %s", otp.Otp)
		htmlTemplate := fmt.Sprintf(fmt.Sprintf("<H2>Your Otp is: %s</H2>", otp.Otp))

		mailer.GenMessage("no-reply@identity.com", email, textTemplate, htmlTemplate)
	}()
	return otp.Otp, nil
}

func (s *OtpService) Otp(ctx context.Context, otp, email string) error {
	defer func() {
		if err := s.redis.client.Del(ctx, otp).Err(); err != nil {
			log.Println("unable to delete an user otp token")
		}
	}()
	data := s.redis.client.Get(ctx, otp)
	if data == nil {
		return identity.ErrNotFound
	}
	var t token
	b, err := data.Bytes()
	if err != nil {
		return fmt.Errorf("unable to get token info: %v: %w", err, identity.ErrInternal)
	}
	if err := json.Unmarshal(b, &t); err != nil {
		return fmt.Errorf("unable to decode token info: %v: %w", err, identity.ErrInternal)
	}
	if t.Email != email {
		return identity.ErrNotFound
	}
	return nil
}
