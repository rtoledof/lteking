package mock

import (
	"context"
	"net/http"

	"github.com/go-chi/oauth"

	"auth.io/models"
)

var _ models.TokenVerifier = &TokenVerifier{}

type TokenVerifier struct {
	AddClaimsFn       func(oauth.TokenType, string, string, string, *http.Request) (map[string]string, error)
	AddPropertiesFn   func(oauth.TokenType, string, string, string, *http.Request) (map[string]string, error)
	RemoveByAccessFn  func(context.Context, string) error
	StoreTokenIDFn    func(oauth.TokenType, string, string, string) error
	ValidateClientFn  func(string, string, string, *http.Request) error
	ValidateTokenIDFn func(oauth.TokenType, string, string, string) error
	ValidateUserFn    func(string, string, string, *http.Request) error
}

// AddClaims implements models.TokenVerifier.
func (s *TokenVerifier) AddClaims(tokenType oauth.TokenType, credential string, tokenID string, scope string, r *http.Request) (map[string]string, error) {
	return s.AddClaimsFn(tokenType, credential, tokenID, scope, r)
}

// AddProperties implements models.TokenVerifier.
func (s *TokenVerifier) AddProperties(tokenType oauth.TokenType, credential string, tokenID string, scope string, r *http.Request) (map[string]string, error) {
	return s.AddPropertiesFn(tokenType, credential, tokenID, scope, r)
}

// RemoveByAccess implements models.TokenVerifier.
func (s *TokenVerifier) RemoveByAccess(ctx context.Context, token string) error {
	return s.RemoveByAccessFn(ctx, token)
}

// StoreTokenID implements models.TokenVerifier.
func (s *TokenVerifier) StoreTokenID(tokenType oauth.TokenType, credential string, tokenID string, refreshTokenID string) error {
	return s.StoreTokenIDFn(tokenType, credential, tokenID, refreshTokenID)
}

// ValidateClient implements models.TokenVerifier.
func (s *TokenVerifier) ValidateClient(clientID string, clientSecret string, scope string, r *http.Request) error {
	return s.ValidateClientFn(clientID, clientSecret, scope, r)
}

// ValidateTokenID implements models.TokenVerifier.
func (s *TokenVerifier) ValidateTokenID(tokenType oauth.TokenType, credential string, tokenID string, refreshTokenID string) error {
	return s.ValidateTokenIDFn(tokenType, credential, tokenID, refreshTokenID)
}

// ValidateUser implements models.TokenVerifier.
func (s *TokenVerifier) ValidateUser(username string, password string, scope string, r *http.Request) error {
	return s.ValidateUserFn(username, password, scope, r)
}
