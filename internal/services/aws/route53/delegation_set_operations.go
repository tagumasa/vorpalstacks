package route53

import (
	"context"

	"vorpalstacks/internal/common/protocol"
	"vorpalstacks/internal/common/request"
)

// ListReusableDelegationSets returns all reusable delegation sets in the account.
func (s *Route53Service) ListReusableDelegationSets(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return map[string]interface{}{
		"DelegationSets": protocol.XMLElements{ElementName: "DelegationSet", Items: []interface{}{}},
	}, nil
}
