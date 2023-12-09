package graph

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/redis/go-redis/v9"

	"cubawheeler.io/pkg/mongo"
	"cubawheeler.io/pkg/oauth"
	otp "cubawheeler.io/pkg/redis"
)

func NewHandler(
	client *redis.Client,
	db *mongo.DB,
) *handler.Server {
	resolver := &Resolver{
		user:     mongo.NewUserService(db),
		token:    oauth.NewTokenStore(client),
		ads:      mongo.NewAdsService(db),
		charge:   mongo.NewChargeService(db),
		coupon:   mongo.NewCouponService(db),
		profile:  mongo.NewProfileService(db),
		order:    mongo.NewOrderService(db),
		vehicle:  mongo.NewVehicleService(db),
		location: mongo.NewLocationService(db),
		plan:     mongo.NewPlanService(db),
		message:  mongo.NewMessageService(db),
		otp:      otp.NewOtpService(client),
	}
	return handler.NewDefaultServer(NewExecutableSchema(Config{Resolvers: resolver}))
}
