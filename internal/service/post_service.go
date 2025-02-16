package service

import (
	"context"
	"echobackend/internal/model"
	"echobackend/internal/repository"
	"echobackend/internal/storage"
	"mime/multipart"
)

type PostService interface {
	GetPosts(limit int, offset int) ([]*model.PostResponse, int64, error)
	GetPostsByUsername(username string, offset int, limit int) ([]*model.PostResponse, int64, error)
	GetPostsRandom(limit int) ([]*model.PostResponse, error)
	GetPostByID(id string) (*model.PostResponse, error)
	GetPostBySlugAndUsername(slug string, username string) (*model.PostResponse, error)
	GetPostsByCreatedBy(createdBy string, offset int, limit int) ([]*model.PostResponse, int64, error)
	DeletePostByID(id string) error
	UploadImagePosts(file *multipart.FileHeader) error
}

type postService struct {
	postRepo     repository.PostRepository
	miniostorage *storage.MinioStorage
}

func NewPostService(postRepo repository.PostRepository, storageclient *storage.MinioStorage) PostService {
	return &postService{postRepo: postRepo, miniostorage: storageclient}
}

func (s *postService) GetPostsByUsername(username string, offset int, limit int) ([]*model.PostResponse, int64, error) {
	posts, total, err := s.postRepo.GetPostByUsername(username, offset, limit)
	if err != nil {
		return nil, 0, err
	}

	var postsResponse []*model.PostResponse

	for _, post := range posts {
		postsResponse = append(postsResponse, post.ToResponse())
	}

	return postsResponse, total, nil
}

func (s *postService) GetPostBySlugAndUsername(slug string, username string) (*model.PostResponse, error) {
	post, err := s.postRepo.GetPostBySlugAndUsername(slug, username)
	if err != nil {
		return nil, err
	}

	return post.ToResponse(), nil
}

func (s *postService) DeletePostByID(id string) error {
	return s.postRepo.DeletePostByID(id)
}

func (s *postService) GetPostByID(id string) (*model.PostResponse, error) {
	post, err := s.postRepo.GetPostByID(id)
	if err != nil {
		return nil, err
	}

	return post.ToResponse(), nil
}

func (s *postService) GetPosts(limit int, offset int) ([]*model.PostResponse, int64, error) {
	posts, total, err := s.postRepo.GetPosts(limit, offset)
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

func (s *postService) GetPostsRandom(limit int) ([]*model.PostResponse, error) {
	posts, err := s.postRepo.GetPostsRandom(limit)
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

func (s *postService) GetPostsByCreatedBy(createdBy string, offset int, limit int) ([]*model.PostResponse, int64, error) {
	posts, total, err := s.postRepo.GetPostsByCreatedBy(createdBy, offset, limit)
	if err != nil {
		return nil, 0, err
	}

	var postsResponse []*model.PostResponse
	for _, post := range posts {
		postsResponse = append(postsResponse, post.ToResponse())
	}

	return postsResponse, total, nil
}

func (s *postService) UploadImagePosts(file *multipart.FileHeader) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	return s.miniostorage.Save(context.Background(), file.Filename, src)
}
