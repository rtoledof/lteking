package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/goccy/go-json"
	"github.com/redis/go-redis/v9"
	r "github.com/redis/go-redis/v9"
	"identity.io/pkg/identity"
)

type Redis struct {
	client *r.Client
}

func NewRedis(client *r.Client) *Redis {
	return &Redis{client: client}
}

func (db *Redis) Ping(ctx context.Context) error {
	err := db.client.Ping(ctx).Err()
	if err != nil {
		return fmt.Errorf("failed to connect to redis: %w", err)
	}
	return nil
}

func (db *Redis) Close() error {
	if err := db.client.Close(); err != nil {
		fmt.Println("failed to close redis", err)
	}
	return nil
}

func (db *Redis) Publish(ctx context.Context, channel string, message interface{}) error {
	data, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %v: %w", err, identity.ErrInternal)
	}
	if err := db.client.Publish(ctx, channel, data).Err(); err != nil {
		return fmt.Errorf("failed to publish message to redis: %v: %w", err, identity.ErrInternal)
	}
	return nil
}

func (db *Redis) Orders(ctx context.Context) ([]string, error) {
	orders, err := db.client.LRange(ctx, "order", 0, -1).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get orders from redis: %v: %w", err, identity.ErrInternal)
	}
	return orders, nil
}

func (db *Redis) Order(ctx context.Context) ([]string, error) {
	order, err := db.client.LRange(ctx, "order:confirmed", 0, -1).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get orders from redis: %v: %w", err, identity.ErrInternal)
	}
	return order, nil
}

func (db *Redis) Subscripe(ctx context.Context, channel string) *redis.PubSub {
	return db.client.Subscribe(ctx, channel)
}

func (db *Redis) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	if err := db.client.Set(ctx, key, value, expiration).Err(); err != nil {
		return fmt.Errorf("failed to set key: %v: %w", err, identity.ErrInternal)
	}
	return nil
}

func (db *Redis) Get(ctx context.Context, key string) (string, error) {
	data := db.client.Get(ctx, key)
	if data == nil {
		return "", identity.ErrNotFound
	}
	value, err := data.Result()
	if err != nil {
		return "", fmt.Errorf("failed to get key: %v: %w", err, identity.ErrInternal)
	}
	return value, nil
}

func (db *Redis) Del(ctx context.Context, key string) error {
	if err := db.client.Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("failed to delete key: %v: %w", err, identity.ErrInternal)
	}
	return nil
}
