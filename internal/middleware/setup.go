package middleware

import (
	"os"
	"time"

	"echobackend/config"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"golang.org/x/time/rate"
)

func InitMiddleware(e *echo.Echo, config *config.Config) {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	logger := log.Output(zerolog.ConsoleWriter{Out: os.Stdout})

	// Set server timeouts
	e.Server.ReadTimeout = 10 * time.Second
	e.Server.WriteTimeout = 15 * time.Second
	e.Server.IdleTimeout = 60 * time.Second

	// Add body limit middleware to prevent memory exhaustion
	e.Use(middleware.BodyLimit("10M")) // Limit request body to 10MB

	// Add security headers
	e.Use(middleware.SecureWithConfig(middleware.SecureConfig{
		XSSProtection:         "1; mode=block",
		ContentTypeNosniff:    "nosniff",
		XFrameOptions:         "SAMEORIGIN",
		HSTSMaxAge:            3600,
		ContentSecurityPolicy: "default-src 'self'",
		ReferrerPolicy:        "strict-origin-when-cross-origin",
	}))

	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:    true,
		LogStatus: true,
		LogMethod: true,
		LogValuesFunc: func(c echo.Context, values middleware.RequestLoggerValues) error {
			logger.Info().
				Str("method", values.Method).
				Str("uri", values.URI).
				Str("request_id", c.Response().Header().Get(echo.HeaderXRequestID)).
				Int("status", values.Status).
				Dur("latency", values.Latency).
				Str("remote_ip", c.RealIP()).
				Msg("handled request")

			return nil
		},
	}))

	// Enhanced rate limiting with custom store and configuration
	if config.RATE_LIMITER_MAX > 0 {
		e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStoreWithConfig(
			// Use the correct RateLimiterMemoryStoreConfig
			middleware.RateLimiterMemoryStoreConfig{
				Rate:  rate.Limit(config.RATE_LIMITER_MAX),
				Burst: config.RATE_LIMITER_MAX * 2,
			},
		)))
	}

	e.Use(middleware.RequestID())
	e.Use(middleware.Gzip())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{AllowOrigins: []string{"*"}}))
}
