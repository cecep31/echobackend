package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	JWT_SECRET   string `mapstructure:"JWT_SECRET"`
	DATABASE_URL string `mapstructure:"DATABASE_URL"`
	PORT         string `mapstructure:"PORT"`
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
