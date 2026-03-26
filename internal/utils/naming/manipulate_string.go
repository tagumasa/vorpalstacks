// Package naming provides naming and string manipulation utilities for vorpalstacks.
package naming

import "strings"

// ContainsAny checks if a string contains any of the substrings.
//
// Parameters:
//   - s: The string to check
//   - substrings: The substrings to look for
//
// Returns:
//   - bool: True if any substring is found
//
// Example:
//
//	ContainsAny("hello world", "foo", "world") // true
//	ContainsAny("hello world", "foo", "bar") // false
func ContainsAny(s string, substrings ...string) bool {
	for _, substr := range substrings {
		if strings.Contains(s, substr) {
			return true
		}
	}
	return false
}

// ReplaceAll replaces all occurrences of old with new in s.
// This is a convenience function for strings.Replace with -1 count.
//
// Parameters:
//   - s: The string to replace in
//   - old: The old substring to replace
//   - new: The new substring to replace with
//
// Returns:
//   - string: The string with all occurrences replaced
//
// Example:
//
//	ReplaceAll("hello world", "o", "O") // "hellO wOrld"
func ReplaceAll(s, old, new string) string {
	return strings.ReplaceAll(s, old, new)
}

// TrimPrefixes removes all specified prefixes from a string.
// It checks each prefix in order and removes it if present.
//
// Parameters:
//   - s: The input string
//   - prefixes: The prefixes to remove
//
// Returns:
//   - string: The string with prefixes removed
//
// Example:
//
//	TrimPrefixes("/api/v1/users", "/api") // "/v1/users"
//	TrimPrefixes("/api/v1/users", "/api", "/v1") // "/users"
func TrimPrefixes(s string, prefixes ...string) string {
	for _, prefix := range prefixes {
		s = strings.TrimPrefix(s, prefix)
	}
	return s
}

// Truncate truncates a string to the specified length, adding an ellipsis if truncated.
//
// Parameters:
//   - s: The input string
//   - maxLen: The maximum length
//
// Returns:
//   - string: The truncated string
//
// Example:
//
//	Truncate("Hello World", 8) // "Hello..."
//	Truncate("Hi", 8) // "Hi"
func Truncate(s string, maxLen int) string {
	if maxLen < 0 {
		maxLen = 0
	}
	if len(s) <= maxLen {
		return s
	}
	if maxLen <= 3 {
		return s[:maxLen]
	}
	return s[:maxLen-3] + "..."
}
