package cubawheeler

import (
	"context"

	"github.com/go-chi/oauth"
)

// A private key for context that only this package can access. This is important
// to prevent collisions between different context uses
var userCtxKey = &contextKey{"user"}
var clientCtxKey = &contextKey{"client"}
var jwtCtxKey = &contextKey{"jwt"}
var tokenCtxKey = &contextKey{"token"}

type contextKey struct {
	name string
}

// UserFromContext finds the user from the context. REQUIRES Middleware to have run.
func UserFromContext(ctx context.Context) *User {
	raw, _ := ctx.Value(userCtxKey).(*User)
	return raw
}

// ClientFromContext finds the user from the context. REQUIRES Middleware to have run.
func ClientFromContext(ctx context.Context) *Application {
	raw, _ := ctx.Value(clientCtxKey).(*Application)
	return raw
}

// NewContextWithUser create a new context adding the user to it
func NewContextWithUser(ctx context.Context, user *User) context.Context {
	return context.WithValue(ctx, userCtxKey, user)
}

// NewContextWithClient create a new context adding the user to it
func NewContextWithClient(ctx context.Context, client *Application) context.Context {
	return context.WithValue(ctx, clientCtxKey, client)
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

func GetClaimsFromContext(ctx context.Context) map[string]string {
	raw, _ := ctx.Value(oauth.ClaimsContext).(map[string]string)
	return raw
}

func GetTokenTypeFromContext(ctx context.Context) string {
	raw, _ := ctx.Value(oauth.TokenTypeContext).(string)
	return raw
}
