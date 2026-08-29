package kinesis

import (
	"vorpalstacks/internal/common/request"
	kinesisstore "vorpalstacks/internal/store/aws/kinesis"
)

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
