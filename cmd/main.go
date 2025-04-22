package main

import (
	"echobackend/config"
	"echobackend/internal/di"
	"echobackend/internal/middleware"
	"echobackend/internal/routes"
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
	container := di.BuildContainer(conf)

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

	// Start server
	// e.Logger.Printf("Starting server on port %s", conf.App_Port)
	e.Logger.Fatal(e.Start(":" + conf.App_Port))
}

func helloWorld(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]any{
		"message": "Hello, World!",
		"success": true,
	})
}
