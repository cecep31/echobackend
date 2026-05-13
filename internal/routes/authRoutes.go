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
		auth.GET("/email/:email", r.authHandler.CheckEmail)
		auth.POST("/forgot-password", r.authHandler.ForgotPassword, echomidleware.RateLimiter(echomidleware.NewRateLimiterMemoryStoreWithConfig(confratelimit)))
		auth.POST("/reset-password", r.authHandler.ResetPassword)
		auth.POST("/refresh", r.authHandler.RefreshToken)
		auth.POST("/logout", r.authHandler.Logout, r.authMiddleware.Auth())
		auth.GET("/profile", r.authHandler.GetProfile, r.authMiddleware.Auth())
		auth.PATCH("/password", r.authHandler.ChangePassword, r.authMiddleware.Auth())
		auth.GET("/activity-logs", r.authHandler.GetActivityLogs, r.authMiddleware.Auth())
		auth.GET("/activity-logs/recent", r.authHandler.GetRecentActivity, r.authMiddleware.Auth())
		auth.GET("/activity-logs/failed-logins", r.authHandler.GetFailedLogins, r.authMiddleware.Auth())
		auth.GET("/oauth/github", r.authHandler.GithubOAuthRedirect)
		auth.GET("/oauth/github/callback", r.authHandler.GithubOAuthCallback)
	}
}