package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/subosito/gotenv"
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

// Config represents the application configuration.
type Config struct {
	// AppPort is the port on which the HTTP server will listen
	AppPort string
	// JWT configuration
	// JWTSecret is the secret key used for JWT token signing and verification
	JWTSecret string
	// Database configuration
	// DatabaseURL is the PostgreSQL connection string
	DatabaseURL string
	// MaxOpenConns is the maximum number of open database connections in the pool
	MaxOpenConns int
	// MaxIdleConns is the maximum number of idle database connections in the pool
	MaxIdleConns int
	// ConnMaxLifetime is the maximum duration a database connection can be reused
	ConnMaxLifetime time.Duration
	// Rate limiter configuration
	// RateLimiterMax is the maximum number of requests allowed per window (0 = disabled)
	RateLimiterMax int
	// RateLimiterTTL is the time window in seconds for rate limiting
	RateLimiterTTL int
	// S3 configuration
	// S3 contains MinIO/S3 storage configuration
	S3 S3Config
	// Debug mode
	// Debug enables verbose logging and debug features
	Debug bool
}

// Load reads configuration from environment variables with defaults.
// It loads a .env file if present, then reads environment variables.
// Returns a validated Config struct or an error if required fields are missing.
func Load() (*Config, error) {
	// Load .env file
	gotenv.Load()

	config := &Config{
		AppPort:         getEnv("PORT", "8080"),
		JWTSecret:       getEnv("JWT_SECRET", "your-secret-key"),
		DatabaseURL:     getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"),
		MaxOpenConns:    getEnvAsInt("MAX_OPEN_CONNS", 30),
		MaxIdleConns:    getEnvAsInt("MAX_IDLE_CONNS", 2),
		ConnMaxLifetime: getEnvAsDuration("CONN_MAX_LIFETIME", 30*time.Minute),
		RateLimiterMax:  getEnvAsInt("RATE_LIMITER_MAX", 0),
		RateLimiterTTL:  getEnvAsInt("RATE_LIMITER_TTL", 60),
		S3: S3Config{
			Endpoint:  getEnv("MINIO_ENDPOINT", "localhost:9000"),
			AccessKey: getEnv("MINIO_ACCESS_KEY", "minioadmin"),
			SecretKey: getEnv("MINIO_SECRET_KEY", "minioadmin"),
			Bucket:    getEnv("MINIO_BUCKET", "minio-bucket"),
			UseSSL:    getEnvAsBool("MINIO_USE_SSL", false),
		},
		Debug: getEnvAsBool("DEBUG", false),
	}

	if err := config.validate(); err != nil {
		return nil, err
	}

	return config, nil
}

// Helper function to get environment variable with default value.
// Returns the environment variable value if it exists, otherwise returns the default value.
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// Helper function to get environment variable as integer with default value.
// Returns the parsed integer value if the environment variable exists and is valid,
// otherwise returns the default value.
func getEnvAsInt(key string, defaultValue int) int {
	if valueStr, exists := os.LookupEnv(key); exists {
		if value, err := strconv.Atoi(valueStr); err == nil {
			return value
		}
	}
	return defaultValue
}

// Helper function to get environment variable as boolean with default value.
// Returns the parsed boolean value if the environment variable exists and is valid,
// otherwise returns the default value.
func getEnvAsBool(key string, defaultValue bool) bool {
	if valueStr, exists := os.LookupEnv(key); exists {
		if value, err := strconv.ParseBool(valueStr); err == nil {
			return value
		}
	}
	return defaultValue
}

// Helper function to get environment variable as duration with default value.
// Returns the parsed duration value if the environment variable exists and is valid,
// otherwise returns the default value.
func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	if valueStr, exists := os.LookupEnv(key); exists {
		if value, err := time.ParseDuration(valueStr); err == nil {
			return value
		}
	}
	return defaultValue
}

// validate ensures that all required configuration fields are present and valid.
// Returns an error if any required field is missing or invalid.
func (c *Config) validate() error {
	if c.JWTSecret == "" {
		return fmt.Errorf("JWT_SECRET is required")
	}
	if c.DatabaseURL == "" {
		return fmt.Errorf("DATABASE_URL is required")
	}
	return nil
}
