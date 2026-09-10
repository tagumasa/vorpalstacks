package request

import (
	"strings"

	svcarn "vorpalstacks/internal/utils/aws/arn"
)

// TagsRouteService reports which service a /tags/{arn} path addresses by
// parsing the ARN embedded in the path. The bare prefix is shared by several
// REST services (API Gateway, EventBridge Scheduler, Neptune Graph), so the
// ARN service field is the only sound discriminator; a path without a
// parseable ARN is claimed by no service. The classifier routes /tags/
// paths through the same discrimination, so both layers agree on ownership
// for every partition and resource shape.
func TagsRouteService(path string) string {
	if !strings.HasPrefix(path, "/tags/") {
		return ""
	}
	parsed, err := svcarn.ParseARN(strings.TrimPrefix(path, "/tags/"))
	if err != nil {
		return ""
	}
	return parsed.Service
}
