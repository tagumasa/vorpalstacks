package generators

// Package generators provides ID generation utilities for vorpalstacks AWS resources.

// GenerateIAMRoleID generates a unique IAM role ID.
//
// Returns:
//   - string: The generated IAM role ID
//   - error: An error if ID generation fails
func GenerateIAMRoleID() (string, error) {
	return generateIDWithPrefix("ARO", 21)
}
