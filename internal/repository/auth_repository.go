package repository

import (
	"echobackend/internal/model"

	"gorm.io/gorm"
)

type AuthRepository interface {
	FindUserByEmail(email string) (*model.User, error)
	CreateUser(user *model.User) error
}

type authRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) AuthRepository {
	return &authRepository{db: db}
}

func (r *authRepository) FindUserByEmail(email string) (*model.User, error) {
	var user model.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *authRepository) CreateUser(user *model.User) error {
	return r.db.Create(user).Error
}
