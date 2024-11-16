package handler

import (
	"echobackend/internal/service"
	"net/http"

	"github.com/labstack/echo/v4"
)

type PostHandler struct {
	userService service.PostService
}

func NewPostHandler(postService service.PostService) *PostHandler {
	return &PostHandler{userService: postService}
}

func (h *PostHandler) GetPosts(c echo.Context) error {
	response, err := h.userService.GetPosts(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, response)
}
