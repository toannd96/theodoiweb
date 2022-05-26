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
