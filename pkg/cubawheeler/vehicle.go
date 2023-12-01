package cubawheeler

import (
	"fmt"
	"gorm.io/gorm"
	"io"
	"strconv"
)

type Vehicle struct {
	gorm.Model
	ID         string          `json:"id" gorm:"privateKey;varchar(36);not null"`
	Plate      *string         `json:"plate,omitempty"`
	Category   VehicleCategory `json:"category"`
	Brand      Brand           `json:"brand"`
	Year       int             `json:"year"`
	CarModel   string          `json:"model"`
	Seats      int             `json:"seats"`
	Color      string          `json:"color"`
	Type       VehicleType     `json:"type"`
	Photos     []string        `json:"photos,omitempty"`
	Facilities []Facilities    `json:"facilities,omitempty"`
	User       string          `json:"-"`
}

func (v *Vehicle) BeforeSave(*gorm.DB) error {
	if v.ID == "" {
		v.ID = NewID().String()
	}
	return nil
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
