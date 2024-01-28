package graph

import (
	"identity.io/graph/model"
	"identity.io/pkg/identity"
)

func assembleUpdateProfile(p model.ProfileInput) *identity.UpdateProfile {
	updateProfile := identity.UpdateProfile{}
	if p.Name != nil {
		updateProfile.Name = *p.Name
	}
	if p.Phone != nil {
		updateProfile.Phone = *p.Phone
	}
	if p.Dob != nil {
		updateProfile.Dob = *p.Dob
	}
	if p.PreferedCurrency != nil {
		updateProfile.PreferedCurrency = *p.PreferedCurrency
	}
	if p.Photo != nil {
		// TODO: process upload file
		updateProfile.Photo = *&p.Photo.Filename
	}
	if p.Dni != nil {
		// TODO: process upload file
		updateProfile.Dni = *&p.Dni.Filename
	}

	if p.Licence != nil {
		// TODO: process upload file
		updateProfile.Licence = *&p.Licence.Filename
	}
	if p.Circulation != nil {
		// TODO: Process upload file
		updateProfile.Circulation = *&p.Circulation.Filename
	}

	return &updateProfile
}

func assembleModelProfile(p *identity.Profile) *model.Profile {
	status := string(p.Status)
	return &model.Profile{
		ID:     p.ID,
		Name:   p.Name,
		Dob:    &p.DOB,
		Phone:  &p.Phone,
		Photo:  &p.Photo,
		Status: &status,
	}
}

func assembleVehicle(v model.VehicleInput) (*identity.Vehicle, error) {
	vehicle := &identity.Vehicle{}
	if v.Brand != nil {
		vehicle.Brand = identity.Brand(*v.Brand)
		if vehicle.Brand.IsValid() {
			return nil, identity.NewInvalidParameter("brand", v.Brand)
		}
	}
	if v.Category != nil {
		vehicle.Category = identity.VehicleCategory(*v.Category)
		if vehicle.Category.IsValid() {
			return nil, identity.NewInvalidParameter("category", v.Category)
		}
	}
	if v.Model != nil {
		vehicle.CarModel = *v.Model
	}
	if v.Colors != nil {
		vehicle.Colors = v.Colors
	}
	if v.Seats != nil {
		vehicle.Seats = *v.Seats
	}
	if v.Year != nil {
		vehicle.Year = *v.Year
	}
	if v.PlateNumber != nil {
		vehicle.Plate = *v.PlateNumber
	}
	if v.Type != nil {
		vehicle.Type = identity.VehicleType(*v.Type)
		if vehicle.Type.IsValid() {
			return nil, identity.NewInvalidParameter("type", v.Type)
		}
	}

	if v.Facilities != nil {
		vehicle.Facilities = []identity.Facilities{}
		for _, f := range v.Facilities {
			vehicle.Facilities = append(vehicle.Facilities, identity.Facilities(f))
		}
	}

	return vehicle, nil
}

func assembleModelVehicle(v *identity.Vehicle) (*model.Vehicle, error) {
	brand := model.Brand(v.Brand.String())
	if !v.Brand.IsValid() {
		return nil, identity.NewInvalidParameter("brand", v.Brand)
	}
	category := model.VechicleCategory(v.Category.String())
	if !v.Category.IsValid() {
		return nil, identity.NewInvalidParameter("category", v.Category)
	}
	vehicleType := model.VechicleType(v.Type.String())
	if !v.Type.IsValid() {
		return nil, identity.NewInvalidParameter("type", v.Type)
	}
	vehicle := &model.Vehicle{
		ID:          v.ID,
		Brand:       brand,
		Model:       v.CarModel,
		Category:    category,
		Colors:      v.Colors,
		PlateNumber: v.Plate,
		Seats:       v.Seats,
		Type:        vehicleType,
		Year:        &v.Year,
	}
	for _, f := range v.Facilities {
		facility := model.Facilities(f.String())
		if !f.IsValid() {
			return nil, identity.NewInvalidParameter("facility", f)
		}
		vehicle.Facilities = append(vehicle.Facilities, facility)
	}

	return vehicle, nil
}

func assembleVehicles(v []*identity.Vehicle) ([]*model.Vehicle, error) {
	vehicles := []*model.Vehicle{}
	for _, vehicle := range v {
		modelVehicle, err := assembleModelVehicle(vehicle)
		if err != nil {
			return nil, err
		}
		vehicles = append(vehicles, modelVehicle)
	}
	return vehicles, nil
}

func assembleLocation(l *identity.Location) (*model.Location, error) {
	loc := &model.Location{
		ID: l.ID,
		Point: &model.Point{
			Lat: l.Geolocation.Lat,
			Lng: l.Geolocation.Long,
		},
		Name: l.Name,
	}
	return loc, nil
}

func assembleLocations(l []*identity.Location) ([]*model.Location, error) {
	locations := []*model.Location{}
	for _, location := range l {
		loc, err := assembleLocation(location)
		if err != nil {
			return nil, err
		}
		locations = append(locations, loc)
	}
	return locations, nil
}