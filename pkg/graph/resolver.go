//go:generate go run github.com/99designs/gqlgen
package graph

import (
	"cubawheeler.io/pkg/cubawheeler"
	"github.com/go-oauth2/oauth2/v4"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	user    cubawheeler.UserService
	token   oauth2.TokenStore
	profile cubawheeler.ProfileService
	otp     cubawheeler.OTPServer
}
