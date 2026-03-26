// Package validators provides AWS-specific utility functions for vorpalstacks.
package validators

// ValidateLogGroupName validates a CloudWatch log group name.
//
// Parameters:
//   - name: The log group name to validate
//
// Returns:
//   - error: An error if validation fails, nil otherwise
//
// AWS CloudWatch log group name rules:
//   - Must be 1-512 characters
//   - Can contain alphanumeric characters, hyphens, underscores, slashes, and colons
func ValidateLogGroupName(name string) error {
	return validateLogName(name, "log group", 512, "-_/:")
}
