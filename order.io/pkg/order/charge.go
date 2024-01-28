package order

import (
	"fmt"
	"io"
	"strconv"
)

type ChargeMethod string

const (
	ChargeMethodCash           ChargeMethod = "Cash"
	ChargeMethodCard           ChargeMethod = "Card"
	ChargeMethodBank           ChargeMethod = "Bank"
	ChargeMethodPaypal         ChargeMethod = "Paypal"
	ChargeMethodBitcoin        ChargeMethod = "Bitcoin"
	ChargeMethodEthereum       ChargeMethod = "Ethereum"
	ChargeMethodBalance        ChargeMethod = "Balance"
	ChargeMethodCUPTransaction ChargeMethod = "CupTransaction"
	ChargeMethodMLCTransaction ChargeMethod = "MlcTransaction"
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
