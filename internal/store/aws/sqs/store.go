package sqs

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"

	"vorpalstacks/internal/core/storage"
	pb "vorpalstacks/internal/pb/storage/storage_sqs"
	"vorpalstacks/internal/store/aws/common"
	svcarn "vorpalstacks/internal/utils/aws/arn"
)

var _ = pb.Queue{}

const (
	maxQueueNameLength    = 80
	maxMessageBodySize    = 262144
	maxBatchSize          = 10
	maxBatchEntryIdLength = 80
	maxTagsPerQueue       = 50
	maxTagKeyLength       = 128
	maxTagValueLength     = 256
)

func messageKey(queueURL, messageID string) string {
	return queueURL + "\x00" + messageID
}

func messagePrefix(queueURL string) string {
	return queueURL + "\x00"
}

const (
	purgeTimeout              = 60 * time.Second
	minVisibilityTimeout      = 0
	maxVisibilityTimeout      = 43200
	minDelaySeconds           = 0
	maxDelaySeconds           = 900
	minMessageRetentionPeriod = 60
	maxMessageRetentionPeriod = 1209600
	minMaximumMessageSize     = 1024
	maxMaximumMessageSize     = 262144
	minReceiveMessageWaitTime = 0
	maxReceiveMessageWaitTime = 20
	minMaxNumberOfMessages    = 1
	maxMaxNumberOfMessages    = 10
	deduplicationWindow       = 5 * time.Minute
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
	purgeInProgress    map[string]time.Time
	storage            storage.BasicStorage
	deduplicationCache map[string]*deduplicationEntry
	deduplicationMu    sync.RWMutex
	sequenceCounter    int64
	sequenceMu         sync.Mutex
}

type deduplicationEntry struct {
	messageID string
	expiresAt time.Time
}

// NewSQSStore creates a new SQS store with the specified storage, account ID, region, and base URL.
func NewSQSStore(store storage.BasicStorage, accountID, region, baseURL string) *SQSStore {
	regionSuffix := "-" + region
	return &SQSStore{
		BaseStore:          common.NewBaseStore(store.Bucket("sqs-queues"+regionSuffix), "sqs-queues"),
		messagesStore:      common.NewBaseStore(store.Bucket("sqs-messages"+regionSuffix), "sqs-messages"),
		tasksStore:         common.NewBaseStore(store.Bucket("sqs-move-tasks"+regionSuffix), "sqs-move-tasks"),
		TagStore:           common.NewTagStoreWithRegion(store, "sqs", region),
		arnBuilder:         svcarn.NewARNBuilder(accountID, region),
		accountID:          accountID,
		region:             region,
		baseURL:            baseURL,
		storage:            store,
		purgeInProgress:    make(map[string]time.Time),
		deduplicationCache: make(map[string]*deduplicationEntry),
	}
}

// Storage returns the underlying storage for the SQS store.
func (s *SQSStore) Storage() storage.BasicStorage {
	return s.storage
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
	parts := strings.Split(arn, ":")
	if len(parts) < 6 {
		return ""
	}
	accountID := parts[4]
	queueName := parts[5]
	return fmt.Sprintf("https://sqs.%s.amazonaws.com/%s/%s", s.region, accountID, queueName)
}

func (s *SQSStore) buildDeduplicationKey(queueURL string, message *Message) string {
	if message.MessageDeduplicationID != "" {
		return queueURL + "#" + message.MessageDeduplicationID
	}
	return queueURL + "#" + calculateMD5(message.Body)
}

func (s *SQSStore) cleanupDeduplicationCache() {
	now := time.Now()
	deleted := 0
	const maxDeletesPerCleanup = 100
	for key, entry := range s.deduplicationCache {
		if now.After(entry.expiresAt) {
			delete(s.deduplicationCache, key)
			deleted++
			if deleted >= maxDeletesPerCleanup {
				break
			}
		}
	}
}

func (s *SQSStore) cleanupDeduplicationCacheForQueue(queueURL string) {
	prefix := queueURL + "#"
	for key := range s.deduplicationCache {
		if strings.HasPrefix(key, prefix) {
			delete(s.deduplicationCache, key)
		}
	}
}

func (s *SQSStore) cleanupDeduplicationCacheForQueueLocked(queueURL string) {
	s.deduplicationMu.Lock()
	defer s.deduplicationMu.Unlock()
	s.cleanupDeduplicationCacheForQueue(queueURL)
}

func generateReceiptHandle() string {
	return uuid.New().String() + "#" + fmt.Sprintf("%d", time.Now().UnixNano())
}

