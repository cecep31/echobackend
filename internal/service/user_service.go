package service

import (
	"context"
	"echobackend/internal/model"
	"echobackend/internal/repository"
)

type UserService interface {
	// Create(ctx context.Context, user *domain.User) (*domain.UserResponse, error)
	GetByID(ctx context.Context, id string) (*model.UserResponse, error)
	GetUsers(ctx context.Context) ([]*model.UserResponse, error)
	// Update(ctx context.Context, user *domain.User) (*domain.UserResponse, error)
	// Delete(ctx context.Context, id uint) error
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

	return &model.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

func (s *userService) GetUsers(ctx context.Context) ([]*model.UserResponse, error) {
	users, err := s.userRepo.GetUsers(ctx)
	if err != nil {
		return nil, err
	}

	var userResponses []*model.UserResponse
	for _, user := range users {
		userResponses = append(userResponses, &model.UserResponse{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		})
	}

	return userResponses, nil
}
