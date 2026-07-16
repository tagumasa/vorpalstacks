// Package validators provides AWS-specific utility functions for vorpalstacks.
package validators

// ValidateRoleName validates an IAM role name.
//
// Parameters:
//   - roleName: The role name to validate
//
// Returns:
//   - error: An error if validation fails, nil otherwise
func ValidateRoleName(roleName string) error {
	return validateIAMName(roleName, 64, "role name")
}
