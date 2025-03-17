package service

import (
	"context"
	"echobackend/internal/model"
	"echobackend/internal/repository"
)

type UserService interface {
	GetByID(ctx context.Context, id string) (*model.UserResponse, error)
	GetUsers(ctx context.Context, offset int, limit int) ([]*model.UserResponse, int64, error)
	Delete(ctx context.Context, id string) error
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

func (s *userService) GetByID(ctx context.Context, id string) (*model.UserResponse, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return user.ToResponse(), nil
}

func (s *userService) GetUsers(ctx context.Context, offset int, limit int) ([]*model.UserResponse, int64, error) {
	users, total, err := s.userRepo.GetUsers(ctx, offset, limit)
	if err != nil {
		return nil, 0, err
	}

	var userResponses []*model.UserResponse
	for _, user := range users {
		userResponses = append(userResponses, user.ToResponse())
	}

	return userResponses, total, nil
}

func (s *userService) Delete(ctx context.Context, id string) error {
	return s.userRepo.SoftDeleteByID(ctx, id)
}
