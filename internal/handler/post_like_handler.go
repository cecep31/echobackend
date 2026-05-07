package handler

import (
	"errors"

	apperrors "echobackend/internal/errors"
	"echobackend/internal/service"
	"echobackend/pkg/response"

	"github.com/labstack/echo/v5"
)

type PostLikeHandler struct {
	postLikeService service.PostLikeService
}

func NewPostLikeHandler(postLikeService service.PostLikeService) *PostLikeHandler {
	return &PostLikeHandler{postLikeService: postLikeService}
}

func (h *PostLikeHandler) LikePost(c *echo.Context) error {
	postID := c.Param("id")
	if postID == "" {
		return response.BadRequest(c, "Post ID is required", nil)
	}

	userID, ok := GetUserIDFromClaims(c)
	if !ok {
		return response.Unauthorized(c, "User authentication required")
	}

	err := h.postLikeService.LikePost(c.Request().Context(), postID, userID)
	if err != nil {
		if errors.Is(err, apperrors.ErrAlreadyLiked) {
			return response.BadRequest(c, "You have already liked this post", nil)
		}
		return response.InternalServerError(c, "Failed to like post", err)
	}

	return response.Success(c, "Post liked successfully", nil)
}

func (h *PostLikeHandler) UnlikePost(c *echo.Context) error {
	postID := c.Param("id")
	if postID == "" {
		return response.BadRequest(c, "Post ID is required", nil)
	}

	userID, ok := GetUserIDFromClaims(c)
	if !ok {
		return response.Unauthorized(c, "User authentication required")
	}

	err := h.postLikeService.UnlikePost(c.Request().Context(), postID, userID)
	if err != nil {
		if errors.Is(err, apperrors.ErrNotLiked) {
			return response.BadRequest(c, "You have not liked this post", nil)
		}
		return response.InternalServerError(c, "Failed to unlike post", err)
	}

	return response.Success(c, "Post unliked successfully", nil)
}

func (h *PostLikeHandler) GetPostLikes(c *echo.Context) error {
	postID := c.Param("id")
	if postID == "" {
		return response.BadRequest(c, "Post ID is required", nil)
	}

	limit, offset := ParsePaginationParams(c, 10)

	likes, total, err := h.postLikeService.GetLikesByPostID(c.Request().Context(), postID, limit, offset)
	if err != nil {
		return response.InternalServerError(c, "Failed to get post likes", err)
	}

	responseData := map[string]any{
		"likes":  likes,
		"total":  total,
		"limit":  limit,
		"offset": offset,
	}

	return response.Success(c, "Post likes retrieved successfully", responseData)
}

func (h *PostLikeHandler) GetPostLikeStats(c *echo.Context) error {
	postID := c.Param("id")
	if postID == "" {
		return response.BadRequest(c, "Post ID is required", nil)
	}

	stats, err := h.postLikeService.GetLikeStats(c.Request().Context(), postID)
	if err != nil {
		return response.InternalServerError(c, "Failed to get like stats", err)
	}

	return response.Success(c, "Like stats retrieved successfully", stats)
}

func (h *PostLikeHandler) CheckUserLiked(c *echo.Context) error {
	postID := c.Param("id")
	if postID == "" {
		return response.BadRequest(c, "Post ID is required", nil)
	}

	userID, ok := GetUserIDFromClaims(c)
	if !ok {
		return response.Unauthorized(c, "User authentication required")
	}

	hasLiked, err := h.postLikeService.HasUserLikedPost(c.Request().Context(), postID, userID)
	if err != nil {
		return response.InternalServerError(c, "Failed to check like status", err)
	}

	responseData := map[string]interface{}{
		"has_liked": hasLiked,
		"post_id":   postID,
		"user_id":   userID,
	}

	return response.Success(c, "Like status retrieved successfully", responseData)
}
