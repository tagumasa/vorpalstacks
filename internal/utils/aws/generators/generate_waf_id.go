package generators

// Package generators provides ID generation utilities for vorpalstacks AWS resources.

// GenerateWAFID generates a unique WAF ID.
//
// Returns:
//   - string: The generated WAF ID
//   - error: An error if ID generation fails
func GenerateWAFID() (string, error) {
	return generateIDWithPrefix("WAF", 28)
}
