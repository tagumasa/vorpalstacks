package generators

// Package generators provides ID generation utilities for vorpalstacks AWS resources.

// GenerateIAMGroupID generates a unique IAM group ID.
//
// Returns:
//   - string: The generated IAM group ID
//   - error: An error if ID generation fails
func GenerateIAMGroupID() (string, error) {
	return generateIDWithPrefix("AGP", 21)
}
