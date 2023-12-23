package processor

import (
	"context"

	"cubawheeler.io/pkg/cubawheeler"
	"cubawheeler.io/pkg/currency"
)

var _ cubawheeler.PaymentMethod = (*Charge)(nil)

type Charge struct {
	config cubawheeler.PaymentmethodConfig
}

func NewCharge(pm cubawheeler.PaymentmethodConfig) *Charge {
	return &Charge{config: pm}
}

// Charge implements cubawheeler.PaymentMethod.
func (p *Charge) Charge(ctx context.Context, amount currency.Amount) (*cubawheeler.Charge, error) {
	charger, err := ChargeRefund(ctx, p.config, cubawheeler.ChargeMethodCard)
	if err != nil {
		return nil, err
	}
	return charger.Charge(ctx, amount)
}

// Refund implements cubawheeler.PaymentMethod.
func (p *Charge) Refund(ctx context.Context, transID string, amount currency.Amount) (*cubawheeler.Charge, error) {
	refunder, err := ChargeRefund(ctx, p.config, cubawheeler.ChargeMethodCard)
	if err != nil {
		return nil, err
	}
	return refunder.Refund(ctx, transID, amount)
}

func ChargeRefund(
	ctx context.Context,
	config cubawheeler.PaymentmethodConfig,
	method cubawheeler.ChargeMethod,
) (cubawheeler.PaymentMethod, error) {
	switch method {
	case cubawheeler.ChargeMethodCard:
		return NewStripe(config.Stripe.APIKey, config.Stripe.Token), nil
	default:
		return nil, cubawheeler.ErrInvalidInput
	}
}
