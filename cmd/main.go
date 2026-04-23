package main

import (
	"context"
	"echobackend/config"
	"echobackend/internal/di"
	"echobackend/internal/middleware"
	"echobackend/pkg/validator"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v5"
)

func main() {
	// load config
	conf, errconf := config.Load()
	if errconf != nil {
		panic(errconf)
	}

	// Initialize dependency container
	container, err := di.NewContainer(conf)
	if err != nil {
		panic(err)
	}

	// Initialize Echo
	e := echo.New()
	if conf.HTTPTrustProxy {
		e.IPExtractor = echo.ExtractIPFromXFFHeader()
	} else {
		e.IPExtractor = echo.ExtractIPDirect()
	}

	// Set custom validator
	e.Validator = validator.NewValidator()

	// Initialize routes with manually wired dependencies
	container.Routes.Setup(e)

	e.GET("/", helloWorld)

	// Setup middleware
	middleware.InitMiddleware(e, conf)

	// Start server in a goroutine
	server := &http.Server{
		Addr:         ":" + conf.HTTPPort,
		Handler:      e,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		e.Logger.Info("starting server", "port", conf.HTTPPort)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			e.Logger.Error("server exited unexpectedly", "error", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	e.Logger.Info("server is shutting down")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Shutdown server
	if err := server.Shutdown(ctx); err != nil {
		e.Logger.Error("server forced to shutdown", "error", err)
	}

	// Cleanup resources
	cleanup, err := di.GetCleanupManager(container)
	if err != nil {
		e.Logger.Error("failed to get cleanup manager", "error", err)
	} else if cleanup != nil {
		if err := cleanup.CleanupWithTimeout(5 * time.Second); err != nil {
			e.Logger.Error("cleanup failed", "error", err)
		} else {
			e.Logger.Info("resources cleaned up successfully")
		}
	}

	e.Logger.Info("server exited")
}

func helloWorld(c *echo.Context) error {
	return c.JSON(http.StatusOK, map[string]any{
		"message": "Hello, World!",
		"success": true,
	})
}
