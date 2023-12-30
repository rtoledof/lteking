package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"cubawheeler.io/pkg/cubawheeler"
)

type OtpHandler struct {
	OTP  cubawheeler.OtpService
	User cubawheeler.UserService
}

func (h *OtpHandler) Otp(w http.ResponseWriter, r *http.Request) error {
	if err := r.ParseForm(); err != nil {
		return fmt.Errorf("error parsing form: %v: %w", err, cubawheeler.ErrInternal)
	}
	client := cubawheeler.ClientFromContext(r.Context())
	userType := cubawheeler.RoleRider
	if client != nil && client.Type == cubawheeler.ApplicationTypeDriver {
		userType = cubawheeler.RoleDriver
	}
	email := r.Form.Get("email")
	if email == "" {
		return cubawheeler.NewMissingParameter("email")
	}
	user, err := h.User.FindByEmail(r.Context(), email)
	if err != nil {
		if errors.Is(err, cubawheeler.ErrNotFound) {
			user = &cubawheeler.User{
				ID:    cubawheeler.NewID().String(),
				Email: email,
				Role:  userType,
				Profile: cubawheeler.Profile{
					Status: cubawheeler.ProfileStatusIncompleted,
				},
				Status: cubawheeler.UserStatusOnReview,
			}
			if err := h.User.CreateUser(r.Context(), user); err != nil {
				return fmt.Errorf("error creating user: %v: %w", err, cubawheeler.ErrInternal)
			}
		} else {
			return fmt.Errorf("error finding user: %v: %w", err, cubawheeler.ErrInternal)
		}
	}
	otp, err := h.OTP.Create(r.Context(), email)
	if err != nil {
		return fmt.Errorf("error creating otp: %w", err)
	}
	user.Otp = otp
	if err := h.User.Update(r.Context(), user); err != nil {
		return fmt.Errorf("error updating user: %v: %w", err, cubawheeler.ErrInternal)
	}
	return json.NewEncoder(w).Encode(otp)
}
