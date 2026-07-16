package generators

// Package generators provides ID generation utilities for vorpalstacks AWS resources.

// GenerateACMCertificateID generates a unique ACM certificate ID.
//
// Returns:
//   - string: The generated certificate ID
//   - error: An error if ID generation fails
func GenerateACMCertificateID() (string, error) {
	return generateIDWithPrefix("CERT", 28)
}
