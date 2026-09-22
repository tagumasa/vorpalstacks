package cloudwatchlogs

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"

	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/core/storage/chunk"
	"vorpalstacks/internal/store/aws/common"
	svcarn "vorpalstacks/internal/utils/aws/arn"
	"vorpalstacks/internal/utils/naming"
)

// Store provides CloudWatch Logs storage operations.
type Store struct {
	*common.BaseStore
	tagStore   *common.TagStore
	ts         storage.TransactionalStorage
	arnBuilder *svcarn.ARNBuilder
	region     string
	bucketName string
	chunksDir  string
	// groupLocks shards the per-group plane. Writers hold their group's
	// write side for the whole mutation (ingestion's metadata commit,
	// stream and group teardown, retention purges, the per-group
	// configuration puts that must not resurrect a torn-down group); event
	// readers hold the read side for the whole gather, so a served page
	// can never interleave a mid-flight chunk delete — the reader sees the
	// chunk set as it was before the delete or after it, never a torn
	// listing whose files are already gone. The lock's domain is the group
	// and the shard is only its address: every operation holds exactly one
	// group's lock for its whole critical section, so different groups
	// never serialise and two groups that share a shard merely share a
	// queue.
	groupLocks   [groupLockShards]sync.RWMutex
	chunkCounter uint64
	subFilterMu  sync.Mutex
}

// groupLockShards is the width of the shard table guarding the per-group
// plane. A fixed table needs no teardown — a map of per-group locks would
// have to free a group's entry while operations may still be reaching for
// it through the very teardown that would do the freeing.
const groupLockShards = 64

// groupLock returns the shard guarding one log group's plane. The hash is
// FNV-1a over the group name; correctness never depends on it — the
// distribution only decides how often unrelated groups share a queue.
func (s *Store) groupLock(logGroupName string) *sync.RWMutex {
	var h uint32 = 2166136261
	for i := 0; i < len(logGroupName); i++ {
		h ^= uint32(logGroupName[i])
		h *= 16777619
	}
	return &s.groupLocks[h%groupLockShards]
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

	s := &Store{
		BaseStore: baseStore,
		// The cumulative tag ceiling answers with TooManyTagsException —
		// the identity TagResource's error list alone declares. The legacy
		// tag writers (TagLogGroup, PutDestination) declare no such error
		// and convert the budget identity to InvalidParameterException at
		// their Core.
		tagStore: common.NewTagStoreWithRegion(store, "cloudwatchlogs", region, common.TagBudget{
			MaxKeys: common.MaxTagsPerResource,
			Exceeded: common.NewAWSError("TooManyTagsException",
				fmt.Sprintf("A resource can have no more than %d tags", common.MaxTagsPerResource),
				http.StatusBadRequest),
		}),
		ts:         ts,
		arnBuilder: svcarn.NewARNBuilder(accountID, region),
		region:     region,
		bucketName: "logs-" + region,
		chunksDir:  chunksDir,
	}
	if err := s.seedChunkIdentity(); err != nil {
		return nil, fmt.Errorf("failed to seed chunk identity from stored state: %w", err)
	}
	return s, nil
}

// seedChunkIdentity derives the floors of both chunk identity counters from
// the durable state a restart inherits, so a replayed write can never
// reuse an identity an earlier process issued. The index-key sequence
// comes from the highest sequence in the chunk index (the keys carry the
// chunk ID), and the file-name sequence from the highest sequence among
// the chunk files themselves (the names carry it). Without this, both
// counters restart at zero and a backfill replayed at the same timestamps
// recomputes an earlier chunk's index key — overwriting its index record,
// which makes the earlier events vanish from reads — and its file name,
// which truncates the earlier file.
func (s *Store) seedChunkIdentity() error {
	var maxSeq uint64
	err := s.ScanPrefix(keyPrefixChunk, func(key string, _ []byte) error {
		// The chunk ID is the last colon-separated segment of the key;
		// the group and stream segments escape their own colons.
		id := key[strings.LastIndexByte(key, ':')+1:]
		parts := strings.Split(id, "-")
		if len(parts) != 3 {
			return nil
		}
		seq, err := strconv.ParseUint(parts[1], 10, 64)
		if err != nil {
			return nil
		}
		if seq > maxSeq {
			maxSeq = seq
		}
		return nil
	})
	if err != nil {
		return err
	}
	atomic.StoreUint64(&s.chunkCounter, maxSeq)

	entries, err := os.ReadDir(s.chunksDir)
	if err != nil {
		return err
	}
	var maxFileSeq uint64
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		parts := strings.Split(strings.TrimSuffix(entry.Name(), ".vlog"), "-")
		if len(parts) != 3 {
			continue
		}
		seq, err := strconv.ParseUint(parts[2], 10, 64)
		if err != nil {
			continue
		}
		if seq > maxFileSeq {
			maxFileSeq = seq
		}
	}
	chunk.SeedFileSequence(maxFileSeq)
	return nil
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

