package handlers

import (
	"fmt"
	"net/http"

	"cubawheeler.io/pkg/cannon"
	"cubawheeler.io/pkg/cubawheeler"
	"cubawheeler.io/pkg/derrors"
	"github.com/go-oauth2/oauth2/v4"
)

type LoginHandler struct {
	User        cubawheeler.UserService
	OTP         cubawheeler.OtpService
	Application cubawheeler.ApplicationService
	Token       oauth2.TokenStore
}

func (h *LoginHandler) Login(w http.ResponseWriter, r *http.Request) (err error) {
	defer derrors.Wrap(&err, "LoginHandler.Login")
	logger := cannon.LoggerFromContext(r.Context())
	logger.Info("login handler")
	if err := r.ParseForm(); err != nil {
		return fmt.Errorf("error parsing form: %v: %w", err, cubawheeler.ErrInternal)
	}
	email := r.FormValue("username")
	if email == "" {
		return cubawheeler.NewMissingParameter("username")
	}
	otp := r.FormValue("password")
	if otp == "" {
		return cubawheeler.NewMissingParameter("password")
	}
	if err := h.OTP.Otp(r.Context(), otp, email); err != nil {
		return fmt.Errorf("error validating otp: %v: %w", err, cubawheeler.ErrAccessDenied)
	}
	_, err = h.User.Login(r.Context(), cubawheeler.LoginRequest{
		Email: email,
		Otp:   otp,
	})
	if err != nil {
		return fmt.Errorf("error logging in: %v: %w", err, cubawheeler.ErrAccessDenied)
	}

	return nil
}
