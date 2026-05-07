package handler

import (
	"echobackend/internal/service"
	"echobackend/pkg/response"

	"github.com/labstack/echo/v5"
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

func (h *UserHandler) GetByID(c *echo.Context) error {
	userID := c.Param("id")

	var currentUserID string
	if uid, ok := GetUserIDFromClaims(c); ok {
		currentUserID = uid
	}

	userResponse, err := h.userFollowService.GetUserWithFollowStatus(c.Request().Context(), userID, currentUserID)
	if err != nil {
		return response.InternalServerError(c, "Failed to retrieve user", err)
	}

	return response.Success(c, "Successfully retrieved user", userResponse)
}

func (h *UserHandler) GetByUsername(c *echo.Context) error {
	username := c.Param("username")

	var currentUserID string
	if uid, ok := GetUserIDFromClaims(c); ok {
		currentUserID = uid
	}

	user, err := h.userService.GetByUsername(c.Request().Context(), username)
	if err != nil {
		return response.InternalServerError(c, "Failed to retrieve user", err)
	}

	userResponse, err := h.userFollowService.GetUserWithFollowStatus(c.Request().Context(), user.ID, currentUserID)
	if err != nil {
		return response.InternalServerError(c, "Failed to retrieve user", err)
	}

	return response.Success(c, "Successfully retrieved user", userResponse)
}

func (h *UserHandler) GetUsers(c *echo.Context) error {
	limit, offset := ParsePaginationParams(c, 10)

	users, total, err := h.userService.GetUsers(c.Request().Context(), offset, limit)
	if err != nil {
		return response.InternalServerError(c, "Failed to retrieve users", err)
	}

	meta := response.CalculatePaginationMeta(total, offset, limit)
	return response.SuccessWithMeta(c, "Successfully retrieved users", users, meta)
}

func (h *UserHandler) DeleteUser(c *echo.Context) error {
	id := c.Param("id")
	err := h.userService.Delete(c.Request().Context(), id)
	if err != nil {
		return response.InternalServerError(c, "Failed to delete user", err)
	}

	return response.Success(c, "Successfully deleted user", nil)
}

func (h *UserHandler) GetMe(c *echo.Context) error {
	userID, ok := GetUserIDFromClaims(c)
	if !ok {
		return response.Unauthorized(c, "User not authenticated")
	}

	userResponse, err := h.userService.GetByID(c.Request().Context(), userID)
	if err != nil {
		return response.InternalServerError(c, "Failed to retrieve user", err)
	}

	return response.Success(c, "Successfully retrieved current user", userResponse)
}
