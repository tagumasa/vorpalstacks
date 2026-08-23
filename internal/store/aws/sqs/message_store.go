package sqs

import (
	"bytes"
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"

	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/core/storage"
	pb "vorpalstacks/internal/pb/storage/storage_sqs"
	"vorpalstacks/internal/store/aws/common"
)

const deduplicationCacheMaxSize = 500

// MessageSize returns the total size of a message in bytes: the body plus
// every message attribute component (name, data type, and value). Per the
// Amazon SQS Developer Guide, all components of a message attribute are
// included in the message size restriction.
func MessageSize(body string, attrs map[string]*MessageAttributeValue) int {
	size := len(body)
	for name, attr := range attrs {
		size += len(name)
		if attr == nil {
			continue
		}
		size += len(attr.DataType)
		if attr.StringValue != nil {
			size += len(*attr.StringValue)
		}
		if attr.BinaryValue != nil {
			size += len(attr.BinaryValue)
		}
	}
	return size
}

// messageSize is the message-based convenience wrapper around MessageSize.
func messageSize(message *Message) int {
	return MessageSize(message.Body, message.MessageAttributes)
}

// SendMessage sends a message to the specified queue.
func (s *SQSStore) SendMessage(queueURL string, message *Message) (*Message, error) {
	queue, err := s.GetQueue(queueURL)
	if err != nil {
		return nil, err
	}

	// All components of a message (body plus attribute names, data types,
	// and values) count towards the queue's MaximumMessageSize.
	if int32(messageSize(message)) > queue.MaximumMessageSize {
		return nil, ErrMessageTooLarge
	}

	if err := ValidateMessageAttributes(message.MessageAttributes); err != nil {
		return nil, err
	}

	if err := validateMessageBody(message.Body); err != nil {
		return nil, err
	}

	if err := validateFifoIdentifier(message.MessageGroupID); err != nil {
		return nil, err
	}
	if err := validateFifoIdentifier(message.MessageDeduplicationID); err != nil {
		return nil, err
	}

	if message.DelaySeconds < MinDelaySeconds || message.DelaySeconds > MaxDelaySeconds {
		return nil, ErrInvalidParameterValue
	}

	if queue.FifoQueue {
		if message.MessageGroupID == "" {
			return nil, ErrMissingMessageGroupId
		}
		if message.MessageDeduplicationID == "" && !queue.ContentBasedDeduplication {
			return nil, ErrMissingDeduplicationId
		}

		dedupKey := s.buildDeduplicationKey(queueURL, message)
		if dedupKey != "" {
			if msgID, ok := s.getDeduplicationMessageID(dedupKey); ok {
				var existingMsgPb pb.Message
				if err := s.messagesStore.GetProto(msgID, &existingMsgPb); err == nil {
					return ProtoToMessage(&existingMsgPb), nil
				}
			}
		}
	}

	message.ID = uuid.New().String()
	message.QueueURL = queueURL
	message.QueueARN = queue.ARN
	message.SentTimestamp = time.Now().UTC()
	message.ReceiptHandle = ""
	message.ApproximateReceiveCount = 0
	message.MD5OfBody = calculateMD5(message.Body)

	if message.DelaySeconds > 0 {
		message.VisibleAfter = message.SentTimestamp.Add(time.Duration(message.DelaySeconds) * time.Second)
	} else if queue.DelaySeconds > 0 {
		message.DelaySeconds = queue.DelaySeconds
		message.VisibleAfter = message.SentTimestamp.Add(time.Duration(queue.DelaySeconds) * time.Second)
	}

	if message.MessageAttributes == nil {
		message.MessageAttributes = make(map[string]*MessageAttributeValue)
	}
	message.MD5OfMessageAttributes = calculateMessageAttributesMD5(message.MessageAttributes)

	if message.Attributes == nil {
		message.Attributes = make(map[string]string)
	}
	message.Attributes["SenderId"] = s.accountID
	message.Attributes["SentTimestamp"] = fmt.Sprintf("%d", message.SentTimestamp.UnixMilli())

	if queue.FifoQueue {
		s.sequenceMu.Lock()
		counter := s.sequenceCounters[queueURL]
		if counter == 0 {
			counter = time.Now().UnixNano()
		}
		counter++
		s.sequenceCounters[queueURL] = counter
		message.SequenceNumber = fmt.Sprintf("%d", counter)
		s.sequenceMu.Unlock()
	}

	if err := s.messagesStore.PutProto(messageKey(queueURL, message.ID), MessageToProto(message)); err != nil {
		return nil, err
	}

	if queue.FifoQueue {
		dedupKey := s.buildDeduplicationKey(queueURL, message)
		if dedupKey != "" {
			s.putDeduplicationEntry(dedupKey, messageKey(queueURL, message.ID))
		}
	}

	return message, nil
}

