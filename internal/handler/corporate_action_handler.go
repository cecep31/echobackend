package handler

import (
	"time"

	"echobackend/internal/service"
	"echobackend/pkg/response"

	"github.com/labstack/echo/v5"
)

const calendarHandlerDateFormat = "2006-01-02"

// CorporateActionHandler serves the GET /api/holdings/calendar endpoint.
type CorporateActionHandler struct {
	service service.CorporateActionService
}

// NewCorporateActionHandler constructs the handler.
func NewCorporateActionHandler(svc service.CorporateActionService) *CorporateActionHandler {
	return &CorporateActionHandler{service: svc}
}

// GetCalendar godoc
//
//	GET /api/holdings/calendar
//
// Query params:
//   - from  YYYY-MM-DD  (optional, default: first day of current month)
//   - to    YYYY-MM-DD  (optional, default: 3 months from now)
//
// Returns dividend and RUPS events for the authenticated user's stock holdings.
// Results are cached server-side for 6 hours.
func (h *CorporateActionHandler) GetCalendar(c *echo.Context) error {
	userID, ok := GetUserIDFromClaims(c)
	if !ok {
		return response.Unauthorized(c, "User not authenticated")
	}

	now := time.Now()

	// Default: from = start of this month
	from := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	// Default: to = 3 months from today
	to := now.AddDate(0, 3, 0)

	if f := c.QueryParam("from"); f != "" {
		if t, err := time.Parse(calendarHandlerDateFormat, f); err == nil {
			from = t
		}
	}
	if t := c.QueryParam("to"); t != "" {
		if parsed, err := time.Parse(calendarHandlerDateFormat, t); err == nil {
			to = parsed
		}
	}

	result, err := h.service.GetCalendar(c.Request().Context(), userID, from, to)
	if err != nil {
		return response.InternalServerError(c, "Failed to fetch corporate actions calendar", err)
	}

	return response.Success(c, "Corporate actions calendar fetched successfully", result)
}
