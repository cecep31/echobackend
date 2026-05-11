package database

import (
	"context"
	"database/sql"
	"echobackend/config"
	"fmt"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// NewDatabase creates a new database connection using the provided configuration
func NewDatabase(config *config.Config) *gorm.DB {
	// Configure GORM logger
	var gormLogLevel logger.LogLevel
	if config.AppDebug {
		gormLogLevel = logger.Info
	} else {
		gormLogLevel = logger.Error
	}

	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(gormLogLevel),
	}

	// Parse database configuration
	pgxConfig, err := pgx.ParseConfig(config.PostgresDSN)
	if err != nil {
		panic(fmt.Errorf("failed to parse database config: %w", err))
	}
	// Use the default extended query protocol for better performance (named statements,
	// binary encoding). If you run behind PgBouncer in transaction-pooling mode, set
	// PGX_QUERY_EXEC_MODE=simple or switch back to QueryExecModeSimpleProtocol here.
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
			slog.Info("retrying database connection", "attempt", attempt+1, "max", maxRetries)
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
			slog.Error("failed to connect to database", "attempt", attempt+1, "max", maxRetries, "error", err)
			if sqldb != nil {
				sqldb.Close()
			}
			continue
		}

		// Verify connection
		sqlDB, err := db.DB()
		if err != nil {
			slog.Error("failed to get underlying sql.DB", "attempt", attempt+1, "max", maxRetries, "error", err)
			continue
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		err = sqlDB.PingContext(ctx)
		cancel()
		if err != nil {
			slog.Error("failed to ping database", "attempt", attempt+1, "max", maxRetries, "error", err)
			sqlDB.Close()
			continue
		}

		// Connection successful
		slog.Info("successfully connected to database")
		break
	}

	if err != nil {
		panic(fmt.Errorf("failed to connect to database after %d attempts: %w", maxRetries, err))
	}

	return db
}
