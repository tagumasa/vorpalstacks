package dynamodb

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/store/aws/common"
)

// errStopScan is a sentinel to stop prefix scans early.
var errStopScan = errors.New("stop scan")

// ShardIDForStream generates a deterministic shard identifier for a given
// stream ARN. Real AWS DynamoDB Streams uses multiple shards derived from
// partitions, but our single-node deployment uses a single shard per stream.
// The shard ID is derived from the stream ARN via SHA-256, ensuring each
// stream has a unique, deterministic identifier that persists across restarts.
func ShardIDForStream(streamArn string) string {
	h := sha256.Sum256([]byte(streamArn))
	seq := binary.BigEndian.Uint64(h[:8]) % 1000000000000000000
	return fmt.Sprintf("shardId-%020d-%012x", seq, h[8:14])
}

// StreamEventName enumerates the DynamoDB Streams event types.
type StreamEventName string

const (
	StreamEventInsert StreamEventName = "INSERT"
	StreamEventModify StreamEventName = "MODIFY"
	StreamEventRemove StreamEventName = "REMOVE"
)

// StreamRecord represents a single DynamoDB Streams record. The JSON
// structure matches the AWS DynamoDB Streams event format so that Lambda
// event source mappings can consume it directly.
type StreamRecord struct {
	EventID        string              `json:"eventID"`
	EventName      StreamEventName     `json:"eventName"`
	EventVersion   string              `json:"eventVersion"`
	EventSource    string              `json:"eventSource"`
	AWSRegion      string              `json:"awsRegion"`
	Dynamodb       StreamRecordData    `json:"dynamodb"`
	EventSourceARN string              `json:"eventSourceARN"`
	UserIdentity   *StreamUserIdentity `json:"userIdentity,omitempty"`
}

// StreamUserIdentity identifies the actor that triggered the stream event.
// For TTL deletions AWS sets this to {type: "Service", principalId:
// "dynamodb.amazonaws.com"} so consumers can distinguish TTL expiry from
// user-initiated deletes.
type StreamUserIdentity struct {
	Type        string `json:"type"`
	PrincipalID string `json:"principalId"`
}

// TTLServiceIdentity is the userIdentity AWS attaches to TTL-deletion records.
var TTLServiceIdentity = &StreamUserIdentity{
	Type:        "Service",
	PrincipalID: "dynamodb.amazonaws.com",
}

// StreamRecordData contains the DynamoDB-specific portion of a stream record.
type StreamRecordData struct {
	ApproximateCreationDateTime float64                `json:"ApproximateCreationDateTime,omitempty"`
	Keys                        map[string]interface{} `json:"Keys"`
	NewImage                    map[string]interface{} `json:"NewImage,omitempty"`
	OldImage                    map[string]interface{} `json:"OldImage,omitempty"`
	SequenceNumber              string                 `json:"SequenceNumber"`
	SizeBytes                   float64                `json:"SizeBytes"`
	StreamViewType              string                 `json:"StreamViewType"`
}

// streamSeqCounter is stored per table to atomically allocate sequence
// numbers. TrimmedFloor records the highest sequence number removed by the
// retention sweep; reads starting at or below it no longer have data.
type streamSeqCounter struct {
	LastSeq      int64 `json:"last_seq"`
	TrimmedFloor int64 `json:"trimmed_floor,omitempty"`
}

// StreamStore manages DynamoDB Streams records. Records are stored in a
// Pebble bucket keyed by "tableName\x00sequenceNumber" where sequenceNumber
// is a zero-padded 20-digit integer for correct lexicographic ordering.
type StreamStore struct {
	*common.BaseStore
	mu        sync.Mutex
	accountID string
	region    string
}

func streamBucketName(region string) string {
	return "dynamodb_streams-" + region
}

// NewStreamStore creates a new stream record store for the given region.
func NewStreamStore(store storage.BasicStorage, accountID, region string) *StreamStore {
	bucketName := streamBucketName(region)
	return &StreamStore{
		BaseStore: common.NewBaseStore(store.Bucket(bucketName), "dynamodb-streams"),
		accountID: accountID,
		region:    region,
	}
}

func streamRecordKey(tableName string, seq int64) string {
	return tableName + keySep + fmt.Sprintf("%020d", seq)
}

func streamSeqKey(tableName string) string {
	return tableName + keySep + "__seq_counter__"
}

