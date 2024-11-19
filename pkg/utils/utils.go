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

var UaList = []string{
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
	"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
}
