package generators

// Package generators provides ID generation utilities for vorpalstacks AWS resources.

// GenerateIAMAccessKeyID generates a unique IAM access key ID.
//
// Returns:
//   - string: The generated IAM access key ID
//   - error: An error if ID generation fails
func GenerateIAMAccessKeyID() (string, error) {
	return generateIDWithPrefix("AKIA", 20)
}
