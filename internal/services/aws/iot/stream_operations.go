package iot

import (
	"context"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/tags"
)

// ---------------------------------------------------------------------------
// Stream operations (MQTT-based file delivery).
// Streams are lightweight metadata records keyed by streamId. The actual
// file payload delivery happens over MQTT; the control-plane API manages
// the stream catalog and versioning only.
// ---------------------------------------------------------------------------

func (s *IoTService) CreateStream(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	tagList := tags.ParseTagsWithQueryFallback(req.Parameters, "tags")
	recTags := make(map[string]string, len(tagList))
	for _, t := range tagList {
		recTags[t.Key] = t.Value
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	result, err := s.createStreamCore(store, CreateStreamInput{
		StreamID:    request.GetParamCaseInsensitive(req.Parameters, "streamId"),
		Description: request.GetParamCaseInsensitive(req.Parameters, "description"),
		Files:       req.Parameters["files"],
		RoleArn:     request.GetParamCaseInsensitive(req.Parameters, "roleArn"),
		Tags:        recTags,
	})
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"streamId":      result.StreamID,
		"streamArn":     result.StreamArn,
		"description":   result.Description,
		"streamVersion": result.StreamVersion,
	}, nil
}

func (s *IoTService) DeleteStream(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := s.deleteStreamCore(store, request.GetParamCaseInsensitive(req.Parameters, "streamId")); err != nil {
		return nil, err
	}
	return map[string]interface{}{}, nil
}

func (s *IoTService) DescribeStream(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	return s.describeStreamCore(store, request.GetParamCaseInsensitive(req.Parameters, "streamId"))
}

func (s *IoTService) ListStreams(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	items, err := s.listStreamsCore(store)
	if err != nil {
		return nil, err
	}
	return paginatedMaps("streams", items, req.Parameters)
}

func (s *IoTService) UpdateStream(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	files, filesProvided := req.Parameters["files"]
	result, err := s.updateStreamCore(store, UpdateStreamInput{
		StreamID:      request.GetParamCaseInsensitive(req.Parameters, "streamId"),
		Description:   request.GetParamCaseInsensitive(req.Parameters, "description"),
		Files:         files,
		FilesProvided: filesProvided,
		RoleArn:       request.GetParamCaseInsensitive(req.Parameters, "roleArn"),
	})
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"streamId":      result.StreamID,
		"streamArn":     result.StreamArn,
		"description":   result.Description,
		"streamVersion": result.StreamVersion,
	}, nil
}
