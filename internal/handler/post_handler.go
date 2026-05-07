package handler

import (
	"echobackend/internal/dto"
	"echobackend/internal/service"
	"echobackend/pkg/response"
	"strconv"
	"strings"

	"github.com/labstack/echo/v5"
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

func (h *PostHandler) GetPosts(c *echo.Context) error {
	filter := &dto.PostQueryFilter{
		Limit:     10,
		Offset:    0,
		Search:    c.QueryParam("search"),
		SortBy:    c.QueryParam("sort_by"),
		SortOrder: c.QueryParam("sort_order"),
		StartDate: c.QueryParam("start_date"),
		EndDate:   c.QueryParam("end_date"),
		CreatedBy: c.QueryParam("created_by"),
	}

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

	if published := c.QueryParam("published"); published != "" {
		if pubBool, err := strconv.ParseBool(published); err == nil {
			filter.Published = &pubBool
		}
	}

	if tags := c.QueryParam("tags"); tags != "" {
		filter.Tags = strings.Split(tags, ",")
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

func (h *PostHandler) CreatePost(c *echo.Context) error {
	var postReq dto.CreatePostRequest
	if err := c.Bind(&postReq); err != nil {
		return response.BadRequest(c, "Failed to create post", err)
	}

	if err := c.Validate(postReq); err != nil {
		return response.FromValidateError(c, err)
	}

	userID, ok := GetUserIDFromClaims(c)
	if !ok {
		return response.Unauthorized(c, "User not authenticated")
	}

	newpost, err := h.postService.CreatePost(c.Request().Context(), &postReq, userID)
	if err != nil {
		return response.InternalServerError(c, "Failed to create post", err)
	}
	return response.Created(c, "Successfully created post", map[string]any{
		"id": newpost.ID,
	})
}

func (h *PostHandler) UpdatePost(c *echo.Context) error {
	id := c.Param("id")
	var updateDTO dto.UpdatePostRequest
	if err := c.Bind(&updateDTO); err != nil {
		return response.BadRequest(c, "Failed to update post", err)
	}

	if err := c.Validate(updateDTO); err != nil {
		return response.FromValidateError(c, err)
	}

	userID, ok := GetUserIDFromClaims(c)
	if !ok {
		return response.Unauthorized(c, "User not authenticated")
	}

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

func (h *PostHandler) GetPostBySlugAndUsername(c *echo.Context) error {
	slug := c.Param("slug")
	username := c.Param("username")
	post, err := h.postService.GetPostBySlugAndUsername(c.Request().Context(), slug, username)
	if err != nil {
		return response.InternalServerError(c, "Failed to get post", err)
	}

	return response.Success(c, "Successfully retrieved post", post)
}

func (h *PostHandler) GetPost(c *echo.Context) error {
	id := c.Param("id")
	post, err := h.postService.GetPostByID(c.Request().Context(), id)
	if err != nil {
		return response.InternalServerError(c, "Failed to get post", err)
	}

	return response.Success(c, "Successfully retrieved post", post)
}

func (h *PostHandler) DeletePost(c *echo.Context) error {
	id := c.Param("id")
	err := h.postService.DeletePostByID(c.Request().Context(), id)
	if err != nil {
		return response.InternalServerError(c, "Failed to delete post", err)
	}

	return response.Success(c, "Successfully deleted post", nil)
}

func (h *PostHandler) GetPostsRandom(c *echo.Context) error {
	limit, _ := ParsePaginationParams(c, 9)
	if limit > 20 {
		limit = 20
	}
	posts, err := h.postService.GetPostsRandom(c.Request().Context(), limit)
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

func (h *PostHandler) GetPostsTrending(c *echo.Context) error {
	limit, offset := ParsePaginationParams(c, 10)

	posts, total, err := h.postService.GetPostsTrending(c.Request().Context(), offset, limit)
	if err != nil {
		return response.InternalServerError(c, "Failed to get trending posts", err)
	}

	for _, post := range posts {
		if post.Body != nil && len(*post.Body) > 250 {
			truncated := (*post.Body)[:250] + " ..."
			post.Body = &truncated
		}
	}

	metaLimit := limit
	if metaLimit <= 0 {
		metaLimit = 10
	}
	if metaLimit > 100 {
		metaLimit = 100
	}

	meta := response.PaginationMeta{
		TotalItems: int(total),
		Offset:     offset,
		Limit:      metaLimit,
		TotalPages: int(total)/metaLimit + 1,
	}
	if int(total)%metaLimit == 0 {
		meta.TotalPages = int(total) / metaLimit
	}

	return response.SuccessWithMeta(c, "Successfully retrieved trending posts", posts, meta)
}

func (h *PostHandler) GetMyPosts(c *echo.Context) error {
	limit, offset := ParsePaginationParams(c, 10)

	userID, ok := GetUserIDFromClaims(c)
	if !ok {
		return response.Unauthorized(c, "User not authenticated")
	}

	posts, total, err := h.postService.GetPostsByCreatedBy(c.Request().Context(), userID, offset, limit)

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
		Offset:     offset,
		Limit:      limit,
		TotalPages: int(total)/limit + 1,
	}
	if int(total)%limit == 0 {
		meta.TotalPages = int(total) / limit
	}

	return response.SuccessWithMeta(c, "success retrieving posts", posts, meta)
}

func (h *PostHandler) GetPostsForYou(c *echo.Context) error {
	limit, offset := ParsePaginationParams(c, 10)

	userID, ok := GetUserIDFromClaims(c)
	if !ok {
		return response.Unauthorized(c, "User not authenticated")
	}

	posts, total, err := h.postService.GetPostsForYou(c.Request().Context(), userID, offset, limit)
	if err != nil {
		return response.InternalServerError(c, "Failed to get posts", err)
	}

	metaLimit := limit
	if metaLimit <= 0 {
		metaLimit = 10
	}
	if metaLimit > 100 {
		metaLimit = 100
	}

	for _, post := range posts {
		if post.Body != nil && len(*post.Body) > 250 {
			truncated := (*post.Body)[:250] + " ..."
			post.Body = &truncated
		}
	}

	meta := response.PaginationMeta{
		TotalItems: int(total),
		Offset:     offset,
		Limit:      metaLimit,
		TotalPages: int(total)/metaLimit + 1,
	}
	if int(total)%metaLimit == 0 {
		meta.TotalPages = int(total) / metaLimit
	}

	return response.SuccessWithMeta(c, "Successfully retrieved for-you posts", posts, meta)
}

func (h *PostHandler) GetPostsByUsername(c *echo.Context) error {
	username := c.Param("username")
	limit, offset := ParsePaginationParams(c, 10)

	posts, total, err := h.postService.GetPostsByUsername(c.Request().Context(), username, offset, limit)

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
		Offset:     offset,
		Limit:      limit,
		TotalPages: int(total)/limit + 1,
	}
	if int(total)%limit == 0 {
		meta.TotalPages = int(total) / limit
	}

	return response.SuccessWithMeta(c, "success retrieving posts", posts, meta)
}

func (h *PostHandler) GetPostsByAuthor(c *echo.Context) error {
	username := c.Param("username")
	limit, offset := ParsePaginationParams(c, 10)

	posts, total, err := h.postService.GetPostsByUsername(c.Request().Context(), username, offset, limit)

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
		Offset:     offset,
		Limit:      limit,
		TotalPages: int(total)/limit + 1,
	}
	if int(total)%limit == 0 {
		meta.TotalPages = int(total) / limit
	}

	return response.SuccessWithMeta(c, "success retrieving posts", posts, meta)
}

func (h *PostHandler) GetPostsByTag(c *echo.Context) error {
	tag := c.Param("tag")
	limit, offset := ParsePaginationParams(c, 10)

	posts, total, err := h.postService.GetPostsByTag(c.Request().Context(), tag, limit, offset)

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
		Offset:     offset,
		Limit:      limit,
		TotalPages: int(total)/limit + 1,
	}
	if int(total)%limit == 0 {
		meta.TotalPages = int(total) / limit
	}

	return response.SuccessWithMeta(c, "success retrieving posts by tag", posts, meta)
}

func (h *PostHandler) UploadImagePosts(c *echo.Context) error {
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

func (h *PostHandler) GetPostsForSitemap(c *echo.Context) error {
	posts, err := h.postService.GetPostsForSitemap(c.Request().Context(), 1000)
	if err != nil {
		return response.InternalServerError(c, "Failed to get posts for sitemap", err)
	}

	return response.Success(c, "Successfully retrieved posts for sitemap", posts)
}
