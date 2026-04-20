package common

import (
	"net/http"

	"vorpalstacks/internal/common/request"
)

// GetRegionFromHeader extracts the AWS region from the X-Aws-Region header,
// falling back to the default region if the header is absent.
func GetRegionFromHeader(headers http.Header) string {
	region := headers.Get("X-Aws-Region")
	if region == "" {
		region = request.DefaultRegion
	}
	return region
}
