package init

import "github.com/labstack/echo/v4"

func InitServer() {
	e := echo.New()
	return e
}
