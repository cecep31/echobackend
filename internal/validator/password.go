package validator

import (
	"unicode"

	apperrors "echobackend/internal/errors"
)

func ValidatePasswordStrength(password string) error {
	if len(password) < 8 {
		return apperrors.ErrPasswordTooShort
	}
	if len(password) > 128 {
		return apperrors.ErrPasswordTooLong
	}

	hasUpper := false
	hasLower := false
	hasDigit := false
	hasSpecial := false

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasDigit = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	if !hasUpper {
		return apperrors.ErrPasswordNoUpper
	}
	if !hasLower {
		return apperrors.ErrPasswordNoLower
	}
	if !hasDigit {
		return apperrors.ErrPasswordNoDigit
	}
	if !hasSpecial {
		return apperrors.ErrPasswordNoSpecial
	}

	return nil
}