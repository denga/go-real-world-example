package auth

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var (
	// ErrInvalidToken is returned when the token is invalid
	ErrInvalidToken = errors.New("invalid token")
	// ErrExpiredToken is returned when the token has expired
	ErrExpiredToken = errors.New("token has expired")
	// ErrInvalidCredentials is returned when the credentials are invalid
	ErrInvalidCredentials = errors.New("invalid credentials")
)

// Config holds the configuration for the auth package
type Config struct {
	// Secret is the secret key used to sign JWT tokens
	Secret string
	// TokenExpiry is the duration for which a token is valid
	TokenExpiry time.Duration
}

// DefaultConfig returns a default configuration
func DefaultConfig() Config {
	return Config{
		Secret:      "your-secret-key", // In production, this should be set via environment variables
		TokenExpiry: 24 * time.Hour,    // 24 hours
	}
}

// Claims represents the JWT claims
type Claims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

// GenerateToken generates a JWT token for the given email
func GenerateToken(email string, config Config) (string, error) {
	expirationTime := time.Now().Add(config.TokenExpiry)
	claims := &Claims{
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(config.Secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateToken validates a JWT token and returns the email if valid
func ValidateToken(tokenString string, config Config) (string, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.Secret), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return "", ErrExpiredToken
		}
		return "", ErrInvalidToken
	}

	if !token.Valid {
		return "", ErrInvalidToken
	}

	return claims.Email, nil
}

// ExtractTokenFromRequest extracts the JWT token from the Authorization header
func ExtractTokenFromRequest(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", ErrInvalidToken
	}

	// Check if the header has the format "Token <token>"
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Token" {
		return "", ErrInvalidToken
	}

	return parts[1], nil
}

// HashPassword hashes a password using bcrypt
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// VerifyPassword verifies a password against a hash
func VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}