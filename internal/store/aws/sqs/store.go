package sqs

import (
	"context"
	"fmt"
	"sync"
	"time"

	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/store/aws/common"
	svcarn "vorpalstacks/internal/utils/aws/arn"
)

const (
	// MaxQueueNameLength is the maximum queue name length in characters
	// (80 per the AWS SQS specification).
	MaxQueueNameLength    = 80
	maxBatchEntryIdLength = 80
)

func messageKey(queueURL, messageID string) string {
	return queueURL + "\x00" + messageID
}

func messagePrefix(queueURL string) string {
	return queueURL + "\x00"
}

const (
	purgeTimeout = 60 * time.Second
	// MinVisibilityTimeout is the VisibilityTimeout lower bound in seconds.
	MinVisibilityTimeout = 0
	// MaxVisibilityTimeout is the SQS VisibilityTimeout maximum in seconds
	// (12 hours per the AWS SQS specification).
	MaxVisibilityTimeout = 43200
	// MinDelaySeconds is the DelaySeconds lower bound in seconds.
	MinDelaySeconds = 0
	// MaxDelaySeconds is the DelaySeconds upper bound in seconds
	// (15 minutes per the AWS SQS specification).
	MaxDelaySeconds = 900
	// MinMessageRetentionPeriod is the MessageRetentionPeriod lower bound in
	// seconds (1 minute per the AWS SQS specification).
	MinMessageRetentionPeriod = 60
	// MaxMessageRetentionPeriod is the MessageRetentionPeriod upper bound in
	// seconds (14 days per the AWS SQS specification).
	MaxMessageRetentionPeriod = 1209600
	// MinMaximumMessageSize is the MaximumMessageSize lower bound in bytes
	// (1 KiB per the AWS SQS specification).
	MinMaximumMessageSize = 1024
	// MaxMaximumMessageSize is the MaximumMessageSize upper bound in bytes
	// (1 MiB per the AWS SQS specification).
	MaxMaximumMessageSize = 1048576
	// MinReceiveMessageWaitTimeSeconds is the ReceiveMessageWaitTimeSeconds
	// lower bound in seconds.
	MinReceiveMessageWaitTimeSeconds = 0
	// MaxReceiveMessageWaitTimeSeconds is the ReceiveMessageWaitTimeSeconds
	// upper bound in seconds (20 per the AWS SQS specification).
	MaxReceiveMessageWaitTimeSeconds = 20
	// MinMaxNumberOfMessages is the ReceiveMessage MaxNumberOfMessages lower
	// bound (1 per the AWS SQS specification).
	MinMaxNumberOfMessages = 1
	// MaxMaxNumberOfMessages is the ReceiveMessage MaxNumberOfMessages
	// upper bound (10 per the AWS SQS specification).
	MaxMaxNumberOfMessages = 10
	// MaxBatchEntries is the maximum number of entries in SendMessageBatch,
	// DeleteMessageBatch and ChangeMessageVisibilityBatch requests.
	MaxBatchEntries = 10
	// MinKmsDataKeyReusePeriodSeconds is the KmsDataKeyReusePeriodSeconds
	// lower bound in seconds (1 minute per the AWS SQS specification).
	MinKmsDataKeyReusePeriodSeconds = 60
	// MaxKmsDataKeyReusePeriodSeconds is the KmsDataKeyReusePeriodSeconds
	// upper bound in seconds (24 hours per the AWS SQS specification).
	MaxKmsDataKeyReusePeriodSeconds = 86400
	// MaxListResults is the maximum and default MaxResults page size for
	// ListQueues and ListDeadLetterSourceQueues (AWS SQS API Reference:
	// "Value range is 1 to 1000").
	MaxListResults = 1000
	// MaxListMessageMoveTasksResults is the ListMessageMoveTasks MaxResults
	// upper bound (AWS SQS API Reference: "The maximum number of results to
	// include in the response. ... a maximum of 10 results").
	MaxListMessageMoveTasksResults = 10
	// DefaultListMessageMoveTasksResults is the ListMessageMoveTasks
	// MaxResults default (AWS SQS API Reference: "The default value is 1.").
	DefaultListMessageMoveTasksResults = 1
	// MaxRedriveAllowPolicySourceQueues is the byQueue sourceQueueArns cap
	// (AWS SQS Developer Guide, dead-letter queues: "you can specify up to
	// 10 source queues using the sourceQueueArn").
	MaxRedriveAllowPolicySourceQueues = 10
	// MaxMessageMoveRate is the StartMessageMoveTask
	// MaxNumberOfMessagesPerSecond ceiling (AWS SQS API Reference: "The
	// maximum task rate is 500 messages per second."). Exported as the
	// single owning definition; the services-layer validation references it.
	MaxMessageMoveRate = 500
	// DefaultMaxReceiveCount is the maxReceiveCount a RedrivePolicy uses when
	// the field is absent (AWS SQS API Reference: "Default: 10.").
	DefaultMaxReceiveCount = 10
	// DefaultVisibilityTimeout is the VisibilityTimeout a new queue starts
	// with (30 seconds per the AWS SQS specification).
	DefaultVisibilityTimeout = 30
	// DefaultMessageRetentionPeriod is the MessageRetentionPeriod a new
	// queue starts with (4 days per the AWS SQS specification).
	DefaultMessageRetentionPeriod = 345600
	// MinMaxReceiveCount is the RedrivePolicy maxReceiveCount lower bound
	// (AWS SQS Developer Guide, configuring a dead-letter queue: "Set the
	// Maximum receives value ... (valid range: 1 to 1,000)").
	MinMaxReceiveCount = 1
	// MaxMaxReceiveCount is the RedrivePolicy maxReceiveCount upper bound
	// (same page; AWS Knowledge Center: "You can increase the Maximum
	// receives value for the DLQ redrive policy up to 1,000").
	MaxMaxReceiveCount = 1000
	// MaxFifoIdLength is the maximum length in characters of MessageGroupId
	// and MessageDeduplicationId (AWS SQS API Reference: "The maximum length
	// of MessageDeduplicationId is 128 characters." / "The length of
	// MessageGroupId is 128 characters.").
	MaxFifoIdLength = 128
	// queueDeletionWindow is the minimum interval between deleting a queue
	// and creating another queue with the same name (AWS SQS API Reference:
	// "You must wait 60 seconds after deleting a queue before you can create
	// another queue with the same name.").
	queueDeletionWindow = 60 * time.Second
	deduplicationWindow = 5 * time.Minute
	// receiptHandleRetention bounds how long a receipt-handle entry stays
	// resolvable after it was issued. A handle can only be the current
	// in-flight handle within its visibility window, whose maximum
	// configurable length is MaxVisibilityTimeout; beyond that the handle is
	// a best-effort old handle ("If you use an old ReceiptHandle, the
	// request will succeed, but the message might not be deleted."), which
	// the platform is free to stop resolving.
	receiptHandleRetention = time.Duration(MaxVisibilityTimeout) * time.Second
)