// ReceiveMessage retrieves messages from the specified queue. A negative
// waitTimeSeconds means "unset" and selects the queue's
// ReceiveMessageWaitTimeSeconds attribute; a non-negative value is used
// directly. When the effective wait is positive the call long-polls: it
// rescans the queue until messages are available or the wait deadline
// expires ("The duration (in seconds) for which the call waits for a message
// to arrive in the queue before returning. If a message is available, the
// call returns sooner than WaitTimeSeconds.").
func (s *SQSStore) ReceiveMessage(queueURL string, maxNumberOfMessages int32, visibilityTimeoutPtr *int32, waitTimeSeconds int32, receiveRequestAttemptId string) ([]*Message, error) {
	queue, err := s.GetQueue(queueURL)
	if err != nil {
		return nil, err
	}

	if maxNumberOfMessages < MinMaxNumberOfMessages || maxNumberOfMessages > MaxMaxNumberOfMessages {
		return nil, ErrInvalidParameterValue
	}

	visibilityTimeout := queue.VisibilityTimeout
	if visibilityTimeoutPtr != nil {
		visibilityTimeout = *visibilityTimeoutPtr
	}
	if err := validateVisibilityTimeout(visibilityTimeout); err != nil {
		return nil, err
	}
	wait := waitTimeSeconds
	if wait < 0 {
		wait = queue.ReceiveMessageWaitTimeSeconds
	}
	if err := validateReceiveMessageWaitTimeSeconds(wait); err != nil {
		return nil, err
	}

	deadline := time.Now().Add(time.Duration(wait) * time.Second)
	for {
		messages, err := s.receiveMessagesOnce(queue, queueURL, maxNumberOfMessages, visibilityTimeout, receiveRequestAttemptId)
		if err != nil {
			return nil, err
		}
		if len(messages) > 0 || !time.Now().Before(deadline) {
			return messages, nil
		}
		// Poll in short slices so a server shutdown is observed promptly
		// instead of blocking out the whole wait.
		select {
		case <-s.ctx.Done():
			return messages, nil
		case <-time.After(receivePollInterval):
		}
	}
}

// receivePollInterval is the rescan interval of the long-polling receive
// loop. The message mutex is released between scans so concurrent senders
// and receivers make progress.
const receivePollInterval = 200 * time.Millisecond

// receiveMessagesOnce performs a single receive scan of the queue.
func (s *SQSStore) receiveMessagesOnce(queue *Queue, queueURL string, maxNumberOfMessages int32, visibilityTimeout int32, receiveRequestAttemptId string) ([]*Message, error) {
	now := time.Now().UTC()
	s.msgMutex.Lock()
	defer s.msgMutex.Unlock()

	// FIFO receive dedup via ReceiveRequestAttemptId
	if queue.FifoQueue && receiveRequestAttemptId != "" {
		if cached := s.checkReceiveAttemptCache(queueURL, receiveRequestAttemptId, now); cached != nil {
			return cached, nil
		}
	}

	candidates := s.selectCandidates(queueURL, queue, maxNumberOfMessages, now)

	var messages []*Message
	for _, msgPb := range candidates {
		if int32(len(messages)) >= maxNumberOfMessages {
			break
		}

		msg := ProtoToMessage(msgPb)
		msg.ApproximateReceiveCount++

		if queue.RedrivePolicy != nil && msg.ApproximateReceiveCount > queue.RedrivePolicy.MaxReceiveCount {
			if err := s.moveToDLQ(msg, queue.RedrivePolicy.DeadLetterTargetARN); err == nil {
				continue
			}
			logs.Warn("Failed to move message to DLQ, delivering to consumer", logs.String("messageId", msg.ID))
		}

		// A new receipt handle is issued for every receive; previously
		// issued handles stay resolvable so that deletes with an older
		// handle still succeed ("If you use an old ReceiptHandle, the
		// request will succeed, but the message might not be deleted.").
		msg.ReceiptHandle = generateReceiptHandle()
		msg.VisibilityTimeout = visibilityTimeout
		msg.ReceivedAt = now
		if msg.ApproximateFirstReceiveTimestamp.IsZero() {
			msg.ApproximateFirstReceiveTimestamp = now
		}

		msg.Attributes["ApproximateReceiveCount"] = fmt.Sprintf("%d", msg.ApproximateReceiveCount)
		msg.Attributes["ApproximateFirstReceiveTimestamp"] = fmt.Sprintf("%d", msg.ApproximateFirstReceiveTimestamp.UnixMilli())
		if msg.SequenceNumber != "" {
			msg.Attributes["SequenceNumber"] = msg.SequenceNumber
		}
		if msg.MessageGroupID != "" {
			msg.Attributes["MessageGroupId"] = msg.MessageGroupID
		}
		if msg.MessageDeduplicationID != "" {
			msg.Attributes["MessageDeduplicationId"] = msg.MessageDeduplicationID
		}

		if err := s.messagesStore.PutProto(messageKey(queueURL, msg.ID), MessageToProto(msg)); err != nil {
			continue
		}

		_ = s.storage.Bucket("sqs-receipts-"+s.region).Put(
			[]byte(msg.ReceiptHandle),
			[]byte(messageKey(queueURL, msg.ID)),
		)

		messages = append(messages, msg)
	}

	// Cache the result for FIFO receive dedup
	if queue.FifoQueue && receiveRequestAttemptId != "" && len(messages) > 0 {
		s.cacheReceiveAttempt(queueURL, receiveRequestAttemptId, messages)
	}

	return messages, nil
}

