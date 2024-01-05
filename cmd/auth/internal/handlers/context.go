package handlers

import (
	"context"
)

var tokenCtxKey = &contextKey{"token"}

type contextKey struct {
	name string
}

func NewContextWithToken(ctx context.Context, token string) context.Context {
	return context.WithValue(ctx, tokenCtxKey, token)
}

// TokenFromContext finds the user from the context. REQUIRES Middleware to have run.
func TokenFromContext(ctx context.Context) string {
	raw, _ := ctx.Value(tokenCtxKey).(string)
	return raw
}
