package graph

import (
	"cubawheeler.io/pkg/mongo"
	"cubawheeler.io/pkg/oauth"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/redis/go-redis/v9"
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
		trip:     mongo.NewTripService(db),
		vehicle:  mongo.NewVehicleService(db),
		location: mongo.NewLocationService(db),
		plan:     mongo.NewPlanService(db),
		message:  mongo.NewMessageService(db),
	}
	return handler.NewDefaultServer(NewExecutableSchema(Config{Resolvers: resolver}))
}
