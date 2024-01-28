package currency

import (
	"math"
	"strconv"
)

// Amount is a struct that contain the amount and currency together
type Amount struct {
	Amount   int64    `json:"amount,omitempty" bson:"amount,omitempty"`
	Currency Currency `json:"currency,omitempty" bson:"currency,omitempty"`
}

// Decimal Return the max amount of based of the currency
func (v Amount) Decimal() float64 {
	scale, _ := v.Currency.Rounding()
	return float64(v.Amount) * math.Pow10(-scale)
}

// ToMinAmount convert the max amount from max amount to the minimum value
func (v *Amount) ToMinAmount(amount float64, cur string) {
	v.Currency.Parse(cur)
	scale, _ := v.Currency.Rounding()
	amount = amount * math.Pow10(scale)
	amount = math.Round(amount)
	v.Amount = int64(amount)
}

func (v *Amount) Equal(a Amount) bool {
	return v.Amount == a.Amount && v.Currency.String() == a.Currency.String()
}

func (v Amount) String() string {
	scale, _ := v.Currency.Rounding()
	if scale == 0 {
		return strconv.FormatInt(v.Amount, 10)
	}
	amount := float64(v.Amount) * math.Pow10(-scale)
	return strconv.FormatFloat(amount, 'f', scale, 64)
}