// AddRecord stores a new stream record for the given table. It atomically
// allocates the next sequence number and writes the record.
func (s *StreamStore) AddRecord(tableName, streamArn, streamViewType string, eventName StreamEventName, keys, newImage, oldImage map[string]interface{}, userIdentity *StreamUserIdentity) (*StreamRecord, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Read and increment the sequence counter.
	var counter streamSeqCounter
	counterKey := streamSeqKey(tableName)
	if err := s.BaseStore.Get(counterKey, &counter); err != nil {
		if !common.IsNotFound(err) {
			return nil, fmt.Errorf("failed to read stream sequence counter: %w", err)
		}
		counter = streamSeqCounter{}
	}
	counter.LastSeq++
	if err := s.BaseStore.Put(counterKey, counter); err != nil {
		return nil, fmt.Errorf("failed to write stream sequence counter: %w", err)
	}

	record := s.buildRecord(tableName, streamArn, streamViewType, eventName, keys, newImage, oldImage, counter.LastSeq, userIdentity)

	recordKey := streamRecordKey(tableName, counter.LastSeq)
	if err := s.BaseStore.Put(recordKey, record); err != nil {
		return nil, fmt.Errorf("failed to write stream record: %w", err)
	}

	return record, nil
}

// AddRecordTxn writes a stream record within the given storage transaction,
// ensuring atomicity with the item mutation that triggered the stream event.
// The caller must commit the transaction for the record to be persisted.
func (s *StreamStore) AddRecordTxn(txn storage.Transaction, tableName, streamArn, streamViewType string, eventName StreamEventName, keys, newImage, oldImage map[string]interface{}, userIdentity *StreamUserIdentity) (*StreamRecord, error) {
	bucket := txn.Bucket(streamBucketName(s.region))

	// Read and increment the sequence counter.
	counterKey := []byte(streamSeqKey(tableName))
	counterData, err := bucket.Get(counterKey)
	if err != nil {
		return nil, fmt.Errorf("failed to read stream sequence counter: %w", err)
	}
	var counter streamSeqCounter
	if len(counterData) > 0 {
		if err := json.Unmarshal(counterData, &counter); err != nil {
			return nil, fmt.Errorf("failed to unmarshal stream sequence counter: %w", err)
		}
	}
	counter.LastSeq++
	counterBytes, err := json.Marshal(counter)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal stream sequence counter: %w", err)
	}
	if err := bucket.Put(counterKey, counterBytes); err != nil {
		return nil, fmt.Errorf("failed to write stream sequence counter: %w", err)
	}

	record := s.buildRecord(tableName, streamArn, streamViewType, eventName, keys, newImage, oldImage, counter.LastSeq, userIdentity)

	recordKey := []byte(streamRecordKey(tableName, counter.LastSeq))
	recordBytes, err := json.Marshal(record)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal stream record: %w", err)
	}
	if err := bucket.Put(recordKey, recordBytes); err != nil {
		return nil, fmt.Errorf("failed to write stream record: %w", err)
	}

	return record, nil
}

// buildRecord constructs a StreamRecord from the given parameters.
func (s *StreamStore) buildRecord(tableName, streamArn, streamViewType string, eventName StreamEventName, keys, newImage, oldImage map[string]interface{}, seq int64, userIdentity *StreamUserIdentity) *StreamRecord {
	sizeBytes := int64(0)
	if newImage != nil {
		if data, err := json.Marshal(newImage); err == nil {
			sizeBytes += int64(len(data))
		}
	}
	if oldImage != nil {
		if data, err := json.Marshal(oldImage); err == nil {
			sizeBytes += int64(len(data))
		}
	}
	if keys != nil {
		if data, err := json.Marshal(keys); err == nil {
			sizeBytes += int64(len(data))
		}
	}

	seqStr := strconv.FormatInt(seq, 10)
	return &StreamRecord{
		EventID:        fmt.Sprintf("%s-%s", tableName, seqStr),
		EventName:      eventName,
		EventVersion:   "1.1",
		EventSource:    "aws:dynamodb",
		AWSRegion:      s.region,
		EventSourceARN: streamArn,
		Dynamodb: StreamRecordData{
			ApproximateCreationDateTime: float64(time.Now().Unix()),
			Keys:                        keys,
			NewImage:                    newImage,
			OldImage:                    oldImage,
			SequenceNumber:              seqStr,
			SizeBytes:                   float64(sizeBytes),
			StreamViewType:              streamViewType,
		},
		UserIdentity: userIdentity,
	}
}

