// Package generators provides ID generation utilities for vorpalstacks AWS resources.
package generators

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
)

const randomIDSize = 16

// GenerateIDWithPrefix generates a random URL-safe ID with the given prefix
// and suffix length. For example, GenerateIDWithPrefix("queue-", 12) returns
// a string like "queue-abcDEF123ghi".
func GenerateIDWithPrefix(prefix string, suffixLen int) (string, error) {
	return generateIDWithPrefix(prefix, suffixLen)
}

func generateIDWithPrefix(prefix string, suffixLen int) (string, error) {
	b := make([]byte, randomIDSize)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}
	encoded := base64.RawURLEncoding.EncodeToString(b)
	if suffixLen > len(encoded) {
		suffixLen = len(encoded)
	}
	return prefix + encoded[:suffixLen], nil
}

func generateHexID() (string, error) {
	b := make([]byte, randomIDSize)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}
	return hex.EncodeToString(b), nil
}