// Key prefixes are declared once: every key a family uses — single record,
// scan prefix or hand-composed range — derives from the same constant, so
// a layout change cannot leave a scan behind.
//
// Key layout discipline: a key is the family prefix followed by the
// record's identity components joined by ':'. Every component that
// precedes a separator is rendered through escapePath, so a name carrying
// ':', '/', '\' or '%' cannot blur a component boundary and two different
// component tuples can never produce one key; the final component is
// stored raw — no separator follows it, and a raw tail keeps a caller's
// name prefix matching the stored form byte-for-byte. No key uses any
// other separator.
const (
	keyPrefixLogGroup                = "log-group:"
	keyPrefixLogStream               = "log-stream:"
	keyPrefixMetricFilter            = "metric-filter:"
	keyPrefixChunk                   = "chunk:"
	keyPrefixSubscriptionFilter      = "subscription-filter:"
	keyPrefixDestination             = "destination:"
	keyPrefixResourcePolicy          = "resource-policy:"
	keyPrefixAccountPolicy           = "account-policy:"
	keyPrefixDataProtectionPolicy    = "data-protection-policy:"
	keyPrefixQueryDefinition         = "query-definition:"
	keyPrefixExportTask              = "export-task:"
	keyPrefixImportTask              = "import-task:"
	keyPrefixScheduledQuery          = "scheduled-query:"
	keyPrefixLookupTable             = "lookup-table:"
	keyPrefixScheduledQueryExecution = "sq-execution:"
	keyPrefixQueryRecord             = "query-record:"
	keyPrefixQueryResultKms          = "query-result-kms:"
	keyPrefixDeliverySource          = "delivery-source:"
	keyPrefixDeliveryDestination     = "delivery-destination:"
	keyPrefixDeliveryDestPolicy      = "delivery-destination-policy:"
	keyPrefixDelivery                = "delivery:"
	keyPrefixTransformer             = "transformer:"
	keyPrefixTransformedMessage      = "transformed-msg:"
	keyPrefixIndexPolicy             = "index-policy:"
	keyPrefixIndexInactiveTrail      = "index-inactive:"
	keyPrefixFieldIndexBounds        = "field-index-bounds:"
	keyPrefixStorageTierPolicy       = "storage-tier-policy:"
	keyPrefixPendingDelivery         = "pending-delivery:"
)

func (s *Store) logGroupKey(name string) string {
	return keyPrefixLogGroup + name
}

func (s *Store) logStreamKey(logGroupName, logStreamName string) string {
	return keyPrefixLogStream + escapePath(logGroupName) + ":" + logStreamName
}

func (s *Store) chunkIndexKey(logGroupName, logStreamName, chunkID string) string {
	return keyPrefixChunk + escapePath(logGroupName) + ":" + escapePath(logStreamName) + ":" + chunkID
}

func (s *Store) subscriptionFilterKey(logGroupName, filterName string) string {
	return keyPrefixSubscriptionFilter + escapePath(logGroupName) + ":" + filterName
}

func (s *Store) destinationKey(name string) string {
	return keyPrefixDestination + name
}

func (s *Store) resourcePolicyKey(policyName string) string {
	return keyPrefixResourcePolicy + policyName
}

func (s *Store) accountPolicyKey(policyType, policyName string) string {
	return keyPrefixAccountPolicy + escapePath(policyType) + ":" + policyName
}

func (s *Store) dataProtectionPolicyKey(logGroupIdentifier string) string {
	return keyPrefixDataProtectionPolicy + escapePath(logGroupIdentifier)
}

func (s *Store) queryDefinitionKey(id string) string {
	return keyPrefixQueryDefinition + id
}

func (s *Store) exportTaskKey(taskId string) string {
	return keyPrefixExportTask + taskId
}

func (s *Store) importTaskKey(importId string) string {
	return keyPrefixImportTask + importId
}

func (s *Store) scheduledQueryKey(id string) string {
	return keyPrefixScheduledQuery + id
}

func (s *Store) lookupTableKey(name string) string {
	return keyPrefixLookupTable + name
}

func (s *Store) scheduledQueryExecutionKey(sqId string, triggerTime int64) string {
	return keyPrefixScheduledQueryExecution + escapePath(sqId) + ":" + strconv.FormatInt(triggerTime, 10)
}

func (s *Store) queryRecordKey(queryId string) string {
	return keyPrefixQueryRecord + queryId
}

func (s *Store) deliverySourceKey(name string) string {
	return keyPrefixDeliverySource + name
}

func (s *Store) deliveryDestinationKey(name string) string {
	return keyPrefixDeliveryDestination + name
}

func (s *Store) deliveryDestinationPolicyKey(name string) string {
	return keyPrefixDeliveryDestPolicy + name
}

func (s *Store) deliveryKey(id string) string {
	return keyPrefixDelivery + id
}

func (s *Store) transformerKey(logGroupName string) string {
	return keyPrefixTransformer + escapePath(logGroupName)
}

func (s *Store) transformedMessageKey(logGroupName, digest string) string {
	return keyPrefixTransformedMessage + escapePath(logGroupName) + ":" + digest
}

func (s *Store) indexPolicyKey(logGroupName string) string {
	return keyPrefixIndexPolicy + escapePath(logGroupName)
}

func (s *Store) indexInactiveTrailKey(logGroupName string) string {
	return keyPrefixIndexInactiveTrail + escapePath(logGroupName)
}

func (s *Store) fieldIndexBoundsKey(logGroupName string) string {
	return keyPrefixFieldIndexBounds + escapePath(logGroupName)
}

// storageTierPolicyKey is the account-level singleton's key: one policy
// record per region store (the policy is account-scoped, the store
// regional).
func (s *Store) storageTierPolicyKey() string {
	return keyPrefixStorageTierPolicy + "account"
}

// escapePath renders a name component injectively into a key component:
// every byte that could begin an escape sequence is itself escaped, so a
// literal "%2f" inside a name cannot alias the escaped form of a name that
// contains '/' (the stream-name pattern permits '%').
func escapePath(path string) string {
	result := make([]byte, 0, len(path)*2)
	for i := 0; i < len(path); i++ {
		c := path[i]
		if c == '/' || c == ':' || c == '\\' || c == '%' {
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
