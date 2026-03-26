package generators

// Package generators provides ID generation utilities for vorpalstacks AWS resources.

// GenerateIAMPolicyID generates a unique IAM policy ID.
//
// Returns:
//   - string: The generated IAM policy ID
//   - error: An error if ID generation fails
func GenerateIAMPolicyID() (string, error) {
	return generateIDWithPrefix("ANP", 21)
}
