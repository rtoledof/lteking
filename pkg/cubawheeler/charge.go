package cubawheeler

import (
	"context"
	"fmt"
	"io"
	"strconv"

	"cubawheeler.io/pkg/currency"
)

type Charge struct {
	ID                string          `json:"id" bson:"_id"`
	Amount            currency.Amount `json:"amount" bson:"amount"`
	Rider             string          `json:"-" bson:"rider"`
	Description       string          `json:"description" bson:"description"`
	Order             string          `json:"-" bson:"trip"`
	Disputed          bool            `json:"disputed,omitempty" bson:"disputed"`
	ReceiptEmail      string          `json:"receipt_email" bson:"receipt_email"`
	Status            ChargeStatus    `json:"status" bson:"status"`
	Paid              *bool           `json:"paid,omitempty" bson:"paid"`
	Method            ChargeMethod    `json:"method,omitempty" bson:"method"`
	ExternalReference string          `json:"external_reference,omitempty" bson:"external_reference"`
	Fees              []*Rate         `json:"fees,omitempty" `
}

func NewCharge(status ChargeStatus, amount currency.Amount) *Charge {
	return &Charge{
		ID:     NewID().String(),
		Amount: amount,
		Status: status,
	}
}

type ChargeRequest struct {
	Limit        int
	Token        string
	Ids          []string
	Amount       int
	Currency     string
	Description  string
	Order        string
	Disputed     *bool
	ReceiptEmail string
	Status       *ChargeStatus
	Paid         *bool
	Method       ChargeMethod
	Reference    *string
	Fees         []string
	Rider        *string
	Driver       *string
	Metod        string
}

type ChargeList struct {
	Token string    `json:"token"`
	Data  []*Charge `json:"data"`
}

type ChargeService interface {
	Create(context.Context, *ChargeRequest) (*Charge, error)
	Update(context.Context, *ChargeRequest) (*Charge, error)
	FindByID(context.Context, string) (*Charge, error)
	FindAll(context.Context, ChargeRequest) (*ChargeList, error)
}

type StripeConfig struct {
	APIKey string
	Token  string
}

type PaymentmethodConfig struct {
	Stripe *StripeConfig
}

type PaymentMethod interface {
	Charge(context.Context, currency.Amount) (*Charge, error)
	Refund(context.Context, string, currency.Amount) (*Charge, error)
}

type ChargeStatus string

const (
	ChargeStatusSucceeded ChargeStatus = "SUCCEEDED"
	ChargeStatusPending   ChargeStatus = "PENDING"
	ChargeStatusFailed    ChargeStatus = "FAILED"
)

var AllChargeStatus = []ChargeStatus{
	ChargeStatusSucceeded,
	ChargeStatusPending,
	ChargeStatusFailed,
}

func (e ChargeStatus) IsValid() bool {
	switch e {
	case ChargeStatusSucceeded, ChargeStatusPending, ChargeStatusFailed:
		return true
	}
	return false
}

func (e ChargeStatus) String() string {
	return string(e)
}

func (e *ChargeStatus) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = ChargeStatus(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid ChargeStatus", str)
	}
	return nil
}

func (e ChargeStatus) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type ChargeMethod string

const (
	ChargeMethodCash           ChargeMethod = "CASH"
	ChargeMethodCard           ChargeMethod = "CARD"
	ChargeMethodBank           ChargeMethod = "BANK"
	ChargeMethodPaypal         ChargeMethod = "Paypal"
	ChargeMethodBitcoin        ChargeMethod = "Bitcoin"
	ChargeMethodEthereum       ChargeMethod = "Ethereum"
	ChargeMethodBalance        ChargeMethod = "Balance"
	ChargeMethodCUPTransaction ChargeMethod = "CUP_TRANSACTION"
	ChargeMethodMLCTransaction ChargeMethod = "MLC_TRANSACTION"
)

var AllChargeMethod = []ChargeMethod{
	ChargeMethodCash,
	ChargeMethodCard,
	ChargeMethodBank,
	ChargeMethodPaypal,
	ChargeMethodBitcoin,
	ChargeMethodEthereum,
	ChargeMethodBalance,
	ChargeMethodCUPTransaction,
	ChargeMethodMLCTransaction,
}

func (e ChargeMethod) IsValid() bool {
	switch e {
	case ChargeMethodCash,
		ChargeMethodCard,
		ChargeMethodBank,
		ChargeMethodPaypal,
		ChargeMethodBitcoin,
		ChargeMethodEthereum,
		ChargeMethodBalance,
		ChargeMethodCUPTransaction,
		ChargeMethodMLCTransaction:
		return true
	}
	return false
}

func (e ChargeMethod) String() string {
	return string(e)
}

func (e *ChargeMethod) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = ChargeMethod(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid ChargeMethod", str)
	}
	return nil
}

func (e ChargeMethod) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
