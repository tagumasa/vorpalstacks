package appsync

import (
	"context"
	"strings"

	"vorpalstacks/internal/utils/aws/arn"
)

// WebACLResourceService returns the ARN service namespace whose resources
// this service hosts for WAFv2 WebACL association.
func (s *AppSyncService) WebACLResourceService() string {
	return "appsync"
}

// WebACLResourceExists reports whether the AppSync API referenced by the ARN
// exists. AppSync API ARNs take the form
// arn:<partition>:appsync:<region>:<account>:apis/<api-id>.
func (s *AppSyncService) WebACLResourceExists(ctx context.Context, region, resourceArn string) bool {
	parsed, err := arn.ParseARN(resourceArn)
	if err != nil {
		return false
	}
	apiId := strings.TrimPrefix(parsed.Resource, "apis/")
	if apiId == "" || apiId == parsed.Resource {
		return false
	}
	_, err = s.lookupStoreByApiId(apiId)
	return err == nil
}
