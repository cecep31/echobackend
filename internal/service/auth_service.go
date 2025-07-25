package service

import (
	"context"
	"echobackend/config"
	"echobackend/internal/model"
	"echobackend/internal/repository"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserExists         = errors.New("user already exists")
)

type AuthService interface {
	Register(ctx context.Context, email, username, password string) (*model.User, error)
	Login(ctx context.Context, email, password string) (string, *model.User, error)
	CheckUsernameAvailability(ctx context.Context, username string) (bool, error)
}

type authService struct {
	authRepo  repository.AuthRepository
	userRepo  repository.UserRepository
	jwtSecret []byte
}

func NewAuthService(authRepo repository.AuthRepository, userRepo repository.UserRepository, config *config.Config) AuthService {
	return &authService{
		authRepo:  authRepo,
		userRepo:  userRepo,
		jwtSecret: []byte(config.JWT_SECRET),
	}
}

// should be error not hanlde yet
func (s *authService) Register(ctx context.Context, email, username, password string) (*model.User, error) {
	_, err := s.authRepo.FindUserByEmail(ctx, email)
	if err == nil {
		return nil, ErrUserExists
	}

	// Check if username is already taken
	err = s.userRepo.CheckUserByUsername(ctx, username)
	if err == repository.ErrUserExists {
		return nil, ErrUserExists
	}
	if err != nil {
		return nil, err
	}

	hashedPasswordBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	newUser := &model.User{
		Email:    email,
		Username: username,
		Password: string(hashedPasswordBytes),
	}

	if err := s.authRepo.CreateUser(ctx, newUser); err != nil {
		return nil, err
	}

	return newUser, nil
}

func (s *authService) Login(ctx context.Context, email, password string) (string, *model.User, error) {

	user, err := s.authRepo.FindUserByEmail(ctx, email)
	if err != nil {
		fmt.Println("email not found")
		fmt.Println(err)
		return "", nil, ErrInvalidCredentials
	}

	if compareErr := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); compareErr != nil {
		return "", nil, ErrInvalidCredentials
	}

	// Generate JWT token
	claims := jwt.MapClaims{
		"user_id":      user.ID,
		"username":     user.Username,
		"email":        user.Email,
		"isSuperAdmin": user.IsSuperAdmin,
		"iat":          time.Now().Unix(),
		"exp":          time.Now().Add(time.Hour * 48).Unix(), // Token expires after 48 hours
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(s.jwtSecret)
	if err != nil {
		return "", nil, err
	}

	return tokenString, user, nil
}

func (s *authService) CheckUsernameAvailability(ctx context.Context, username string) (bool, error) {
	err := s.userRepo.CheckUserByUsername(ctx, username)
	if err == repository.ErrUserExists {
		return false, nil // Username is taken
	}
	if err != nil {
		return false, err // Database error
	}
	return true, nil // Username is available
}
