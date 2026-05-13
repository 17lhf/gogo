package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"gogo/internal/pkg"
)

// SessionCache manages JWT session storage in Redis.
type SessionCache struct {
	rdb *redis.Client
	ttl time.Duration
}

// NewSessionCache creates a new SessionCache.
func NewSessionCache(rdb *redis.Client, ttl time.Duration) *SessionCache {
	return &SessionCache{rdb: rdb, ttl: ttl}
}

func (s *SessionCache) key(userID int64, jti string) string {
	return fmt.Sprintf("session:%d:%s", userID, jti)
}

// Set stores a session in Redis with TTL.
func (s *SessionCache) Set(ctx context.Context, userID int64, jti string, claims *pkg.Claims) error {
	data, err := json.Marshal(claims)
	if err != nil {
		return err
	}
	return s.rdb.Set(ctx, s.key(userID, jti), data, s.ttl).Err()
}

// Get retrieves a session from Redis. Returns nil if not found.
func (s *SessionCache) Get(ctx context.Context, userID int64, jti string) (*pkg.Claims, error) {
	data, err := s.rdb.Get(ctx, s.key(userID, jti)).Bytes()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var claims pkg.Claims
	if err := json.Unmarshal(data, &claims); err != nil {
		return nil, err
	}
	return &claims, nil
}

// Delete removes a session from Redis (logout).
func (s *SessionCache) Delete(ctx context.Context, userID int64, jti string) error {
	return s.rdb.Del(ctx, s.key(userID, jti)).Err()
}
