package dynamodb

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"google.golang.org/protobuf/proto"
	"vorpalstacks/internal/core/storage"
	pb "vorpalstacks/internal/pb/storage/storage_dynamodb"
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

// ReplicationServiceIdentity is the userIdentity attached to records written
// by global-table replication: consumers distinguish replicated writes from
// direct client writes by the Service principal, the same principal AWS
// attaches to its own service-initiated records.
var ReplicationServiceIdentity = &StreamUserIdentity{
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
	LastSeq      int64
	TrimmedFloor int64
}

// readSeqCounter loads a table's persisted sequence allocator state; exists
// reports whether a counter record is present at all.
func (s *StreamStore) readSeqCounter(counterKey string) (counter streamSeqCounter, exists bool, err error) {
	var pbCounter pb.StreamSequenceCounter
	if err := s.BaseStore.GetProto(counterKey, &pbCounter); err != nil {
		if common.IsNotFound(err) {
			return streamSeqCounter{}, false, nil
		}
		return streamSeqCounter{}, false, fmt.Errorf("failed to read stream sequence counter: %w", err)
	}
	return protoToStreamCounter(&pbCounter), true, nil
}

// StreamStore manages DynamoDB Streams records. Records are stored in a
// Pebble bucket keyed by "tableName\x00sequenceNumber" where sequenceNumber
// is a zero-padded 20-digit integer for correct lexicographic ordering.
type StreamStore struct {
	*common.BaseStore
	mu        sync.Mutex
	accountID string
	region    string
	// seqNext holds the last sequence number allocated per table in this
	// process. Allocation must be unique across every writer, so it happens
	// outside the caller's transaction under mu: the storage layer is
	// read-committed, so allocating by counter read-modify-write inside a
	// transaction makes two records in one transaction (or two concurrent
	// transactions) collide on the same record key and silently overwrite
	// each other. The map seeds lazily from the persisted counter and the
	// highest record key (a committed counter can lag the records when
	// transactions commit out of allocation order). A table deleted and
	// recreated under the same name leaves a stale entry, which only
	// produces a gap in the new stream's numbering — sequence numbers stay
	// unique and increasing.
	seqNext map[string]int64
	// iteratorKey caches the shard-iterator signing key read (or generated)
	// by IteratorSigningKey.
	iteratorKey []byte
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
		seqNext:   make(map[string]int64),
	}
}

func streamRecordKey(tableName string, seq int64) string {
	return tableName + KeySep + fmt.Sprintf("%020d", seq)
}

func streamSeqKey(tableName string) string {
	return tableName + KeySep + "__seq_counter__"
}

// allocateSeq reserves the next sequence number for a table. The caller
// must hold s.mu. Allocation is a purely in-memory increment seeded once
// per table per process from the persisted counter and the highest record
// key — the counter alone is not enough, because transactions commit out
// of allocation order and can leave it behind the records. The returned
// counter is the value the caller persists alongside the record; its
// TrimmedFloor is read under the lock so retention trimming and allocation
// cannot clobber each other's floor within the process.
func (s *StreamStore) allocateSeq(tableName string) (int64, streamSeqCounter, error) {
	if _, seeded := s.seqNext[tableName]; !seeded {
		counter, _, err := s.readSeqCounter(streamSeqKey(tableName))
		if err != nil {
			return 0, streamSeqCounter{}, err
		}
		highest := counter.LastSeq
		prefix := tableName + KeySep
		err = s.BaseStore.ScanPrefix(prefix, func(key string, _ []byte) error {
			if strings.HasSuffix(key, "__seq_counter__") {
				return nil
			}
			digits := key[len(key)-20:]
			seq, err := strconv.ParseInt(digits, 10, 64)
			if err != nil {
				return fmt.Errorf("invalid sequence number in stream record key %q: %w", key, err)
			}
			if seq > highest {
				highest = seq
			}
			return nil
		})
		if err != nil {
			return 0, streamSeqCounter{}, err
		}
		s.seqNext[tableName] = highest
	}
	s.seqNext[tableName]++
	seq := s.seqNext[tableName]

	current, _, err := s.readSeqCounter(streamSeqKey(tableName))
	if err != nil {
		return 0, streamSeqCounter{}, err
	}
	return seq, streamSeqCounter{LastSeq: seq, TrimmedFloor: current.TrimmedFloor}, nil
}

// iteratorSigningKeyKey is the fixed record the stream bucket persists the
// shard-iterator signing key under. It carries no table prefix, so it never
// collides with per-table record or counter keys.
const iteratorSigningKeyKey = "__iterator_signing_key__"

// IteratorSigningKey returns the store-wide key shard iterators are signed
// with, generating and persisting a fresh random key on first use. Shard
// iterators are opaque server-issued tokens and the signature is what stops
// a client from crafting or tampering with one, so the key lives in the
// store and stays stable across restarts — a per-process key would silently
// invalidate every iterator issued before a restart.
func (s *StreamStore) IteratorSigningKey() ([]byte, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.iteratorKey != nil {
		return s.iteratorKey, nil
	}
	// The bucket contract reports a missing record as empty data with no
	// error, so absence is the empty case here, not a not-found error.
	key, err := s.BaseStore.GetRaw(iteratorSigningKeyKey)
	if err != nil {
		return nil, fmt.Errorf("failed to read iterator signing key: %w", err)
	}
	if len(key) > 0 {
		s.iteratorKey = key
		return key, nil
	}
	generated := make([]byte, 32)
	if _, err := rand.Read(generated); err != nil {
		return nil, fmt.Errorf("failed to generate iterator signing key: %w", err)
	}
	if err := s.BaseStore.PutRaw(iteratorSigningKeyKey, generated); err != nil {
		return nil, fmt.Errorf("failed to persist iterator signing key: %w", err)
	}
	s.iteratorKey = generated
	return generated, nil
}

