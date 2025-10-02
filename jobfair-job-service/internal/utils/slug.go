package utils

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

// GenerateSlug generates a URL-friendly slug from a string
func GenerateSlug(text string) string {
	// Convert to lowercase
	slug := strings.ToLower(text)

	// Replace spaces with hyphens
	slug = strings.ReplaceAll(slug, " ", "-")

	// Remove special characters
	reg := regexp.MustCompile("[^a-z0-9-]+")
	slug = reg.ReplaceAllString(slug, "")

	// Remove multiple consecutive hyphens
	reg = regexp.MustCompile("-+")
	slug = reg.ReplaceAllString(slug, "-")

	// Trim hyphens from start and end
	slug = strings.Trim(slug, "-")

	// Add timestamp to ensure uniqueness
	timestamp := time.Now().Unix()
	slug = fmt.Sprintf("%s-%d", slug, timestamp)

	return slug
}
