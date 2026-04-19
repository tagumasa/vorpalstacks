package common

import (
	"net/http"

	"vorpalstacks/internal/common/request"
)

func GetRegionFromHeader(headers http.Header) string {
	region := headers.Get("X-Aws-Region")
	if region == "" {
		region = request.DefaultRegion
	}
	return region
}
