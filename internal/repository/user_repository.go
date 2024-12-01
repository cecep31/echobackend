package repository

import (
	"context"
	"errors"
	"fmt"

	"echobackend/internal/model"

	"gorm.io/gorm"
)

// Errors that can be returned by the repository
var (
	ErrUserNotFound = errors.New("user not found")
	ErrUserExists   = errors.New("user already exists")
)

type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	GetByID(ctx context.Context, id string) (*model.User, error)
	GetUsers(ctx context.Context, offset int, limit int) ([]*model.User, int64, error)
	GetUsersByEmail(ctx context.Context, email string) ([]*model.User, error)
	Update(ctx context.Context, user *model.User) error
	SoftDeleteByID(ctx context.Context, id string) error
	Exists(ctx context.Context, email string) (bool, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *model.User) error {
	exists, err := r.Exists(ctx, user.Email)
	if err != nil {
		return fmt.Errorf("failed to check user existence: %w", err)
	}
	if exists {
		return ErrUserExists
	}

	return r.db.WithContext(ctx).Create(user).Error
}

func (r *userRepository) Update(ctx context.Context, user *model.User) error {
	result := r.db.WithContext(ctx).Model(&model.User{}).Where("id = ?", user.ID).Updates(user)
	if result.Error != nil {
		return fmt.Errorf("failed to update user: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return ErrUserNotFound
	}
	return nil
}

func (r *userRepository) GetByID(ctx context.Context, id string) (*model.User, error) {
	var user model.User
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return &user, nil
}

func (r *userRepository) GetUsers(ctx context.Context, offset, limit int) ([]*model.User, int64, error) {
	var users []*model.User
	var total int64

	// Count total records
	if err := r.db.WithContext(ctx).Model(&model.User{}).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count users: %w", err)
	}

	// Get paginated records
	if err := r.db.WithContext(ctx).Offset(offset).Limit(limit).Find(&users).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to get users: %w", err)
	}

	return users, total, nil
}

func (r *userRepository) GetUsersByEmail(ctx context.Context, email string) ([]*model.User, error) {
	var users []*model.User
	if err := r.db.WithContext(ctx).Where("email = ?", email).Find(&users).Error; err != nil {
		return nil, fmt.Errorf("failed to get users by email: %w", err)
	}
	return users, nil
}

func (r *userRepository) SoftDeleteByID(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Where("id = ?", id).Delete(&model.User{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete user: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return ErrUserNotFound
	}
	return nil
}

func (r *userRepository) Exists(ctx context.Context, email string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.User{}).Where("email = ?", email).Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("failed to check user existence: %w", err)
	}
	return count > 0, nil
}
