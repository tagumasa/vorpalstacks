// Package generators provides AWS-specific utility functions for vorpalstacks.
package generators

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

// GenerateIAMSecretAccessKey generates an IAM secret access key.
//
// Returns:
//   - string: The generated secret access key
//   - error: An error if generation fails
func GenerateIAMSecretAccessKey() (string, error) {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, 40)
	for i := range result {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", fmt.Errorf("failed to generate random number: %w", err)
		}
		result[i] = charset[num.Int64()]
	}
	return string(result), nil
}
