package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"echobackend/internal/dto"
	"echobackend/internal/model"
	"echobackend/internal/service"
	"echobackend/pkg/validator"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v5"
)

var _ service.HoldingService = (*mockHoldingService)(nil)

type mockHoldingService struct {
	getMonthlyDataFn func(ctx context.Context, userID string, q *dto.HoldingMonthlyQuery) ([]dto.HoldingMonthlyDataResponse, error)
}

func (m *mockHoldingService) GetHoldings(ctx context.Context, userID string, filter *dto.HoldingQueryFilter) ([]model.Holding, error) {
	return nil, nil
}

func (m *mockHoldingService) GetHoldingByID(ctx context.Context, id int64, userID string) (*model.Holding, error) {
	return nil, nil
}

func (m *mockHoldingService) CreateHolding(ctx context.Context, userID string, req *dto.CreateHoldingRequest) (*model.Holding, error) {
	return nil, nil
}

func (m *mockHoldingService) UpdateHolding(ctx context.Context, id int64, userID string, req *dto.UpdateHoldingRequest) (*model.Holding, error) {
	return nil, nil
}

func (m *mockHoldingService) DeleteHolding(ctx context.Context, id int64, userID string) error {
	return nil
}

func (m *mockHoldingService) GetHoldingTypes(ctx context.Context) ([]model.HoldingType, error) {
	return nil, nil
}

func (m *mockHoldingService) GetSummary(ctx context.Context, userID string, q *dto.HoldingSummaryQuery) (*dto.HoldingSummaryResponse, error) {
	return nil, nil
}

func (m *mockHoldingService) GetTrends(ctx context.Context, userID string, q *dto.HoldingTrendsQuery) ([]dto.HoldingTrendResponse, error) {
	return nil, nil
}

func (m *mockHoldingService) CompareMonths(ctx context.Context, userID string, q *dto.HoldingCompareQuery) (*dto.HoldingMonthComparisonResponse, error) {
	return nil, nil
}

func (m *mockHoldingService) GetMonthlyData(ctx context.Context, userID string, q *dto.HoldingMonthlyQuery) ([]dto.HoldingMonthlyDataResponse, error) {
	if m.getMonthlyDataFn != nil {
		return m.getMonthlyDataFn(ctx, userID, q)
	}
	return nil, nil
}

func (m *mockHoldingService) SyncPrices(ctx context.Context, userID string) (*dto.HoldingSyncResponse, error) {
	return nil, nil
}

func (m *mockHoldingService) DuplicateHoldings(ctx context.Context, userID string, req *dto.DuplicateHoldingRequest) ([]dto.DuplicateResultItem, error) {
	return nil, nil
}

func newHoldingTestContext(t *testing.T, method, target string) (*echo.Context, *httptest.ResponseRecorder) {
	t.Helper()

	e := echo.New()
	e.Validator = validator.NewValidator()

	req := httptest.NewRequest(method, target, nil)
	rec := httptest.NewRecorder()
	c := echo.NewContext(req, rec, e)
	c.Set("user", jwt.MapClaims{"user_id": "user-1"})
	return c, rec
}

func TestHoldingHandler_GetMonthlyData_NormalizesNaturalRange(t *testing.T) {
	var captured *dto.HoldingMonthlyQuery
	h := NewHoldingHandler(&mockHoldingService{
		getMonthlyDataFn: func(ctx context.Context, userID string, q *dto.HoldingMonthlyQuery) ([]dto.HoldingMonthlyDataResponse, error) {
			captured = q
			return []dto.HoldingMonthlyDataResponse{}, nil
		},
	})

	c, rec := newHoldingTestContext(t, http.MethodGet, "/api/holdings/monthly?startMonth=1&startYear=2025&endMonth=12&endYear=2025")
	if err := h.GetMonthlyData(c); err != nil {
		t.Fatalf("GetMonthlyData returned error: %v", err)
	}
	if rec.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", rec.Code, rec.Body.String())
	}
	if captured == nil {
		t.Fatal("expected service to be called with a query")
	}
	if captured.StartMonth != 12 || captured.StartYear != 2025 {
		t.Fatalf("expected Start to be normalized to Dec 2025, got %d/%d", captured.StartMonth, captured.StartYear)
	}
	if captured.EndMonth != 1 || captured.EndYear != 2025 {
		t.Fatalf("expected End to be normalized to Jan 2025, got %d/%d", captured.EndMonth, captured.EndYear)
	}
}

