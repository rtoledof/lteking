package graph

import (
	"github.com/99designs/gqlgen/graphql/handler"

	"cubawheeler.io/pkg/ably"
	"cubawheeler.io/pkg/cubawheeler"
	"cubawheeler.io/pkg/mongo"
	"cubawheeler.io/pkg/realtime"
	"cubawheeler.io/pkg/redis"
)

func NewHandler(
	client *redis.Redis,
	db *mongo.DB,
	user cubawheeler.UserService,
	exit chan struct{},
	connectionString string,
	abyKey string,
) *handler.Server {
	resolver := &Resolver{
		user:       user,
		ads:        mongo.NewAdsService(db),
		charge:     mongo.NewChargeService(db),
		coupon:     mongo.NewCouponService(db),
		profile:    mongo.NewProfileService(db),
		order:      mongo.NewOrderService(db),
		vehicle:    mongo.NewVehicleService(db),
		location:   mongo.NewLocationService(db),
		plan:       mongo.NewPlanService(db),
		message:    mongo.NewMessageService(db),
		otp:        redis.NewOtpService(client),
		ablyClient: ably.NewClient(connectionString, exit, abyKey),
	}
	resolver.realTimeLocation = realtime.NewRealTimeService(
		redis.NewRealTimeService(client),
		resolver.ablyClient.Notifier,
		resolver.user,
	)
	return handler.NewDefaultServer(NewExecutableSchema(Config{Resolvers: resolver}))
}
