package service

import (
	"context"
	"fmt"
	"mime/multipart"

	"echobackend/internal/dto"
	apperrors "echobackend/internal/errors"
	"echobackend/internal/model"
	"echobackend/internal/repository"
	"echobackend/pkg/cache"
	"echobackend/pkg/storage"
)

type PostService interface {
	GetPosts(ctx context.Context, limit int, offset int) ([]*dto.PostResponse, int64, error)
	GetPostsFiltered(ctx context.Context, filter *dto.PostQueryFilter) ([]*dto.PostResponse, int64, error)
	GetPostsByUsername(ctx context.Context, username string, offset int, limit int) ([]*dto.PostResponse, int64, error)
	GetPostsRandom(ctx context.Context, limit int) ([]*dto.PostResponse, error)
	GetPostsTrending(ctx context.Context, offset int, limit int) ([]*dto.PostResponse, int64, error)
	GetPostByID(ctx context.Context, id string) (*dto.PostResponse, error)
	GetPostBySlugAndUsername(ctx context.Context, slug string, username string) (*dto.PostResponse, error)
	GetPostsByCreatedBy(ctx context.Context, createdBy string, offset int, limit int) ([]*dto.PostResponse, int64, error)
	GetPostsByTag(ctx context.Context, tag string, limit int, offset int) ([]*dto.PostResponse, int64, error)
	GetPostsForYou(ctx context.Context, userID string, offset int, limit int) ([]*dto.PostResponse, int64, error)
	DeletePostByID(ctx context.Context, id string) error
	UploadImagePosts(ctx context.Context, file *multipart.FileHeader) error
	CreatePost(ctx context.Context, req *dto.CreatePostRequest, creatorID string) (*dto.PostResponse, error)
	UpdatePost(ctx context.Context, id string, req *dto.UpdatePostRequest) (*dto.PostResponse, error)
	IsAuthor(ctx context.Context, id string, userid string) error
	GetPostsForSitemap(ctx context.Context, limit int) ([]*dto.SitemapPost, error)
}

type postService struct {
	postRepo   repository.PostRepository
	tagService TagService
	s3storage  *storage.S3Storage
	cache      *cache.ValkeyCache
}

func NewPostService(postRepo repository.PostRepository, tagService TagService, storageclient *storage.S3Storage, valkeyCache *cache.ValkeyCache) PostService {
	return &postService{postRepo: postRepo, tagService: tagService, s3storage: storageclient, cache: valkeyCache}
}

func (s *postService) IsAuthor(ctx context.Context, id string, userid string) error {
	post, err := s.postRepo.GetPostByID(ctx, id)
	if err != nil {
		return err
	}
	if post.CreatedBy == nil || *post.CreatedBy != userid {
		return apperrors.ErrNotAuthor
	}
	return nil
}

func (s *postService) GetPostsByUsername(ctx context.Context, username string, offset int, limit int) ([]*dto.PostResponse, int64, error) {
	if limit < 0 {
		limit = 0
	}
	if offset < 0 {
		offset = 0
	}
	if username == "" {
		return []*dto.PostResponse{}, 0, nil
	}

	posts, total, err := s.postRepo.GetPostByUsername(ctx, username, offset, limit)
	if err != nil {
		return nil, 0, err
	}

	postsResponse := make([]*dto.PostResponse, 0, len(posts))
	for _, post := range posts {
		postsResponse = append(postsResponse, dto.PostToResponse(post))
	}

	return postsResponse, total, nil
}

func (s *postService) CreatePost(ctx context.Context, req *dto.CreatePostRequest, creatorID string) (*dto.PostResponse, error) {
	var tags []model.Tag
	if len(req.Tags) > 0 {
		for _, tagName := range req.Tags {
			if tagName == "" {
				continue
			}

			tag, err := s.findOrCreateTagByName(ctx, tagName)
			if err != nil {
				return nil, err
			}
			tags = append(tags, *tag)
		}
	}

	post := &model.Post{
		Title:     &req.Title,
		Slug:      &req.Slug,
		Body:      &req.Body,
		CreatedBy: &creatorID,
		Photo_url: &req.PhotoURL,
		Published: &req.Published,
	}

	created, err := s.postRepo.CreatePostWithTags(ctx, post, tags)
	if err != nil {
		return nil, err
	}

	return dto.PostToResponse(created), nil
}

func (s *postService) GetPostBySlugAndUsername(ctx context.Context, slug string, username string) (*dto.PostResponse, error) {
	post, err := s.postRepo.GetPostBySlugAndUsername(ctx, slug, username)
	if err != nil {
		return nil, err
	}

	return dto.PostToResponse(post), nil
}

func (s *postService) DeletePostByID(ctx context.Context, id string) error {
	return s.postRepo.DeletePostByID(ctx, id)
}

func (s *postService) UpdatePost(ctx context.Context, id string, req *dto.UpdatePostRequest) (*dto.PostResponse, error) {
	updates := make(map[string]interface{})
	if req.Title != "" {
		updates["title"] = req.Title
	}
	if req.Body != "" {
		updates["body"] = req.Body
	}
	if req.Slug != "" {
		updates["slug"] = req.Slug
	}
	if req.PhotoURL != "" {
		updates["photo_url"] = req.PhotoURL
	}
	if req.Published != nil {
		updates["published"] = *req.Published
	}

	if len(updates) == 0 && len(req.Tags) == 0 {
		post, err := s.postRepo.GetPostByID(ctx, id)
		if err != nil {
			return nil, err
		}
		return dto.PostToResponse(post), nil
	}

	updatedPost, err := s.postRepo.UpdatePost(ctx, id, updates)
	if err != nil {
		return nil, err
	}
	return dto.PostToResponse(updatedPost), nil
}

