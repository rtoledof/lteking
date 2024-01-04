package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"cubawheeler.io/pkg/cubawheeler"
)

type ProfileHandler struct {
	User cubawheeler.UserService
}

func NewProfileHandler(user cubawheeler.UserService) *ProfileHandler {
	return &ProfileHandler{
		User: user,
	}
}

func (h *ProfileHandler) Update(w http.ResponseWriter, r *http.Request) error {
	user := cubawheeler.UserFromContext(r.Context())
	if user == nil {
		return cubawheeler.ErrUnauthorized
	}
	if err := r.ParseForm(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return cubawheeler.NewError(err, http.StatusBadRequest, "invalid form")
	}
	if str := r.FormValue("name"); str != "" {
		user.Name = str
	}
	if str := r.FormValue("phone"); str != "" {
		user.Profile.Phone = str
	}
	if str := r.FormValue("last_name"); str != "" {
		user.Profile.LastName = str
	}
	if str := r.FormValue("dob"); str != "" {
		user.Profile.DOB = str
	}
	if str := r.FormValue("photo"); str != "" {
		user.Profile.Photo = str
	}
	if str := r.FormValue("gender"); str != "" {
		user.Profile.Gender = cubawheeler.Gender(str)
		if !user.Profile.Gender.IsValid() {
			w.WriteHeader(http.StatusBadRequest)
			return cubawheeler.NewError(cubawheeler.ErrInvalidInput, http.StatusBadRequest, "invalid gender")
		}
	}
	if err := h.User.Update(r.Context(), user); err != nil {
		e := &cubawheeler.Error{}
		if errors.As(err, &e) {
			w.WriteHeader(e.StatusCode)
			return e
		}
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}
	w.WriteHeader(http.StatusNoContent)
	return nil
}

func (h *ProfileHandler) Get(w http.ResponseWriter, r *http.Request) error {
	user := cubawheeler.UserFromContext(r.Context())
	if user == nil {
		return cubawheeler.ErrUnauthorized
	}
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(user.Profile)
}

func (h *ProfileHandler) AddDevice(w http.ResponseWriter, r *http.Request) error {
	user := cubawheeler.UserFromContext(r.Context())
	if user == nil {
		return cubawheeler.ErrUnauthorized
	}
	if err := r.ParseForm(); err != nil {
		return cubawheeler.NewError(err, http.StatusBadRequest, "invalid form")
	}
	if str := r.FormValue("device_id"); str != "" {
		if err := h.User.AddDevice(r.Context(), str); err != nil {
			return err
		}
	}
	w.WriteHeader(http.StatusNoContent)
	return nil
}
