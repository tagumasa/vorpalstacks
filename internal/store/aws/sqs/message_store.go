package sqs

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"

	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/core/storage"
	pb "vorpalstacks/internal/pb/storage/storage_sqs"
	"vorpalstacks/internal/store/aws/common"
)

const deduplicationCacheMaxSize = 500

// SendMessage sends a message to the specified queue.
func (s *SQSStore) SendMessage(queueURL string, message *Message) (*Message, error) {
	queue, err := s.GetQueue(queueURL)
	if err != nil {
		return nil, err
	}

	if err := validateMessageBodySize(message.Body); err != nil {
		return nil, err
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
		s.sequenceCounter++
		message.SequenceNumber = fmt.Sprintf("%d", s.sequenceCounter)
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

// ReceiveMessage retrieves messages from the specified queue.
func (s *SQSStore) ReceiveMessage(queueURL string, maxNumberOfMessages int32, visibilityTimeoutPtr *int32, waitTimeSeconds int32) ([]*Message, error) {
	queue, err := s.GetQueue(queueURL)
	if err != nil {
		return nil, err
	}

	if maxNumberOfMessages <= 0 {
		maxNumberOfMessages = 1
	}
	if maxNumberOfMessages > maxMaxNumberOfMessages {
		maxNumberOfMessages = maxMaxNumberOfMessages
	}

	var visibilityTimeout int32
	if visibilityTimeoutPtr != nil {
		visibilityTimeout = *visibilityTimeoutPtr
	} else {
		visibilityTimeout = queue.VisibilityTimeout
	}
	if err := validateVisibilityTimeout(visibilityTimeout); err != nil {
		return nil, err
	}

	if waitTimeSeconds > 0 {
		if err := validateReceiveMessageWaitTimeSeconds(waitTimeSeconds); err != nil {
			return nil, err
		}
	}

	now := time.Now().UTC()
	opts := common.ListOptions{Prefix: messagePrefix(queueURL), MaxItems: int(maxNumberOfMessages) * 3}

	s.msgMutex.Lock()
	defer s.msgMutex.Unlock()

	result, err := common.ListProto[*pb.Message](s.messagesStore, opts, func() *pb.Message { return &pb.Message{} }, func(m *pb.Message) bool {
		visibleAfter := protoToTime(m.VisibleAfter)
		if !visibleAfter.IsZero() && now.Before(visibleAfter) {
			return false
		}
		if m.ReceiptHandle != "" {
			receivedAt := protoToTime(m.ReceivedAt)
			if !receivedAt.IsZero() && now.Before(receivedAt.Add(time.Duration(m.VisibilityTimeout)*time.Second)) {
				return false
			}
		}
		return true
	})
	if err != nil {
		return nil, err
	}

	var messages []*Message
	for _, msgPb := range result.Items {
		if int32(len(messages)) >= maxNumberOfMessages {
			break
		}

		msg := ProtoToMessage(msgPb)
		msg.ApproximateReceiveCount++

		if queue.RedrivePolicy != nil && msg.ApproximateReceiveCount > queue.RedrivePolicy.MaxReceiveCount {
			if err := s.moveToDLQ(msg, queue.RedrivePolicy.DeadLetterTargetARN); err != nil {
				logs.Warn("Failed to move message to DLQ", logs.String("messageId", msg.ID), logs.Err(err))
			}
			continue
		}

		msg.ReceiptHandle = generateReceiptHandle()
		msg.VisibilityTimeout = visibilityTimeout
		msg.ReceivedAt = now
		if msg.ApproximateFirstReceiveTimestamp.IsZero() {
			msg.ApproximateFirstReceiveTimestamp = now
		}

		msg.Attributes["ApproximateReceiveCount"] = fmt.Sprintf("%d", msg.ApproximateReceiveCount)
		msg.Attributes["ApproximateFirstReceiveTimestamp"] = fmt.Sprintf("%d", msg.ApproximateFirstReceiveTimestamp.UnixMilli())
		if msg.MessageGroupID != "" {
			msg.Attributes["MessageGroupId"] = msg.MessageGroupID
		}
		if msg.MessageDeduplicationID != "" {
			msg.Attributes["MessageDeduplicationId"] = msg.MessageDeduplicationID
		}

		if err := s.messagesStore.PutProto(messageKey(queueURL, msg.ID), MessageToProto(msg)); err != nil {
			continue
		}

		if s.storage != nil {
			_ = s.storage.Bucket("sqs-receipts-"+s.region).Put(
				[]byte(msg.ReceiptHandle),
				[]byte(messageKey(queueURL, msg.ID)),
			)
		}

		messages = append(messages, msg)
	}

	return messages, nil
}

// DeleteMessage deletes a message from the queue using the receipt handle.
func (s *SQSStore) DeleteMessage(queueURL, receiptHandle string) error {
	s.msgMutex.Lock()
	defer s.msgMutex.Unlock()

	if s.storage != nil {
		receiptsBucket := s.storage.Bucket("sqs-receipts-" + s.region)
		msgKeyBytes, err := receiptsBucket.Get([]byte(receiptHandle))
		if err != nil || len(msgKeyBytes) == 0 {
			return ErrInvalidReceiptHandle
		}
		messagesBucket := "sqs-messages-" + s.region
		err = s.storage.Update(context.Background(), func(txn storage.Transaction) error {
			if err := txn.Bucket(messagesBucket).Delete(msgKeyBytes); err != nil {
				return err
			}
			return txn.Bucket("sqs-receipts-" + s.region).Delete([]byte(receiptHandle))
		})
		return err
	}

	opts := common.ListOptions{Prefix: messagePrefix(queueURL), MaxItems: 1000}
	result, err := common.ListProto[*pb.Message](s.messagesStore, opts, func() *pb.Message { return &pb.Message{} }, func(m *pb.Message) bool {
		return m.ReceiptHandle == receiptHandle
	})
	if err != nil {
		return err
	}
	if len(result.Items) == 0 {
		return ErrInvalidReceiptHandle
	}
	return s.messagesStore.Delete(messageKey(queueURL, result.Items[0].Id))
}

// ChangeMessageVisibility changes the visibility timeout of a message.
func (s *SQSStore) ChangeMessageVisibility(queueURL, receiptHandle string, visibilityTimeout int32) error {
	if visibilityTimeout < 0 || visibilityTimeout > 43200 {
		return fmt.Errorf("InvalidParameterValue: VisibilityTimeout must be between 0 and 43200 seconds")
	}

	s.msgMutex.Lock()
	defer s.msgMutex.Unlock()

	if s.storage != nil {
		receiptsBucket := s.storage.Bucket("sqs-receipts-" + s.region)
		msgKeyBytes, err := receiptsBucket.Get([]byte(receiptHandle))
		if err != nil || len(msgKeyBytes) == 0 {
			return ErrInvalidReceiptHandle
		}

		var msgPb pb.Message
		if err := s.messagesStore.GetProto(string(msgKeyBytes), &msgPb); err != nil {
			return ErrInvalidReceiptHandle
		}

		msg := ProtoToMessage(&msgPb)
		msg.VisibilityTimeout = visibilityTimeout
		msg.ReceivedAt = time.Now()

		if visibilityTimeout == 0 {
			msg.ReceiptHandle = ""
			messagesBucket := "sqs-messages-" + s.region
			msgData, marshalErr := proto.Marshal(MessageToProto(msg))
			if marshalErr != nil {
				return marshalErr
			}
			return s.storage.Update(context.Background(), func(txn storage.Transaction) error {
				if err := txn.Bucket(messagesBucket).Put(msgKeyBytes, msgData); err != nil {
					return err
				}
				return txn.Bucket("sqs-receipts-" + s.region).Delete([]byte(receiptHandle))
			})
		}

		return s.messagesStore.PutProto(string(msgKeyBytes), MessageToProto(msg))
	}

	opts := common.ListOptions{Prefix: messagePrefix(queueURL), MaxItems: 1000}
	result, err := common.ListProto[*pb.Message](s.messagesStore, opts, func() *pb.Message { return &pb.Message{} }, func(m *pb.Message) bool {
		return m.ReceiptHandle == receiptHandle
	})
	if err != nil {
		return err
	}

	if len(result.Items) == 0 {
		return ErrInvalidReceiptHandle
	}

	msg := ProtoToMessage(result.Items[0])
	msg.VisibilityTimeout = visibilityTimeout
	msg.ReceivedAt = time.Now()

	if visibilityTimeout == 0 {
		msg.ReceiptHandle = ""
	}

	return s.messagesStore.PutProto(messageKey(queueURL, msg.ID), MessageToProto(msg))
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

	defer func() {
		s.purgeMutex.Lock()
		delete(s.purgeInProgress, queueURL)
		s.purgeMutex.Unlock()
	}()

	s.msgMutex.Lock()
	defer s.msgMutex.Unlock()

	err := s.messagesStore.DeleteByPrefix(messagePrefix(queueURL))
	if err != nil {
		return err
	}

	if s.storage != nil {
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
	}

	return nil
}
