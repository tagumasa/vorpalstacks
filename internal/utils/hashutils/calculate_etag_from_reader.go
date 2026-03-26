// Package hashutils provides hash utility functions for vorpalstacks.
package hashutils

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
)

// CalculateETagFromReader calculates SHA-256 hash from reader for ETag.
//
// Parameters:
//   - r: The reader to read data from
//
// Returns:
//   - string: The SHA-256 hash as hex string
//   - error: An error if reading fails
func CalculateETagFromReader(r io.Reader) (string, error) {
	if r == nil {
		return "", fmt.Errorf("reader is nil")
	}
	hash := sha256.New()
	if _, err := io.Copy(hash, r); err != nil {
		return "", fmt.Errorf("failed to calculate ETag: %w", err)
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}
