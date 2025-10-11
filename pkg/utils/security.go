package utils

import (
	"regexp"
	"strings"
)

// SanitizeInput sanitizes user input to prevent injection attacks
func SanitizeInput(input string) string {
	// Remove potentially dangerous characters/sequences
	re := regexp.MustCompile(`<script[^>]*>.*?</script>`)
	sanitized := re.ReplaceAllString(input, "")
	
	// Remove other potentially dangerous patterns
	re = regexp.MustCompile(`(?i)<iframe[^>]*>.*?</iframe>`)
	sanitized = re.ReplaceAllString(sanitized, "")
	
	re = regexp.MustCompile(`(?i)<object[^>]*>.*?</object>`)
	sanitized = re.ReplaceAllString(sanitized, "")
	
	re = regexp.MustCompile(`(?i)<embed[^>]*>.*?</embed>`)
	sanitized = re.ReplaceAllString(sanitized, "")
	
	// Remove JavaScript event handlers and other potentially dangerous attributes
	re = regexp.MustCompile(`(?i)(on\w+\s*=)`)
	sanitized = re.ReplaceAllString(sanitized, "")
	
	// Additional sanitization as needed
	sanitized = strings.TrimSpace(sanitized)
	return sanitized
}

// IsValidEmail validates email format
func IsValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

// IsValidUsername validates username format (alphanumeric, underscore, hyphen, 3-30 chars)
func IsValidUsername(username string) bool {
	usernameRegex := regexp.MustCompile(`^[a-zA-Z0-9_-]{3,30}$`)
	return usernameRegex.MatchString(username)
}

// IsValidPassword checks if password meets basic security requirements
func IsValidPassword(password string) bool {
	// At least one uppercase, one lowercase, one digit, one special char, at least 8 chars
	passwordRegex := regexp.MustCompile(`^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[@$!%*?&])[A-Za-z\d@$!%*?&]{8,}$`)
	return passwordRegex.MatchString(password)
}

// SanitizeForSQL sanitizes input for SQL queries
func SanitizeForSQL(input string) string {
	// Escape SQL special characters
	input = strings.ReplaceAll(input, "'", "''") // Escape single quotes
	input = strings.ReplaceAll(input, "\\", "\\\\") // Escape backslashes
	return input
}