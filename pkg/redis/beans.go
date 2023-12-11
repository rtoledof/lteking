package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"cubawheeler.io/pkg/errors"
	"cubawheeler.io/pkg/pusher"
)

type BeansToken struct {
	redis *Redis
	beans *pusher.PushNotification
}

func NewBeansToken(client *Redis, notification *pusher.PushNotification) *BeansToken {
	return &BeansToken{
		redis: client,
		beans: notification,
	}
}

func (db *BeansToken) StoreBeansToken(ctx context.Context, userID string, token map[string]any) error {
	data, err := json.Marshal(token)
	if err != nil {
		return fmt.Errorf("unable to encode token: %v: %w", err, errors.ErrInternal)
	}
	if err = db.redis.client.Set(ctx, fmt.Sprintf("BEANS_%s", userID), data, time.Hour*24).Err(); err != nil {
		return fmt.Errorf("unable to store beans token on database: %v: %w", err, errors.ErrInternal)
	}
	return nil
}

func (db *BeansToken) GetBeansToken(ctx context.Context, userID string) (map[string]any, error) {
	var token map[string]any
	data := db.redis.client.Get(ctx, fmt.Sprintf("BEANS_%s", userID))
	if data == nil {
		return nil, errors.ErrNotFound
	}
	b, err := data.Bytes()
	if err != nil {
		return nil, fmt.Errorf("unable to get token data: %v: %w", err, errors.ErrNotFound)
	}
	if err := json.Unmarshal(b, &token); err != nil {
		return nil, fmt.Errorf("unable to decode the token data: %v: %w", err, errors.ErrNotFound)
	}
	return token, nil
}
