package kinesis

import (
	"vorpalstacks/internal/common/request"
	kinesisstore "vorpalstacks/internal/store/aws/kinesis"
)

// resolveStreamName resolves a stream name from either StreamName or StreamARN
// parameters, returning the store and stream name. Returns ErrInvalidArgument
// if neither is provided.
func (s *KinesisService) resolveStreamName(reqCtx *request.RequestContext, params map[string]interface{}) (*kinesisstore.KinesisStore, string, error) {
	streamName := request.GetParamLowerFirst(params, "StreamName")
	streamARN := request.GetParamLowerFirst(params, "StreamARN")

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, "", err
	}

	if streamARN != "" {
		stream, err := store.GetStreamByARN(streamARN)
		if err != nil {
			return nil, "", s.mapStoreError(err)
		}
		streamName = stream.StreamName
	}

	if streamName == "" {
		return nil, "", ErrInvalidArgument
	}

	if !validateStreamName(streamName) {
		return nil, "", ErrInvalidArgument
	}

	return store, streamName, nil
}

// resolveStreamNameOptional resolves a stream name but does not require one
// (used by ListShards which accepts optional stream identification).
func (s *KinesisService) resolveStreamNameOptional(reqCtx *request.RequestContext, params map[string]interface{}) (*kinesisstore.KinesisStore, string, error) {
	streamName := request.GetParamLowerFirst(params, "StreamName")
	streamARN := request.GetParamLowerFirst(params, "StreamARN")

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, "", err
	}

	if streamARN != "" {
		stream, err := store.GetStreamByARN(streamARN)
		if err != nil {
			return nil, "", s.mapStoreError(err)
		}
		streamName = stream.StreamName
	}

	if streamName != "" && !validateStreamName(streamName) {
		return nil, "", ErrInvalidArgument
	}

	return store, streamName, nil
}

// parseStreamModeDetails extracts the stream mode from StreamModeDetails parameter.
func parseStreamModeDetails(params map[string]interface{}) kinesisstore.StreamMode {
	streamMode := kinesisstore.StreamModeProvisioned
	streamModeDetails := request.GetMapParam(params, "StreamModeDetails")
	if streamModeDetails == nil {
		streamModeDetails = request.GetMapParam(params, "streamModeDetails")
	}
	if streamModeDetails != nil {
		if v, ok := streamModeDetails["StreamMode"].(string); ok {
			streamMode = kinesisstore.StreamMode(v)
		} else if v, ok := streamModeDetails["streamMode"].(string); ok {
			streamMode = kinesisstore.StreamMode(v)
		}
	}
	return streamMode
}

// formatConsumer converts a Consumer to its API response map.
func formatConsumer(c *kinesisstore.Consumer) map[string]interface{} {
	return map[string]interface{}{
		"ConsumerName":              c.ConsumerName,
		"ConsumerARN":               c.ConsumerARN,
		"StreamARN":                 c.StreamARN,
		"ConsumerStatus":            c.ConsumerStatus,
		"ConsumerCreationTimestamp": float64(c.ConsumerCreationTimestamp.Unix()),
	}
}

// resolveEncryptionType returns the encryption type string, defaulting to "NONE".
func resolveEncryptionType(stream *kinesisstore.Stream) string {
	if stream.EncryptionType != "" {
		return stream.EncryptionType
	}
	return "NONE"
}
