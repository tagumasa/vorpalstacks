// Package validators provides AWS-specific utility functions for vorpalstacks.
package validators

// ValidateLogStreamName validates a CloudWatch log stream name.
//
// Parameters:
//   - name: The log stream name to validate
//
// Returns:
//   - error: An error if validation fails, nil otherwise
//
// AWS CloudWatch log stream name rules:
//   - Must be 1-512 characters
//   - Can contain alphanumeric characters, hyphens, underscores, slashes, and colons
func ValidateLogStreamName(name string) error {
	return validateLogName(name, "log stream", 512, "-_/#")
}
