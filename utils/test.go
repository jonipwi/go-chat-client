package utils

import (
	"testing"
)

func TestSanitizeInput(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"  Hello World  ", "hello world"},
		{"  TEST  ", "test"},
		{"  ", ""},
	}

	for _, test := range tests {
		result := SanitizeInput(test.input)
		if result != test.expected {
			t.Errorf("SanitizeInput(%q) = %q; expected %q", test.input, result, test.expected)
		}
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
		err := ValidateUsername(test.username)
		if (err == nil) != test.valid {
			t.Errorf("ValidateUsername(%q) = %v; expected valid=%v", test.username, err, test.valid)
		}
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
		{"", 5, ""},
	}

	for _, test := range tests {
		result := TruncateMessage(test.message, test.maxLen)
		if result != test.expected {
			t.Errorf("TruncateMessage(%q, %d) = %q; expected %q", test.message, test.maxLen, result, test.expected)
		}
	}
}

func TestGenerateRandomID(t *testing.T) {
	id1 := GenerateRandomID()
	id2 := GenerateRandomID()

	if id1 == id2 {
		t.Errorf("GenerateRandomID() produced duplicate IDs: %q and %q", id1, id2)
	}
}
