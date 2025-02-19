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
