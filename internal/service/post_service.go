package service

import (
	"echobackend/internal/model"
	"echobackend/internal/repository"
)

type PostService interface {
	GetPosts(limit int, offset int) ([]*model.Post, int64, error)
	GetPostsRandom(limit int) ([]*model.Post, error)
}

type postService struct {
	postRepo repository.PostRepository
}

func NewPostService(postRepo repository.PostRepository) PostService {
	return &postService{postRepo: postRepo}
}

func (s *postService) GetPosts(limit int, offset int) ([]*model.Post, int64, error) {
	var posts []*model.Post
	var total int64
	var err error

	total, err = s.postRepo.GetTotalPosts()
	if err != nil {
		return nil, 0, err
	}

	posts, err = s.postRepo.GetPosts(limit, offset)
	if err != nil {
		return nil, 0, err
	}

	return posts, total, nil
}

func (s *postService) GetPostsRandom(limit int) ([]*model.Post, error) {
	return s.postRepo.GetPostsRandom(limit)
}
