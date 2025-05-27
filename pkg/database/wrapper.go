package database

import (
	"context"
	"fmt"
	"sync"

	"github.com/uptrace/bun"
)

// DatabaseWrapper wraps bun.DB with cleanup functionality
type DatabaseWrapper struct {
	*bun.DB
	mu     sync.RWMutex
	closed bool
}

// NewDatabaseWrapper creates a new database wrapper
func NewDatabaseWrapper(db *bun.DB) *DatabaseWrapper {
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

	// Close the underlying SQL database
	if err := dw.DB.Close(); err != nil {
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

	return dw.DB.PingContext(ctx)
}
