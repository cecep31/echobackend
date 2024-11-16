package routes

import (
	"echobackend/internal/handler"
	"log"

	"github.com/labstack/echo/v4"
)

type Routes struct {
	userHandler *handler.UserHandler
	// Add other handlers
}

func NewRoutes(
	userHandler *handler.UserHandler,
) *Routes {
	return &Routes{
		userHandler: userHandler,
	}
}

func (r *Routes) Setup(e *echo.Echo) {
	// API Group
	log.Println("setup v1 routes")
	v1 := e.Group("/v1")
	r.setupV1Routes(v1)
}

func (r *Routes) setupV1Routes(v1 *echo.Group) {
	r.setupUserRoutes(v1)
}

func (r *Routes) setupUserRoutes(v1 *echo.Group) {
	users := v1.Group("/users")

	{
		users.GET("/:id", r.userHandler.GetByID)
		users.GET("", r.userHandler.GetUsers)
	}
}
