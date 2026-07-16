package generators

// Package generators provides ID generation utilities for vorpalstacks AWS resources.

// GenerateCloudFrontDistributionID generates a unique CloudFront distribution ID.
//
// Returns:
//   - string: The generated distribution ID
//   - error: An error if ID generation fails
func GenerateCloudFrontDistributionID() (string, error) {
	return generateIDWithPrefix("E", 28)
}
