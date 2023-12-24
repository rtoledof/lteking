package cubawheeler

import (
	"context"
	"time"
)

type Wallet struct {
	ID        string        `json:"id" bson:"_id"`
	Owner     string        `json:"owner" bson:"owner"`
	Balance   int64         `json:"balance" bson:"balance"`
	CreatedAt uint          `json:"-" bson:"created_at"`
	UpdatedAt uint          `json:"updated_at" bson:"updated_at"`
	Events    []interface{} `json:"-" bson:"events"`
}

func (w *Wallet) Deposit(amount int64) {
	w.Balance += amount
	w.UpdatedAt = uint(time.Now().Unix())
	w.Events = append(w.Events, DepositEvent{
		Amount:    amount,
		CreatedAt: uint(time.Now().Unix()),
	})
}

func (w *Wallet) Withdraw(amount int64) {
	w.Balance -= amount
	w.UpdatedAt = uint(time.Now().Unix())
	w.Events = append(w.Events, WithdrawEvent{
		Amount:    amount,
		CreatedAt: uint(time.Now().Unix()),
	})
}

func NewWallet() *Wallet {
	return &Wallet{
		ID:        NewID().String(),
		Balance:   0,
		CreatedAt: uint(time.Now().Unix()),
		UpdatedAt: uint(time.Now().Unix()),
	}
}

type DepositEvent struct {
	Amount    int64 `json:"amount" bson:"amount"`
	CreatedAt uint  `json:"-" bson:"created_at"`
}

type WithdrawEvent struct {
	Amount    int64 `json:"amount" bson:"amount"`
	CreatedAt uint  `json:"-" bson:"created_at"`
}

type WalletService interface {
	Create(context.Context, string) (*Wallet, error)
	Deposit(context.Context, string, int64) (*Wallet, error)
	Withdraw(context.Context, string, int64) (*Wallet, error)
	Transfer(context.Context, string, string, int64) (*Wallet, *Wallet, error)
	FindByOwner(context.Context, string) (*Wallet, error)
	Balance(context.Context, string) (int64, error)
}
