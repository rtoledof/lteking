package order

import (
	"context"

	"github.com/go-chi/jwtauth"
	"github.com/go-chi/oauth"
)

// A private key for context that only this package can access. This is important
// to prevent collisions between different context uses
var jwtCtxKey = &contextKey{"jwt"}

type contextKey struct {
	name string
}

func NewContextWithJWT(ctx context.Context, jwt string) context.Context {
	return context.WithValue(ctx, jwtCtxKey, jwt)
}

// UserFromContext finds the user from the context. REQUIRES Middleware to have run.
func JWTFromContext(ctx context.Context) string {
	raw, _ := ctx.Value(jwtCtxKey).(string)
	return raw
}

func ClaimsFromContext(ctx context.Context) Claim {
	raw, _ := ctx.Value(oauth.ClaimsContext).(Claim)
	return raw
}

func GetTokenTypeFromContext(ctx context.Context) string {
	raw, _ := ctx.Value(oauth.TokenTypeContext).(string)
	return raw
}

func UserFromContext(ctx context.Context) *User {
	_, claim, err := jwtauth.FromContext(ctx)
	if err != nil || claim == nil {
		return nil
	}
	userData, ok := claim["user"]
	if !ok {
		return nil
	}
	var user User
	for key, v := range userData.(map[string]interface{}) {
		switch key {
		case "id":
			user.ID = v.(string)
		case "email":
			user.Email = v.(string)
		case "name":
			user.Name = v.(string)
		case "role":
			user.Role = Role(v.(string))
		}
	}
	return &user
}
