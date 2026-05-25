package service

import (
	"context"
	"echobackend/internal/dto"
	"echobackend/internal/repository"
	"fmt"
)

type UserService interface {
	GetByID(ctx context.Context, id string) (*dto.UserResponse, error)
	GetAdminByID(ctx context.Context, id string, deletedOnly bool) (*dto.UserResponse, error)
	GetMe(ctx context.Context, id string) (*dto.CurrentUserResponse, error)
	GetByUsername(ctx context.Context, username string) (*dto.UserResponse, error)
	GetUsers(ctx context.Context, offset int, limit int, deletedFilter repository.UserDeletedFilter) ([]*dto.UserResponse, int64, error)
	Delete(ctx context.Context, id string) error
	Restore(ctx context.Context, id string) (*dto.UserResponse, error)
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

func (s *userService) GetByID(ctx context.Context, id string) (*dto.UserResponse, error) {
	user, err := s.userRepo.GetByID(ctx, id, false)
	if err != nil {
		return nil, err
	}
	return dto.UserToResponse(user), nil
}

func (s *userService) GetAdminByID(ctx context.Context, id string, deletedOnly bool) (*dto.UserResponse, error) {
	user, err := s.userRepo.GetByID(ctx, id, deletedOnly)
	if err != nil {
		return nil, err
	}
	return dto.UserToAdminResponse(user), nil
}

func (s *userService) GetMe(ctx context.Context, id string) (*dto.CurrentUserResponse, error) {
	user, err := s.userRepo.GetByID(ctx, id, false)
	if err != nil {
		return nil, err
	}
	return dto.UserToCurrentUserResponse(user), nil
}

func (s *userService) GetByUsername(ctx context.Context, username string) (*dto.UserResponse, error) {
	user, err := s.userRepo.GetByUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	return dto.UserToResponse(user), nil
}

func (s *userService) GetUsers(ctx context.Context, offset int, limit int, deletedFilter repository.UserDeletedFilter) ([]*dto.UserResponse, int64, error) {
	users, total, err := s.userRepo.GetUsers(ctx, offset, limit, deletedFilter)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to retrieve users from repository: %w", err)
	}

	var userResponses []*dto.UserResponse
	for _, user := range users {
		if user == nil {
			continue
		}
		userResponses = append(userResponses, dto.UserToAdminResponse(user))
	}

	return userResponses, total, nil
}

func (s *userService) Delete(ctx context.Context, id string) error {
	return s.userRepo.SoftDeleteByID(ctx, id)
}

func (s *userService) Restore(ctx context.Context, id string) (*dto.UserResponse, error) {
	if err := s.userRepo.RestoreByID(ctx, id); err != nil {
		return nil, err
	}

	user, err := s.userRepo.GetByID(ctx, id, false)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve restored user: %w", err)
	}

	return dto.UserToAdminResponse(user), nil
}
