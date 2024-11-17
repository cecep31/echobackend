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
	posts, err := h.userService.GetPosts(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error":   err.Error(),
			"message": "Failed to get posts",
			"success": false,
		})
	}

	for _, post := range posts {
		if len(post.Body) > 250 {
			post.Body = post.Body[:250] + " ..."
		}
	}

	return c.JSON(http.StatusOK, echo.Map{
		"data":    posts,
		"message": "Successfully retrieved posts",
		"success": true,
	})
}

func (h *PostHandler) GetPostsRandom(c echo.Context) error {
	const limit = 6
	posts, err := h.userService.GetPostsRandom(limit)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error":   err.Error(),
			"message": "Failed to get posts",
			"success": false,
		})
	}

	for _, post := range posts {
		if len(post.Body) > 250 {
			post.Body = post.Body[:250] + " ..."
		}
	}

	return c.JSON(http.StatusOK, echo.Map{
		"data":    posts,
		"message": "Successfully retrieved posts",
		"success": true,
	})
}
