package routes

import (
	"echobackend/internal/handler"

	"github.com/labstack/echo/v4"
)

type Routes struct {
	userHandler *handler.UserHandler
	postHandler *handler.PostHandler
	// Add other handlers
}

func NewRoutes(
	userHandler *handler.UserHandler,
	postHandler *handler.PostHandler,
) *Routes {
	return &Routes{
		userHandler: userHandler,
		postHandler: postHandler,
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
