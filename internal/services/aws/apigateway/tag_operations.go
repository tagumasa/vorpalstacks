package apigateway

import (
	"context"

	"vorpalstacks/internal/common/request"
	tagutil "vorpalstacks/internal/common/tags"
)

// TagResource adds tags to a resource in API Gateway.
func (s *APIGatewayService) TagResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	return tagutil.HandleTag(ctx, req, s.apiGatewayTagConfig(stores, req))
}

// UntagResource removes tags from a resource in API Gateway.
func (s *APIGatewayService) UntagResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	return tagutil.HandleUntag(ctx, req, s.apiGatewayTagConfig(stores, req))
}

// ListTagsForResource lists tags for a resource in API Gateway.
func (s *APIGatewayService) ListTagsForResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	return tagutil.HandleList(ctx, req, s.apiGatewayTagConfig(stores, req))
}
