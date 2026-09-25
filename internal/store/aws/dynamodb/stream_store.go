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
// shardIdSeqModulus bounds the hash-derived shard-id sequence component so
// it stays below 10^18 — the rendered id keeps its fixed 20-digit field
// (%020d) with the leading two digits as padding, leaving room for the
// scheme to grow without renumbering existing shard ids.
const shardIdSeqModulus = 1_000_000_000_000_000_000

func ShardIDForStream(streamArn string) string {
	h := sha256.Sum256([]byte(streamArn))
	seq := binary.BigEndian.Uint64(h[:8]) % shardIdSeqModulus
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
// attaches to its own service-initiated records — one principal value with
// the TTL identity, two documented roles.
var ReplicationServiceIdentity = TTLServiceIdentity

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

// streamSeqCounter is the retention sweep's own per-table record.
// TrimmedFloor records the highest sequence number removed by the retention
// sweep; reads starting at or below it no longer have data. The counter
// deliberately persists no sequence extent: a value written inside each
// record's carrying transaction inherits the transaction's commit order, so
// a slower transaction would regress a persisted maximum below committed
// records — the stream's extent is derived from the record keys, which
// never regress. With no extent member, the record path has nothing to add
// to the counter and stops writing it, which leaves the sweep as the
// counter's single writer: the floor can only be raised, never rewritten by
// a racing record commit. Allocation still seeds from the floor — it is
// what survives a trim that removed every record, keeping removed numbers
// from being re-allocated after a restart.
type streamSeqCounter struct {
	TrimmedFloor int64
}

// readSeqCounter loads a table's persisted retention record; a table that
// has never trimmed carries none, which reads as a zero floor.
func (s *StreamStore) readSeqCounter(counterKey string) (counter streamSeqCounter, err error) {
	var pbCounter pb.StreamSequenceCounter
	if err := s.BaseStore.GetProto(counterKey, &pbCounter); err != nil {
		if common.IsNotFound(err) {
			return streamSeqCounter{}, nil
		}
		return streamSeqCounter{}, fmt.Errorf("failed to read stream sequence counter: %w", err)
	}
	return protoToStreamCounter(&pbCounter), nil
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
	// each other. The map seeds lazily from the highest record key and the
	// persisted trim floor (the two things that survive a restart: the
	// records themselves, and the floor when the records have been trimmed
	// away). A table deleted and recreated under the same name leaves a
	// stale entry, which only produces a gap in the new stream's numbering —
	// sequence numbers stay unique and increasing.
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
// per table per process from the highest record key and the persisted trim
// floor — together they bound every number any record has used or any trim
// has removed, so the numbering never restarts below a used number after a
// restart. Nothing is persisted here: the record's carrying transaction
// writes only the record (a counter value written inside it would inherit
// the transaction's commit order and regress below committed records), and
// the trim floor belongs to the retention sweep alone.
func (s *StreamStore) allocateSeq(tableName string) (int64, error) {
	if _, seeded := s.seqNext[tableName]; !seeded {
		recordHighest, err := s.highestRecordSequenceForStream(tableName, "")
		if err != nil {
			return 0, err
		}
		counter, err := s.readSeqCounter(streamSeqKey(tableName))
		if err != nil {
			return 0, err
		}
		highest := recordHighest
		if counter.TrimmedFloor > highest {
			highest = counter.TrimmedFloor
		}
		s.seqNext[tableName] = highest
	}
	s.seqNext[tableName]++
	return s.seqNext[tableName], nil
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

// AddRecordTxn writes a stream record within the given storage transaction,
// ensuring atomicity with the item mutation that triggered the stream
// event. The caller must commit the transaction for the record to be
// persisted. The sequence number is allocated outside the transaction
// (see allocateSeq), so two records added in one transaction — or in two
// concurrent transactions — receive distinct consecutive numbers and never
// collide on the same record key. The transaction writes only the record:
// the sequence counter is the retention sweep's own record (a counter value
// written here would inherit the transaction's commit order — regressing an
// extent below committed records, or rewriting a trim floor the sweep
// raised between this allocation and the commit), and the stream's extent
// is derived from the highest record key (highestRecordSequence), which
// never regresses.
func (s *StreamStore) AddRecordTxn(txn storage.Transaction, tableName, streamArn, streamViewType string, eventName StreamEventName, keys, newImage, oldImage map[string]interface{}, userIdentity *StreamUserIdentity) (*StreamRecord, error) {
	s.mu.Lock()
	seq, err := s.allocateSeq(tableName)
	s.mu.Unlock()
	if err != nil {
		return nil, err
	}

	bucket := txn.Bucket(streamBucketName(s.region))

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

// FormatStreamSequenceNumber renders a stream sequence number in the wire
// form the model defines: the SequenceNumber type is a 21-40 character
// numeric string, so the platform's compact integers are zero-padded to the
// 21-digit minimum. Parsing the padded form with ParseInt is unchanged, and
// records, shard ranges, and iterator positions all share the one format.
func FormatStreamSequenceNumber(seq int64) string {
	return fmt.Sprintf("%021d", seq)
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

	seqStr := FormatStreamSequenceNumber(seq)
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

// GetRecords retrieves the stream records of the named stream generation —
// the records whose eventSourceARN is streamArn — starting from the given
// sequence number (exclusive). Returns up to limit records and the next
// sequence number for pagination. If no more records are available, returns
// an empty slice and the current position. A generation change empties the
// record space, but a write that raced the change can still commit a
// superseded generation's record into it; the ARN comparison is what
// guarantees such a record is never served to the successor generation's
// readers.
//
// streamRecordsDefaultLimit bounds a reader that passes no explicit limit
// to one bounded batch of the shard, not the whole retained history.
const streamRecordsDefaultLimit = 100

func (s *StreamStore) GetRecords(tableName, streamArn string, fromSeq int64, limit int) ([]*StreamRecord, int64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if limit <= 0 {
		limit = streamRecordsDefaultLimit
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
		rec := streamRecordFromProto(&pbRec)
		if rec.EventSourceARN != streamArn {
			// A record of another generation of the same table — only
			// reachable when a write raced a generation change — never
			// serves the generation this read belongs to.
			return nil
		}
		records = append(records, rec)
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

// highestRecordSequenceForStream returns the highest committed record
// sequence number for the table, 0 when the stream holds no records. The
// caller holds s.mu. A non-empty streamArn restricts the count to the
// records of that stream generation; the empty string counts every record
// of the table, which is what allocateSeq's restart seeding needs — its
// numbering must stay unique and increasing across generations, residue
// included. Committed record keys are the truth about a stream's extent:
// the persisted counter carries no extent member (a transaction-committed
// value would regress under commit-order inversion), and the reverse scan
// reads only committed state, so the derived value is monotonic by
// construction.
func (s *StreamStore) highestRecordSequenceForStream(tableName, streamArn string) (int64, error) {
	var highest int64
	prefix := tableName + KeySep
	scanErr := s.BaseStore.ScanPrefixReverse(prefix, "", func(key string, value []byte) error {
		// The counter key sorts after every record key of the table, so
		// the reverse scan meets it first; it is not a record.
		if strings.HasSuffix(key, "__seq_counter__") {
			return nil
		}
		if len(key) < 20 {
			// A record key always ends in the 20-digit zero-padded
			// sequence, so anything shorter under the table prefix is
			// foreign or corrupt rather than a record — the scan
			// continues past it instead of slicing out of bounds.
			return nil
		}
		digits := key[len(key)-20:]
		seq, err := strconv.ParseInt(digits, 10, 64)
		if err != nil {
			return fmt.Errorf("invalid sequence number in stream record key %q: %w", key, err)
		}
		if streamArn != "" {
			var pbRec pb.StoredStreamRecord
			if err := proto.Unmarshal(value, &pbRec); err != nil {
				return fmt.Errorf("decode stream record %q during extent scan: %w", key, err)
			}
			if streamRecordFromProto(&pbRec).EventSourceARN != streamArn {
				// Residue of a superseded generation the best-effort
				// generation sweep missed: it shares the table's key space,
				// but no read of this generation may serve it, so it must
				// not advance this generation's extent either.
				return nil
			}
		}
		highest = seq
		return errStopScan
	})
	if scanErr != nil && scanErr != errStopScan {
		return 0, scanErr
	}
	return highest, nil
}

// GetLatestSequenceForStream returns the latest committed sequence number
// among the records of the named stream generation, derived from the
// highest matching record key (see highestRecordSequenceForStream for why
// the persisted counter is not the extent). The extent a stream reports —
// its EndingSequenceNumber and the LATEST iterator position — must count
// only the generation's own records.
func (s *StreamStore) GetLatestSequenceForStream(tableName, streamArn string) (int64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.highestRecordSequenceForStream(tableName, streamArn)
}

// DeleteTableRecords empties the table's whole stream record space — every
// generation's records and the sequence counter — the state reset a stream
// generation change owes: the successor generation starts with an empty
// space. The in-process sequence allocator keeps its high-water mark, which
// only leaves a gap in the successor's numbering (see allocateSeq); the
// persisted counter leaves with the records because it shares the table
// prefix. Keys are collected before any delete so the prefix scan never
// mutates while iterating.
func (s *StreamStore) DeleteTableRecords(tableName string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	prefix := tableName + KeySep
	var doomed []string
	err := s.BaseStore.ScanPrefix(prefix, func(key string, _ []byte) error {
		doomed = append(doomed, key)
		return nil
	})
	if err != nil {
		return err
	}
	for _, key := range doomed {
		if err := s.BaseStore.Delete(key); err != nil {
			return fmt.Errorf("failed to delete stream record: %w", err)
		}
	}
	return nil
}

// StreamRetention is the documented DynamoDB Streams retention window:
// records older than 24 hours are subject to trimming.
const StreamRetention = 24 * time.Hour

// TrimOlderThan removes the stream records whose approximate creation time
// is before the cutoff and advances the table's trimmed floor to the
// highest removed sequence number. Keys are collected before any delete so
// the prefix scan never mutates while iterating. The counter is this
// sweep's own record: it is created by the first trim that removes records
// (a table that has never trimmed carries none), and a missing counter is
// simply a zero floor — the records themselves decide what is doomed.
func (s *StreamStore) TrimOlderThan(tableName string, cutoff time.Time) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	counter, err := s.readSeqCounter(streamSeqKey(tableName))
	if err != nil {
		return err
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

	counter, err := s.readSeqCounter(streamSeqKey(tableName))
	if err != nil {
		return 0, err
	}
	return counter.TrimmedFloor, nil
}
