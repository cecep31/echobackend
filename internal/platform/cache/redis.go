package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"echobackend/config"
	"echobackend/pkg/applog"

	"github.com/redis/go-redis/v9"
)

var log = applog.Component("cache")

// RedisCache is a small JSON cache wrapper for Redis/Valkey.
type RedisCache struct {
	client    *redis.Client
	keyPrefix string
	ttl       time.Duration
}

var fixedWindowIncrementScript = redis.NewScript(`
local current = redis.call("INCR", KEYS[1])
if current == 1 then
	redis.call("PEXPIRE", KEYS[1], ARGV[1])
end
local ttl = redis.call("PTTL", KEYS[1])
return {current, ttl}
`)

// NewRedisCache creates a fail-open Redis client. If config is missing or invalid,
// it returns nil so the application can continue without caching.
func NewRedisCache(cfg *config.Config) *RedisCache {
	if cfg == nil || cfg.Cache.RedisURL == "" {
		return nil
	}

	opts, err := redis.ParseURL(cfg.Cache.RedisURL)
	if err != nil {
		log.Warn("cache: invalid CACHE_REDIS_URL, caching disabled", "error", err)
		return nil
	}

	// Cap connect timeout: use configured value but never more than 2s so startup
	// fails fast when Redis is unreachable (it's optional/fail-open anyway).
	connectTimeout := cfg.Cache.ConnectTimeout
	if connectTimeout <= 0 || connectTimeout > 2*time.Second {
		connectTimeout = 2 * time.Second
	}
	opts.DialTimeout = connectTimeout

	client := redis.NewClient(opts)

	pingCtx, cancel := context.WithTimeout(context.Background(), connectTimeout)
	defer cancel()
	if err := client.Ping(pingCtx).Err(); err != nil {
		log.Warn("cache: failed to connect to Redis, caching disabled", "error", err)
		_ = client.Close()
		return nil
	}

	if cfg.Cache.TTL <= 0 {
		log.Warn("cache: CACHE_TTL_SECONDS is 0 — SetJSON will be skipped, caching effectively disabled")
	}

	log.Info("cache: connected", "ttl", cfg.Cache.TTL, "key_prefix", strings.TrimSpace(cfg.Cache.KeyPrefix))

	return &RedisCache{
		client:    client,
		keyPrefix: strings.TrimSpace(cfg.Cache.KeyPrefix),
		ttl:       cfg.Cache.TTL,
	}
}

func (c *RedisCache) Close() error {
	if c == nil || c.client == nil {
		return nil
	}
	return c.client.Close()
}

func (c *RedisCache) BuildKey(parts ...string) string {
	if c == nil {
		return ""
	}

	if c.keyPrefix == "" {
		return strings.Join(parts, ":")
	}

	return c.keyPrefix + ":" + strings.Join(parts, ":")
}

func (c *RedisCache) GetJSON(ctx context.Context, key string, dest any) (bool, error) {
	value, found, err := c.getBytes(ctx, key, false)
	if err != nil || !found {
		return found, err
	}

	if err := json.Unmarshal(value, dest); err != nil {
		log.Warn("cache: GetJSON unmarshal error", "key", key, "error", err)
		return false, err
	}

	return true, nil
}

func (c *RedisCache) GetJSONAndDelete(ctx context.Context, key string, dest any) (bool, error) {
	value, found, err := c.getBytes(ctx, key, true)
	if err != nil || !found {
		return found, err
	}

	if err := json.Unmarshal(value, dest); err != nil {
		log.Warn("cache: GetJSONAndDelete unmarshal error", "key", key, "error", err)
		return false, err
	}

	return true, nil
}

func (c *RedisCache) SetJSON(ctx context.Context, key string, value any) error {
	if c == nil || c.client == nil || key == "" || c.ttl <= 0 {
		return nil
	}

	return c.SetJSONWithTTL(ctx, key, value, c.ttl)
}

func (c *RedisCache) SetJSONWithTTL(ctx context.Context, key string, value any, ttl time.Duration) error {
	if c == nil || c.client == nil || key == "" || ttl <= 0 {
		return nil
	}

	payload, err := json.Marshal(value)
	if err != nil {
		log.Warn("cache: SetJSON marshal error", "key", key, "error", err)
		return err
	}

	if err := c.client.Set(ctx, key, payload, ttl).Err(); err != nil {
		log.Warn("cache: SetJSON write error", "key", key, "error", err)
		return err
	}
	return nil
}

func (c *RedisCache) IncrementFixedWindow(ctx context.Context, key string, window time.Duration) (int, time.Duration, error) {
	if c == nil || c.client == nil || key == "" || window <= 0 {
		return 0, 0, nil
	}

	result, err := fixedWindowIncrementScript.Run(ctx, c.client, []string{key}, window.Milliseconds()).Result()
	if err != nil {
		log.Warn("cache: IncrementFixedWindow error", "key", key, "error", err)
		return 0, 0, err
	}

	values, ok := result.([]any)
	if !ok || len(values) != 2 {
		err := fmt.Errorf("unexpected redis script result: %T", result)
		log.Warn("cache: IncrementFixedWindow invalid result", "key", key, "error", err)
		return 0, 0, err
	}

	count, err := redisValueToInt64(values[0])
	if err != nil {
		log.Warn("cache: IncrementFixedWindow invalid count", "key", key, "error", err)
		return 0, 0, err
	}

	ttlMillis, err := redisValueToInt64(values[1])
	if err != nil {
		log.Warn("cache: IncrementFixedWindow invalid ttl", "key", key, "error", err)
		return 0, 0, err
	}
	if ttlMillis < 0 {
		ttlMillis = 0
	}

	return int(count), time.Duration(ttlMillis) * time.Millisecond, nil
}

func (c *RedisCache) getBytes(ctx context.Context, key string, deleteAfterRead bool) ([]byte, bool, error) {
	if c == nil || c.client == nil || key == "" {
		return nil, false, nil
	}

	cmd := c.client.Get(ctx, key)
	if deleteAfterRead {
		cmd = c.client.GetDel(ctx, key)
	}

	value, err := cmd.Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, false, nil
		}
		operation := "GetJSON"
		if deleteAfterRead {
			operation = "GetJSONAndDelete"
		}
		log.Warn("cache: "+operation+" error", "key", key, "error", err)
		return nil, false, err
	}

	return value, true, nil
}

func redisValueToInt64(value any) (int64, error) {
	switch v := value.(type) {
	case int64:
		return v, nil
	case int:
		return int64(v), nil
	case string:
		return strconv.ParseInt(v, 10, 64)
	case []byte:
		return strconv.ParseInt(string(v), 10, 64)
	default:
		return 0, fmt.Errorf("unexpected integer value type: %T", value)
	}
}
