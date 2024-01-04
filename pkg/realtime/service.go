package realtime

import (
	"context"
	"log/slog"

	"cubawheeler.io/pkg/cubawheeler"
	"cubawheeler.io/pkg/redis"
)

var (
	DriverLocations        = make(chan cubawheeler.Location, 10000)
	UserAvailabilityStatus = make(chan UserStatus, 1000)
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

type Order struct {
	ID       string
	Price    int
	Currency string
	Points   []cubawheeler.GeoLocation
}

type Notifier interface {
	NotifyToDevices(context.Context, []string, string) error
}

type UserUpdateService interface {
	SetAvailability(context.Context, string, bool) error
}

type RealTimeService struct {
	finder     FinderUpdater
	notifier   Notifier
	userUpdate UserUpdateService
	user       cubawheeler.UserService
	redis      *redis.Redis
	order      cubawheeler.OrderService
}

func NewRealTimeService(
	finder FinderUpdater,
	notifier Notifier,
	user cubawheeler.UserService,
	redis *redis.Redis,
	order cubawheeler.OrderService,
) *RealTimeService {

	s := &RealTimeService{
		finder:   finder,
		notifier: notifier,
		user:     user,
		redis:    redis,
		order:    order,
	}
	go storeOrUpdateDriversLocation(finder)
	go processNewOrder(s)
	go updateUserStatus(user)
	go notifyDrivers(s)

	return s
}
func (s *RealTimeService) FindNearByDrivers(ctx context.Context, location cubawheeler.GeoLocation) ([]*cubawheeler.Location, error) {
	locations, err := s.finder.FindNearByDrivers(ctx, location)
	if err != nil {
		return nil, err
	}
	return locations, nil
}

func (s *RealTimeService) NotifyToDevices(ctx context.Context, users []string, order string) error {
	return s.notifier.NotifyToDevices(ctx, users, order)
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

func notifyDrivers(s *RealTimeService) {
	ctx := context.Background()
	orders, err := s.redis.Orders(ctx)
	if err != nil {
		slog.Info("unable to get orders")
		return
	}
	ctx = cubawheeler.NewContextWithUser(ctx, &cubawheeler.User{Role: cubawheeler.RoleAdmin})

	for _, orderID := range orders {
		order, err := s.order.FindByID(ctx, orderID)
		if err != nil {
			slog.Info("unable to get order")
			continue
		}
		startPoint := cubawheeler.GeoLocation{
			Type:        "Point",
			Coordinates: []float64{order.Items.Points[0].Lng, order.Items.Points[0].Lat},
		}
		locations, err := s.finder.FindNearByDrivers(ctx, startPoint)
		if err != nil {
			slog.Info("unable to get drivers")
			continue
		}
		var users []string
		for _, l := range locations {
			users = append(users, l.User)
		}
		devices, err := s.user.GetUserDevices(ctx, users)
		if err != nil {
			slog.Info("unable to get devices")
			continue
		}
		if err := s.notifier.NotifyToDevices(ctx, devices, order.ID); err != nil {
			slog.Info("unable to notify drivers")
			continue
		}
	}
}
