package mongo

import (
	"auth.io/models"
	"context"
	"go.mongodb.org/mongo-driver/bson"
)

var _ models.DeviceManager = &DeviceService{}

const deviceCollection Collections = "devices"

type DeviceService struct {
	*DB
}

func NewDeviceService(db *DB) *DeviceService {
	return &DeviceService{db}
}

func (d DeviceService) Remove(ctx context.Context, id string) error {
	//TODO implement me
	panic("implement me")
}

func (d DeviceService) FindAll(ctx context.Context) ([]models.Device, error) {
	//TODO implement me
	panic("implement me")
}

func (d DeviceService) AddDevice(ctx context.Context, device models.Device) error {
	//TODO implement me
	panic("implement me")
}

func findDevices(ctx context.Context, db *DB, filter models.DeviceFilter) ([]models.Device, string, error) {
	var token string
	var devices []models.Device
	user := models.UserFromContext(ctx)
	if user == nil {
		return nil, token, models.ErrAccessDenied
	}
	f := bson.D{}
	if len(filter.Name) > 0 {
		f = append(f, bson.E{Key: "name", Value: bson.D{{Key: "$in", Value: filter.Name}}})
	}
	if filter.Active != nil {
		f = append(f, bson.E{Key: "active", Value: filter.Active})
	}
	if user.Role != models.RoleAdmin {

	}
	cursor, err := db.Collection(deviceCollection).Find(ctx, f)
	if err != nil {
		return nil, "", models.NewInternalError(err)
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var device models.Device
		if err := cursor.Decode(&device); err != nil {
			return nil, token, err
		}
		devices = append(devices, device)
		if filter.Limit > 0 && len(devices) == filter.Limit+1 {
			devices = devices[:filter.Limit]
			return devices, device.ID.Hex(), nil
		}
	}
	if err := cursor.Err(); err != nil {
		return nil, token, models.NewInternalError(err)
	}
	return devices, token, nil
}

func findDeviceByName(ctx context.Context, name string) {

}
