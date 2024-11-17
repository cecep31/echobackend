package service

import (
	"context"
	"echobackend/internal/domain"
	"echobackend/internal/repository"
)

type PostService interface {
	GetPosts(ctx context.Context) ([]*domain.Post, error)
	GetPostsRandom(limit int) ([]*domain.Post, error)
}

type postService struct {
	postRepo repository.PostRepository
}

func NewPostService(postRepo repository.PostRepository) PostService {
	return &postService{postRepo: postRepo}
}

func (s *postService) GetPosts(ctx context.Context) ([]*domain.Post, error) {
	return s.postRepo.GetPosts(ctx)
}

func (s *postService) GetPostsRandom(limit int) ([]*domain.Post, error) {
	return s.postRepo.GetPostsRandom(limit)
}
