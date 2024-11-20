package service

import (
	"echobackend/internal/domain"
	"echobackend/internal/repository"
)

type PostService interface {
	GetPosts(limit int, offset int) ([]*domain.Post, int64, error)
	GetPostsRandom(limit int) ([]*domain.Post, error)
}

type postService struct {
	postRepo repository.PostRepository
}

func NewPostService(postRepo repository.PostRepository) PostService {
	return &postService{postRepo: postRepo}
}

func (s *postService) GetPosts(limit int, offset int) ([]*domain.Post, int64, error) {
	total, err := s.postRepo.GetTotalPosts()
	if err != nil {
		return nil, 0, err
	}
	posts, err := s.postRepo.GetPosts(limit, offset)
	if err != nil {
		return nil, 0, err
	}
	return posts, total, err
}

func (s *postService) GetPostsRandom(limit int) ([]*domain.Post, error) {
	return s.postRepo.GetPostsRandom(limit)
}
