package generators

// Package generators provides ID generation utilities for vorpalstacks AWS resources.

// GenerateIAMUserID generates a unique IAM user ID.
//
// Returns:
//   - string: The generated IAM user ID
//   - error: An error if ID generation fails
func GenerateIAMUserID() (string, error) {
	return generateIDWithPrefix("AID", 21)
}
