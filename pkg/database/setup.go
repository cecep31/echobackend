package database

import (
	"context"
	"echobackend/config"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
)

// NewDatabase creates a new database connection using the provided configuration
func NewDatabase(config *config.Config) *bun.DB {
	pgxConfig, err := pgx.ParseConfig(config.Database_URL)
	if err != nil {
		panic(fmt.Errorf("failed to parse database config: %w", err))
	}
	pgxConfig.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol

	sqldb := stdlib.OpenDB(*pgxConfig)

	// Configure connection pool
	sqldb.SetMaxOpenConns(config.MaxOpenConns)
	sqldb.SetMaxIdleConns(config.MaxIdleConns)
	sqldb.SetConnMaxLifetime(config.ConnMaxLifetime)
	sqldb.SetConnMaxIdleTime(10 * time.Minute) // Add max idle time to recycle connections

	// Create Bun DB instance
	db := bun.NewDB(sqldb, pgdialect.New())

	// Verify connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		panic(fmt.Errorf("failed to ping database: %w", err))
	}

	// Register models with Bun ORM
	RegisterModels(db)

	return db
}
