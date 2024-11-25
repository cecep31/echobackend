package service

import (
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
	Register(email, password string) (*model.User, error)
	Login(email, password string) (string, *model.User, error)
}

type authService struct {
	repo      repository.AuthRepository
	jwtSecret []byte
}

func NewAuthService(repo repository.AuthRepository, config *config.Config) AuthService {
	return &authService{
		repo:      repo,
		jwtSecret: []byte(config.JWT_SECRET),
	}
}

func (s *authService) Register(email, password string) (*model.User, error) {
	// Check if user exists
	if _, err := s.repo.FindUserByEmail(email); err == nil {
		return nil, ErrUserExists
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Email:    email,
		Password: string(hashedPassword),
	}

	if err := s.repo.CreateUser(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *authService) Login(email, password string) (string, *model.User, error) {
	user, err := s.repo.FindUserByEmail(email)
	if err != nil {
		fmt.Println("email not found")
		fmt.Println(err)
		return "", nil, ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", nil, ErrInvalidCredentials
	}

	// Generate JWT token
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(s.jwtSecret)
	if err != nil {
		return "", nil, err
	}

	return tokenString, user, nil
}
