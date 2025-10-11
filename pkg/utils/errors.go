package utils

import (
	"errors"
	"fmt"
)

// AppError represents a structured application error
type AppError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Err     error  `json:"-"` // Not serialized to JSON
}

// Error returns the error message
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

// Unwrap returns the underlying error
func (e *AppError) Unwrap() error {
	return e.Err
}

// NewAppError creates a new AppError
func NewAppError(code, message string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// Error codes
const (
	ErrCodeValidationFailed = "VALIDATION_FAILED"
	ErrCodeNotFound         = "NOT_FOUND"
	ErrCodeConflict         = "CONFLICT"
	ErrCodeInternal         = "INTERNAL_ERROR"
	ErrCodeUnauthorized     = "UNAUTHORIZED"
	ErrCodeForbidden        = "FORBIDDEN"
)

// IsAppError checks if an error is an AppError
func IsAppError(err error) bool {
	var appErr *AppError
	return errors.As(err, &appErr)
}

// AsAppError tries to convert an error to AppError
func AsAppError(err error) (*AppError, bool) {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr, true
	}
	return nil, false
}

// ValidationError creates a validation error
func ValidationError(message string, err error) *AppError {
	return NewAppError(ErrCodeValidationFailed, message, err)
}

// NotFoundError creates a not found error
func NotFoundError(message string, err error) *AppError {
	return NewAppError(ErrCodeNotFound, message, err)
}

// ConflictError creates a conflict error
func ConflictError(message string, err error) *AppError {
	return NewAppError(ErrCodeConflict, message, err)
}

// InternalError creates an internal error
func InternalError(message string, err error) *AppError {
	return NewAppError(ErrCodeInternal, message, err)
}

// UnauthorizedError creates an unauthorized error
func UnauthorizedError(message string, err error) *AppError {
	return NewAppError(ErrCodeUnauthorized, message, err)
}

// ForbiddenError creates a forbidden error
func ForbiddenError(message string, err error) *AppError {
	return NewAppError(ErrCodeForbidden, message, err)
}