func (s *SQSStore) moveToDLQ(msg *Message, dlqARN string) error {
	dlqURL := s.arnToQueueURL(dlqARN)
	if dlqURL == "" {
		return fmt.Errorf("invalid DLQ ARN: %s", dlqARN)
	}

	newMsg := &Message{
		ID:                               uuid.New().String(),
		Body:                             msg.Body,
		MD5OfBody:                        msg.MD5OfBody,
		MD5OfMessageAttributes:           msg.MD5OfMessageAttributes,
		MessageAttributes:                msg.MessageAttributes,
		QueueURL:                         dlqURL,
		QueueARN:                         dlqARN,
		SentTimestamp:                    time.Now().UTC(),
		ApproximateReceiveCount:          0,
		ApproximateFirstReceiveTimestamp: time.Time{},
		Attributes:                       make(map[string]string),
		MessageDeduplicationID:           msg.MessageDeduplicationID,
		MessageGroupID:                   msg.MessageGroupID,
	}
	newMsg.Attributes["SenderId"] = s.accountID
	newMsg.Attributes["SentTimestamp"] = fmt.Sprintf("%d", newMsg.SentTimestamp.UnixMilli())

	if err := s.messagesStore.PutProto(messageKey(dlqURL, newMsg.ID), MessageToProto(newMsg)); err != nil {
		return fmt.Errorf("failed to put message to DLQ: %w", err)
	}

	if err := s.messagesStore.Delete(messageKey(msg.QueueURL, msg.ID)); err != nil {
		return fmt.Errorf("failed to delete original message after DLQ move: %w", err)
	}

	return nil
}

func calculateMessageAttributesMD5(attrs map[string]*MessageAttributeValue) string {
	if len(attrs) == 0 {
		return calculateMD5("")
	}

	keys := make([]string, 0, len(attrs))
	for k := range attrs {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var buf bytes.Buffer
	for _, k := range keys {
		v := attrs[k]
		if v == nil {
			continue
		}

		buf.Write(uint32ToBytes(uint32(len(k))))
		buf.WriteString(k)

		buf.Write(uint32ToBytes(uint32(len(v.DataType))))
		buf.WriteString(v.DataType)

		if v.StringValue != nil {
			buf.WriteByte(1)
			buf.Write(uint32ToBytes(uint32(len(*v.StringValue))))
			buf.WriteString(*v.StringValue)
		} else if v.BinaryValue != nil {
			buf.WriteByte(2)
			buf.Write(uint32ToBytes(uint32(len(v.BinaryValue))))
			buf.Write(v.BinaryValue)
		} else if len(v.StringListValues) > 0 {
			buf.WriteByte(3)
			buf.Write(uint32ToBytes(uint32(len(v.StringListValues))))
			for _, sv := range v.StringListValues {
				buf.Write(uint32ToBytes(uint32(len(sv))))
				buf.WriteString(sv)
			}
		} else if len(v.BinaryListValues) > 0 {
			buf.WriteByte(4)
			buf.Write(uint32ToBytes(uint32(len(v.BinaryListValues))))
			for _, bv := range v.BinaryListValues {
				buf.Write(uint32ToBytes(uint32(len(bv))))
				buf.Write(bv)
			}
		} else {
			buf.WriteByte(1)
			buf.Write(uint32ToBytes(0))
		}
	}

	return calculateMD5(buf.String())
}

func uint32ToBytes(n uint32) []byte {
	return []byte{
		byte(n >> 24),
		byte(n >> 16),
		byte(n >> 8),
		byte(n),
	}
}

// DecodeBinaryValue decodes a base64-encoded string into a byte slice.
func DecodeBinaryValue(encoded string) []byte {
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return nil
	}
	return decoded
}

// EncodeBinaryValue encodes a byte slice into a base64-encoded string.
func EncodeBinaryValue(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

// ListDeadLetterSourceQueues returns all queues that have the specified dead letter queue as their target.
func (s *SQSStore) ListDeadLetterSourceQueues(dlqARN string) ([]*Queue, error) {
	opts := common.ListOptions{MaxItems: 1000}
	result, err := common.ListProto[*pb.Queue](s.BaseStore, opts, func() *pb.Queue { return &pb.Queue{} }, func(q *pb.Queue) bool {
		return q.GetRedrivePolicy() != nil && q.GetRedrivePolicy().GetDeadLetterTargetArn() == dlqARN
	})
	if err != nil {
		return nil, err
	}

	queues := make([]*Queue, 0, len(result.Items))
	for _, pbQueue := range result.Items {
		queues = append(queues, ProtoToQueue(pbQueue))
	}
	return queues, nil
}

// GetMessageCounts returns the count of visible, not visible, and delayed messages for a queue.
func (s *SQSStore) GetMessageCounts(queueURL string) (visible, notVisible, delayed int32) {
	s.msgMutex.RLock()
	defer s.msgMutex.RUnlock()

	now := time.Now().UTC()
	opts := common.ListOptions{Prefix: messagePrefix(queueURL), MaxItems: 10000}

	result, err := common.ListProto[*pb.Message](s.messagesStore, opts, func() *pb.Message { return &pb.Message{} }, func(m *pb.Message) bool {
		return true
	})
	if err != nil {
		return 0, 0, 0
	}

	for _, msgPb := range result.Items {
		visibleAfter := protoToTime(msgPb.VisibleAfter)
		receivedAt := protoToTime(msgPb.ReceivedAt)
		if !visibleAfter.IsZero() && now.Before(visibleAfter) {
			delayed++
		} else if msgPb.ReceiptHandle != "" && !receivedAt.IsZero() && now.Before(receivedAt.Add(time.Duration(msgPb.VisibilityTimeout)*time.Second)) {
			notVisible++
		} else {
			visible++
		}
	}
	return visible, notVisible, delayed
}
