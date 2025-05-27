package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"echobackend/internal/model"

	"github.com/uptrace/bun"
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
	db *bun.DB
}

func NewUserRepository(db *bun.DB) UserRepository {
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

	_, err = r.db.NewInsert().
		Model(user).
		Exec(ctx)
	return err
}

func (r *userRepository) Update(ctx context.Context, user *model.User) error {
	res, err := r.db.NewUpdate().
		Model(user).
		Column("name", "email", "username", "bio", "avatar", "updated_at").
		Where("id = ?", user.ID).
		Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return ErrUserNotFound
	}
	return nil
}

func (r *userRepository) GetByID(ctx context.Context, id string) (*model.User, error) {
	var user model.User
	err := r.db.NewSelect().
		Model(&user).
		Where("id = ?", id).
		Limit(1).
		Scan(ctx)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return &user, nil
}

func (r *userRepository) GetUsers(ctx context.Context, offset, limit int) ([]*model.User, int64, error) {
	var users []*model.User

	// Count total records
	totalCount, err := r.db.NewSelect().
		Model((*model.User)(nil)).
		Count(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count users: %w", err)
	}
	total := int64(totalCount)

	// Get paginated records
	err = r.db.NewSelect().
		Model(&users).
		Offset(offset).
		Limit(limit).
		Scan(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get users: %w", err)
	}

	return users, total, nil
}

func (r *userRepository) GetUsersByEmail(ctx context.Context, email string) ([]*model.User, error) {
	var users []*model.User
	err := r.db.NewSelect().
		Model(&users).
		Where("email = ?", email).
		Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get users by email: %w", err)
	}
	return users, nil
}

func (r *userRepository) SoftDeleteByID(ctx context.Context, id string) error {
	res, err := r.db.NewDelete().
		Model(&model.User{}).
		Where("id = ?", id).
		Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return ErrUserNotFound
	}
	return nil
}

func (r *userRepository) Exists(ctx context.Context, email string) (bool, error) {
	count, err := r.db.NewSelect().Model((*model.User)(nil)).Where("email = ?", email).Count(ctx)
	if err != nil {
		return false, fmt.Errorf("failed to check user existence: %w", err)
	}
	return count > 0, nil
}
