package repository

import (
	"context"
	"echobackend/internal/model"

	"github.com/uptrace/bun"
)

type AuthRepository interface {
	FindUserByEmail(ctx context.Context, email string) (*model.User, error)
	CreateUser(ctx context.Context, user *model.User) error
}

type authRepository struct {
	db *bun.DB
}

func NewAuthRepository(db *bun.DB) AuthRepository {
	return &authRepository{db: db}
}

func (r *authRepository) FindUserByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	err := r.db.NewSelect().
		Model(&user).
		Where("email = ?", email).
		Limit(1).
		Scan(ctx)
	
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *authRepository) CreateUser(ctx context.Context, user *model.User) error {
	_, err := r.db.NewInsert().
		Model(user).
		Exec(ctx)
	return err
}
