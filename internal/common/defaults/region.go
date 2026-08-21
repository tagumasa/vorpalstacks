// Package defaults provides shared default values for vorpalstacks services.
package defaults

import "net/http"

// DefaultRegion is the default AWS region used when none is specified.
const DefaultRegion = "us-east-1"

// GetRegionFromHeader extracts the AWS region from the X-Aws-Region
// header, falling back to the default region if the header is absent.
func GetRegionFromHeader(headers http.Header) string {
	region := headers.Get("X-Aws-Region")
	if region == "" {
		region = DefaultRegion
	}
	return region
}
