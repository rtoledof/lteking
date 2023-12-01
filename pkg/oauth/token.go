package oauth

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-oauth2/oauth2/v4"
	"github.com/redis/go-redis/v9"

	"cubawheeler.io/pkg/cubawheeler"
)

var _ oauth2.TokenStore = &TokenStore{}

type TokenStore struct {
	redis *redis.Client
}

func NewTokenStore(client *redis.Client) *TokenStore {
	return &TokenStore{redis: client}
}

func (s *TokenStore) Create(ctx context.Context, info oauth2.TokenInfo) error {
	return storeToken(ctx, s.redis, info)
}

func (s *TokenStore) RemoveByCode(ctx context.Context, code string) error {
	token, err := getByToken(ctx, s.redis, code)
	if err != nil {
		return err
	}
	return removeToken(ctx, s.redis, token)
}

func (s *TokenStore) RemoveByAccess(ctx context.Context, access string) error {
	token, err := s.GetByAccess(ctx, access)
	if err != nil {
		return err
	}
	return removeToken(ctx, s.redis, token)
}

func (s *TokenStore) RemoveByRefresh(ctx context.Context, refresh string) error {
	token, err := getByToken(ctx, s.redis, refresh)
	if err != nil {
		return err
	}
	return removeToken(ctx, s.redis, token)
}

func (s *TokenStore) GetByCode(ctx context.Context, code string) (oauth2.TokenInfo, error) {
	return getByToken(ctx, s.redis, code)
}

func (s *TokenStore) GetByAccess(ctx context.Context, access string) (oauth2.TokenInfo, error) {
	data, err := s.redis.Get(ctx, access).Bytes()
	if err != nil {
		return nil, fmt.Errorf("token not found: %w", err)
	}
	var token cubawheeler.Token
	if err := json.Unmarshal(data, &token); err != nil {
		return nil, fmt.Errorf("unable to unmarshal token info: %w", err)
	}
	return &token, nil
}

func (s *TokenStore) GetByRefresh(ctx context.Context, refresh string) (oauth2.TokenInfo, error) {
	return getByToken(ctx, s.redis, refresh)
}

func storeToken(ctx context.Context, client *redis.Client, info oauth2.TokenInfo) error {
	data, err := json.Marshal(info)
	if err != nil {
		return fmt.Errorf("unable to marshal token info: %w", err)
	}
	if err := client.Set(ctx, info.GetAccess(), data, info.GetAccessExpiresIn()).Err(); err != nil {
		return fmt.Errorf("unable to store token info %w", err)
	}
	if err := client.Set(ctx, info.GetCode(), info.GetAccess(), info.GetCodeExpiresIn()).Err(); err != nil {
		return fmt.Errorf("unable to store token by code: %w", err)
	}
	if err := client.Set(ctx, info.GetRefresh(), info.GetAccess(), info.GetRefreshExpiresIn()).Err(); err != nil {
		return fmt.Errorf("unable to store token by refresh: %w", err)
	}

	return nil
}

func removeToken(ctx context.Context, client *redis.Client, info oauth2.TokenInfo) error {
	if err := client.Del(ctx, []string{info.GetCode(), info.GetRefresh(), info.GetAccess()}...).Err(); err != nil {
		return fmt.Errorf("unable to remove token from redis: %w", err)
	}
	return nil
}

func getByToken(ctx context.Context, client *redis.Client, token string) (oauth2.TokenInfo, error) {
	accessToken := client.Get(ctx, token).String()
	if len(accessToken) == 0 {
		return nil, fmt.Errorf("unable to find the token info by code")
	}
	data, err := client.Get(ctx, accessToken).Bytes()
	if err != nil {
		return nil, fmt.Errorf("token not found: %w", err)
	}
	var info cubawheeler.Token
	if err := json.Unmarshal(data, &token); err != nil {
		return nil, fmt.Errorf("unable to unmarshal token info: %w", err)
	}
	return &info, nil
}
