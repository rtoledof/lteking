package handlers

import (
	"net/http"

	"cubawheeler.io/pkg/cubawheeler"
	"github.com/go-chi/chi/v5"
)

type ChargeHandler struct {
	Service cubawheeler.ChargeService
}

func NewChargeHandler(service cubawheeler.ChargeService) *ChargeHandler {
	return &ChargeHandler{
		Service: service,
	}
}

func (h *ChargeHandler) Charge(w http.ResponseWriter, r *http.Request) error {
	if !canDo(r, cubawheeler.RoleAdmin, cubawheeler.RoleDriver) {
		return cubawheeler.ErrUnauthorized
	}
	user := cubawheeler.UserFromContext(r.Context())
	if user == nil {
		return cubawheeler.ErrUnauthorized
	}

	idParam := chi.URLParam(r, "id")
	if idParam == "" {
		return cubawheeler.NewError(nil, http.StatusBadRequest, "invalid id")
	}
	charge, err := h.Service.FindByID(r.Context(), idParam)
	if err != nil {
		return err
	}
	if charge == nil {
		return cubawheeler.NewError(nil, http.StatusNotFound, "charge not found")
	}

	if charge.Rider != user.ID && user.Role != cubawheeler.RoleAdmin {
		return cubawheeler.ErrUnauthorized
	}

	return writeJSON(w, http.StatusOK, charge)
}

func (h *ChargeHandler) Charges(w http.ResponseWriter, r *http.Request) error {
	if !canDo(r, cubawheeler.RoleAdmin) {
		return cubawheeler.ErrUnauthorized
	}

	charges, err := h.Service.FindAll(r.Context(), cubawheeler.ChargeRequest{
		Limit: 100,
	})
	if err != nil {
		return err
	}

	return writeJSON(w, http.StatusOK, charges)
}
