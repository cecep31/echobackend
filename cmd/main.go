package main

import (
	"echobackend/config"
	"echobackend/internal/di"
	"echobackend/internal/middleware"
	"echobackend/internal/routes"
	"echobackend/pkg/validator"
	"net/http"
	"net/http/pprof"
	"os"
	"runtime"

	"github.com/labstack/echo/v4"
)

func main() {
	// Set GOMAXPROCS to match available CPU cores
	runtime.GOMAXPROCS(runtime.NumCPU())

	// load config
	conf, errconf := config.Load()
	if errconf != nil {
		panic(errconf)
	}

	// Initialize dependency container
	container := di.BuildContainer(conf)

	// Initialize Echo
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

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

	// Add debug endpoints in non-production environments
	if os.Getenv("APP_ENV") != "production" {
		debug := e.Group("/debug")
		debug.GET("/pprof/*", echo.WrapHandler(http.HandlerFunc(pprof.Index)))
		debug.GET("/pprof/cmdline", echo.WrapHandler(http.HandlerFunc(pprof.Cmdline)))
		debug.GET("/pprof/profile", echo.WrapHandler(http.HandlerFunc(pprof.Profile)))
		debug.GET("/pprof/symbol", echo.WrapHandler(http.HandlerFunc(pprof.Symbol)))
		debug.GET("/pprof/trace", echo.WrapHandler(http.HandlerFunc(pprof.Trace)))
		debug.GET("/pprof/heap", echo.WrapHandler(http.HandlerFunc(pprof.Handler("heap").ServeHTTP)))
		debug.GET("/pprof/goroutine", echo.WrapHandler(http.HandlerFunc(pprof.Handler("goroutine").ServeHTTP)))
		debug.GET("/pprof/allocs", echo.WrapHandler(http.HandlerFunc(pprof.Handler("allocs").ServeHTTP)))
		debug.GET("/pprof/block", echo.WrapHandler(http.HandlerFunc(pprof.Handler("block").ServeHTTP)))
		debug.GET("/pprof/mutex", echo.WrapHandler(http.HandlerFunc(pprof.Handler("mutex").ServeHTTP)))
	}

	// Start server
	e.Logger.Printf("Starting server on port %s", conf.App_Port)
	e.Logger.Fatal(e.Start(":" + conf.App_Port))
}

func hellworld(c echo.Context) error {
	return c.JSON(http.StatusOK, &echo.Map{
		"message": "Hello, World!",
	})
}
