package realtime

import (
	"context"
	"log/slog"

	"cubawheeler.io/pkg/cubawheeler"
)

var (
	DriverLocations        = make(chan cubawheeler.Location, 10000)
	UserAvailabilityStatus = make(chan UserStatus, 1000)
	OrderChan              = make(chan *cubawheeler.Order, 10000)
)

type UserStatus struct {
	User      string
	Available bool
}

type Finder interface {
	FindNearByDrivers(context.Context, cubawheeler.GeoLocation) ([]*cubawheeler.Location, error)
}

type Updater interface {
	UpdateLocation(context.Context, string, cubawheeler.GeoLocation) error
}

type FinderUpdater interface {
	Finder
	Updater
}

type Notifier interface {
	NotifyToDevices(context.Context, []string) error
}

type UserUpdateService interface {
	SetAvailability(context.Context, string, bool) error
}

type RealTimeService struct {
	finder     FinderUpdater
	notifier   Notifier
	userUpdate UserUpdateService
	user       cubawheeler.UserService
}

func NewRealTimeService(
	finder FinderUpdater,
	notifier Notifier,
	user cubawheeler.UserService,
) *RealTimeService {

	s := &RealTimeService{
		finder:   finder,
		notifier: notifier,
		user:     user,
	}
	go storeOrUpdateDriversLocation(finder)
	go processNewOrder(s)
	go updateUserStatus(user)

	return s
}
func (s *RealTimeService) FindNearByDrivers(ctx context.Context, location cubawheeler.GeoLocation) ([]*cubawheeler.Location, error) {
	locations, err := s.finder.FindNearByDrivers(ctx, location)
	if err != nil {
		return nil, err
	}
	return locations, nil
}

func (s *RealTimeService) NotifyToDevices(ctx context.Context, users []string) error {
	return s.notifier.NotifyToDevices(ctx, users)
}

func storeOrUpdateDriversLocation(s Updater) {
	ctx := context.Background()
	for v := range DriverLocations {
		err := s.UpdateLocation(ctx, v.User, v.Geolocation)
		if err != nil {
			slog.Info("unable to update user real time")
		}
	}
}

func processNewOrder(s *RealTimeService) {

}

func updateUserStatus(s UserUpdateService) {
	ctx := context.Background()
	for v := range UserAvailabilityStatus {
		if err := s.SetAvailability(ctx, v.User, v.Available); err != nil {
			slog.Info("unable to update use availability")
		}
	}
}
