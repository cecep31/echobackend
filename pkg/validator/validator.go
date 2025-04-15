package validator

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/go-playground/validator/v10"
)

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
			var errors ValidationErrors
			for _, e := range validationErrors {
				errors.Errors = append(errors.Errors, ValidationError{
					Field:   e.Field(),
					Message: getErrorMessage(e),
					Value:   e.Value(),
					Tag:     e.Tag(),
				})
			}
			return errors
		}
		return err
	}
	return nil
}

func getErrorMessage(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return getMessage("required", e.Field())
	case "email":
		return getMessage("email", e.Field())
	case "min":
		return getMessage("min", e.Field(), e.Param())
	case "max":
		return getMessage("max", e.Field(), e.Param())
	case "oneof":
		return getMessage("oneof", e.Field(), e.Param())
	default:
		return getMessage("default", e.Field(), e.Tag())
	}
}

func getMessage(tag string, params ...string) string {
	messages := map[string]string{
		"required": "%s is required",
		"email":    "%s must be a valid email address",
		"min":      "%s must be at least %s characters long",
		"max":      "%s must not exceed %s characters",
		"oneof":    "%s must be one of [%s]",
		"default":  "%s failed validation for tag %s",
	}

	msg, ok := messages[tag]
	if !ok {
		msg = messages["default"]
	}

	switch len(params) {
	case 1:
		return sprintf(msg, params[0])
	case 2:
		return sprintf(msg, params[0], params[1])
	default:
		return sprintf(messages["default"], params[0], "unknown")
	}
}

func sprintf(format string, args ...any) string {
	// Convert field names from camelCase/PascalCase to user-friendly format
	if len(args) > 0 {
		if fieldName, ok := args[0].(string); ok {
			args[0] = toReadableFieldName(fieldName)
		}
	}
	return fmt.Sprintf(format, args...)
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
