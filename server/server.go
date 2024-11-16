package server

import "github.com/labstack/echo/v4"

func InitServer() *echo.Echo {
	e := echo.New()
	return e
}
