package validator

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"

	"github.com/go-playground/validator/v10"
)

var messages = map[string]string{
	"required": "%s is required",
	"email":    "%s must be a valid email address",
	"min":      "%s must be at least %s characters long",
	"max":      "%s must not exceed %s characters",
	"oneof":    "%s must be one of [%s]",
	"default":  "%s failed validation for tag %s",
}

type CustomValidator struct {
	validator *validator.Validate
}

// ValidationErrors represents multiple validation errors
type ValidationErrors struct {
	Errors []ValidationError `json:"errors"`
}

// ValidationError represents a single validation error
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Value   any    `json:"value,omitempty"`
	Tag     string `json:"tag,omitempty"`
}

func (v ValidationErrors) Error() string {
	if len(v.Errors) == 0 {
		return "validation failed"
	}
	return v.Errors[0].Message
}

func NewValidator() *CustomValidator {
	v := validator.New()
	return &CustomValidator{validator: v}
}

func (cv *CustomValidator) Validate(i any) error {
	if err := cv.validator.Struct(i); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errors := ValidationErrors{
				Errors: make([]ValidationError, len(validationErrors)),
			}
			for i, e := range validationErrors {
				errors.Errors[i] = ValidationError{
					Field:   e.Field(),
					Message: getErrorMessage(e),
					Value:   e.Value(),
					Tag:     e.Tag(),
				}
			}
			return errors
		}
		return err
	}
	return nil
}

func getErrorMessage(e validator.FieldError) string {
	fieldName := toReadableFieldName(e.Field())
	switch e.Tag() {
	case "required":
		return fmt.Sprintf(messages["required"], fieldName)
	case "email":
		return fmt.Sprintf(messages["email"], fieldName)
	case "min":
		return fmt.Sprintf(messages["min"], fieldName, e.Param())
	case "max":
		return fmt.Sprintf(messages["max"], fieldName, e.Param())
	case "oneof":
		return fmt.Sprintf(messages["oneof"], fieldName, e.Param())
	default:
		return fmt.Sprintf(messages["default"], fieldName, e.Tag())
	}
}

func toReadableFieldName(field string) string {
	var result strings.Builder
	for i, r := range field {
		if i > 0 && unicode.IsUpper(r) {
			result.WriteRune(' ')
		}
		if i == 0 {
			result.WriteRune(unicode.ToUpper(r))
		} else {
			result.WriteRune(unicode.ToLower(r))
		}
	}
	return result.String()
}

// IsValidUUID validates if a string is a valid UUID v7 format
func IsValidUUID(uuid string) bool {
	if uuid == "" {
		return false
	}
	// UUID v7 pattern: xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
	pattern := `^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`
	matched, _ := regexp.MatchString(pattern, uuid)
	return matched
}

// ValidatePagination validates pagination parameters
func ValidatePagination(limit, offset int) error {
	if limit <= 0 {
		return fmt.Errorf("limit must be greater than 0")
	}
	if limit > 100 {
		return fmt.Errorf("limit must not exceed 100")
	}
	if offset < 0 {
		return fmt.Errorf("offset must be non-negative")
	}
	return nil
}

// ValidatePostLikeInput validates input for post like operations
func ValidatePostLikeInput(postID, userID string) error {
	if !IsValidUUID(postID) {
		return fmt.Errorf("invalid post ID format")
	}
	if !IsValidUUID(userID) {
		return fmt.Errorf("invalid user ID format")
	}
	return nil
}
