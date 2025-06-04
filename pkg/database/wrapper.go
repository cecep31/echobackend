package database

import (
	"context"
	"fmt"
	"sync"

	"gorm.io/gorm"
)

// DatabaseWrapper wraps gorm.DB with cleanup functionality
type DatabaseWrapper struct {
	*gorm.DB
	mu     sync.RWMutex
	closed bool
}

// NewDatabaseWrapper creates a new database wrapper
func NewDatabaseWrapper(db *gorm.DB) *DatabaseWrapper {
	return &DatabaseWrapper{
		DB: db,
	}
}

// Close gracefully closes the database connection
func (dw *DatabaseWrapper) Close() error {
	dw.mu.Lock()
	defer dw.mu.Unlock()

	if dw.closed {
		return nil
	}

	// Get the underlying sql.DB and close it
	sqlDB, err := dw.DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB for closing: %w", err)
	}
	if err := sqlDB.Close(); err != nil {
		return fmt.Errorf("failed to close database: %w", err)
	}

	dw.closed = true
	return nil
}

// IsClosed returns whether the database connection is closed
func (dw *DatabaseWrapper) IsClosed() bool {
	dw.mu.RLock()
	defer dw.mu.RUnlock()
	return dw.closed
}

// Ping checks if the database connection is still alive
func (dw *DatabaseWrapper) Ping(ctx context.Context) error {
	dw.mu.RLock()
	defer dw.mu.RUnlock()

	if dw.closed {
		return fmt.Errorf("database connection is closed")
	}

	// Get the underlying sql.DB and ping it
	sqlDB, err := dw.DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB for ping: %w", err)
	}
	return sqlDB.PingContext(ctx)
}
