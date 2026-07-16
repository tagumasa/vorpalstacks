package generators

// Package generators provides ID generation utilities for vorpalstacks AWS resources.

// GenerateSQSMessageID generates a unique SQS message ID.
//
// Returns:
//   - string: The generated message ID
//   - error: An error if ID generation fails
func GenerateSQSMessageID() (string, error) {
	return generateIDWithPrefix("MSG", 28)
}
