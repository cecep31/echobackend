package utils

import "testing"

func TestIsValidEmail(t *testing.T) {
	tests := []struct {
		name  string
		email string
		want  bool
	}{
		{
			name:  "Valid email",
			email: "test@example.com",
			want:  true,
		},
		{
			name:  "Invalid email - no @",
			email: "testexample.com",
			want:  false,
		},
		{
			name:  "Invalid email - no domain",
			email: "test@",
			want:  false,
		},
		{
			name:  "Invalid email - no username",
			email: "@example.com",
			want:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidEmail(tt.email); got != tt.want {
				t.Errorf("IsValidEmail() = %v, want %v", got, tt.want)
			}
		})
	}
}
