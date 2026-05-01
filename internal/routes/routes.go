package routes

import (
	"echobackend/config"
	"echobackend/internal/handler"
	"echobackend/internal/middleware"

	"github.com/labstack/echo/v5"
)

type Routes struct {
	config                  *config.Config
	userHandler             *handler.UserHandler
	postHandler             *handler.PostHandler
	authHandler             *handler.AuthHandler
	authMiddleware          *middleware.AuthMiddleware
	tagHandler              *handler.TagHandler
	commentHandler          *handler.CommentHandler
	postViewHandler         *handler.PostViewHandler
	postLikeHandler         *handler.PostLikeHandler
	userFollowHandler       *handler.UserFollowHandler
	chatConversationHandler *handler.ChatConversationHandler
}

func NewRoutes(
	config *config.Config,
	userHandler *handler.UserHandler,
	postHandler *handler.PostHandler,
	authHandler *handler.AuthHandler,
	authMiddleware *middleware.AuthMiddleware,
	tagHandler *handler.TagHandler,
	commentHandler *handler.CommentHandler,
	postViewHandler *handler.PostViewHandler,
	postLikeHandler *handler.PostLikeHandler,
	userFollowHandler *handler.UserFollowHandler,
	chatConversationHandler *handler.ChatConversationHandler,
) *Routes {
	return &Routes{
		config:                  config,
		userHandler:             userHandler,
		postHandler:             postHandler,
		authHandler:             authHandler,
		authMiddleware:          authMiddleware,
		tagHandler:              tagHandler,
		commentHandler:          commentHandler,
		postViewHandler:         postViewHandler,
		postLikeHandler:         postLikeHandler,
		userFollowHandler:       userFollowHandler,
		chatConversationHandler: chatConversationHandler,
	}
}

func (r *Routes) Setup(e *echo.Echo) {
	// API Group
	api := e.Group("/api")
	r.setupAPIRoutes(api)
}

func (r *Routes) setupAPIRoutes(api *echo.Group) {
	r.setupUserRoutes(api)
	r.setupPostRoutes(api)
	r.setupAuthRoutes(api)
	r.setupTagRoutes(api)
	r.setupChatConversationRoutes(api)
	if r.config.AppDebug {
		r.setupDebugRoutes(api)
	}
}

func (r *Routes) setupChatConversationRoutes(api *echo.Group) {
	conversations := api.Group("/chat/conversations")
	{
		conversations.POST("", r.chatConversationHandler.CreateConversation, r.authMiddleware.Auth())
		conversations.GET("", r.chatConversationHandler.GetConversations, r.authMiddleware.Auth())
		conversations.GET("/:id", r.chatConversationHandler.GetConversation, r.authMiddleware.Auth())
		conversations.PUT("/:id", r.chatConversationHandler.UpdateConversation, r.authMiddleware.Auth())
		conversations.DELETE("/:id", r.chatConversationHandler.DeleteConversation, r.authMiddleware.Auth())
	}
}
