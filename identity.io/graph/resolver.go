//go:generate go run github.com/99designs/gqlgen generate
package graph

import (
	"github.com/go-chi/jwtauth"
	"identity.io/pkg/identity"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	identity  identity.UserService
	otp       identity.OtpService
	tokenAuth *jwtauth.JWTAuth
}
