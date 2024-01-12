package wallet

import (
	"context"
	"net/http"
	"net/url"

	"cubawheeler.io/pkg/cannon"
	"cubawheeler.io/pkg/cubawheeler"
)

type WalletService service

type CreateRequest struct {
	Owner string `url:"owner"`
}

func (s *WalletService) Create(ctx context.Context, owner CreateRequest) (*cubawheeler.Wallet, error) {
	logger := cannon.LoggerFromContext(ctx)
	logger.Info("WalletService.Create")
	value := url.Values{
		"owner": []string{owner.Owner},
	}
	url := "/v1/wallet"

	req, err := s.client.NewRequest(http.MethodPost, url, value)
	if err != nil {
		return nil, err
	}
	var w cubawheeler.Wallet
	if _, err := s.client.Do(req, &w); err != nil {
		return nil, err
	}
	logger.Info("WalletService.Create", "wallet", w)
	return &w, nil
}
