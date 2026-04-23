package middleware

import (
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

	e.Use(middleware.RequestID())

	// Enhanced request logging with structured format (Echo v5 logger is *slog.Logger)
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:    true,
		LogStatus: true,
		LogMethod: true,
		LogValuesFunc: func(c *echo.Context, values middleware.RequestLoggerValues) error {
			start := c.Get("start").(time.Time)
			latency := time.Since(start)
			e.Logger.Info("handled request",
				"method", values.Method,
				"uri", values.URI,
				"request_id", c.Response().Header().Get(echo.HeaderXRequestID),
				"status", values.Status,
				"latency_ms", float64(latency.Nanoseconds())/1e6,
				"remote_ip", c.RealIP(),
			)
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

	e.Use(middleware.Gzip())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{AllowOrigins: []string{"*"}}))
}
