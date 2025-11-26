package service

import (
	"context"
	"echobackend/config"
	"echobackend/internal/model"
	"echobackend/internal/repository"
	"errors"
	"fmt"
	"time"

	"crypto/rand"
	"encoding/base64"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserExists         = errors.New("user already exists")
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidToken       = errors.New("invalid or expired token")
	ErrTokenExpired       = errors.New("token has expired")
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
		fmt.Println("email not found")
		fmt.Println(err)
		return "", "", nil, ErrInvalidCredentials
	}

	if user.Password == nil {
		return "", "", nil, ErrInvalidCredentials
	}

	if compareErr := bcrypt.CompareHashAndPassword([]byte(*user.Password), []byte(password)); compareErr != nil {
		return "", "", nil, ErrInvalidCredentials
	}

	// Generate JWT token
	claims := jwt.MapClaims{
		"user_id":        user.ID,
		"username":       user.Username,
		"email":          user.Email,
		"is_super_admin": user.IsSuperAdmin,
		"iat":            time.Now().Unix(),
		"exp":            time.Now().Add(6 * time.Hour).Unix(), // Token expires after 6 hours
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(s.jwtSecret)
	if err != nil {
		return "", "", nil, err
	}

	// Generate secure refresh token
	refreshBytes := make([]byte, 64)
	if _, err := rand.Read(refreshBytes); err != nil {
		return "", "", nil, err
	}
	refreshToken := "pl_" + base64.RawURLEncoding.EncodeToString(refreshBytes)

	// Persist session
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
	if err == repository.ErrUserExists {
		return false, nil // Username is taken
	}
	if err != nil {
		return false, err // Database error
	}
	return true, nil // Username is available
}

func (s *authService) ForgotPassword(ctx context.Context, email string) error {
	// In a real implementation, you'd check if the user exists
	// For security, we don't reveal whether the email exists or not

	// Generate password reset token
	resetToken := "pr_" + base64.RawURLEncoding.EncodeToString(generateRandomBytes(32))
	expiresAt := time.Now().Add(1 * time.Hour) // Token expires in 1 hour

	// Store the reset token (in a real app, you'd send an email)
	// For now, we'll store it in the session table or create a password_reset_tokens table
	// This is a simplified implementation - in production, you'd want to store this properly
	fmt.Printf("Password reset token for %s: %s (expires at %s)\n", email, resetToken, expiresAt.Format(time.RFC3339))

	return nil
}

func (s *authService) ResetPassword(ctx context.Context, token, password string) error {
	// In a real implementation, you'd verify the token against a stored hash
	// and check if it has expired
	if !isValidPasswordResetToken(token) {
		return ErrInvalidToken
	}

	// Extract email from token (simplified - in reality you'd decode from a JWT or look up in DB)
	// For this implementation, we'll need to store the token mapping
	// This is a simplified version - you'd want to implement proper token verification

	// For now, let's implement a basic version that finds user by token
	// In production, you'd have a password_reset_tokens table
	return ErrInvalidToken // Simplified for now
}

func (s *authService) RefreshToken(ctx context.Context, refreshToken string) (string, string, *model.User, error) {
	// Verify the refresh token exists and is valid
	session, err := s.sessionRepo.GetSessionByRefreshToken(ctx, refreshToken)
	if err != nil {
		return "", "", nil, ErrInvalidToken
	}

	// Check if session has expired
	if session.ExpiresAt != nil && time.Now().After(*session.ExpiresAt) {
		// Clean up expired session
		s.sessionRepo.DeleteSession(ctx, refreshToken)
		return "", "", nil, ErrTokenExpired
	}

	// Get user information
	user, err := s.userRepo.GetByID(ctx, session.UserID)
	if err != nil {
		return "", "", nil, err
	}

	// Generate new JWT token
	claims := jwt.MapClaims{
		"user_id":        user.ID,
		"username":       user.Username,
		"email":          user.Email,
		"is_super_admin": user.IsSuperAdmin,
		"iat":            time.Now().Unix(),
		"exp":            time.Now().Add(6 * time.Hour).Unix(), // Token expires after 6 hours
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(s.jwtSecret)
	if err != nil {
		return "", "", nil, err
	}

	// Generate new refresh token and update session
	newRefreshBytes := make([]byte, 64)
	if _, err := rand.Read(newRefreshBytes); err != nil {
		return "", "", nil, err
	}
	newRefreshToken := "pl_" + base64.RawURLEncoding.EncodeToString(newRefreshBytes)

	// Update session with new refresh token
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
		return ErrUserNotFound
	}

	if user.Password == nil {
		return ErrInvalidCredentials
	}

	// Verify current password
	if err := bcrypt.CompareHashAndPassword([]byte(*user.Password), []byte(currentPassword)); err != nil {
		return ErrInvalidCredentials
	}

	// Hash new password
	hashedPasswordBytes, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	hashedPassword := string(hashedPasswordBytes)

	// Update password
	user.Password = &hashedPassword
	return s.userRepo.Update(ctx, user)
}

// Helper functions
func generateRandomBytes(n int) []byte {
	b := make([]byte, n)
	rand.Read(b)
	return b
}

func isValidPasswordResetToken(token string) bool {
	// In a real implementation, you'd verify this against stored tokens
	// This is a simplified check for the format
	return len(token) > 10 && token[:3] == "pr_"
}
