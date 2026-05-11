package routes

import (
	"time"

	"github.com/labstack/echo/v5"
	echomidleware "github.com/labstack/echo/v5/middleware"
)

func (r *Routes) setupAuthRoutes(api *echo.Group) {
	auth := api.Group("/auth")
	confratelimit := echomidleware.RateLimiterMemoryStoreConfig{Rate: 5, ExpiresIn: 5 * time.Minute, Burst: 5}
	{
		auth.POST("/register", r.authHandler.Register)
		auth.POST("/login", r.authHandler.Login, echomidleware.RateLimiter(echomidleware.NewRateLimiterMemoryStoreWithConfig(confratelimit)))
		auth.POST("/check-username", r.authHandler.CheckUsername)
		auth.POST("/forgot-password", r.authHandler.ForgotPassword, echomidleware.RateLimiter(echomidleware.NewRateLimiterMemoryStoreWithConfig(confratelimit)))
		auth.POST("/reset-password", r.authHandler.ResetPassword)
		auth.POST("/refresh", r.authHandler.RefreshToken)
		auth.PATCH("/password", r.authHandler.ChangePassword, r.authMiddleware.Auth())
	}
}