// selectCandidates returns candidate messages for ReceiveMessage, applying
// standard or FIFO-specific filtering depending on queue type.
func (s *SQSStore) selectCandidates(queueURL string, queue *Queue, maxItems int32, now time.Time) []*pb.Message {
	if !queue.FifoQueue {
		// Standard queues have no cross-message ordering constraints, so a
		// bounded scan that returns the first visible messages suffices.
		scanLimit := int(maxItems) * 10
		return s.selectStandardCandidates(queueURL, queue, scanLimit, now)
	}
	return s.selectFIFOCandidates(queueURL, queue, maxItems, now)
}

// selectStandardCandidates returns visible, non-expired messages for a
// standard (non-FIFO) queue using a single filtered scan.
func (s *SQSStore) selectStandardCandidates(queueURL string, queue *Queue, scanLimit int, now time.Time) []*pb.Message {
	retentionCutoff := now.Add(-time.Duration(queue.MessageRetentionPeriod) * time.Second)
	opts := common.ListOptions{Prefix: messagePrefix(queueURL), MaxItems: scanLimit}
	result, err := common.ListProto[*pb.Message](s.messagesStore, opts, func() *pb.Message { return &pb.Message{} }, func(m *pb.Message) bool {
		return s.isMessageVisible(m, now) && !s.isMessageExpired(m, retentionCutoff)
	})
	if err != nil {
		return nil
	}
	return result.Items
}

// selectFIFOCandidates returns visible, non-expired messages for a FIFO queue,
// sorted by SentTimestamp, with at most one message per MessageGroupID, and
// excluding groups that already have an in-flight message.
func (s *SQSStore) selectFIFOCandidates(queueURL string, queue *Queue, maxItems int32, now time.Time) []*pb.Message {
	retentionCutoff := now.Add(-time.Duration(queue.MessageRetentionPeriod) * time.Second)

	// "Messages within the same message group are always processed one at a
	// time, in strict order" — both the in-flight group exclusion and the
	// send-order sort are queue-wide guarantees, so the whole queue is
	// scanned rather than a bounded window.
	allMessages := make([]*pb.Message, 0)
	_ = common.ForEachAllProto(s.messagesStore, messagePrefix(queueURL),
		func() *pb.Message { return &pb.Message{} }, nil,
		func(m *pb.Message) error {
			allMessages = append(allMessages, m)
			return nil
		},
	)

	inFlightGroups := make(map[string]bool)
	for _, m := range allMessages {
		if s.isMessageInFlight(m, now) && m.MessageGroupId != "" {
			inFlightGroups[m.MessageGroupId] = true
		}
	}

	var visible []*pb.Message
	for _, m := range allMessages {
		if !s.isMessageVisible(m, now) || s.isMessageExpired(m, retentionCutoff) {
			continue
		}
		if m.MessageGroupId != "" && inFlightGroups[m.MessageGroupId] {
			continue
		}
		visible = append(visible, m)
	}

	sort.Slice(visible, func(i, j int) bool {
		return protoToTime(visible[i].SentTimestamp).Before(protoToTime(visible[j].SentTimestamp))
	})

	seenGroups := make(map[string]bool)
	var candidates []*pb.Message
	for _, m := range visible {
		if int32(len(candidates)) >= maxItems {
			break
		}
		if m.MessageGroupId != "" {
			if seenGroups[m.MessageGroupId] {
				continue
			}
			seenGroups[m.MessageGroupId] = true
		}
		candidates = append(candidates, m)
	}
	return candidates
}

