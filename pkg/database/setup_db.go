package database

import (
	"echobackend/internal/config"
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func SetupDatabase(conf *config.Config) (*gorm.DB, error) {
	dsn := conf.GetDSN()

	var gormConfig gorm.Config
	if os.Getenv("ENABLE_GORM_LOGGER") != "" {
		gormConfig = gorm.Config{} // Default to verbose logging
	} else {
		gormConfig = gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent), // Disable logging
		}
	}

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true,
	}), &gormConfig)

	if err != nil {
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}

	return db, nil
}
