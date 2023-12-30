package handlers

import (
	"net/http"

	"cubawheeler.io/pkg/cubawheeler"
)

type CarHandler struct {
	User cubawheeler.UserService
}

func NewCarHandler(user cubawheeler.UserService) *CarHandler {
	return &CarHandler{
		User: user,
	}
}

func (h *CarHandler) Car(w http.ResponseWriter, r *http.Request) error {
	user := cubawheeler.UserFromContext(r.Context())
	if user == nil || user.Role != cubawheeler.RoleDriver {
		return cubawheeler.ErrUnauthorized
	}
	if err := r.ParseForm(); err != nil {
		return cubawheeler.NewError(err, http.StatusBadRequest, "invalid form")
	}

	if str := r.FormValue("car"); str != "" {
		for _, v := range user.Vehicles {
			if v.ID == str {
				user.ActiveVehicle = v.ID
				break
			}
		}
	}
	if user.ActiveVehicle == "" {
		return cubawheeler.NewError(nil, http.StatusBadRequest, "invalid car")
	}
	if err := h.User.Update(r.Context(), user); err != nil {
		return err
	}
	w.WriteHeader(http.StatusNoContent)
	return nil
}

func (h *CarHandler) Add(w http.ResponseWriter, r *http.Request) error {
	user := cubawheeler.UserFromContext(r.Context())
	if user == nil || user.Role != cubawheeler.RoleDriver {
		return cubawheeler.ErrUnauthorized
	}
	if err := r.ParseForm(); err != nil {
		return cubawheeler.NewError(err, http.StatusBadRequest, "invalid form")
	}
	panic("implement me")
}
