package middleware

import (
	"log"
	"time"

	"echobackend/config"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

func InitMiddleware(e *echo.Echo, config *config.Config) {

	// Middleware to set start time for latency measurement
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			c.Set("start", time.Now())
			return next(c)
		}
	})

	// Add body limit middleware to prevent memory exhaustion
	e.Use(middleware.BodyLimit(10 * 1024 * 1024)) // Limit request body to 10MB

	// Add security headers
	e.Use(middleware.SecureWithConfig(middleware.SecureConfig{
		XSSProtection:         "1; mode=block",
		ContentTypeNosniff:    "nosniff",
		XFrameOptions:         "SAMEORIGIN",
		HSTSMaxAge:            3600,
		ContentSecurityPolicy: "default-src 'self'",
		ReferrerPolicy:        "strict-origin-when-cross-origin",
	}))

	// Enhanced request logging with structured format
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:    true,
		LogStatus: true,
		LogMethod: true,
		LogValuesFunc: func(c *echo.Context, values middleware.RequestLoggerValues) error {
			start := c.Get("start").(time.Time)
			latency := time.Since(start)
			log.Printf("handled request method=%s uri=%s request_id=%s status=%d latency=%.3f ms remote_ip=%s", values.Method, values.URI, c.Response().Header().Get(echo.HeaderXRequestID), values.Status, float64(latency.Nanoseconds())/1000000, c.RealIP())

			return nil
		},
	}))

	// Global HTTP rate limit (sustained RPS, token bucket; 0 = disabled)
	if config.HTTPRateLimitRPS > 0 {
		storeCfg := middleware.RateLimiterMemoryStoreConfig{
			Rate:  float64(config.HTTPRateLimitRPS),
			Burst: config.HTTPRateLimitRPS * 2,
		}
		if config.HTTPRateLimitWindowSec > 0 {
			storeCfg.ExpiresIn = time.Duration(config.HTTPRateLimitWindowSec) * time.Second
		}
		e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStoreWithConfig(storeCfg)))
	}

	e.Use(middleware.RequestID())
	e.Use(middleware.Gzip())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{AllowOrigins: []string{"*"}}))
}
