package handler

import (
	"echobackend/internal/service"
	"echobackend/pkg/response"
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



func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Register(c echo.Context) error {
	var req RegisterRequest
	if err := c.Bind(&req); err != nil {
		return response.BadRequest(c, "Invalid request format", err)
	}

	if err := c.Validate(req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			return c.JSON(http.StatusBadRequest, response.APIResponse{
				Success: false,
				Message: "Validation failed",
				Error:   validationErrors.Error(),
				Data:    validationErrors.Errors,
			})
		}
		return response.ValidationError(c, "Validation failed", err)
	}

	user, err := h.authService.Register(c.Request().Context(), req.Email, req.Username, req.Password)
	if err == service.ErrUserExists {
		return c.JSON(http.StatusConflict, response.APIResponse{
			Success: false,
			Message: "Registration failed",
			Error:   "Email or username already exists",
		})
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

func (h *AuthHandler) Login(c echo.Context) error {
	var loginReq LoginRequest
	if err := c.Bind(&loginReq); err != nil {
		return response.BadRequest(c, "Invalid request format", err)
	}

	if err := c.Validate(loginReq); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			return c.JSON(http.StatusBadRequest, response.APIResponse{
				Success: false,
				Message: "Validation failed",
				Error:   validationErrors.Error(),
				Data:    validationErrors.Errors,
			})
		}
		return response.ValidationError(c, "Validation failed", err)
	}

	token, user, err := h.authService.Login(c.Request().Context(), loginReq.Email, loginReq.Password)
	if err == service.ErrInvalidCredentials {
		return response.Unauthorized(c, "Invalid email or password")
	}
	if err != nil {
		return response.InternalServerError(c, "Login failed", err)
	}

	return response.Success(c, "Login successful", map[string]any{
		"access_token": token,
		"user": map[string]any{
			"id":       user.ID,
			"email":    user.Email,
			"username": user.Username,
		},
	})
}

func (h *AuthHandler) CheckUsername(c echo.Context) error {
	var req CheckUsernameRequest
	if err := c.Bind(&req); err != nil {
		return response.BadRequest(c, "Invalid request format", err)
	}

	if err := c.Validate(req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			return c.JSON(http.StatusBadRequest, response.APIResponse{
				Success: false,
				Message: "Validation failed",
				Error:   validationErrors.Error(),
				Data:    validationErrors.Errors,
			})
		}
		return response.ValidationError(c, "Validation failed", err)
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
