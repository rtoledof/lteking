//go:generate go run github.com/99designs/gqlgen generate
package graph

import "wallet.io/pkg/wallet"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	wallet wallet.WalletService
}
