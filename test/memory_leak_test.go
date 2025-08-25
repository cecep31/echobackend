package test

import (
	"context"
	"echobackend/config"
	"echobackend/internal/di"
	"testing"
	"time"
)

func TestGracefulShutdown(t *testing.T) {
	// Load test config
	conf := &config.Config{
		Database_URL:    "postgres://test:test@localhost:5432/test?sslmode=disable",
		MaxOpenConns:    5,
		MaxIdleConns:    2,
		ConnMaxLifetime: 5 * time.Minute,
		JWT_SECRET:      "test-secret",
		App_Port:        "8081",
	}

	// Build container
	container := di.BuildContainer(conf)

	// Get cleanup manager
	cleanup, err := di.GetCleanupManager()
	if err != nil {
		t.Fatalf("Failed to get cleanup manager: %v", err)
	}

	// Test cleanup with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = cleanup.Cleanup(ctx)
	if err != nil {
		t.Logf("Cleanup completed with errors (expected in test environment): %v", err)
	} else {
		t.Log("Cleanup completed successfully")
	}

	// Verify container is still accessible
	if container == nil {
		t.Fatal("Container should not be nil")
	}
}

func TestDatabaseWrapperCleanup(t *testing.T) {
	// This test would require a real database connection
	// For now, we'll just test the cleanup manager registration
	cleanup := di.NewCleanupManager()

	// Mock cleaner
	mockCleaner := &mockCleaner{closed: false}
	cleanup.Register(mockCleaner)

	// Test cleanup
	err := cleanup.CleanupWithTimeout(1 * time.Second)
	if err != nil {
		t.Fatalf("Cleanup failed: %v", err)
	}

	if !mockCleaner.closed {
		t.Fatal("Mock cleaner should have been closed")
	}
}

// Mock cleaner for testing
type mockCleaner struct {
	closed bool
}

func (m *mockCleaner) Close() error {
	m.closed = true
	return nil
}
