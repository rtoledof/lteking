package handlers

import (
	"net/http"
	"strconv"

	"cubawheeler.io/pkg/cannon"
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

func (h *CarHandler) SetActiveVehicle(w http.ResponseWriter, r *http.Request) error {
	user := cubawheeler.UserFromContext(r.Context())
	if user == nil || user.Role != cubawheeler.RoleDriver {
		return cubawheeler.ErrUnauthorized
	}
	logger := cannon.LoggerFromContext(r.Context())
	logger.Info("car handler set active vehicle")
	if err := r.ParseForm(); err != nil {
		return cubawheeler.NewError(err, http.StatusBadRequest, "invalid form")
	}

	if str := r.FormValue("car"); str != "" {
		if !user.HasVehicle(str) {
			return cubawheeler.NewError(nil, http.StatusNotFound, "invalid car")
		}
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
	logger := cannon.LoggerFromContext(r.Context())
	logger.Info("add car handler")
	user := cubawheeler.UserFromContext(r.Context())
	if user == nil || user.Role != cubawheeler.RoleDriver {
		w.WriteHeader(http.StatusUnauthorized)
		return cubawheeler.ErrUnauthorized
	}
	if err := r.ParseForm(); err != nil {
		return cubawheeler.NewError(err, http.StatusBadRequest, "invalid form")
	}
	id := r.FormValue("id")
	if id == "" {
		id = cubawheeler.NewID().String()
	}

	vehicle := &cubawheeler.Vehicle{
		ID:       id,
		Plate:    r.FormValue("plate"),
		Name:     r.FormValue("name"),
		Category: cubawheeler.VehicleCategory(r.FormValue("category")),
		Brand:    cubawheeler.Brand(r.FormValue("brand")),
		Year: func() int {
			year, _ := strconv.Atoi(r.FormValue("year"))
			return year
		}(),
		CarModel: r.FormValue("model"),
		Seats: func() int {
			seat, _ := strconv.Atoi(r.FormValue("seats"))
			return seat
		}(),
		Color:       r.FormValue("color"),
		Status:      cubawheeler.VehicleStatusNew,
		Type:        cubawheeler.VehicleType(r.FormValue("type")),
		Photos:      r.Form["photos"],
		Circulation: r.FormValue("circulation"),
	}
	if !vehicle.IsValid() {
		return cubawheeler.NewError(nil, http.StatusBadRequest, "invalid vehicle")
	}

	user.Vehicles = append(user.Vehicles, vehicle)

	if err := h.User.Update(r.Context(), user); err != nil {
		return err
	}
	w.WriteHeader(http.StatusNoContent)
	return nil
}
