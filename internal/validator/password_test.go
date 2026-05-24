package validator

import (
	"errors"
	"strings"
	"testing"

	apperrors "echobackend/internal/errors"
)

func TestValidatePasswordStrength(t *testing.T) {
	tests := []struct {
		name    string
		pwd     string
		wantErr error
	}{
		{"valid mix", "GoodPass1!", nil},
		{"valid with symbol", "Aa1@aaaa", nil},
		{"too short", "Aa1!", apperrors.ErrPasswordTooShort},
		{"too long", strings.Repeat("Aa1!aaaa", 17), apperrors.ErrPasswordTooLong},
		{"missing upper", "lowerpass1!", apperrors.ErrPasswordNoUpper},
		{"missing lower", "UPPERPASS1!", apperrors.ErrPasswordNoLower},
		{"missing digit", "NoDigits!!!", apperrors.ErrPasswordNoDigit},
		{"missing special", "NoSpecial1A", apperrors.ErrPasswordNoSpecial},
		// 8 chars minimum boundary
		{"exactly 8 chars valid", "Aa1!aaaa", nil},
		{"7 chars too short", "Aa1!aaa", apperrors.ErrPasswordTooShort},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePasswordStrength(tt.pwd)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("ValidatePasswordStrength(%q) = %v, want %v", tt.pwd, err, tt.wantErr)
			}
		})
	}
}
