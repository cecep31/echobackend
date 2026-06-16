package middleware

import (
	"log/slog"
	"strconv"
	"sync"
	"time"

	"echobackend/internal/platform/cache"
	"echobackend/pkg/response"

	"github.com/labstack/echo/v5"
)

type fixedWindowVisitor struct {
	count      int
	windowEnds time.Time
}

type fixedWindowStore struct {
	mu       sync.Mutex
	visitors map[string]fixedWindowVisitor
}

// FixedWindowRateLimiter limits each client IP to maxRequests within window.
// It is intended for low-volume abuse protection on sensitive routes.
func FixedWindowRateLimiter(maxRequests int, window time.Duration) echo.MiddlewareFunc {
	return FixedWindowRateLimiterWithCache(nil, "", maxRequests, window)
}

// FixedWindowRateLimiterWithCache uses Redis when available so limits are
// shared across app instances. If cache is nil or errors, it falls back to memory.
func FixedWindowRateLimiterWithCache(redisCache *cache.RedisCache, name string, maxRequests int, window time.Duration) echo.MiddlewareFunc {
	store := &fixedWindowStore{visitors: make(map[string]fixedWindowVisitor)}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			if maxRequests <= 0 || window <= 0 {
				return next(c)
			}

			now := time.Now()
			identifier := c.RealIP()
			if identifier == "" {
				identifier = c.Request().RemoteAddr
			}

			allowed, retryAfter := allowFixedWindow(c, redisCache, store, name, identifier, maxRequests, window, now)
			if !allowed {
				seconds := max(int(retryAfter.Seconds()), 1)
				c.Response().Header().Set("Retry-After", strconv.Itoa(seconds))
				return response.TooManyRequests(c, "Terlalu banyak percobaan. Coba lagi nanti.")
			}

			return next(c)
		}
	}
}

func allowFixedWindow(
	c *echo.Context,
	redisCache *cache.RedisCache,
	store *fixedWindowStore,
	name string,
	identifier string,
	maxRequests int,
	window time.Duration,
	now time.Time,
) (bool, time.Duration) {
	if redisCache == nil {
		return store.allow(identifier, maxRequests, window, now)
	}

	cacheKey := redisCache.BuildKey("rate_limit", name, identifier)
	count, retryAfter, err := redisCache.IncrementFixedWindow(c.Request().Context(), cacheKey, window)
	if err != nil {
		slog.Warn("rate limit: falling back to in-memory store", "name", name, "error", err)
		return store.allow(identifier, maxRequests, window, now)
	}
	if count == 0 {
		return store.allow(identifier, maxRequests, window, now)
	}

	if count > maxRequests {
		return false, retryAfter
	}
	return true, 0
}

func (s *fixedWindowStore) allow(identifier string, maxRequests int, window time.Duration, now time.Time) (bool, time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()

	visitor, ok := s.visitors[identifier]
	if !ok || !now.Before(visitor.windowEnds) {
		s.cleanup(now)
		s.visitors[identifier] = fixedWindowVisitor{
			count:      1,
			windowEnds: now.Add(window),
		}
		return true, 0
	}

	if visitor.count >= maxRequests {
		return false, visitor.windowEnds.Sub(now)
	}

	visitor.count++
	s.visitors[identifier] = visitor
	return true, 0
}

func (s *fixedWindowStore) cleanup(now time.Time) {
	for identifier, visitor := range s.visitors {
		if !now.Before(visitor.windowEnds) {
			delete(s.visitors, identifier)
		}
	}
}
