package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"git.homelab.lan/denga/go-real-world-example/api"
	"git.homelab.lan/denga/go-real-world-example/internal/auth"
	"git.homelab.lan/denga/go-real-world-example/internal/db"
	"git.homelab.lan/denga/go-real-world-example/internal/middleware"
)

// setupTestHandler creates a new Handler with a test database and auth config
func setupTestHandler() (*Handler, *db.InMemoryDB) {
	// Create a test database
	testDB := db.NewInMemoryDB()

	// Create auth config
	authConfig := auth.Config{
		Secret:      "test-secret-key",
		TokenExpiry: 1 * time.Hour,
	}

	// Create handler
	handler := NewHandler(testDB, authConfig)

	return handler, testDB
}

// setupTestUser creates a test user in the database and returns the user and token
func setupTestUser(testDB *db.InMemoryDB, authConfig auth.Config) (api.User, string) {
	user := api.User{
		Username: "testuser",
		Email:    "test@example.com",
		Bio:      "Test bio",
		Image:    "test-image.jpg",
	}

	// Generate token
	token, _ := auth.GenerateToken(user.Email, authConfig)
	user.Token = token

	// Create user in database
	testDB.CreateUser(user, "password123")

	return user, token
}

// addAuthHeader adds an Authorization header with the given token to the request
func addAuthHeader(req *http.Request, token string) {
	req.Header.Set("Authorization", "Token "+token)
}

// addUserToContext adds a user email to the request context
func addUserToContext(req *http.Request, email string) *http.Request {
	ctx := req.Context()
	ctx = context.WithValue(ctx, middleware.UserEmailKey, email)
	return req.WithContext(ctx)
}

func TestCreateUser(t *testing.T) {
	handler, _ := setupTestHandler()

	// Create request body
	reqBody := api.NewUserRequest{
		User: api.NewUser{
			Username: "newuser",
			Email:    "new@example.com",
			Password: "password123",
		},
	}
	body, _ := json.Marshal(reqBody)

	// Create request
	req := httptest.NewRequest("POST", "/api/users", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	// Create response recorder
	rr := httptest.NewRecorder()

	// Call handler
	handler.CreateUser(rr, req)

	// Check response
	if rr.Code != http.StatusCreated {
		t.Errorf("Expected status code %d, got %d", http.StatusCreated, rr.Code)
	}

	// Parse response
	var resp api.UserResponse
	err := json.Unmarshal(rr.Body.Bytes(), &resp)
	if err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	// Check user data
	if resp.User.Username != reqBody.User.Username {
		t.Errorf("Expected username %s, got %s", reqBody.User.Username, resp.User.Username)
	}
	if resp.User.Email != reqBody.User.Email {
		t.Errorf("Expected email %s, got %s", reqBody.User.Email, resp.User.Email)
	}
	if resp.User.Token == "" {
		t.Error("Expected non-empty token")
	}
}

func TestLogin(t *testing.T) {
	handler, testDB := setupTestHandler()

	// Create a test user
	user := api.User{
		Username: "testuser",
		Email:    "test@example.com",
		Bio:      "Test bio",
		Image:    "test-image.jpg",
	}
	testDB.CreateUser(user, "password123")

	// Create request body
	reqBody := api.LoginUserRequest{
		User: api.LoginUser{
			Email:    user.Email,
			Password: "password123",
		},
	}
	body, _ := json.Marshal(reqBody)

	// Create request
	req := httptest.NewRequest("POST", "/api/users/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	// Create response recorder
	rr := httptest.NewRecorder()

	// Call handler
	handler.Login(rr, req)

	// Check response
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, rr.Code)
	}

	// Parse response
	var resp api.UserResponse
	err := json.Unmarshal(rr.Body.Bytes(), &resp)
	if err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	// Check user data
	if resp.User.Username != user.Username {
		t.Errorf("Expected username %s, got %s", user.Username, resp.User.Username)
	}
	if resp.User.Email != user.Email {
		t.Errorf("Expected email %s, got %s", user.Email, resp.User.Email)
	}
	if resp.User.Token == "" {
		t.Error("Expected non-empty token")
	}
}

