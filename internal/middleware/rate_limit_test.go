package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/labstack/echo/v5"
)

func TestFixedWindowRateLimiterDeniesAfterLimit(t *testing.T) {
	e := echo.New()
	limiter := FixedWindowRateLimiter(2, time.Minute)
	handler := limiter(func(c *echo.Context) error {
		return c.NoContent(http.StatusNoContent)
	})

	for i := 0; i < 2; i++ {
		req := httptest.NewRequest(http.MethodPost, "/login", nil)
		req.RemoteAddr = "192.0.2.10:1234"
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		if err := handler(c); err != nil {
			t.Fatalf("request %d returned error: %v", i+1, err)
		}
		if rec.Code != http.StatusNoContent {
			t.Fatalf("request %d status = %d, want %d", i+1, rec.Code, http.StatusNoContent)
		}
	}

	req := httptest.NewRequest(http.MethodPost, "/login", nil)
	req.RemoteAddr = "192.0.2.10:1234"
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if err := handler(c); err != nil {
		t.Fatalf("limited request returned error: %v", err)
	}
	if rec.Code != http.StatusTooManyRequests {
		t.Fatalf("limited request status = %d, want %d", rec.Code, http.StatusTooManyRequests)
	}
	if rec.Header().Get("Retry-After") == "" {
		t.Fatal("Retry-After header is empty")
	}
}

func TestFixedWindowStoreAllowsAfterWindowExpires(t *testing.T) {
	store := &fixedWindowStore{visitors: make(map[string]fixedWindowVisitor)}
	now := time.Date(2026, 5, 30, 12, 0, 0, 0, time.UTC)

	allowed, _ := store.allow("192.0.2.10", 1, time.Minute, now)
	if !allowed {
		t.Fatal("first request should be allowed")
	}

	allowed, retryAfter := store.allow("192.0.2.10", 1, time.Minute, now.Add(30*time.Second))
	if allowed {
		t.Fatal("second request inside window should be denied")
	}
	if retryAfter != 30*time.Second {
		t.Fatalf("retryAfter = %s, want 30s", retryAfter)
	}

	allowed, _ = store.allow("192.0.2.10", 1, time.Minute, now.Add(time.Minute))
	if !allowed {
		t.Fatal("request after window expiry should be allowed")
	}
}
