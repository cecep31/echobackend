package routes

import (
	"time"

	"github.com/labstack/echo/v4"
	echomidleware "github.com/labstack/echo/v4/middleware"
	"golang.org/x/time/rate"
)

func (r *Routes) setupAuthRoutes(v1 *echo.Group) {
	auth := v1.Group("/auth")
	confratelimit := echomidleware.RateLimiterMemoryStoreConfig{Rate: rate.Limit(5), ExpiresIn: 5 * time.Minute, Burst: 5}
	{
		auth.POST("/register", r.authHandler.Register)
		auth.POST("/login", r.authHandler.Login, echomidleware.RateLimiter(echomidleware.NewRateLimiterMemoryStoreWithConfig(confratelimit)))
		auth.POST("/check-username", r.authHandler.CheckUsername)
	}
}
