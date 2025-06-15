package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"git.homelab.lan/denga/go-real-world-example/internal/auth"
)

func TestAuth(t *testing.T) {
	// Create auth config
	config := auth.Config{
		Secret:      "test-secret-key",
		TokenExpiry: 1 * time.Hour,
	}

	// Create a test handler that checks if the user email is in the context
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		email, ok := GetUserEmail(r)
		if !ok {
			t.Error("Expected user email in context, but it was not found")
		}
		if email != "test@example.com" {
			t.Errorf("Expected email 'test@example.com', got '%s'", email)
		}
		w.WriteHeader(http.StatusOK)
	})

	// Create the middleware
	middleware := Auth(config)

	// Create a test server with the middleware
	ts := httptest.NewServer(middleware(testHandler))
	defer ts.Close()

	// Generate a valid token
	token, err := auth.GenerateToken("test@example.com", config)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	// Create a request with the token
	req, err := http.NewRequest("GET", ts.URL, nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Authorization", "Token "+token)

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// Check the response
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}
}

func TestAuthWithInvalidToken(t *testing.T) {
	// Create auth config
	config := auth.Config{
		Secret:      "test-secret-key",
		TokenExpiry: 1 * time.Hour,
	}

	// Create a test handler that should not be called
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("Handler should not be called with invalid token")
		w.WriteHeader(http.StatusOK)
	})

	// Create the middleware
	middleware := Auth(config)

	// Create a test server with the middleware
	ts := httptest.NewServer(middleware(testHandler))
	defer ts.Close()

	// Create a request with an invalid token
	req, err := http.NewRequest("GET", ts.URL, nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Authorization", "Token invalid-token")

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// Check the response
	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("Expected status code %d, got %d", http.StatusUnauthorized, resp.StatusCode)
	}
}

func TestAuthWithMissingToken(t *testing.T) {
	// Create auth config
	config := auth.Config{
		Secret:      "test-secret-key",
		TokenExpiry: 1 * time.Hour,
	}

	// Create a test handler that should not be called
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("Handler should not be called with missing token")
		w.WriteHeader(http.StatusOK)
	})

	// Create the middleware
	middleware := Auth(config)

	// Create a test server with the middleware
	ts := httptest.NewServer(middleware(testHandler))
	defer ts.Close()

	// Create a request without a token
	req, err := http.NewRequest("GET", ts.URL, nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// Check the response
	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("Expected status code %d, got %d", http.StatusUnauthorized, resp.StatusCode)
	}
}

func TestAuthWithPublicEndpoints(t *testing.T) {
	// Create auth config
	config := auth.Config{
		Secret:      "test-secret-key",
		TokenExpiry: 1 * time.Hour,
	}

	// Create a test handler that should be called for public endpoints
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Create the middleware
	middleware := Auth(config)

	// Test public endpoints
	publicEndpoints := []struct {
		method string
		path   string
	}{
		{"POST", "/api/users"},
		{"POST", "/api/users/login"},
		{"GET", "/api/tags"},
		{"GET", "/api/articles"},
		{"GET", "/openapi.yml"},
	}

	for _, endpoint := range publicEndpoints {
		t.Run(endpoint.method+" "+endpoint.path, func(t *testing.T) {
			// Create a test server with the middleware
			ts := httptest.NewServer(middleware(testHandler))
			defer ts.Close()

			// Create a request for the public endpoint
			req, err := http.NewRequest(endpoint.method, ts.URL+endpoint.path, nil)
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}

			// Send the request
			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				t.Fatalf("Failed to send request: %v", err)
			}
			defer resp.Body.Close()

			// Check the response
			if resp.StatusCode != http.StatusOK {
				t.Errorf("Expected status code %d for public endpoint %s %s, got %d",
					http.StatusOK, endpoint.method, endpoint.path, resp.StatusCode)
			}
		})
	}
}

func TestGetUserEmail(t *testing.T) {
	// Create a request with a context containing a user email
	req, _ := http.NewRequest("GET", "/", nil)
	ctx := context.WithValue(req.Context(), UserEmailKey, "test@example.com")
	req = req.WithContext(ctx)

	// Test GetUserEmail
	email, ok := GetUserEmail(req)
	if !ok {
		t.Error("Expected GetUserEmail to return true, got false")
	}
	if email != "test@example.com" {
		t.Errorf("Expected email 'test@example.com', got '%s'", email)
	}

	// Test with missing email
	req, _ = http.NewRequest("GET", "/", nil)
	email, ok = GetUserEmail(req)
	if ok {
		t.Error("Expected GetUserEmail to return false for missing email, got true")
	}
	if email != "" {
		t.Errorf("Expected empty email for missing email, got '%s'", email)
	}
}
