package graph

import (
	"auth.io/graph/model"
	"auth.io/models"
)

func assembleUpdateProfile(p model.ProfileInput) *models.UpdateProfile {
	updateProfile := models.UpdateProfile{}
	if p.FirstName != nil {
		updateProfile.Name = *p.FirstName
	}
	if p.LastName != nil {
		updateProfile.LastName = *p.LastName
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
		updateProfile.Photo = *p.Photo
	}
	if p.Dni != nil {
		// TODO: process upload file
		updateProfile.Dni = *p.Dni
	}

	if p.Licence != nil {
		// TODO: process upload file
		updateProfile.Licence = *&p.Licence.Filename
	}
	if p.Circulation != nil {
		// TODO: Process upload file
		updateProfile.Circulation = *&p.Circulation.Filename
	}
	if p.PreferedCurrency != nil {
		updateProfile.PreferedCurrency = *p.PreferedCurrency
	}

	return &updateProfile
}

func assembleModelProfile(p *models.User) *model.Profile {
	profile := &model.Profile{}
	if p == nil {
		return profile
	}
	profile.ID = p.ID
	profile.Email = p.Email
	status := p.Status.String()
	profile.Status = &status
	profile.ReferalCode = &p.Referer
	if p.Profile != nil {
		profile.LastName = p.Profile.LastName
		profile.FirstName = p.Profile.Name
		profile.Dob = &p.Profile.DOB
		profile.Phone = &p.Profile.Phone
		profile.Photo = &p.Profile.Photo
		profile.PreferedCurrency = &p.Profile.PreferedCurrency
		if activeVehicle := p.GetActiveVehicle(); activeVehicle != nil {
			vehicle, err := assembleModelVehicle(p.GetActiveVehicle())
			if err != nil {
				return profile
			}
			profile.ActiveVehicle = vehicle
		}
	}

	return profile
}

func assembleVehicle(v model.VehicleInput) (*models.Vehicle, error) {
	vehicle := &models.Vehicle{}
	if v.Brand != nil {
		vehicle.Brand = models.Brand(*v.Brand)
		if vehicle.Brand.IsValid() {
			return nil, models.NewInvalidParameter("brand", v.Brand)
		}
	}
	if v.Category != nil {
		vehicle.Category = models.VehicleCategory(*v.Category)
		if vehicle.Category.IsValid() {
			return nil, models.NewInvalidParameter("category", v.Category)
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
		vehicle.Type = models.VehicleType(*v.Type)
		if vehicle.Type.IsValid() {
			return nil, models.NewInvalidParameter("type", v.Type)
		}
	}

	if v.Facilities != nil {
		vehicle.Facilities = []models.Facilities{}
		for _, f := range v.Facilities {
			vehicle.Facilities = append(vehicle.Facilities, models.Facilities(f))
		}
	}

	return vehicle, nil
}

func assembleModelVehicle(v *models.Vehicle) (*model.Vehicle, error) {
	brand := model.Brand(v.Brand.String())
	if !v.Brand.IsValid() {
		return nil, models.NewInvalidParameter("brand", v.Brand)
	}
	category := model.VechicleCategory(v.Category.String())
	if !v.Category.IsValid() {
		return nil, models.NewInvalidParameter("category", v.Category)
	}
	vehicleType := model.VechicleType(v.Type.String())
	if !v.Type.IsValid() {
		return nil, models.NewInvalidParameter("type", v.Type)
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
			return nil, models.NewInvalidParameter("facility", f)
		}
		vehicle.Facilities = append(vehicle.Facilities, facility)
	}

	return vehicle, nil
}

func assembleVehicles(v []*models.Vehicle) ([]*model.Vehicle, error) {
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

func assembleLocation(l *models.Location) (*model.Location, error) {
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

func assembleLocations(l []*models.Location) ([]*model.Location, error) {
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
