package main

import (
	"echobackend/internal/di"
	"echobackend/internal/middleware"
	"net/http"

	"github.com/labstack/echo/v4"
)

func main() {
	// Initialize dependency container
	container := di.NewContainer()

	// Initialize Echo
	e := echo.New()

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
	return c.JSON(http.StatusOK, map[string]string{"message": "Hello World!"})
}
