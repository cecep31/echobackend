package database

import (
	"context"
	"echobackend/config"

	// "echobackend/internal/model" // No longer needed as AutoMigrate is removed
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// NewDatabase creates a new database connection using the provided configuration
func NewDatabase(config *config.Config) *DatabaseWrapper {
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

	// Create GORM DB instance
	db, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqldb,
	}), &gorm.Config{})
	if err != nil {
		panic(fmt.Errorf("failed to connect to database: %w", err))
	}

	// Verify connection (GORM does this implicitly on Open, but we can do an explicit check if needed)
	// For example, get the underlying sql.DB and ping it
	sqlDB, err := db.DB()
	if err != nil {
		panic(fmt.Errorf("failed to get underlying sql.DB: %w", err))
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := sqlDB.PingContext(ctx); err != nil {
		panic(fmt.Errorf("failed to ping database: %w", err))
	}

	// AutoMigrate models (GORM way of registering/updating schema)
	// User has requested to remove AutoMigrate. Schema migrations should be handled externally.
	// All known models have been updated for GORM.
	// err = db.AutoMigrate(
	// 	&model.User{},
	// 	&model.Post{},
	// 	&model.Tag{},
	// 	&model.Page{},
	// 	&model.Block{},
	// 	&model.Workspace{},
	// 	&model.WorkspaceMember{},
	// 	&model.Comment{},
	// )
	// if err != nil {
	// 	panic(fmt.Errorf("failed to auto migrate models: %w", err))
	// }

	// Note: GORM's AutoMigrate creates tables, missing foreign keys, constraints, columns and indexes.
	// It will NOT delete unused columns to protect your data.
	// For more complex migrations (e.g., renaming columns, changing types with data loss), a dedicated migration tool is recommended.

	return NewDatabaseWrapper(db)
}
