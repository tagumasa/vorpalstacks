// Package kinesis provides AWS Kinesis data stream storage functionality for vorpalstacks.
package kinesis

import (
	"crypto/md5"
	"encoding/binary"
	"fmt"
	"strings"
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
	arnBuilder      *svcarn.ARNBuilder
	accountID       string
	region          string
	mu              sync.Mutex
	sequenceCounter int64
	shardIDCounter  int64
	storage         storage.TransactionalStorageWith2PC
}

// NewKinesisStore creates a new KinesisStore instance with the specified storage, account ID, and region.
func NewKinesisStore(store storage.TransactionalStorageWith2PC, accountID, region string) *KinesisStore {
	return &KinesisStore{
		BaseStore:      common.NewBaseStore(store.Bucket("kinesis-streams-"+region), "kinesis-streams"),
		shardsStore:    common.NewBaseStore(store.Bucket("kinesis-shards-"+region), "kinesis-shards"),
		recordsStore:   common.NewBaseStore(store.Bucket("kinesis-records-"+region), "kinesis-records"),
		consumersStore: common.NewBaseStore(store.Bucket("kinesis-consumers-"+region), "kinesis-consumers"),
		iteratorsStore: common.NewBaseStore(store.Bucket("kinesis-iterators-"+region), "kinesis-iterators"),
		TagStore:       common.NewTagStoreWithRegion(store, "kinesis", region),
		arnBuilder:     svcarn.NewARNBuilder(accountID, region),
		accountID:      accountID,
		region:         region,
		storage:        store,
	}
}

// GetAccountID returns the account ID associated with this store.
func (s *KinesisStore) GetAccountID() string {
	return s.accountID
}

// GetRegion returns the region associated with this store.
func (s *KinesisStore) GetRegion() string {
	return s.region
}

func (s *KinesisStore) buildStreamARN(streamName string) string {
	return s.arnBuilder.Kinesis().Stream(streamName)
}

// BuildStreamARN constructs the ARN for a Kinesis stream.
func (s *KinesisStore) BuildStreamARN(streamName string) string {
	return s.buildStreamARN(streamName)
}

func (s *KinesisStore) buildConsumerARN(streamName, consumerName string) string {
	return s.arnBuilder.Kinesis().Consumer(streamName, consumerName)
}

func (s *KinesisStore) generateSequenceNumber(shardID string) string {
	ts := time.Now().UTC().UnixNano()
	counter := atomic.AddInt64(&s.sequenceCounter, 1)
	return fmt.Sprintf("%d%012d%012d", ts, counter, s.hashToInt(shardID))
}

func (s *KinesisStore) hashToInt(str string) int64 {
	h := md5.Sum([]byte(str))
	return int64(binary.BigEndian.Uint64(h[:8]))
}

// PutRecordRequest represents a request to put a record into a Kinesis stream.
type PutRecordRequest struct {
	Data         string `json:"data"`
	PartitionKey string `json:"partitionKey"`
}

// PutRecordResult represents the result of putting a record into a Kinesis stream.
type PutRecordResult struct {
	SequenceNumber string `json:"sequenceNumber,omitempty"`
	ShardID        string `json:"shardId,omitempty"`
	ErrorCode      string `json:"errorCode,omitempty"`
	ErrorMessage   string `json:"errorMessage,omitempty"`
}

func normalizeARN(arn string, accountID string) string {
	parts := strings.SplitN(arn, ":", 6)
	if len(parts) >= 5 && parts[4] == "" {
		parts[4] = accountID
		return strings.Join(parts, ":")
	}
	return arn
}
