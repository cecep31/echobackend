package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"echobackend/config"
	"echobackend/internal/dto"
	apperrors "echobackend/internal/errors"
	"echobackend/internal/model"
	"echobackend/internal/service"
	"echobackend/pkg/response"

	"github.com/labstack/echo/v5"
)

type AuthHandler struct {
	authService     service.AuthService
	activityService service.AuthActivityService
	frontendConfig  config.FrontendConfig
}

func NewAuthHandler(authService service.AuthService, activityService service.AuthActivityService, frontendConfig config.FrontendConfig) *AuthHandler {
	return &AuthHandler{
		authService:     authService,
		activityService: activityService,
		frontendConfig:  frontendConfig,
	}
}

func (h *AuthHandler) Register(c *echo.Context) error {
	var req dto.RegisterRequest
	if err := c.Bind(&req); err != nil {
		return response.BadRequest(c, "Invalid request format", err)
	}

	if err := c.Validate(req); err != nil {
		return response.FromValidateError(c, err)
	}

	user, err := h.authService.Register(c.Request().Context(), req.Email, req.Username, req.Password)
	if err == apperrors.ErrUserExists {
		return response.Conflict(c, "Registration failed", "Email or username already exists")
	}
	if err != nil {
		return response.InternalServerError(c, "Registration failed", err)
	}

	ipAddress := c.RealIP()
	userAgent := c.Request().UserAgent()
	h.activityService.LogActivity(c.Request().Context(), &user.ID, model.ActivityRegister, model.StatusSuccess, ipAddress, userAgent, nil, nil)

	return response.Created(c, "User registered successfully", map[string]any{
		"id":       user.ID,
		"email":    user.Email,
		"username": user.Username,
	})
}

func (h *AuthHandler) Login(c *echo.Context) error {
	var loginReq dto.LoginRequest
	if err := c.Bind(&loginReq); err != nil {
		return response.BadRequest(c, "Invalid request format", err)
	}

	if err := c.Validate(loginReq); err != nil {
		return response.FromValidateError(c, err)
	}

	ipAddress := c.RealIP()
	userAgent := c.Request().UserAgent()

	token, refreshToken, user, err := h.authService.Login(c.Request().Context(), loginReq.Identifier, loginReq.Password, ipAddress, userAgent)
	if err == apperrors.ErrInvalidCredentials {
		return response.Unauthorized(c, "Invalid identifier or password")
	}
	if err != nil {
		return response.InternalServerError(c, "Login failed", err)
	}

	return response.Success(c, "Login successful", map[string]any{
		"access_token":  token,
		"refresh_token": refreshToken,
		"user": map[string]any{
			"id":       user.ID,
			"email":    user.Email,
			"username": user.Username,
		},
	})
}

func (h *AuthHandler) CheckUsername(c *echo.Context) error {
	var req dto.CheckUsernameRequest
	if err := c.Bind(&req); err != nil {
		return response.BadRequest(c, "Invalid request format", err)
	}

	if err := c.Validate(req); err != nil {
		return response.FromValidateError(c, err)
	}

	isAvailable, err := h.authService.CheckUsernameAvailability(c.Request().Context(), req.Username)
	if err != nil {
		return response.InternalServerError(c, "Failed to check username availability", err)
	}

	return response.Success(c, "Username availability checked", map[string]any{
		"username":  req.Username,
		"available": isAvailable,
	})
}

func (h *AuthHandler) ForgotPassword(c *echo.Context) error {
	var req dto.ForgotPasswordRequest
	if err := c.Bind(&req); err != nil {
		return response.BadRequest(c, "Invalid request format", err)
	}

	if err := c.Validate(req); err != nil {
		return response.FromValidateError(c, err)
	}

	ipAddress := c.RealIP()
	userAgent := c.Request().UserAgent()

	err := h.authService.ForgotPassword(c.Request().Context(), req.Email, ipAddress, userAgent)
	if err != nil {
		return response.Success(c, "If the email exists, a password reset link has been sent", nil)
	}

	return response.Success(c, "If the email exists, a password reset link has been sent", nil)
}

func (h *AuthHandler) ResetPassword(c *echo.Context) error {
	var req dto.ResetPasswordRequest
	if err := c.Bind(&req); err != nil {
		return response.BadRequest(c, "Invalid request format", err)
	}

	if err := c.Validate(req); err != nil {
		return response.FromValidateError(c, err)
	}

	ipAddress := c.RealIP()
	userAgent := c.Request().UserAgent()

	err := h.authService.ResetPassword(c.Request().Context(), req.Token, req.Password, ipAddress, userAgent)
	if err == apperrors.ErrInvalidToken {
		return response.BadRequest(c, "Invalid or expired reset token", err)
	}
	if err == apperrors.ErrPasswordResetTokenUsed {
		return response.BadRequest(c, "Reset token has already been used", err)
	}
	if err == apperrors.ErrPasswordResetTokenExpired {
		return response.BadRequest(c, "Reset token has expired", err)
	}
	if err != nil {
		return response.InternalServerError(c, "Failed to reset password", err)
	}

	return response.Success(c, "Password reset successful", nil)
}

