package utils

import (
	"regexp"
	"strconv"
	"strings"
)

var RE *regexp.Regexp

func init() {
	RE = regexp.MustCompile("[^a-z0-9]+")
}

func Slugify(s string) string {
	slug := strings.Trim(RE.ReplaceAllString(strings.ToLower(s), "-"), "-")
	dashes := 0
	for i, char := range slug {
		if string(char) == "-" {
			dashes++
		}
		if dashes > 5 {
			slug = slug[0:i]
			break
		}
	}
	sanitizedSlug := ""
	for _, r := range slug {
		asciiChar := strconv.QuoteRuneToASCII(r)
		if len(asciiChar) == 1 {
			sanitizedSlug += string(asciiChar)
		}
	}
	if len(strings.TrimSpace(slug)) <= 0 {
		slug = NewOctid()
	}
	return slug
}
