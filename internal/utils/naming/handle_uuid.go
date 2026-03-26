// Package naming provides naming and ID generation utilities for vorpalstacks.
package naming

import (
	"github.com/google/uuid"
)

// GenerateUUID generates a new UUID v4.
//
// Returns:
//   - string: The generated UUID
//
// Example:
//
//	id := GenerateUUID() // "550e8400-e29b-41d4-a716-4466554400000"
func GenerateUUID() string {
	return uuid.New().String()
}

// GenerateUUIDBytes generates a new UUID v4 as bytes.
//
// Returns:
//   - []byte: The generated UUID as bytes
func GenerateUUIDBytes() []byte {
	id := uuid.New()
	return id[:]
}

// UUIDNil represents the nil UUID.
var UUIDNil = uuid.Nil

// IsNilUUID checks if a UUID is nil.
//
// Parameters:
//   - id: The UUID to check
//
// Returns:
//   - bool: True if the UUID is nil
//
// Example:
//
//	IsNilUUID(uuid.Nil) // true
//	IsNilUUID(uuid.New()) // false
func IsNilUUID(id uuid.UUID) bool {
	return id == UUIDNil
}

// IsValidUUID checks if a string is a valid UUID.
//
// Parameters:
//   - s: The string to check
//
// Returns:
//   - bool: True if the string is a valid UUID
//
// Example:
//
//	IsValidUUID("550e8400-e29b-41d4-a716-4466554400000") // true
//	IsValidUUID("invalid") // false
func IsValidUUID(s string) bool {
	_, err := uuid.Parse(s)
	return err == nil
}

// ParseUUID parses a UUID string.
//
// Parameters:
//   - s: The UUID string to parse
//
// Returns:
//   - uuid.UUID: The parsed UUID
//   - error: An error if parsing fails
//
// Example:
//
//	id, err := ParseUUID("550e8400-e29b-41d4-a716-4466554400000")
func ParseUUID(s string) (uuid.UUID, error) {
	return uuid.Parse(s)
}
