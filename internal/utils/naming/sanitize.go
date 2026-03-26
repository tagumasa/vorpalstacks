// Package naming provides naming and string manipulation utilities for vorpalstacks.
package naming

import (
	"strings"
	"unicode"
)

// SanitizeFilename removes or replaces characters that are invalid in filenames.
//
// Parameters:
//   - s: The input filename
//
// Returns:
//   - string: The sanitized filename
//
// Example:
//
//	SanitizeFilename("my/file:name.txt") // "my_file_name_.txt"
//	SanitizeFilename("my-file.txt") // "my-file.txt"
func SanitizeFilename(s string) string {
	replacer := strings.NewReplacer(
		"/", "_",
		"\\", "_",
		":", "_",
		"*", "_",
		"?", "_",
		"\"", "_",
		"<", "_",
		">", "_",
		"|", "_",
	)
	return replacer.Replace(s)
}

// SanitizeName removes invalid characters from a name.
//
// Parameters:
//   - inputName: The input name to sanitize
//
// Returns:
//   - string: The sanitized name containing only letters, digits, and underscores
//
// Example:
//
//	SanitizeName("my-name_123") // "myname_123"
//	SanitizeName("name@#$%") // "name"
func SanitizeName(inputName string) string {
	var result []rune
	for _, r := range inputName {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' {
			result = append(result, r)
		}
	}
	return string(result)
}
