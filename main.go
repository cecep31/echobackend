package main

import (
	"echobackend/middleware"
	"echobackend/server"
	"net/http"

	"github.com/labstack/echo/v4"
)

func main() {
	e := server.InitServer()
	middleware.InitMiddleware(e)
	e.GET("/", hellworld)
	e.Logger.Fatal(e.Start(":1323"))
}

func hellworld(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{"message": "Hello World!"})
}
