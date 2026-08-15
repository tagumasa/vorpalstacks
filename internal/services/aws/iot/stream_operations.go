package iot

import (
	"context"
	"time"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/tags"
	iotstore "vorpalstacks/internal/store/aws/iot"
)

// ---------------------------------------------------------------------------
// Stream operations (MQTT-based file delivery).
// Streams are lightweight metadata records keyed by streamId. The actual
// file payload delivery happens over MQTT; the control-plane API manages
// the stream catalog and versioning only.
// ---------------------------------------------------------------------------

func (s *IoTService) CreateStream(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	streamID := request.GetParamCaseInsensitive(req.Parameters, "streamId")
	if streamID == "" {
		return nil, iotstore.ErrMissingParam
	}
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	// A streamId is the primary key of the stream catalog: recreating an
	// existing id must not overwrite the record (resetting its version to
	// one), matching the ResourceAlreadyExistsException the AWS API
	// documents for CreateStream.
	exists, err := store.GetGenericExists("stream/"+streamID, &map[string]interface{}{})
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, iotstore.ErrStreamAlreadyExists
	}
	tagList := tags.ParseTagsWithQueryFallback(req.Parameters, "tags")
	recTags := make(map[string]string, len(tagList))
	for _, t := range tagList {
		recTags[t.Key] = t.Value
	}
	now := time.Now().UTC().Unix()
	rec := map[string]interface{}{
		"streamId":      streamID,
		"streamArn":     iotstore.BuildStreamARN(reqCtx.GetAccountID(), reqCtx.GetRegion(), streamID),
		"streamVersion": int64(1),
		"description":   request.GetParamCaseInsensitive(req.Parameters, "description"),
		"files":         req.Parameters["files"],
		"roleArn":       request.GetParamCaseInsensitive(req.Parameters, "roleArn"),
		"tags":          recTags,
		"createdAt":     now,
		"lastUpdatedAt": now,
	}
	if err := store.PutGeneric("stream/"+streamID, rec); err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"streamId":      streamID,
		"streamArn":     rec["streamArn"],
		"description":   rec["description"],
		"streamVersion": rec["streamVersion"],
	}, nil
}

func (s *IoTService) DeleteStream(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	streamID := request.GetParamCaseInsensitive(req.Parameters, "streamId")
	if streamID == "" {
		return nil, iotstore.ErrMissingParam
	}
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	exists, err := store.GetGenericExists("stream/"+streamID, &map[string]interface{}{})
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, iotstore.ErrStreamNotFound
	}
	if err := store.DeleteGeneric("stream/" + streamID); err != nil {
		return nil, err
	}
	return map[string]interface{}{}, nil
}

func (s *IoTService) DescribeStream(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	streamID := request.GetParamCaseInsensitive(req.Parameters, "streamId")
	if streamID == "" {
		return nil, iotstore.ErrMissingParam
	}
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	rec := map[string]interface{}{}
	exists, err := store.GetGenericExists("stream/"+streamID, &rec)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, iotstore.ErrStreamNotFound
	}
	return map[string]interface{}{"streamInfo": rec}, nil
}

func (s *IoTService) ListStreams(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	items, err := store.ListGeneric("stream/")
	if err != nil {
		return nil, err
	}
	return paginatedMaps("streams", items, req.Parameters), nil
}

func (s *IoTService) UpdateStream(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	streamID := request.GetParamCaseInsensitive(req.Parameters, "streamId")
	if streamID == "" {
		return nil, iotstore.ErrMissingParam
	}
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	rec := map[string]interface{}{}
	exists, err := store.GetGenericExists("stream/"+streamID, &rec)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, iotstore.ErrStreamNotFound
	}
	if desc := request.GetParamCaseInsensitive(req.Parameters, "description"); desc != "" {
		rec["description"] = desc
	}
	if files, ok := req.Parameters["files"]; ok {
		rec["files"] = files
	}
	if roleArn := request.GetParamCaseInsensitive(req.Parameters, "roleArn"); roleArn != "" {
		rec["roleArn"] = roleArn
	}
	if v, ok := rec["streamVersion"].(int64); ok {
		rec["streamVersion"] = v + 1
	} else if v, ok := rec["streamVersion"].(float64); ok {
		rec["streamVersion"] = int64(v) + 1
	} else {
		rec["streamVersion"] = int64(2)
	}
	rec["lastUpdatedAt"] = time.Now().UTC().Unix()
	if err := store.PutGeneric("stream/"+streamID, rec); err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"streamId":      streamID,
		"streamArn":     rec["streamArn"],
		"description":   rec["description"],
		"streamVersion": rec["streamVersion"],
	}, nil
}
