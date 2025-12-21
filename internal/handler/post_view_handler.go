package handler

import (
	"echobackend/internal/service"
	"echobackend/pkg/response"
	"strconv"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type PostViewHandler struct {
	postViewService service.PostViewService
}

func NewPostViewHandler(postViewService service.PostViewService) *PostViewHandler {
	return &PostViewHandler{postViewService: postViewService}
}

// RecordView records a view for a post
func (h *PostViewHandler) RecordView(c echo.Context) error {
	postID := c.Param("id")
	if postID == "" {
		return response.BadRequest(c, "Post ID is required", nil)
	}

	// Get user ID from JWT if authenticated
	var userID string
	if userClaims := c.Get("user"); userClaims != nil {
		if claims, ok := userClaims.(jwt.MapClaims); ok {
			if uid, exists := claims["user_id"]; exists {
				if uidStr, ok := uid.(string); ok {
					userID = uidStr
				}
			}
		}
	}

	// Extract IP address and user agent from request
	ipAddress := c.RealIP()
	userAgent := c.Request().UserAgent()

	// Record the view with additional metadata
	err := h.postViewService.RecordView(c.Request().Context(), postID, userID, &ipAddress, &userAgent)
	if err != nil {
		return response.InternalServerError(c, "Failed to record view", err)
	}

	return response.Success(c, "View recorded successfully", nil)
}

// GetPostViews gets views for a specific post
func (h *PostViewHandler) GetPostViews(c echo.Context) error {
	postID := c.Param("id")
	if postID == "" {
		return response.BadRequest(c, "Post ID is required", nil)
	}

	// Parse pagination parameters
	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil || limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100 // Max limit
	}

	offset, err := strconv.Atoi(c.QueryParam("offset"))
	if err != nil || offset < 0 {
		offset = 0
	}

	// Get views
	views, total, err := h.postViewService.GetViewsByPostID(c.Request().Context(), postID, limit, offset)
	if err != nil {
		return response.InternalServerError(c, "Failed to get post views", err)
	}

	// Calculate pagination meta
	meta := response.CalculatePaginationMeta(total, offset, limit)

	return response.SuccessWithMeta(c, "Successfully retrieved post views", views, meta)
}

// GetPostViewStats gets view statistics for a specific post
func (h *PostViewHandler) GetPostViewStats(c echo.Context) error {
	postID := c.Param("id")
	if postID == "" {
		return response.BadRequest(c, "Post ID is required", nil)
	}

	// Get view statistics
	stats, err := h.postViewService.GetViewStats(c.Request().Context(), postID)
	if err != nil {
		return response.InternalServerError(c, "Failed to get view statistics", err)
	}

	return response.Success(c, "Successfully retrieved view statistics", stats)
}

// CheckUserViewed checks if the current user has viewed a specific post
func (h *PostViewHandler) CheckUserViewed(c echo.Context) error {
	postID := c.Param("id")
	if postID == "" {
		return response.BadRequest(c, "Post ID is required", nil)
	}

	// Get user ID from JWT
	userClaims := c.Get("user")
	if userClaims == nil {
		return response.Unauthorized(c, "Authentication required")
	}

	claims, ok := userClaims.(jwt.MapClaims)
	if !ok {
		return response.InternalServerError(c, "Invalid user context", nil)
	}

	userID, exists := claims["user_id"]
	if !exists {
		return response.InternalServerError(c, "User ID not found in token", nil)
	}

	userIDStr, ok := userID.(string)
	if !ok {
		return response.InternalServerError(c, "Invalid user ID format", nil)
	}

	// Check if user has viewed the post
	hasViewed, err := h.postViewService.HasUserViewedPost(c.Request().Context(), postID, userIDStr)
	if err != nil {
		return response.InternalServerError(c, "Failed to check view status", err)
	}

	return response.Success(c, "Successfully checked view status", map[string]bool{
		"has_viewed": hasViewed,
	})
}
