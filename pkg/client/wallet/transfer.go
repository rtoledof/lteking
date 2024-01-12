package wallet

import (
	"context"
	"net/http"
	"net/url"

	"cubawheeler.io/pkg/cannon"
	"cubawheeler.io/pkg/cubawheeler"
)

type Transfer struct {
	From     string `url:"from"`
	To       string `url:"to"`
	Amount   int64  `url:"amount"`
	Currency string `url:"currency"`
}

type TransferService service

func (s *TransferService) Transfer(ctx context.Context, transfer *Transfer) (*cubawheeler.TransferEvent, error) {
	logger := cannon.LoggerFromContext(ctx)
	logger.Info("TransferService.Transfer")
	value := url.Values{
		"from":     []string{transfer.From},
		"to":       []string{transfer.To},
		"amount":   []string{string(transfer.Amount)},
		"currency": []string{transfer.Currency},
	}
	req, err := s.client.NewRequest(http.MethodPost, "/transfer", value)
	if err != nil {
		return nil, err
	}
	var transaction cubawheeler.TransferEvent
	if _, err := s.client.Do(req, &transaction); err != nil {
		return nil, err
	}
	return &transaction, nil
}
