package middleware

import (
	"echobackend/internal/config"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

type AuthMiddleware struct {
	cfg *config.Config
}

func NewAuthMiddleware(cfg *config.Config) *AuthMiddleware {
	return &AuthMiddleware{cfg: cfg}
}

func (a *AuthMiddleware) Middleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		return next(c)
	}
}

func (a *AuthMiddleware) Auth() echo.MiddlewareFunc {
	return echojwt.WithConfig(echojwt.Config{
		SigningMethod: "HS256",
		SigningKey:    []byte(a.cfg.GetJWTSecret()),
	})
}
