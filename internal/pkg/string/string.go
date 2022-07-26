package pkg

import (
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
