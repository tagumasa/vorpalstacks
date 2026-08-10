package sqs

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"

	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/core/storage"
	pb "vorpalstacks/internal/pb/storage/storage_sqs"
	"vorpalstacks/internal/store/aws/common"
	svcarn "vorpalstacks/internal/utils/aws/arn"
)

const (
	maxQueueNameLength    = 80
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
	MinMaximumMessageSize     = 1024
	MaxMaximumMessageSize     = 1048576
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
	arnBuilder          *svcarn.ARNBuilder
	accountID           string
	region              string
	baseURL             string
	msgMutex            sync.RWMutex
	purgeMutex          sync.Mutex
	queueMutex          sync.RWMutex
	taskMu              sync.Mutex
	purgeInProgress     map[string]time.Time
	storage             storage.TransactionalStorage
	deduplicationCache  map[string]*deduplicationEntry
	deduplicationMu     sync.RWMutex
	sequenceCounters    map[string]int64
	sequenceMu          sync.Mutex
	receiveAttemptCache map[string]*receiveAttemptEntry
	receiveAttemptMu    sync.Mutex
	ctx                 context.Context
	cancel              context.CancelFunc
	wg                  sync.WaitGroup
}

type deduplicationEntry struct {
	messageID string
	expiresAt time.Time
}

type receiveAttemptEntry struct {
	messageIDs []string
	createdAt  time.Time
}

// NewSQSStore creates a new SQS store with the specified storage, account ID, region, and base URL.
func NewSQSStore(store storage.BasicStorage, accountID, region, baseURL string) *SQSStore {
	ts, _ := store.(storage.TransactionalStorage)
	regionSuffix := "-" + region
	ctx, cancel := context.WithCancel(context.Background())
	s := &SQSStore{
		BaseStore:           common.NewBaseStore(store.Bucket("sqs-queues"+regionSuffix), "sqs-queues"),
		messagesStore:       common.NewBaseStore(store.Bucket("sqs-messages"+regionSuffix), "sqs-messages"),
		tasksStore:          common.NewBaseStore(store.Bucket("sqs-move-tasks"+regionSuffix), "sqs-move-tasks"),
		TagStore:            common.NewTagStoreWithRegion(store, "sqs", region),
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
	}
	s.wg.Add(1)
	go s.cleanupExpiredMessages()
	return s
}

// cleanupExpiredMessages periodically scans all queues and deletes messages
// that have exceeded their MessageRetentionPeriod. Expired messages are never
// delivered by ReceiveMessage (filtered by isMessageExpired), so this cleanup
// can run without holding msgMutex.
func (s *SQSStore) cleanupExpiredMessages() {
	defer s.wg.Done()
	defer func() {
		if r := recover(); r != nil {
			logs.Error("SQS: panic in cleanupExpiredMessages goroutine", logs.Any("panic", r))
		}
	}()
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()
	for {
		select {
		case <-s.ctx.Done():
			return
		case <-ticker.C:
			s.doMessageRetentionCleanup()
		}
	}
}

// doMessageRetentionCleanup scans ALL queues (with pagination) and deletes
// expired messages from each. Expired messages are never delivered by
// ReceiveMessage (filtered by isMessageExpired), so this cleanup can run
// without holding msgMutex. Errors are logged but do not stop the cleanup.
func (s *SQSStore) doMessageRetentionCleanup() {
	now := time.Now().UTC()
	const pageSize = 100
	var marker string

	for {
		if s.ctx.Err() != nil {
			return
		}

		opts := common.ListOptions{MaxItems: pageSize}
		if marker != "" {
			opts.Marker = marker
		}

		result, err := s.ListQueues(opts)
		if err != nil {
			logs.Warn("SQS retention cleanup: ListQueues failed", logs.Err(err))
			return
		}

		for _, queue := range result.Items {
			if s.ctx.Err() != nil {
				return
			}
			retentionCutoff := now.Add(-time.Duration(queue.MessageRetentionPeriod) * time.Second)
			prefix := messagePrefix(queue.URL)

			var expiredKeys []string
			if err := common.ForEachAllProto[*pb.Message](s.messagesStore, prefix, func() *pb.Message { return &pb.Message{} }, nil, func(msgPb *pb.Message) error {
				if s.isMessageExpired(msgPb, retentionCutoff) {
					expiredKeys = append(expiredKeys, messageKey(queue.URL, msgPb.Id))
				}
				return nil
			}); err != nil {
				logs.Warn("SQS retention cleanup: message scan failed", logs.String("queue", queue.URL), logs.Err(err))
				continue
			}

			for _, key := range expiredKeys {
				if err := s.messagesStore.Delete(key); err != nil {
					logs.Warn("SQS retention cleanup: delete failed", logs.String("key", key), logs.Err(err))
				}
			}
		}

		if !result.IsTruncated || result.NextMarker == "" {
			break
		}
		marker = result.NextMarker
	}
}

