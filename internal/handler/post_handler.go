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
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error":   err.Error(),
			"message": "Failed to get posts",
			"success": false,
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"data":    response,
		"message": "Successfully get posts",
		"success": true,
	})
}

func (h *PostHandler) GetPostsRandom(c echo.Context) error {
	response, err := h.userService.GetPostsRandom(6)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error":   err.Error(),
			"message": "Failed to get posts",
			"success": false,
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"data":    response,
		"message": "Successfully get posts",
		"success": true,
	})
}
