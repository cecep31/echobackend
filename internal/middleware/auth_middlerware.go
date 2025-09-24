package middleware

import (
	"echobackend/config"
	"fmt"
	"strings"

	"github.com/golang-jwt/jwt/v5"
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
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return echo.NewHTTPError(401, "missing authorization header")
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || parts[0] != "Bearer" {
				return echo.NewHTTPError(401, "invalid token format")
			}
			tokenString := parts[1]

			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				return []byte(a.conf.JWT_SECRET), nil
			})

			if err != nil {
				return echo.NewHTTPError(401, "invalid or expired token")
			}

			if !token.Valid {
				return echo.NewHTTPError(401, "invalid token")
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				return echo.NewHTTPError(401, "invalid token claims")
			}

			c.Set("user", claims)
			return next(c)
		}
	}
}

func (a *AuthMiddleware) AuthAdmin() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userClaims := c.Get("user")
			if userClaims == nil {
				return echo.NewHTTPError(401, "unauthorized: missing user context")
			}

			claims, ok := userClaims.(jwt.MapClaims)
			if !ok {
				return echo.NewHTTPError(401, "unauthorized: invalid user context")
			}

			isSuperAdmin, exists := claims["is_super_admin"]
			if !exists {
				return echo.NewHTTPError(403, "forbidden: missing admin privileges")
			}

			if isSuperAdmin != "true" && isSuperAdmin != true {
				return echo.NewHTTPError(403, "forbidden: insufficient privileges")
			}

			return next(c)
		}
	}
}
