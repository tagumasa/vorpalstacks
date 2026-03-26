// Package validators provides AWS-specific utility functions for vorpalstacks.
package validators

import (
	"fmt"
)

// ValidateKMSPolicyName validates a KMS policy name.
//
// Parameters:
//   - policyName: The policy name to validate
//
// Returns:
//   - error: An error if validation fails, nil otherwise
func ValidateKMSPolicyName(policyName string) error {
	if policyName == "" {
		return fmt.Errorf("policy name cannot be empty")
	}
	if policyName != "default" {
		return fmt.Errorf("only 'default' policy name is supported")
	}
	return nil
}
