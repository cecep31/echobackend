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

var DB *bun.DB

// Init initializes the database connection with proper configuration
func Init(conf *config.Config) error {
	config, err := pgx.ParseConfig(conf.Database_URL)
	if err != nil {
		return fmt.Errorf("failed to parse database config: %w", err)
	}
	config.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol

	sqldb := stdlib.OpenDB(*config)

	// Configure connection pool
	sqldb.SetMaxOpenConns(25)
	sqldb.SetMaxIdleConns(25)
	sqldb.SetConnMaxLifetime(5 * time.Minute)

	// Create Bun DB instance
	db := bun.NewDB(sqldb, pgdialect.New())

	// Add query hook for debugging
	// db.AddQueryHook(bundebug.NewQueryHook(
	// 	bundebug.WithVerbose(true),
	// 	bundebug.FromEnv("BUNDEBUG"),
	// ))

	// Verify connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	DB = db
	return nil
}
