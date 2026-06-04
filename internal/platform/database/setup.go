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
	var gormLogLevel logger.LogLevel
	if config.App.Debug {
		gormLogLevel = logger.Info
	} else {
		gormLogLevel = logger.Error
	}

	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(gormLogLevel),
	}

	pgxConfig, err := pgx.ParseConfig(config.Database.DSN)
	if err != nil {
		panic(fmt.Errorf("failed to parse database config: %w", err))
	}
	// Use the default extended query protocol for better performance (named statements,
	// binary encoding). If you run behind PgBouncer in transaction-pooling mode, set
	// PGX_QUERY_EXEC_MODE=simple or switch back to QueryExecModeSimpleProtocol here.
	pgxConfig.ConnectTimeout = 10 * time.Second

	// Create database connection with retry logic
	var db *gorm.DB

	// Retry connection with exponential backoff
	maxRetries := 3
	baseDelay := 1 * time.Second
	poolConfig := connectionPoolConfig{
		maxOpenConns:    defaultInt(config.Database.MaxOpenConns, 25),
		maxIdleConns:    defaultInt(config.Database.MaxIdleConns, 5),
		connMaxLifetime: defaultDuration(config.Database.ConnMaxLifetime, time.Hour),
		connMaxIdleTime: 30 * time.Minute,
	}

	for attempt := 0; attempt < maxRetries; attempt++ {
		if attempt > 0 {
			delay := baseDelay * time.Duration(1<<uint(attempt-1))
			time.Sleep(delay)
			slog.Info("retrying database connection", "attempt", attempt+1, "max", maxRetries)
		}

		sqldb := stdlib.OpenDB(*pgxConfig)
		configureConnectionPool(sqldb, poolConfig)

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
			_ = sqldb.Close()
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
		slog.Info("database: connected", "max_open", poolConfig.maxOpenConns, "max_idle", poolConfig.maxIdleConns, "conn_lifetime", poolConfig.connMaxLifetime)
		break
	}

	if err != nil {
		panic(fmt.Errorf("failed to connect to database after %d attempts: %w", maxRetries, err))
	}

	return db
}

type connectionPoolConfig struct {
	maxOpenConns    int
	maxIdleConns    int
	connMaxLifetime time.Duration
	connMaxIdleTime time.Duration
}

func configureConnectionPool(db *sql.DB, cfg connectionPoolConfig) {
	db.SetMaxOpenConns(cfg.maxOpenConns)
	db.SetMaxIdleConns(cfg.maxIdleConns)
	db.SetConnMaxLifetime(cfg.connMaxLifetime)
	db.SetConnMaxIdleTime(cfg.connMaxIdleTime)
}

func defaultInt(value, fallback int) int {
	if value == 0 {
		return fallback
	}
	return value
}

func defaultDuration(value, fallback time.Duration) time.Duration {
	if value == 0 {
		return fallback
	}
	return value
}
