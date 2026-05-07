package service

import (
	"context"
	"echobackend/internal/dto"
	"echobackend/internal/repository"
	"fmt"
)

type UserService interface {
	GetByID(ctx context.Context, id string) (*dto.UserResponse, error)
	GetByUsername(ctx context.Context, username string) (*dto.UserResponse, error)
	GetUsers(ctx context.Context, offset int, limit int) ([]*dto.UserResponse, int64, error)
	Delete(ctx context.Context, id string) error
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

func (s *userService) GetByID(ctx context.Context, id string) (*dto.UserResponse, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return dto.UserToResponse(user), nil
}

func (s *userService) GetByUsername(ctx context.Context, username string) (*dto.UserResponse, error) {
	user, err := s.userRepo.GetByUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	return dto.UserToResponse(user), nil
}

func (s *userService) GetUsers(ctx context.Context, offset int, limit int) ([]*dto.UserResponse, int64, error) {
	users, total, err := s.userRepo.GetUsers(ctx, offset, limit)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to retrieve users from repository: %w", err)
	}

	var userResponses []*dto.UserResponse
	for _, user := range users {
		if user == nil {
			continue
		}
		userResponses = append(userResponses, dto.UserToResponse(user))
	}

	return userResponses, total, nil
}

func (s *userService) Delete(ctx context.Context, id string) error {
	return s.userRepo.SoftDeleteByID(ctx, id)
}
