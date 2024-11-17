package handler

import (
	"echobackend/internal/service"
	"net/http"

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
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error":   err.Error(),
			"message": "Failed to retrieve user",
			"success": false,
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"data":    userResponse,
		"success": true,
	})
}

func (h *UserHandler) GetUsers(c echo.Context) error {
	users, err := h.userService.GetUsers(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error":   err.Error(),
			"success": false,
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"data":    users,
		"success": true,
	})
}
