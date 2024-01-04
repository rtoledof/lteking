package cubawheeler

import (
	"context"
	"fmt"
	"io"
	"strconv"
)

type Vehicle struct {
	ID          string          `json:"id" bson:"_id"`
	Plate       string          `json:"plate,omitempty" bson:"plate"`
	Name        string          `json:"name,omitempty" bson:"name,omitmepty"`
	Category    VehicleCategory `json:"category" bson:"category"`
	Brand       Brand           `json:"brand" bson:"brand"`
	Year        int             `json:"year" bson:"year"`
	CarModel    string          `json:"model" bson:"car_model"`
	Seats       int             `json:"seats" bson:"seats"`
	Color       string          `json:"color" bson:"color"`
	Status      VehicleStatus   `json:"status" bson:"status"`
	Type        VehicleType     `json:"type" bson:"type"`
	Photos      []string        `json:"photos,omitempty" bson:"photos"`
	Facilities  []Facilities    `json:"facilities,omitempty" bson:"facilities"`
	User        string          `json:"-" bson:"user"`
	CreatedAt   int64           `json:"-" bson:"created_at"`
	UpdatedAt   int64           `json:"-" bson:"updated_at"`
	Circulation string          `json:"circulation,omitempty" bson:"circulation,omitempty"`
}

func (v Vehicle) IsValid() bool {
	return v.Plate != "" &&
		v.Category.IsValid() &&
		v.Brand.IsValid() &&
		v.Year > 0 &&
		v.CarModel != "" &&
		v.Seats > 0 &&
		v.Color != "" &&
		v.Status.IsValid() &&
		v.Type.IsValid()
}

type VehicleFilter struct {
	Limit int
	Ids   []string
	Token string
	Plate string
	Brand Brand
	Model string
	Color string
	User  string
	Type  VehicleType
}

type UpdateVehicle struct {
	ID         string
	Plate      string
	Category   VehicleCategory
	Year       int
	Type       VehicleType
	Facilities []Facilities
	User       string
}

type VehicleService interface {
	Store(context.Context, *Vehicle) error
	Update(context.Context, UpdateVehicle) (*Vehicle, error)
	FindByID(context.Context, string) (*Vehicle, error)
	FindByPlate(context.Context, string) (*Vehicle, error)
	FindAll(context.Context, *VehicleFilter) ([]*Vehicle, string, error)
}

type VehicleCategory string

const (
	VehicleCategoryX        VehicleCategory = "X"
	VehicleCategoryXl       VehicleCategory = "XL"
	VehicleCategoryConfort  VehicleCategory = "CONFORT"
	VehicleCategoryGreen    VehicleCategory = "GREEN"
	VehicleCategoryPets     VehicleCategory = "PETS"
	VehicleCategoryPackage  VehicleCategory = "PACKAGE"
	VehicleCategoryPriority VehicleCategory = "PRIORITY"
)

var AllVehicleCategory = []VehicleCategory{
	VehicleCategoryX,
	VehicleCategoryXl,
	VehicleCategoryConfort,
	VehicleCategoryGreen,
	VehicleCategoryPets,
	VehicleCategoryPackage,
	VehicleCategoryPriority,
}

func (e VehicleCategory) IsValid() bool {
	switch e {
	case VehicleCategoryX, VehicleCategoryXl, VehicleCategoryConfort, VehicleCategoryGreen, VehicleCategoryPets, VehicleCategoryPackage, VehicleCategoryPriority:
		return true
	}
	return false
}

func (e VehicleCategory) String() string {
	return string(e)
}

func (e *VehicleCategory) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = VehicleCategory(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid VehicleCategory", str)
	}
	return nil
}

func (e VehicleCategory) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type Facilities string

const (
	FacilitiesAirConditioner Facilities = "AIR_CONDITIONER"
	FacilitiesPetsAllowed    Facilities = "PETS_ALLOWED"
	FacilitiesSmokeAllowed   Facilities = "SMOKE_ALLOWED"
)

var AllFacilities = []Facilities{
	FacilitiesAirConditioner,
	FacilitiesPetsAllowed,
	FacilitiesSmokeAllowed,
}

func (e Facilities) IsValid() bool {
	switch e {
	case FacilitiesAirConditioner, FacilitiesPetsAllowed, FacilitiesSmokeAllowed:
		return true
	}
	return false
}

func (e Facilities) String() string {
	return string(e)
}

func (e *Facilities) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = Facilities(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid Facilities", str)
	}
	return nil
}

func (e Facilities) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type Brand string

const (
	BrandBmw Brand = "BMW"
)

var AllBrand = []Brand{
	BrandBmw,
}

func (e Brand) IsValid() bool {
	switch e {
	case BrandBmw:
		return true
	}
	return false
}

func (e Brand) String() string {
	return string(e)
}

func (e *Brand) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = Brand(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid Brand", str)
	}
	return nil
}

func (e Brand) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type VehicleType string

const (
	TypeAuto VehicleType = "AUTO"
	TypeMoto VehicleType = "MOTO"
	TypeBike VehicleType = "BIKE"
)

var AllType = []VehicleType{
	TypeAuto,
	TypeMoto,
	TypeBike,
}

func (e VehicleType) IsValid() bool {
	switch e {
	case TypeAuto, TypeMoto, TypeBike:
		return true
	}
	return false
}

func (e VehicleType) String() string {
	return string(e)
}

func (e *VehicleType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = VehicleType(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid Type", str)
	}
	return nil
}

func (e VehicleType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type VehicleStatus string

const (
	VehicleStatusNew       VehicleStatus = "NEW"
	VehicleStatusActive    VehicleStatus = "ACTIVE"
	VehicleStatusInactive  VehicleStatus = "INACTIVE"
	VehicleStatusSuspended VehicleStatus = "SUSPENDED"
)

var AllVehicleStatus = []VehicleStatus{
	VehicleStatusNew,
	VehicleStatusActive,
	VehicleStatusInactive,
	VehicleStatusSuspended,
}

func (e VehicleStatus) IsValid() bool {
	switch e {
	case VehicleStatusNew, VehicleStatusActive, VehicleStatusInactive, VehicleStatusSuspended:
		return true
	}
	return false
}

func (e VehicleStatus) String() string {
	return string(e)
}

func (e *VehicleStatus) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = VehicleStatus(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid VehicleStatus", str)
	}
	return nil
}

func (e VehicleStatus) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
