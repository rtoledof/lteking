package order

import (
	"context"
	"encoding/json"

	"github.com/go-chi/jwtauth"
	"github.com/go-chi/oauth"
)

// A private key for context that only this package can access. This is important
// to prevent collisions between different context uses
var userCtxKey = &contextKey{"user_object"}
var clientCtxKey = &contextKey{"client"}
var jwtCtxKey = &contextKey{"jwt"}
var tokenCtxKey = &contextKey{"token"}

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

func NewContextWithToken(ctx context.Context, token *oauth.Token) context.Context {
	return context.WithValue(ctx, jwtCtxKey, token)
}

// TokenFromContext finds the user from the context. REQUIRES Middleware to have run.
func TokenFromContext(ctx context.Context) *oauth.Token {
	raw, _ := ctx.Value(tokenCtxKey).(*oauth.Token)
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
	if err := json.Unmarshal(userData.([]byte), &user); err != nil {
		return nil
	}
	return &user
}