// GetRecords retrieves stream records for the given table starting from
// the given sequence number (exclusive). Returns up to limit records and
// the next sequence number for pagination. If no more records are
// available, returns an empty slice and the current latest sequence number.
func (s *StreamStore) GetRecords(tableName string, fromSeq int64, limit int) ([]*StreamRecord, int64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Determine the latest sequence number.
	var counter streamSeqCounter
	counterKey := streamSeqKey(tableName)
	if err := s.BaseStore.Get(counterKey, &counter); err != nil {
		if !common.IsNotFound(err) {
			return nil, 0, fmt.Errorf("failed to read stream sequence counter: %w", err)
		}
		return nil, 0, nil // No records yet.
	}

	if limit <= 0 {
		limit = 100
	}

	var records []*StreamRecord
	startKey := streamRecordKey(tableName, fromSeq+1)
	prefix := tableName + keySep

	err := s.BaseStore.ScanPrefix(prefix, func(key string, value []byte) error {
		// Skip the sequence counter key (it is not a stream record).
		if strings.HasSuffix(key, "__seq_counter__") {
			return nil
		}
		if key < startKey {
			// Skip records before the requested start position.
			return nil
		}
		var rec StreamRecord
		if err := json.Unmarshal(value, &rec); err != nil {
			return err
		}
		records = append(records, &rec)
		if len(records) >= limit {
			return errStopScan
		}
		return nil
	})
	if err != nil && err != errStopScan {
		return nil, 0, err
	}

	nextSeq := fromSeq
	if len(records) > 0 {
		lastSeq, err := strconv.ParseInt(records[len(records)-1].Dynamodb.SequenceNumber, 10, 64)
		if err != nil {
			return nil, 0, fmt.Errorf("invalid sequence number in stream record: %w", err)
		}
		nextSeq = lastSeq
	}

	return records, nextSeq, nil
}

// GetLatestSequence returns the latest sequence number for the given table.
func (s *StreamStore) GetLatestSequence(tableName string) (int64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var counter streamSeqCounter
	counterKey := streamSeqKey(tableName)
	if err := s.BaseStore.Get(counterKey, &counter); err != nil {
		if common.IsNotFound(err) {
			return 0, nil
		}
		return 0, err
	}
	return counter.LastSeq, nil
}

// StreamRetention is the documented DynamoDB Streams retention window:
// records older than 24 hours are subject to trimming.
const StreamRetention = 24 * time.Hour

// TrimOlderThan removes the stream records whose approximate creation time
// is before the cutoff and advances the table's trimmed floor to the
// highest removed sequence number. Keys are collected before any delete so
// the prefix scan never mutates while iterating.
func (s *StreamStore) TrimOlderThan(tableName string, cutoff time.Time) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	var counter streamSeqCounter
	counterKey := streamSeqKey(tableName)
	if err := s.BaseStore.Get(counterKey, &counter); err != nil {
		if common.IsNotFound(err) {
			return nil
		}
		return fmt.Errorf("failed to read stream sequence counter: %w", err)
	}

	prefix := tableName + keySep
	var doomed []string
	var maxTrimmed int64
	err := s.BaseStore.ScanPrefix(prefix, func(key string, value []byte) error {
		if strings.HasSuffix(key, "__seq_counter__") {
			return nil
		}
		var rec StreamRecord
		if err := json.Unmarshal(value, &rec); err != nil {
			return err
		}
		if int64(rec.Dynamodb.ApproximateCreationDateTime) >= cutoff.Unix() {
			return nil
		}
		seq, err := strconv.ParseInt(rec.Dynamodb.SequenceNumber, 10, 64)
		if err != nil {
			return fmt.Errorf("invalid sequence number in stream record: %w", err)
		}
		doomed = append(doomed, key)
		if seq > maxTrimmed {
			maxTrimmed = seq
		}
		return nil
	})
	if err != nil {
		return err
	}
	if len(doomed) == 0 {
		return nil
	}
	for _, key := range doomed {
		if err := s.BaseStore.Delete(key); err != nil {
			return fmt.Errorf("failed to trim stream record: %w", err)
		}
	}
	if maxTrimmed > counter.TrimmedFloor {
		counter.TrimmedFloor = maxTrimmed
		if err := s.BaseStore.Put(counterKey, counter); err != nil {
			return fmt.Errorf("failed to write stream sequence counter: %w", err)
		}
	}
	return nil
}

// OldestSequence returns the trimmed floor: sequence numbers at or below
// it have already been removed by retention trimming.
func (s *StreamStore) OldestSequence(tableName string) (int64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var counter streamSeqCounter
	counterKey := streamSeqKey(tableName)
	if err := s.BaseStore.Get(counterKey, &counter); err != nil {
		if common.IsNotFound(err) {
			return 0, nil
		}
		return 0, fmt.Errorf("failed to read stream sequence counter: %w", err)
	}
	return counter.TrimmedFloor, nil
}

// DeleteAllForTable removes all stream records and the sequence counter
// for the given table. Called during table deletion.
func (s *StreamStore) DeleteAllForTable(tableName string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	prefix := tableName + keySep
	var keys []string
	err := s.BaseStore.ScanPrefix(prefix, func(key string, value []byte) error {
		keys = append(keys, key)
		return nil
	})
	if err != nil {
		return err
	}
	for _, key := range keys {
		if err := s.BaseStore.Delete(key); err != nil {
			return err
		}
	}
	return nil
}
