// Package cloudfront provides CloudFront storage functionality for vorpalstacks.
package cloudfront

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"
)

func generateDistributionID() (string, error) {
	bytes := make([]byte, 8)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate distribution ID: %w", err)
	}
	return hex.EncodeToString(bytes), nil
}

func generateETag() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate ETag: %w", err)
	}
	return hex.EncodeToString(bytes), nil
}

func generateDomainName(id string) string {
	return fmt.Sprintf("%s.cloudfront.net", id)
}

func normalizeCallerReference(ref string) string {
	return strings.ReplaceAll(ref, " ", "-")
}
