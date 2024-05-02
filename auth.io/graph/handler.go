package graph

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/transport"

	"auth.io/models"
)

func NewHandler(
	identity models.UserService,
	otp models.OtpService,
) *handler.Server {
	resolver := &Resolver{
		identity: identity,
		otp:      otp,
	}
	srv := handler.NewDefaultServer(NewExecutableSchema(Config{Resolvers: resolver}))
	srv.AddTransport(&transport.Websocket{})

	return srv
}
