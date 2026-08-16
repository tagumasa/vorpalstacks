package apigateway

import (
	"context"
	"strings"

	"vorpalstacks/internal/utils/aws/arn"
)

// WebACLResourceService returns the ARN service namespace whose resources
// this service hosts for WAFv2 WebACL association.
func (s *APIGatewayService) WebACLResourceService() string {
	return "apigateway"
}

// WebACLResourceExists reports whether the API Gateway stage referenced by
// the ARN exists. API Gateway stage ARNs take the form
// arn:<partition>:apigateway:<region>::/restapis/<api-id>/stages/<stage>.
func (s *APIGatewayService) WebACLResourceExists(ctx context.Context, region, resourceArn string) bool {
	parsed, err := arn.ParseARN(resourceArn)
	if err != nil {
		return false
	}
	rest := strings.TrimPrefix(parsed.Resource, "/restapis/")
	parts := strings.Split(rest, "/stages/")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return false
	}
	apiId, stageName := parts[0], parts[1]

	stores, err := s.GetStoreForRegion(region)
	if err != nil {
		return false
	}
	if _, err := stores.restApis.Get(apiId); err != nil {
		return false
	}
	_, err = stores.restApis.GetStage(apiId, stageName)
	return err == nil
}
