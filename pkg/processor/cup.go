package processor

import (
	"context"

	"cubawheeler.io/pkg/cubawheeler"
	"cubawheeler.io/pkg/currency"
)

var _ cubawheeler.PaymentMethod = (*CUP)(nil)

type CUP struct {
}

func NewCUP() *CUP { return &CUP{} }

// Charge implements cubawheeler.PaymentMethod.
func (p *CUP) Charge(_ context.Context, pm cubawheeler.ChargeMethod, amount currency.Amount) (*cubawheeler.Charge, error) {
	return cubawheeler.NewCharge(cubawheeler.ChargeStatusSucceeded, pm, amount), nil
}

// Refund implements cubawheeler.PaymentMethod.
func (*CUP) Refund(context.Context, string, currency.Amount) (*cubawheeler.Charge, error) {
	panic("unimplemented")
}
