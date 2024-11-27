package repository

import (
	"context"
	"echobackend/internal/model"

	"gorm.io/gorm"
)

type UserRepository interface {
	// Create(ctx context.Context, user *domain.User) error
	GetByID(ctx context.Context, id string) (*model.User, error)
	GetUsers(ctx context.Context) ([]*model.User, error)
	GetUsersByEmail(ctx context.Context, email string) ([]*model.User, error)
	// Update(ctx context.Context, user *domain.User) error
	// Delete(ctx context.Context, id uint) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) GetUsersByEmail(ctx context.Context, email string) ([]*model.User, error) {
	var users []*model.User
	return users, r.db.WithContext(ctx).Where("email = ?", email).Find(&users).Error
}

func (r *userRepository) GetByID(ctx context.Context, id string) (*model.User, error) {
	var user model.User
	return &user, r.db.WithContext(ctx).Where("id = ?", id).First(&user).Error
}

func (r *userRepository) GetUsers(ctx context.Context) ([]*model.User, error) {
	var users []*model.User
	return users, r.db.WithContext(ctx).Find(&users).Error
}
