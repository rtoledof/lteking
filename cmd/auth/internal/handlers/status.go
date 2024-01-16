package handlers

import (
	"net/http"
	"strconv"

	"cubawheeler.io/pkg/cubawheeler"
)

type StatusHandler struct {
	User cubawheeler.UserService
}

func NewStatusHandler(user cubawheeler.UserService) *StatusHandler {
	return &StatusHandler{
		User: user,
	}
}

func (h *StatusHandler) Availability(w http.ResponseWriter, r *http.Request) error {
	user := cubawheeler.UserFromContext(r.Context())
	if user == nil || user.Role != cubawheeler.RoleDriver {
		return cubawheeler.ErrUnauthorized
	}
	if err := r.ParseForm(); err != nil {
		return cubawheeler.NewError(err, http.StatusBadRequest, "invalid form")
	}

	if str := r.FormValue("active"); str != "" {
		status, err := strconv.ParseBool(str)
		if err != nil {
			return cubawheeler.NewError(err, http.StatusBadRequest, "invalid status")
		}
		user.Available = status
		if err := h.User.Update(r.Context(), user); err != nil {
			return err
		}
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}
