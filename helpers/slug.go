package helpers

import (
	"regexp"
	"strings"
)

func Slugify(text string) string {
	slug := strings.ToLower(text)

	re := regexp.MustCompile(`[^a-z0-9\s-]`)
	slug = re.ReplaceAllString(slug, "")

	slug = strings.ReplaceAll(slug, " ", "-")
	slug = regexp.MustCompile(`-+`).ReplaceAllString(slug, "-")

	return slug
}