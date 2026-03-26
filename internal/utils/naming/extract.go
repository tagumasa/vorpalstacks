// Package naming provides naming and string manipulation utilities for vorpalstacks.
package naming

import "strings"

// ExtractServiceName extracts the service name from a shape ID.
//
// Parameters:
//   - inputName: The input shape ID (e.g., "com.amazonaws.lambda#DeleteFunction")
//
// Returns:
//   - string: The extracted service name (e.g., "DeleteFunction")
//
// Example:
//
//	ExtractServiceName("com.amazonaws.lambda#DeleteFunction") // "DeleteFunction"
//	ExtractServiceName("simple_name") // "simple_name"
func ExtractServiceName(inputName string) string {
	parts := strings.Split(inputName, "#")
	if len(parts) > 1 {
		return parts[1]
	}
	return inputName
}

// ExtractServiceNameFromNamespace extracts the service name from a namespace.
//
// Parameters:
//   - inputName: The input namespace (e.g., "com.amazonaws.lambda")
//
// Returns:
//   - string: The extracted service name (e.g., "lambda")
//
// Example:
//
//	ExtractServiceNameFromNamespace("com.amazonaws.lambda") // "lambda"
//	ExtractServiceNameFromNamespace("simple") // "simple"
func ExtractServiceNameFromNamespace(inputName string) string {
	parts := strings.Split(inputName, ".")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}
	return inputName
}

// ExtractNamespace extracts the namespace from a shape ID.
//
// Parameters:
//   - inputName: The input shape ID (e.g., "com.amazonaws.lambda#DeleteFunction")
//
// Returns:
//   - string: The extracted namespace (e.g., "com.amazonaws.lambda")
//
// Example:
//
//	ExtractNamespace("com.amazonaws.lambda#DeleteFunction") // "com.amazonaws.lambda"
//	ExtractNamespace("simple_name") // ""
func ExtractNamespace(inputName string) string {
	parts := strings.Split(inputName, "#")
	if len(parts) > 1 {
		return parts[0]
	}
	return ""
}
