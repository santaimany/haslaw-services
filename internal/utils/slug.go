package utils

import (
	"crypto/rand"
	"encoding/hex"
	"regexp"
	"strings"
	"unicode"

	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

func GenerateSlug(title string) string {

	slug := strings.ToLower(title)

	t := transform.Chain(norm.NFD, transform.RemoveFunc(isMn), norm.NFC)
	slug, _, _ = transform.String(t, slug)

	reg := regexp.MustCompile(`[^a-z0-9]+`)
	slug = reg.ReplaceAllString(slug, "-")

	slug = strings.Trim(slug, "-")

	reg = regexp.MustCompile(`-+`)
	slug = reg.ReplaceAllString(slug, "-")

	return slug
}

func GenerateSlugWithRandomID(title string) string {
	baseSlug := GenerateSlug(title)
	randomID := GenerateRandomID(10)
	return baseSlug + "-" + randomID
}

func GenerateRandomID(length int) string {
	bytes := make([]byte, length/2+1)
	rand.Read(bytes)
	randomStr := hex.EncodeToString(bytes)

	result := ""
	for i, char := range randomStr {
		if i < length {
			if i%2 == 0 {
				result += strings.ToLower(string(char))
			} else {
				result += strings.ToUpper(string(char))
			}
		}
	}

	if len(result) > length {
		result = result[:length]
	}

	return strings.ToLower(result)
}

func isMn(r rune) bool {
	return unicode.Is(unicode.Mn, r)
}

func ValidateSlug(slug string) bool {
	if slug == "" {
		return false
	}

	reg := regexp.MustCompile(`^[a-z0-9-]+$`)
	return reg.MatchString(slug)
}

func CleanSlug(slug string) string {
	if slug == "" {
		return ""
	}

	cleaned := GenerateSlug(slug)

	if ValidateSlug(cleaned) {
		return cleaned
	}

	return ""
}
