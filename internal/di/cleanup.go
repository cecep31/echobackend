package di

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// Cleaner interface for resources that need cleanup
type Cleaner interface {
	Close() error
}

// CleanupManager manages cleanup of resources
type CleanupManager struct {
	cleaners []Cleaner
	mu       sync.RWMutex
}

// NewCleanupManager creates a new cleanup manager
func NewCleanupManager() *CleanupManager {
	return &CleanupManager{
		cleaners: make([]Cleaner, 0),
	}
}

// Register adds a cleaner to the cleanup manager
func (cm *CleanupManager) Register(cleaner Cleaner) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.cleaners = append(cm.cleaners, cleaner)
}

// Cleanup closes all registered cleaners in reverse order
func (cm *CleanupManager) Cleanup(ctx context.Context) error {
	cm.mu.RLock()
	cleaners := make([]Cleaner, len(cm.cleaners))
	copy(cleaners, cm.cleaners)
	cm.mu.RUnlock()

	var errors []error

	// Cleanup in reverse order (LIFO)
	for i := len(cleaners) - 1; i >= 0; i-- {
		select {
		case <-ctx.Done():
			return fmt.Errorf("cleanup timeout: %w", ctx.Err())
		default:
			if err := cleaners[i].Close(); err != nil {
				errors = append(errors, fmt.Errorf("cleanup error: %w", err))
			}
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("cleanup completed with %d errors: %v", len(errors), errors)
	}

	return nil
}

// CleanupWithTimeout performs cleanup with a timeout
func (cm *CleanupManager) CleanupWithTimeout(timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return cm.Cleanup(ctx)
}
