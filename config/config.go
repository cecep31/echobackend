package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type S3Config struct {
	// Endpoint is the S3 service endpoint (e.g., "localhost:9000")
	Endpoint string
	// AccessKey is the S3 access key for authentication
	AccessKey string
	// SecretKey is the S3 secret key for authentication
	SecretKey string
	// Bucket is the default S3 bucket name
	Bucket string
	// UseSSL determines whether to use SSL/TLS for S3 connections
	UseSSL bool
}

type CacheConfig struct {
	// ValkeyURL is the connection URL for Valkey/Redis (for example redis://localhost:6379/0).
	ValkeyURL string
	// KeyPrefix is prepended to cache keys to avoid collisions across apps/environments.
	KeyPrefix string
	// TTL is the default cache lifetime.
	TTL time.Duration
	// ConnectTimeout is used when establishing a new cache connection.
	ConnectTimeout time.Duration
}

// Config represents the application configuration.
type Config struct {
	// HTTPPort is the TCP port the HTTP server listens on (e.g. "8080").
	HTTPPort string
	// JWTSecret is the secret key used for JWT token signing and verification.
	JWTSecret string
	// PostgresDSN is the PostgreSQL connection string (lib/pgx / GORM DSN).
	PostgresDSN string
	// MaxOpenConns is the maximum number of open database connections in the pool.
	MaxOpenConns int
	// MaxIdleConns is the maximum number of idle database connections in the pool.
	MaxIdleConns int
	// ConnMaxLifetime is the maximum duration a database connection can be reused.
	ConnMaxLifetime time.Duration
	// HTTPRateLimitRPS is the global Echo rate limiter sustained rate in requests per second (0 = disabled).
	HTTPRateLimitRPS int
	// HTTPRateLimitWindowSec is, when >0, the Echo memory store visitor ExpiresIn in seconds. When 0, Echo defaults (3m) apply.
	HTTPRateLimitWindowSec int
	// HTTPTrustProxy when true sets Echo IP extraction from X-Forwarded-For with trust rules (typical behind nginx/ALB).
	// When false, only the direct TCP peer address is used (safer when the app faces the internet without a trusted proxy).
	HTTPTrustProxy bool
	// HTTPAllowOrigins is the value(s) for the Access-Control-Allow-Origin CORS header.
	// Comma-separated list of origins, or "*" to allow any origin.
	HTTPAllowOrigins string
	// S3 contains S3-compatible object storage (MinIO, AWS S3, etc.).
	S3 S3Config
	// Cache contains Valkey/Redis configuration.
	Cache CacheConfig
	// AppDebug enables verbose logging, GORM info logs, and debug routes.
	AppDebug bool
}

// Load reads configuration from environment variables with defaults.
// It loads a .env file if present, then reads environment variables.
// Returns a validated Config struct or an error if required fields are missing.
func Load() (*Config, error) {
	// Load .env file
	_ = godotenv.Load()

	// New env names are scoped where noted (DB_POOL_*, S3_*, etc.).
	config := &Config{
		HTTPPort:  envString([]string{"PORT"}, "8080"),
		JWTSecret: envString([]string{"JWT_SECRET"}, "your-secret-key"),
		// Default requires TLS; for local Postgres without SSL, set DATABASE_URL (see .env.example).
		PostgresDSN:            envString([]string{"DATABASE_URL"}, "postgres://postgres:postgres@127.0.0.1:5432/postgres?sslmode=require"),
		MaxOpenConns:           envInt([]string{"DB_POOL_MAX_OPEN", "MAX_OPEN_CONNS"}, 25),
		MaxIdleConns:           envInt([]string{"DB_POOL_MAX_IDLE", "MAX_IDLE_CONNS"}, 10),
		ConnMaxLifetime:        envDuration([]string{"DB_POOL_CONN_LIFETIME", "CONN_MAX_LIFETIME"}, 15*time.Minute),
		HTTPRateLimitRPS:       envInt([]string{"HTTP_RATE_LIMIT_RPS", "RATE_LIMITER_MAX"}, 0),
		HTTPRateLimitWindowSec: envInt([]string{"HTTP_RATE_LIMIT_WINDOW_SEC", "RATE_LIMITER_TTL"}, 0),
		HTTPTrustProxy:         envBool([]string{"HTTP_TRUST_PROXY", "TRUST_PROXY"}, false),
		HTTPAllowOrigins:       envString([]string{"HTTP_ALLOW_ORIGINS"}, "*"),
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
		AppDebug: envBool([]string{"APP_DEBUG", "DEBUG"}, false),
	}

	if err := config.validate(); err != nil {
		return nil, err
	}

	return config, nil
}

// validate ensures that all required configuration fields are present and valid.
// Returns an error if any required field is missing or invalid.
func (c *Config) validate() error {
	if c.JWTSecret == "" {
		return fmt.Errorf("JWT_SECRET is required")
	}
	if c.PostgresDSN == "" {
		return fmt.Errorf("DATABASE_URL is required")
	}
	if c.Cache.TTL < 0 {
		return fmt.Errorf("CACHE_TTL_SECONDS must be >= 0")
	}
	if c.Cache.ConnectTimeout < 0 {
		return fmt.Errorf("VALKEY_CONNECT_TIMEOUT_MS must be >= 0")
	}
	return nil
}

// envString returns the first set environment variable from keys, or defaultValue.
func envString(keys []string, defaultValue string) string {
	for _, k := range keys {
		if v, ok := os.LookupEnv(k); ok {
			return v
		}
	}
	return defaultValue
}

// envInt returns the first successfully parsed int from set env keys, or defaultValue.
func envInt(keys []string, defaultValue int) int {
	for _, k := range keys {
		if s, ok := os.LookupEnv(k); ok {
			if n, err := strconv.Atoi(s); err == nil {
				return n
			}
			return defaultValue
		}
	}
	return defaultValue
}

// envBool returns the first successfully parsed bool from set env keys, or defaultValue.
func envBool(keys []string, defaultValue bool) bool {
	for _, k := range keys {
		if s, ok := os.LookupEnv(k); ok {
			if b, err := strconv.ParseBool(s); err == nil {
				return b
			}
			return defaultValue
		}
	}
	return defaultValue
}

// envDuration returns the first successfully parsed duration from set env keys, or defaultValue.
func envDuration(keys []string, defaultValue time.Duration) time.Duration {
	for _, k := range keys {
		if s, ok := os.LookupEnv(k); ok {
			if d, err := time.ParseDuration(s); err == nil {
				return d
			}
			return defaultValue
		}
	}
	return defaultValue
}