func TestLoginWithInvalidCredentials(t *testing.T) {
	handler, testDB := setupTestHandler()

	// Create a test user
	user := api.User{
		Username: "testuser",
		Email:    "test@example.com",
		Bio:      "Test bio",
		Image:    "test-image.jpg",
	}
	testDB.CreateUser(user, "password123")

	// Create request body with wrong password
	reqBody := api.LoginUserRequest{
		User: api.LoginUser{
			Email:    user.Email,
			Password: "wrongpassword",
		},
	}
	body, _ := json.Marshal(reqBody)

	// Create request
	req := httptest.NewRequest("POST", "/api/users/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	// Create response recorder
	rr := httptest.NewRecorder()

	// Call handler
	handler.Login(rr, req)

	// Check response
	if rr.Code != http.StatusUnauthorized {
		t.Errorf("Expected status code %d, got %d", http.StatusUnauthorized, rr.Code)
	}
}

func TestGetCurrentUser(t *testing.T) {
	handler, testDB := setupTestHandler()

	// Create a test user
	user, token := setupTestUser(testDB, handler.AuthConfig)

	// Create request
	req := httptest.NewRequest("GET", "/api/user", nil)
	addAuthHeader(req, token)

	// Add user to context (simulating middleware)
	req = addUserToContext(req, user.Email)

	// Create response recorder
	rr := httptest.NewRecorder()

	// Call handler
	handler.GetCurrentUser(rr, req)

	// Check response
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, rr.Code)
	}

	// Parse response
	var resp api.UserResponse
	err := json.Unmarshal(rr.Body.Bytes(), &resp)
	if err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	// Check user data
	if resp.User.Username != user.Username {
		t.Errorf("Expected username %s, got %s", user.Username, resp.User.Username)
	}
	if resp.User.Email != user.Email {
		t.Errorf("Expected email %s, got %s", user.Email, resp.User.Email)
	}
	if resp.User.Token == "" {
		t.Error("Expected non-empty token")
	}
}

func TestUpdateCurrentUser(t *testing.T) {
	handler, testDB := setupTestHandler()

	// Create a test user
	user, token := setupTestUser(testDB, handler.AuthConfig)

	// Create request body
	newBio := "Updated bio"
	newImage := "updated-image.jpg"
	reqBody := api.UpdateUserRequest{
		User: api.UpdateUser{
			Bio:   &newBio,
			Image: &newImage,
		},
	}
	body, _ := json.Marshal(reqBody)

	// Create request
	req := httptest.NewRequest("PUT", "/api/user", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	addAuthHeader(req, token)

	// Add user to context (simulating middleware)
	req = addUserToContext(req, user.Email)

	// Create response recorder
	rr := httptest.NewRecorder()

	// Call handler
	handler.UpdateCurrentUser(rr, req)

	// Check response
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, rr.Code)
	}

	// Parse response
	var resp api.UserResponse
	err := json.Unmarshal(rr.Body.Bytes(), &resp)
	if err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	// Check updated user data
	if resp.User.Bio != newBio {
		t.Errorf("Expected bio %s, got %s", newBio, resp.User.Bio)
	}
	if resp.User.Image != newImage {
		t.Errorf("Expected image %s, got %s", newImage, resp.User.Image)
	}
}

func TestGetTags(t *testing.T) {
	handler, testDB := setupTestHandler()

	// Create a test user
	user, _ := setupTestUser(testDB, handler.AuthConfig)

	// Create author profile
	author := api.Profile{
		Username:  user.Username,
		Bio:       user.Bio,
		Image:     user.Image,
		Following: false,
	}

	// Create test articles with tags
	article1 := api.Article{
		Title:       "Test Article 1",
		Description: "Test description 1",
		Body:        "Test body 1",
		TagList:     []string{"tag1", "tag2"},
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Favorited:   false,
		Author:      author,
		Slug:        "test-article-1",
	}
	testDB.CreateArticle(article1)

	article2 := api.Article{
		Title:       "Test Article 2",
		Description: "Test description 2",
		Body:        "Test body 2",
		TagList:     []string{"tag2", "tag3"},
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Favorited:   false,
		Author:      author,
		Slug:        "test-article-2",
	}
	testDB.CreateArticle(article2)

	// Create request
	req := httptest.NewRequest("GET", "/api/tags", nil)

	// Create response recorder
	rr := httptest.NewRecorder()

	// Call handler
	handler.GetTags(rr, req)

	// Check response
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, rr.Code)
	}

	// Parse response
	var resp api.TagsResponse
	err := json.Unmarshal(rr.Body.Bytes(), &resp)
	if err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	// Check tags
	if len(resp.Tags) != 3 {
		t.Errorf("Expected 3 tags, got %d", len(resp.Tags))
	}

	// Check that all expected tags are present
	expectedTags := map[string]bool{
		"tag1": true,
		"tag2": true,
		"tag3": true,
	}
	for _, tag := range resp.Tags {
		if !expectedTags[tag] {
			t.Errorf("Unexpected tag: %s", tag)
		}
		delete(expectedTags, tag)
	}
	if len(expectedTags) != 0 {
		t.Errorf("Missing tags: %v", expectedTags)
	}
}
