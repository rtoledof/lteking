package handlers

import (
	"fmt"
	"net/http"

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
	if err := r.ParseForm(); err != nil {
		return fmt.Errorf("error parsing form: %v: %w", err, cubawheeler.ErrInternal)
	}
	grandType := r.Form.Get("grant_type")
	switch grandType {
	case "password":
		email := r.Form.Get("username")
		if email == "" {
			return cubawheeler.NewMissingParameter("username")
		}
		otp := r.Form.Get("password")
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
	case "client_credentials":
		clientID := r.Form.Get("client_id")
		if clientID == "" {
			return cubawheeler.NewMissingParameter("client_id")
		}
		clientSecret := r.Form.Get("client_secret")
		if clientSecret == "" {
			return cubawheeler.NewMissingParameter("client_secret")
		}
		app, err := h.Application.FindByID(r.Context(), clientID)
		if err != nil || app.Secret != clientSecret {
			return fmt.Errorf("error finding application: %v: %w", err, cubawheeler.ErrAccessDenied)
		}
	default:
		return fmt.Errorf("error invalid grant_type: %v: %w", grandType, cubawheeler.ErrAccessDenied)
	}

	return nil
}
