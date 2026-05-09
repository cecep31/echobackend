package repository

import (
	"context"
	"fmt"

	apperrors "echobackend/internal/errors"
	"echobackend/internal/model"

	"gorm.io/gorm"
)

type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	GetByID(ctx context.Context, id string) (*model.User, error)
	GetUsers(ctx context.Context, offset int, limit int) ([]*model.User, int64, error)
	GetUsersByEmail(ctx context.Context, email string) ([]*model.User, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	GetByUsername(ctx context.Context, username string) (*model.User, error)
	Update(ctx context.Context, user *model.User) error
	SoftDeleteByID(ctx context.Context, id string) error
	Exists(ctx context.Context, email string) (bool, error)
	CheckUserByUsername(ctx context.Context, username string) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) CheckUserByUsername(ctx context.Context, username string) error {
	var exists bool
	err := r.db.WithContext(ctx).Model(&model.User{}).
		Select("1").
		Where("username = ?", username).
		Limit(1).
		Scan(&exists).Error
	if err != nil {
		return err
	}
	if exists {
		return apperrors.ErrUserExists
	}
	return nil
}

func (r *userRepository) Create(ctx context.Context, user *model.User) error {
	exists, err := r.Exists(ctx, user.Email)
	if err != nil {
		return fmt.Errorf("failed to check user existence: %w", err)
	}
	if exists {
		return apperrors.ErrUserExists
	}

	result := r.db.WithContext(ctx).Create(user)
	return result.Error
}

func (r *userRepository) Update(ctx context.Context, user *model.User) error {
	result := r.db.WithContext(ctx).Model(user).
		Select("Email", "FirstName", "LastName", "Username", "IsSuperAdmin", "Password").
		Where("id = ?", user.ID).
		Updates(user)

	if result.Error != nil {
		return fmt.Errorf("failed to update user: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return apperrors.ErrUserNotFound
	}
	return nil
}

func (r *userRepository) GetByID(ctx context.Context, id string) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).Preload("Profile").Where("id = ?", id).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, apperrors.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return &user, nil
}

func (r *userRepository) GetUsers(ctx context.Context, offset, limit int) ([]*model.User, int64, error) {
	var users []*model.User
	var totalCount int64

	if offset < 0 {
		return nil, 0, fmt.Errorf("offset cannot be negative")
	}
	if limit <= 0 || limit > 100 {
		return nil, 0, fmt.Errorf("limit must be between 1 and 100")
	}

	err := r.db.WithContext(ctx).Model((*model.User)(nil)).Count(&totalCount).Error
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count users: %w", err)
	}

	err = r.db.WithContext(ctx).Preload("Profile").Offset(offset).Limit(limit).Find(&users).Error
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get users: %w", err)
	}

	return users, totalCount, nil
}

func (r *userRepository) GetUsersByEmail(ctx context.Context, email string) ([]*model.User, error) {
	var users []*model.User
	err := r.db.WithContext(ctx).Where("email = ?", email).Find(&users).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get users by email: %w", err)
	}
	return users, nil
}

func (r *userRepository) SoftDeleteByID(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Where("id = ?", id).Delete(&model.User{})
	if result.Error != nil {
		return fmt.Errorf("failed to soft delete user: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return apperrors.ErrUserNotFound
	}
	return nil
}

func (r *userRepository) Exists(ctx context.Context, email string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model((*model.User)(nil)).Where("email = ?", email).Count(&count).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, fmt.Errorf("failed to check user existence: %w", err)
	}
	return count > 0, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, apperrors.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}
	return &user, nil
}

func (r *userRepository) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).Preload("Profile").Where("username = ?", username).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, apperrors.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user by username: %w", err)
	}
	return &user, nil
}
