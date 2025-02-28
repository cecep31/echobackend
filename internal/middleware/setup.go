package middleware

import (
	"os"

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

	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(rate.Limit(config.RATE_LIMITER_MAX))))
	e.Use(middleware.RequestID())
	e.Use(middleware.Gzip())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{AllowOrigins: []string{"*"}}))
}
