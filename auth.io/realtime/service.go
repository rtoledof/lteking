package realtime

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/ably/ably-go/ably"
	"github.com/go-chi/jwtauth"
	"github.com/lestrrat-go/jwx/jwt"

	"auth.io/models"
	"auth.io/redis"
)

var (
	DriverLocations        = make(chan models.Location, 10000)
	UserAvailabilityStatus = make(chan UserStatus, 1000)
)

type UserStatus struct {
	User      string
	Available bool
}

type Finder interface {
	FindNearByDrivers(context.Context, models.GeoLocation) ([]*models.Location, error)
}

type Updater interface {
	UpdateLocation(context.Context, string, models.GeoLocation) error
}

type FinderUpdater interface {
	Finder
	Updater
}
type OrderNotification struct {
	ID       string `json:"id"`
	Cost     int64  `json:"cost"`
	Currency string `json:"currency"`
	Distance int64  `json:"distance"`
	Duration int64  `json:"duration"`
}

func AssambleOrderNotification(order *models.Order) OrderNotification {
	return OrderNotification{
		ID:       order.ID,
		Cost:     int64(order.Price),
		Currency: order.Currency,
		Distance: int64(order.Distance),
		Duration: int64(order.Duration),
	}
}

type Notifier interface {
	NotifyToDevices(context.Context, []string, OrderNotification, *ably.Realtime, *ably.REST) error
	NotifyRiderOrderAccepted(context.Context, []string, OrderNotification) error
}

type UserUpdateService interface {
	SetAvailability(context.Context, string, bool) error
}

type RealTimeService struct {
	finder       FinderUpdater
	notifier     Notifier
	user         models.UserService
	redis        *redis.Redis
	ablyRealTime *ably.Realtime
	rest         *ably.REST
}

func NewRealTimeService(
	finder FinderUpdater,
	notifier Notifier,
	redis *redis.Redis,
	ablyRealTime *ably.Realtime,
	user models.UserService,
) *RealTimeService {

	s := &RealTimeService{
		finder:       finder,
		notifier:     notifier,
		redis:        redis,
		ablyRealTime: ablyRealTime,
		user:         user,
	}

	go storeOrUpdateDriversLocation(finder)
	go processNewOrder(s)
	go updateUserStatus(user)
	go notifyDrivers(s)

	return s
}
func (s *RealTimeService) FindNearByDrivers(ctx context.Context, location models.GeoLocation) ([]*models.Location, error) {
	locations, err := s.finder.FindNearByDrivers(ctx, location)
	if err != nil {
		return nil, err
	}
	return locations, nil
}

func (s *RealTimeService) NotifyToDevices(ctx context.Context, users []string, order OrderNotification, realTime *ably.Realtime, rest *ably.REST) error {
	return s.notifier.NotifyToDevices(ctx, users, order, realTime, rest)
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

func updateUserStatus(service models.UserService) {
	ctx := adminContext()
	for v := range UserAvailabilityStatus {
		usr, err := service.FindByID(ctx, v.User)
		if err != nil {
			slog.Info("unable to find user")
			continue
		}
		usr.Available = v.Available

		if err := service.Update(ctx, usr); err != nil {
			slog.Info("unable to update use availability")
		}
	}
}

func notifyDrivers(s *RealTimeService) {
	ctx := adminContext()
	pubsub := s.redis.Subscripe(ctx, "orders")
	defer pubsub.Close()
	for {
		msg, err := pubsub.ReceiveMessage(ctx)
		if err != nil {
			slog.Info("unable to receive message", "%v", err)
			continue
		}

		var order models.Order
		if err := json.Unmarshal([]byte(msg.Payload), &order); err != nil {
			slog.Info("unable to get order")
			continue
		}

		startPoint := models.GeoLocation{
			Type:        "Point",
			Coordinates: []float64{order.Item.Points[0].Lng, order.Item.Points[0].Lat},
			Lat:         order.Item.Points[0].Lat,
			Long:        order.Item.Points[0].Lng,
		}
		locations, err := s.finder.FindNearByDrivers(ctx, startPoint)
		if err != nil || len(locations) == 0 {
			slog.Info("unable to get drivers")
			continue
		}
		var users []string
		for _, l := range locations {
			if _, ok := order.BannedDrivers[l.User]; !ok {
				users = append(users, l.User)
			}
		}
		devices, err := s.user.GetUserDevices(ctx, models.UserFilter{
			Ids:  users,
			Role: models.RoleDriver,
		})
		if err != nil || len(devices) == 0 {
			slog.Info("unable to get devices")
			continue
		}
		if err := s.notifier.NotifyToDevices(ctx, devices, AssambleOrderNotification(&order), s.ablyRealTime, s.rest); err != nil {
			slog.Info("unable to notify drivers")
			continue
		}
	}
}

func adminContext() context.Context {
	ctx := context.Background()
	token := jwt.New()
	token.Set("id", models.NewID().String())
	user := models.User{
		ID:   models.NewID().String(),
		Role: models.RoleAdmin,
	}
	userData, _ := json.Marshal(user)
	token.Set("user", userData)
	return jwtauth.NewContext(ctx, token, nil)
}
