package middleware

import (
	"echobackend/internal/config"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

type AuthMiddleware struct {
	conf *config.Config
}

func NewAuthMiddleware(conf *config.Config) *AuthMiddleware {
	return &AuthMiddleware{conf: conf}
}

func (a *AuthMiddleware) ExampleMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		return next(c)
	}
}

func (a *AuthMiddleware) Auth() echo.MiddlewareFunc {
	return echojwt.WithConfig(echojwt.Config{
		SigningMethod: "HS256",
		SigningKey:    []byte(a.conf.GetJWTSecret()),
	})
}
