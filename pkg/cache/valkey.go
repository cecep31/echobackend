package cache

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"echobackend/config"

	"github.com/redis/go-redis/v9"
)

// ValkeyCache is a small JSON cache wrapper for Valkey/Redis.
type ValkeyCache struct {
	client    *redis.Client
	keyPrefix string
	ttl       time.Duration
}

// NewValkeyCache creates a fail-open Valkey client. If config is missing or invalid,
// it returns nil so the application can continue without caching.
func NewValkeyCache(cfg *config.Config) *ValkeyCache {
	if cfg == nil || cfg.Cache.ValkeyURL == "" {
		return nil
	}

	opts, err := redis.ParseURL(cfg.Cache.ValkeyURL)
	if err != nil {
		return nil
	}

	if cfg.Cache.ConnectTimeout > 0 {
		opts.DialTimeout = cfg.Cache.ConnectTimeout
	}

	client := redis.NewClient(opts)

	pingCtx := context.Background()
	if cfg.Cache.ConnectTimeout > 0 {
		var cancel context.CancelFunc
		pingCtx, cancel = context.WithTimeout(pingCtx, cfg.Cache.ConnectTimeout)
		defer cancel()
	}

	if err := client.Ping(pingCtx).Err(); err != nil {
		_ = client.Close()
		return nil
	}

	return &ValkeyCache{
		client:    client,
		keyPrefix: strings.TrimSpace(cfg.Cache.KeyPrefix),
		ttl:       cfg.Cache.TTL,
	}
}

func (c *ValkeyCache) Close() error {
	if c == nil || c.client == nil {
		return nil
	}
	return c.client.Close()
}

func (c *ValkeyCache) BuildKey(parts ...string) string {
	if c == nil {
		return ""
	}

	if c.keyPrefix == "" {
		return strings.Join(parts, ":")
	}

	return c.keyPrefix + ":" + strings.Join(parts, ":")
}

func (c *ValkeyCache) GetJSON(ctx context.Context, key string, dest any) (bool, error) {
	if c == nil || c.client == nil || key == "" {
		return false, nil
	}

	value, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return false, nil
		}
		return false, err
	}

	if err := json.Unmarshal(value, dest); err != nil {
		return false, err
	}

	return true, nil
}

func (c *ValkeyCache) SetJSON(ctx context.Context, key string, value any) error {
	if c == nil || c.client == nil || key == "" || c.ttl <= 0 {
		return nil
	}

	payload, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return c.client.Set(ctx, key, payload, c.ttl).Err()
}
