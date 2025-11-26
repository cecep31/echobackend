package database

import (
	"context"
	"database/sql"
	"echobackend/config"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// NewDatabase creates a new database connection using the provided configuration
func NewDatabase(config *config.Config) *DatabaseWrapper {
	// Configure GORM logger
	var gormLogLevel logger.LogLevel
	if config.Debug {
		gormLogLevel = logger.Info
	} else {
		gormLogLevel = logger.Error
	}

	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(gormLogLevel),
	}

	// Parse database configuration
	pgxConfig, err := pgx.ParseConfig(config.DatabaseURL)
	if err != nil {
		panic(fmt.Errorf("failed to parse database config: %w", err))
	}
	pgxConfig.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol
	pgxConfig.ConnectTimeout = 10 * time.Second

	// Create database connection with retry logic
	var db *gorm.DB
	var sqldb *sql.DB

	// Retry connection with exponential backoff
	maxRetries := 3
	baseDelay := 1 * time.Second

	for attempt := 0; attempt < maxRetries; attempt++ {
		if attempt > 0 {
			delay := baseDelay * time.Duration(1<<uint(attempt-1))
			time.Sleep(delay)
			log.Printf("Retrying database connection (attempt %d/%d)", attempt+1, maxRetries)
		}

		sqldb = stdlib.OpenDB(*pgxConfig)

		// Configure connection pool with better defaults
		maxOpenConns := config.MaxOpenConns
		if maxOpenConns == 0 {
			maxOpenConns = 25
		}

		maxIdleConns := config.MaxIdleConns
		if maxIdleConns == 0 {
			maxIdleConns = 5
		}

		connMaxLifetime := config.ConnMaxLifetime
		if connMaxLifetime == 0 {
			connMaxLifetime = 1 * time.Hour
		}

		sqldb.SetMaxOpenConns(maxOpenConns)
		sqldb.SetMaxIdleConns(maxIdleConns)
		sqldb.SetConnMaxLifetime(connMaxLifetime)
		sqldb.SetConnMaxIdleTime(30 * time.Minute)

		// Create GORM DB instance
		db, err = gorm.Open(postgres.New(postgres.Config{
			Conn: sqldb,
		}), gormConfig)
		if err != nil {
			log.Printf("Failed to connect to database (attempt %d/%d): %v", attempt+1, maxRetries, err)
			if sqldb != nil {
				sqldb.Close()
			}
			continue
		}

		// Verify connection
		sqlDB, err := db.DB()
		if err != nil {
			log.Printf("Failed to get underlying sql.DB (attempt %d/%d): %v", attempt+1, maxRetries, err)
			continue
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := sqlDB.PingContext(ctx); err != nil {
			log.Printf("Failed to ping database (attempt %d/%d): %v", attempt+1, maxRetries, err)
			sqlDB.Close()
			continue
		}

		// Connection successful
		log.Printf("Successfully connected to database")
		break
	}

	if err != nil {
		panic(fmt.Errorf("failed to connect to database after %d attempts: %w", maxRetries, err))
	}

	return NewDatabaseWrapper(db)
}
