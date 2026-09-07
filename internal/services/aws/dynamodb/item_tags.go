package dynamodb

import (
	"context"

	"vorpalstacks/internal/common/request"
)

// TagResource adds or overwrites tags on a DynamoDB table.
func (s *DynamoDBService) TagResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	return s.tagResourceCore(ctx, req, store)
}

// UntagResource removes the specified tags from a DynamoDB table.
func (s *DynamoDBService) UntagResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	return s.untagResourceCore(ctx, req, store)
}

// ListTagsForResource lists tags assigned to a DynamoDB table with pagination.
// AWS paginates tags; the marker is the tag key of the last returned entry.
func (s *DynamoDBService) ListTagsForResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	return s.listTagsForResourceCore(ctx, req, store)
}
