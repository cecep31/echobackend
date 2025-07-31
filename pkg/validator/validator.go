package validator

import (
	"fmt"
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
