package handler

import (
	"strconv"
	"time"

	"echobackend/internal/service"
	"echobackend/pkg/response"

	"github.com/labstack/echo/v5"
)

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
//   - month  1-12  (optional, default: current month)
//   - year        (optional, default: current year)
//
// Returns dividend and RUPS events for the given month, backed by Postgres.
func (h *CorporateActionHandler) GetCalendar(c *echo.Context) error {
	userID, ok := GetUserIDFromClaims(c)
	if !ok {
		return response.Unauthorized(c, "User not authenticated")
	}

	now := time.Now()
	month := int(now.Month())
	year := now.Year()

	if m := c.QueryParam("month"); m != "" {
		if v, err := strconv.Atoi(m); err == nil && v >= 1 && v <= 12 {
			month = v
		}
	}
	if y := c.QueryParam("year"); y != "" {
		if v, err := strconv.Atoi(y); err == nil {
			year = v
		}
	}

	result, err := h.service.GetCalendar(c.Request().Context(), userID, year, month)
	if err != nil {
		return response.InternalServerError(c, "Failed to fetch corporate actions calendar", err)
	}

	return response.Success(c, "Corporate actions calendar fetched successfully", result)
}
