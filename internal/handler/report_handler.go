package handler

import (
	"strconv"

	"echobackend/internal/dto"
	"echobackend/internal/service"
	"echobackend/pkg/response"

	"github.com/labstack/echo/v5"
)

type ReportHandler struct {
	reportService service.ReportService
}

func NewReportHandler(reportService service.ReportService) *ReportHandler {
	return &ReportHandler{reportService: reportService}
}

func (h *ReportHandler) GetOverview(c *echo.Context) error {
	query := dto.DateRangeQuery{
		StartDate: c.QueryParam("startDate"),
		EndDate:   c.QueryParam("endDate"),
	}
	overview, err := h.reportService.GetOverviewStats(c.Request().Context())
	if err != nil {
		return response.InternalServerError(c, "Failed to fetch overview report", err)
	}
	engagement, err := h.reportService.GetEngagementMetrics(c.Request().Context(), query)
	if err != nil {
		return response.InternalServerError(c, "Failed to fetch engagement metrics", err)
	}
	return response.Success(c, "Overview report fetched successfully", map[string]any{
		"overview":   overview,
		"engagement": engagement,
	})
}

func (h *ReportHandler) GetUsers(c *echo.Context) error {
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	report, err := h.reportService.GetUserReport(c.Request().Context(), dto.DateRangeQuery{
		StartDate: c.QueryParam("startDate"),
		EndDate:   c.QueryParam("endDate"),
	}, limit)
	if err != nil {
		return response.InternalServerError(c, "Failed to fetch user report", err)
	}
	return response.Success(c, "User report fetched successfully", report)
}

func (h *ReportHandler) GetPosts(c *echo.Context) error {
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	var tagID *int
	if raw := c.QueryParam("tagId"); raw != "" {
		if parsed, err := strconv.Atoi(raw); err == nil {
			tagID = &parsed
		}
	}
	report, err := h.reportService.GetPostReport(c.Request().Context(), dto.DateRangeQuery{
		StartDate: c.QueryParam("startDate"),
		EndDate:   c.QueryParam("endDate"),
	}, limit, tagID)
	if err != nil {
		return response.InternalServerError(c, "Failed to fetch post report", err)
	}
	return response.Success(c, "Post report fetched successfully", report)
}

func (h *ReportHandler) GetEngagement(c *echo.Context) error {
	metrics, err := h.reportService.GetEngagementMetrics(c.Request().Context(), dto.DateRangeQuery{
		StartDate: c.QueryParam("startDate"),
		EndDate:   c.QueryParam("endDate"),
	})
	if err != nil {
		return response.InternalServerError(c, "Failed to fetch engagement metrics", err)
	}
	return response.Success(c, "Engagement metrics fetched successfully", metrics)
}
