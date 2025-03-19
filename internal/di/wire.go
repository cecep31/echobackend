//go:build wireinject
// +build wireinject

package di

import (
	"context"
	"echobackend/config"
	"echobackend/internal/handler"
	"echobackend/internal/middleware"
	"echobackend/internal/repository"
	"echobackend/internal/routes"
	"echobackend/internal/service"
	"echobackend/internal/storage"
	"echobackend/pkg/database"
	"echobackend/pkg/validator"
	"fmt"
	"net/http"

	"github.com/google/wire"
	"github.com/labstack/echo/v4"
)

// Application represents the application with all its dependencies
type Application struct {
	Echo     *echo.Echo
	Config   *config.Config
	Routes   *routes.Routes
	Shutdown func(context.Context) error
}

// Start starts the application server
func (app *Application) Start() {
	// Start server in a goroutine so it doesn't block
	go func() {
		addr := fmt.Sprintf(":%s", app.Config.App_Port)
		if err := app.Echo.Start(addr); err != nil && err != http.ErrServerClosed {
			app.Echo.Logger.Fatalf("shutting down the server: %v", err)
		}
	}()
	app.Echo.Logger.Printf("Starting server on port %s", app.Config.App_Port)
}

// Stop gracefully shuts down the application
func (app *Application) Stop(ctx context.Context) error {
	return app.Shutdown(ctx)
}

// Run starts the application and blocks until it's stopped
func (app *Application) Run() {
	app.Start()
	// Block forever
	select {}
}

// ProvideEcho provides an Echo instance
func ProvideEcho() *echo.Echo {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true
	e.Validator = validator.NewValidator()
	return e
}

// ProvideApplication provides the Application struct
func ProvideApplication(
	e *echo.Echo,
	conf *config.Config,
	routes *routes.Routes,
) *Application {
	// Setup routes
	routes.Setup(e)

	// Add root route
	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, &echo.Map{
			"message": "Hello, World!",
		})
	})

	// Setup middleware
	middleware.InitMiddleware(e, conf)

	return &Application{
		Echo:     e,
		Config:   conf,
		Routes:   routes,
		Shutdown: e.Shutdown,
	}
}

// RepositorySet is a Wire provider set for repositories
var RepositorySet = wire.NewSet(
	repository.NewUserRepository,
	repository.NewPostRepository,
	repository.NewAuthRepository,
	repository.NewTagRepository,
	repository.NewPageRepository,
	repository.NewWorkspaceRepository,
)

// ServiceSet is a Wire provider set for services
var ServiceSet = wire.NewSet(
	service.NewUserService,
	service.NewPostService,
	service.NewAuthService,
	service.NewTagService,
	service.NewPageService,
	service.NewWorkspaceService,
)

// HandlerSet is a Wire provider set for handlers
var HandlerSet = wire.NewSet(
	handler.NewUserHandler,
	handler.NewPostHandler,
	handler.NewAuthHandler,
	handler.NewTagHandler,
	handler.NewPageHandler,
	handler.NewWorkspaceHandler,
)

// InfrastructureSet is a Wire provider set for infrastructure components
var InfrastructureSet = wire.NewSet(
	ProvideEcho,
	database.NewDatabase,
	storage.NewMinioStorage,
	middleware.NewAuthMiddleware,
	routes.NewRoutes,
)

// ApplicationSet is a Wire provider set for the application
var ApplicationSet = wire.NewSet(
	ProvideApplication,
)

// BuildApplication builds the application with all dependencies
func BuildApplication(conf *config.Config) (*Application, error) {
	wire.Build(
		InfrastructureSet,
		RepositorySet,
		ServiceSet,
		HandlerSet,
		ApplicationSet,
	)
	return nil, nil
}
