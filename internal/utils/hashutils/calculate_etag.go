// Package hashutils provides hash utility functions for vorpalstacks.
package hashutils

import (
	"crypto/sha256"
	"encoding/hex"
)

// CalculateETag calculates SHA-256 hash of data for ETag.
//
// Parameters:
//   - data: The data to calculate ETag for
//
// Returns:
//   - string: The SHA-256 hash as hex string
func CalculateETag(data []byte) string {
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])
}
