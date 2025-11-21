package handler

import (
	"echobackend/internal/model"
	"echobackend/internal/service"
	"echobackend/pkg/response"
	"echobackend/pkg/validator"
	"net/http"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type PostHandler struct {
	postService     service.PostService
	postViewService service.PostViewService
}

func NewPostHandler(postService service.PostService, postViewService service.PostViewService) *PostHandler {
	return &PostHandler{
		postService:     postService,
		postViewService: postViewService,
	}
}

func (h *PostHandler) GetPosts(c echo.Context) error {
	// Parse query parameters into filter struct
	filter := &model.PostQueryFilter{
		Limit:     10, // Default limit
		Offset:    0,  // Default offset
		Search:    c.QueryParam("search"),
		SortBy:    c.QueryParam("sort_by"),
		SortOrder: c.QueryParam("sort_order"),
		StartDate: c.QueryParam("start_date"),
		EndDate:   c.QueryParam("end_date"),
		CreatedBy: c.QueryParam("created_by"),
	}

	// Parse limit and offset
	if limit := c.QueryParam("limit"); limit != "" {
		if limitInt, err := strconv.Atoi(limit); err == nil && limitInt > 0 {
			filter.Limit = limitInt
		}
	}

	if offset := c.QueryParam("offset"); offset != "" {
		if offsetInt, err := strconv.Atoi(offset); err == nil && offsetInt > 0 {
			filter.Offset = offsetInt
		}
	}

	// Parse published filter
	if published := c.QueryParam("published"); published != "" {
		if pubBool, err := strconv.ParseBool(published); err == nil {
			filter.Published = &pubBool
		}
	}

	// Parse tags filter
	if tags := c.QueryParam("tags"); tags != "" {
		filter.Tags = strings.Split(tags, ",")
		// Trim whitespace from each tag
		for i, tag := range filter.Tags {
			filter.Tags[i] = strings.TrimSpace(tag)
		}
	}

	posts, total, err := h.postService.GetPostsFiltered(c.Request().Context(), filter)
	if err != nil {
		return response.InternalServerError(c, "Failed to get posts", err)
	}

	for _, post := range posts {
		if post.Body != nil && len(*post.Body) > 250 {
			truncated := (*post.Body)[:250] + " ..."
			post.Body = &truncated
		}
	}

	meta := response.PaginationMeta{
		TotalItems: int(total),
		Offset:     filter.Offset,
		Limit:      filter.Limit,
		TotalPages: int(total)/filter.Limit + 1,
	}
	if int(total)%filter.Limit == 0 {
		meta.TotalPages = int(total) / filter.Limit
	}

	return response.SuccessWithMeta(c, "Successfully retrieved posts", posts, meta)
}

func (h *PostHandler) CreatePost(c echo.Context) error {
	var postReq model.CreatePostDTO
	if err := c.Bind(&postReq); err != nil {
		return response.BadRequest(c, "Failed to create post", err)
	}

	if err := c.Validate(postReq); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			return c.JSON(http.StatusBadRequest, response.APIResponse{
				Success: false,
				Message: "Validation failed",
				Error:   validationErrors.Error(),
				Data:    validationErrors.Errors,
			})
		}
		return response.ValidationError(c, "Validation failed", err)
	}

	claims := c.Get("user").(jwt.MapClaims)
	userID := (claims)["user_id"].(string)

	newpost, err := h.postService.CreatePost(c.Request().Context(), &postReq, userID)

	if err != nil {
		return response.InternalServerError(c, "Failed to create post", err)
	}
	return response.Created(c, "Successfully created post", map[string]any{
		"id": newpost.ID,
	})
}

func (h *PostHandler) UpdatePost(c echo.Context) error {
	id := c.Param("id")
	var updateDTO model.UpdatePostDTO
	if err := c.Bind(&updateDTO); err != nil {
		return response.BadRequest(c, "Failed to update post", err)
	}

	if err := c.Validate(updateDTO); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			return c.JSON(http.StatusBadRequest, response.APIResponse{
				Success: false,
				Message: "Validation failed",
				Error:   validationErrors.Error(),
				Data:    validationErrors.Errors,
			})
		}
		return response.ValidationError(c, "Validation failed", err)
	}

	// Get the user ID from the JWT token
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userID := claims["user_id"].(string)

	// Check if the user is the author of the post
	err := h.postService.IsAuthor(c.Request().Context(), id, userID)
	if err != nil {
		return response.InternalServerError(c, "Failed to check post ownership", err)
	}

	updatedPost, err := h.postService.UpdatePost(c.Request().Context(), id, &updateDTO)
	if err != nil {
		return response.InternalServerError(c, "Failed to update post", err)
	}

	return response.Success(c, "Post updated successfully", updatedPost)
}

