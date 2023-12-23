package processor

import (
	"context"

	"cubawheeler.io/pkg/cubawheeler"
	"cubawheeler.io/pkg/currency"
)

var _ cubawheeler.PaymentMethod = (*Cash)(nil)

type Cash struct{}

func NewCash() *Cash {
	return &Cash{}
}

// Charge implements cubawheeler.PaymentMethod.
func (p *Cash) Charge(_ context.Context, amount currency.Amount) (*cubawheeler.Charge, error) {
	return cubawheeler.NewCharge(cubawheeler.ChargeStatusSucceeded, amount), nil
}

// Refund implements cubawheeler.PaymentMethod.
func (*Cash) Refund(context.Context, string, currency.Amount) (*cubawheeler.Charge, error) {
	panic("unsupported")
}
