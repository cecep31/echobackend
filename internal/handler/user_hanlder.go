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
	id := c.Param("id")

	response, err := h.userService.GetByID(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error":   err.Error(),
			"message": "Failed to get user",
			"success": false,
		})
	}

	return c.JSON(http.StatusOK, response)
}

func (h *UserHandler) GetUsers(c echo.Context) error {
	response, err := h.userService.GetUsers(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error":   err.Error(),
			"message": "Failed to get users",
			"success": false,
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"data":    response,
		"message": "Successfully get users",
		"success": true,
	})
}
