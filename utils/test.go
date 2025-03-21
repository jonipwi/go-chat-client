package utils

import (
	"fmt"
	"testing"
	"time"
)

func TestSanitizeInput(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"  Hello World  ", "hello world"},
		{"  TEST  ", "test"},
		{"  ", ""}, // edge case for empty string
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("SanitizeInput(%q)", test.input), func(t *testing.T) {
			result := SanitizeInput(test.input)
			if result != test.expected {
				t.Errorf("SanitizeInput(%q) = %q; expected %q", test.input, result, test.expected)
			}
		})
	}
}

func TestValidateUsername(t *testing.T) {
	tests := []struct {
		username string
		valid    bool
	}{
		{"user", true},
		{"u", false},                     // Too short
		{"thisusernameistoolong", false}, // Too long
		{"user@name", false},             // Invalid character
		{"user_name", true},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("ValidateUsername(%q)", test.username), func(t *testing.T) {
			err := ValidateUsername(test.username)
			if (err == nil) != test.valid {
				t.Errorf("ValidateUsername(%q) = %v; expected valid=%v", test.username, err, test.valid)
			}
		})
	}
}

func TestTruncateMessage(t *testing.T) {
	tests := []struct {
		message  string
		maxLen   int
		expected string
	}{
		{"Hello World", 5, "Hello..."},
		{"Short", 10, "Short"},
		{"", 5, ""}, // edge case for empty message
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("TruncateMessage(%q, %d)", test.message, test.maxLen), func(t *testing.T) {
			result := TruncateMessage(test.message, test.maxLen)
			if result != test.expected {
				t.Errorf("TruncateMessage(%q, %d) = %q; expected %q", test.message, test.maxLen, result, test.expected)
			}
		})
	}
}

func TestGenerateRandomID(t *testing.T) {
	id1 := GenerateRandomID()
	id2 := GenerateRandomID()

	if id1 == id2 {
		t.Errorf("GenerateRandomID() produced duplicate IDs: %q and %q", id1, id2)
	}
}

func TestEscapeSpecialChars(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"<div>Hello & Welcome</div>", "&lt;div&gt;Hello &amp; Welcome&lt;/div&gt;"},
		{"\"Quoted text\"", "&quot;Quoted text&quot;"},
		{"It's a test", "It&#39;s a test"},
		{"<script>alert('XSS')</script>", "&lt;script&gt;alert(&#39;XSS&#39;)&lt;/script&gt;"},
		{"Hello World", "Hello World"},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("EscapeSpecialChars(%q)", test.input), func(t *testing.T) {
			result := EscapeSpecialChars(test.input)
			if result != test.expected {
				t.Errorf("EscapeSpecialChars(%q) = %q; expected %q", test.input, result, test.expected)
			}
		})
	}
}

func TestLoggerOutput(t *testing.T) {
	// Test that Logger outputs messages to stdout without error
	t.Run("Logger outputs without error", func(t *testing.T) {
		Logger.Println("Test log message")

		// Manually check output (this is usually done with log capture tools, but simple verification can be done)
		// We can capture the output here or verify if the log works.
		// For now, this test will just verify if no panic occurs when logging.
	})
}

func TestFormatTimestamp(t *testing.T) {
	tests := []struct {
		timestamp time.Time
		expected  string
	}{
		{time.Date(2023, 3, 21, 15, 4, 5, 0, time.UTC), "2023-03-21 15:04:05"},
		{time.Time{}, "0001-01-01 00:00:00"}, // Empty time edge case
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("FormatTimestamp(%v)", test.timestamp), func(t *testing.T) {
			result := FormatTimestamp(test.timestamp)
			if result != test.expected {
				t.Errorf("FormatTimestamp(%v) = %v; expected %v", test.timestamp, result, test.expected)
			}
		})
	}
}
