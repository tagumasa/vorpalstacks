package cloudwatchlogs

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"sync"

	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/store/aws/common"
	svcarn "vorpalstacks/internal/utils/aws/arn"
	"vorpalstacks/internal/utils/naming"
)

// Store provides CloudWatch Logs storage operations.
type Store struct {
	*common.BaseStore
	tagStore     *common.TagStore
	ts           storage.TransactionalStorage
	arnBuilder   *svcarn.ARNBuilder
	region       string
	bucketName   string
	chunksDir    string
	chunkMutex   sync.Mutex
	activeChunks map[string]*activeChunk
	chunkCounter uint64
	subFilterMu  sync.Mutex
}

type activeChunk struct {
	entries []LogEntry
}

// NewStore creates a new CloudWatch Logs store.
func NewStore(store storage.BasicStorage, bucket storage.Bucket, accountID, region, dataPath string) (*Store, error) {
	baseStore := common.NewBaseStore(bucket, "logs")
	chunksDir := filepath.Join(dataPath, "logs-chunks")
	if err := os.MkdirAll(chunksDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create logs chunks directory: %w", err)
	}

	var ts storage.TransactionalStorage
	if txnStorage, ok := store.(storage.TransactionalStorage); ok {
		ts = txnStorage
	}

	return &Store{
		BaseStore:    baseStore,
		tagStore:     common.NewTagStoreWithRegion(store, "cloudwatchlogs", region),
		ts:           ts,
		arnBuilder:   svcarn.NewARNBuilder(accountID, region),
		region:       region,
		bucketName:   "logs-" + region,
		chunksDir:    chunksDir,
		activeChunks: make(map[string]*activeChunk),
	}, nil
}

// ARNBuilder returns the ARN builder for the store.
func (s *Store) ARNBuilder() *svcarn.ARNBuilder {
	return s.arnBuilder
}

// Tags returns the tag store for CloudWatch Logs log groups.
func (s *Store) Tags() *common.TagStore {
	return s.tagStore
}

func (s *Store) safeChunkPath(chunkPath string) (string, error) {
	return naming.ValidatePathWithinDir(s.chunksDir, chunkPath)
}

// --- key layout helpers ---

func (s *Store) logGroupKey(name string) string {
	return "log-group:" + name
}

func (s *Store) logStreamKey(logGroupName, logStreamName string) string {
	return "log-stream:" + escapePath(logGroupName) + ":" + logStreamName
}

func (s *Store) chunkIndexKey(logGroupName, logStreamName, chunkID string) string {
	return "chunk:" + escapePath(logGroupName) + ":" + escapePath(logStreamName) + ":" + chunkID
}

func (s *Store) subscriptionFilterKey(logGroupName, filterName string) string {
	return "subscription-filter:" + escapePath(logGroupName) + ":" + filterName
}

func (s *Store) destinationKey(name string) string {
	return "destination:" + name
}

func (s *Store) resourcePolicyKey(policyName string) string {
	return "resource-policy:" + policyName
}

func (s *Store) accountPolicyKey(policyType, policyName string) string {
	return "account-policy:" + policyType + ":" + policyName
}

func (s *Store) dataProtectionPolicyKey(logGroupIdentifier string) string {
	return "data-protection-policy:" + escapePath(logGroupIdentifier)
}

func (s *Store) queryDefinitionKey(id string) string {
	return "query-definition:" + id
}

func (s *Store) exportTaskKey(taskId string) string {
	return "export-task:" + taskId
}

func (s *Store) importTaskKey(importId string) string {
	return "import-task:" + importId
}

func (s *Store) scheduledQueryKey(id string) string {
	return "scheduled-query:" + id
}

func (s *Store) lookupTableKey(name string) string {
	return "lookup-table:" + name
}

func (s *Store) scheduledQueryExecutionKey(sqId string, triggerTime int64) string {
	return "sq-execution:" + sqId + ":" + strconv.FormatInt(triggerTime, 10)
}

func escapePath(path string) string {
	result := make([]byte, 0, len(path)*2)
	for i := 0; i < len(path); i++ {
		c := path[i]
		if c == '/' || c == ':' || c == '\\' {
			result = append(result, '%')
			result = append(result, hexChar(c>>4))
			result = append(result, hexChar(c&0x0F))
		} else {
			result = append(result, c)
		}
	}
	return string(result)
}

func hexChar(b byte) byte {
	if b < 10 {
		return '0' + b
	}
	return 'a' + b - 10
}
