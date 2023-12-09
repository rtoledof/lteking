//go:generate go run github.com/99designs/gqlgen
package graph

import (
	"github.com/go-oauth2/oauth2/v4"

	"cubawheeler.io/pkg/cubawheeler"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	user     cubawheeler.UserService
	token    oauth2.TokenStore
	otp      cubawheeler.OtpService
	order    cubawheeler.OrderService
	review   cubawheeler.ReviewService
	rate     cubawheeler.RateService
	profile  cubawheeler.ProfileService
	plan     cubawheeler.PlanService
	message  cubawheeler.MessageService
	location cubawheeler.LocationService
	coupon   cubawheeler.CouponService
	client   cubawheeler.ClientService
	charge   cubawheeler.ChargeService
	ads      cubawheeler.AdsService
	vehicle  cubawheeler.VehicleService
	app      cubawheeler.ApplicationService
}
