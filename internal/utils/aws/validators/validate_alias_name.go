// Package validators provides AWS-specific utility functions for vorpalstacks.
package validators

import (
	"fmt"
	"strings"
)

// ValidateAliasName validates a KMS alias name according to AWS KMS rules.
// KMS alias names must start with "alias/" and can contain alphanumeric
// characters, underscores, and hyphens.
//
// Parameters:
//   - aliasName: The alias name to validate
//
// Returns:
//   - error: An error if validation fails, nil if valid
//
// AWS KMS alias name rules:
//   - Must start with "alias/"
//   - Must be between 7 and 256 characters
//   - Can only contain alphanumeric characters, underscores, and hyphens
//
// Example:
//
//	err := ValidateAliasName("alias/my-key")
func ValidateAliasName(aliasName string) error {
	if !strings.HasPrefix(aliasName, "alias/") {
		return fmt.Errorf("alias name must start with 'alias/'")
	}
	if len(aliasName) < 7 || len(aliasName) > 256 {
		return fmt.Errorf("alias name must be between 7 and 256 characters")
	}
	if !aliasNameRegex.MatchString(aliasName) {
		return fmt.Errorf("alias name can only contain alphanumeric characters, underscores, and hyphens")
	}
	return nil
}
