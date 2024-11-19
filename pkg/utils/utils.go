package utils

import (
	"regexp"
	"strings"
)

func HasPathAfterNews(url string) bool {
	if strings.Contains(url, "finance.yahoo.com") {
		newsIndex := strings.Index(url, "/news/")
		if newsIndex != -1 {
			return len(url) > newsIndex+len("/news/")
		}
	}

	if strings.Contains(url, "benzinga.com") {
		if strings.Contains(url, "/recent") {
			return false
		}
		return true
	}

	return false
}

// ContainsURLLink checks if a URL matches any of the patterns in the slice.
// The patterns can contain wildcards (*) and will be converted to regex patterns.
// For example, patterns like:
// - "https://www.benzinga.com/trading-ideas/*"
// - "https://www.benzinga.com/news/*"
// - "https://www.benzinga.com/[0-9][0-9]/[0-9][0-9]/*"
// - "https://www.benzinga.com/markets/*"
func ContainsURLLink(slice []string, item string) bool {
	for _, s := range slice {
		// Convert pattern to regex by escaping dots and replacing * with .*
		pattern := strings.ReplaceAll(strings.ReplaceAll(s, ".", "\\."), "*", ".*")
		matched, _ := regexp.MatchString("^"+pattern+"$", item)
		if matched {
			return true
		}
	}
	return false
}
