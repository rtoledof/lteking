package cubawheeler

import "context"

type RealTimeService interface {
	FindNearByDrivers(context.Context, GeoLocation) ([]*Location, error)
	NotifyDrivers(context.Context, []*User) error
}

type NearByResponse struct {
	Driver   *User     `json:"driver"`
	Location *Location `json:"location"`
}
