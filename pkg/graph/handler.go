package graph

import (
	"github.com/99designs/gqlgen/graphql/handler"

	"cubawheeler.io/pkg/cubawheeler"
	"cubawheeler.io/pkg/mongo"
	"cubawheeler.io/pkg/redis"
)

func NewHandler(
	client *redis.Redis,
	db *mongo.DB,
	user cubawheeler.UserService,
) *handler.Server {
	resolver := &Resolver{
		user:     user,
		ads:      mongo.NewAdsService(db),
		charge:   mongo.NewChargeService(db),
		coupon:   mongo.NewCouponService(db),
		profile:  mongo.NewProfileService(db),
		order:    mongo.NewOrderService(db),
		vehicle:  mongo.NewVehicleService(db),
		location: mongo.NewLocationService(db),
		plan:     mongo.NewPlanService(db),
		message:  mongo.NewMessageService(db),
		otp:      redis.NewOtpService(client),
	}
	return handler.NewDefaultServer(NewExecutableSchema(Config{Resolvers: resolver}))
}
