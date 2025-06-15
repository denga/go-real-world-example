package db

import (
	"testing"
	"time"

	"github.com/denga/go-real-world-example/api"
)

func TestUserOperations(t *testing.T) {
	// Create a new in-memory database
	db := NewInMemoryDB()

	// Test user data
	user := api.User{
		Username: "testuser",
		Email:    "test@example.com",
		Bio:      "Test bio",
		Image:    "test-image.jpg",
		Token:    "test-token",
	}
	password := "password123"

	// Test CreateUser
	err := db.CreateUser(user, password)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	// Test GetUserByEmail
	retrievedUser, err := db.GetUserByEmail(user.Email)
	if err != nil {
		t.Fatalf("Failed to get user by email: %v", err)
	}
	if retrievedUser.Username != user.Username {
		t.Errorf("Expected username %s, got %s", user.Username, retrievedUser.Username)
	}

	// Test GetUserByUsername
	retrievedUser, err = db.GetUserByUsername(user.Username)
	if err != nil {
		t.Fatalf("Failed to get user by username: %v", err)
	}
	if retrievedUser.Email != user.Email {
		t.Errorf("Expected email %s, got %s", user.Email, retrievedUser.Email)
	}

	// Test VerifyUserPassword
	err = db.VerifyUserPassword(user.Email, password)
	if err != nil {
		t.Errorf("Failed to verify correct password: %v", err)
	}

	err = db.VerifyUserPassword(user.Email, "wrongpassword")
	if err != ErrInvalidCredentials {
		t.Errorf("Expected ErrInvalidCredentials for wrong password, got %v", err)
	}

	// Test UpdateUser
	newEmail := "newemail@example.com"
	newUsername := "newusername"
	newBio := "New bio"
	newImage := "new-image.jpg"

	updates := api.UpdateUser{
		Email:    &newEmail,
		Username: &newUsername,
		Bio:      &newBio,
		Image:    &newImage,
	}

	updatedUser, err := db.UpdateUser(user.Email, updates)
	if err != nil {
		t.Fatalf("Failed to update user: %v", err)
	}
	if updatedUser.Email != newEmail {
		t.Errorf("Expected updated email %s, got %s", newEmail, updatedUser.Email)
	}
	if updatedUser.Username != newUsername {
		t.Errorf("Expected updated username %s, got %s", newUsername, updatedUser.Username)
	}
	if updatedUser.Bio != newBio {
		t.Errorf("Expected updated bio %s, got %s", newBio, updatedUser.Bio)
	}
	if updatedUser.Image != newImage {
		t.Errorf("Expected updated image %s, got %s", newImage, updatedUser.Image)
	}

	// Test conflict errors
	conflictUser := api.User{
		Username: "conflictuser",
		Email:    newEmail, // Same email as updated user
		Bio:      "Conflict bio",
		Image:    "conflict-image.jpg",
		Token:    "conflict-token",
	}
	err = db.CreateUser(conflictUser, "conflictpassword")
	if err != ErrConflict {
		t.Errorf("Expected ErrConflict for duplicate email, got %v", err)
	}

	conflictUser.Email = "unique@example.com"
	conflictUser.Username = newUsername // Same username as updated user
	err = db.CreateUser(conflictUser, "conflictpassword")
	if err != ErrConflict {
		t.Errorf("Expected ErrConflict for duplicate username, got %v", err)
	}
}

