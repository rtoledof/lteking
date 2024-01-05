package processor

import (
	"context"

	"cubawheeler.io/pkg/cubawheeler"
	"cubawheeler.io/pkg/currency"
)

var _ cubawheeler.PaymentMethod = (*Stripe)(nil)

type Stripe struct {
	apiKey string
	token  string
}

func NewStripe(apiKey, token string) *Stripe {
	return &Stripe{apiKey: apiKey, token: token}
}

// Charge implements cubawheeler.PaymentMethod.
func (p *Stripe) Charge(ctx context.Context, pm cubawheeler.ChargeMethod, amount currency.Amount) (*cubawheeler.Charge, error) {
	if pm != cubawheeler.ChargeMethodCard {
		return nil, cubawheeler.ErrInvalidInput
	}
	panic("unimplemented")
}

// Refund implements cubawheeler.PaymentMethod.
func (*Stripe) Refund(context.Context, string, currency.Amount) (*cubawheeler.Charge, error) {
	panic("unimplemented")
}
