package mock

import (
	"context"

	"cubawheeler.io/pkg/cubawheeler"
)

var _ cubawheeler.WalletService = &WalletService{}

type WalletService struct {
	BalanceFn         func(context.Context, string) (int64, error)
	ConfirmTransferFn func(context.Context, string, string) error
	CreateFn          func(context.Context, string) (*cubawheeler.Wallet, error)
	DepositFn         func(context.Context, string, int64, string) (*cubawheeler.Wallet, error)
	FindByOwnerFn     func(context.Context, string) (*cubawheeler.Wallet, error)
	TransactionsFn    func(context.Context, string) ([]cubawheeler.TransferEvent, error)
	TransferFn        func(context.Context, string, int64, string) (*cubawheeler.TransferEvent, error)
	WithdrawFn        func(context.Context, string, int64, string) (*cubawheeler.Wallet, error)
}

// Balance implements cubawheeler.WalletService.
func (s *WalletService) Balance(ctx context.Context, owner string) (int64, error) {
	return s.BalanceFn(ctx, owner)
}

// ConfirmTransfer implements cubawheeler.WalletService.
func (s *WalletService) ConfirmTransfer(ctx context.Context, transfer, pin string) error {
	return s.ConfirmTransferFn(ctx, transfer, pin)
}

// Create implements cubawheeler.WalletService.
func (s *WalletService) Create(ctx context.Context, owner string) (*cubawheeler.Wallet, error) {
	return s.CreateFn(ctx, owner)
}

// Deposit implements cubawheeler.WalletService.
func (s *WalletService) Deposit(ctx context.Context, owner string, amount int64, currency string) (*cubawheeler.Wallet, error) {
	return s.DepositFn(ctx, owner, amount, currency)
}

// FindByOwner implements cubawheeler.WalletService.
func (s *WalletService) FindByOwner(ctx context.Context, owner string) (*cubawheeler.Wallet, error) {
	return s.FindByOwnerFn(ctx, owner)
}

// Transactions implements cubawheeler.WalletService.
func (s *WalletService) Transactions(ctx context.Context, owner string) ([]cubawheeler.TransferEvent, error) {
	return s.TransactionsFn(ctx, owner)
}

// Transfer implements cubawheeler.WalletService.
func (s *WalletService) Transfer(ctx context.Context, to string, amount int64, currency string) (*cubawheeler.TransferEvent, error) {
	return s.TransferFn(ctx, to, amount, currency)
}

// Withdraw implements cubawheeler.WalletService.
func (s *WalletService) Withdraw(ctx context.Context, owner string, amount int64, currency string) (*cubawheeler.Wallet, error) {
	return s.WithdrawFn(ctx, owner, amount, currency)
}
