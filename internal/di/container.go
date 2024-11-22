package di

import (
	"echobackend/internal/config"
	"echobackend/internal/handler"
	"echobackend/internal/repository"
	"echobackend/internal/routes"
	"echobackend/internal/service"
	"echobackend/pkg/database"
	"sync"

	"gorm.io/gorm"
)

// Container holds all dependencies
type Container struct {
	config      *config.Config
	db          *gorm.DB
	postRepo    repository.PostRepository
	postService service.PostService
	userRepo    repository.UserRepository
	userService service.UserService
	routes      *routes.Routes
	userHandler *handler.UserHandler
	postHandler *handler.PostHandler
	// Add other dependencies here

	once sync.Once
}

// NewContainer creates a new dependency injection container
func NewContainer() *Container {
	return &Container{}
}

func (c *Container) UserHandler() *handler.UserHandler {
	if c.userHandler == nil {
		c.userHandler = handler.NewUserHandler(c.UserServices())
	}
	return c.userHandler
}

func (c *Container) PostHandler() *handler.PostHandler {
	if c.postHandler == nil {
		c.postHandler = handler.NewPostHandler(c.PostService())
	}
	return c.postHandler
}

func (c *Container) Routes() *routes.Routes {
	if c.routes == nil {
		c.routes = routes.NewRoutes(c.UserHandler(), c.PostHandler())
	}
	return c.routes
}

func (c *Container) UserServices() service.UserService {
	if c.userService == nil {
		c.userService = service.NewUserService(c.UserRepository())
	}
	return c.userService
}

// Config returns the application configuration
func (c *Container) Config() *config.Config {
	c.once.Do(func() {
		conf, err := config.Load()
		if err != nil {
			panic(err)
		}
		c.config = conf
	})
	return c.config
}

// Database returns the database instance
func (c *Container) Database() *gorm.DB {
	if c.db == nil {
		db, err := database.SetupDatabase(c.Config())
		if err != nil {
			panic(err)
		}
		c.db = db
	}
	return c.db
}

// PostRepository returns the post repository instance
func (c *Container) PostRepository() repository.PostRepository {
	if c.postRepo == nil {
		c.postRepo = repository.NewPostRepository(c.Database())
	}
	return c.postRepo
}

// PostService returns the post service instance
func (c *Container) PostService() service.PostService {
	if c.postService == nil {
		c.postService = service.NewPostService(c.PostRepository())
	}
	return c.postService
}

// UserRepository returns the user repository instance
func (c *Container) UserRepository() repository.UserRepository {
	if c.userRepo == nil {
		c.userRepo = repository.NewUserRepository(c.Database())
	}
	return c.userRepo
}
