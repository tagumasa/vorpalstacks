package common

import "net/http"

// GetRegionFromHeader extracts the AWS region from request headers.
// Defaults to "us-east-1" if not specified.
func GetRegionFromHeader(headers http.Header) string {
	region := headers.Get("X-Aws-Region")
	if region == "" {
		region = "us-east-1"
	}
	return region
}
