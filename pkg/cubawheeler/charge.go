package cubawheeler

import (
	"fmt"
	"io"
	"strconv"
)

type Charge struct {
	ID                string       `json:"id"`
	Amount            int          `json:"amount"`
	Currency          string       `json:"currency"`
	Rider             *User        `json:"rider"`
	Description       string       `json:"description"`
	Trip              *Trip        `json:"trip"`
	Disputed          *bool        `json:"disputed,omitempty"`
	ReceiptEmail      string       `json:"receipt_email"`
	Status            ChargeStatus `json:"status"`
	Paid              *bool        `json:"paid,omitempty"`
	Method            *string      `json:"method,omitempty"`
	ExternalReference *string      `json:"external_reference,omitempty"`
	Fees              []*Rate      `json:"fees,omitempty"`
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
