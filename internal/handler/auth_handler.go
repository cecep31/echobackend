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
	Password string `json:"password" validate:"required,min=6"`
}

// Response represents a standard API response
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Errors  interface{} `json:"errors,omitempty"`
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

	user, err := h.authService.Register(req.Email, req.Password)
	if err == service.ErrUserExists {
		return c.JSON(http.StatusConflict, Response{
			Success: false,
			Message: "Registration failed",
			Errors:  []string{"Email already registered"},
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
		Data: map[string]interface{}{
			"id":    user.ID,
			"email": user.Email,
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

	token, user, err := h.authService.Login(loginReq.Email, loginReq.Password)
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
		Data: map[string]interface{}{
			"access_token": token,
			"user": map[string]interface{}{
				"id":    user.ID,
				"email": user.Email,
			},
		},
	})
}
