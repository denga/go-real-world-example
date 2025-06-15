package util

import (
	"testing"
)

func TestGenerateSlug(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Simple title",
			input:    "Hello World",
			expected: "hello-world",
		},
		{
			name:     "Title with special characters",
			input:    "Hello, World! How are you?",
			expected: "hello-world-how-are-you",
		},
		{
			name:     "Title with numbers",
			input:    "Article 123",
			expected: "article-123",
		},
		{
			name:     "Title with extra spaces",
			input:    "  Hello   World  ",
			expected: "hello-world",
		},
		{
			name:     "Empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "Only special characters",
			input:    "!@#$%^&*()",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GenerateSlug(tt.input)
			if result != tt.expected {
				t.Errorf("GenerateSlug(%q) = %q, expected %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestGenerateUniqueSlug(t *testing.T) {
	existingSlugs := map[string]bool{
		"hello-world":   true,
		"hello-world-1": true,
	}

	tests := []struct {
		name     string
		input    string
		exists   func(string) bool
		expected string
	}{
		{
			name:  "Non-existing slug",
			input: "New Article",
			exists: func(s string) bool {
				return existingSlugs[s]
			},
			expected: "new-article",
		},
		{
			name:  "Existing slug",
			input: "Hello World",
			exists: func(s string) bool {
				return existingSlugs[s]
			},
			expected: "hello-world-2",
		},
		{
			name:  "Existing slug with counter",
			input: "Hello World 1",
			exists: func(s string) bool {
				return existingSlugs[s]
			},
			expected: "hello-world-1-1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GenerateUniqueSlug(tt.input, tt.exists)
			if result != tt.expected {
				t.Errorf("GenerateUniqueSlug(%q) = %q, expected %q", tt.input, result, tt.expected)
			}
		})
	}
}