// isMessageVisible returns true if the message is not delayed and not
// currently in-flight (its visibility timeout has not expired).
func (s *SQSStore) isMessageVisible(m *pb.Message, now time.Time) bool {
	visibleAfter := protoToTime(m.VisibleAfter)
	if !visibleAfter.IsZero() && now.Before(visibleAfter) {
		return false
	}
	return !s.isMessageInFlight(m, now)
}

// isMessageInFlight returns true if the message has been received and its
// visibility timeout has not yet expired.
func (s *SQSStore) isMessageInFlight(m *pb.Message, now time.Time) bool {
	if m.ReceiptHandle == "" {
		return false
	}
	receivedAt := protoToTime(m.ReceivedAt)
	if receivedAt.IsZero() {
		return false
	}
	return now.Before(receivedAt.Add(time.Duration(m.VisibilityTimeout) * time.Second))
}

// isMessageExpired returns true if the message was sent before the retention
// cutoff and should no longer be delivered.
func (s *SQSStore) isMessageExpired(m *pb.Message, retentionCutoff time.Time) bool {
	sentTime := protoToTime(m.SentTimestamp)
	if sentTime.IsZero() {
		return false
	}
	return sentTime.Before(retentionCutoff)
}

// DeleteMessage deletes a message from the queue using the receipt handle.
func (s *SQSStore) DeleteMessage(queueURL, receiptHandle string) error {
	if !s.Exists(queueURL) {
		return ErrQueueNotFound
	}

	s.msgMutex.Lock()
	defer s.msgMutex.Unlock()

	receiptsBucket := s.storage.Bucket("sqs-receipts-" + s.region)
	msgKeyBytes, err := receiptsBucket.Get([]byte(receiptHandle))
	if err != nil || len(msgKeyBytes) == 0 {
		return ErrInvalidReceiptHandle
	}
	// Receipt handles are scoped to their queue: a handle issued by another
	// queue must not delete this queue's messages.
	if !bytes.HasPrefix(msgKeyBytes, []byte(messagePrefix(queueURL))) {
		return ErrInvalidReceiptHandle
	}
	messagesBucket := "sqs-messages-" + s.region
	return s.storage.Update(context.Background(), func(txn storage.Transaction) error {
		if err := txn.Bucket(messagesBucket).Delete(msgKeyBytes); err != nil {
			return err
		}
		return txn.Bucket("sqs-receipts-" + s.region).Delete([]byte(receiptHandle))
	})
}

// ChangeMessageVisibility changes the visibility timeout of a message.
func (s *SQSStore) ChangeMessageVisibility(queueURL, receiptHandle string, visibilityTimeout int32) error {
	if !s.Exists(queueURL) {
		return ErrQueueNotFound
	}

	if err := validateVisibilityTimeout(visibilityTimeout); err != nil {
		return err
	}

	s.msgMutex.Lock()
	defer s.msgMutex.Unlock()

	receiptsBucket := s.storage.Bucket("sqs-receipts-" + s.region)
	msgKeyBytes, err := receiptsBucket.Get([]byte(receiptHandle))
	if err != nil || len(msgKeyBytes) == 0 {
		return ErrInvalidReceiptHandle
	}
	// Receipt handles are scoped to their queue: a handle issued by another
	// queue must not mutate this queue's messages.
	if !bytes.HasPrefix(msgKeyBytes, []byte(messagePrefix(queueURL))) {
		return ErrInvalidReceiptHandle
	}

	var msgPb pb.Message
	if err := s.messagesStore.GetProto(string(msgKeyBytes), &msgPb); err != nil {
		return ErrInvalidReceiptHandle
	}

	// Only in-flight messages can have their visibility changed: the AWS
	// contract returns MessageNotInflight once the visibility timeout has
	// expired ("The specified message isn't in flight.").
	if !s.isMessageInFlight(&msgPb, time.Now()) {
		return ErrMessageNotInflight
	}

	msg := ProtoToMessage(&msgPb)
	msg.VisibilityTimeout = visibilityTimeout
	msg.ReceivedAt = time.Now()

	if visibilityTimeout == 0 {
		msg.ReceiptHandle = ""
	}

	msgData, marshalErr := proto.Marshal(MessageToProto(msg))
	if marshalErr != nil {
		return marshalErr
	}

	messagesBucket := "sqs-messages-" + s.region
	return s.storage.Update(context.Background(), func(txn storage.Transaction) error {
		if err := txn.Bucket(messagesBucket).Put(msgKeyBytes, msgData); err != nil {
			return err
		}
		if visibilityTimeout == 0 {
			return txn.Bucket("sqs-receipts-" + s.region).Delete([]byte(receiptHandle))
		}
		return nil
	})
}

