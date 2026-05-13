package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"net/url"
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
	Login(ctx context.Context, identifier, password, ipAddress, userAgent string) (string, string, *model.User, error)
	CheckUsernameAvailability(ctx context.Context, username string) (bool, error)
	CheckEmailAvailability(ctx context.Context, email string) (bool, error)
	ForgotPassword(ctx context.Context, email, ipAddress, userAgent string) error
	ResetPassword(ctx context.Context, token, password, ipAddress, userAgent string) error
	RefreshToken(ctx context.Context, refreshToken, ipAddress, userAgent string) (string, string, *model.User, error)
	ChangePassword(ctx context.Context, userID, currentPassword, newPassword, ipAddress, userAgent string) error
	Logout(ctx context.Context, refreshToken string) error
	GetProfile(ctx context.Context, userID string) (*model.User, error)
	GetGithubOAuthURL() string
	GetGithubToken(code string) (string, error)
	SignInWithGithub(ctx context.Context, githubUser *GithubUser, ipAddress, userAgent string) (string, string, *model.User, error)
}

type GithubUser struct {
	Login      string  `json:"login"`
	ID         int64   `json:"id"`
	AvatarURL  string  `json:"avatar_url"`
	Email      *string `json:"email"`
	Name       string  `json:"name"`
	HTMLURL    string  `json:"html_url"`
}

type authService struct {
	authRepo               repository.AuthRepository
	userRepo               repository.UserRepository
	sessionRepo            repository.SessionRepository
	passwordResetTokenRepo repository.PasswordResetTokenRepository
	activityService        AuthActivityService
	jwtSecret              []byte
	jwtExpiry              time.Duration
	githubConfig           config.GitHubConfig
}

