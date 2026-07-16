package generators

// Package generators provides ID generation utilities for vorpalstacks AWS resources.

// GenerateACMCertificateSerial generates a unique ACM certificate serial number.
//
// Returns:
//   - string: The generated certificate serial
//   - error: An error if ID generation fails
func GenerateACMCertificateSerial() (string, error) {
	return generateHexID()
}
