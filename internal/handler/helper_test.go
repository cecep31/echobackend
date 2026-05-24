package handler

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v5"
)

func newCtx(t *testing.T, rawQuery string) *echo.Context {
	t.Helper()
	e := echo.New()
	target := "/?" + rawQuery
	if rawQuery == "" {
		target = "/"
	}
	req := httptest.NewRequest(http.MethodGet, target, nil)
	if rawQuery != "" {
		// httptest.NewRequest already parses the query, but ensure normalized form
		req.URL.RawQuery = url.Values(req.URL.Query()).Encode()
	}
	rec := httptest.NewRecorder()
	return echo.NewContext(req, rec, e)
}

func TestGetUserIDFromClaims_NoUser(t *testing.T) {
	c := newCtx(t, "")
	if id, ok := GetUserIDFromClaims(c); ok || id != "" {
		t.Fatalf("expected empty/false, got (%q,%v)", id, ok)
	}
}

func TestGetUserIDFromClaims_MapClaims(t *testing.T) {
	c := newCtx(t, "")
	c.Set("user", jwt.MapClaims{"user_id": "user-123"})
	id, ok := GetUserIDFromClaims(c)
	if !ok || id != "user-123" {
		t.Fatalf("got (%q,%v)", id, ok)
	}
}

func TestGetUserIDFromClaims_MapClaims_MissingKey(t *testing.T) {
	c := newCtx(t, "")
	c.Set("user", jwt.MapClaims{"sub": "user-123"})
	if _, ok := GetUserIDFromClaims(c); ok {
		t.Fatal("expected missing key to return false")
	}
}

func TestGetUserIDFromClaims_MapClaims_NonString(t *testing.T) {
	c := newCtx(t, "")
	c.Set("user", jwt.MapClaims{"user_id": 42})
	if _, ok := GetUserIDFromClaims(c); ok {
		t.Fatal("expected non-string user_id to return false")
	}
}

func TestGetUserIDFromClaims_JWTToken(t *testing.T) {
	c := newCtx(t, "")
	tok := &jwt.Token{Claims: jwt.MapClaims{"user_id": "user-456"}}
	c.Set("user", tok)
	id, ok := GetUserIDFromClaims(c)
	if !ok || id != "user-456" {
		t.Fatalf("got (%q,%v)", id, ok)
	}
}

func TestGetUserIDFromClaims_JWTToken_BadClaims(t *testing.T) {
	c := newCtx(t, "")
	tok := &jwt.Token{Claims: jwt.RegisteredClaims{}}
	c.Set("user", tok)
	if _, ok := GetUserIDFromClaims(c); ok {
		t.Fatal("expected non-MapClaims to return false")
	}
}

func TestGetUserIDFromClaims_PlainMap(t *testing.T) {
	c := newCtx(t, "")
	c.Set("user", map[string]interface{}{"user_id": "user-789"})
	id, ok := GetUserIDFromClaims(c)
	if !ok || id != "user-789" {
		t.Fatalf("got (%q,%v)", id, ok)
	}
}

func TestGetUserIDFromClaims_UnknownType(t *testing.T) {
	c := newCtx(t, "")
	c.Set("user", "just-a-string")
	if _, ok := GetUserIDFromClaims(c); ok {
		t.Fatal("expected unknown type to return false")
	}
}

func TestParsePaginationParams(t *testing.T) {
	tests := []struct {
		name       string
		query      string
		def        int
		wantLimit  int
		wantOffset int
	}{
		{"defaults", "", 25, 25, 0},
		{"valid both", "limit=50&offset=10", 25, 50, 10},
		{"limit cap at 100", "limit=500", 25, 100, 0},
		{"invalid limit -> default", "limit=abc", 25, 25, 0},
		{"non-positive limit -> default", "limit=0", 25, 25, 0},
		{"negative offset ignored", "offset=-5", 25, 25, 0},
		{"non-numeric offset ignored", "offset=foo", 25, 25, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := newCtx(t, tt.query)
			limit, offset := ParsePaginationParams(c, tt.def)
			if limit != tt.wantLimit || offset != tt.wantOffset {
				t.Fatalf("got (%d,%d), want (%d,%d)", limit, offset, tt.wantLimit, tt.wantOffset)
			}
		})
	}
}
