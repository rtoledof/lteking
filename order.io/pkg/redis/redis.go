package redis

import (
	"context"
	"fmt"

	"github.com/goccy/go-json"
	"github.com/redis/go-redis/v9"
	r "github.com/redis/go-redis/v9"

	"order.io/pkg/order"
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
		return fmt.Errorf("failed to marshal message: %v: %w", err, order.ErrInternal)
	}
	if err := db.client.Publish(ctx, channel, data).Err(); err != nil {
		return fmt.Errorf("failed to publish message to redis: %v: %w", err, order.ErrInternal)
	}
	return nil
}

func (db *Redis) Orders(ctx context.Context) ([]string, error) {
	orders, err := db.client.LRange(ctx, "order", 0, -1).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get orders from redis: %v: %w", err, order.ErrInternal)
	}
	return orders, nil
}

func (db *Redis) Subscripe(ctx context.Context, channel string) *redis.PubSub {
	return db.client.Subscribe(ctx, channel)
}
