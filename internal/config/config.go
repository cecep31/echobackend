package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	JWTSecret string `mapstructure:"JWT_SECRET"`
	DSN       string `mapstructure:"DATABASE_URL"`
}

// Load reads configuration from environment variables with defaults
func Load() (*Config, error) {
	config := &Config{}

	viper.SetConfigName(".env") // name of config file (without extension)
	viper.SetConfigType("env")  // type of config file
	viper.AddConfigPath(".")    // optionally look for config in the working directory

	err := viper.ReadInConfig()
	// Read config file
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
		// Config file not found; ignore error if desired
		fmt.Println("No configuration file found. Using defaults and environment variables.")
	}

	viper.AutomaticEnv()

	// Set defaults
	setDefaults()

	// Unmarshal config
	if err := viper.Unmarshal(config); err != nil {
		return nil, fmt.Errorf("unable to decode config into struct: %w", err)
	}

	// Validate config
	if err := config.validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return config, nil
}

func setDefaults() {

	// Database defaults
	viper.SetDefault("database.DSN", "")

	// Auth defaults
	viper.SetDefault("auth.JWTSecret", "")
}

func (c *Config) validate() error {
	if c.DSN == "" {
		return fmt.Errorf("database url is required")
	}
	if c.JWTSecret == "" {
		return fmt.Errorf("JWT secret is required")
	}
	return nil
}

// GetDSN returns the database connection string
func (c *Config) GetDSN() string {
	return c.DSN
}

func (c *Config) GetJWTSecret() string {
	return c.JWTSecret
}