func TestArticleOperations(t *testing.T) {
	// Create a new in-memory database
	db := NewInMemoryDB()

	// Create a test user
	user := api.User{
		Username: "testuser",
		Email:    "test@example.com",
		Bio:      "Test bio",
		Image:    "test-image.jpg",
		Token:    "test-token",
	}
	db.CreateUser(user, "password123")

	// Create author profile
	author := api.Profile{
		Username:  user.Username,
		Bio:       user.Bio,
		Image:     user.Image,
		Following: false,
	}

	// Test article data
	article := api.Article{
		Title:       "Test Article",
		Description: "Test description",
		Body:        "Test body",
		TagList:     []string{"test", "article"},
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Favorited:   false,
		Author:      author,
		Slug:        "test-article",
	}

	// Test CreateArticle
	err := db.CreateArticle(article)
	if err != nil {
		t.Fatalf("Failed to create article: %v", err)
	}

	// Test GetArticle
	retrievedArticle, err := db.GetArticle(article.Slug)
	if err != nil {
		t.Fatalf("Failed to get article: %v", err)
	}
	if retrievedArticle.Title != article.Title {
		t.Errorf("Expected title %s, got %s", article.Title, retrievedArticle.Title)
	}

	// Test UpdateArticle
	newTitle := "Updated Title"
	newDescription := "Updated description"
	newBody := "Updated body"

	updates := api.UpdateArticle{
		Title:       &newTitle,
		Description: &newDescription,
		Body:        &newBody,
	}

	updatedArticle, err := db.UpdateArticle(article.Slug, updates)
	if err != nil {
		t.Fatalf("Failed to update article: %v", err)
	}
	if updatedArticle.Title != newTitle {
		t.Errorf("Expected updated title %s, got %s", newTitle, updatedArticle.Title)
	}
	if updatedArticle.Description != newDescription {
		t.Errorf("Expected updated description %s, got %s", newDescription, updatedArticle.Description)
	}
	if updatedArticle.Body != newBody {
		t.Errorf("Expected updated body %s, got %s", newBody, updatedArticle.Body)
	}

	// Test ListArticles
	articles, count, err := db.ListArticles("", "", "", 10, 0)
	if err != nil {
		t.Fatalf("Failed to list articles: %v", err)
	}
	if count != 1 {
		t.Errorf("Expected 1 article, got %d", count)
	}
	if len(articles) != 1 {
		t.Errorf("Expected 1 article in list, got %d", len(articles))
	}

	// Test filtering by tag
	articles, count, err = db.ListArticles("test", "", "", 10, 0)
	if err != nil {
		t.Fatalf("Failed to list articles filtered by tag: %v", err)
	}
	if count != 1 {
		t.Errorf("Expected 1 article with tag 'test', got %d", count)
	}

	// Test filtering by author
	articles, count, err = db.ListArticles("", user.Username, "", 10, 0)
	if err != nil {
		t.Fatalf("Failed to list articles filtered by author: %v", err)
	}
	if count != 1 {
		t.Errorf("Expected 1 article by author '%s', got %d", user.Username, count)
	}

	// Test DeleteArticle
	err = db.DeleteArticle(article.Slug)
	if err != nil {
		t.Fatalf("Failed to delete article: %v", err)
	}

	// Verify article is deleted
	_, err = db.GetArticle(article.Slug)
	if err != ErrNotFound {
		t.Errorf("Expected ErrNotFound after deletion, got %v", err)
	}
}

func TestCommentOperations(t *testing.T) {
	// Create a new in-memory database
	db := NewInMemoryDB()

	// Create a test user
	user := api.User{
		Username: "testuser",
		Email:    "test@example.com",
		Bio:      "Test bio",
		Image:    "test-image.jpg",
		Token:    "test-token",
	}
	db.CreateUser(user, "password123")

	// Create author profile
	author := api.Profile{
		Username:  user.Username,
		Bio:       user.Bio,
		Image:     user.Image,
		Following: false,
	}

	// Create a test article
	article := api.Article{
		Title:       "Test Article",
		Description: "Test description",
		Body:        "Test body",
		TagList:     []string{"test", "article"},
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Favorited:   false,
		Author:      author,
		Slug:        "test-article",
	}
	db.CreateArticle(article)

	// Test comment data
	comment := api.Comment{
		Body:      "Test comment",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Author:    author,
	}

	// Test AddComment
	commentID, err := db.AddComment(article.Slug, comment)
	if err != nil {
		t.Fatalf("Failed to add comment: %v", err)
	}
	if commentID != 1 {
		t.Errorf("Expected comment ID 1, got %d", commentID)
	}

	// Test GetComments
	comments, err := db.GetComments(article.Slug)
	if err != nil {
		t.Fatalf("Failed to get comments: %v", err)
	}
	if len(comments) != 1 {
		t.Errorf("Expected 1 comment, got %d", len(comments))
	}
	if comments[0].Body != comment.Body {
		t.Errorf("Expected comment body %s, got %s", comment.Body, comments[0].Body)
	}

	// Test DeleteComment
	err = db.DeleteComment(article.Slug, commentID)
	if err != nil {
		t.Fatalf("Failed to delete comment: %v", err)
	}

	// Verify comment is deleted
	comments, err = db.GetComments(article.Slug)
	if err != nil {
		t.Fatalf("Failed to get comments after deletion: %v", err)
	}
	if len(comments) != 0 {
		t.Errorf("Expected 0 comments after deletion, got %d", len(comments))
	}
}

