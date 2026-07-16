// Package validators provides AWS-specific utility functions for vorpalstacks.
package validators

// ValidatePolicyName validates an IAM policy name.
//
// Parameters:
//   - policyName: The policy name to validate
//
// Returns:
//   - error: An error if validation fails, nil otherwise
func ValidatePolicyName(policyName string) error {
	return validateIAMName(policyName, 128, "policy name")
}
