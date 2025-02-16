package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	JWT_SECRET   string `mapstructure:"JWT_SECRET"`
	DATABASE_URL string `mapstructure:"DATABASE_URL"`
	PORT         string `mapstructure:"PORT"`

	RATE_LIMITER_MAX int `mapstructure:"RATE_LIMITER_MAX"`
	RATE_LIMITER_TTL int `mapstructure:"RATE_LIMITER_TTL"`

	MINIO_ENDPOINT   string `mapstructure:"MINIO_ENDPOINT"`
	MINIO_ACCESS_KEY string `mapstructure:"MINIO_ACCESS_KEY"`
	MINIO_SECRET_KEY string `mapstructure:"MINIO_SECRET_KEY"`
	MINIO_BUCKET     string `mapstructure:"MINIO_BUCKET"`
	MINIO_USE_SSL    bool   `mapstructure:"MINIO_USE_SSL"`
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
	viper.SetDefault("RATE_LIMITER_MAX", 1000)
	viper.SetDefault("RATE_LIMITER_TTL", 60)
	viper.SetDefault("MINIO_ENDPOINT", "localhost:9000")
	viper.SetDefault("MINIO_ACCESS_KEY", "minioadmin")
	viper.SetDefault("MINIO_SECRET_KEY", "minioadmin")
}

func (c *Config) validate() error {
	if c.JWT_SECRET == "" {
		return fmt.Errorf("JWT_SECRET is required")
	}
	if c.DATABASE_URL == "" {
		return fmt.Errorf("DATABASE_URL is required")
	}
	return nil
}

// GetDSN returns the database connection string
func (c *Config) GetDSN() string {
	return c.DATABASE_URL
}

func (c *Config) GetJWTSecret() string {
	return c.JWT_SECRET
}

func (c *Config) GetAppPort() string {
	return c.PORT
}

func (c *Config) GetRateLimiterMax() int {
	return c.RATE_LIMITER_MAX
}

func (c *Config) GetMinioEndpoint() string {
	return c.MINIO_ENDPOINT
}

func (c *Config) GetMinioAccessKey() string {
	return c.MINIO_ACCESS_KEY
}

func (c *Config) GetMinioSecretKey() string {
	return c.MINIO_SECRET_KEY
}
