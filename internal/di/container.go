package di

import (
	"echobackend/config"
	"echobackend/internal/handler"
	"echobackend/internal/middleware"
	"echobackend/internal/repository"
	"echobackend/internal/routes"
	"echobackend/internal/service"
	"echobackend/internal/storage"
	"echobackend/pkg/database"

	"go.uber.org/dig"
)

var container *dig.Container

// BuildContainer initializes the dependency injection container
func BuildContainer(configgure *config.Config) *dig.Container {
	container = dig.New()

	// Provide config
	container.Provide(func() *config.Config {
		return configgure
	})

	// Provide database
	container.Provide(database.NewDatabase)

	// Provide repositories
	container.Provide(repository.NewUserRepository)
	container.Provide(repository.NewPostRepository)
	container.Provide(repository.NewAuthRepository)
	container.Provide(repository.NewTagRepository)
	container.Provide(repository.NewPageRepository)
	container.Provide(repository.NewWorkspaceRepository)

	// Provide services
	container.Provide(service.NewUserService)
	container.Provide(service.NewPostService)
	container.Provide(service.NewAuthService)
	container.Provide(service.NewTagService)
	container.Provide(service.NewPageService)
	container.Provide(service.NewWorkspaceService)

	// Provide storage
	container.Provide(storage.NewMinioStorage)

	// Provide middleware
	container.Provide(middleware.NewAuthMiddleware)

	// Provide handlers
	container.Provide(handler.NewUserHandler)
	container.Provide(handler.NewPostHandler)
	container.Provide(handler.NewAuthHandler)
	container.Provide(handler.NewTagHandler)
	container.Provide(handler.NewPageHandler)
	container.Provide(handler.NewWorkspaceHandler)

	// Provide routes
	container.Provide(routes.NewRoutes)

	return container
}

// GetContainer returns the dig container instance
func GetContainer() *dig.Container {
	return container
}
