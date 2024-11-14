package utils

import "strings"

func HasPathAfterNews(url string) bool {
	if !strings.Contains(url, "finance.yahoo.com") {
		return false
	}

	newsIndex := strings.Index(url, "/news/")
	if newsIndex == -1 {
		return false
	}
	return len(url) > newsIndex+len("/news/")
}
