package routes

import (
	"echobackend/internal/handler"
	"echobackend/internal/middleware"

	"github.com/labstack/echo/v4"
)

type Routes struct {
	userHandler    *handler.UserHandler
	postHandler    *handler.PostHandler
	authMiddleware *middleware.AuthMiddleware
	// Add other handlers
}

func NewRoutes(
	userHandler *handler.UserHandler,
	postHandler *handler.PostHandler,
	authMiddleware *middleware.AuthMiddleware,
) *Routes {
	return &Routes{
		userHandler:    userHandler,
		postHandler:    postHandler,
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

}

func (r *Routes) setupUserRoutes(v1 *echo.Group) {
	users := v1.Group("/users")

	{
		users.GET("/:id", r.userHandler.GetByID)
		users.GET("", r.userHandler.GetUsers, r.authMiddleware.Auth())
	}
}

func (r *Routes) setupPostRoutes(v1 *echo.Group) {
	posts := v1.Group("/posts")

	{
		posts.GET("", r.postHandler.GetPosts)
		posts.GET("/random", r.postHandler.GetPostsRandom)
	}
}
