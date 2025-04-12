package handler

import (
	"echobackend/internal/model"
	"echobackend/internal/service"
	"echobackend/pkg/validator"
	"net/http"
	"strconv"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type PostHandler struct {
	postService service.PostService
}

func NewPostHandler(postService service.PostService) *PostHandler {
	return &PostHandler{postService: postService}
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
	posts, total, err := h.postService.GetPosts(c.Request().Context(), limitInt, offsetInt)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]any{
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

	return c.JSON(http.StatusOK, map[string]any{
		"data":    posts,
		"message": "Successfully retrieved posts",
		"success": true,
		"metadata": map[string]any{
			"totalItems": total,
			"limit":      limitInt,
			"offset":     offsetInt,
		},
	})
}

func (h *PostHandler) CreatePost(c echo.Context) error {
	var postReq model.CreatePostDTO
	if err := c.Bind(&postReq); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{
			"success": false,
			"message": "Failed to create post",
		})
	}

	if err := c.Validate(postReq); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			return c.JSON(http.StatusBadRequest, Response{
				Success: false,
				Message: "Validation failed",
				Errors:  validationErrors.Errors,
			})
		}
		return c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: "Validation failed",
			Errors:  []string{err.Error()},
		})
	}

	claims := c.Get("user").(jwt.MapClaims)
	userID := (claims)["user_id"].(string)

	newpost, err := h.postService.CreatePost(c.Request().Context(), &postReq, userID)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]any{
			"error":   err.Error(),
			"message": "Failed to create post",
			"success": false,
		})
	}
	return c.JSON(http.StatusCreated, map[string]any{
		"data": map[string]any{
			"id": newpost.ID,
		},
		"message": "Successfully created post",
		"success": true,
	})
}

func (h *PostHandler) UpdatePost(c echo.Context) error {
	return nil
}

func (h *PostHandler) GetPostBySlugAndUsername(c echo.Context) error {
	slug := c.Param("slug")
	username := c.Param("username")
	post, err := h.postService.GetPostBySlugAndUsername(c.Request().Context(), slug, username)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]any{
			"error":   err.Error(),
			"message": "Failed to get post",
			"success": false,
		})
	}

	return c.JSON(http.StatusOK, map[string]any{
		"data":    post,
		"message": "Successfully retrieved post",
		"success": true,
	})
}

func (h *PostHandler) GetPost(c echo.Context) error {
	id := c.Param("id")
	post, err := h.postService.GetPostByID(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]any{
			"error":   err.Error(),
			"message": "Failed to get post",
			"success": false,
		})
	}

	return c.JSON(http.StatusOK, map[string]any{
		"data":    post,
		"message": "Successfully retrieved post",
		"success": true,
	})
}

func (h *PostHandler) DeletePost(c echo.Context) error {
	id := c.Param("id")
	err := h.postService.DeletePostByID(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]any{
			"error":   err.Error(),
			"message": "Failed to delete post",
			"success": false,
		})
	}

	return c.JSON(http.StatusOK, map[string]any{
		"message": "Successfully deleted post",
		"success": true,
	})
}

func (h *PostHandler) GetPostsRandom(c echo.Context) error {
	limit := c.QueryParam("limit") // Default limit if not provided or invalid
	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		limitInt = 6 // Default limit if not provided or invalid
	}
	posts, err := h.postService.GetPostsRandom(c.Request().Context(), limitInt)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]any{
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

	return c.JSON(http.StatusOK, map[string]any{
		"data":    posts,
		"message": "Successfully retrieved posts",
		"success": true,
	})
}

func (h *PostHandler) GetMyPosts(c echo.Context) error {
	offset := c.QueryParam("offset")
	limit := c.QueryParam("limit")
	offsetInt, err := strconv.Atoi(offset)
	if err != nil {
		offsetInt = 0 // Default offset if not provided or invalid
	}
	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		limitInt = 10 // Default limit if not provided or invalid
	}
	claims := c.Get("user").(jwt.MapClaims)
	userID := (claims)["user_id"].(string)
	posts, total, err := h.postService.GetPostsByCreatedBy(c.Request().Context(), userID, offsetInt, limitInt)

	for _, post := range posts {
		if len(post.Body) > 250 {
			post.Body = post.Body[:250] + " ..."
		}
	}

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]any{
			"error":   err.Error(),
			"message": "Failed to get posts",
			"success": false,
		})
	}

	return c.JSON(http.StatusOK, map[string]any{
		"success": true,
		"message": "success retrieving posts",
		"data":    posts,
		"metadata": map[string]any{
			"totalItems": total,
			"limit":      limitInt,
			"offset":     offsetInt,
		},
	})
}

func (h *PostHandler) GetPostsByUsername(c echo.Context) error {
	username := c.Param("username")
	offset := c.QueryParam("offset")
	limit := c.QueryParam("limit")
	offsetInt, err := strconv.Atoi(offset)
	if err != nil {
		offsetInt = 0 // Default offset if not provided or invalid
	}
	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		limitInt = 10 // Default limit if not provided or invalid
	}
	posts, total, err := h.postService.GetPostsByUsername(c.Request().Context(), username, offsetInt, limitInt)

	for _, post := range posts {
		if len(post.Body) > 250 {
			post.Body = post.Body[:250] + " ..."
		}
	}

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]any{
			"error":   err.Error(),
			"message": "Failed to get posts",
			"success": false,
		})
	}

	return c.JSON(http.StatusOK, map[string]any{
		"success": true,
		"message": "success retrieving posts",
		"data":    posts,
		"metadata": map[string]any{
			"totalItems": total,
			"limit":      limitInt,
			"offset":     offsetInt,
		},
	})
}

func (h *PostHandler) UploadImagePosts(c echo.Context) error {
	file, err := c.FormFile("image")
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{
			"success": false,
			"message": "Failed to upload image",
			"data":    nil,
			"error":   err.Error(),
		})
	}

	if file == nil {
		return c.JSON(http.StatusBadRequest, map[string]any{
			"success": false,
			"message": "No file uploaded",
			"data":    nil,
			"error":   nil,
		})
	}

	if err := h.postService.UploadImagePosts(c.Request().Context(), file); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]any{
			"success": false,
			"message": "Failed to upload image",
			"data":    nil,
			"error":   err.Error(),
		})
	}
	return c.JSON(http.StatusOK, map[string]any{
		"success": true,
		"message": "success uploading image",
		"data":    nil,
		"error":   nil,
	})
}
