// Package naming provides naming and ID generation utilities for vorpalstacks.
package naming

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"time"
)

// ErrInvalidLength is returned when an invalid length is provided to a generation function.
var ErrInvalidLength = errors.New("length must be greater than 0")

// GenerateRandomHex generates a random hexadecimal string of the given length.
// Uses crypto/rand for secure random generation.
//
// Parameters:
//   - length: The length of the hex string to generate
//
// Returns:
//   - string: The generated random hex string
//
// Example:
//
//	hexStr := GenerateRandomHex(32) // "a1b2c3d4e5f6..."
func GenerateRandomHex(length int) (string, error) {
	if length <= 0 {
		return "", ErrInvalidLength
	}
	const hexChars = "0123456789abcdef"
	result := make([]byte, length)
	for i := range result {
		num, err := rand.Int(rand.Reader, big.NewInt(16))
		if err != nil {
			return "", fmt.Errorf("failed to generate random bytes: %w", err)
		}
		result[i] = hexChars[num.Int64()]
	}
	return string(result), nil
}

// GenerateRandomString generates a random alphanumeric string of the given length.
// Uses crypto/rand for secure random generation.
//
// Parameters:
//   - length: The length of the string to generate
//
// Returns:
//   - string: The generated random string
//   - error: An error if generation fails
//
// Example:
//
//	randomStr, _ := GenerateRandomString(16) // "aB3dF7gH9jK2mN5p"
func GenerateRandomString(length int) (string, error) {
	if length <= 0 {
		return "", ErrInvalidLength
	}
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", fmt.Errorf("failed to generate random number: %w", err)
		}
		result[i] = charset[num.Int64()]
	}
	return string(result), nil
}

// GenerateRequestID generates a unique request ID.
// This is useful for tracing requests across services.
//
// Returns:
//   - string: The generated request ID
//
// Example:
//
//	requestID := GenerateRequestID()
func GenerateRequestID() string {
	return GenerateUUID()
}

// GenerateResourceID generates a resource ID with a prefix.
//
// Parameters:
//   - prefix: The prefix for the resource ID
//
// Returns:
//   - string: The generated resource ID
//
// Example:
//
//	resourceID := GenerateResourceID("bucket") // "bucket-550e8400"
func GenerateResourceID(prefix string) string {
	return prefix + "-" + GenerateShortID()
}

// GenerateShortID generates a short ID based on UUID.
// It returns the first 8 characters of a UUID.
//
// Returns:
//   - string: The generated short ID
//
// Example:
//
//	shortID := GenerateShortID() // "550e8400"
func GenerateShortID() string {
	id := GenerateUUID()
	if len(id) < 8 {
		return id
	}
	return id[:8]
}

// GenerateTimestampID generates a unique ID based on Unix timestamp.
//
// Parameters:
//   - prefix: The prefix for the ID
//
// Returns:
//   - string: The generated timestamp-based ID
//
// Example:
//
//	id := GenerateTimestampID("api") // "api-1707260800123456789"
func GenerateTimestampID(prefix string) string {
	return fmt.Sprintf("%s-%d", prefix, time.Now().UnixNano())
}
