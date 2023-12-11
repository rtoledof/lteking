package redis

import (
	"context"
	"fmt"

	r "github.com/redis/go-redis/v9"
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