func (s *postService) GetPostByID(ctx context.Context, id string) (*dto.PostResponse, error) {
	post, err := s.postRepo.GetPostByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return dto.PostToResponse(post), nil
}

func (s *postService) GetPosts(ctx context.Context, limit int, offset int) ([]*dto.PostResponse, int64, error) {
	if limit < 0 {
		limit = 0
	}
	if offset < 0 {
		offset = 0
	}

	posts, total, err := s.postRepo.GetPosts(ctx, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	postsResponse := make([]*dto.PostResponse, 0, len(posts))
	for _, post := range posts {
		postResponse := dto.PostToResponse(post)
		postsResponse = append(postsResponse, postResponse)
	}

	return postsResponse, total, nil
}

func (s *postService) GetPostsRandom(ctx context.Context, limit int) ([]*dto.PostResponse, error) {
	if limit < 0 {
		limit = 0
	}

	cacheKey := ""
	if s.cache != nil {
		cacheKey = s.cache.BuildKey("posts", "random", fmt.Sprintf("limit:%d", limit))
		var cachedPosts []*dto.PostResponse
		found, err := s.cache.GetJSON(ctx, cacheKey, &cachedPosts)
		if err == nil && found {
			return cachedPosts, nil
		}
	}

	posts, err := s.postRepo.GetPostsRandom(ctx, limit)
	if err != nil {
		return nil, err
	}

	postsResponse := make([]*dto.PostResponse, 0, len(posts))
	for _, post := range posts {
		postResponse := dto.PostToResponse(post)
		postsResponse = append(postsResponse, postResponse)
	}

	if cacheKey != "" {
		_ = s.cache.SetJSON(ctx, cacheKey, postsResponse)
	}

	return postsResponse, nil
}

func (s *postService) GetPostsTrending(ctx context.Context, offset int, limit int) ([]*dto.PostResponse, int64, error) {
	if limit < 0 {
		limit = 10
	}
	if limit == 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	posts, total, err := s.postRepo.GetPostsTrending(ctx, offset, limit)
	if err != nil {
		return nil, 0, err
	}

	postsResponse := make([]*dto.PostResponse, 0, len(posts))
	for _, post := range posts {
		postsResponse = append(postsResponse, dto.PostToResponse(post))
	}

	return postsResponse, total, nil
}

func (s *postService) GetPostsByCreatedBy(ctx context.Context, createdBy string, offset int, limit int) ([]*dto.PostResponse, int64, error) {
	if limit < 0 {
		limit = 0
	}
	if offset < 0 {
		offset = 0
	}
	if createdBy == "" {
		return []*dto.PostResponse{}, 0, nil
	}

	posts, total, err := s.postRepo.GetPostsByCreatedBy(ctx, createdBy, offset, limit)
	if err != nil {
		return nil, 0, err
	}

	postsResponse := make([]*dto.PostResponse, 0, len(posts))
	for _, post := range posts {
		postsResponse = append(postsResponse, dto.PostToResponse(post))
	}

	return postsResponse, total, nil
}

func (s *postService) GetPostsByTag(ctx context.Context, tag string, limit int, offset int) ([]*dto.PostResponse, int64, error) {
	if limit < 0 {
		limit = 0
	}
	if offset < 0 {
		offset = 0
	}
	if tag == "" {
		return []*dto.PostResponse{}, 0, nil
	}

	posts, total, err := s.postRepo.GetPostsByTag(ctx, tag, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	postsResponse := make([]*dto.PostResponse, 0, len(posts))
	for _, post := range posts {
		postsResponse = append(postsResponse, dto.PostToResponse(post))
	}

	return postsResponse, total, nil
}

func (s *postService) GetPostsFiltered(ctx context.Context, filter *dto.PostQueryFilter) ([]*dto.PostResponse, int64, error) {
	if filter.Limit < 0 {
		filter.Limit = 10
	}
	if filter.Limit > 100 {
		filter.Limit = 100
	}
	if filter.Offset < 0 {
		filter.Offset = 0
	}

	posts, total, err := s.postRepo.GetPostsFiltered(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	postsResponse := make([]*dto.PostResponse, 0, len(posts))
	for _, post := range posts {
		postResponse := dto.PostToResponse(post)
		postsResponse = append(postsResponse, postResponse)
	}

	return postsResponse, total, nil
}

func (s *postService) GetPostsForYou(ctx context.Context, userID string, offset int, limit int) ([]*dto.PostResponse, int64, error) {
	if limit < 0 {
		limit = 10
	}
	if limit == 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}
	if userID == "" {
		return []*dto.PostResponse{}, 0, nil
	}

	posts, total, err := s.postRepo.GetPostsForYou(ctx, userID, offset, limit)
	if err != nil {
		return nil, 0, err
	}

	postsResponse := make([]*dto.PostResponse, 0, len(posts))
	for _, post := range posts {
		postsResponse = append(postsResponse, dto.PostToResponse(post))
	}

	return postsResponse, total, nil
}

func (s *postService) UploadImagePosts(ctx context.Context, file *multipart.FileHeader) error {
	if file == nil {
		return apperrors.ErrFileNil
	}

	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	return s.s3storage.Save(ctx, file.Filename, src)
}

func (s *postService) findOrCreateTagByName(ctx context.Context, tagName string) (*model.Tag, error) {
	return s.tagService.FindOrCreateByName(ctx, tagName)
}

func (s *postService) GetPostsForSitemap(ctx context.Context, limit int) ([]*dto.SitemapPost, error) {
	if limit < 0 {
		limit = 0
	}
	return s.postRepo.GetPostsForSitemap(ctx, limit)
}
