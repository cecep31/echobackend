package handler

import (
	"echobackend/internal/dto"
	apperrors "echobackend/internal/errors"
	"echobackend/internal/service"
	"echobackend/pkg/response"

	"github.com/labstack/echo/v5"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
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

	token, refreshToken, user, err := h.authService.Login(c.Request().Context(), loginReq.Identifier, loginReq.Password)
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

	err := h.authService.ForgotPassword(c.Request().Context(), req.Email)
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

	err := h.authService.ResetPassword(c.Request().Context(), req.Token, req.Password)
	if err == apperrors.ErrInvalidToken {
		return response.BadRequest(c, "Invalid or expired reset token", err)
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

	token, refreshToken, user, err := h.authService.RefreshToken(c.Request().Context(), req.RefreshToken)
	if err == apperrors.ErrInvalidToken {
		return response.Unauthorized(c, "Invalid or expired refresh token")
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

	err := h.authService.ChangePassword(c.Request().Context(), userID, req.CurrentPassword, req.NewPassword)
	if err == apperrors.ErrInvalidCredentials {
		return response.Unauthorized(c, "Current password is incorrect")
	}
	if err != nil {
		return response.InternalServerError(c, "Failed to change password", err)
	}

	return response.Success(c, "Password changed successfully", nil)
}
