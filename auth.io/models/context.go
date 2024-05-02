package models

import (
	"context"

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

func ClaimsFromContext(ctx context.Context) Claim {
	raw, _ := ctx.Value(oauth.ClaimsContext).(Claim)
	return raw
}

func GetTokenTypeFromContext(ctx context.Context) string {
	raw, _ := ctx.Value(oauth.TokenTypeContext).(string)
	return raw
}

func NewContextWithClient(ctx context.Context, client *Client) context.Context {
	return context.WithValue(ctx, clientCtxKey, client)
}

func ClientFromContext(ctx context.Context) *Client {
	raw, _ := ctx.Value(clientCtxKey).(*Client)
	return raw
}

func NewContextWithUser(ctx context.Context, user *User) context.Context {
	return context.WithValue(ctx, userCtxKey, user)
}

func UserFromContext(ctx context.Context) *User {
	raw, _ := ctx.Value(userCtxKey).(*User)
	return raw
}

func NewContextWithToken(ctx context.Context, token string) context.Context {
	return context.WithValue(ctx, tokenCtxKey, token)
}

func TokenFromContext(ctx context.Context) string {
	raw, _ := ctx.Value(tokenCtxKey).(string)
	return raw
}
