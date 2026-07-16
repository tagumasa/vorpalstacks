// Package validators provides AWS-specific utility functions for vorpalstacks.
package validators

import (
	"fmt"
)

// ValidateTopicName validates an SNS topic name.
//
// Parameters:
//   - name: The topic name to validate
//
// Returns:
//   - error: An error if validation fails, nil otherwise
//
// AWS SNS topic name rules:
//   - Must be 1-256 characters
//   - Can contain alphanumeric characters, hyphens, and underscores
func ValidateTopicName(name string) error {
	if len(name) == 0 {
		return fmt.Errorf("topic name cannot be empty")
	}
	if len(name) > 256 {
		return fmt.Errorf("topic name cannot exceed 256 characters")
	}
	// Check for valid characters
	for _, c := range name {
		if !((c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '-' || c == '_') {
			return fmt.Errorf("topic name contains invalid character: %c", c)
		}
	}
	return nil
}
