package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"log"
	"time"

	"echobackend/config"
	apperrors "echobackend/internal/errors"
	"echobackend/internal/model"
	"echobackend/internal/repository"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Register(ctx context.Context, email, username, password string) (*model.User, error)
	Login(ctx context.Context, email, password string) (string, string, *model.User, error)
	CheckUsernameAvailability(ctx context.Context, username string) (bool, error)
	ForgotPassword(ctx context.Context, email string) error
	ResetPassword(ctx context.Context, token, password string) error
	RefreshToken(ctx context.Context, refreshToken string) (string, string, *model.User, error)
	ChangePassword(ctx context.Context, userID, currentPassword, newPassword string) error
}

type authService struct {
	authRepo    repository.AuthRepository
	userRepo    repository.UserRepository
	sessionRepo repository.SessionRepository
	jwtSecret   []byte
}

func NewAuthService(authRepo repository.AuthRepository, userRepo repository.UserRepository, sessionRepo repository.SessionRepository, config *config.Config) AuthService {
	return &authService{
		authRepo:    authRepo,
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
		jwtSecret:   []byte(config.JWTSecret),
	}
}

func (s *authService) Register(ctx context.Context, email, username, password string) (*model.User, error) {
	_, err := s.authRepo.FindUserByEmail(ctx, email)
	if err == nil {
		return nil, apperrors.ErrUserExists
	}
	if err != nil && err != apperrors.ErrUserNotFound {
		return nil, err
	}

	err = s.userRepo.CheckUserByUsername(ctx, username)
	if err != nil {
		return nil, err
	}

	hashedPasswordBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	hashedPassword := string(hashedPasswordBytes)

	newUser := &model.User{
		Email:    email,
		Username: &username,
		Password: &hashedPassword,
	}

	if err := s.authRepo.CreateUser(ctx, newUser); err != nil {
		return nil, err
	}

	return newUser, nil
}

func (s *authService) Login(ctx context.Context, email, password string) (string, string, *model.User, error) {
	user, err := s.authRepo.FindUserByEmail(ctx, email)
	if err != nil {
		return "", "", nil, apperrors.ErrInvalidCredentials
	}

	if user.Password == nil {
		return "", "", nil, apperrors.ErrInvalidCredentials
	}

	if compareErr := bcrypt.CompareHashAndPassword([]byte(*user.Password), []byte(password)); compareErr != nil {
		return "", "", nil, apperrors.ErrInvalidCredentials
	}

	claims := jwt.MapClaims{
		"user_id":        user.ID,
		"username":       user.Username,
		"email":          user.Email,
		"is_super_admin": user.IsSuperAdmin,
		"iat":            time.Now().Unix(),
		"exp":            time.Now().Add(3 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(s.jwtSecret)
	if err != nil {
		return "", "", nil, err
	}

	refreshBytes := make([]byte, 64)
	if _, err := rand.Read(refreshBytes); err != nil {
		return "", "", nil, err
	}
	refreshToken := "pl_" + base64.RawURLEncoding.EncodeToString(refreshBytes)

	sess := &model.Session{
		RefreshToken: refreshToken,
		UserID:       user.ID,
	}
	if err := s.sessionRepo.CreateSession(ctx, sess); err != nil {
		return "", "", nil, err
	}

	return tokenString, refreshToken, user, nil
}

func (s *authService) CheckUsernameAvailability(ctx context.Context, username string) (bool, error) {
	err := s.userRepo.CheckUserByUsername(ctx, username)
	if err == nil {
		return true, nil
	}
	if err == apperrors.ErrUserExists {
		return false, nil
	}
	return false, err
}

func (s *authService) ForgotPassword(ctx context.Context, email string) error {
	resetToken := "pr_" + base64.RawURLEncoding.EncodeToString(generateRandomBytes(32))
	expiresAt := time.Now().Add(1 * time.Hour)

	log.Printf("Password reset token for %s: %s (expires at %s)", email, resetToken, expiresAt.Format(time.RFC3339))

	return nil
}

func (s *authService) ResetPassword(ctx context.Context, token, password string) error {
	if !isValidPasswordResetToken(token) {
		return apperrors.ErrInvalidToken
	}

	return apperrors.ErrInvalidToken
}

func (s *authService) RefreshToken(ctx context.Context, refreshToken string) (string, string, *model.User, error) {
	session, err := s.sessionRepo.GetSessionByRefreshToken(ctx, refreshToken)
	if err != nil {
		return "", "", nil, apperrors.ErrInvalidToken
	}

	if session.ExpiresAt != nil && time.Now().After(*session.ExpiresAt) {
		s.sessionRepo.DeleteSession(ctx, refreshToken)
		return "", "", nil, apperrors.ErrTokenExpired
	}

	user, err := s.userRepo.GetByID(ctx, session.UserID)
	if err != nil {
		return "", "", nil, err
	}

	claims := jwt.MapClaims{
		"user_id":        user.ID,
		"username":       user.Username,
		"email":          user.Email,
		"is_super_admin": user.IsSuperAdmin,
		"iat":            time.Now().Unix(),
		"exp":            time.Now().Add(3 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(s.jwtSecret)
	if err != nil {
		return "", "", nil, err
	}

	newRefreshBytes := make([]byte, 64)
	if _, err := rand.Read(newRefreshBytes); err != nil {
		return "", "", nil, err
	}
	newRefreshToken := "pl_" + base64.RawURLEncoding.EncodeToString(newRefreshBytes)

	session.RefreshToken = newRefreshToken
	session.CreatedAt = &time.Time{}
	if err := s.sessionRepo.UpdateSession(ctx, session); err != nil {
		return "", "", nil, err
	}

	return tokenString, newRefreshToken, user, nil
}

func (s *authService) ChangePassword(ctx context.Context, userID, currentPassword, newPassword string) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return apperrors.ErrUserNotFound
	}

	if user.Password == nil {
		return apperrors.ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(*user.Password), []byte(currentPassword)); err != nil {
		return apperrors.ErrInvalidCredentials
	}

	hashedPasswordBytes, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	hashedPassword := string(hashedPasswordBytes)

	user.Password = &hashedPassword
	return s.userRepo.Update(ctx, user)
}

func generateRandomBytes(n int) []byte {
	b := make([]byte, n)
	rand.Read(b)
	return b
}

func isValidPasswordResetToken(token string) bool {
	return len(token) > 10 && token[:3] == "pr_"
}
