//go:generate go run github.com/99designs/gqlgen
package graph

import "cubawheeler.io/pkg/uploader"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	orderService  string
	authService   string
	walletService string

	uploader uploader.Uploader
}

func NewResolver(orderService, authService, walletService string) *Resolver {
	return &Resolver{
		orderService:  orderService,
		authService:   authService,
		walletService: walletService,
	}
}