// Close stops the background cleanup goroutine.
func (s *SQSStore) Close() {
	if s.cancel != nil {
		s.cancel()
		s.wg.Wait()
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
	return s.buildQueueURL(parts[5])
}

func (s *SQSStore) buildDeduplicationKey(queueURL string, message *Message) string {
	if message.MessageDeduplicationID != "" {
		return queueURL + "#" + message.MessageDeduplicationID
	}
	return queueURL + "#" + calculateMD5(message.Body)
}

func (s *SQSStore) getDeduplicationMessageID(dedupKey string) (string, bool) {
	s.deduplicationMu.RLock()
	entry, exists := s.deduplicationCache[dedupKey]
	s.deduplicationMu.RUnlock()

	if exists && time.Now().Before(entry.expiresAt) {
		return entry.messageID, true
	}

	data, err := s.storage.Bucket("sqs-dedup-" + s.region).Get([]byte(dedupKey))
	if err == nil && data != nil {
		idx := bytes.IndexByte(data, '\x01')
		if idx > 0 {
			msgKey := string(data[:idx])
			expiryStr := string(data[idx+1:])
			if expiryMs, err := strconv.ParseInt(expiryStr, 10, 64); err == nil {
				if time.Now().UnixMilli() >= expiryMs {
					_ = s.storage.Bucket("sqs-dedup-" + s.region).Delete([]byte(dedupKey))
					s.deduplicationMu.Lock()
					delete(s.deduplicationCache, dedupKey)
					s.deduplicationMu.Unlock()
					return "", false
				}
			}
			s.deduplicationMu.Lock()
			s.deduplicationCache[dedupKey] = &deduplicationEntry{
				messageID: msgKey,
				expiresAt: time.Now().Add(deduplicationWindow),
			}
			s.deduplicationMu.Unlock()
			return msgKey, true
		}
	}

	return "", false
}

func (s *SQSStore) putDeduplicationEntry(dedupKey, messageID string) {
	s.deduplicationMu.Lock()
	s.deduplicationCache[dedupKey] = &deduplicationEntry{
		messageID: messageID,
		expiresAt: time.Now().Add(deduplicationWindow),
	}
	if len(s.deduplicationCache) > deduplicationCacheMaxSize {
		s.cleanupDeduplicationCache()
	}
	s.deduplicationMu.Unlock()

	expiry := time.Now().Add(deduplicationWindow).UnixMilli()
	val := append([]byte(messageID), '\x01')
	val = append(val, []byte(strconv.FormatInt(expiry, 10))...)
	_ = s.storage.Bucket("sqs-dedup-"+s.region).Put([]byte(dedupKey), val)
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

func generateReceiptHandle() string {
	return uuid.New().String() + "#" + fmt.Sprintf("%d", time.Now().UnixNano())
}

func (s *SQSStore) moveToDLQ(msg *Message, dlqARN string) error {
	dlqURL := s.arnToQueueURL(dlqARN)
	if dlqURL == "" {
		return fmt.Errorf("invalid DLQ ARN: %s", dlqARN)
	}

	// Preserve original SentTimestamp for standard queues. FIFO messages
	// always have a MessageGroupID; only reset SentTimestamp for those.
	sentTimestamp := msg.SentTimestamp
	if msg.MessageGroupID != "" {
		sentTimestamp = time.Now().UTC()
	}

	newMsg := &Message{
		ID:                               msg.ID,
		Body:                             msg.Body,
		MD5OfBody:                        msg.MD5OfBody,
		MD5OfMessageAttributes:           msg.MD5OfMessageAttributes,
		MessageAttributes:                msg.MessageAttributes,
		QueueURL:                         dlqURL,
		QueueARN:                         dlqARN,
		SentTimestamp:                    sentTimestamp,
		ApproximateReceiveCount:          0,
		ApproximateFirstReceiveTimestamp: time.Time{},
		Attributes:                       make(map[string]string),
		MessageDeduplicationID:           msg.MessageDeduplicationID,
		MessageGroupID:                   msg.MessageGroupID,
	}
	newMsg.Attributes["SenderId"] = s.accountID
	newMsg.Attributes["SentTimestamp"] = fmt.Sprintf("%d", newMsg.SentTimestamp.UnixMilli())

	dlqKey := messageKey(dlqURL, newMsg.ID)
	srcKey := messageKey(msg.QueueURL, msg.ID)
	dlqData, marshalErr := proto.Marshal(MessageToProto(newMsg))
	if marshalErr != nil {
		return fmt.Errorf("failed to marshal DLQ message: %w", marshalErr)
	}
	messagesBucket := "sqs-messages-" + s.region
	receiptsBucket := "sqs-receipts-" + s.region
	return s.storage.Update(context.Background(), func(txn storage.Transaction) error {
		if err := txn.Bucket(messagesBucket).Put([]byte(dlqKey), dlqData); err != nil {
			return err
		}
		if err := txn.Bucket(messagesBucket).Delete([]byte(srcKey)); err != nil {
			return err
		}
		if msg.ReceiptHandle != "" {
			_ = txn.Bucket(receiptsBucket).Delete([]byte(msg.ReceiptHandle))
		}
		return nil
	})
}

func calculateMessageAttributesMD5(attrs map[string]*MessageAttributeValue) string {
	return CalculateMessageAttributesMD5(attrs)
}

// CalculateMessageAttributesMD5 computes the MD5 digest of message attributes
// per the AWS SQS specification. Exported for cross-package use.
func CalculateMessageAttributesMD5(attrs map[string]*MessageAttributeValue) string {
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
func DecodeBinaryValue(encoded string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(encoded)
}

// EncodeBinaryValue encodes a byte slice into a base64-encoded string.
func EncodeBinaryValue(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

// ListDeadLetterSourceQueues returns all queues that have the specified dead letter queue as their target.
func (s *SQSStore) ListDeadLetterSourceQueues(dlqARN string) ([]*Queue, error) {
	items, err := common.ListMatchingProto[*pb.Queue](s.BaseStore, "", func() *pb.Queue { return &pb.Queue{} }, func(q *pb.Queue) bool {
		return q.GetRedrivePolicy() != nil && q.GetRedrivePolicy().GetDeadLetterTargetArn() == dlqARN
	})
	if err != nil {
		return nil, err
	}

	queues := make([]*Queue, 0, len(items))
	for _, pbQueue := range items {
		queues = append(queues, ProtoToQueue(pbQueue))
	}
	return queues, nil
}

// GetMessageCounts returns the count of visible, not visible, and delayed messages for a queue.
func (s *SQSStore) GetMessageCounts(queueURL string) (visible, notVisible, delayed int32) {
	s.msgMutex.RLock()
	defer s.msgMutex.RUnlock()

	now := time.Now().UTC()
	prefix := messagePrefix(queueURL)

	err := common.ForEachAllProto[*pb.Message](s.messagesStore, prefix, func() *pb.Message { return &pb.Message{} }, nil, func(msgPb *pb.Message) error {
		visibleAfter := protoToTime(msgPb.VisibleAfter)
		receivedAt := protoToTime(msgPb.ReceivedAt)
		if !visibleAfter.IsZero() && now.Before(visibleAfter) {
			delayed++
		} else if msgPb.ReceiptHandle != "" && !receivedAt.IsZero() && now.Before(receivedAt.Add(time.Duration(msgPb.VisibilityTimeout)*time.Second)) {
			notVisible++
		} else {
			visible++
		}
		return nil
	})
	if err != nil {
		logs.Error("SQS: failed to count messages", logs.String("queueUrl", queueURL), logs.Err(err))
		return 0, 0, 0
	}
	return visible, notVisible, delayed
}
