package mock

import (
	"context"

	"github.com/go-oauth2/oauth2/v4"
)

var _ oauth2.TokenStore = (*TokenStore)(nil)

type TokenStore struct {
	CreateFn          func(info oauth2.TokenInfo) error
	GetByAccessFn     func(access string) (oauth2.TokenInfo, error)
	GetByCodeFn       func(code string) (oauth2.TokenInfo, error)
	GetByRefreshFn    func(refresh string) (oauth2.TokenInfo, error)
	RemoveByAccessFn  func(access string) error
	RemoveByCodeFn    func(code string) error
	RemoveByRefreshFn func(refresh string) error
}

// Create implements oauth2.TokenStore.
func (s *TokenStore) Create(_ context.Context, info oauth2.TokenInfo) error {
	return s.CreateFn(info)
}

// GetByAccess implements oauth2.TokenStore.
func (s *TokenStore) GetByAccess(_ context.Context, access string) (oauth2.TokenInfo, error) {
	return s.GetByAccessFn(access)
}

// GetByCode implements oauth2.TokenStore.
func (s *TokenStore) GetByCode(_ context.Context, code string) (oauth2.TokenInfo, error) {
	return s.GetByCodeFn(code)
}

// GetByRefresh implements oauth2.TokenStore.
func (s *TokenStore) GetByRefresh(_ context.Context, refresh string) (oauth2.TokenInfo, error) {
	return s.GetByRefreshFn(refresh)
}

// RemoveByAccess implements oauth2.TokenStore.
func (s *TokenStore) RemoveByAccess(_ context.Context, access string) error {
	return s.RemoveByAccessFn(access)
}

// RemoveByCode implements oauth2.TokenStore.
func (s *TokenStore) RemoveByCode(_ context.Context, code string) error {
	return s.RemoveByCodeFn(code)
}

// RemoveByRefresh implements oauth2.TokenStore.
func (s *TokenStore) RemoveByRefresh(ctx context.Context, refresh string) error {
	return s.RemoveByRefreshFn(refresh)
}
