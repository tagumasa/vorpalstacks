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
// Every record reports the stream's encryption type — the same value the
// GetRecords path derives, so a KMS-encrypted stream never reports NONE
// over the subscribe transport.
func (s *SubscribeToShardEventStreamWriter) WriteSubscribeToShardEvent(records []*kinesisstore.Record, continuationSeqNum string, millisBehindLatest int64, childShards []interface{}, encryptionType string) error {
	formattedRecords := make([]map[string]interface{}, 0)
	for _, r := range records {
		formattedRecords = append(formattedRecords, map[string]interface{}{
			"SequenceNumber":              r.SequenceNumber,
			"ApproximateArrivalTimestamp": formatEpochSeconds(r.ApproximateArrivalTimestamp),
			"Data":                        r.Data,
			"PartitionKey":                r.PartitionKey,
			"EncryptionType":              encryptionType,
		})
	}

	eventPayload := map[string]interface{}{
		"Records":                    formattedRecords,
		"ContinuationSequenceNumber": continuationSeqNum,
		"MillisBehindLatest":         millisBehindLatest,
	}
	// ChildShards is present when the shard closes with children; idle
	// events carry no null member.
	if len(childShards) > 0 {
		eventPayload["ChildShards"] = childShards
	}

	payload, err := json.Marshal(eventPayload)
	if err != nil {
		return err
	}

	return s.encoder.WriteEvent("SubscribeToShardEvent", "application/json", payload)
}

// WriteErrorEvent writes a modelled-exception frame carrying a service
// error code — the mapped identity of the failure that ended the stream,
// typed-decodable at the client.
func (s *SubscribeToShardEventStreamWriter) WriteErrorEvent(code, message string) error {
	return s.encoder.WriteErrorEvent(code, message)
}

// WriteEndEvent writes an end event to the stream to signal completion.
func (s *SubscribeToShardEventStreamWriter) WriteEndEvent() error {
	return s.encoder.WriteEndEvent()
}

// WriteInitialResponse writes an initial-response event to the stream.
func (s *SubscribeToShardEventStreamWriter) WriteInitialResponse() error {
	return s.encoder.WriteInitialResponse([]byte("{}"))
}
