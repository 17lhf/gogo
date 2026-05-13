package cache

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

// HeartbeatCache manages terminal heartbeat TTL in Redis.
type HeartbeatCache struct {
	rdb *redis.Client
	ttl time.Duration
}

// NewHeartbeatCache creates a new HeartbeatCache.
func NewHeartbeatCache(rdb *redis.Client, ttl time.Duration) *HeartbeatCache {
	return &HeartbeatCache{rdb: rdb, ttl: ttl}
}

const heartbeatKeyPrefix = "heartbeat:"

func (h *HeartbeatCache) key(sn string) string {
	return heartbeatKeyPrefix + sn
}

// Set sets the heartbeat key with TTL. Called on each heartbeat.
func (h *HeartbeatCache) Set(ctx context.Context, sn string) error {
	return h.rdb.Set(ctx, h.key(sn), time.Now().Unix(), h.ttl).Err()
}

// Exists checks if a heartbeat key exists.
func (h *HeartbeatCache) Exists(ctx context.Context, sn string) (bool, error) {
	n, err := h.rdb.Exists(ctx, h.key(sn)).Result()
	if err != nil {
		return false, err
	}
	return n > 0, nil
}

// Delete removes a heartbeat key (on terminal disable/delete).
func (h *HeartbeatCache) Delete(ctx context.Context, sn string) error {
	return h.rdb.Del(ctx, h.key(sn)).Err()
}

// GetTokenKey returns the Redis key for device token.
func GetTokenKey(token string) string {
	return fmt.Sprintf("device_token:%s", token)
}

// SetDeviceToken stores a device token → SN mapping in Redis.
func (h *HeartbeatCache) SetDeviceToken(ctx context.Context, token, sn string) error {
	return h.rdb.Set(ctx, GetTokenKey(token), sn, h.ttl).Err()
}

// GetSNByDeviceToken retrieves the SN for a device token.
func (h *HeartbeatCache) GetSNByDeviceToken(ctx context.Context, token string) (string, error) {
	sn, err := h.rdb.Get(ctx, GetTokenKey(token)).Result()
	if err == redis.Nil {
		return "", nil
	}
	return sn, err
}

// DeleteDeviceToken removes a device token from Redis.
func (h *HeartbeatCache) DeleteDeviceToken(ctx context.Context, token string) error {
	return h.rdb.Del(ctx, GetTokenKey(token)).Err()
}

// ListenForExpiry subscribes to Redis keyspace expired events and invokes the callback.
// Requires Redis configured with: notify-keyspace-events Ex
func ListenForExpiry(ctx context.Context, rdb *redis.Client, onExpiry func(sn string)) {
	// Ensure keyspace notifications are enabled
	rdb.ConfigSet(ctx, "notify-keyspace-events", "Ex")

	pubsub := rdb.PSubscribe(ctx, "__keyevent@0__:expired")
	ch := pubsub.Channel()

	slog.Info("listening for heartbeat expiry events")

	for {
		select {
		case <-ctx.Done():
			pubsub.Close()
			return
		case msg, ok := <-ch:
			if !ok {
				return
			}
			if strings.HasPrefix(msg.Payload, heartbeatKeyPrefix) {
				sn := strings.TrimPrefix(msg.Payload, heartbeatKeyPrefix)
				onExpiry(sn)
			}
		}
	}
}
