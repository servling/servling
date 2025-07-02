package util

import (
	"regexp"
	"strings"
)

func NormalizeContainerName(s string) string {
	s = strings.ToLower(s)

	invalidCharsRegex := regexp.MustCompile("[^a-z0-9_.]+")
	s = invalidCharsRegex.ReplaceAllString(s, "_")

	multipleUnderscoresRegex := regexp.MustCompile("__+")
	s = multipleUnderscoresRegex.ReplaceAllString(s, "_")

	s = strings.Trim(s, "_.")

	return s
}
