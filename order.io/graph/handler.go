package graph

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/transport"

	"order.io/pkg/order"
)

func NewHandler(
	order order.OrderService,
) *handler.Server {
	resolver := &Resolver{
		order: order,
	}
	srv := handler.NewDefaultServer(NewExecutableSchema(Config{Resolvers: resolver}))
	srv.AddTransport(&transport.Websocket{})

	return srv
}
