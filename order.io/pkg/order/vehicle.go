package order

import (
	"fmt"
	"io"
	"strconv"
)

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
