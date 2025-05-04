package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	App_Port string `mapstructure:"PORT"`
	// JWT configuration
	JWT_SECRET string `mapstructure:"JWT_SECRET"`
	// Database configuration
	Database_URL    string        `mapstructure:"DATABASE_URL"`
	MaxOpenConns    int           `mapstructure:"MAX_OPEN_CONNS"`
	MaxIdleConns    int           `mapstructure:"MAX_IDLE_CONNS"`
	ConnMaxLifetime time.Duration `mapstructure:"CONN_MAX_LIFETIME"`
	// Rate limiter configuration
	RATE_LIMITER_MAX int `mapstructure:"RATE_LIMITER_MAX"`
	RATE_LIMITER_TTL int `mapstructure:"RATE_LIMITER_TTL"`
	// Minio configuration
	MINIO_ENDPOINT   string `mapstructure:"MINIO_ENDPOINT"`
	MINIO_ACCESS_KEY string `mapstructure:"MINIO_ACCESS_KEY"`
	MINIO_SECRET_KEY string `mapstructure:"MINIO_SECRET_KEY"`
	MINIO_BUCKET     string `mapstructure:"MINIO_BUCKET"`
	MINIO_USE_SSL    bool   `mapstructure:"MINIO_USE_SSL"`
	// Debug mode
	DEBUG bool `mapstructure:"DEBUG"`
}

// Load reads configuration from environment variables with defaults
func Load() (*Config, error) {
	config := &Config{}

	setDefaults()

	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
	}

	viper.AutomaticEnv()

	if err := viper.Unmarshal(config); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}

	if err := config.validate(); err != nil {
		return nil, err
	}

	return config, nil
}

func setDefaults() {
	viper.SetDefault("PORT", "8080")
	viper.SetDefault("JWT_SECRET", "your-secret-key")
	viper.SetDefault("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable")
	viper.SetDefault("MAX_OPEN_CONNS", 30)
	viper.SetDefault("MAX_IDLE_CONNS", 2)
	viper.SetDefault("CONN_MAX_LIFETIME", 30*time.Minute)
	viper.SetDefault("RATE_LIMITER_MAX", 0)
	viper.SetDefault("RATE_LIMITER_TTL", 60)
	viper.SetDefault("MINIO_ENDPOINT", "localhost:9000")
	viper.SetDefault("MINIO_ACCESS_KEY", "minioadmin")
	viper.SetDefault("MINIO_SECRET_KEY", "minioadmin")
	viper.SetDefault("MINIO_BUCKET", "minio-bucket")
	viper.SetDefault("MINIO_USE_SSL", false)
	viper.SetDefault("DEBUG", false)
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
