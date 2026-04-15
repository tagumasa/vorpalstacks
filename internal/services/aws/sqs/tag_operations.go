package sqs

import (
	"context"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	tagutil "vorpalstacks/internal/common/tags"
)

// TagQueue adds tags to an SQS queue.
func (s *SQSService) TagQueue(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	queueURL := tagutil.GetResourceKey(req.Parameters, tagutil.SQSConfig)
	if queueURL == "" {
		return nil, ErrMissingParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if _, err := store.GetQueue(queueURL); err != nil {
		return nil, convertStoreError(err)
	}

	tags := tagutil.GetTags(req.Parameters, tagutil.SQSConfig)
	if len(tags) > 0 {
		if err := store.TagQueue(queueURL, tagutil.ToMap(tags)); err != nil {
			return nil, convertStoreError(err)
		}
	}

	return response.EmptyResponse(), nil
}

// UntagQueue removes tags from an SQS queue.
func (s *SQSService) UntagQueue(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	queueURL := tagutil.GetResourceKey(req.Parameters, tagutil.SQSConfig)
	if queueURL == "" {
		return nil, ErrMissingParameter
	}

	tagKeys := tagutil.GetTagKeys(req.Parameters, tagutil.SQSConfig)
	if len(tagKeys) > 0 {
		store, err := s.store(reqCtx)
		if err != nil {
			return nil, err
		}
		if err := store.UntagQueue(queueURL, tagKeys); err != nil {
			return nil, convertStoreError(err)
		}
	}

	return response.EmptyResponse(), nil
}

// ListQueueTags lists tags for an SQS queue.
func (s *SQSService) ListQueueTags(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	queueURL := tagutil.GetResourceKey(req.Parameters, tagutil.SQSConfig)
	if queueURL == "" {
		return nil, ErrMissingParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	tags, err := store.ListQueueTags(queueURL)
	if err != nil {
		return nil, convertStoreError(err)
	}

	return map[string]interface{}{
		"Tags": tags,
	}, nil
}
