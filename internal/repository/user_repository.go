package repository

import (
	"context"
	"echobackend/internal/domain"

	"gorm.io/gorm"
)

type UserRepository interface {
	// Create(ctx context.Context, user *domain.User) error
	GetByID(ctx context.Context, id string) (*domain.User, error)
	GetUsers(ctx context.Context) ([]*domain.User, error)
	GetUsersByEmail(ctx context.Context, email string) ([]*domain.User, error)
	// Update(ctx context.Context, user *domain.User) error
	// Delete(ctx context.Context, id uint) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) GetUsersByEmail(ctx context.Context, email string) ([]*domain.User, error) {
	var users []*domain.User
	return users, r.db.WithContext(ctx).Where("email = ?", email).Find(&users).Error
}

func (r *userRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {
	var user domain.User
	return &user, r.db.WithContext(ctx).Where("id = ?", id).First(&user).Error
}

func (r *userRepository) GetUsers(ctx context.Context) ([]*domain.User, error) {
	var users []*domain.User
	return users, r.db.WithContext(ctx).Find(&users).Error
}
