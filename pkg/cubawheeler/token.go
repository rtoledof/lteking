package cubawheeler

import (
	"context"
	"time"

	"github.com/go-chi/oauth"
)

type Token struct {
	AccessToken  string        `json:"access_token" bson:"access_token"`
	RefreshToken string        `json:"refresh_token,omitempty" bson:"refresh_token,omitempty"`
	ExpiresAt    Time          `json:"expires_at,omitempty" bson:"expires_at"`
	ExpiresIn    time.Duration `json:"expires_in,omitempty"  bson:"expires_in"`
	ClientID     string        `json:"client_id,omitempty" bson:"client_id,omitempty"`

	Code                  string        `json:"code,omitempty" bson:"code,omitempty"`
	CodeExpiresIn         time.Duration `json:"code_expires_in,omitempty" bson:"code_expires_in,omitempty"`
	CodeChallenge         string        `json:"code_challenge,omitempty" bson:"code_challenge,omitempty"`
	RedirectUrl           string        `json:"redirect_url,omitempty" bson:"redirect_url,omitempty"`
	RefreshTokenCreateAt  time.Time     `json:"refresh_token_create_at,omitempty" bson:"refresh_token_create_at,omitempty"`
	RefreshTokenExpiresIn time.Duration `json:"refresh_token_expires_in,omitempty" bson:"refresh_token_expires_in,omitempty"`
	Scope                 string        `json:"scope,omitempty" bson:"scope,omitempty"`
	UserID                string        `json:"user_id,omitempty" bson:"user_id,omitempty"`
}

type TokenVerifier interface {
	oauth.CredentialsVerifier
	RemoveByAccess(ctx context.Context, token string) error
}
