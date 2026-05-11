package di

import (
	"context"
	"echobackend/config"
	"echobackend/internal/handler"
	"echobackend/internal/middleware"
	"echobackend/internal/repository"
	"echobackend/internal/routes"
	"echobackend/internal/service"
	"echobackend/pkg/cache"
	"echobackend/pkg/database"
	"echobackend/pkg/market"
	"echobackend/pkg/storage"

	"gorm.io/gorm"
)

// Container holds the manually wired application dependencies.
type Container struct {
	Config  *config.Config
	Cleanup *CleanupManager
	Routes  *routes.Routes
	db      *gorm.DB
}

// NewContainer creates a manually wired application container.
func NewContainer(cfg *config.Config) (*Container, error) {
	cleanup := NewCleanupManager()

	db := database.NewDatabase(cfg)
	cleanup.Register(func() error {
		sqlDB, err := db.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	})

	valkeyCache := cache.NewValkeyCache(cfg)
	if valkeyCache != nil {
		cleanup.Register(func() error {
			return valkeyCache.Close()
		})
	}

	s3Storage := storage.NewS3Storage(cfg)

	userRepo := repository.NewUserRepository(db)
	postRepo := repository.NewPostRepository(db)
	authRepo := repository.NewAuthRepository(db)
	sessionRepo := repository.NewSessionRepository(db)
	tagRepo := repository.NewTagRepository(db)
	commentRepo := repository.NewCommentRepository(db)
	postViewRepo := repository.NewPostViewRepository(db)
	postLikeRepo := repository.NewPostLikeRepository(db)
	userFollowRepo := repository.NewUserFollowRepository(db)
	chatConversationRepo := repository.NewChatConversationRepository(db)
	holdingRepo := repository.NewHoldingRepository(db)

	userService := service.NewUserService(userRepo)
	tagService := service.NewTagService(tagRepo)
	postService := service.NewPostService(postRepo, tagService, s3Storage, valkeyCache)
	authService := service.NewAuthService(authRepo, userRepo, sessionRepo, cfg)
	commentService := service.NewCommentService(commentRepo, postRepo)
	postViewService := service.NewPostViewService(postViewRepo, postRepo)
	postLikeService := service.NewPostLikeService(postLikeRepo, postRepo)
	userFollowService := service.NewUserFollowService(userFollowRepo, userRepo)
	chatConversationService := service.NewChatConversationService(chatConversationRepo)
	yahooClient := market.NewYahooClient(nil)
	holdingService := service.NewHoldingService(holdingRepo, yahooClient)

	userHandler := handler.NewUserHandler(userService, userFollowService)
	postHandler := handler.NewPostHandler(postService, postViewService)
	authHandler := handler.NewAuthHandler(authService)
	tagHandler := handler.NewTagHandler(tagService)
	commentHandler := handler.NewCommentHandler(commentService)
	postViewHandler := handler.NewPostViewHandler(postViewService)
	postLikeHandler := handler.NewPostLikeHandler(postLikeService)
	userFollowHandler := handler.NewUserFollowHandler(userFollowService)
	chatConversationHandler := handler.NewChatConversationHandler(chatConversationService)
	holdingHandler := handler.NewHoldingHandler(holdingService)

	authMiddleware := middleware.NewAuthMiddleware(cfg, userRepo)
	appRoutes := routes.NewRoutes(
		cfg,
		userHandler,
		postHandler,
		authHandler,
		authMiddleware,
		tagHandler,
		commentHandler,
		postViewHandler,
		postLikeHandler,
		userFollowHandler,
		chatConversationHandler,
		holdingHandler,
	)

	return &Container{
		Config:  cfg,
		Cleanup: cleanup,
		Routes:  appRoutes,
		db:      db,
	}, nil
}

// PingDB checks that the database connection is alive.
func (c *Container) PingDB(ctx context.Context) error {
	sqlDB, err := c.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.PingContext(ctx)
}

// GetCleanupManager retrieves the cleanup manager from the container.
func GetCleanupManager(container *Container) (*CleanupManager, error) {
	if container == nil {
		return nil, nil
	}
	return container.Cleanup, nil
}