// SQSStore provides SQS queue storage functionality.
type SQSStore struct {
	*common.BaseStore
	messagesStore *common.BaseStore
	tasksStore    *common.BaseStore
	*common.TagStore
	arnBuilder         *svcarn.ARNBuilder
	accountID          string
	region             string
	baseURL            string
	msgMutex           sync.RWMutex
	purgeMutex         sync.Mutex
	queueMutex         sync.RWMutex
	taskMu             sync.Mutex
	purgeInProgress    map[string]time.Time
	storage            storage.TransactionalStorage
	deduplicationCache map[string]*deduplicationEntry
	deduplicationMu    sync.RWMutex
	// dedupKeyLocker serialises the FIFO send's dedup check → message
	// persist → entry registration sequence per dedup key, so two concurrent
	// same-key sends cannot both miss and both persist.
	dedupKeyLocker      common.KeyLocker
	sequenceCounters    map[string]int64
	sequenceMu          sync.Mutex
	receiveAttemptCache map[string]*receiveAttemptEntry
	receiveAttemptMu    sync.Mutex
	ctx                 context.Context
	cancel              context.CancelFunc
	wg                  sync.WaitGroup
	// closing refuses new move-task workers during shutdown; guarded by
	// taskMu (see Close).
	closing bool
	// Auxiliary bucket names, derived once in the constructor through
	// auxiliaryBucketName: every access to these buckets goes through these
	// fields so the write paths and the background sweeps can never diverge
	// to differently spelled buckets.
	receiptsBucket     string
	dedupBucket        string
	deletionsBucket    string
	messagesBucketName string
}

