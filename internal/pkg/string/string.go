package pkg

import (
	"net/url"
	"strings"
	"unicode"
)

// RemoveSubstring remove substring in string
func RemoveSubstring(s string, substr string) string {
	if n := strings.Index(s, substr); n >= 0 {
		return strings.TrimRightFunc(s[:n], unicode.IsSpace)
	}
	return s
}

// RemoveDuplicateValues remove duplicate values from string slice
func RemoveDuplicateValues(strSlice []string) []string {
	keys := make(map[string]bool)
	list := []string{}

	for _, entry := range strSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func ParseURL(input string) (string, error) {
	url, err := url.Parse(input)
	if err != nil {
		return "", err
	}
	hostname := strings.TrimPrefix(url.Hostname(), "www.")
	return hostname, nil
}
