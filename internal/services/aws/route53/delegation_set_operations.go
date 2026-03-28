package route53

import (
	"context"

	"vorpalstacks/internal/services/aws/common/protocol"
	"vorpalstacks/internal/services/aws/common/request"
)

func (s *Route53Service) ListReusableDelegationSets(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return map[string]interface{}{
		"DelegationSets": protocol.XMLElements{ElementName: "DelegationSet", Items: []interface{}{}},
	}, nil
}