func TestFollowOperations(t *testing.T) {
	// Create a new in-memory database
	db := NewInMemoryDB()

	// Create test users
	follower := api.User{
		Username: "follower",
		Email:    "follower@example.com",
		Bio:      "Follower bio",
		Image:    "follower-image.jpg",
		Token:    "follower-token",
	}
	db.CreateUser(follower, "password123")

	followed := api.User{
		Username: "followed",
		Email:    "followed@example.com",
		Bio:      "Followed bio",
		Image:    "followed-image.jpg",
		Token:    "followed-token",
	}
	db.CreateUser(followed, "password123")

	// Test FollowUser
	err := db.FollowUser(follower.Username, followed.Username)
	if err != nil {
		t.Fatalf("Failed to follow user: %v", err)
	}

	// Test IsFollowing
	isFollowing := db.IsFollowing(follower.Username, followed.Username)
	if !isFollowing {
		t.Errorf("Expected follower to be following followed, but IsFollowing returned false")
	}

	// Test UnfollowUser
	err = db.UnfollowUser(follower.Username, followed.Username)
	if err != nil {
		t.Fatalf("Failed to unfollow user: %v", err)
	}

	// Verify user is unfollowed
	isFollowing = db.IsFollowing(follower.Username, followed.Username)
	if isFollowing {
		t.Errorf("Expected follower to not be following followed after unfollowing, but IsFollowing returned true")
	}
}

func TestFavoriteOperations(t *testing.T) {
	// Create a new in-memory database
	db := NewInMemoryDB()

	// Create a test user
	user := api.User{
		Username: "testuser",
		Email:    "test@example.com",
		Bio:      "Test bio",
		Image:    "test-image.jpg",
		Token:    "test-token",
	}
	db.CreateUser(user, "password123")

	// Create author profile
	author := api.Profile{
		Username:  user.Username,
		Bio:       user.Bio,
		Image:     user.Image,
		Following: false,
	}

	// Create a test article
	article := api.Article{
		Title:       "Test Article",
		Description: "Test description",
		Body:        "Test body",
		TagList:     []string{"test", "article"},
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Favorited:   false,
		Author:      author,
		Slug:        "test-article",
	}
	db.CreateArticle(article)

	// Test FavoriteArticle
	err := db.FavoriteArticle(article.Slug, user.Username)
	if err != nil {
		t.Fatalf("Failed to favorite article: %v", err)
	}

	// Test IsFavorite
	isFavorite := db.IsFavorite(article.Slug, user.Username)
	if !isFavorite {
		t.Errorf("Expected article to be favorited, but IsFavorite returned false")
	}

	// Check that favorites count was updated
	retrievedArticle, err := db.GetArticle(article.Slug)
	if err != nil {
		t.Fatalf("Failed to get article after favoriting: %v", err)
	}
	if retrievedArticle.FavoritesCount != 1 {
		t.Errorf("Expected favorites count to be 1, got %d", retrievedArticle.FavoritesCount)
	}

	// Test UnfavoriteArticle
	err = db.UnfavoriteArticle(article.Slug, user.Username)
	if err != nil {
		t.Fatalf("Failed to unfavorite article: %v", err)
	}

	// Verify article is unfavorited
	isFavorite = db.IsFavorite(article.Slug, user.Username)
	if isFavorite {
		t.Errorf("Expected article to not be favorited after unfavoriting, but IsFavorite returned true")
	}

	// Check that favorites count was updated
	retrievedArticle, err = db.GetArticle(article.Slug)
	if err != nil {
		t.Fatalf("Failed to get article after unfavoriting: %v", err)
	}
	if retrievedArticle.FavoritesCount != 0 {
		t.Errorf("Expected favorites count to be 0 after unfavoriting, got %d", retrievedArticle.FavoritesCount)
	}
}

func TestGetTags(t *testing.T) {
	// Create a new in-memory database
	db := NewInMemoryDB()

	// Create a test user
	user := api.User{
		Username: "testuser",
		Email:    "test@example.com",
		Bio:      "Test bio",
		Image:    "test-image.jpg",
		Token:    "test-token",
	}
	db.CreateUser(user, "password123")

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
	db.CreateArticle(article1)

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
	db.CreateArticle(article2)

	// Test GetTags
	tags := db.GetTags()
	if len(tags) != 3 {
		t.Errorf("Expected 3 unique tags, got %d", len(tags))
	}

	// Check that all expected tags are present
	expectedTags := map[string]bool{
		"tag1": true,
		"tag2": true,
		"tag3": true,
	}
	for _, tag := range tags {
		if !expectedTags[tag] {
			t.Errorf("Unexpected tag: %s", tag)
		}
		delete(expectedTags, tag)
	}
	if len(expectedTags) != 0 {
		t.Errorf("Missing tags: %v", expectedTags)
	}
}
