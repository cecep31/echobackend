// Package config loads and validates the application configuration.
//
// All values come from environment variables (a .env file is loaded first if
// present) and are grouped into section configs on the root Config struct:
//
//	cfg.App       // app-level toggles (debug)
//	cfg.HTTP      // listener, rate limit, CORS, proxy trust
//	cfg.Auth      // JWT secret
//	cfg.Database  // PostgreSQL DSN and pool tuning
//	cfg.S3        // S3-compatible object storage
//	cfg.Cache     // Valkey/Redis cache
//
// Some env keys have fallback aliases (legacy names). The first-set key wins;
// see Load() for the full list.
package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

// Config is the root application configuration.
type Config struct {
	App      AppConfig
	HTTP     HTTPConfig
	Auth     AuthConfig
	Database DatabaseConfig
	S3       S3Config
	Cache    CacheConfig
}

// AppConfig contains application-level toggles.
type AppConfig struct {
	// Debug enables verbose logging, GORM info logs, and debug routes (/api/debug/pprof/*).
	Debug bool
}

// HTTPConfig controls the HTTP server, rate limiting, CORS and proxy trust.
type HTTPConfig struct {
	// Port is the TCP port the HTTP server listens on (e.g. "8080").
	Port string
	// RateLimitRPS is the sustained rate in requests per second for the global
	// Echo rate limiter (0 = disabled). Burst is set to 2x this value.
	RateLimitRPS int
	// RateLimitWindow, when > 0, sets the Echo memory-store visitor ExpiresIn.
	// When 0, Echo defaults (3 minutes) apply.
	RateLimitWindow time.Duration
	// TrustProxy, when true, extracts client IP from X-Forwarded-For.
	// Use only behind a trusted reverse proxy.
	TrustProxy bool
	// AllowOrigins is the parsed CORS allow-list (always non-empty; defaults to ["*"]).
	AllowOrigins []string
}

// AuthConfig contains authentication secrets.
type AuthConfig struct {
	// JWTSecret is the secret key used for JWT token signing and verification.
	JWTSecret string
}

// DatabaseConfig contains the PostgreSQL DSN and connection pool tuning.
type DatabaseConfig struct {
	// DSN is the PostgreSQL connection string (pgx / GORM DSN).
	DSN string
	// MaxOpenConns is the maximum number of open database connections in the pool.
	MaxOpenConns int
	// MaxIdleConns is the maximum number of idle database connections in the pool.
	MaxIdleConns int
	// ConnMaxLifetime is the maximum duration a database connection can be reused.
	ConnMaxLifetime time.Duration
}

// S3Config contains S3-compatible object storage settings (MinIO, AWS S3, etc.).
type S3Config struct {
	// Endpoint is the S3 service endpoint (e.g. "localhost:9000").
	Endpoint string
	// AccessKey is the S3 access key for authentication.
	AccessKey string
	// SecretKey is the S3 secret key for authentication.
	SecretKey string
	// Bucket is the default S3 bucket name.
	Bucket string
	// UseSSL determines whether to use SSL/TLS for S3 connections.
	UseSSL bool
}

// CacheConfig contains Valkey/Redis cache settings.
type CacheConfig struct {
	// ValkeyURL is the connection URL for Valkey/Redis (e.g. redis://localhost:6379/0).
	// Empty disables caching (the app runs fail-open).
	ValkeyURL string
	// KeyPrefix is prepended to cache keys to avoid collisions across apps/environments.
	KeyPrefix string
	// TTL is the default cache lifetime. 0 effectively disables writes.
	TTL time.Duration
	// ConnectTimeout is used when establishing a new cache connection.
	ConnectTimeout time.Duration
}

