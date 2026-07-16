// Package validators provides AWS-specific validation utility functions for vorpalstacks.
package validators

import (
	"fmt"
	"strings"
)

// ValidateBucketName validates an S3 bucket name according to AWS S3 naming rules.
//
// Parameters:
//   - name: The bucket name to validate
//
// Returns:
//   - error: An error if validation fails, nil otherwise
//
// AWS S3 bucket name rules:
//   - Must be 3-63 characters
//   - Can contain lowercase letters, numbers, hyphens, and periods
//   - Must start and end with a letter or number
//   - Must not be formatted as an IP address
//   - Must not contain consecutive periods
//   - Must not contain a period adjacent to a hyphen
func ValidateBucketName(name string) error {
	if len(name) < 3 || len(name) > 63 {
		return fmt.Errorf("bucket name must be between 3 and 63 characters")
	}

	// Check for valid characters
	for _, c := range name {
		if !((c >= 'a' && c <= 'z') || (c >= '0' && c <= '9') || c == '-' || c == '.') {
			return fmt.Errorf("bucket name contains invalid characters")
		}
	}

	// Check for IP address format
	if IsIPAddress(name) {
		return fmt.Errorf("bucket name cannot be formatted as an IP address")
	}

	// Check for consecutive periods
	if strings.Contains(name, "..") {
		return fmt.Errorf("bucket name cannot contain consecutive periods")
	}

	// Check for period adjacent to hyphen
	if strings.Contains(name, "-.") || strings.Contains(name, ".-") {
		return fmt.Errorf("bucket name cannot have a period adjacent to a hyphen")
	}

	// Check start and end characters
	if name[0] == '-' || name[0] == '.' {
		return fmt.Errorf("bucket name must start with a letter or number")
	}
	if name[len(name)-1] == '-' || name[len(name)-1] == '.' {
		return fmt.Errorf("bucket name must end with a letter or number")
	}

	return nil
}