func (h *PostHandler) GetPostBySlugAndUsername(c echo.Context) error {
	slug := c.Param("slug")
	username := c.Param("username")
	post, err := h.postService.GetPostBySlugAndUsername(c.Request().Context(), slug, username)
	if err != nil {
		return response.InternalServerError(c, "Failed to get post", err)
	}

	return response.Success(c, "Successfully retrieved post", post)
}

func (h *PostHandler) GetPost(c echo.Context) error {
	id := c.Param("id")
	post, err := h.postService.GetPostByID(c.Request().Context(), id)
	if err != nil {
		return response.InternalServerError(c, "Failed to get post", err)
	}

	return response.Success(c, "Successfully retrieved post", post)
}

func (h *PostHandler) DeletePost(c echo.Context) error {
	id := c.Param("id")
	err := h.postService.DeletePostByID(c.Request().Context(), id)
	if err != nil {
		return response.InternalServerError(c, "Failed to delete post", err)
	}

	return response.Success(c, "Successfully deleted post", nil)
}

func (h *PostHandler) GetPostsRandom(c echo.Context) error {
	limit := c.QueryParam("limit") // Default limit if not provided or invalid
	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		limitInt = 9 // Default limit if not provided or invalid
	}
	// Ensure limit doesn't exceed 20
	if limitInt > 20 {
		limitInt = 20 // Limit to 20
	}
	posts, err := h.postService.GetPostsRandom(c.Request().Context(), limitInt)
	if err != nil {
		return response.InternalServerError(c, "Failed to get posts", err)
	}

	for _, post := range posts {
		if post.Body != nil && len(*post.Body) > 250 {
			truncated := (*post.Body)[:250] + " ..."
			post.Body = &truncated
		}
	}

	return response.Success(c, "Successfully retrieved posts", posts)
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
		if post.Body != nil && len(*post.Body) > 250 {
			truncated := (*post.Body)[:250] + " ..."
			post.Body = &truncated
		}
	}

	if err != nil {
		return response.InternalServerError(c, "Failed to get posts", err)
	}

	meta := response.PaginationMeta{
		TotalItems: int(total),
		Offset:     offsetInt,
		Limit:      limitInt,
		TotalPages: int(total)/limitInt + 1,
	}
	if int(total)%limitInt == 0 {
		meta.TotalPages = int(total) / limitInt
	}

	return response.SuccessWithMeta(c, "success retrieving posts", posts, meta)
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
		if post.Body != nil && len(*post.Body) > 250 {
			truncated := (*post.Body)[:250] + " ..."
			post.Body = &truncated
		}
	}

	if err != nil {
		return response.InternalServerError(c, "Failed to get posts", err)
	}

	meta := response.PaginationMeta{
		TotalItems: int(total),
		Offset:     offsetInt,
		Limit:      limitInt,
		TotalPages: int(total)/limitInt + 1,
	}
	if int(total)%limitInt == 0 {
		meta.TotalPages = int(total) / limitInt
	}

	return response.SuccessWithMeta(c, "success retrieving posts", posts, meta)
}

func (h *PostHandler) GetPostsByTag(c echo.Context) error {
	tag := c.Param("tag")
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
	posts, total, err := h.postService.GetPostsByTag(c.Request().Context(), tag, limitInt, offsetInt)

	for _, post := range posts {
		if post.Body != nil && len(*post.Body) > 250 {
			truncated := (*post.Body)[:250] + " ..."
			post.Body = &truncated
		}
	}

	if err != nil {
		return response.InternalServerError(c, "Failed to get posts by tag", err)
	}

	meta := response.PaginationMeta{
		TotalItems: int(total),
		Offset:     offsetInt,
		Limit:      limitInt,
		TotalPages: int(total)/limitInt + 1,
	}
	if int(total)%limitInt == 0 {
		meta.TotalPages = int(total) / limitInt
	}

	return response.SuccessWithMeta(c, "success retrieving posts by tag", posts, meta)
}

func (h *PostHandler) UploadImagePosts(c echo.Context) error {
	file, err := c.FormFile("image")
	if err != nil {
		return response.BadRequest(c, "Failed to upload image", err)
	}

	if file == nil {
		return response.BadRequest(c, "No file uploaded", nil)
	}

	if err := h.postService.UploadImagePosts(c.Request().Context(), file); err != nil {
		return response.InternalServerError(c, "Failed to upload image", err)
	}
	return response.Success(c, "success uploading image", nil)
}
