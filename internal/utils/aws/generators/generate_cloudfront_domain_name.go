package generators

// Package generators provides ID generation utilities for vorpalstacks AWS resources.

import "fmt"

// GenerateCloudFrontDomainName generates a CloudFront domain name for a distribution.
//
// Parameters:
//   - distributionID: The CloudFront distribution ID
//
// Returns:
//   - string: The generated domain name
//   - error: An error if distribution ID is empty
func GenerateCloudFrontDomainName(distributionID string) (string, error) {
	if distributionID == "" {
		return "", fmt.Errorf("distributionID cannot be empty")
	}
	return distributionID + ".cloudfront.net", nil
}
