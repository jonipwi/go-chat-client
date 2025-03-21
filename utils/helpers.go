package utils

import (
	"fmt"
	"strings"
	"time"
)

// SanitizeInput removes leading/trailing whitespace and converts to lowercase
func SanitizeInput(input string) string {
	return strings.TrimSpace(strings.ToLower(input))
}

// ValidateUsername checks if a username is valid
func ValidateUsername(username string) error {
	// Remove leading/trailing whitespace
	username = strings.TrimSpace(username)

	// Check length
	if len(username) < 3 {
		return fmt.Errorf("username must be at least 3 characters long")
	}

	if len(username) > 20 {
		return fmt.Errorf("username cannot be longer than 20 characters")
	}

	// Check for invalid characters
	for _, char := range username {
		if !((char >= 'a' && char <= 'z') ||
			(char >= 'A' && char <= 'Z') ||
			(char >= '0' && char <= '9') ||
			char == '_' || char == '-') {
			return fmt.Errorf("username can only contain letters, numbers, underscores, and hyphens")
		}
	}

	return nil
}

// TruncateMessage limits the length of a message
func TruncateMessage(message string, maxLength int) string {
	if len(message) <= maxLength {
		return message
	}
	return message[:maxLength] + "..."
}

// EscapeSpecialChars escapes special characters that might cause issues in messaging
func EscapeSpecialChars(input string) string {
	// Replace potentially problematic characters
	replacements := map[string]string{
		"&":  "&amp;",
		"<":  "&lt;",
		">":  "&gt;",
		"\"": "&quot;",
		"'":  "&#39;",
	}

	for orig, replacement := range replacements {
		input = strings.ReplaceAll(input, orig, replacement)
	}

	return input
}

// GenerateRandomID creates a simple unique identifier
func GenerateRandomID() string {
	// This is a naive implementation. In a real-world scenario,
	// you'd use a more robust method like UUID
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

// FormatTimestamp converts a timestamp to a human-readable format
func FormatTimestamp(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}
