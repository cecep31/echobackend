package main

import (
	"echobackend/internal/di"
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
	conf := container.Config()
	e.GET("/", hellworld)
	// Setup middleware

	// Start server
	// Start server
	port := conf.GetAppPort()
	e.Logger.Fatal(e.Start(":" + port))
}

func hellworld(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{"message": "Hello World!"})
}
