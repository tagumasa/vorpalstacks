package kinesis

import (
	"context"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	"vorpalstacks/internal/common/tags"
	kinesisstore "vorpalstacks/internal/store/aws/kinesis"
)

func (s *KinesisService) resolveStreamName(store *kinesisstore.KinesisStore, req *request.ParsedRequest) (string, error) {
	name := request.GetParamLowerFirst(req.Parameters, "StreamName")
	if name != "" {
		return name, nil
	}

	arn := request.GetParamLowerFirst(req.Parameters, "StreamARN")
	if arn == "" {
		arn = request.GetParamLowerFirst(req.Parameters, "ResourceARN")
	}
	if arn == "" {
		return "", nil
	}

	stream, err := store.GetStreamByARN(arn)
	if err != nil {
		return "", s.mapStoreError(err)
	}
	return stream.StreamName, nil
}

// AddTagsToStream adds or overwrites tags on a Kinesis data stream.
func (s *KinesisService) AddTagsToStream(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	streamName, err := s.resolveStreamName(store, req)
	if err != nil {
		return nil, err
	}
	if streamName == "" {
		return nil, ErrInvalidArgument
	}

	tagList := tags.ParseTags(req.Parameters, "Tags")
	tagMap := tags.ToMap(tagList)

	if err := store.Tag(streamName, tagMap); err != nil {
		return nil, s.mapStoreError(err)
	}

	return response.EmptyResponse(), nil
}

// RemoveTagsFromStream removes the specified tags from a Kinesis data stream.
func (s *KinesisService) RemoveTagsFromStream(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	streamName, err := s.resolveStreamName(store, req)
	if err != nil {
		return nil, err
	}
	if streamName == "" {
		return nil, ErrInvalidArgument
	}

	tagKeys := tags.ParseTagKeysAsSlice(req.Parameters, "TagKeys")

	if err := store.Untag(streamName, tagKeys); err != nil {
		return nil, s.mapStoreError(err)
	}

	return response.EmptyResponse(), nil
}

// ListTagsForStream lists all tags assigned to a Kinesis data stream.
func (s *KinesisService) ListTagsForStream(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	streamName, err := s.resolveStreamName(store, req)
	if err != nil {
		return nil, err
	}
	if streamName == "" {
		return nil, ErrInvalidArgument
	}

	tagSlice, err := store.ListAsSlice(streamName)
	if err != nil {
		return nil, s.mapStoreError(err)
	}

	return map[string]interface{}{
		"Tags":        tags.ToResponse(tagSlice),
		"HasMoreTags": false,
	}, nil
}

// TagResource adds or overwrites tags on a Kinesis resource identified by ARN.
func (s *KinesisService) TagResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	resourceARN := request.GetParamLowerFirst(req.Parameters, "ResourceARN")

	if resourceARN == "" {
		return nil, ErrInvalidArgument
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	stream, err := store.GetStreamByARN(resourceARN)
	if err != nil {
		return nil, s.mapStoreError(err)
	}

	tagList := tags.ParseTags(req.Parameters, "Tags")
	tagMap := tags.ToMap(tagList)

	if err := store.Tag(stream.StreamName, tagMap); err != nil {
		return nil, s.mapStoreError(err)
	}

	return response.EmptyResponse(), nil
}

// UntagResource removes the specified tags from a Kinesis resource identified by ARN.
func (s *KinesisService) UntagResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	resourceARN := request.GetParamLowerFirst(req.Parameters, "ResourceARN")

	if resourceARN == "" {
		return nil, ErrInvalidArgument
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	stream, err := store.GetStreamByARN(resourceARN)
	if err != nil {
		return nil, s.mapStoreError(err)
	}

	tagKeys := tags.ParseTagKeysAsSlice(req.Parameters, "TagKeys")

	if err := store.Untag(stream.StreamName, tagKeys); err != nil {
		return nil, s.mapStoreError(err)
	}

	return response.EmptyResponse(), nil
}

// ListTagsForResource lists all tags assigned to a Kinesis resource identified by ARN.
func (s *KinesisService) ListTagsForResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	resourceARN := request.GetParamLowerFirst(req.Parameters, "ResourceARN")

	if resourceARN == "" {
		return nil, ErrInvalidArgument
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	stream, err := store.GetStreamByARN(resourceARN)
	if err != nil {
		return nil, s.mapStoreError(err)
	}

	tagSlice, err := store.ListAsSlice(stream.StreamName)
	if err != nil {
		return nil, s.mapStoreError(err)
	}

	return map[string]interface{}{
		"Tags":        tags.ToResponse(tagSlice),
		"HasMoreTags": false,
	}, nil
}
