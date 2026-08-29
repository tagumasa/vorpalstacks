package kinesis

import (
	"context"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/tags"
)

// AddTagsToStream adds or overwrites tags on a Kinesis data stream.
func (s *KinesisService) AddTagsToStream(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	return tags.HandleTag(ctx, req, s.kinesisTagConfig(store, req))
}

// RemoveTagsFromStream removes the specified tags from a Kinesis data stream.
func (s *KinesisService) RemoveTagsFromStream(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	return tags.HandleUntag(ctx, req, s.kinesisTagConfig(store, req))
}

// ListTagsForStream lists all tags assigned to a Kinesis data stream.
func (s *KinesisService) ListTagsForStream(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	return tags.HandleList(ctx, req, s.kinesisTagConfig(store, req))
}

// TagResource adds or overwrites tags on a Kinesis resource identified by ARN.
func (s *KinesisService) TagResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	return tags.HandleTag(ctx, req, s.kinesisTagConfig(store, req))
}

// UntagResource removes the specified tags from a Kinesis resource identified by ARN.
func (s *KinesisService) UntagResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	return tags.HandleUntag(ctx, req, s.kinesisTagConfig(store, req))
}

// ListTagsForResource lists all tags assigned to a Kinesis resource identified by ARN.
func (s *KinesisService) ListTagsForResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	return tags.HandleList(ctx, req, s.kinesisTagConfig(store, req))
}
