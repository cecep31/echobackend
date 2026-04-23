package database

import (
	"fmt"
	"sync"

	"gorm.io/gorm"
)

// DatabaseWrapper wraps gorm.DB so it can be registered as a [Cleaner] (Close on shutdown)
// while still exposing *gorm.DB via embedding for repositories.
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
