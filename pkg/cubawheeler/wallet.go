package cubawheeler

import (
	"context"
	"time"
)

type Wallet struct {
	ID               string          `json:"id" bson:"_id"`
	Owner            string          `json:"owner" bson:"owner"`
	Balance          Balance         `json:"balance" bson:"balance"`
	Currency         string          `json:"currency" bson:"currency"`
	CreatedAt        uint            `json:"-" bson:"created_at"`
	UpdatedAt        uint            `json:"updated_at" bson:"updated_at"`
	Events           []interface{}   `json:"-" bson:"events"`
	TransferEvent    []TransferEvent `json:"-" bson:"transfer_event"`
	PendingTransfers []TransferEvent `json:"-" bson:"pending_transfers"`
}

func (w *Wallet) CanWithdraw(amount int64, currency string) bool {
	return w.Balance.Amount[currency]-amount >= 0
}

func (w *Wallet) CanTransfer(amount int64, currency string) bool {
	return w.Balance.Amount[currency]-amount >= 0
}

func (w *Wallet) FindPendingTransfer(id string) (*TransferEvent, int) {
	for i, t := range w.PendingTransfers {
		if t.ID == id {
			return &t, i
		}
	}
	return nil, 0
}

func (w *Wallet) Deposit(amount int64, currency string) {
	w.Balance.Amount[currency] += amount
	w.UpdatedAt = uint(time.Now().Unix())
	w.Events = append(w.Events, DepositEvent{
		Amount:    amount,
		Currency:  currency,
		CreatedAt: uint(time.Now().Unix()),
	})
	w.TransferEvent = append(w.TransferEvent, TransferEvent{
		Type:      TransferTypeDeposit,
		Amount:    amount,
		Currency:  currency,
		CreatedAt: uint(time.Now().Unix()),
	})
}

func (w *Wallet) Withdraw(amount int64, currency string) {
	w.Balance.Amount[currency] -= amount
	w.UpdatedAt = uint(time.Now().Unix())
	w.Events = append(w.Events, WithdrawEvent{
		Amount:    amount,
		Currency:  currency,
		CreatedAt: uint(time.Now().Unix()),
	})
	w.TransferEvent = append(w.TransferEvent, TransferEvent{
		Type:      TransferTypeWithdraw,
		Amount:    amount,
		Currency:  currency,
		CreatedAt: uint(time.Now().Unix()),
	})
}

func NewWallet() *Wallet {
	return &Wallet{
		ID:        NewID().String(),
		Balance:   Balance{Amount: make(map[string]int64)},
		CreatedAt: uint(time.Now().Unix()),
		UpdatedAt: uint(time.Now().Unix()),
	}
}

type DepositEvent struct {
	Amount    int64  `json:"amount" bson:"amount"`
	Currency  string `json:"currency" bson:"currency"`
	CreatedAt uint   `json:"-" bson:"created_at"`
}

type WithdrawEvent struct {
	Amount    int64  `json:"amount" bson:"amount"`
	Currency  string `json:"currency" bson:"currency"`
	CreatedAt uint   `json:"-" bson:"created_at"`
}

type TransferType string

const (
	TransferTypeDeposit  TransferType = "deposit"
	TransferTypeWithdraw TransferType = "withdraw"
	TransferTypeTransfer TransferType = "transfer"
)

type TransferStatus int

const (
	TransferStatusPending TransferStatus = iota
	TransferStatusConfirmed
	TransferStatusFailed
	TransferStatusCancelled
)

func (s TransferStatus) IsValid() bool {
	return s == TransferStatusPending ||
		s == TransferStatusConfirmed ||
		s == TransferStatusFailed ||
		s == TransferStatusCancelled
}

func (s TransferStatus) String() string {
	switch s {
	case TransferStatusPending:
		return "PENDING"
	case TransferStatusConfirmed:
		return "CONFIRMED"
	case TransferStatusFailed:
		return "FAILED"
	case TransferStatusCancelled:
		return "CANCELLED"
	default:
		return "UNKNOWN"
	}
}

func (s *TransferStatus) UnmarshalJSON(b []byte) error {
	switch string(b) {
	case `"PENDING"`:
		*s = TransferStatusPending
	case `"CONFIRMED"`:
		*s = TransferStatusConfirmed
	case `"FAILED"`:
		*s = TransferStatusFailed
	case `"CANCELLED"`:
		*s = TransferStatusCancelled
	default:
		return ErrInvalidInput
	}
	return nil
}

func (s TransferStatus) MarshalJSON() ([]byte, error) {
	switch s {
	case TransferStatusPending:
		return []byte(`"PENDING"`), nil
	case TransferStatusConfirmed:
		return []byte(`"CONFIRMED"`), nil
	case TransferStatusFailed:
		return []byte(`"FAILED"`), nil
	case TransferStatusCancelled:
		return []byte(`"CANCELLED"`), nil
	default:
		return nil, ErrInvalidInput
	}
}

type TransferEvent struct {
	ID        string         `json:"id,omitempty" bson:"_id,omitempty"`
	From      string         `json:"from,omitempty" bson:"from,omitempty"`
	To        string         `json:"to,omitempty" bson:"to,omitempty"`
	Type      TransferType   `json:"type" bson:"type"`
	Status    TransferStatus `json:"status" bson:"status"`
	Amount    int64          `json:"amount" bson:"amount"`
	Currency  string         `json:"currency" bson:"currency"`
	CreatedAt uint           `json:"created_at" bson:"created_at"`
}

type WalletService interface {
	Create(context.Context, string) (*Wallet, error)
	Deposit(context.Context, string, int64, string) (*Wallet, error)
	Withdraw(context.Context, string, int64, string) (*Wallet, error)
	// TODO: add confirmation to transfer
	Transfer(context.Context, string, int64, string) (*TransferEvent, error)
	ConfirmTransfer(context.Context, string, string) error
	FindByOwner(context.Context, string) (*Wallet, error)
	Balance(context.Context, string) (Balance, error)
	Transactions(context.Context, string) ([]TransferEvent, error)
}
