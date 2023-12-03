package cubawheeler

import (
	"context"
	"fmt"
	"io"
	"strconv"
)

type Charge struct {
	ID                string       `json:"id" bson:"_id"`
	Amount            int          `json:"amount" bson:"amount"`
	Currency          string       `json:"currency" bson:"currency"`
	Rider             string       `json:"-" bson:"rider"`
	Description       string       `json:"description" bson:"description"`
	Trip              string       `json:"-" bson:"trip"`
	Disputed          *bool        `json:"disputed,omitempty" bson:"disputed"`
	ReceiptEmail      string       `json:"receipt_email" bson:"receipt_email"`
	Status            ChargeStatus `json:"status" bson:"status"`
	Paid              *bool        `json:"paid,omitempty" bson:"paid"`
	Method            *string      `json:"method,omitempty" bson:"method"`
	ExternalReference *string      `json:"external_reference,omitempty" bson:"external_reference"`
	Fees              []*Rate      `json:"fees,omitempty" `
}

type ChargeRequest struct {
	Limit        int
	Token        string
	Ids          []string
	Amount       int
	Currency     string
	Description  string
	Trip         string
	Disputed     *bool
	ReceiptEmail string
	Status       *ChargeStatus
	Paid         *bool
	Method       *string
	Reference    *string
	Fees         []string
	Rider        *string
	Driver       *string
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