func (h *AuthHandler) RefreshToken(c *echo.Context) error {
	var req dto.RefreshTokenRequest
	if err := c.Bind(&req); err != nil {
		return response.BadRequest(c, "Invalid request format", err)
	}

	if err := c.Validate(req); err != nil {
		return response.FromValidateError(c, err)
	}

	ipAddress := c.RealIP()
	userAgent := c.Request().UserAgent()

	token, refreshToken, user, err := h.authService.RefreshToken(c.Request().Context(), req.RefreshToken, ipAddress, userAgent)
	if err == apperrors.ErrInvalidToken {
		return response.Unauthorized(c, "Invalid or expired refresh token")
	}
	if err == apperrors.ErrTokenExpired {
		return response.Unauthorized(c, "Refresh token has expired")
	}
	if err != nil {
		return response.InternalServerError(c, "Failed to refresh token", err)
	}

	return response.Success(c, "Token refreshed successfully", map[string]any{
		"access_token":  token,
		"refresh_token": refreshToken,
		"user": map[string]any{
			"id":       user.ID,
			"email":    user.Email,
			"username": user.Username,
		},
	})
}

func (h *AuthHandler) ChangePassword(c *echo.Context) error {
	userID, ok := GetUserIDFromClaims(c)
	if !ok {
		return response.Unauthorized(c, "User not authenticated")
	}

	var req dto.ChangePasswordRequest
	if err := c.Bind(&req); err != nil {
		return response.BadRequest(c, "Invalid request format", err)
	}

	if err := c.Validate(req); err != nil {
		return response.FromValidateError(c, err)
	}

	ipAddress := c.RealIP()
	userAgent := c.Request().UserAgent()

	err := h.authService.ChangePassword(c.Request().Context(), userID, req.CurrentPassword, req.NewPassword, ipAddress, userAgent)
	if err == apperrors.ErrInvalidCredentials {
		return response.Unauthorized(c, "Current password is incorrect")
	}
	if err != nil {
		return response.InternalServerError(c, "Failed to change password", err)
	}

	return response.Success(c, "Password changed successfully", nil)
}

func (h *AuthHandler) Logout(c *echo.Context) error {
	var req dto.LogoutRequest
	if err := c.Bind(&req); err != nil {
		return response.BadRequest(c, "Invalid request format", err)
	}

	if err := c.Validate(req); err != nil {
		return response.FromValidateError(c, err)
	}

	userID, _ := GetUserIDFromClaims(c)
	ipAddress := c.RealIP()
	userAgent := c.Request().UserAgent()

	err := h.authService.Logout(c.Request().Context(), req.RefreshToken)
	if err != nil {
		return response.Success(c, "Logout successful", nil)
	}

	h.activityService.LogActivity(c.Request().Context(), &userID, model.ActivityLogout, model.StatusSuccess, ipAddress, userAgent, nil, nil)

	return response.Success(c, "Logout successful", nil)
}

func (h *AuthHandler) GetProfile(c *echo.Context) error {
	userID, ok := GetUserIDFromClaims(c)
	if !ok {
		return response.Unauthorized(c, "User not authenticated")
	}

	user, err := h.authService.GetProfile(c.Request().Context(), userID)
	if err != nil {
		return response.InternalServerError(c, "Failed to get profile", err)
	}

	return response.Success(c, "Profile retrieved successfully", map[string]any{
		"id":               user.ID,
		"email":            user.Email,
		"username":         user.Username,
		"first_name":       user.FirstName,
		"last_name":        user.LastName,
		"image":            user.Image,
		"is_super_admin":   user.IsSuperAdmin,
		"followers_count":  user.FollowersCount,
		"following_count":  user.FollowingCount,
	})
}

func (h *AuthHandler) GetActivityLogs(c *echo.Context) error {
	userID, ok := GetUserIDFromClaims(c)
	if !ok {
		return response.Unauthorized(c, "User not authenticated")
	}

	limit, offset := ParsePaginationParams(c, 20)
	activityType := c.QueryParam("activity_type")

	logs, totalCount, err := h.activityService.GetActivityLogs(c.Request().Context(), userID, activityType, limit, offset)
	if err != nil {
		return response.InternalServerError(c, "Failed to get activity logs", err)
	}

	meta := response.CalculatePaginationMeta(totalCount, offset, limit)
	return response.SuccessWithMeta(c, "Activity logs retrieved successfully", logs, meta)
}

