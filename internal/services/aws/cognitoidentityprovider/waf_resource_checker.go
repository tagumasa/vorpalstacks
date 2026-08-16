package cognitoidentityprovider

import (
	"context"
	"strings"

	"vorpalstacks/internal/utils/aws/arn"
)

// WebACLResourceService returns the ARN service namespace whose resources
// this service hosts for WAFv2 WebACL association.
func (s *CognitoService) WebACLResourceService() string {
	return "cognito-idp"
}

// WebACLResourceExists reports whether the Cognito user pool referenced by
// the ARN exists. User pool ARNs take the form
// arn:<partition>:cognito-idp:<region>:<account>:userpool/<user-pool-id>.
func (s *CognitoService) WebACLResourceExists(ctx context.Context, region, resourceArn string) bool {
	parsed, err := arn.ParseARN(resourceArn)
	if err != nil {
		return false
	}
	poolID := strings.TrimPrefix(parsed.Resource, "userpool/")
	if poolID == "" || poolID == parsed.Resource {
		return false
	}
	store, err := s.GetStoreForRegion(region)
	if err != nil {
		return false
	}
	_, err = store.GetUserPool(poolID)
	return err == nil
}
