// Package kinesis provides AWS Kinesis data stream storage functionality for vorpalstacks.
package kinesis

import (
	"crypto/md5"
	"encoding/binary"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/store/aws/common"
	svcarn "vorpalstacks/internal/utils/aws/arn"
)

// KinesisStore provides data storage operations for AWS Kinesis streams.
type KinesisStore struct {
	*common.BaseStore
	shardsStore    *common.BaseStore
	recordsStore   *common.BaseStore
	consumersStore *common.BaseStore
	iteratorsStore *common.BaseStore
	*common.TagStore
	arnIndexBucket      storage.Bucket
	arnBuilder          *svcarn.ARNBuilder
	accountID           string
	region              string
	mu                  sync.Mutex
	sequenceCounter     int64
	storage             storage.TransactionalStorageWith2PC
	nextIteratorCleanup time.Time
	nextRetentionTrim   time.Time
}

// NewKinesisStore creates a new KinesisStore instance with the specified storage, account ID, and region.
func NewKinesisStore(store storage.TransactionalStorageWith2PC, accountID, region string) *KinesisStore {
	ks := &KinesisStore{
		accountID:      accountID,
		region:         region,
		storage:        store,
		arnBuilder:     svcarn.NewARNBuilder(accountID, region),
		arnIndexBucket: store.Bucket("kinesis-streams-arn-" + region),
	}
	ks.BaseStore = common.NewBaseStore(store.Bucket(ks.streamsBucketName()), "kinesis-streams")
	ks.shardsStore = common.NewBaseStore(store.Bucket(ks.shardsBucketName()), "kinesis-shards")
	ks.recordsStore = common.NewBaseStore(store.Bucket(ks.recordsBucketName()), "kinesis-records")
	ks.consumersStore = common.NewBaseStore(store.Bucket(ks.consumersBucketName()), "kinesis-consumers")
	ks.iteratorsStore = common.NewBaseStore(store.Bucket(ks.iteratorsBucketName()), "kinesis-iterators")
	ks.TagStore = common.NewTagStoreWithRegion(store, "kinesis", region, common.TagBudget{
		MaxKeys:  MaxTagsPerResource,
		Exceeded: common.NewAWSError("InvalidArgumentException", "Too many tags.", http.StatusBadRequest),
	})
	return ks
}

// Bucket names — the single definition sites. The transactional paths
// resolve buckets by name inside their transactions, so the names travel
// as methods rather than pre-opened handles.
func (s *KinesisStore) streamsBucketName() string { return "kinesis-streams-" + s.region }
func (s *KinesisStore) shardsBucketName() string  { return "kinesis-shards-" + s.region }
func (s *KinesisStore) recordsBucketName() string { return "kinesis-records-" + s.region }
func (s *KinesisStore) consumersBucketName() string {
	return "kinesis-consumers-" + s.region
}
func (s *KinesisStore) iteratorsBucketName() string {
	return "kinesis-iterators-" + s.region
}
func (s *KinesisStore) policiesBucketName() string { return "kinesis-policies-" + s.region }
func (s *KinesisStore) arnIndexBucketName() string { return "kinesis-streams-arn-" + s.region }

func (s *KinesisStore) buildStreamARN(streamName string) string {
	return s.arnBuilder.Kinesis().Stream(streamName)
}

// BuildStreamARN constructs the ARN for a Kinesis stream.
func (s *KinesisStore) BuildStreamARN(streamName string) string {
	return s.buildStreamARN(streamName)
}

func (s *KinesisStore) buildConsumerARN(streamName, consumerName string, createdEpochSeconds int64) string {
	return s.arnBuilder.Kinesis().Consumer(streamName, consumerName, createdEpochSeconds)
}

// formatSequenceNumber is the single definition of the sequence-number
// format: arrival-time nanoseconds, a monotonic counter, and a shard hash,
// all decimal with the counter and hash zero-padded. The retention trim
// relies on the timestamp leading the key.
func formatSequenceNumber(ts, counter, shardHash int64) string {
	return fmt.Sprintf("%d%012d%012d", ts, counter, shardHash)
}

// sequenceNumberTime extracts the arrival-time component of a sequence
// number: the leading decimal digits ahead of the fixed-width counter and
// shard hash (formatSequenceNumber writes both to exactly twelve digits —
// the hash value is masked into that width at generation). A number that
// does not carry the fixed tail decodes as not-ok and the caller falls back
// to its own timestamp.
func sequenceNumberTime(seq string) (time.Time, bool) {
	if len(seq) <= sequenceNumberCounterHashDigits {
		return time.Time{}, false
	}
	nanos, err := strconv.ParseInt(seq[:len(seq)-sequenceNumberCounterHashDigits], 10, 64)
	if err != nil {
		return time.Time{}, false
	}
	return time.Unix(0, nanos).UTC(), true
}

// sequenceNumberCounterHashDigits is the combined width of the counter and
// shard-hash fields every sequence number carries after its timestamp.
const sequenceNumberCounterHashDigits = 24

func (s *KinesisStore) generateSequenceNumber(shardID string) string {
	return s.generateSequenceNumberAt(shardID, time.Now().UTC())
}

// generateSequenceNumberAt mints one sequence number from the caller's
// stamp: the record-write paths pass the same time they record as the
// arrival timestamp, so the number's embedded time and the record's
// arrival can never straddle a clock boundary.
func (s *KinesisStore) generateSequenceNumberAt(shardID string, now time.Time) string {
	return formatSequenceNumber(now.UnixNano(), atomic.AddInt64(&s.sequenceCounter, 1), s.hashToInt(shardID))
}

// hashToInt derives a shard's tie-breaking hash. The sequence-number
// format reserves exactly twelve decimal digits for the field, and
// zero-padding carries a minimum width only, so the value is masked into
// the reserved width — twelve digits hold at most 2^39-1.
func (s *KinesisStore) hashToInt(str string) int64 {
	h := md5.Sum([]byte(str))
	return int64(binary.BigEndian.Uint64(h[:8]) >> 25)
}

// PutRecordRequest represents a request to put a record into a Kinesis stream.
type PutRecordRequest struct {
	Data            string `json:"data"`
	PartitionKey    string `json:"partitionKey"`
	ExplicitHashKey string `json:"explicitHashKey,omitempty"`
}

// PutRecordResult represents the result of putting a record into a Kinesis stream.
type PutRecordResult struct {
	SequenceNumber string `json:"sequenceNumber,omitempty"`
	ShardID        string `json:"shardId,omitempty"`
	ErrorCode      string `json:"errorCode,omitempty"`
	ErrorMessage   string `json:"errorMessage,omitempty"`
}

// normalizeARN restores an omitted account segment using this store's
// account, so the ARN index and equality comparisons run on the completed
// form. Parsing goes through the shared ARN utilities; an unparseable
// input is returned untouched and surfaces as not-found downstream.
func (s *KinesisStore) normalizeARN(arn string) string {
	partition, service, region, accountID, resource := svcarn.SplitARN(arn)
	if service == "" {
		return arn
	}
	if accountID == "" {
		accountID = s.accountID
	}
	return "arn:" + partition + ":" + service + ":" + region + ":" + accountID + ":" + resource
}
