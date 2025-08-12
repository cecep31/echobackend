package service

import (
	"context"
	"echobackend/internal/model"
	"echobackend/internal/repository"
	"echobackend/pkg/storage"
	"errors"
	"mime/multipart"
)

type PostService interface {
	GetPosts(ctx context.Context, limit int, offset int) ([]*model.PostResponse, int64, error)
	GetPostsByUsername(ctx context.Context, username string, offset int, limit int) ([]*model.PostResponse, int64, error)
	GetPostsRandom(ctx context.Context, limit int) ([]*model.PostResponse, error)
	GetPostByID(ctx context.Context, id string) (*model.PostResponse, error)
	GetPostBySlugAndUsername(ctx context.Context, slug string, username string) (*model.PostResponse, error)
	GetPostsByCreatedBy(ctx context.Context, createdBy string, offset int, limit int) ([]*model.PostResponse, int64, error)
	GetPostsByTag(ctx context.Context, tag string, limit int, offset int) ([]*model.PostResponse, int64, error)
	DeletePostByID(ctx context.Context, id string) error
	UploadImagePosts(ctx context.Context, file *multipart.FileHeader) error
	CreatePost(ctx context.Context, post *model.CreatePostDTO, creator_id string) (*model.Post, error)
	UpdatePost(ctx context.Context, id string, post *model.UpdatePostDTO) (*model.Post, error)
	IsAuthor(ctx context.Context, id string, userid string) error
}

type postService struct {
	postRepo   repository.PostRepository
	tagService TagService
	s3storage  *storage.S3Storage
}

func NewPostService(postRepo repository.PostRepository, tagService TagService, storageclient *storage.S3Storage) PostService {
	return &postService{postRepo: postRepo, tagService: tagService, s3storage: storageclient}
}

func (s *postService) IsAuthor(ctx context.Context, id string, userid string) error {
	post, err := s.postRepo.GetPostByID(ctx, id)
	if err != nil {
		return err
	}
	if post.CreatedBy != userid {
		return errors.New("not author")
	}
	return nil
}

func (s *postService) GetPostsByUsername(ctx context.Context, username string, offset int, limit int) ([]*model.PostResponse, int64, error) {
	// Input validation
	if limit < 0 {
		limit = 0
	}
	if offset < 0 {
		offset = 0
	}
	if username == "" {
		return []*model.PostResponse{}, 0, nil
	}

	posts, total, err := s.postRepo.GetPostByUsername(ctx, username, offset, limit)
	if err != nil {
		return nil, 0, err
	}

	// Pre-allocate slice with known capacity to reduce memory allocations
	postsResponse := make([]*model.PostResponse, 0, len(posts))

	for _, post := range posts {
		postsResponse = append(postsResponse, post.ToResponse())
	}

	return postsResponse, total, nil
}

func (s *postService) CreatePost(ctx context.Context, post *model.CreatePostDTO, creator_id string) (*model.Post, error) {
	// Handle tags if they exist
	var tags []model.Tag
	if len(post.Tags) > 0 {
		for _, tagName := range post.Tags {
			if tagName == "" {
				continue // Skip empty tag names
			}

			// Try to find existing tag by name
			tag, err := s.findOrCreateTagByName(ctx, tagName)
			if err != nil {
				return nil, err
			}
			tags = append(tags, *tag)
		}
	}

	// Create the post with tags
	return s.postRepo.CreatePostWithTags(ctx, post, creator_id, tags)
}

func (s *postService) GetPostBySlugAndUsername(ctx context.Context, slug string, username string) (*model.PostResponse, error) {
	post, err := s.postRepo.GetPostBySlugAndUsername(ctx, slug, username)
	if err != nil {
		return nil, err
	}

	return post.ToResponse(), nil
}

func (s *postService) DeletePostByID(ctx context.Context, id string) error {
	return s.postRepo.DeletePostByID(ctx, id)
}

func (s *postService) UpdatePost(ctx context.Context, id string, post *model.UpdatePostDTO) (*model.Post, error) {
	return s.postRepo.UpdatePost(ctx, id, post)
}

func (s *postService) GetPostByID(ctx context.Context, id string) (*model.PostResponse, error) {
	post, err := s.postRepo.GetPostByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return post.ToResponse(), nil
}

func (s *postService) GetPosts(ctx context.Context, limit int, offset int) ([]*model.PostResponse, int64, error) {
	// Input validation
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

	// Pre-allocate slice with known capacity to reduce memory allocations
	postsResponse := make([]*model.PostResponse, 0, len(posts))

	for _, post := range posts {
		postResponse := post.ToResponse()
		postsResponse = append(postsResponse, postResponse)
	}

	return postsResponse, total, nil
}

func (s *postService) GetPostsRandom(ctx context.Context, limit int) ([]*model.PostResponse, error) {
	// Input validation
	if limit < 0 {
		limit = 0
	}

	posts, err := s.postRepo.GetPostsRandom(ctx, limit)
	if err != nil {
		return nil, err
	}

	// Pre-allocate slice with known capacity to reduce memory allocations
	postsResponse := make([]*model.PostResponse, 0, len(posts))

	for _, post := range posts {
		postResponse := post.ToResponse()
		postsResponse = append(postsResponse, postResponse)
	}

	return postsResponse, nil
}

func (s *postService) GetPostsByCreatedBy(ctx context.Context, createdBy string, offset int, limit int) ([]*model.PostResponse, int64, error) {
	// Input validation
	if limit < 0 {
		limit = 0
	}
	if offset < 0 {
		offset = 0
	}
	if createdBy == "" {
		return []*model.PostResponse{}, 0, nil
	}

	posts, total, err := s.postRepo.GetPostsByCreatedBy(ctx, createdBy, offset, limit)
	if err != nil {
		return nil, 0, err
	}

	// Pre-allocate slice with known capacity to reduce memory allocations
	postsResponse := make([]*model.PostResponse, 0, len(posts))
	for _, post := range posts {
		postsResponse = append(postsResponse, post.ToResponse())
	}

	return postsResponse, total, nil
}

func (s *postService) GetPostsByTag(ctx context.Context, tag string, limit int, offset int) ([]*model.PostResponse, int64, error) {
	// Input validation
	if limit < 0 {
		limit = 0
	}
	if offset < 0 {
		offset = 0
	}
	if tag == "" {
		return []*model.PostResponse{}, 0, nil
	}

	posts, total, err := s.postRepo.GetPostsByTag(ctx, tag, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	// Pre-allocate slice with known capacity to reduce memory allocations
	postsResponse := make([]*model.PostResponse, 0, len(posts))
	for _, post := range posts {
		postsResponse = append(postsResponse, post.ToResponse())
	}

	return postsResponse, total, nil
}

func (s *postService) UploadImagePosts(ctx context.Context, file *multipart.FileHeader) error {
	// Input validation
	if file == nil {
		return errors.New("file cannot be nil")
	}

	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	// Use the passed context instead of context.Background() to respect cancellation/timeout
	return s.s3storage.Save(ctx, file.Filename, src)
}

// findOrCreateTagByName finds an existing tag by name or creates a new one
func (s *postService) findOrCreateTagByName(ctx context.Context, tagName string) (*model.Tag, error) {
	return s.tagService.FindOrCreateByName(ctx, tagName)
}
