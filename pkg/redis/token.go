package redis

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"cubawheeler.io/pkg/cubawheeler"
	"cubawheeler.io/pkg/derrors"
	"github.com/go-chi/oauth"
)

var _ oauth.CredentialsVerifier = &TokenVerifier{}

type Token struct {
	AccessToken           string        `json:"access_token"`
	RefreshToken          string        `json:"refresh_token"`
	CreatedAt             time.Time     `json:"created_at"`
	AccessTokenExpiresIn  time.Duration `json:"access_token_expires_in"`
	RefreshTokenExpiresIn time.Duration `json:"refresh_token_expires_in"`
	Credentials           string        `json:"credentials"`
}

type Client struct {
	ID string `json:"id"`
}

type TokenVerifier struct {
	redis       *Redis
	user        cubawheeler.UserService
	application cubawheeler.ApplicationService
}

func NewTokenVerifier(
	redis *Redis,
	user cubawheeler.UserService,
	application cubawheeler.ApplicationService,
) *TokenVerifier {
	return &TokenVerifier{
		redis:       redis,
		user:        user,
		application: application,
	}
}

// AddClaims implements oauth.CredentialsVerifier.
func (s *TokenVerifier) AddClaims(tokenType oauth.TokenType, credential string, tokenID string, scope string, r *http.Request) (_ map[string]string, err error) {
	defer derrors.Wrap(&err, "redis.TokenVerifier.AddClaims")
	claim := make(map[string]string)

	switch tokenType {
	case oauth.BearerToken, oauth.AuthToken, oauth.UserToken:
		user, err := s.user.FindByEmail(context.Background(), credential)
		if err != nil {
			return nil, err
		}
		data, err := json.Marshal(user)
		if err != nil {
			return nil, fmt.Errorf("error marshaling user: %v: %w", err, cubawheeler.ErrInternal)
		}
		claim["user"] = string(data)
	case oauth.ClientToken:
		app, err := s.application.FindByID(context.Background(), credential)
		if err != nil {
			return nil, err
		}
		data, err := json.Marshal(app)
		if err != nil {
			return nil, fmt.Errorf("error marshaling application: %v: %w", err, cubawheeler.ErrInternal)
		}
		claim["client"] = string(data)
	}

	return claim, nil
}

// AddProperties implements oauth.CredentialsVerifier.
func (s *TokenVerifier) AddProperties(tokenType oauth.TokenType, credential string, tokenID string, scope string, r *http.Request) (_ map[string]string, err error) {
	defer derrors.Wrap(&err, "redis.TokenVerifier.AddProperties")
	return nil, nil
}

// StoreTokenID implements oauth.CredentialsVerifier.
func (s *TokenVerifier) StoreTokenID(tokenType oauth.TokenType, credential string, tokenID string, refreshTokenID string) (err error) {
	defer derrors.Wrap(&err, "redis.TokenVerifier.StoreTokenID")
	token := Token{
		AccessToken:           tokenID,
		RefreshToken:          refreshTokenID,
		CreatedAt:             time.Now().UTC(),
		AccessTokenExpiresIn:  time.Hour * 24 * 30,
		RefreshTokenExpiresIn: time.Hour * 24 * 180,
		Credentials:           credential,
	}
	return storeToken(context.Background(), s.redis, token)
}

// ValidateClient implements oauth.CredentialsVerifier.
func (s *TokenVerifier) ValidateClient(clientID string, clientSecret string, scope string, r *http.Request) (err error) {
	defer derrors.Wrap(&err, "redis.TokenVerifier.ValidateClient")
	app, err := s.application.FindByClient(context.Background(), clientID)
	if err != nil {
		if err == cubawheeler.ErrNotFound {
			app, err = s.application.FindByID(context.Background(), clientID)
		}
		if err != nil {
			return err
		}
	}
	if app.Secret != clientSecret {
		return cubawheeler.ErrAccessDenied
	}
	return nil
}

// ValidateTokenID implements oauth.CredentialsVerifier.
func (s *TokenVerifier) ValidateTokenID(tokenType oauth.TokenType, credential string, tokenID string, refreshTokenID string) error {
	token, err := getByToken(context.Background(), s.redis, tokenID)
	if err != nil {
		return err
	}
	if token.Credentials != credential || token.RefreshToken != refreshTokenID {
		return cubawheeler.ErrAccessDenied
	}

	return nil
}

// ValidateUser implements oauth.CredentialsVerifier.
func (s *TokenVerifier) ValidateUser(username string, password string, scope string, r *http.Request) (err error) {
	defer derrors.Wrap(&err, "redis.TokenVerifier.ValidateUser")
	user, err := s.user.FindByEmail(context.Background(), username)
	if err != nil {
		return err
	}
	if (user.Otp != password && bytes.Equal(user.Password, []byte(password))) || !user.IsActive() {
		return cubawheeler.ErrAccessDenied
	}

	return nil
}

func storeToken(ctx context.Context, redis *Redis, token Token) error {
	data, err := json.Marshal(token)
	if err != nil {
		slog.Info(fmt.Sprintf("error marshaling token: %v", err))
		return fmt.Errorf("error marshaling token: %w", err)
	}
	if err := redis.client.Set(ctx, token.AccessToken, data, token.AccessTokenExpiresIn).Err(); err != nil {
		slog.Info(fmt.Sprintf("error storing token: %v", err))
		return err
	}
	if err := redis.client.Set(ctx, token.RefreshToken, data, token.RefreshTokenExpiresIn).Err(); err != nil {
		slog.Info(fmt.Sprintf("error storing token: %v", err))
		return err
	}
	return nil
}

func getByToken(ctx context.Context, redis *Redis, tokenId string) (*Token, error) {
	data, err := redis.client.Get(ctx, tokenId).Result()
	if err != nil {
		return nil, fmt.Errorf("token not found: %v: %w", err, cubawheeler.ErrNotFound)
	}
	var token Token
	if err := json.Unmarshal([]byte(data), &token); err != nil {
		return nil, fmt.Errorf("unable to unmarshal token info: %w", err)
	}
	return &token, nil
}
