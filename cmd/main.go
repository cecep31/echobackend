package main

import (
	"context"
	"echobackend/config"
	"echobackend/internal/di"
	"echobackend/internal/middleware"
	"echobackend/pkg/response"
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
	if conf.HTTP.TrustProxy {
		e.IPExtractor = echo.ExtractIPFromXFFHeader()
	} else {
		e.IPExtractor = echo.ExtractIPDirect()
	}

	// Set custom validator
	e.Validator = validator.NewValidator()

	// Initialize routes with manually wired dependencies
	container.Routes.Setup(e)

	e.GET("/", helloWorld)

	// Health check endpoint — used by Fly.io, Docker HEALTHCHECK, and load balancers.
	// Returns 200 when the DB is reachable, 503 otherwise.
	e.GET("/health", func(c *echo.Context) error {
		return healthCheck(c, container)
	})

	// Setup middleware
	middleware.InitMiddleware(e, conf)

	// Start server in a goroutine.
	// ReadTimeout covers the full request body read. For most endpoints 10 s is fine,
	// but file uploads can be slow on poor connections. We raise it to 60 s here and
	// rely on the 10 MB body limit (middleware) to bound abuse.
	// If you need tighter control per-route, wrap individual handlers with a context
	// deadline instead of changing the global server timeout.
	server := &http.Server{
		Addr:              ":" + conf.HTTP.Port,
		Handler:           e,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       60 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	go func() {
		e.Logger.Info("starting server", "port", conf.HTTP.Port)
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
	return response.Success(c, "Hello, World!", nil)
}

// healthCheck pings the database and returns 200 OK or 503 Service Unavailable.
func healthCheck(c *echo.Context, container *di.Container) error {
	ctx, cancel := context.WithTimeout(c.Request().Context(), 3*time.Second)
	defer cancel()

	if err := container.PingDB(ctx); err != nil {
		return c.JSON(http.StatusServiceUnavailable, map[string]string{
			"status": "unhealthy",
			"reason": "database unreachable",
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"status": "ok",
	})
}
