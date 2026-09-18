package cloudtrail

import (
	"strings"

	svcarn "vorpalstacks/internal/utils/aws/arn"
)

// isCloudTrailResourceARN reports whether value parses as a CloudTrail ARN
// of the given resource form ("trail/", "eventdatastore/", "channel/").
// The DeleteTrail reference words the shared trigger: "This exception is
// thrown when an operation is called with an ARN that is not valid", and
// lists the trail, event data store, dashboard and channel ARN formats —
// an ARN-prefixed selector for one resource family must carry that
// family's resource form.
func isCloudTrailResourceARN(value, resourcePrefix string) bool {
	parsed, err := svcarn.ParseARN(value)
	return err == nil && parsed.Service == "cloudtrail" && strings.HasPrefix(parsed.Resource, resourcePrefix)
}
