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

	result := r.db.WithContext(ctx).Create(user)
	return result.Error
}

func (r *userRepository) Update(ctx context.Context, user *model.User) error {
	// Ensure `UpdatedAt` is set if you have such a field and GORM doesn't handle it automatically
	// based on your model definition or hooks. GORM typically handles `UpdatedAt` automatically.
	result := r.db.WithContext(ctx).Model(user).
		// Select specific columns to update. GORM updates non-zero fields by default.
		// If you want to update all fields including zero values, use Select("*")
		// or specify all columns. For partial updates, this is good.
		Select("Email", "FirstName", "LastName", "Username", "IsSuperAdmin"). // Note: GORM uses struct field names. UpdatedAt is handled by GORM.
		Where("id = ?", user.ID).
		Updates(user) // Updates will only update non-zero fields of the user struct unless specified in Select.
		// If you want to update specific fields to their zero values, use a map:
		// Updates(map[string]interface{}{"name": user.Name, "email": user.Email, ...})

	if result.Error != nil {
		return fmt.Errorf("failed to update user: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return ErrUserNotFound // Or handle as a specific case where no update occurred
	}
	return nil
}

func (r *userRepository) GetByID(ctx context.Context, id string) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return &user, nil
}

func (r *userRepository) GetUsers(ctx context.Context, offset, limit int) ([]*model.User, int64, error) {
	var users []*model.User
	var totalCount int64

	// Count total records
	err := r.db.WithContext(ctx).Model((*model.User)(nil)).Count(&totalCount).Error
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count users: %w", err)
	}

	// Get paginated records
	err = r.db.WithContext(ctx).Offset(offset).Limit(limit).Find(&users).Error
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
	// Assumes model.User has a gorm.DeletedAt field for soft delete
	result := r.db.WithContext(ctx).Where("id = ?", id).Delete(&model.User{})
	if result.Error != nil {
		return fmt.Errorf("failed to soft delete user: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return ErrUserNotFound
	}
	return nil
}

func (r *userRepository) Exists(ctx context.Context, email string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model((*model.User)(nil)).Where("email = ?", email).Count(&count).Error
	if err != nil {
		// If the error is record not found, it means no user with that email exists.
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, fmt.Errorf("failed to check user existence: %w", err)
	}
	return count > 0, nil
}
