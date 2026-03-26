// Package validators provides AWS-specific utility functions for vorpalstacks.
package validators

import (
	"fmt"
	"strings"
)

// ValidateHostedZoneName validates a Route53 hosted zone name.
//
// Parameters:
//   - name: The hosted zone name to validate
//
// Returns:
//   - error: An error if validation fails, nil otherwise
func ValidateHostedZoneName(name string) error {
	if name == "" {
		return fmt.Errorf("hosted zone name cannot be empty")
	}

	// Maximum length is 255 characters
	if len(name) > 255 {
		return fmt.Errorf("hosted zone name cannot exceed 255 characters")
	}

	// Check for valid characters (letters, numbers, hyphens, and dots)
	for _, c := range name {
		if !((c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '-' || c == '.') {
			return fmt.Errorf("hosted zone name contains invalid character: %c", c)
		}
	}

	// Must start and end with a letter or number
	if len(name) > 0 {
		if name[0] == '-' || name[0] == '.' {
			return fmt.Errorf("hosted zone name must start with a letter or number")
		}
		if name[len(name)-1] == '-' || name[len(name)-1] == '.' {
			return fmt.Errorf("hosted zone name must end with a letter or number")
		}
	}

	// Must not contain consecutive hyphens or dots
	if strings.Contains(name, "--") || strings.Contains(name, "..") {
		return fmt.Errorf("hosted zone name cannot contain consecutive hyphens or dots")
	}

	return nil
}
