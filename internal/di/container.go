package di

import (
	"echobackend/config"
	"echobackend/internal/handler"
	"echobackend/internal/middleware"
	"echobackend/internal/repository"
	"echobackend/internal/routes"
	"echobackend/internal/service"
	"echobackend/pkg/database"
	"echobackend/pkg/storage"

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

	// Register repositories
	repositories := []any{
		repository.NewUserRepository,
		repository.NewPostRepository,
		repository.NewAuthRepository,
		repository.NewTagRepository,
		repository.NewPageRepository,
		repository.NewWorkspaceRepository,
		repository.NewCommentRepository,
		repository.NewPostViewRepository,
		repository.NewPostLikeRepository,
		repository.NewUserFollowRepository,
	}
	for _, repo := range repositories {
		container.Provide(repo)
	}

	// Register services
	services := []any{
		service.NewUserService,
		service.NewPostService,
		service.NewAuthService,
		service.NewTagService,
		service.NewPageService,
		service.NewWorkspaceService,
		service.NewCommentService,
		service.NewPostViewService,
		service.NewPostLikeService,
		service.NewUserFollowService,
	}
	for _, svc := range services {
		container.Provide(svc)
	}

	// Register infrastructure
	container.Provide(storage.NewS3Storage)
	container.Provide(middleware.NewAuthMiddleware)
	container.Provide(routes.NewRoutes)

	// Register handlers
	handlers := []any{
		handler.NewUserHandler,
		handler.NewPostHandler,
		handler.NewAuthHandler,
		handler.NewTagHandler,
		handler.NewPageHandler,
		handler.NewWorkspaceHandler,
		handler.NewCommentHandler,
		handler.NewPostViewHandler,
		handler.NewPostLikeHandler,
		handler.NewUserFollowHandler,
	}
	for _, hdl := range handlers {
		container.Provide(hdl)
	}

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