// Load reads configuration from environment variables with defaults.
//
// It loads a .env file if present, then reads environment variables.
// Returns a validated *Config or an error if required fields are missing.
//
// Some keys accept legacy aliases; the first-set key wins.
func Load() (*Config, error) {
	// Best-effort .env loading; missing file is not an error.
	_ = godotenv.Load()

	cfg := &Config{
		App: AppConfig{
			Debug: envBool([]string{"APP_DEBUG", "DEBUG"}, false),
		},
		HTTP: HTTPConfig{
			Port:            envString([]string{"PORT"}, "8080"),
			RateLimitRPS:    envInt([]string{"HTTP_RATE_LIMIT_RPS", "RATE_LIMITER_MAX"}, 0),
			RateLimitWindow: time.Duration(envInt([]string{"HTTP_RATE_LIMIT_WINDOW_SEC", "RATE_LIMITER_TTL"}, 0)) * time.Second,
			TrustProxy:      envBool([]string{"HTTP_TRUST_PROXY", "TRUST_PROXY"}, false),
			AllowOrigins:    parseOrigins(envString([]string{"HTTP_ALLOW_ORIGINS"}, "*")),
		},
		Auth: AuthConfig{
			JWTSecret: envString([]string{"JWT_SECRET"}, "your-secret-key"),
		},
		Database: DatabaseConfig{
			// Default requires TLS; for local Postgres without SSL, set DATABASE_URL (see .env.example).
			DSN:             envString([]string{"DATABASE_URL"}, "postgres://postgres:postgres@127.0.0.1:5432/postgres?sslmode=require"),
			MaxOpenConns:    envInt([]string{"DB_POOL_MAX_OPEN", "MAX_OPEN_CONNS"}, 25),
			MaxIdleConns:    envInt([]string{"DB_POOL_MAX_IDLE", "MAX_IDLE_CONNS"}, 10),
			ConnMaxLifetime: envDuration([]string{"DB_POOL_CONN_LIFETIME", "CONN_MAX_LIFETIME"}, 15*time.Minute),
		},
		S3: S3Config{
			Endpoint:  envString([]string{"S3_ENDPOINT", "MINIO_ENDPOINT"}, "localhost:9000"),
			AccessKey: envString([]string{"S3_ACCESS_KEY", "MINIO_ACCESS_KEY"}, "minioadmin"),
			SecretKey: envString([]string{"S3_SECRET_KEY", "MINIO_SECRET_KEY"}, "minioadmin"),
			Bucket:    envString([]string{"S3_BUCKET", "MINIO_BUCKET"}, "minio-bucket"),
			UseSSL:    envBool([]string{"S3_USE_SSL", "MINIO_USE_SSL"}, false),
		},
		Cache: CacheConfig{
			ValkeyURL:      envString([]string{"VALKEY_URL"}, ""),
			KeyPrefix:      envString([]string{"CACHE_KEY_PREFIX"}, "pilput"),
			TTL:            time.Duration(envInt([]string{"CACHE_TTL_SECONDS"}, 60)) * time.Second,
			ConnectTimeout: time.Duration(envInt([]string{"VALKEY_CONNECT_TIMEOUT_MS"}, 5000)) * time.Millisecond,
		},
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// validate checks cross-section invariants and required fields.
func (c *Config) validate() error {
	if c.Auth.JWTSecret == "" {
		return fmt.Errorf("JWT_SECRET is required")
	}
	if c.Database.DSN == "" {
		return fmt.Errorf("DATABASE_URL is required")
	}
	if c.Cache.TTL < 0 {
		return fmt.Errorf("CACHE_TTL_SECONDS must be >= 0")
	}
	if c.Cache.ConnectTimeout < 0 {
		return fmt.Errorf("VALKEY_CONNECT_TIMEOUT_MS must be >= 0")
	}
	if c.HTTP.RateLimitRPS < 0 {
		return fmt.Errorf("HTTP_RATE_LIMIT_RPS must be >= 0")
	}
	return nil
}

// parseOrigins splits a comma-separated CORS origin list, trims whitespace,
// and drops empty entries. Returns ["*"] when the raw value is empty or "*".
func parseOrigins(raw string) []string {
	raw = strings.TrimSpace(raw)
	if raw == "" || raw == "*" {
		return []string{"*"}
	}
	parts := strings.Split(raw, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if v := strings.TrimSpace(p); v != "" {
			out = append(out, v)
		}
	}
	if len(out) == 0 {
		return []string{"*"}
	}
	return out
}
