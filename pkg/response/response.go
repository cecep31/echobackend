package response

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// APIResponse represents the standard API response format
type APIResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
	Error   string `json:"error,omitempty"`
	Meta    any    `json:"meta,omitempty"`
}

// PaginationMeta represents pagination metadata
type PaginationMeta struct {
	TotalItems int `json:"total_items"`
	Offset     int `json:"offset"`
	Limit      int `json:"limit"`
	TotalPages int `json:"total_pages"`
}

// Success sends a successful response
func Success(c echo.Context, message string, data any) error {
	return c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// SuccessWithMeta sends a successful response with metadata
func SuccessWithMeta(c echo.Context, message string, data any, meta any) error {
	return c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Message: message,
		Data:    data,
		Meta:    meta,
	})
}

// Created sends a created response
func Created(c echo.Context, message string, data any) error {
	return c.JSON(http.StatusCreated, APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// BadRequest sends a bad request error response
func BadRequest(c echo.Context, message string, err error) error {
	errorMsg := ""
	if err != nil {
		errorMsg = err.Error()
	}
	return c.JSON(http.StatusBadRequest, APIResponse{
		Success: false,
		Message: message,
		Error:   errorMsg,
	})
}

// Unauthorized sends an unauthorized error response
func Unauthorized(c echo.Context, message string) error {
	return c.JSON(http.StatusUnauthorized, APIResponse{
		Success: false,
		Message: message,
		Error:   "Unauthorized access",
	})
}

// Forbidden sends a forbidden error response
func Forbidden(c echo.Context, message string) error {
	return c.JSON(http.StatusForbidden, APIResponse{
		Success: false,
		Message: message,
		Error:   "Access forbidden",
	})
}

// NotFound sends a not found error response
func NotFound(c echo.Context, message string, err error) error {
	errorMsg := "Resource not found"
	if err != nil {
		errorMsg = err.Error()
	}
	return c.JSON(http.StatusNotFound, APIResponse{
		Success: false,
		Message: message,
		Error:   errorMsg,
	})
}

// InternalServerError sends an internal server error response
func InternalServerError(c echo.Context, message string, err error) error {
	errorMsg := ""
	if err != nil {
		errorMsg = err.Error()
	}
	return c.JSON(http.StatusInternalServerError, APIResponse{
		Success: false,
		Message: message,
		Error:   errorMsg,
	})
}

// ValidationError sends a validation error response
func ValidationError(c echo.Context, message string, err error) error {
	errorMsg := ""
	if err != nil {
		errorMsg = err.Error()
	}
	return c.JSON(http.StatusUnprocessableEntity, APIResponse{
		Success: false,
		Message: message,
		Error:   errorMsg,
	})
}

// CalculatePaginationMeta calculates pagination metadata
func CalculatePaginationMeta(totalItems int64, offset, limit int) PaginationMeta {
	totalPages := int(totalItems) / limit
	if int(totalItems)%limit > 0 {
		totalPages++
	}

	return PaginationMeta{
		TotalItems: int(totalItems),
		Offset:     offset,
		Limit:      limit,
		TotalPages: totalPages,
	}
}
