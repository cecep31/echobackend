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
	"sync"

	"gorm.io/gorm"
)

// Container holds all dependencies
type Container struct {
	sync.Mutex
	config         *config.Config
	db             *gorm.DB
	userRepo       repository.UserRepository
	postRepo       repository.PostRepository
	postService    service.PostService
	userService    service.UserService
	authRepo       repository.AuthRepository
	authService    service.AuthService
	tagRepo        repository.TagRepository
	tagService     service.TagService
	routes         *routes.Routes
	userHandler    *handler.UserHandler
	postHandler    *handler.PostHandler
	authHandler    *handler.AuthHandler
	tagHandler     *handler.TagHandler
	authMiddleware *middleware.AuthMiddleware
	miniostorage   *storage.MinioStorage
	// Add other dependencies here
}

// NewContainer creates a new dependency injection container
func NewContainer(config *config.Config) *Container {
	return &Container{config: config}
}

func (c *Container) AuthMiddleware() *middleware.AuthMiddleware {
	if c.authMiddleware == nil {
		c.authMiddleware = middleware.NewAuthMiddleware(c.Config())
	}
	return c.authMiddleware
}

func (c *Container) UserHandler() *handler.UserHandler {
	if c.userHandler == nil {
		c.userHandler = handler.NewUserHandler(c.UserServices())
	}
	return c.userHandler
}

func (c *Container) TagHandler() *handler.TagHandler {
	if c.tagHandler == nil {
		c.tagHandler = handler.NewTagHandler(c.TagService())
	}
	return c.tagHandler
}

func (c *Container) PostHandler() *handler.PostHandler {
	if c.postHandler == nil {
		c.postHandler = handler.NewPostHandler(c.PostService())
	}
	return c.postHandler
}

func (c *Container) AuthHandler() *handler.AuthHandler {

	if c.authHandler == nil {
		c.authHandler = handler.NewAuthHandler(c.AuthService())
	}
	return c.authHandler
}

func (c *Container) Routes() *routes.Routes {

	if c.routes == nil {
		c.routes = routes.NewRoutes(c.UserHandler(), c.PostHandler(), c.AuthHandler(), c.AuthMiddleware(), c.TagHandler())
	}
	return c.routes
}

func (c *Container) UserServices() service.UserService {

	if c.userService == nil {
		c.userService = service.NewUserService(c.UserRepository())
	}
	return c.userService
}

func (c *Container) TagService() service.TagService {

	if c.tagService == nil {
		c.tagService = service.NewTagService(c.TagRepository())
	}
	return c.tagService
}

func (c *Container) TagRepository() repository.TagRepository {

	if c.tagRepo == nil {
		c.tagRepo = repository.NewTagRepository(c.Database())
	}
	return c.tagRepo
}

// Config returns the application configuration
func (c *Container) Config() *config.Config {

	if c.config == nil {
		conf, err := config.Load()
		if err != nil {
			panic(err)
		}
		c.config = conf
	}
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
		c.postService = service.NewPostService(c.PostRepository(), c.MinioStorage())
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

func (c *Container) AuthRepository() repository.AuthRepository {

	if c.authRepo == nil {
		c.authRepo = repository.NewAuthRepository(c.Database())
	}
	return c.authRepo
}

func (c *Container) AuthService() service.AuthService {

	if c.authService == nil {
		c.authService = service.NewAuthService(c.AuthRepository(), c.Config())
	}
	return c.authService
}

func (c *Container) MinioStorage() *storage.MinioStorage {

	if c.miniostorage == nil {
		c.miniostorage = storage.NewMinioStorage(c.Config())
	}
	return c.miniostorage
}
