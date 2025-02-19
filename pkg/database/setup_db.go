package database

import (
	"echobackend/config"
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// NewDatabase creates a new database connection using the provided configuration
func NewDatabase(config *config.Config) (*gorm.DB, error) {
	gormConfig := gorm.Config{
		Logger:      logger.Default.LogMode(logger.Silent),
		PrepareStmt: true,
	}

	db, err := gorm.Open(postgres.Open(config.GetDSN()), &gormConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	sqlDB.SetMaxOpenConns(20)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)

	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("database connection failed: %w", err)
	}

	return db, nil
}
