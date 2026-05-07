package cache

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

// New creates a Redis client.
// addr format: host:port (e.g. localhost:6379)
func New(addr, password string) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       0,
	})
}

// Ping verifies the connection to Redis is alive.
func Ping(ctx context.Context, client *redis.Client) error {
	if err := client.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("redis ping: %w", err)
	}
	return nil
}
