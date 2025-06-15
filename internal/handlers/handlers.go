package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/denga/go-real-world-example/api"
	"github.com/denga/go-real-world-example/internal/auth"
	"github.com/denga/go-real-world-example/internal/db"
	"github.com/denga/go-real-world-example/internal/middleware"
	"github.com/denga/go-real-world-example/internal/util"
)

// Handler implements the ServerInterface from the generated API code
type Handler struct {
	DB         *db.InMemoryDB
	AuthConfig auth.Config
}

// NewHandler creates a new Handler
func NewHandler(db *db.InMemoryDB, authConfig auth.Config) *Handler {
	return &Handler{
		DB:         db,
		AuthConfig: authConfig,
	}
}

// GetArticles returns a list of articles
func (h *Handler) GetArticles(w http.ResponseWriter, r *http.Request, params api.GetArticlesParams) {
	// Set default values for limit and offset
	limit := 20
	if params.Limit != nil {
		limit = int(*params.Limit)
	}

	offset := 0
	if params.Offset != nil {
		offset = int(*params.Offset)
	}

	// Extract filter parameters
	var tag, author, favorited string
	if params.Tag != nil {
		tag = *params.Tag
	}
	if params.Author != nil {
		author = *params.Author
	}
	if params.Favorited != nil {
		favorited = *params.Favorited
	}

	// Get articles from database
	articles, count, err := h.DB.ListArticles(tag, author, favorited, limit, offset)
	if err != nil {
		http.Error(w, "Error retrieving articles", http.StatusInternalServerError)
		return
	}

	// Prepare response
	response := api.MultipleArticlesResponse{
		ArticlesCount: count,
	}

	// Convert articles to response format
	for _, article := range articles {
		// Create an anonymous struct that matches the expected type
		responseArticle := struct {
			Author         api.Profile `json:"author"`
			CreatedAt      time.Time   `json:"createdAt"`
			Description    string      `json:"description"`
			Favorited      bool        `json:"favorited"`
			FavoritesCount int         `json:"favoritesCount"`
			Slug           string      `json:"slug"`
			TagList        []string    `json:"tagList"`
			Title          string      `json:"title"`
			UpdatedAt      time.Time   `json:"updatedAt"`
		}{
			Author:         article.Author,
			CreatedAt:      article.CreatedAt,
			Description:    article.Description,
			Favorited:      article.Favorited,
			FavoritesCount: article.FavoritesCount,
			Slug:           article.Slug,
			TagList:        article.TagList,
			Title:          article.Title,
			UpdatedAt:      article.UpdatedAt,
		}
		response.Articles = append(response.Articles, responseArticle)
	}

	// Write response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// CreateArticle creates a new article
func (h *Handler) CreateArticle(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var request api.NewArticleRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Get authenticated user from context
	email, ok := middleware.GetUserEmail(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get user from database
	user, err := h.DB.GetUserByEmail(email)
	if err != nil {
		if err == db.ErrNotFound {
			http.Error(w, "User not found", http.StatusNotFound)
		} else {
			http.Error(w, "Error retrieving user", http.StatusInternalServerError)
		}
		return
	}

	// Create author profile
	author := api.Profile{
		Username:  user.Username,
		Bio:       user.Bio,
		Image:     user.Image,
		Following: false, // User can't follow themselves
	}

	// Generate slug from title
	slug := util.GenerateUniqueSlug(request.Article.Title, func(s string) bool {
		article, err := h.DB.GetArticle(s)
		return err == nil && article != nil
	})

	// Create article
	article := api.Article{
		Title:       request.Article.Title,
		Description: request.Article.Description,
		Body:        request.Article.Body,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Author:      author,
		Slug:        slug,
		Favorited:   false,
	}

	// Handle TagList which is a pointer to a slice
	if request.Article.TagList != nil {
		article.TagList = *request.Article.TagList
	} else {
		article.TagList = []string{}
	}

	// Save article to database
	if err := h.DB.CreateArticle(article); err != nil {
		if err == db.ErrConflict {
			http.Error(w, "Article with this title already exists", http.StatusConflict)
		} else {
			http.Error(w, "Error creating article", http.StatusInternalServerError)
		}
		return
	}

	// Prepare response
	response := api.SingleArticleResponse{
		Article: article,
	}

	// Write response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// GetArticlesFeed returns articles from followed users
func (h *Handler) GetArticlesFeed(w http.ResponseWriter, r *http.Request, params api.GetArticlesFeedParams) {
	// Set default values for limit and offset
	limit := 20
	if params.Limit != nil {
		limit = int(*params.Limit)
	}

	offset := 0
	if params.Offset != nil {
		offset = int(*params.Offset)
	}

	// Get authenticated user from context
	email, ok := middleware.GetUserEmail(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get user from database
	user, err := h.DB.GetUserByEmail(email)
	if err != nil {
		if err == db.ErrNotFound {
			http.Error(w, "User not found", http.StatusNotFound)
		} else {
			http.Error(w, "Error retrieving user", http.StatusInternalServerError)
		}
		return
	}

	// Get articles from database
	articles, count, err := h.DB.GetArticlesFeed(user.Username, limit, offset)
	if err != nil {
		if err == db.ErrNotFound {
			http.Error(w, "User not found", http.StatusNotFound)
		} else {
			http.Error(w, "Error retrieving articles", http.StatusInternalServerError)
		}
		return
	}

	// Prepare response
	response := api.MultipleArticlesResponse{
		ArticlesCount: count,
	}

	// Convert articles to response format
	for _, article := range articles {
		// Create an anonymous struct that matches the expected type
		responseArticle := struct {
			Author         api.Profile `json:"author"`
			CreatedAt      time.Time   `json:"createdAt"`
			Description    string      `json:"description"`
			Favorited      bool        `json:"favorited"`
			FavoritesCount int         `json:"favoritesCount"`
			Slug           string      `json:"slug"`
			TagList        []string    `json:"tagList"`
			Title          string      `json:"title"`
			UpdatedAt      time.Time   `json:"updatedAt"`
		}{
			Author:         article.Author,
			CreatedAt:      article.CreatedAt,
			Description:    article.Description,
			Favorited:      article.Favorited,
			FavoritesCount: article.FavoritesCount,
			Slug:           article.Slug,
			TagList:        article.TagList,
			Title:          article.Title,
			UpdatedAt:      article.UpdatedAt,
		}
		response.Articles = append(response.Articles, responseArticle)
	}

	// Write response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetTags returns all tags
func (h *Handler) GetTags(w http.ResponseWriter, r *http.Request) {
	// Get tags from database
	tags := h.DB.GetTags()

	// Prepare response
	response := api.TagsResponse{
		Tags: tags,
	}

	// Write response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// CreateUser creates a new user
func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var request api.NewUserRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Generate token
	token, err := auth.GenerateToken(request.User.Email, h.AuthConfig)
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	// Create user
	user := api.User{
		Username: request.User.Username,
		Email:    request.User.Email,
		Bio:      "",
		Image:    "",
		Token:    token,
	}

	// Save user to database
	if err := h.DB.CreateUser(user, request.User.Password); err != nil {
		if err == db.ErrConflict {
			http.Error(w, "User with this email or username already exists", http.StatusConflict)
		} else {
			http.Error(w, "Error creating user", http.StatusInternalServerError)
		}
		return
	}

	// Prepare response
	response := api.UserResponse{
		User: user,
	}

	// Write response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// Login authenticates a user
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var request api.LoginUserRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Get user from database
	user, err := h.DB.GetUserByEmail(request.User.Email)
	if err != nil {
		if err == db.ErrNotFound {
			http.Error(w, "User not found", http.StatusNotFound)
		} else {
			http.Error(w, "Error retrieving user", http.StatusInternalServerError)
		}
		return
	}

	// Verify password
	if err := h.DB.VerifyUserPassword(request.User.Email, request.User.Password); err != nil {
		if err == db.ErrInvalidCredentials {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		} else {
			http.Error(w, "Error verifying credentials", http.StatusInternalServerError)
		}
		return
	}

	// Generate token
	authConfig := h.AuthConfig
	token, err := auth.GenerateToken(user.Email, authConfig)
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	// Update user with token
	user.Token = token

	// Prepare response
	response := api.UserResponse{
		User: *user,
	}

	// Write response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetCurrentUser returns the current user
func (h *Handler) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	// Get authenticated user from context
	email, ok := middleware.GetUserEmail(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get user from database
	user, err := h.DB.GetUserByEmail(email)
	if err != nil {
		if err == db.ErrNotFound {
			http.Error(w, "User not found", http.StatusNotFound)
		} else {
			http.Error(w, "Error retrieving user", http.StatusInternalServerError)
		}
		return
	}

	// Generate a fresh token
	token, err := auth.GenerateToken(email, h.AuthConfig)
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}
	user.Token = token

	// Prepare response
	response := api.UserResponse{
		User: *user,
	}

	// Write response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// UpdateCurrentUser updates the current user
func (h *Handler) UpdateCurrentUser(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var request api.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Get authenticated user from context
	email, ok := middleware.GetUserEmail(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Update user in database
	user, err := h.DB.UpdateUser(email, request.User)
	if err != nil {
		if err == db.ErrNotFound {
			http.Error(w, "User not found", http.StatusNotFound)
		} else if err == db.ErrConflict {
			http.Error(w, "User with this email or username already exists", http.StatusConflict)
		} else {
			http.Error(w, "Error updating user", http.StatusInternalServerError)
		}
		return
	}

	// Generate a fresh token
	token, err := auth.GenerateToken(user.Email, h.AuthConfig)
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}
	user.Token = token

	// Prepare response
	response := api.UserResponse{
		User: *user,
	}

	// Write response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Implement the remaining methods of the ServerInterface
// These are just stubs for now, but they satisfy the interface

func (h *Handler) DeleteArticle(w http.ResponseWriter, r *http.Request, slug string) {
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}

func (h *Handler) GetArticle(w http.ResponseWriter, r *http.Request, slug string) {
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}

func (h *Handler) UpdateArticle(w http.ResponseWriter, r *http.Request, slug string) {
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}

func (h *Handler) GetArticleComments(w http.ResponseWriter, r *http.Request, slug string) {
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}

func (h *Handler) CreateArticleComment(w http.ResponseWriter, r *http.Request, slug string) {
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}

func (h *Handler) DeleteArticleComment(w http.ResponseWriter, r *http.Request, slug string, id int) {
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}

func (h *Handler) DeleteArticleFavorite(w http.ResponseWriter, r *http.Request, slug string) {
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}

func (h *Handler) CreateArticleFavorite(w http.ResponseWriter, r *http.Request, slug string) {
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}

func (h *Handler) GetProfileByUsername(w http.ResponseWriter, r *http.Request, username string) {
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}

func (h *Handler) UnfollowUserByUsername(w http.ResponseWriter, r *http.Request, username string) {
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}

func (h *Handler) FollowUserByUsername(w http.ResponseWriter, r *http.Request, username string) {
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}
