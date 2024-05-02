//go:generate go run github.com/99designs/gqlgen generate
package graph

import (
	"auth.io/models"
	"github.com/go-chi/jwtauth"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	identity  models.UserService
	otp       models.OtpService
	tokenAuth *jwtauth.JWTAuth
}
