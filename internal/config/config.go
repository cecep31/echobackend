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

	err := viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
		// Ignore error if config file not found
	}

	err = viper.Unmarshal(config)
	if err != nil {
		return nil, fmt.Errorf("unable to decode config into struct: %w", err)
	}

	err = config.validate()
	if err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return config, nil
}

func setDefaults() {
	viper.SetDefault("PORT", "1323")
	// Database defaults
	viper.SetDefault("DATABASE_URL", "")
	// Auth defaults
	viper.SetDefault("JWT_SECRET", "")
}

func (c *Config) validate() error {
	if c.DATABASE_URL == "" {
		return fmt.Errorf("database url is required")
	}
	if c.JWT_SECRET == "" {
		return fmt.Errorf("JWT secret is required")
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
