package handler

import (
	"echobackend/internal/service"
	"echobackend/pkg/response"
	"strconv"

	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	userService       service.UserService
	userFollowService service.UserFollowService
}

func NewUserHandler(userService service.UserService, userFollowService service.UserFollowService) *UserHandler {
	return &UserHandler{
		userService:       userService,
		userFollowService: userFollowService,
	}
}

func (h *UserHandler) GetByID(c echo.Context) error {
	userID := c.Param("id")

	// Get current user ID from JWT if authenticated
	var currentUserID string
	if userClaims := c.Get("user"); userClaims != nil {
		if claims, ok := userClaims.(map[string]interface{}); ok {
			if uid, exists := claims["user_id"]; exists {
				if uidStr, ok := uid.(string); ok {
					currentUserID = uidStr
				}
			}
		}
	}

	// Get user with follow status
	userResponse, err := h.userFollowService.GetUserWithFollowStatus(c.Request().Context(), userID, currentUserID)
	if err != nil {
		return response.InternalServerError(c, "Failed to retrieve user", err)
	}

	return response.Success(c, "Successfully retrieved user", userResponse)
}

func (h *UserHandler) GetUsers(c echo.Context) error {
	offset, err := strconv.Atoi(c.QueryParam("offset"))
	if err != nil {
		offset = 0
	}
	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil {
		limit = 10
	}
	users, total, err := h.userService.GetUsers(c.Request().Context(), offset, limit)
	if err != nil {
		return response.InternalServerError(c, "Failed to retrieve users", err)
	}

	meta := response.CalculatePaginationMeta(total, offset, limit)
	return response.SuccessWithMeta(c, "Successfully retrieved users", users, meta)
}

// delete user
func (h *UserHandler) DeleteUser(c echo.Context) error {
	id := c.Param("id")
	err := h.userService.Delete(c.Request().Context(), id)
	if err != nil {
		return response.InternalServerError(c, "Failed to delete user", err)
	}

	return response.Success(c, "Successfully deleted user", nil)
}

// GetMe returns the current authenticated user's information
func (h *UserHandler) GetMe(c echo.Context) error {
	userClaims := c.Get("user")
	if userClaims == nil {
		return response.InternalServerError(c, "User context not found", nil)
	}

	claims, ok := userClaims.(map[string]interface{})
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

	userResponse, err := h.userService.GetByID(c.Request().Context(), userIDStr)
	if err != nil {
		return response.InternalServerError(c, "Failed to retrieve user", err)
	}

	return response.Success(c, "Successfully retrieved current user", userResponse)
}
