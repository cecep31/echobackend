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
	}
}

func (r *Routes) setupPostRoutes(v1 *echo.Group) {
	posts := v1.Group("/posts")

	{
		posts.GET("", r.postHandler.GetPosts)
		posts.GET("/random", r.postHandler.GetPostsRandom)
	}
}

func (r *Routes) setupAuthRoutes(v1 *echo.Group) {
	auth := v1.Group("/auth")
	{
		auth.POST("/register", r.authHandler.Register)
		auth.POST("/login", r.authHandler.Login)
	}
}
