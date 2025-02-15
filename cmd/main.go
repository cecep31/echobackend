package main

import (
	"echobackend/config"
	"echobackend/internal/di"
	"echobackend/internal/middleware"
	"echobackend/pkg/validator"
	"net/http"

	"github.com/labstack/echo/v4"
)

func main() {

	// load config
	conf, errconf := config.Load()
	if errconf != nil {
		panic(errconf)
	}

	// Initialize dependency container
	container := di.NewContainer(conf)

	// Initialize Echo
	e := echo.New()

	// Set custom validator
	e.Validator = validator.NewValidator()

	// Initialize handlers with dependencies
	routes := container.Routes()
	routes.Setup(e)

	e.GET("/", hellworld)

	// Setup middleware
	middleware.InitMiddleware(e, conf)

	// Start server
	port := conf.GetAppPort()
	e.Logger.Fatal(e.Start(":" + port))
}

func hellworld(c echo.Context) error {
	return c.JSON(http.StatusOK, &echo.Map{
		"message": "Hello, World!",
	})
}
