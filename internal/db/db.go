package db

import (
	"errors"
	"sync"

	"github.com/denga/go-real-world-example/api"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrNotFound           = errors.New("not found")
	ErrConflict           = errors.New("conflict")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

// InternalUser extends api.User with a password field for internal use
type InternalUser struct {
	api.User
	Password string // Hashed password
}

// InMemoryDB is a simple in-memory database implementation
type InMemoryDB struct {
	users     map[string]*InternalUser        // key: email
	usernames map[string]string               // key: username, value: email
	articles  map[string]*api.Article         // key: slug
	comments  map[string]map[int]*api.Comment // key: article slug, value: map of comments by ID
	follows   map[string]map[string]bool      // key: follower username, value: map of followed usernames
	favorites map[string]map[string]bool      // key: article slug, value: map of usernames who favorited
	tags      map[string]bool                 // set of unique tags
	mutex     sync.RWMutex
}

// NewInMemoryDB creates a new in-memory database
func NewInMemoryDB() *InMemoryDB {
	return &InMemoryDB{
		users:     make(map[string]*InternalUser),
		usernames: make(map[string]string),
		articles:  make(map[string]*api.Article),
		comments:  make(map[string]map[int]*api.Comment),
		follows:   make(map[string]map[string]bool),
		favorites: make(map[string]map[string]bool),
		tags:      make(map[string]bool),
	}
}

// CreateUser creates a new user
func (db *InMemoryDB) CreateUser(user api.User, password string) error {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	// Check if email already exists
	if _, exists := db.users[user.Email]; exists {
		return ErrConflict
	}

	// Check if username already exists
	if _, exists := db.usernames[user.Username]; exists {
		return ErrConflict
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Store user
	internalUser := &InternalUser{
		User:     user,
		Password: string(hashedPassword),
	}
	db.users[user.Email] = internalUser
	db.usernames[user.Username] = user.Email

	return nil
}

// GetUserByEmail retrieves a user by email
func (db *InMemoryDB) GetUserByEmail(email string) (*api.User, error) {
	db.mutex.RLock()
	defer db.mutex.RUnlock()

	internalUser, exists := db.users[email]
	if !exists {
		return nil, ErrNotFound
	}

	// Return a copy of the User field
	user := internalUser.User
	return &user, nil
}

// GetInternalUserByEmail retrieves an internal user (including password) by email
func (db *InMemoryDB) GetInternalUserByEmail(email string) (*InternalUser, error) {
	db.mutex.RLock()
	defer db.mutex.RUnlock()

	internalUser, exists := db.users[email]
	if !exists {
		return nil, ErrNotFound
	}

	// Return a copy of the internal user
	user := *internalUser
	return &user, nil
}

// VerifyUserPassword verifies a user's password
func (db *InMemoryDB) VerifyUserPassword(email, password string) error {
	internalUser, err := db.GetInternalUserByEmail(email)
	if err != nil {
		return err
	}

	// Import the auth package for password verification
	// This is a circular dependency, so we'll use bcrypt directly
	err = bcrypt.CompareHashAndPassword([]byte(internalUser.Password), []byte(password))
	if err != nil {
		return ErrInvalidCredentials
	}

	return nil
}

// GetUserByUsername retrieves a user by username
func (db *InMemoryDB) GetUserByUsername(username string) (*api.User, error) {
	db.mutex.RLock()
	defer db.mutex.RUnlock()

	email, exists := db.usernames[username]
	if !exists {
		return nil, ErrNotFound
	}

	internalUser := db.users[email]
	// Return a copy of the User field
	user := internalUser.User
	return &user, nil
}

// UpdateUser updates an existing user
func (db *InMemoryDB) UpdateUser(email string, updates api.UpdateUser) (*api.User, error) {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	internalUser, exists := db.users[email]
	if !exists {
		return nil, ErrNotFound
	}

	// Update fields if provided
	if updates.Email != nil {
		// Check if new email already exists
		if *updates.Email != email {
			if _, exists := db.users[*updates.Email]; exists {
				return nil, ErrConflict
			}
			// Update email
			delete(db.users, email)
			db.users[*updates.Email] = internalUser
			db.usernames[internalUser.Username] = *updates.Email
			internalUser.Email = *updates.Email
		}
	}

	if updates.Username != nil {
		// Check if new username already exists
		if *updates.Username != internalUser.Username {
			if _, exists := db.usernames[*updates.Username]; exists {
				return nil, ErrConflict
			}
			// Update username
			delete(db.usernames, internalUser.Username)
			db.usernames[*updates.Username] = internalUser.Email
			internalUser.Username = *updates.Username
		}
	}

	if updates.Password != nil {
		internalUser.Password = *updates.Password
	}

	if updates.Bio != nil {
		internalUser.Bio = *updates.Bio
	}

	if updates.Image != nil {
		internalUser.Image = *updates.Image
	}

	// Return a copy of the User field
	user := internalUser.User
	return &user, nil
}

// CreateArticle creates a new article
func (db *InMemoryDB) CreateArticle(article api.Article) error {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	// Check if slug already exists
	if _, exists := db.articles[article.Slug]; exists {
		return ErrConflict
	}

	// Store article
	db.articles[article.Slug] = &article

	// Add tags to the set
	for _, tag := range article.TagList {
		db.tags[tag] = true
	}

	// Initialize comments and favorites for this article
	db.comments[article.Slug] = make(map[int]*api.Comment)
	db.favorites[article.Slug] = make(map[string]bool)

	return nil
}

// GetArticle retrieves an article by slug
func (db *InMemoryDB) GetArticle(slug string) (*api.Article, error) {
	db.mutex.RLock()
	defer db.mutex.RUnlock()

	article, exists := db.articles[slug]
	if !exists {
		return nil, ErrNotFound
	}

	return article, nil
}

// UpdateArticle updates an existing article
func (db *InMemoryDB) UpdateArticle(slug string, updates api.UpdateArticle) (*api.Article, error) {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	article, exists := db.articles[slug]
	if !exists {
		return nil, ErrNotFound
	}

	// Update fields if provided
	if updates.Title != nil {
		article.Title = *updates.Title
	}

	if updates.Description != nil {
		article.Description = *updates.Description
	}

	if updates.Body != nil {
		article.Body = *updates.Body
	}

	return article, nil
}

// DeleteArticle deletes an article by slug
func (db *InMemoryDB) DeleteArticle(slug string) error {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	if _, exists := db.articles[slug]; !exists {
		return ErrNotFound
	}

	// Delete article
	delete(db.articles, slug)
	// Delete comments
	delete(db.comments, slug)
	// Delete favorites
	delete(db.favorites, slug)

	return nil
}

// ListArticles returns a list of articles with optional filtering
func (db *InMemoryDB) ListArticles(tag, author, favorited string, limit, offset int) ([]api.Article, int, error) {
	db.mutex.RLock()
	defer db.mutex.RUnlock()

	var articles []api.Article
	for _, article := range db.articles {
		// Filter by tag if provided
		if tag != "" {
			found := false
			for _, t := range article.TagList {
				if t == tag {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}

		// Filter by author if provided
		if author != "" && article.Author.Username != author {
			continue
		}

		// Filter by favorited if provided
		if favorited != "" {
			if _, exists := db.favorites[article.Slug][favorited]; !exists {
				continue
			}
		}

		articles = append(articles, *article)
	}

	// Calculate total count
	totalCount := len(articles)

	// Apply pagination
	if offset >= len(articles) {
		return []api.Article{}, totalCount, nil
	}

	end := offset + limit
	if end > len(articles) {
		end = len(articles)
	}

	return articles[offset:end], totalCount, nil
}

// GetTags returns all unique tags
func (db *InMemoryDB) GetTags() []string {
	db.mutex.RLock()
	defer db.mutex.RUnlock()

	tags := make([]string, 0, len(db.tags))
	for tag := range db.tags {
		tags = append(tags, tag)
	}

	return tags
}

// AddComment adds a comment to an article
func (db *InMemoryDB) AddComment(slug string, comment api.Comment) (int, error) {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	// Check if article exists
	if _, exists := db.articles[slug]; !exists {
		return 0, ErrNotFound
	}

	// Generate ID for the comment
	id := len(db.comments[slug]) + 1
	comment.Id = id

	// Store comment
	db.comments[slug][id] = &comment

	return id, nil
}

// GetComments returns all comments for an article
func (db *InMemoryDB) GetComments(slug string) ([]api.Comment, error) {
	db.mutex.RLock()
	defer db.mutex.RUnlock()

	// Check if article exists
	if _, exists := db.articles[slug]; !exists {
		return nil, ErrNotFound
	}

	comments := make([]api.Comment, 0, len(db.comments[slug]))
	for _, comment := range db.comments[slug] {
		comments = append(comments, *comment)
	}

	return comments, nil
}

// DeleteComment deletes a comment from an article
func (db *InMemoryDB) DeleteComment(slug string, id int) error {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	// Check if article exists
	if _, exists := db.articles[slug]; !exists {
		return ErrNotFound
	}

	// Check if comment exists
	if _, exists := db.comments[slug][id]; !exists {
		return ErrNotFound
	}

	// Delete comment
	delete(db.comments[slug], id)

	return nil
}

// FollowUser makes one user follow another
func (db *InMemoryDB) FollowUser(follower, followed string) error {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	// Check if both users exist
	if _, exists := db.usernames[follower]; !exists {
		return ErrNotFound
	}
	if _, exists := db.usernames[followed]; !exists {
		return ErrNotFound
	}

	// Initialize follows map for follower if it doesn't exist
	if _, exists := db.follows[follower]; !exists {
		db.follows[follower] = make(map[string]bool)
	}

	// Add follow relationship
	db.follows[follower][followed] = true

	return nil
}

// UnfollowUser makes one user unfollow another
func (db *InMemoryDB) UnfollowUser(follower, followed string) error {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	// Check if both users exist
	if _, exists := db.usernames[follower]; !exists {
		return ErrNotFound
	}
	if _, exists := db.usernames[followed]; !exists {
		return ErrNotFound
	}

	// Check if follows map exists for follower
	if _, exists := db.follows[follower]; !exists {
		return nil // Already not following
	}

	// Remove follow relationship
	delete(db.follows[follower], followed)

	return nil
}

// IsFollowing checks if one user is following another
func (db *InMemoryDB) IsFollowing(follower, followed string) bool {
	db.mutex.RLock()
	defer db.mutex.RUnlock()

	// Check if follows map exists for follower
	if _, exists := db.follows[follower]; !exists {
		return false
	}

	// Check if follow relationship exists
	return db.follows[follower][followed]
}

// FavoriteArticle adds an article to a user's favorites
func (db *InMemoryDB) FavoriteArticle(slug, username string) error {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	// Check if article exists
	article, exists := db.articles[slug]
	if !exists {
		return ErrNotFound
	}

	// Check if user exists
	if _, exists := db.usernames[username]; !exists {
		return ErrNotFound
	}

	// Initialize favorites map for article if it doesn't exist
	if _, exists := db.favorites[slug]; !exists {
		db.favorites[slug] = make(map[string]bool)
	}

	// Add favorite relationship
	db.favorites[slug][username] = true

	// Update favorites count
	article.FavoritesCount = len(db.favorites[slug])

	return nil
}

// UnfavoriteArticle removes an article from a user's favorites
func (db *InMemoryDB) UnfavoriteArticle(slug, username string) error {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	// Check if article exists
	article, exists := db.articles[slug]
	if !exists {
		return ErrNotFound
	}

	// Check if user exists
	if _, exists := db.usernames[username]; !exists {
		return ErrNotFound
	}

	// Check if favorites map exists for article
	if _, exists := db.favorites[slug]; !exists {
		return nil // Already not favorited
	}

	// Remove favorite relationship
	delete(db.favorites[slug], username)

	// Update favorites count
	article.FavoritesCount = len(db.favorites[slug])

	return nil
}

// IsFavorite checks if a user has favorited an article
func (db *InMemoryDB) IsFavorite(slug, username string) bool {
	db.mutex.RLock()
	defer db.mutex.RUnlock()

	// Check if favorites map exists for article
	if _, exists := db.favorites[slug]; !exists {
		return false
	}

	// Check if favorite relationship exists
	return db.favorites[slug][username]
}

// GetArticlesFeed returns articles from followed users
func (db *InMemoryDB) GetArticlesFeed(username string, limit, offset int) ([]api.Article, int, error) {
	db.mutex.RLock()
	defer db.mutex.RUnlock()

	// Check if user exists
	if _, exists := db.usernames[username]; !exists {
		return nil, 0, ErrNotFound
	}

	// Get followed users
	followed, exists := db.follows[username]
	if !exists {
		return []api.Article{}, 0, nil // Not following anyone
	}

	// Collect articles from followed users
	var articles []api.Article
	for _, article := range db.articles {
		if followed[article.Author.Username] {
			articles = append(articles, *article)
		}
	}

	// Calculate total count
	totalCount := len(articles)

	// Apply pagination
	if offset >= len(articles) {
		return []api.Article{}, totalCount, nil
	}

	end := offset + limit
	if end > len(articles) {
		end = len(articles)
	}

	return articles[offset:end], totalCount, nil
}
