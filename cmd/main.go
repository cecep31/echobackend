package main

import (
	"echobackend/internal/config"
	"echobackend/internal/handler"
	"echobackend/internal/middleware"
	"echobackend/internal/repository"
	"echobackend/internal/routes"
	"echobackend/internal/service"
	"echobackend/pkg/database"
	"net/http"

	"github.com/labstack/echo/v4"
)

func main() {
	cfg, errconf := config.Load()
	if errconf != nil {
		panic(errconf)
	}

	// Initialize database
	db := database.SetupDatabase(cfg)
	// Initialize dependencies
	userRepo := repository.NewUserRepository(db)
	postrepo := repository.NewPostRepository(db)
	// Initialize services
	userService := service.NewUserService(userRepo)
	postService := service.NewPostService(postrepo)
	// Initialize handlers
	userHandler := handler.NewUserHandler(userService) //
	postHandler := handler.NewPostHandler(postService)
	//setupMiddleware
	authmid := middleware.NewAuthMiddleware(cfg)

	// Initialize server
	routes := routes.NewRoutes(userHandler, postHandler, authmid)

	e := echo.New()

	routes.Setup(e)

	middleware.InitMiddleware(e)
	e.GET("/", hellworld)
	port := cfg.GetAppPort()
	e.Logger.Fatal(e.Start(":" + port))
}

func hellworld(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{"message": "Hello World!"})
}
