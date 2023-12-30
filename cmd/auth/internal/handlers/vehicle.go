package handlers

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"cubawheeler.io/pkg/cubawheeler"
	"github.com/go-chi/chi/v5"
)

type VehicleHandler struct {
	Service cubawheeler.UserService
}

func NewVehicleHandler(service cubawheeler.UserService) *VehicleHandler {
	return &VehicleHandler{
		Service: service,
	}
}

func (h *VehicleHandler) Add(w http.ResponseWriter, r *http.Request) error {
	if err := r.ParseForm(); err != nil {
		return cubawheeler.ErrBadRequest
	}
	user := cubawheeler.UserFromContext(r.Context())
	vehicle, err := parseVehicleForm(r.Form)
	if err != nil {
		return err
	}
	user.Vehicles = append(user.Vehicles, vehicle)
	if err := h.Service.Update(r.Context(), user); err != nil {
		return err
	}
	return json.NewEncoder(w).Encode(vehicle)
}

func (h *VehicleHandler) Update(w http.ResponseWriter, r *http.Request) error {
	if err := r.ParseForm(); err != nil {
		return cubawheeler.ErrBadRequest
	}
	user := cubawheeler.UserFromContext(r.Context())
	vehicle, err := parseVehicleForm(r.Form)
	if err != nil {
		return err
	}
	for i, v := range user.Vehicles {
		if v.ID == vehicle.ID {
			user.Vehicles[i] = vehicle
			break
		}
	}
	if err := h.Service.Update(r.Context(), user); err != nil {
		return err
	}
	return json.NewEncoder(w).Encode(vehicle)
}

func (h *VehicleHandler) Remove(w http.ResponseWriter, r *http.Request) error {
	if err := r.ParseForm(); err != nil {
		return cubawheeler.ErrBadRequest
	}
	user := cubawheeler.UserFromContext(r.Context())
	id := chi.URLParam(r, "id")
	for i, v := range user.Vehicles {
		if v.ID == id {
			user.Vehicles = append(user.Vehicles[:i], user.Vehicles[i+1:]...)
			break
		}
	}
	if err := h.Service.Update(r.Context(), user); err != nil {
		return err
	}
	return json.NewEncoder(w).Encode(user.Vehicles)
}

func (h *VehicleHandler) List(w http.ResponseWriter, r *http.Request) error {
	user := cubawheeler.UserFromContext(r.Context())
	return json.NewEncoder(w).Encode(user.Vehicles)
}

func (h *VehicleHandler) FindByID(w http.ResponseWriter, r *http.Request) error {
	user := cubawheeler.UserFromContext(r.Context())
	id := chi.URLParam(r, "id")
	for _, v := range user.Vehicles {
		if v.ID == id {
			return json.NewEncoder(w).Encode(v)
		}
	}
	return cubawheeler.ErrNotFound
}

func (h *VehicleHandler) FindByPlate(w http.ResponseWriter, r *http.Request) error {
	user := cubawheeler.UserFromContext(r.Context())
	plate := chi.URLParam(r, "plate")
	for _, v := range user.Vehicles {
		if v.Plate == plate {
			return json.NewEncoder(w).Encode(v)
		}
	}
	return cubawheeler.ErrNotFound
}

func (h *VehicleHandler) SetActiveVehicle(w http.ResponseWriter, r *http.Request) error {
	user := cubawheeler.UserFromContext(r.Context())
	id := chi.URLParam(r, "id")
	for _, v := range user.Vehicles {
		if v.ID == id {
			user.ActiveVehicle = v.ID
			break
		}
	}
	if err := h.Service.Update(r.Context(), user); err != nil {
		return err
	}
	return json.NewEncoder(w).Encode(user)
}

func parseVehicleForm(form url.Values) (_ *cubawheeler.Vehicle, err error) {
	vehicle := &cubawheeler.Vehicle{
		CarModel:  form.Get("model"),
		Plate:     form.Get("plate"),
		Color:     form.Get("color"),
		Status:    cubawheeler.VehicleStatusNew,
		Type:      cubawheeler.VehicleType(form.Get("type")),
		CreatedAt: cubawheeler.Now().UTC().Unix(),
	}

	if brand := form.Get("brand"); brand != "" {
		vehicle.Brand = cubawheeler.Brand(brand)
		if !vehicle.Brand.IsValid() {
			return nil, cubawheeler.NewInvalidParameter("brand", brand)
		}
	}
	if year := form.Get("year"); year != "" {
		vehicle.Year, err = strconv.Atoi(year)
		if err != nil {
			return nil, cubawheeler.NewInvalidParameter("year", year)
		}
	}
	if seats := form.Get("seats"); seats != "" {
		vehicle.Seats, err = strconv.Atoi(seats)
		if err != nil {
			return nil, cubawheeler.NewInvalidParameter("seats", seats)
		}
	}
	if !vehicle.Type.IsValid() {
		return nil, cubawheeler.NewInvalidParameter("type", vehicle.Type)
	}
	if cat := form.Get("category"); cat != "" {
		vehicle.Category = cubawheeler.VehicleCategory(cat)
		if !vehicle.Category.IsValid() {
			return nil, cubawheeler.NewInvalidParameter("category", cat)
		}
	}
	if facilities := form.Get("facilities"); facilities != "" {
		for _, facility := range strings.Split(facilities, ",") {
			if !cubawheeler.Facilities(facility).IsValid() {
				return nil, cubawheeler.NewInvalidParameter("facilities", facility)
			}
			vehicle.Facilities = append(vehicle.Facilities, cubawheeler.Facilities(facility))
		}
	}

	return vehicle, nil
}