func TestHoldingHandler_GetMonthlyData_DefaultsEndTo11MonthsBeforeStart(t *testing.T) {
	var captured *dto.HoldingMonthlyQuery
	h := NewHoldingHandler(&mockHoldingService{
		getMonthlyDataFn: func(ctx context.Context, userID string, q *dto.HoldingMonthlyQuery) ([]dto.HoldingMonthlyDataResponse, error) {
			captured = q
			return []dto.HoldingMonthlyDataResponse{}, nil
		},
	})

	c, rec := newHoldingTestContext(t, http.MethodGet, "/api/holdings/monthly?startMonth=5&startYear=2026")
	if err := h.GetMonthlyData(c); err != nil {
		t.Fatalf("GetMonthlyData returned error: %v", err)
	}
	if rec.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", rec.Code, rec.Body.String())
	}
	if captured == nil {
		t.Fatal("expected service to be called with a query")
	}
	if captured.StartMonth != 5 || captured.StartYear != 2026 {
		t.Fatalf("expected Start to be May 2026, got %d/%d", captured.StartMonth, captured.StartYear)
	}
	if captured.EndMonth != 6 || captured.EndYear != 2025 {
		t.Fatalf("expected End to be 11 months before start (Jun 2025), got %d/%d", captured.EndMonth, captured.EndYear)
	}
}

func TestHoldingHandler_GetMonthlyData_RequiresAuth(t *testing.T) {
	h := NewHoldingHandler(&mockHoldingService{})

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/holdings/monthly", nil)
	rec := httptest.NewRecorder()
	c := echo.NewContext(req, rec, e)

	if err := h.GetMonthlyData(c); err != nil {
		t.Fatalf("GetMonthlyData returned error: %v", err)
	}
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", rec.Code)
	}
}

func runParseMonthlyQuery(target string) (*dto.HoldingMonthlyQuery, error) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, target, nil)
	rec := httptest.NewRecorder()
	c := echo.NewContext(req, rec, e)
	return parseMonthlyQuery(c)
}

func TestParseMonthlyQuery_DefaultsToCurrentDate(t *testing.T) {
	q, err := runParseMonthlyQuery("/api/holdings/monthly")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	now := time.Now()
	if q.StartMonth != int(now.Month()) || q.StartYear != now.Year() {
		t.Fatalf("expected start to default to current date, got %+v", q)
	}
	if q.EndYear > q.StartYear || (q.EndYear == q.StartYear && q.EndMonth > q.StartMonth) {
		t.Fatalf("expected end to be on or before start, got %+v", q)
	}
}

func TestParseMonthlyQuery_DerivesPartialEndFromStart(t *testing.T) {
	q, err := runParseMonthlyQuery("/api/holdings/monthly?startMonth=5&startYear=2026&endYear=2025")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if q.StartMonth != 5 || q.StartYear != 2026 || q.EndMonth != 5 || q.EndYear != 2025 {
		t.Fatalf("unexpected query: %+v", q)
	}
}

func TestParseMonthlyQuery_NormalizesInvertedDifferentYear(t *testing.T) {
	q, err := runParseMonthlyQuery("/api/holdings/monthly?startMonth=1&startYear=2024&endMonth=12&endYear=2025")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if q.StartMonth != 12 || q.StartYear != 2025 || q.EndMonth != 1 || q.EndYear != 2024 {
		t.Fatalf("unexpected query: %+v", q)
	}
}
