// Package validators provides AWS-specific utility functions for vorpalstacks.
package validators

import (
	"strings"
)

// IsAWSAlias returns true if this is an AWS-managed alias.
//
// Parameters:
//   - aliasName: The alias name to check
//
// Returns:
//   - bool: True if this is an AWS-managed alias
func IsAWSAlias(aliasName string) bool {
	return strings.HasPrefix(aliasName, "alias/aws/")
}
