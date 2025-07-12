package handler

import (
	"echobackend/internal/service"
	"echobackend/pkg/validator"
	"net/http"

	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	authService service.AuthService
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Username string `json:"username" validate:"required,min=3,max=30"`
	Password string `json:"password" validate:"required,min=6"`
}

type CheckUsernameRequest struct {
	Username string `json:"username" validate:"required,min=3,max=30"`
}

// Response represents a standard API response
type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
	Errors  any    `json:"errors,omitempty"`
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Register(c echo.Context) error {
	var req RegisterRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: "Invalid request format",
			Errors:  []string{err.Error()},
		})
	}

	if err := c.Validate(req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			return c.JSON(http.StatusBadRequest, Response{
				Success: false,
				Message: "Validation failed",
				Errors:  validationErrors.Errors,
			})
		}
		return c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: "Validation failed",
			Errors:  []string{err.Error()},
		})
	}

	user, err := h.authService.Register(c.Request().Context(), req.Email, req.Username, req.Password)
	if err == service.ErrUserExists {
		return c.JSON(http.StatusConflict, Response{
			Success: false,
			Message: "Registration failed",
			Errors:  []string{"Email or username already exists"},
		})
	}
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Message: "Registration failed",
			Errors:  []string{"Failed to register user"},
		})
	}

	return c.JSON(http.StatusCreated, Response{
		Success: true,
		Message: "User registered successfully",
		Data: map[string]any{
			"id":       user.ID,
			"email":    user.Email,
			"username": user.Username,
		},
	})
}

func (h *AuthHandler) Login(c echo.Context) error {
	var loginReq LoginRequest
	if err := c.Bind(&loginReq); err != nil {
		return c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: "Invalid request format",
			Errors:  []string{err.Error()},
		})
	}

	if err := c.Validate(loginReq); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			return c.JSON(http.StatusBadRequest, Response{
				Success: false,
				Message: "Validation failed",
				Errors:  validationErrors.Errors,
			})
		}
		return c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: "Validation failed",
			Errors:  []string{err.Error()},
		})
	}

	token, user, err := h.authService.Login(c.Request().Context(), loginReq.Email, loginReq.Password)
	if err == service.ErrInvalidCredentials {
		return c.JSON(http.StatusUnauthorized, Response{
			Success: false,
			Message: "Login failed",
			Errors:  []string{"Invalid email or password"},
		})
	}
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Message: "Login failed",
			Errors:  []string{"Failed to process login"},
		})
	}

	return c.JSON(http.StatusOK, Response{
		Success: true,
		Message: "Login successful",
		Data: map[string]any{
			"access_token": token,
			"user": map[string]any{
				"id":       user.ID,
				"email":    user.Email,
				"username": user.Username,
			},
		},
	})
}

func (h *AuthHandler) CheckUsername(c echo.Context) error {
	var req CheckUsernameRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: "Invalid request format",
			Errors:  []string{err.Error()},
		})
	}

	if err := c.Validate(req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			return c.JSON(http.StatusBadRequest, Response{
				Success: false,
				Message: "Validation failed",
				Errors:  validationErrors.Errors,
			})
		}
		return c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: "Validation failed",
			Errors:  []string{err.Error()},
		})
	}

	isAvailable, err := h.authService.CheckUsernameAvailability(c.Request().Context(), req.Username)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Message: "Failed to check username availability",
			Errors:  []string{"Internal server error"},
		})
	}

	return c.JSON(http.StatusOK, Response{
		Success: true,
		Message: "Username availability checked",
		Data: map[string]any{
			"username":  req.Username,
			"available": isAvailable,
		},
	})
}
