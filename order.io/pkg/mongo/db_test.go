package mongo

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/go-chi/jwtauth"
	"github.com/lestrrat-go/jwx/jwt"

	"order.io/pkg/order"
)

func NewTestDB() *DB {
	return NewDB("mongodb://localhost:27017", "test")
}

func prepareContext(t *testing.T, roles ...order.Role) context.Context {
	t.Helper()
	ctx := context.Background()

	token := jwt.New()
	token.Set("id", order.NewID().String())
	user := order.User{
		ID:    order.NewID().String(),
		Name:  "test",
		Email: "test",
		Role:  "rider",
	}
	if roles != nil {
		user.Role = roles[0]
	}
	userData, _ := json.Marshal(user)
	token.Set("user", userData)
	return jwtauth.NewContext(ctx, token, nil)
}
