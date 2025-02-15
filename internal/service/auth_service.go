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
	authRepo  repository.AuthRepository
	jwtSecret []byte
}

func NewAuthService(repo repository.AuthRepository, config *config.Config) AuthService {
	return &authService{
		authRepo:  repo,
		jwtSecret: []byte(config.JWT_SECRET),
	}
}

// should be error not hanlde yet
func (s *authService) Register(email, password string) (*model.User, error) {
	_, err := s.authRepo.FindUserByEmail(email)
	if err == nil {
		return nil, ErrUserExists
	}

	hashedPasswordBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	newUser := &model.User{
		Email:    email,
		Password: string(hashedPasswordBytes),
	}

	if err := s.authRepo.CreateUser(newUser); err != nil {
		return nil, err
	}

	return newUser, nil
}

func (s *authService) Login(email, password string) (string, *model.User, error) {

	user, err := s.authRepo.FindUserByEmail(email)
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
		"user_id":      user.ID,
		"email":        user.Email,
		"isSuperAdmin": user.IsSuperAdmin,
		"exp":          time.Now().Add(time.Hour * 48).Unix(), // Token expires after 48 hours
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(s.jwtSecret)
	if err != nil {
		return "", nil, err
	}

	return tokenString, user, nil
}
