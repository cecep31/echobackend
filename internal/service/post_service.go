package service

import (
	"context"
	"echobackend/internal/model"
	"echobackend/internal/repository"
	"echobackend/internal/storage"
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
	postRepo     repository.PostRepository
	tagService   TagService
	miniostorage *storage.MinioStorage
}

func NewPostService(postRepo repository.PostRepository, tagService TagService, storageclient *storage.MinioStorage) PostService {
	return &postService{postRepo: postRepo, tagService: tagService, miniostorage: storageclient}
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
	posts, total, err := s.postRepo.GetPostByUsername(ctx, username, offset, limit)
	if err != nil {
		return nil, 0, err
	}

	var postsResponse []*model.PostResponse

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
	posts, total, err := s.postRepo.GetPosts(ctx, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	var postsResponse []*model.PostResponse

	for _, post := range posts {
		postResponse := post.ToResponse()
		postsResponse = append(postsResponse, postResponse)
	}

	return postsResponse, total, nil
}

func (s *postService) GetPostsRandom(ctx context.Context, limit int) ([]*model.PostResponse, error) {
	posts, err := s.postRepo.GetPostsRandom(ctx, limit)
	if err != nil {
		return nil, err
	}

	var postsResponse []*model.PostResponse

	for _, post := range posts {
		postResponse := post.ToResponse()
		postsResponse = append(postsResponse, postResponse)
	}

	return postsResponse, nil
}

func (s *postService) GetPostsByCreatedBy(ctx context.Context, createdBy string, offset int, limit int) ([]*model.PostResponse, int64, error) {
	posts, total, err := s.postRepo.GetPostsByCreatedBy(ctx, createdBy, offset, limit)
	if err != nil {
		return nil, 0, err
	}

	var postsResponse []*model.PostResponse
	for _, post := range posts {
		postsResponse = append(postsResponse, post.ToResponse())
	}

	return postsResponse, total, nil
}

func (s *postService) GetPostsByTag(ctx context.Context, tag string, limit int, offset int) ([]*model.PostResponse, int64, error) {
	posts, total, err := s.postRepo.GetPostsByTag(ctx, tag, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	var postsResponse []*model.PostResponse
	for _, post := range posts {
		postsResponse = append(postsResponse, post.ToResponse())
	}

	return postsResponse, total, nil
}

func (s *postService) UploadImagePosts(ctx context.Context, file *multipart.FileHeader) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	return s.miniostorage.Save(context.Background(), file.Filename, src)
}

// findOrCreateTagByName finds an existing tag by name or creates a new one
func (s *postService) findOrCreateTagByName(ctx context.Context, tagName string) (*model.Tag, error) {
	return s.tagService.FindOrCreateByName(ctx, tagName)
}
