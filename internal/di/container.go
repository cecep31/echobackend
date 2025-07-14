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
	"gorm.io/gorm"
)

var container *dig.Container

// BuildContainer initializes the dependency injection container
func BuildContainer(configgure *config.Config) *dig.Container {
	container = dig.New()

	// Provide cleanup manager
	container.Provide(NewCleanupManager)

	// Provide config
	container.Provide(func() *config.Config {
		return configgure
	})

	// Provide database with cleanup registration
	container.Provide(func(config *config.Config, cleanup *CleanupManager) *database.DatabaseWrapper {
		db := database.NewDatabase(config)
		cleanup.Register(db)
		return db
	})

	// Provide *gorm.DB from the wrapper for repositories
	container.Provide(func(wrapper *database.DatabaseWrapper) *gorm.DB {
		return wrapper.DB
	})

	// Provide repositories
	container.Provide(repository.NewUserRepository)
	container.Provide(repository.NewPostRepository)
	container.Provide(repository.NewAuthRepository)
	container.Provide(repository.NewTagRepository)
	container.Provide(repository.NewPageRepository)
	container.Provide(repository.NewWorkspaceRepository)
	container.Provide(repository.NewCommentRepository)

	// Provide services
	container.Provide(service.NewUserService)
	container.Provide(service.NewPostService)
	container.Provide(service.NewAuthService)
	container.Provide(service.NewTagService)
	container.Provide(service.NewPageService)
	container.Provide(service.NewWorkspaceService)
	container.Provide(service.NewCommentService)

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
	container.Provide(handler.NewCommentHandler)

	// Provide routes
	container.Provide(routes.NewRoutes)

	return container
}

// GetContainer returns the dig container instance
func GetContainer() *dig.Container {
	return container
}

// GetCleanupManager retrieves the cleanup manager from the container
func GetCleanupManager() (*CleanupManager, error) {
	var cleanup *CleanupManager
	if err := container.Invoke(func(c *CleanupManager) {
		cleanup = c
	}); err != nil {
		return nil, err
	}
	return cleanup, nil
}
