package middleware

import (
	"echobackend/config"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

// AuthMiddleware provides authentication middleware for Echo
type AuthMiddleware struct {
	conf *config.Config
}

// NewAuthMiddleware creates a new instance of AuthMiddleware
func NewAuthMiddleware(conf *config.Config) *AuthMiddleware {
	return &AuthMiddleware{conf: conf}
}

// Auth validates JWT tokens and sets user claims in the context
func (a *AuthMiddleware) Auth() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "missing authorization header")
			}

			tokenString, err := extractBearerToken(authHeader)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
			}

			claims, err := validateToken(tokenString, a.conf.JWT_SECRET)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, fmt.Sprintf("invalid token: %v", err))
			}

			c.Set("user", claims)
			return next(c)
		}
	}
}

// AuthAdmin validates that the user has admin privileges
func (a *AuthMiddleware) AuthAdmin() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userClaims := c.Get("user")
			if userClaims == nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized: missing user context")
			}

			claims, ok := userClaims.(jwt.MapClaims)
			if !ok {
				return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized: invalid user context")
			}

			if !isAdmin(claims) {
				return echo.NewHTTPError(http.StatusForbidden, "forbidden: insufficient privileges")
			}

			return next(c)
		}
	}
}

// extractBearerToken extracts the token from the Authorization header
func extractBearerToken(authHeader string) (string, error) {
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", fmt.Errorf("invalid token format, expected 'Bearer <token>'")
	}
	return parts[1], nil
}

// validateToken validates the JWT token and returns the claims
func validateToken(tokenString, jwtSecret string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtSecret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("token parsing failed: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}

// isAdmin checks if the user has admin privileges
func isAdmin(claims jwt.MapClaims) bool {
	isSuperAdmin, exists := claims["is_super_admin"]
	if !exists {
		return false
	}

	// Check for both string and boolean representations
	switch v := isSuperAdmin.(type) {
	case string:
		return v == "true"
	case bool:
		return v
	default:
		return false
	}
}
