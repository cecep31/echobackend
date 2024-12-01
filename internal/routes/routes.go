package routes

import (
	"echobackend/internal/handler"
	"echobackend/internal/middleware"

	"github.com/labstack/echo/v4"
)

type Routes struct {
	userHandler    *handler.UserHandler
	postHandler    *handler.PostHandler
	authHandler    *handler.AuthHandler
	authMiddleware *middleware.AuthMiddleware
	// Add other handlers
}

func NewRoutes(
	userHandler *handler.UserHandler,
	postHandler *handler.PostHandler,
	authHandler *handler.AuthHandler,
	authMiddleware *middleware.AuthMiddleware,
) *Routes {
	return &Routes{
		userHandler:    userHandler,
		postHandler:    postHandler,
		authHandler:    authHandler,
		authMiddleware: authMiddleware,
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
}

func (r *Routes) setupUserRoutes(v1 *echo.Group) {
	users := v1.Group("/users", r.authMiddleware.Auth())
	{
		users.GET("/:id", r.userHandler.GetByID)
		users.GET("", r.userHandler.GetUsers)
		users.DELETE("/:id", r.userHandler.DeleteUser, r.authMiddleware.AuthAdmin())
	}
}

func (r *Routes) setupPostRoutes(v1 *echo.Group) {
	posts := v1.Group("/posts")
	{
		posts.GET("/u/:username/:slug", r.postHandler.GetPostBySlugAndUsername)
		posts.GET("", r.postHandler.GetPosts)
		posts.DELETE("/:id", r.postHandler.DeletePost, r.authMiddleware.Auth())
		posts.GET("/random", r.postHandler.GetPostsRandom)
		posts.GET("/:id", r.postHandler.GetPost)
		posts.GET("/username/:username", r.postHandler.GetPostsByUsername)
	}
}

func (r *Routes) setupAuthRoutes(v1 *echo.Group) {
	auth := v1.Group("/auth")
	{
		auth.POST("/register", r.authHandler.Register)
		auth.POST("/login", r.authHandler.Login)
	}
}
