package database

import (
	"echobackend/config"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// NewDatabase creates a new database connection using the provided configuration
func NewDatabase(config *config.Config) *gorm.DB {
	gormConfig := gorm.Config{
		Logger:      logger.Default.LogMode(logger.Silent),
		PrepareStmt: true,
		SkipDefaultTransaction: true, // Improves performance for non-transactional operations
	}

	db, err := gorm.Open(postgres.Open(config.Database_URL), &gormConfig)
	if err != nil {
		panic(err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}

	sqlDB.SetMaxOpenConns(config.MaxOpenConns)
	sqlDB.SetMaxIdleConns(config.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(config.ConnMaxLifetime)
	sqlDB.SetConnMaxIdleTime(10 * time.Minute) // Add max idle time to recycle connections

	if err := sqlDB.Ping(); err != nil {
		panic(err)
	}

	return db
}
