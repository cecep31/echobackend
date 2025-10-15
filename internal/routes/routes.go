package routes

import (
	"echobackend/config"
	"echobackend/internal/handler"
	"echobackend/internal/middleware"

	"github.com/labstack/echo/v4"
)

type Routes struct {
	config            *config.Config
	userHandler       *handler.UserHandler
	postHandler       *handler.PostHandler
	authHandler       *handler.AuthHandler
	authMiddleware    *middleware.AuthMiddleware
	tagHandler        *handler.TagHandler
	pageHandler       *handler.PageHandler
	workspaceHandler  *handler.WorkspaceHandler
	commentHandler    *handler.CommentHandler
	postViewHandler   *handler.PostViewHandler
	postLikeHandler   *handler.PostLikeHandler
	userFollowHandler *handler.UserFollowHandler
}

func NewRoutes(
	config *config.Config,
	userHandler *handler.UserHandler,
	postHandler *handler.PostHandler,
	authHandler *handler.AuthHandler,
	authMiddleware *middleware.AuthMiddleware,
	tagHandler *handler.TagHandler,
	pageHandler *handler.PageHandler,
	workspaceHandler *handler.WorkspaceHandler,
	commentHandler *handler.CommentHandler,
	postViewHandler *handler.PostViewHandler,
	postLikeHandler *handler.PostLikeHandler,
	userFollowHandler *handler.UserFollowHandler,
) *Routes {
	return &Routes{
		config:            config,
		userHandler:       userHandler,
		postHandler:       postHandler,
		authHandler:       authHandler,
		authMiddleware:    authMiddleware,
		tagHandler:        tagHandler,
		pageHandler:       pageHandler,
		workspaceHandler:  workspaceHandler,
		commentHandler:    commentHandler,
		postViewHandler:   postViewHandler,
		postLikeHandler:   postLikeHandler,
		userFollowHandler: userFollowHandler,
	}
}

func (r *Routes) Setup(e *echo.Echo) {
	// API Group
	v1 := e.Group("/v1")
	r.setupV1Routes(v1)
}

func (r *Routes) setupV1Routes(v1 *echo.Group) {
	r.setupUserRoutes(v1)
	r.setupPostRoutes(v1)
	r.setupAuthRoutes(v1)
	r.setupTagRoutes(v1)
	r.setupPageRoutes(v1)
	r.setupWorkspaceRoutes(v1)
	if r.config.DEBUG {
		r.setupDebugRoutes(v1)
	}
}
