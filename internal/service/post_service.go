package service

import (
	"echobackend/internal/model"
	"echobackend/internal/repository"
)

type PostService interface {
	GetPosts(limit int, offset int) ([]*model.PostResponse, int64, error)
	GetPostsRandom(limit int) ([]*model.PostResponse, error)
	GetPostByID(id string) (*model.PostResponse, error)
	DeletePostByID(id string) error
}

type postService struct {
	postRepo repository.PostRepository
}

func NewPostService(postRepo repository.PostRepository) PostService {
	return &postService{postRepo: postRepo}
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
	var total int64
	var err error

	total, err = s.postRepo.GetTotalPosts()
	if err != nil {
		return nil, 0, err
	}

	posts, err := s.postRepo.GetPosts(limit, offset)
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
