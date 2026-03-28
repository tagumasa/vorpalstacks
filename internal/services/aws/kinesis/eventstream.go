package kinesis

import (
	"encoding/json"
	"io"

	kinesisstore "vorpalstacks/internal/store/aws/kinesis"
	"vorpalstacks/internal/utils/aws/eventstream"
)

// SubscribeToShardEventStreamWriter writes events to a SubscribeToShard event stream.
type SubscribeToShardEventStreamWriter struct {
	encoder *eventstream.Encoder
}

// NewSubscribeToShardEventStreamWriter creates a new writer for SubscribeToShard events.
func NewSubscribeToShardEventStreamWriter(w io.Writer) *SubscribeToShardEventStreamWriter {
	return &SubscribeToShardEventStreamWriter{
		encoder: eventstream.NewEncoder(w),
	}
}

// WriteSubscribeToShardEvent writes a SubscribeToShardEvent to the stream.
func (s *SubscribeToShardEventStreamWriter) WriteSubscribeToShardEvent(records []*kinesisstore.Record, continuationSeqNum string, millisBehindLatest int64, childShards []interface{}) error {
	formattedRecords := make([]map[string]interface{}, 0)
	for _, r := range records {
		formattedRecords = append(formattedRecords, map[string]interface{}{
			"SequenceNumber":              r.SequenceNumber,
			"ApproximateArrivalTimestamp": r.ApproximateArrivalTimestamp.Unix(),
			"Data":                        r.Data,
			"PartitionKey":                r.PartitionKey,
			"EncryptionType":              "NONE",
		})
	}

	eventPayload := map[string]interface{}{
		"Records":                    formattedRecords,
		"ContinuationSequenceNumber": continuationSeqNum,
		"MillisBehindLatest":         millisBehindLatest,
		"ChildShards":                childShards,
	}

	payload, err := json.Marshal(eventPayload)
	if err != nil {
		return err
	}

	return s.encoder.WriteEvent("SubscribeToShardEvent", "application/json", payload)
}

// WriteResourceNotFoundException writes a ResourceNotFoundException error to the stream.
func (s *SubscribeToShardEventStreamWriter) WriteResourceNotFoundException(message string) error {
	return s.encoder.WriteErrorEvent("ResourceNotFoundException", message)
}

// WriteInvalidArgumentException writes an InvalidArgumentException error to the stream.
func (s *SubscribeToShardEventStreamWriter) WriteInvalidArgumentException(message string) error {
	return s.encoder.WriteErrorEvent("InvalidArgumentException", message)
}

// WriteEndEvent writes an end event to the stream to signal completion.
func (s *SubscribeToShardEventStreamWriter) WriteEndEvent() error {
	return s.encoder.WriteEndEvent()
}

// WriteInitialResponse writes an initial-response event to the stream.
func (s *SubscribeToShardEventStreamWriter) WriteInitialResponse() error {
	return s.encoder.WriteInitialResponse([]byte("{}"))
}
