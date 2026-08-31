package sqs

import (
	"context"

	"vorpalstacks/internal/common/request"
	tagutil "vorpalstacks/internal/common/tags"
)

// TagQueue adds or overwrites tags on an SQS queue.
func (s *SQSService) TagQueue(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	return tagutil.HandleTag(ctx, req, sqsTagConfig(store))
}

// UntagQueue removes the specified tags from an SQS queue.
func (s *SQSService) UntagQueue(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	return tagutil.HandleUntag(ctx, req, sqsTagConfig(store))
}

// ListQueueTags lists all tags assigned to an SQS queue.
func (s *SQSService) ListQueueTags(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	return tagutil.HandleList(ctx, req, sqsTagConfig(store))
}
