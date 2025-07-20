package handler

import (
	"echobackend/internal/service"
	"echobackend/pkg/response"
	"strconv"

	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) GetByID(c echo.Context) error {
	userID := c.Param("id")

	userResponse, err := h.userService.GetByID(c.Request().Context(), userID)
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
