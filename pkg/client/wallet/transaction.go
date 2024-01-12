package wallet

import (
	"context"
	"net/http"

	"cubawheeler.io/pkg/cannon"
	"cubawheeler.io/pkg/cubawheeler"
)

type TransactionService service

func (s *TransactionService) Transaction(ctx context.Context, currency string) (*cubawheeler.TransferEvent, error) {
	logger := cannon.LoggerFromContext(ctx)
	logger.Info("TransactionService.Transaction")
	url := "/transaction"
	if currency != "" {
		url += "?currency=" + currency
	}
	req, err := s.client.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	var transaction cubawheeler.TransferEvent
	if _, err := s.client.Do(req, &transaction); err != nil {
		return nil, err
	}
	return &transaction, nil
}
