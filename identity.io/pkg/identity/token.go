package identity

import (
	"context"
	"encoding/base64"
	"time"

	"github.com/go-chi/oauth"
)

type Token struct {
	AccessToken           string        `json:"access_token" bson:"access_token"`
	RefreshToken          string        `json:"refresh_token,omitempty" bson:"refresh_token,omitempty"`
	ExpiresAt             Time          `json:"expires_at,omitempty" bson:"expires_at"`
	CreatedAt             Time          `json:"created_at,omitempty" bson:"created_at"`
	ExpiresIn             time.Duration `json:"expires_in,omitempty"  bson:"expires_in"`
	ClientID              ID            `json:"-" bson:"client_id,omitempty"`
	RedirectUrl           string        `json:"redirect_url,omitempty" bson:"redirect_url,omitempty"`
	RefreshTokenExpiresIn time.Duration `json:"refresh_token_expires_in,omitempty" bson:"refresh_token_expires_in,omitempty"`
	Scope                 []Scope       `json:"scope,omitempty" bson:"scope,omitempty"`
	UserID                string        `json:"-" bson:"user_id,omitempty"`
}

type TokenVerifier interface {
	oauth.CredentialsVerifier
	RemoveByAccess(ctx context.Context, token string) error
}

func NewToken(userID string, client ID) *Token {
	return &Token{
		AccessToken:           base64.RawURLEncoding.EncodeToString([]byte(NewID().String())),
		RefreshToken:          base64.RawURLEncoding.EncodeToString([]byte(NewID().String())),
		ExpiresAt:             Time{time.Now().Add(time.Hour * 24 * 30)},
		ExpiresIn:             time.Hour * 24 * 30 / time.Second,
		RefreshTokenExpiresIn: time.Hour * 24 * 90 / time.Second,
		CreatedAt:             Time{time.Now()},
		UserID:                userID,
		ClientID:              client,
	}
}
