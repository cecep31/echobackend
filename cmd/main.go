package main

import (
	"echobackend/internal/handler"
	"echobackend/internal/middleware"
	"echobackend/internal/repository"
	"echobackend/internal/routes"
	"echobackend/internal/service"
	"echobackend/pkg/database"
	"echobackend/server"
	"log"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalln("Error loading .env file")
	}

	// Initialize database
	db := database.SetupDatabase()
	// Initialize dependencies
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService) //

	// Initialize server
	e := server.InitServer()

	routes := routes.NewRoutes(userHandler)
	routes.Setup(e)

	api := e.Group("/api")

	api.GET("", hellworld)

	middleware.InitMiddleware(e)
	e.GET("/", hellworld)
	e.Logger.Fatal(e.Start(":1323"))
}

func hellworld(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{"message": "Hello World!"})
}
