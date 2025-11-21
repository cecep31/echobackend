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

// Container holds the dependency injection container.
// It's recommended to use a struct to hold the container to avoid global variables.
type Container struct {
	*dig.Container
}

// NewContainer creates a new dependency injection container and registers all the dependencies.
func NewContainer(cfg *config.Config) (*Container, error) {
	container := &Container{dig.New()}

	if err := container.Provide(func() *config.Config { return cfg }); err != nil {
		return nil, err
	}

	if err := container.registerDatabase(); err != nil {
		return nil, err
	}

	if err := container.registerRepositories(); err != nil {
		return nil, err
	}

	if err := container.registerServices(); err != nil {
		return nil, err
	}

	if err := container.registerHandlers(); err != nil {
		return nil, err
	}

	if err := container.registerRoutes(); err != nil {
		return nil, err
	}

	// Provide cleanup manager
	if err := container.Provide(NewCleanupManager); err != nil {
		return nil, err
	}

	// Provide storage
	if err := container.Provide(storage.NewS3Storage); err != nil {
		return nil, err
	}

	return container, nil
}

func (c *Container) registerDatabase() error {
	if err := c.Provide(func(config *config.Config, cleanup *CleanupManager) *database.DatabaseWrapper {
		db := database.NewDatabase(config)
		cleanup.Register(db)
		return db
	}); err != nil {
		return err
	}

	return c.Provide(func(wrapper *database.DatabaseWrapper) *gorm.DB {
		return wrapper.DB
	})
}

func (c *Container) registerRepositories() error {
	if err := c.Provide(repository.NewUserRepository); err != nil {
		return err
	}
	if err := c.Provide(repository.NewPostRepository); err != nil {
		return err
	}
	if err := c.Provide(repository.NewAuthRepository); err != nil {
		return err
	}
	if err := c.Provide(repository.NewSessionRepository); err != nil {
		return err
	}
	if err := c.Provide(repository.NewTagRepository); err != nil {
		return err
	}
	if err := c.Provide(repository.NewPageRepository); err != nil {
		return err
	}
	if err := c.Provide(repository.NewWorkspaceRepository); err != nil {
		return err
	}
	if err := c.Provide(repository.NewCommentRepository); err != nil {
		return err
	}
	if err := c.Provide(repository.NewPostViewRepository); err != nil {
		return err
	}
	if err := c.Provide(repository.NewPostLikeRepository); err != nil {
		return err
	}
	if err := c.Provide(repository.NewUserFollowRepository); err != nil {
		return err
	}
	return c.Provide(repository.NewChatConversationRepository)
}

func (c *Container) registerServices() error {
	if err := c.Provide(service.NewUserService); err != nil {
		return err
	}
	if err := c.Provide(service.NewPostService); err != nil {
		return err
	}
	if err := c.Provide(service.NewAuthService); err != nil {
		return err
	}
	if err := c.Provide(service.NewTagService); err != nil {
		return err
	}
	if err := c.Provide(service.NewPageService); err != nil {
		return err
	}
	if err := c.Provide(service.NewWorkspaceService); err != nil {
		return err
	}
	if err := c.Provide(service.NewCommentService); err != nil {
		return err
	}
	if err := c.Provide(service.NewPostViewService); err != nil {
		return err
	}
	if err := c.Provide(service.NewPostLikeService); err != nil {
		return err
	}
	if err := c.Provide(service.NewUserFollowService); err != nil {
		return err
	}
	return c.Provide(service.NewChatConversationService)
}

func (c *Container) registerHandlers() error {
	if err := c.Provide(handler.NewUserHandler); err != nil {
		return err
	}
	if err := c.Provide(handler.NewPostHandler); err != nil {
		return err
	}
	if err := c.Provide(handler.NewAuthHandler); err != nil {
		return err
	}
	if err := c.Provide(handler.NewTagHandler); err != nil {
		return err
	}
	if err := c.Provide(handler.NewPageHandler); err != nil {
		return err
	}
	if err := c.Provide(handler.NewWorkspaceHandler); err != nil {
		return err
	}
	if err := c.Provide(handler.NewCommentHandler); err != nil {
		return err
	}
	if err := c.Provide(handler.NewPostViewHandler); err != nil {
		return err
	}
	if err := c.Provide(handler.NewPostLikeHandler); err != nil {
		return err
	}
	if err := c.Provide(handler.NewUserFollowHandler); err != nil {
		return err
	}
	return c.Provide(handler.NewChatConversationHandler)
}

func (c *Container) registerRoutes() error {
	if err := c.Provide(middleware.NewAuthMiddleware); err != nil {
		return err
	}
	return c.Provide(routes.NewRoutes)
}

// GetCleanupManager retrieves the cleanup manager from the container.
// This function is kept for convenience, but it's recommended to pass the container
// instance around instead of using this function.
func GetCleanupManager(container *Container) (*CleanupManager, error) {
	var cleanup *CleanupManager
	if err := container.Invoke(func(c *CleanupManager) {
		cleanup = c
	}); err != nil {
		return nil, err
	}
	return cleanup, nil
}
