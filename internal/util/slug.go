package util

import (
	"regexp"
	"strings"
)

var (
	// nonAlphanumericRegex matches any character that is not a letter, number, or space
	nonAlphanumericRegex = regexp.MustCompile(`[^a-zA-Z0-9\s]`)
	// multipleSpacesRegex matches multiple spaces
	multipleSpacesRegex = regexp.MustCompile(`\s+`)
)

// GenerateSlug generates a URL-friendly slug from a string
func GenerateSlug(s string) string {
	// Convert to lowercase
	s = strings.ToLower(s)

	// Remove non-alphanumeric characters
	s = nonAlphanumericRegex.ReplaceAllString(s, "")

	// Replace spaces with hyphens
	s = multipleSpacesRegex.ReplaceAllString(s, " ")
	s = strings.TrimSpace(s)
	s = strings.ReplaceAll(s, " ", "-")

	return s
}

// GenerateUniqueSlug generates a unique slug by appending a number if necessary
func GenerateUniqueSlug(title string, exists func(string) bool) string {
	baseSlug := GenerateSlug(title)
	slug := baseSlug
	counter := 1

	// Keep incrementing counter until we find a unique slug
	for exists(slug) {
		slug = baseSlug + "-" + string(rune('0'+counter))
		counter++
	}

	return slug
}