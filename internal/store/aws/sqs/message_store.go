package sqs

import (
	"bytes"
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"

	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/core/storage"
	pb "vorpalstacks/internal/pb/storage/storage_sqs"
	"vorpalstacks/internal/store/aws/common"
)

// MessageSize returns the total size of a message in bytes: the body plus
// every message attribute component (name, data type, and value). Per the
// Amazon SQS Developer Guide, all components of a message attribute are
// included in the message size restriction. List values contribute nothing
// here while CalculateMessageAttributesMD5 encodes them: attribute list
// values are rejected on both wire paths (UnsupportedOperation), so both the
// omission and the encoding arms are unreachable for list-carrying requests.
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
	// The existence read, the persist and the deduplication registration run
	// under the queue lock's read side: DeleteQueue holds the write side
	// across its whole wipe, so a send that has resolved the queue cannot
	// land a message (or re-arm deduplication state) on a queue the wipe has
	// already emptied. Such an orphan would never be swept — retention
	// cleanup scans only live queues' prefixes — and would resurrect into a
	// recreated same-name queue. The read side serialises only with queue
	// mutators, never with other sends or receives.
	s.queueMutex.RLock()
	defer s.queueMutex.RUnlock()

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

	// The deduplication key is computed once here and reused after the
	// persist below to register the entry; both sites need the same value.
	var dedupKey string
	if queue.FifoQueue {
		if message.MessageGroupID == "" {
			return nil, ErrMissingMessageGroupId
		}
		if message.MessageDeduplicationID == "" && !queue.ContentBasedDeduplication {
			return nil, ErrMissingDeduplicationId
		}

		dedupKey = s.buildDeduplicationKey(queueURL, message)
	} else if message.MessageDeduplicationID != "" {
		// "This parameter applies only to FIFO (first-in-first-out)
		// queues" (SendMessage API reference) — a deduplication id on a
		// standard queue is a parameter that is not valid for the queue
		// type and is rejected instead of silently ignored. MessageGroupId
		// is different: it is documented for standard queues as well
		// (fair queues), so it stays accepted there.
		return nil, ErrInvalidParameterValue
	}
	// Deduplication applies to FIFO queues only; dedupKey is built solely
	// on that path (a standard-queue send carrying a deduplication id was
	// rejected above). The dedup check, the message persist and the entry
	// registration form one check-then-act sequence. Per-key mutual
	// exclusion closes the window where two concurrent same-key sends both
	// miss the lookup and both persist — duplicate delivery, the exact
	// failure the deduplication interval exists to prevent.
	if queue.FifoQueue {
		s.dedupKeyLocker.Lock(dedupKey)
		defer s.dedupKeyLocker.Unlock(dedupKey)

		if recordedKey, recordedSeq, ok := s.getDeduplicationMessageID(dedupKey); ok {
			var existingMsgPb pb.Message
			if err := s.messagesStore.GetProto(recordedKey, &existingMsgPb); err == nil {
				return ProtoToMessage(&existingMsgPb), nil
			}
			// The deduplication interval is independent of message
			// existence: "Amazon SQS continues to keep track of the message
			// deduplication ID even after the message is received and
			// deleted" (SendMessage API reference). The resend is
			// acknowledged with the original's MessageId and the request's
			// digests, and nothing is persisted or delivered. The recorded
			// sequence number rides along — the acknowledgment surface of a
			// resend does not depend on whether the original still exists.
			ack := *message
			// Every registered value is a messageKey ("<queueURL>\x00<id>"),
			// so the recorded key's tail is the original's MessageId.
			ack.ID = recordedKey[strings.LastIndexByte(recordedKey, '\x00')+1:]
			ack.QueueURL = queueURL
			ack.QueueARN = queue.ARN
			ack.MD5OfBody = calculateMD5(ack.Body)
			ack.SequenceNumber = recordedSeq
			if ack.MessageAttributes == nil {
				ack.MessageAttributes = make(map[string]*MessageAttributeValue)
			}
			ack.MD5OfMessageAttributes = CalculateMessageAttributesMD5(ack.MessageAttributes)
			return &ack, nil
		}
	}

	message.ID = uuid.New().String()
	message.QueueURL = queueURL
	message.QueueARN = queue.ARN
	if queue.FifoQueue {
		message.SentTimestamp, message.SequenceNumber = s.nextFIFOStamp(queueURL)
	} else {
		message.SentTimestamp = time.Now().UTC()
	}
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
	message.MD5OfMessageAttributes = CalculateMessageAttributesMD5(message.MessageAttributes)

	if message.Attributes == nil {
		message.Attributes = make(map[string]string)
	}
	message.Attributes["SenderId"] = s.accountID
	message.Attributes["SentTimestamp"] = fmt.Sprintf("%d", message.SentTimestamp.UnixMilli())

	// The message persist and the FIFO deduplication entry commit as one
	// transaction, and the failure of either surfaces as an error: a message
	// can never be queued with its deduplication window silently missing
	// (a same-key resend inside the window would then be delivered as new).
	msgData, marshalErr := proto.Marshal(MessageToProto(message))
	if marshalErr != nil {
		return nil, marshalErr
	}
	msgKey := messageKey(queueURL, message.ID)
	putMessage := func(txn storage.Transaction) error {
		if err := txn.Bucket(s.messagesBucketName).Put([]byte(msgKey), msgData); err != nil {
			return err
		}
		if queue.FifoQueue {
			return txn.Bucket(s.dedupBucket).Put([]byte(dedupKey), deduplicationEntryValue(msgKey, message.SequenceNumber))
		}
		return nil
	}
	if err := s.storage.Update(context.Background(), putMessage); err != nil {
		return nil, err
	}
	if queue.FifoQueue {
		s.registerDeduplicationEntry(dedupKey, msgKey, message.SequenceNumber)
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

	// "This parameter applies only to FIFO (first-in-first-out) queues"
	// (ReceiveMessage API reference) — the same queue-type contract that
	// rejects MessageDeduplicationId on standard queues: a receive-request
	// attempt id on a standard queue is a parameter that is not valid for
	// the queue type and is rejected instead of silently ignored.
	if !queue.FifoQueue && receiveRequestAttemptId != "" {
		return nil, ErrInvalidParameterValue
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
			if err := s.moveToDLQ(queue, msg); err == nil {
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
		// SqsManagedSseEnabled is a documented receivable system attribute:
		// its value is the queue's own encryption setting, surfaced per
		// message when the deprecated AttributeNames channel asks for it.
		if v, ok := queue.Attributes["SqsManagedSseEnabled"]; ok && v != "" {
			msg.Attributes["SqsManagedSseEnabled"] = v
		}
		if msg.SequenceNumber != "" {
			msg.Attributes["SequenceNumber"] = msg.SequenceNumber
		}
		if msg.MessageGroupID != "" {
			msg.Attributes["MessageGroupId"] = msg.MessageGroupID
		}
		if msg.MessageDeduplicationID != "" {
			msg.Attributes["MessageDeduplicationId"] = msg.MessageDeduplicationID
		}

		// The in-flight state and the handle that resolves it are one
		// transaction, like every other multi-write path (DeleteMessage,
		// ChangeMessageVisibility, moveToDLQ, moveMessageBatch): a failure
		// between two independent puts would leave an in-flight message
		// whose handle resolves to nothing — the consumer cannot delete it
		// and it redelivers after the visibility timeout, inflating
		// ApproximateReceiveCount toward the DLQ threshold.
		msgData, marshalErr := proto.Marshal(MessageToProto(msg))
		if marshalErr != nil {
			continue
		}
		msgKey := messageKey(queueURL, msg.ID)
		if err := s.storage.Update(context.Background(), func(txn storage.Transaction) error {
			if err := txn.Bucket(s.messagesBucketName).Put([]byte(msgKey), msgData); err != nil {
				return err
			}
			return txn.Bucket(s.receiptsBucket).Put([]byte(msg.ReceiptHandle), []byte(msgKey))
		}); err != nil {
			continue
		}

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
//
// Known scalability limit (evaluated, kept): this whole-queue scan runs under
// the store-global msgMutex on every receivePollInterval tick of every active
// FIFO long-poll, so an empty FIFO long-poll rescans the queue every 200 ms
// while serialising with every other queue's receive/delete/visibility/purge
// in the process. The scan depth is the ordering guarantee itself — bounding
// it breaks the queue-wide in-flight-group exclusion — and the two bounded
// alternatives were rejected as not low-risk: backing the rescan interval
// off while the queue is empty trades arrival-to-delivery latency on FIFO
// long-polls, and per-queue locking re-opens the lock matrix (dedup key
// locker, DeleteQueue wipe, purge) the delete-path sequencing closed. At
// this platform's single-node scale the limit stands as documented.
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

	msgKeyBytes, err := s.resolveQueueReceiptHandle(queueURL, receiptHandle)
	if err != nil {
		return err
	}
	messagesBucket := s.messagesBucketName
	return s.storage.Update(context.Background(), func(txn storage.Transaction) error {
		if err := txn.Bucket(messagesBucket).Delete(msgKeyBytes); err != nil {
			return err
		}
		return txn.Bucket(s.receiptsBucket).Delete([]byte(receiptHandle))
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

	msgKeyBytes, err := s.resolveQueueReceiptHandle(queueURL, receiptHandle)
	if err != nil {
		return err
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

	messagesBucket := s.messagesBucketName
	return s.storage.Update(context.Background(), func(txn storage.Transaction) error {
		if err := txn.Bucket(messagesBucket).Put(msgKeyBytes, msgData); err != nil {
			return err
		}
		if visibilityTimeout == 0 {
			return txn.Bucket(s.receiptsBucket).Delete([]byte(receiptHandle))
		}
		return nil
	})
}

// resolveQueueReceiptHandle resolves a receipt handle to its message key,
// enforcing the queue scoping: a handle resolves to a message key carrying
// the queue's prefix, so a handle issued by another queue is invalid here
// rather than a cross-queue delete or visibility change.
func (s *SQSStore) resolveQueueReceiptHandle(queueURL, receiptHandle string) ([]byte, error) {
	msgKeyBytes, err := s.storage.Bucket(s.receiptsBucket).Get([]byte(receiptHandle))
	if err != nil || len(msgKeyBytes) == 0 {
		return nil, ErrInvalidReceiptHandle
	}
	if !bytes.HasPrefix(msgKeyBytes, []byte(messagePrefix(queueURL))) {
		return nil, ErrInvalidReceiptHandle
	}
	return msgKeyBytes, nil
}

// deleteReceiptsForPrefix removes every receipt-handle entry that resolves
// into the given queue's message prefix — the receipts leg of a purge or a
// queue wipe.
func (s *SQSStore) deleteReceiptsForPrefix(queueURL string) {
	receiptsBucket := s.storage.Bucket(s.receiptsBucket)
	prefix := []byte(messagePrefix(queueURL))
	var toDelete [][]byte
	if err := receiptsBucket.ForEach(func(k, v []byte) error {
		if bytes.HasPrefix(v, prefix) {
			keyCopy := make([]byte, len(k))
			copy(keyCopy, k)
			toDelete = append(toDelete, keyCopy)
		}
		return nil
	}); err != nil {
		// Fail-open by design, never silently: the purge has already removed
		// the messages these handles resolve to, so a missed handle is inert
		// residue the 12 h age sweep reclaims.
		logs.Warn("SQS purge: receipt scan failed; stale handles stay until the age sweep", logs.String("queueURL", queueURL), logs.Err(err))
	}
	for _, k := range toDelete {
		if err := receiptsBucket.Delete(k); err != nil {
			// Same inert-residue reasoning as a missed scan entry.
			logs.Warn("SQS purge: receipt delete failed; the age sweep reclaims it", logs.String("handle", string(k)), logs.Err(err))
		}
	}
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

	// The helper derives the match prefix from the queue URL itself — the
	// raw URL, not the message prefix (DeleteQueue's receipt collection in
	// deleteQueueRecordLocked matches on messagePrefix(queueURL) the same
	// way; prefixing here would double the prefix and match nothing).
	s.deleteReceiptsForPrefix(queueURL)

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

func generateReceiptHandle() string {
	return uuid.New().String() + "#" + fmt.Sprintf("%d", time.Now().UnixNano())
}

// moveToDLQ moves a message past its receive-count limit into the queue's
// dead-letter target. The target is re-resolved at move time — a target that
// was deleted, was recreated as the other queue type, or is the source queue
// itself makes the move impossible, and the caller delivers the message to
// the consumer instead of losing it. The identity policy matches
// moveMessageBatch and the documented dead-letter-queue semantics: the
// message keeps its ID, a standard-queue message keeps its original enqueue
// timestamp ("For standard queues, the expiration of a message is always
// based on its original enqueue timestamp. When a message is moved to a
// dead-letter queue, the enqueue timestamp is unchanged."), and a FIFO
// message has the enqueue timestamp reset ("For FIFO queues, the enqueue
// timestamp resets when the message is moved to a dead-letter queue.")
// (AWS SQS Developer Guide). Note the standard-queue behaviour is deliberate:
// such a message is swept once it outruns the DLQ's own retention period,
// exactly as AWS documents on the same page.
func (s *SQSStore) moveToDLQ(queue *Queue, msg *Message) error {
	rdp := queue.RedrivePolicy
	if rdp == nil {
		return fmt.Errorf("queue %s has no redrive policy", queue.URL)
	}
	dlqURL := s.arnToQueueURL(rdp.DeadLetterTargetARN)
	if dlqURL == "" {
		return fmt.Errorf("invalid dead-letter target ARN: %s", rdp.DeadLetterTargetARN)
	}
	if dlqURL == queue.URL {
		return fmt.Errorf("dead-letter target is the source queue itself: %s", dlqURL)
	}
	dlq, err := s.GetQueue(dlqURL)
	if err != nil {
		return fmt.Errorf("resolving dead-letter target %s: %w", rdp.DeadLetterTargetARN, err)
	}
	if dlq.FifoQueue != queue.FifoQueue {
		return fmt.Errorf("dead-letter target %s has a mismatched queue type", rdp.DeadLetterTargetARN)
	}

	sentTimestamp := msg.SentTimestamp
	var dlqSequence string
	if dlq.FifoQueue {
		// The enqueue timestamp and the sequence number are one atomic
		// stamp on FIFO destinations (see nextFIFOStamp).
		sentTimestamp, dlqSequence = s.nextFIFOStamp(dlqURL)
	}

	newMsg := &Message{
		ID:                               msg.ID,
		Body:                             msg.Body,
		MD5OfBody:                        msg.MD5OfBody,
		MD5OfMessageAttributes:           msg.MD5OfMessageAttributes,
		MessageAttributes:                msg.MessageAttributes,
		QueueURL:                         dlqURL,
		QueueARN:                         dlq.ARN,
		SentTimestamp:                    sentTimestamp,
		ApproximateReceiveCount:          msg.ApproximateReceiveCount,
		ApproximateFirstReceiveTimestamp: time.Time{},
		Attributes:                       make(map[string]string),
		MessageDeduplicationID:           msg.MessageDeduplicationID,
		MessageGroupID:                   msg.MessageGroupID,
	}
	newMsg.Attributes["SenderId"] = s.accountID
	newMsg.Attributes["SentTimestamp"] = fmt.Sprintf("%d", newMsg.SentTimestamp.UnixMilli())
	// DeadLetterQueueSourceArn identifies the source queue on messages
	// delivered through a redrive policy; it is a documented receive system
	// attribute and the move task's omitted-destination routing key.
	newMsg.Attributes[deadLetterQueueSourceArnAttr] = msg.QueueARN
	// "ApproximateReceiveCount – Returns the number of times a message has
	// been received across all queues but not deleted" (ReceiveMessage API
	// reference): the count accumulates across the relocation, it does not
	// restart at the destination. ApproximateFirstReceiveTimestamp does
	// restart — its doc wording is queue-scoped ("first received from the
	// queue").
	newMsg.Attributes["ApproximateReceiveCount"] = fmt.Sprintf("%d", newMsg.ApproximateReceiveCount)
	// Every enqueued FIFO message carries a sequence number from its own
	// queue's counter — the policy redrive is an enqueue into the DLQ, like
	// the move-task path.
	newMsg.SequenceNumber = dlqSequence

	// Deduplication registration at a FIFO destination follows SendMessage's
	// gating (the destination opted in, or the message carries an explicit
	// deduplication ID): the redriven copy opens the window, so a same-key
	// send into the DLQ inside the interval is suppressed as the duplicate
	// of a message the queue already holds. The redrive itself does not
	// CONSULT the window — suppressing a redrive on an existing entry would
	// drop the message outright, a loss no source licenses — so the two
	// redrive paths register symmetrically while only the send paths
	// suppress.
	var dedupKey string
	if dlq.FifoQueue && (dlq.ContentBasedDeduplication || newMsg.MessageDeduplicationID != "") {
		dedupKey = s.buildDeduplicationKey(dlqURL, newMsg)
		s.dedupKeyLocker.Lock(dedupKey)
		defer s.dedupKeyLocker.Unlock(dedupKey)
	}

	dlqKey := messageKey(dlqURL, newMsg.ID)
	srcKey := messageKey(msg.QueueURL, msg.ID)
	dlqData, marshalErr := proto.Marshal(MessageToProto(newMsg))
	if marshalErr != nil {
		return fmt.Errorf("failed to marshal DLQ message: %w", marshalErr)
	}
	messagesBucket := s.messagesBucketName
	receiptsBucket := s.receiptsBucket
	dedupBucket := s.dedupBucket
	var staleReceiptErr error
	if err := s.storage.Update(context.Background(), func(txn storage.Transaction) error {
		if err := txn.Bucket(messagesBucket).Put([]byte(dlqKey), dlqData); err != nil {
			return err
		}
		if err := txn.Bucket(messagesBucket).Delete([]byte(srcKey)); err != nil {
			return err
		}
		if msg.ReceiptHandle != "" {
			// Inert on failure and never worth aborting the redrive: the
			// stale handle resolves to a deleted key and the 12 h age sweep
			// reclaims it. Logged after the commit, never silent.
			staleReceiptErr = txn.Bucket(receiptsBucket).Delete([]byte(msg.ReceiptHandle))
		}
		if dedupKey != "" {
			return txn.Bucket(dedupBucket).Put([]byte(dedupKey), deduplicationEntryValue(dlqKey, newMsg.SequenceNumber))
		}
		return nil
	}); err != nil {
		return err
	}
	if staleReceiptErr != nil {
		logs.Warn("SQS: redrive left a stale receipt entry behind; the age sweep reclaims it",
			logs.String("sourceUrl", msg.QueueURL), logs.Err(staleReceiptErr))
	}
	if dedupKey != "" {
		s.registerDeduplicationEntry(dedupKey, dlqKey, newMsg.SequenceNumber)
	}
	return nil
}

// nextFIFOStamp stamps a FIFO enqueue's timestamp and sequence number as
// one atomic section: FIFO delivery order is the SentTimestamp sort
// (selectFIFOCandidates), so the sequence must be monotonic with it —
// reading the timestamp and bumping the counter in separate sections lets
// two concurrent sends observe the two orderings disagreeing. (A backwards
// wall-clock step between two sends is an external hazard common to every
// timestamp the platform emits.) The counter seeds from the wall clock on
// first use so numbers stay large and increasing across restarts.
func (s *SQSStore) nextFIFOStamp(queueURL string) (time.Time, string) {
	s.sequenceMu.Lock()
	defer s.sequenceMu.Unlock()
	now := time.Now().UTC()
	counter := s.sequenceCounters[queueURL]
	if counter == 0 {
		counter = now.UnixNano()
	}
	counter++
	s.sequenceCounters[queueURL] = counter
	return now, fmt.Sprintf("%d", counter)
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
