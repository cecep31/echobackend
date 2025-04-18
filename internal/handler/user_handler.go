package handler

import (
	"echobackend/internal/service"
	"net/http"
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
		return c.JSON(http.StatusInternalServerError, map[string]any{
			"error":   err.Error(),
			"message": "Failed to retrieve user",
			"success": false,
		})
	}

	return c.JSON(http.StatusOK, map[string]any{
		"data":    userResponse,
		"message": "Successfully retrieved user",
		"success": true,
	})
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
		return c.JSON(http.StatusInternalServerError, map[string]any{
			"data":    nil,
			"error":   err.Error(),
			"success": false,
		})
	}

	return c.JSON(http.StatusOK, map[string]any{
		"messsage": "Success",
		"success":  true,
		"data":     users,
		"metadata": map[string]any{
			"totalItems": total,
		},
		"error": nil,
	})
}

// delete user
func (h *UserHandler) DeleteUser(c echo.Context) error {
	id := c.Param("id")
	err := h.userService.Delete(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]any{
			"error":   err.Error(),
			"message": "Failed to delete user",
			"success": false,
		})
	}

	return c.JSON(http.StatusOK, map[string]any{
		"message": "Successfully deleted user",
		"success": true,
	})
}
