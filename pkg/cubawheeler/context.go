package cubawheeler

import "context"

// A private key for context that only this package can access. This is important
// to prevent collisions between different context uses
var userCtxKey = &contextKey{"user"}
var clientCtxKey = &contextKey{"client"}

type contextKey struct {
	name string
}

// UserForContext finds the user from the context. REQUIRES Middleware to have run.
func UserFromContext(ctx context.Context) *User {
	raw, _ := ctx.Value(userCtxKey).(*User)
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