// auxiliaryBucketName builds the name of a region-suffixed SQS bucket. Every
// auxiliary bucket name in the package — the constructor's BaseStore wraps
// and the ad-hoc bucket handles alike — is derived through this helper, so a
// misspelling cannot address a different bucket than the sweeps read.
func auxiliaryBucketName(base, region string) string {
	return base + "-" + region
}

type receiveAttemptEntry struct {
	messageIDs []string
	createdAt  time.Time
}

// NewSQSStore creates a new SQS store with the specified storage, account ID, region, and base URL.
func NewSQSStore(store storage.BasicStorage, accountID, region, baseURL string) *SQSStore {
	ts, _ := store.(storage.TransactionalStorage)
	ctx, cancel := context.WithCancel(context.Background())
	s := &SQSStore{
		BaseStore:           common.NewBaseStore(store.Bucket(auxiliaryBucketName("sqs-queues", region)), "sqs-queues"),
		messagesStore:       common.NewBaseStore(store.Bucket(auxiliaryBucketName("sqs-messages", region)), "sqs-messages"),
		tasksStore:          common.NewBaseStore(store.Bucket(auxiliaryBucketName("sqs-move-tasks", region)), "sqs-move-tasks"),
		TagStore:            common.NewTagStoreWithRegion(store, "sqs", region, common.TagBudget{MaxKeys: common.MaxTagsPerResource, Exceeded: ErrTooManyTags}),
		arnBuilder:          svcarn.NewARNBuilder(accountID, region),
		accountID:           accountID,
		region:              region,
		baseURL:             baseURL,
		storage:             ts,
		purgeInProgress:     make(map[string]time.Time),
		deduplicationCache:  make(map[string]*deduplicationEntry),
		sequenceCounters:    make(map[string]int64),
		receiveAttemptCache: make(map[string]*receiveAttemptEntry),
		ctx:                 ctx,
		cancel:              cancel,
		receiptsBucket:      auxiliaryBucketName("sqs-receipts", region),
		dedupBucket:         auxiliaryBucketName("sqs-dedup", region),
		deletionsBucket:     auxiliaryBucketName("sqs-queue-deletions", region),
		messagesBucketName:  auxiliaryBucketName("sqs-messages", region),
	}
	// An unclean shutdown leaves persisted move tasks in RUNNING/CANCELLING
	// with no worker alive to finish them; recovery finalises them before
	// any API call can observe the stale states.
	s.recoverInterruptedMoveTasks()
	s.wg.Add(1)
	go s.cleanupExpiredMessages()
	return s
}

// Close stops the background cleanup goroutine. The closing flag is set
// under taskMu so a StartMessageMoveTask racing the shutdown either
// completes its wg.Add before this critical section (Wait counts the
// worker) or is refused — never a positive-delta Add against a Wait that
// already observes zero.
func (s *SQSStore) Close() {
	if s.cancel != nil {
		s.taskMu.Lock()
		s.closing = true
		s.taskMu.Unlock()
		s.cancel()
		s.wg.Wait()
	}
}

// GetAccountID returns the AWS account ID associated with this SQS store.
func (s *SQSStore) GetAccountID() string {
	return s.accountID
}

// GetRegion returns the AWS region associated with this SQS store.
func (s *SQSStore) GetRegion() string {
	return s.region
}

func (s *SQSStore) buildQueueURL(queueName string) string {
	return fmt.Sprintf("%s/%s/%s", s.baseURL, s.accountID, queueName)
}

func (s *SQSStore) buildQueueARN(queueName string) string {
	return s.arnBuilder.SQS().Queue(queueName)
}

func (s *SQSStore) arnToQueueURL(arn string) string {
	_, _, _, _, queueName := svcarn.SplitARN(arn)
	if queueName == "" {
		return ""
	}
	return s.buildQueueURL(queueName)
}
