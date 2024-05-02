package redis

import (
	"context"
	"fmt"
	"time"

	"auth.io/derrors"
	"github.com/go-chi/oauth"
)

type TokenStorer struct {
	redis *Redis
}

func NewTokenStorer(redis *Redis) *TokenStorer {
	return &TokenStorer{
		redis: redis,
	}
}

// StoreToken implements oauth.TokenStorer.
func (s *TokenStorer) StoreToken(tokenType oauth.TokenType, credential string, tokenID string, refreshTokenID string, expiresIn time.Duration, scope string, claims map[string]string, properties map[string]string) (err error) {
	defer derrors.Wrap(&err, "redis.TokenStorer.StoreToken")
	ctx := context.Background()
	key := fmt.Sprintf("%s:%s", tokenType, tokenID)
	pipe := s.redis.client.Pipeline()
	pipe.HSet(ctx, key, "credential", credential)
	pipe.HSet(ctx, key, "refresh_token", refreshTokenID)
	pipe.HSet(ctx, key, "expires_in", expiresIn.Seconds())
	pipe.HSet(ctx, key, "scope", scope)
	for k, v := range claims {
		pipe.HSet(ctx, key, k, v)
	}
	for k, v := range properties {
		pipe.HSet(ctx, key, k, v)
	}
	_, err = pipe.Exec(ctx)
	return err
}
