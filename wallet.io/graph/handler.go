package graph

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/transport"

	"wallet.io/pkg/wallet"
)

func NewHandler(
	wallet wallet.WalletService,
) *handler.Server {
	resolver := &Resolver{
		wallet: wallet,
	}
	srv := handler.NewDefaultServer(NewExecutableSchema(Config{Resolvers: resolver}))
	srv.AddTransport(&transport.Websocket{})

	return srv
}
