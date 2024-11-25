package main

import (
	"echobackend/internal/di"
	"echobackend/internal/middleware"
	"echobackend/pkg/validator"
	"net/http"

	"github.com/labstack/echo/v4"
)

func main() {
	// Initialize dependency container
	container := di.NewContainer()

	// Initialize Echo
	e := echo.New()

	// Set custom validator
	e.Validator = validator.NewValidator()

	// Initialize handlers with dependencies
	routes := container.Routes()
	routes.Setup(e)

	e.GET("/", hellworld)

	// Setup middleware
	middleware.InitMiddleware(e)

	// load config
	conf := container.Config()
	// Start server
	port := conf.GetAppPort()
	e.Logger.Fatal(e.Start(":" + port))
}

func hellworld(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}
