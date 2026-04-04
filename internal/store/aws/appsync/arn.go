package appsync

import (
	"strings"

	"vorpalstacks/internal/utils/aws/arn"
)

// extractApiSuffix extracts the resource identifier from an AppSync ARN
// by stripping the given prefix from the resource portion.
func extractApiSuffix(arnStr, prefix string) string {
	_, _, _, _, resource := arn.SplitARN(arnStr)
	if strings.HasPrefix(resource, prefix) {
		return strings.TrimPrefix(resource, prefix)
	}
	return ""
}
