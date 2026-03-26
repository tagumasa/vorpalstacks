// Package naming provides naming and string manipulation utilities for vorpalstacks.
package naming

import (
	"strings"
	"unicode"
)

// ToCamelCase converts a snake_case string to camelCase.
//
// Parameters:
//   - inputName: The input string in snake_case
//
// Returns:
//   - string: The converted string in camelCase
//
// Example:
//
//	ToCamelCase("hello_world") // "helloWorld"
func ToCamelCase(inputName string) string {
	pascal := ToPascalCase(inputName)
	if len(pascal) == 0 {
		return ""
	}
	if len(pascal) == 1 {
		return strings.ToLower(pascal)
	}
	return strings.ToLower(pascal[:1]) + pascal[1:]
}

// ToPascalCase converts a snake_case string to PascalCase.
//
// Parameters:
//   - inputName: The input string in snake_case
//
// Returns:
//   - string: The converted string in PascalCase
//
// Example:
//
//	ToPascalCase("hello_world") // "HelloWorld"
func ToPascalCase(inputName string) string {
	parts := strings.Split(inputName, "_")
	var result []string
	for _, part := range parts {
		if len(part) > 0 {
			if len(part) == 1 {
				result = append(result, strings.ToUpper(part))
			} else {
				result = append(result, strings.ToUpper(part[:1])+strings.ToLower(part[1:]))
			}
		}
	}
	return strings.Join(result, "")
}

// ToSnakeCase converts a PascalCase or camelCase string to snake_case.
//
// Parameters:
//   - inputName: The input string in PascalCase or camelCase
//
// Returns:
//   - string: The converted string in snake_case
//
// Example:
//
//	ToSnakeCase("HelloWorld") // "hello_world"
//	ToSnakeCase("helloWorld") // "hello_world"
func ToSnakeCase(inputName string) string {
	var result []rune
	for i, r := range inputName {
		if i > 0 && unicode.IsUpper(r) {
			result = append(result, '_')
		}
		result = append(result, unicode.ToLower(r))
	}
	return string(result)
}