// PurgeQueue removes all messages from the specified queue.
func (s *SQSStore) PurgeQueue(queueURL string) error {
	if !s.Exists(queueURL) {
		return ErrQueueNotFound
	}

	s.purgeMutex.Lock()

	if startTime, inProgress := s.purgeInProgress[queueURL]; inProgress {
		if time.Since(startTime) < purgeTimeout {
			s.purgeMutex.Unlock()
			return ErrPurgeQueueInProgress
		}
	}

	for key, startTime := range s.purgeInProgress {
		if time.Since(startTime) >= purgeTimeout {
			delete(s.purgeInProgress, key)
		}
	}

	s.purgeInProgress[queueURL] = time.Now()
	s.purgeMutex.Unlock()

	// The cooldown marker is deliberately NOT removed when the purge
	// completes: the documented window is measured from the previous
	// PurgeQueue request ("the specified queue previously received a
	// PurgeQueue request within the last 60 seconds"), so it expires via the
	// stale-entry sweep above instead.

	s.msgMutex.Lock()
	defer s.msgMutex.Unlock()

	err := s.messagesStore.DeleteByPrefix(messagePrefix(queueURL))
	if err != nil {
		return err
	}

	receiptsBucket := s.storage.Bucket("sqs-receipts-" + s.region)
	prefix := []byte(messagePrefix(queueURL))
	var toDelete [][]byte
	_ = receiptsBucket.ForEach(func(k, v []byte) error {
		if bytes.HasPrefix(v, prefix) {
			keyCopy := make([]byte, len(k))
			copy(keyCopy, k)
			toDelete = append(toDelete, keyCopy)
		}
		return nil
	})
	for _, k := range toDelete {
		_ = receiptsBucket.Delete(k)
	}

	return nil
}

// ---------------------------------------------------------------------------
// ReceiveRequestAttemptId FIFO receive dedup
// ---------------------------------------------------------------------------

const receiveAttemptCacheMaxSize = 500

// checkReceiveAttemptCache returns cached messages from a previous receive with
// the same attempt ID, if all messages are still in-flight. Returns nil on miss.
func (s *SQSStore) checkReceiveAttemptCache(queueURL, attemptId string, now time.Time) []*Message {
	key := queueURL + "#" + attemptId
	s.receiveAttemptMu.Lock()
	entry, exists := s.receiveAttemptCache[key]
	s.receiveAttemptMu.Unlock()

	if !exists {
		return nil
	}

	var messages []*Message
	for _, msgID := range entry.messageIDs {
		var msgPb pb.Message
		if err := s.messagesStore.GetProto(messageKey(queueURL, msgID), &msgPb); err != nil {
			return nil
		}
		if !s.isMessageInFlight(&msgPb, now) {
			return nil
		}
		messages = append(messages, ProtoToMessage(&msgPb))
	}

	return messages
}

// cacheReceiveAttempt stores the result of a FIFO receive for dedup.
func (s *SQSStore) cacheReceiveAttempt(queueURL, attemptId string, messages []*Message) {
	key := queueURL + "#" + attemptId
	msgIDs := make([]string, len(messages))
	for i, msg := range messages {
		msgIDs[i] = msg.ID
	}
	s.receiveAttemptMu.Lock()
	s.receiveAttemptCache[key] = &receiveAttemptEntry{
		messageIDs: msgIDs,
		createdAt:  time.Now(),
	}
	if len(s.receiveAttemptCache) > receiveAttemptCacheMaxSize {
		s.cleanupReceiveAttemptCache()
	}
	s.receiveAttemptMu.Unlock()
}

func (s *SQSStore) cleanupReceiveAttemptCache() {
	cutoff := time.Now().Add(-time.Duration(MaxVisibilityTimeout) * time.Second)
	for k, entry := range s.receiveAttemptCache {
		if entry.createdAt.Before(cutoff) {
			delete(s.receiveAttemptCache, k)
		}
	}
}
