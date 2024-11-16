package service

import (
	"context"
	"echobackend/internal/domain"
	"echobackend/internal/repository"
)

type UserService interface {
	// Create(ctx context.Context, user *domain.User) (*domain.UserResponse, error)
	GetByID(ctx context.Context, id string) (*domain.UserResponse, error)
	GetUsers(ctx context.Context) ([]*domain.UserResponse, error)
	// Update(ctx context.Context, user *domain.User) (*domain.UserResponse, error)
	// Delete(ctx context.Context, id uint) error
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

func (s *userService) GetByID(ctx context.Context, id string) (*domain.UserResponse, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &domain.UserResponse{
		Id:        user.Id,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

func (s *userService) GetUsers(ctx context.Context) ([]*domain.UserResponse, error) {
	users, err := s.userRepo.GetUsers(ctx)
	if err != nil {
		return nil, err
	}

	var userResponses []*domain.UserResponse
	for _, user := range users {
		userResponses = append(userResponses, &domain.UserResponse{
			Id:        user.Id,
			Name:      user.Name,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		})
	}

	return userResponses, nil
}
