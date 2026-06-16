package database

import (
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

	// Open connection pool — connections are established lazily on first use.
	// We intentionally skip a blocking Ping here to keep startup fast;
	// the /health endpoint verifies liveness on demand.
	poolConfig := connectionPoolConfig{
		maxOpenConns:    defaultInt(config.Database.MaxOpenConns, 25),
		maxIdleConns:    defaultInt(config.Database.MaxIdleConns, 5),
		connMaxLifetime: defaultDuration(config.Database.ConnMaxLifetime, time.Hour),
		connMaxIdleTime: 30 * time.Minute,
	}

	sqldb := stdlib.OpenDB(*pgxConfig)
	configureConnectionPool(sqldb, poolConfig)

	db, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqldb,
	}), gormConfig)
	if err != nil {
		_ = sqldb.Close()
		panic(fmt.Errorf("failed to open database: %w", err))
	}

	slog.Info("database: pool ready", "max_open", poolConfig.maxOpenConns, "max_idle", poolConfig.maxIdleConns, "conn_lifetime", poolConfig.connMaxLifetime)
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
