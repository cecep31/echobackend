package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/subosito/gotenv"
)

type MinioConfig struct {
	Endpoint  string
	AccessKey string
	SecretKey string
	Bucket    string
	UseSSL    bool
}

type Config struct {
	App_Port string
	// JWT configuration
	JWT_SECRET string
	// Database configuration
	Database_URL    string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	// Rate limiter configuration
	RATE_LIMITER_MAX int
	RATE_LIMITER_TTL int
	// Minio configuration
	Minio MinioConfig
	// Debug mode
	DEBUG bool
}

// Load reads configuration from environment variables with defaults
func Load() (*Config, error) {
	// Load .env file
	gotenv.Load()

	config := &Config{
		App_Port:         getEnv("PORT", "8080"),
		JWT_SECRET:       getEnv("JWT_SECRET", "your-secret-key"),
		Database_URL:     getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"),
		MaxOpenConns:     getEnvAsInt("MAX_OPEN_CONNS", 30),
		MaxIdleConns:     getEnvAsInt("MAX_IDLE_CONNS", 2),
		ConnMaxLifetime:  getEnvAsDuration("CONN_MAX_LIFETIME", 30*time.Minute),
		RATE_LIMITER_MAX: getEnvAsInt("RATE_LIMITER_MAX", 0),
		RATE_LIMITER_TTL: getEnvAsInt("RATE_LIMITER_TTL", 60),
		Minio: MinioConfig{
			Endpoint:  getEnv("MINIO_ENDPOINT", "localhost:9000"),
			AccessKey: getEnv("MINIO_ACCESS_KEY", "minioadmin"),
			SecretKey: getEnv("MINIO_SECRET_KEY", "minioadmin"),
			Bucket:    getEnv("MINIO_BUCKET", "minio-bucket"),
			UseSSL:    getEnvAsBool("MINIO_USE_SSL", false),
		},
		DEBUG: getEnvAsBool("DEBUG", false),
	}

	if err := config.validate(); err != nil {
		return nil, err
	}

	return config, nil
}

// Helper function to get environment variable with default value
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// Helper function to get environment variable as integer with default value
func getEnvAsInt(key string, defaultValue int) int {
	if valueStr, exists := os.LookupEnv(key); exists {
		if value, err := strconv.Atoi(valueStr); err == nil {
			return value
		}
	}
	return defaultValue
}

// Helper function to get environment variable as boolean with default value
func getEnvAsBool(key string, defaultValue bool) bool {
	if valueStr, exists := os.LookupEnv(key); exists {
		if value, err := strconv.ParseBool(valueStr); err == nil {
			return value
		}
	}
	return defaultValue
}

// Helper function to get environment variable as duration with default value
func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	if valueStr, exists := os.LookupEnv(key); exists {
		if value, err := time.ParseDuration(valueStr); err == nil {
			return value
		}
	}
	return defaultValue
}

func (c *Config) validate() error {
	if c.JWT_SECRET == "" {
		return fmt.Errorf("JWT_SECRET is required")
	}
	if c.Database_URL == "" {
		return fmt.Errorf("DATABASE_URL is required")
	}
	return nil
}
