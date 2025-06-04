package repository

import (
	"context"
	"echobackend/internal/model"
	"errors" // Added for errors.Is

	"gorm.io/gorm" // Added for gorm.DB and gorm.ErrRecordNotFound
)

// Re-define or import shared errors if necessary. For now, let's assume ErrUserNotFound is a concept.
// We can use gorm.ErrRecordNotFound directly or wrap it.
var (
	ErrUserNotFoundAuth = errors.New("user not found for auth") // Specific error for this context or use a shared one
)

type AuthRepository interface {
	FindUserByEmail(ctx context.Context, email string) (*model.User, error)
	CreateUser(ctx context.Context, user *model.User) error
}

type authRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) AuthRepository {
	return &authRepository{db: db}
}

func (r *authRepository) FindUserByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Option 1: Return a specific error for this repository
			// return nil, ErrUserNotFoundAuth
			// Option 2: Return gorm.ErrRecordNotFound directly if services handle it
			// return nil, err
			// Option 3: For consistency with UserRepository, let's use a similar pattern.
			// If ErrUserNotFound is meant to be a general "user not found", it should be defined in a shared place.
			// For now, returning a wrapped error or gorm.ErrRecordNotFound.
			// Let's return a new error instance for clarity or the specific gorm error.
			return nil, gorm.ErrRecordNotFound // Or a custom error like ErrUserNotFound from user_repository
		}
		return nil, err // Other errors
	}
	return &user, nil
}

func (r *authRepository) CreateUser(ctx context.Context, user *model.User) error {
	// This CreateUser does not check for existence like in UserRepository.
	// It will rely on DB constraints (e.g., unique email) to prevent duplicates.
	// GORM's Create will return an error if a constraint is violated.
	result := r.db.WithContext(ctx).Create(user)
	return result.Error
}
