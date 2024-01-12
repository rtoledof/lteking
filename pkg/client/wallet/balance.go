package wallet

import (
	"context"
	"net/http"

	"cubawheeler.io/pkg/cannon"
	"cubawheeler.io/pkg/cubawheeler"
)

type BalanceService service

func (s *BalanceService) Balance(ctx context.Context, currency string) (*cubawheeler.Balance, error) {
	logger := cannon.LoggerFromContext(ctx)
	logger.Info("BalanceService.Balance")
	url := "/balance"
	if currency != "" {
		url += "?currency=" + currency
	}
	req, err := s.client.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	var balance cubawheeler.Balance
	if _, err := s.client.Do(req, &balance); err != nil {
		return nil, err
	}
	return &balance, nil
}
