package main

import (
	"context"
	"echobackend/config"
	"echobackend/internal/di"
	"echobackend/internal/middleware"
	"echobackend/internal/routes"
	"echobackend/pkg/validator"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
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

	// Set custom validator
	e.Validator = validator.NewValidator()

	// Initialize routes with dependencies
	var newroutes *routes.Routes
	if err := container.Invoke(func(r *routes.Routes) {
		newroutes = r
	}); err != nil {
		panic(err)
	}
	newroutes.Setup(e)

	e.GET("/", helloWorld)

	// Setup middleware
	middleware.InitMiddleware(e, conf)

	// Start server in a goroutine
	go func() {
		e.Logger.Printf("Starting server on port %s", conf.AppPort)
		if err := e.Start(":" + conf.AppPort); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	e.Logger.Print("Server is shutting down...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Shutdown Echo server
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal("Server forced to shutdown:", err)
	}

	// Cleanup resources
	cleanup, err := di.GetCleanupManager(container)
	if err != nil {
		e.Logger.Error("Failed to get cleanup manager:", err)
	} else {
		if err := cleanup.CleanupWithTimeout(5 * time.Second); err != nil {
			e.Logger.Error("Cleanup failed:", err)
		} else {
			e.Logger.Print("Resources cleaned up successfully")
		}
	}

	e.Logger.Print("Server exited")
}

func helloWorld(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]any{
		"message": "Hello, World!",
		"success": true,
	})
}
