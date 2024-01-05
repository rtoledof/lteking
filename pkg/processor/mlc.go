package processor

import (
	"context"

	"cubawheeler.io/pkg/cubawheeler"
	"cubawheeler.io/pkg/currency"
)

var _ cubawheeler.PaymentMethod = (*MLC)(nil)

type MLC struct{}

func NewMLC() *MLC { return &MLC{} }

// Charge implements cubawheeler.PaymentMethod.
func (p *MLC) Charge(_ context.Context, pm cubawheeler.ChargeMethod, amount currency.Amount) (*cubawheeler.Charge, error) {
	return cubawheeler.NewCharge(cubawheeler.ChargeStatusSucceeded, pm, amount), nil
}

// Refund implements cubawheeler.PaymentMethod.
func (*MLC) Refund(context.Context, string, currency.Amount) (*cubawheeler.Charge, error) {
	panic("unimplemented")
}
