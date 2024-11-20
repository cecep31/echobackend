package handler

import (
	"echobackend/internal/service"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type PostHandler struct {
	userService service.PostService
}

func NewPostHandler(postService service.PostService) *PostHandler {
	return &PostHandler{userService: postService}
}

func (h *PostHandler) GetPosts(c echo.Context) error {
	limit := c.QueryParam("limit")
	offset := c.QueryParam("offset")
	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		limitInt = 10 // Default limit if not provided or invalid
	}

	offsetInt, err := strconv.Atoi(offset)
	if err != nil {
		offsetInt = 0 // Default offset if not provided or invalid
	}
	posts, total, err := h.userService.GetPosts(limitInt, offsetInt)
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
		"metadata": echo.Map{
			"totalItems": total,
			"limit":      limitInt,
			"offset":     offsetInt,
		},
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
