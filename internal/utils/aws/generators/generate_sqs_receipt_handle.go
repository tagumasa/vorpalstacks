package generators

// Package generators provides ID generation utilities for vorpalstacks AWS resources.

// GenerateSQSReceiptHandle generates a unique SQS receipt handle.
//
// Returns:
//   - string: The generated receipt handle
//   - error: An error if ID generation fails
func GenerateSQSReceiptHandle() (string, error) {
	return generateIDWithPrefix("RECEIPT", 32)
}