func NewAuthService(
	authRepo repository.AuthRepository,
	userRepo repository.UserRepository,
	sessionRepo repository.SessionRepository,
	passwordResetTokenRepo repository.PasswordResetTokenRepository,
	activityService AuthActivityService,
	config *config.Config,
) AuthService {
	return &authService{
		authRepo:               authRepo,
		userRepo:               userRepo,
		sessionRepo:            sessionRepo,
		passwordResetTokenRepo: passwordResetTokenRepo,
		activityService:        activityService,
		jwtSecret:              []byte(config.Auth.JWTSecret),
		jwtExpiry:              config.Auth.JWTExpiry,
		githubConfig:           config.GitHub,
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

func (s *authService) Login(ctx context.Context, identifier, password, ipAddress, userAgent string) (string, string, *model.User, error) {
	user, err := s.authRepo.FindUserByIdentifier(ctx, identifier)
	if err != nil {
		s.activityService.LogActivity(ctx, nil, model.ActivityLoginFailed, model.StatusFailure, ipAddress, userAgent, nil, nil)
		return "", "", nil, apperrors.ErrInvalidCredentials
	}

	if user.Password == nil {
		s.activityService.LogActivity(ctx, &user.ID, model.ActivityLoginFailed, model.StatusFailure, ipAddress, userAgent, nil, nil)
		return "", "", nil, apperrors.ErrInvalidCredentials
	}

	if compareErr := bcrypt.CompareHashAndPassword([]byte(*user.Password), []byte(password)); compareErr != nil {
		s.activityService.LogActivity(ctx, &user.ID, model.ActivityLoginFailed, model.StatusFailure, ipAddress, userAgent, nil, nil)
		return "", "", nil, apperrors.ErrInvalidCredentials
	}

	tokenString, refreshToken, err := s.createTokenAndSession(ctx, user)
	if err != nil {
		return "", "", nil, err
	}

	s.activityService.LogActivity(ctx, &user.ID, model.ActivityLogin, model.StatusSuccess, ipAddress, userAgent, nil, nil)

	now := time.Now()
	user.LastLoggedAt = &now
	_ = s.userRepo.Update(ctx, user)

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

func (s *authService) CheckEmailAvailability(ctx context.Context, email string) (bool, error) {
	_, err := s.authRepo.FindUserByEmail(ctx, email)
	if err == nil {
		return false, nil
	}
	if err == apperrors.ErrUserNotFound {
		return true, nil
	}
	return false, err
}

func (s *authService) ForgotPassword(ctx context.Context, email, ipAddress, userAgent string) error {
	user, err := s.authRepo.FindUserByEmail(ctx, email)
	if err != nil {
		return nil
	}

	resetToken := "pr_" + base64.RawURLEncoding.EncodeToString(generateRandomBytes(32))
	expiresAt := time.Now().Add(1 * time.Hour)

	tokenEntry := &model.PasswordResetToken{
		UserID:    user.ID,
		Token:     resetToken,
		ExpiresAt: expiresAt,
	}

	if err := s.passwordResetTokenRepo.DeleteByUserID(ctx, user.ID); err != nil {
		_ = err
	}

	if err := s.passwordResetTokenRepo.Create(ctx, tokenEntry); err != nil {
		return err
	}

	s.activityService.LogActivity(ctx, &user.ID, model.ActivityPasswordResetReq, model.StatusSuccess, ipAddress, userAgent, nil, nil)

	fmt.Printf("Password reset token for %s: %s (expires at %s)\n", email, resetToken, expiresAt.Format(time.RFC3339))
	return nil
}

func (s *authService) ResetPassword(ctx context.Context, token, password, ipAddress, userAgent string) error {
	tokenEntry, err := s.passwordResetTokenRepo.FindByToken(ctx, token)
	if err != nil {
		return apperrors.ErrInvalidToken
	}
	if tokenEntry == nil {
		return apperrors.ErrInvalidToken
	}

	if tokenEntry.UsedAt != nil {
		return apperrors.ErrPasswordResetTokenUsed
	}

	if time.Now().After(tokenEntry.ExpiresAt) {
		return apperrors.ErrPasswordResetTokenExpired
	}

	user, err := s.userRepo.GetByID(ctx, tokenEntry.UserID)
	if err != nil {
		return apperrors.ErrUserNotFound
	}

	hashedPasswordBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	hashedPassword := string(hashedPasswordBytes)

	user.Password = &hashedPassword
	if err := s.userRepo.Update(ctx, user); err != nil {
		return err
	}

	if err := s.passwordResetTokenRepo.MarkUsed(ctx, tokenEntry.ID); err != nil {
		_ = err
	}

	_ = s.sessionRepo.DeleteByUserID(ctx, user.ID)

	s.activityService.LogActivity(ctx, &user.ID, model.ActivityPasswordReset, model.StatusSuccess, ipAddress, userAgent, nil, nil)

	return nil
}

func (s *authService) RefreshToken(ctx context.Context, refreshToken, ipAddress, userAgent string) (string, string, *model.User, error) {
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

	tokenString, newRefreshToken, err := s.createTokenAndSession(ctx, user)
	if err != nil {
		return "", "", nil, err
	}

	_ = s.sessionRepo.DeleteSession(ctx, refreshToken)

	s.activityService.LogActivity(ctx, &user.ID, model.ActivityTokenRefresh, model.StatusSuccess, ipAddress, userAgent, nil, nil)

	return tokenString, newRefreshToken, user, nil
}

func (s *authService) ChangePassword(ctx context.Context, userID, currentPassword, newPassword, ipAddress, userAgent string) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return apperrors.ErrUserNotFound
	}

	if user.Password == nil {
		return apperrors.ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(*user.Password), []byte(currentPassword)); err != nil {
		s.activityService.LogActivity(ctx, &userID, model.ActivityPasswordChange, model.StatusFailure, ipAddress, userAgent, nil, nil)
		return apperrors.ErrInvalidCredentials
	}

	hashedPasswordBytes, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	hashedPassword := string(hashedPasswordBytes)

	user.Password = &hashedPassword
	if err := s.userRepo.Update(ctx, user); err != nil {
		return err
	}

	s.activityService.LogActivity(ctx, &userID, model.ActivityPasswordChange, model.StatusSuccess, ipAddress, userAgent, nil, nil)

	return nil
}

func (s *authService) Logout(ctx context.Context, refreshToken string) error {
	return s.sessionRepo.DeleteSession(ctx, refreshToken)
}

func (s *authService) GetProfile(ctx context.Context, userID string) (*model.User, error) {
	return s.userRepo.GetByID(ctx, userID)
}

func (s *authService) GetGithubOAuthURL() string {
	authURL, _ := url.Parse("https://github.com/login/oauth/authorize")
	q := authURL.Query()
	q.Set("client_id", s.githubConfig.ClientID)
	q.Set("redirect_uri", s.githubConfig.RedirectURI)
	q.Set("scope", "user:email")
	authURL.RawQuery = q.Encode()
	return authURL.String()
}

func (s *authService) GetGithubToken(code string) (string, error) {
	data := url.Values{}
	data.Set("client_id", s.githubConfig.ClientID)
	data.Set("client_secret", s.githubConfig.ClientSecret)
	data.Set("code", code)
	data.Set("redirect_uri", s.githubConfig.RedirectURI)

	resp, err := http.PostForm("https://github.com/login/oauth/access_token", data)
	if err != nil {
		return "", fmt.Errorf("failed to exchange code for token: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read token response: %w", err)
	}

	values, err := url.ParseQuery(string(body))
	if err != nil {
		return "", fmt.Errorf("failed to parse token response: %w", err)
	}

	accessToken := values.Get("access_token")
	if accessToken == "" {
		return "", fmt.Errorf("no access_token in GitHub response")
	}

	return accessToken, nil
}

func (s *authService) SignInWithGithub(ctx context.Context, githubUser *GithubUser, ipAddress, userAgent string) (string, string, *model.User, error) {
	var user *model.User

	githubID := githubUser.ID
	user, err := s.authRepo.FindUserByGithubID(ctx, githubID)
	if err != nil && err != apperrors.ErrUserNotFound {
		s.activityService.LogActivity(ctx, nil, model.ActivityOAuthLoginFailed, model.StatusFailure, ipAddress, userAgent, nil, map[string]any{"provider": "github"})
		return "", "", nil, err
	}

	if user == nil || err == apperrors.ErrUserNotFound {
		email := ""
		if githubUser.Email != nil {
			email = *githubUser.Email
		} else {
			email = fmt.Sprintf("%d@github.placeholder", githubUser.ID)
		}

		username := githubUser.Login
		newUser := &model.User{
			Email:        email,
			Username:     &username,
			GithubID:     &githubID,
			Image:        &githubUser.AvatarURL,
		}

		if err := s.authRepo.CreateUser(ctx, newUser); err != nil {
			s.activityService.LogActivity(ctx, nil, model.ActivityOAuthLoginFailed, model.StatusFailure, ipAddress, userAgent, nil, map[string]any{"provider": "github", "error": err.Error()})
			return "", "", nil, err
		}
		user = newUser
	}

	tokenString, refreshToken, err := s.createTokenAndSession(ctx, user)
	if err != nil {
		s.activityService.LogActivity(ctx, &user.ID, model.ActivityOAuthLoginFailed, model.StatusFailure, ipAddress, userAgent, nil, map[string]any{"provider": "github"})
		return "", "", nil, err
	}

	s.activityService.LogActivity(ctx, &user.ID, model.ActivityOAuthLogin, model.StatusSuccess, ipAddress, userAgent, nil, map[string]any{"provider": "github"})

	now := time.Now()
	user.LastLoggedAt = &now
	_ = s.userRepo.Update(ctx, user)

	return tokenString, refreshToken, user, nil
}

func (s *authService) createTokenAndSession(ctx context.Context, user *model.User) (string, string, error) {
	claims := jwt.MapClaims{
		"user_id":        user.ID,
		"username":       user.Username,
		"email":          user.Email,
		"is_super_admin": user.IsSuperAdmin,
		"iat":            time.Now().Unix(),
		"exp":            time.Now().Add(s.jwtExpiry).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(s.jwtSecret)
	if err != nil {
		return "", "", err
	}

	refreshBytes := make([]byte, 64)
	if _, err := rand.Read(refreshBytes); err != nil {
		return "", "", err
	}
	refreshToken := "pl_" + base64.RawURLEncoding.EncodeToString(refreshBytes)

	sess := &model.Session{
		RefreshToken: refreshToken,
		UserID:       user.ID,
	}
	if err := s.sessionRepo.CreateSession(ctx, sess); err != nil {
		return "", "", err
	}

	return tokenString, refreshToken, nil
}

func generateRandomBytes(n int) []byte {
	b := make([]byte, n)
	rand.Read(b)
	return b
}