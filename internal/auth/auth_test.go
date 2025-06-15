package auth

import (
	"net/http"
	"testing"
	"time"
)

func TestGenerateAndValidateToken(t *testing.T) {
	// Create a test config with a short expiry for testing
	config := Config{
		Secret:      "test-secret-key",
		TokenExpiry: 1 * time.Hour,
	}

	// Test email
	email := "test@example.com"

	// Generate token
	token, err := GenerateToken(email, config)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	// Validate token
	validatedEmail, err := ValidateToken(token, config)
	if err != nil {
		t.Fatalf("Failed to validate token: %v", err)
	}

	// Check if the validated email matches the original
	if validatedEmail != email {
		t.Errorf("Validated email does not match original. Got %s, expected %s", validatedEmail, email)
	}
}

func TestValidateTokenWithInvalidToken(t *testing.T) {
	config := Config{
		Secret:      "test-secret-key",
		TokenExpiry: 1 * time.Hour,
	}

	// Test with invalid token
	_, err := ValidateToken("invalid-token", config)
	if err == nil {
		t.Error("Expected error for invalid token, got nil")
	}
	if err != ErrInvalidToken {
		t.Errorf("Expected ErrInvalidToken, got %v", err)
	}
}

func TestValidateTokenWithExpiredToken(t *testing.T) {
	// Create a config with a very short expiry
	config := Config{
		Secret:      "test-secret-key",
		TokenExpiry: 1 * time.Nanosecond,
	}

	// Generate token that will expire immediately
	email := "test@example.com"
	token, err := GenerateToken(email, config)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	// Wait for token to expire
	time.Sleep(10 * time.Millisecond)

	// Validate expired token
	_, err = ValidateToken(token, config)
	if err == nil {
		t.Error("Expected error for expired token, got nil")
	}
	if err != ErrExpiredToken {
		t.Errorf("Expected ErrExpiredToken, got %v", err)
	}
}

func TestExtractTokenFromRequest(t *testing.T) {
	tests := []struct {
		name          string
		authHeader    string
		expectedToken string
		expectedError error
	}{
		{
			name:          "Valid token",
			authHeader:    "Token valid-token",
			expectedToken: "valid-token",
			expectedError: nil,
		},
		{
			name:          "Missing token",
			authHeader:    "",
			expectedToken: "",
			expectedError: ErrInvalidToken,
		},
		{
			name:          "Invalid format",
			authHeader:    "Bearer valid-token",
			expectedToken: "",
			expectedError: ErrInvalidToken,
		},
		{
			name:          "Missing token value",
			authHeader:    "Token ",
			expectedToken: "",
			expectedError: ErrInvalidToken,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/", nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}

			token, err := ExtractTokenFromRequest(req)
			if err != tt.expectedError {
				t.Errorf("Expected error %v, got %v", tt.expectedError, err)
			}
			if token != tt.expectedToken {
				t.Errorf("Expected token %q, got %q", tt.expectedToken, token)
			}
		})
	}
}

func TestHashAndVerifyPassword(t *testing.T) {
	password := "secure-password"

	// Hash the password
	hashedPassword, err := HashPassword(password)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	// Verify correct password
	err = VerifyPassword(hashedPassword, password)
	if err != nil {
		t.Errorf("Failed to verify correct password: %v", err)
	}

	// Verify incorrect password
	err = VerifyPassword(hashedPassword, "wrong-password")
	if err == nil {
		t.Error("Expected error for incorrect password, got nil")
	}
}

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()

	if config.Secret == "" {
		t.Error("Expected non-empty secret in default config")
	}

	if config.TokenExpiry != 24*time.Hour {
		t.Errorf("Expected token expiry of 24 hours, got %v", config.TokenExpiry)
	}
}
