package graph

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/go-oauth2/oauth2/v4"

	"cubawheeler.io/pkg/cubawheeler"
	"cubawheeler.io/pkg/database"
)

func NewHandler(
	token oauth2.TokenStore,
	user cubawheeler.UserService,
) *handler.Server {
	resolver := &Resolver{
		user:    user,
		token:   token,
		profile: database.NewProfileService(database.Db),
	}
	return handler.NewDefaultServer(NewExecutableSchema(Config{Resolvers: resolver}))
}