func (h *AuthHandler) GetRecentActivity(c *echo.Context) error {
	userID, ok := GetUserIDFromClaims(c)
	if !ok {
		return response.Unauthorized(c, "User not authenticated")
	}

	limit := 10
	if l := c.QueryParam("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 && parsed <= 50 {
			limit = parsed
		}
	}

	logs, err := h.activityService.GetRecentActivity(c.Request().Context(), userID, limit)
	if err != nil {
		return response.InternalServerError(c, "Failed to get recent activity", err)
	}

	return response.Success(c, "Recent activity retrieved successfully", logs)
}

func (h *AuthHandler) GetFailedLogins(c *echo.Context) error {
	_, ok := GetUserIDFromClaims(c)
	if !ok {
		return response.Unauthorized(c, "User not authenticated")
	}

	limit, offset := ParsePaginationParams(c, 20)

	hours := 24
	if hParam := c.QueryParam("since_hours"); hParam != "" {
		if parsed, err := strconv.Atoi(hParam); err == nil && parsed > 0 {
			hours = parsed
		}
	}
	since := time.Now().Add(-time.Duration(hours) * time.Hour)

	logs, totalCount, err := h.activityService.GetFailedLogins(c.Request().Context(), since, limit, offset)
	if err != nil {
		return response.InternalServerError(c, "Failed to get failed logins", err)
	}

	meta := response.CalculatePaginationMeta(totalCount, offset, limit)
	return response.SuccessWithMeta(c, "Failed logins retrieved successfully", logs, meta)
}

func (h *AuthHandler) GithubOAuthRedirect(c *echo.Context) error {
	authURL := h.authService.GetGithubOAuthURL()
	return c.Redirect(http.StatusTemporaryRedirect, authURL)
}

func (h *AuthHandler) GithubOAuthCallback(c *echo.Context) error {
	code := c.QueryParam("code")
	if code == "" {
		return c.Redirect(http.StatusTemporaryRedirect, h.frontendConfig.URL+"?error=missing_code")
	}

	ipAddress := c.RealIP()
	userAgent := c.Request().UserAgent()

	githubToken, err := h.authService.GetGithubToken(code)
	if err != nil {
		h.activityService.LogActivity(c.Request().Context(), nil, model.ActivityOAuthLoginFailed, model.StatusFailure, ipAddress, userAgent, nil, map[string]any{"provider": "github", "error": err.Error()})
		return c.Redirect(http.StatusTemporaryRedirect, h.frontendConfig.URL+"?error=github_token_failed")
	}

	githubUser, err := fetchGithubUser(githubToken)
	if err != nil {
		h.activityService.LogActivity(c.Request().Context(), nil, model.ActivityOAuthLoginFailed, model.StatusFailure, ipAddress, userAgent, nil, map[string]any{"provider": "github", "error": err.Error()})
		return c.Redirect(http.StatusTemporaryRedirect, h.frontendConfig.URL+"?error=github_user_failed")
	}

	if githubUser.Email == nil || *githubUser.Email == "" {
		email, err := fetchGithubUserEmail(githubToken)
		if err == nil && email != "" {
			githubUser.Email = &email
		}
	}

	accessToken, refreshToken, user, err := h.authService.SignInWithGithub(c.Request().Context(), githubUser, ipAddress, userAgent)
	if err != nil {
		return c.Redirect(http.StatusTemporaryRedirect, h.frontendConfig.URL+"?error=oauth_login_failed")
	}

	redirectURL := fmt.Sprintf("%s?access_token=%s&refresh_token=%s", h.frontendConfig.URL, accessToken, refreshToken)
	_ = user
	return c.Redirect(http.StatusTemporaryRedirect, redirectURL)
}

func fetchGithubUser(token string) (*service.GithubUser, error) {
	req, err := http.NewRequest("GET", "https://api.github.com/user", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var user service.GithubUser
	if err := json.Unmarshal(body, &user); err != nil {
		return nil, err
	}

	return &user, nil
}

func fetchGithubUserEmail(token string) (string, error) {
	req, err := http.NewRequest("GET", "https://api.github.com/user/emails", nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("GitHub API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var emails []struct {
		Email   string `json:"email"`
		Primary bool   `json:"primary"`
	}
	if err := json.Unmarshal(body, &emails); err != nil {
		return "", err
	}

	for _, e := range emails {
		if e.Primary {
			return e.Email, nil
		}
	}
	if len(emails) > 0 {
		return emails[0].Email, nil
	}
	return "", nil
}

func (h *AuthHandler) CheckEmail(c *echo.Context) error {
	email := c.Param("email")
	if email == "" {
		return response.BadRequest(c, "Email is required", nil)
	}

	isAvailable, err := h.authService.CheckEmailAvailability(c.Request().Context(), email)
	if err != nil {
		return response.InternalServerError(c, "Failed to check email availability", err)
	}

	return response.Success(c, "Email availability checked", map[string]any{
		"email":     email,
		"available": isAvailable,
	})
}