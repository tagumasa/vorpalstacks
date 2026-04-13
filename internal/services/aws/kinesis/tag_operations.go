package kinesis

import (
	"context"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	"vorpalstacks/internal/common/tags"
)

// AddTagsToStream adds tags to a Kinesis stream.
// Populates a tag map from the request parameters and calls the store to tag the stream.
func (s *KinesisService) AddTagsToStream(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	streamName := request.GetParamLowerFirst(req.Parameters, "StreamName")
	streamARN := request.GetParamLowerFirst(req.Parameters, "StreamARN")

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if streamARN != "" {
		stream, err := store.GetStreamByARN(streamARN)
		if err != nil {
			return nil, s.mapStoreError(err)
		}
		streamName = stream.StreamName
	}

	if streamName == "" {
		return nil, ErrInvalidArgument
	}

	tagList := tags.ParseTags(req.Parameters, "Tags")
	tagMap := make(map[string]string)
	for _, t := range tagList {
		tagMap[t.Key] = t.Value
	}

	if err := store.TagResource(streamName, tagMap); err != nil {
		return nil, s.mapStoreError(err)
	}

	return response.EmptyResponse(), nil
}

// RemoveTagsFromStream removes tags from a Kinesis stream.
// Extracts tag keys from the request and calls the store to untag the stream.
func (s *KinesisService) RemoveTagsFromStream(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	streamName := request.GetParamLowerFirst(req.Parameters, "StreamName")
	streamARN := request.GetParamLowerFirst(req.Parameters, "StreamARN")

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if streamARN != "" {
		stream, err := store.GetStreamByARN(streamARN)
		if err != nil {
			return nil, s.mapStoreError(err)
		}
		streamName = stream.StreamName
	}

	if streamName == "" {
		return nil, ErrInvalidArgument
	}

	tagKeys := tags.ParseTagKeysAsSlice(req.Parameters, "TagKeys")

	if err := store.UntagResource(streamName, tagKeys); err != nil {
		return nil, s.mapStoreError(err)
	}

	return response.EmptyResponse(), nil
}

// ListTagsForStream lists the tags for a Kinesis stream.
// Retrieves tags from the store and returns them as a formatted slice.
func (s *KinesisService) ListTagsForStream(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	streamName := request.GetParamLowerFirst(req.Parameters, "StreamName")
	streamARN := request.GetParamLowerFirst(req.Parameters, "StreamARN")

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if streamARN != "" {
		stream, err := store.GetStreamByARN(streamARN)
		if err != nil {
			return nil, s.mapStoreError(err)
		}
		streamName = stream.StreamName
	}

	if streamName == "" {
		return nil, ErrInvalidArgument
	}

	tagSlice, err := store.ListTagsAsSlice(streamName)
	if err != nil {
		return nil, s.mapStoreError(err)
	}

	var formattedTags []map[string]interface{}
	for _, t := range tagSlice {
		formattedTags = append(formattedTags, map[string]interface{}{
			"Key":   t.Key,
			"Value": t.Value,
		})
	}

	return map[string]interface{}{
		"Tags":        formattedTags,
		"HasMoreTags": false,
	}, nil
}

// TagResource adds tags to a Kinesis stream specified by ARN.
// Extracts the resource ARN and tag list from the request.
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
	tagMap := make(map[string]string)
	for _, t := range tagList {
		tagMap[t.Key] = t.Value
	}

	if err := store.TagResource(stream.StreamName, tagMap); err != nil {
		return nil, s.mapStoreError(err)
	}

	return response.EmptyResponse(), nil
}

// UntagResource removes tags from a Kinesis stream specified by ARN.
// Extracts the resource ARN and tag keys from the request.
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

	if err := store.UntagResource(stream.StreamName, tagKeys); err != nil {
		return nil, s.mapStoreError(err)
	}

	return response.EmptyResponse(), nil
}

// ListTagsForResource lists the tags for a Kinesis stream specified by ARN.
// Retrieves tags from the store using the stream ARN and returns them as a formatted slice.
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

	tagSlice, err := store.ListTagsAsSlice(stream.StreamName)
	if err != nil {
		return nil, s.mapStoreError(err)
	}

	var formattedTags []map[string]interface{}
	for _, t := range tagSlice {
		formattedTags = append(formattedTags, map[string]interface{}{
			"Key":   t.Key,
			"Value": t.Value,
		})
	}

	return map[string]interface{}{
		"Tags": formattedTags,
	}, nil
}
