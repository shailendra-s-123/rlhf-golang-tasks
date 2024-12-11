package crypto

import (
	"testing"
)

func TestHash(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"", "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"},
		{"hello", "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824"},
		{"world", "da39a3ee5e6b4b0d3255bfef95601890afd80709"},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := Hash([]byte(tt.input))
			if result != tt.expected {
				t.Errorf("Hash failed. Expected %s, got %s", tt.expected, result)
			}
		})
	}
}