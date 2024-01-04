package handlers

import (
	"fmt"
	"net/http"

	"cubawheeler.io/pkg/cannon"
	"cubawheeler.io/pkg/cubawheeler"
	"cubawheeler.io/pkg/derrors"
)

type AuthorizeHandler struct {
	Service cubawheeler.ApplicationService
}

func NewAuthorizeHandler(service cubawheeler.ApplicationService) *AuthorizeHandler {
	return &AuthorizeHandler{Service: service}
}

func (h *AuthorizeHandler) Authorize(w http.ResponseWriter, r *http.Request) (err error) {
	defer derrors.Wrap(&err, "AuthorizeHandler.Login")
	logger := cannon.LoggerFromContext(r.Context())
	logger.Info("start authorize handler")
	if err := r.ParseForm(); err != nil {
		return fmt.Errorf("error parsing form: %v: %w", err, cubawheeler.ErrInternal)
	}
	clientID := r.Form.Get("client_id")
	if clientID == "" {
		return cubawheeler.NewMissingParameter("client_id")
	}
	clientSecret := r.Form.Get("client_secret")
	if clientSecret == "" {
		return cubawheeler.NewMissingParameter("client_secret")
	}
	app, err := h.Service.FindByClient(r.Context(), clientID)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return fmt.Errorf("error finding application: %v: %w", err, cubawheeler.ErrAccessDenied)
	}
	if app.Secret != clientSecret {
		logger.Info("invalid credentials")
		w.WriteHeader(http.StatusUnauthorized)
		return fmt.Errorf("invalid credentials: %w", cubawheeler.ErrAccessDenied)
	}
	logger.Info("end authorize handler")
	return nil
}