// AddRecord stores a new stream record for the given table. It allocates
// the next sequence number outside any caller transaction and writes the
// record and the counter immediately.
func (s *StreamStore) AddRecord(tableName, streamArn, streamViewType string, eventName StreamEventName, keys, newImage, oldImage map[string]interface{}, userIdentity *StreamUserIdentity) (*StreamRecord, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	seq, counter, err := s.allocateSeq(tableName)
	if err != nil {
		return nil, err
	}
	if err := s.BaseStore.PutProto(streamSeqKey(tableName), streamCounterToProto(counter)); err != nil {
		return nil, fmt.Errorf("failed to write stream sequence counter: %w", err)
	}

	record := s.buildRecord(tableName, streamArn, streamViewType, eventName, keys, newImage, oldImage, seq, userIdentity)

	recordKey := streamRecordKey(tableName, seq)
	if err := s.BaseStore.PutProto(recordKey, streamRecordToProto(record)); err != nil {
		return nil, fmt.Errorf("failed to write stream record: %w", err)
	}

	return record, nil
}

// AddRecordTxn writes a stream record within the given storage transaction,
// ensuring atomicity with the item mutation that triggered the stream
// event. The caller must commit the transaction for the record to be
// persisted. The sequence number is allocated outside the transaction
// (see allocateSeq), so two records added in one transaction — or in two
// concurrent transactions — receive distinct consecutive numbers and never
// collide on the same record key. The counter is written inside the
// carrying transaction with the allocated number: it therefore never runs
// ahead of the records it numbers, and because allocation order and commit
// order can differ, restart seeding takes the highest record key into
// account. A retention trim that lands between allocation and commit can
// have its floor rewritten by this write; the next trim recomputes the
// floor from the records, so the value is advisory and self-healing.
func (s *StreamStore) AddRecordTxn(txn storage.Transaction, tableName, streamArn, streamViewType string, eventName StreamEventName, keys, newImage, oldImage map[string]interface{}, userIdentity *StreamUserIdentity) (*StreamRecord, error) {
	s.mu.Lock()
	seq, counter, err := s.allocateSeq(tableName)
	s.mu.Unlock()
	if err != nil {
		return nil, err
	}

	bucket := txn.Bucket(streamBucketName(s.region))

	counterBytes, err := proto.Marshal(streamCounterToProto(counter))
	if err != nil {
		return nil, fmt.Errorf("failed to marshal stream sequence counter: %w", err)
	}
	if err := bucket.Put([]byte(streamSeqKey(tableName)), counterBytes); err != nil {
		return nil, fmt.Errorf("failed to write stream sequence counter: %w", err)
	}

	record := s.buildRecord(tableName, streamArn, streamViewType, eventName, keys, newImage, oldImage, seq, userIdentity)

	recordKey := []byte(streamRecordKey(tableName, seq))
	recordBytes, err := proto.Marshal(streamRecordToProto(record))
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
	_, exists, err := s.readSeqCounter(streamSeqKey(tableName))
	if err != nil {
		return nil, 0, err
	}
	if !exists {
		return nil, 0, nil // No records yet.
	}

	if limit <= 0 {
		limit = 100
	}

	var records []*StreamRecord
	startKey := streamRecordKey(tableName, fromSeq+1)
	prefix := tableName + KeySep

	scanErr := s.BaseStore.ScanPrefix(prefix, func(key string, value []byte) error {
		// Skip the sequence counter key (it is not a stream record).
		if strings.HasSuffix(key, "__seq_counter__") {
			return nil
		}
		if key < startKey {
			// Skip records before the requested start position.
			return nil
		}
		var pbRec pb.StoredStreamRecord
		if err := proto.Unmarshal(value, &pbRec); err != nil {
			return err
		}
		records = append(records, streamRecordFromProto(&pbRec))
		if len(records) >= limit {
			return errStopScan
		}
		return nil
	})
	if scanErr != nil && scanErr != errStopScan {
		return nil, 0, scanErr
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

	counter, _, err := s.readSeqCounter(streamSeqKey(tableName))
	if err != nil {
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

	counter, exists, err := s.readSeqCounter(streamSeqKey(tableName))
	if err != nil {
		return err
	}
	if !exists {
		return nil
	}
	counterKey := streamSeqKey(tableName)

	prefix := tableName + KeySep
	var doomed []string
	var maxTrimmed int64
	err = s.BaseStore.ScanPrefix(prefix, func(key string, value []byte) error {
		if strings.HasSuffix(key, "__seq_counter__") {
			return nil
		}
		var pbRec pb.StoredStreamRecord
		if err := proto.Unmarshal(value, &pbRec); err != nil {
			return err
		}
		rec := streamRecordFromProto(&pbRec)
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
		if err := s.BaseStore.PutProto(counterKey, streamCounterToProto(counter)); err != nil {
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

	counter, _, err := s.readSeqCounter(streamSeqKey(tableName))
	if err != nil {
		return 0, err
	}
	return counter.TrimmedFloor, nil
}
