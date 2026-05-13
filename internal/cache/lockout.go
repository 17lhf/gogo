package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// LockoutCache manages login failure tracking to prevent brute-force attacks.
type LockoutCache struct {
	rdb       *redis.Client
	threshold int
	duration  time.Duration
}

// NewLockoutCache creates a new LockoutCache.
func NewLockoutCache(rdb *redis.Client, threshold int, duration time.Duration) *LockoutCache {
	return &LockoutCache{rdb: rdb, threshold: threshold, duration: duration}
}

func (l *LockoutCache) failKey(username string) string {
	return fmt.Sprintf("login_fail:%s", username)
}

func (l *LockoutCache) lockKey(username string) string {
	return fmt.Sprintf("login_lock:%s", username)
}

// RecordFailure increments the failure counter for a username.
// If the counter reaches the threshold, sets a lock key.
func (l *LockoutCache) RecordFailure(ctx context.Context, username string) (bool, error) {
	key := l.failKey(username)
	count, err := l.rdb.Incr(ctx, key).Result()
	if err != nil {
		return false, err
	}
	// Set expiry on the counter key
	l.rdb.Expire(ctx, key, l.duration)

	if count >= int64(l.threshold) {
		// Lock the account
		l.rdb.Set(ctx, l.lockKey(username), "1", l.duration)
		return true, nil
	}
	return false, nil
}

// IsLocked checks if an account is currently locked.
func (l *LockoutCache) IsLocked(ctx context.Context, username string) (bool, error) {
	exists, err := l.rdb.Exists(ctx, l.lockKey(username)).Result()
	if err != nil {
		return false, err
	}
	return exists > 0, nil
}

// Reset clears all lockout state for a username (on successful login).
func (l *LockoutCache) Reset(ctx context.Context, username string) error {
	pipe := l.rdb.Pipeline()
	pipe.Del(ctx, l.failKey(username))
	pipe.Del(ctx, l.lockKey(username))
	_, err := pipe.Exec(ctx)
	return err
}
