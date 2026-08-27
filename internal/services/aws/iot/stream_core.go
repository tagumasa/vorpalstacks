package iot

import (
	"time"

	iotstore "vorpalstacks/internal/store/aws/iot"
)

// ---------------------------------------------------------------------------
// Stream Core (MQTT-based file delivery).
// Streams are lightweight metadata records keyed by streamId under the
// "stream/<streamId>" key. The actual file payload delivery happens over
// MQTT; the control-plane API manages the stream catalog and versioning
// only.
// ---------------------------------------------------------------------------

// CreateStreamInput carries the fields for CreateStream. Files keeps the
// raw wire value (nested structure).
type CreateStreamInput struct {
	StreamID    string
	Description string
	Files       interface{}
	RoleArn     string
	Tags        map[string]string
}

// StreamMutationResult is the transport-agnostic result of the stream
// create and update operations.
type StreamMutationResult struct {
	StreamID      string
	StreamArn     interface{}
	Description   interface{}
	StreamVersion interface{}
}

// UpdateStreamInput carries the fields for UpdateStream. FilesProvided
// distinguishes an explicitly supplied files member from an omitted one.
type UpdateStreamInput struct {
	StreamID      string
	Description   string
	Files         interface{}
	FilesProvided bool
	RoleArn       string
}

// createStreamCore validates and persists a stream record. A streamId is
// the primary key of the stream catalog: recreating an existing id must not
// overwrite the record (resetting its version to one), matching the
// ResourceAlreadyExistsException the AWS API documents for CreateStream.
func (s *IoTService) createStreamCore(store iotstore.IotStoreInterface, in CreateStreamInput) (*StreamMutationResult, error) {
	if in.StreamID == "" {
		return nil, iotstore.ErrMissingParam
	}
	exists, err := store.GetGenericExists("stream/"+in.StreamID, &map[string]interface{}{})
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, iotstore.ErrStreamAlreadyExists
	}
	now := time.Now().UTC().Unix()
	rec := map[string]interface{}{
		"streamId":      in.StreamID,
		"streamArn":     iotstore.BuildStreamARN(store.GetAccountID(), store.GetRegion(), in.StreamID),
		"streamVersion": int64(1),
		"description":   in.Description,
		"files":         in.Files,
		"roleArn":       in.RoleArn,
		"tags":          in.Tags,
		"createdAt":     now,
		"lastUpdatedAt": now,
	}
	if err := store.PutGeneric("stream/"+in.StreamID, rec); err != nil {
		return nil, err
	}
	return &StreamMutationResult{
		StreamID:      in.StreamID,
		StreamArn:     rec["streamArn"],
		Description:   rec["description"],
		StreamVersion: rec["streamVersion"],
	}, nil
}

// deleteStreamCore removes a stream record.
func (s *IoTService) deleteStreamCore(store iotstore.IotStoreInterface, streamID string) error {
	if streamID == "" {
		return iotstore.ErrMissingParam
	}
	exists, err := store.GetGenericExists("stream/"+streamID, &map[string]interface{}{})
	if err != nil {
		return err
	}
	if !exists {
		return iotstore.ErrStreamNotFound
	}
	return store.DeleteGeneric("stream/" + streamID)
}

// describeStreamCore retrieves a stream record wrapped in the streamInfo
// response member.
func (s *IoTService) describeStreamCore(store iotstore.IotStoreInterface, streamID string) (map[string]interface{}, error) {
	if streamID == "" {
		return nil, iotstore.ErrMissingParam
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

// listStreamsCore lists all stream records.
func (s *IoTService) listStreamsCore(store iotstore.IotStoreInterface) ([]map[string]interface{}, error) {
	return store.ListGeneric("stream/")
}

// updateStreamCore applies the supplied fields to an existing stream and
// bumps its version.
func (s *IoTService) updateStreamCore(store iotstore.IotStoreInterface, in UpdateStreamInput) (*StreamMutationResult, error) {
	if in.StreamID == "" {
		return nil, iotstore.ErrMissingParam
	}
	rec := map[string]interface{}{}
	exists, err := store.GetGenericExists("stream/"+in.StreamID, &rec)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, iotstore.ErrStreamNotFound
	}
	if in.Description != "" {
		rec["description"] = in.Description
	}
	if in.FilesProvided {
		rec["files"] = in.Files
	}
	if in.RoleArn != "" {
		rec["roleArn"] = in.RoleArn
	}
	if v, ok := rec["streamVersion"].(int64); ok {
		rec["streamVersion"] = v + 1
	} else if v, ok := rec["streamVersion"].(float64); ok {
		rec["streamVersion"] = int64(v) + 1
	} else {
		rec["streamVersion"] = int64(2)
	}
	rec["lastUpdatedAt"] = time.Now().UTC().Unix()
	if err := store.PutGeneric("stream/"+in.StreamID, rec); err != nil {
		return nil, err
	}
	return &StreamMutationResult{
		StreamID:      in.StreamID,
		StreamArn:     rec["streamArn"],
		Description:   rec["description"],
		StreamVersion: rec["streamVersion"],
	}, nil
}
