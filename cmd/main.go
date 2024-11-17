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
	conf, err := config.Load()
	if err != nil {
		panic(err)
	}

	// Initialize database
	databaseConnection, err := database.SetupDatabase(conf)
	if err != nil {
		panic(err)
	}

	// Initialize repositories
	userRepository := repository.NewUserRepository(databaseConnection)
	postRepository := repository.NewPostRepository(databaseConnection)

	// Initialize services
	userService := service.NewUserService(userRepository)
	postService := service.NewPostService(postRepository)

	// Initialize handlers
	userHandler := handler.NewUserHandler(userService)
	postHandler := handler.NewPostHandler(postService)

	// Initialize middleware
	authMiddleware := middleware.NewAuthMiddleware(conf)

	// Initialize server
	e := echo.New()
	routes := routes.NewRoutes(userHandler, postHandler, authMiddleware)
	routes.Setup(e)
	middleware.InitMiddleware(e)

	// Define routes
	e.GET("/", hellworld)

	// Start server
	port := conf.GetAppPort()
	e.Logger.Fatal(e.Start(":" + port))
}

func hellworld(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{"message": "Hello World!"})
}
