package sqs

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"

	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/core/storage"
	pb "vorpalstacks/internal/pb/storage/storage_sqs"
	"vorpalstacks/internal/store/aws/common"
)

const moveTaskBatchSize = 100

// StartMessageMoveTask starts a task to move messages from a source queue to a
// destination queue. A background goroutine performs the actual transfer.
func (s *SQSStore) StartMessageMoveTask(sourceARN, destARN string, maxMessages int32) (*MessageMoveTask, error) {
	// AWS SQS allows only one active message-move task per source queue.
	if s.hasActiveMoveTask(sourceARN) {
		return nil, ErrOverLimit
	}

	sourceURL := s.arnToQueueURL(sourceARN)
	if sourceURL == "" || !s.Exists(sourceURL) {
		return nil, ErrQueueNotFound
	}

	destURL := s.arnToQueueURL(destARN)
	if destURL == "" {
		destURL = s.findSourceQueueForDLQ(sourceURL)
		if destURL == "" {
			return nil, ErrQueueNotFound
		}
		if idx := strings.LastIndex(destURL, "/"); idx >= 0 {
			destARN = s.buildQueueARN(destURL[idx+1:])
		}
	}
	if destURL != sourceURL && !s.Exists(destURL) {
		return nil, ErrQueueNotFound
	}

	var totalToMove int32
	prefix := messagePrefix(sourceURL)
	_ = common.ForEachAllProto(s.messagesStore, prefix,
		func() *pb.Message { return &pb.Message{} }, nil,
		func(_ *pb.Message) error {
			totalToMove++
			return nil
		},
	)

	task := &MessageMoveTask{
		TaskId:                            uuid.New().String(),
		SourceQueueARN:                    sourceARN,
		DestinationQueueARN:               destARN,
		Status:                            "RUNNING",
		MaxNumberOfMessages:               maxMessages,
		StartTime:                         time.Now().UTC(),
		ApproximateNumberOfMessagesToMove: totalToMove,
	}

	if err := s.tasksStore.PutProto(task.TaskId, MessageMoveTaskToProto(task)); err != nil {
		return nil, err
	}

	// Clean up old terminal tasks for this source ARN to prevent unbounded
	// accumulation in PebbleDB.
	s.cleanupTerminalTasks(sourceARN)

	s.wg.Add(1)
	go s.runMessageMoveTask(task.TaskId, sourceURL, destURL, maxMessages)

	return task, nil
}

// runMessageMoveTask is the background worker that iterates source queue
// messages, copies them to the destination, and deletes from source.
func (s *SQSStore) runMessageMoveTask(taskID, sourceURL, destURL string, maxRate int32) {
	defer s.wg.Done()
	defer func() {
		if r := recover(); r != nil {
			logs.Error("SQS: panic in runMessageMoveTask goroutine",
				logs.String("taskId", taskID), logs.Any("panic", r))
			s.finalizeMoveTask(taskID, "FAILED", 0, 0, fmt.Sprintf("internal error: %v", r))
		}
	}()

	// Yield before the first batch so callers have a window to cancel or
	// list the task while it is still RUNNING. Without this, an empty source
	// queue causes the goroutine to finalise instantly, producing a race with
	// CancelMessageMoveTask and ListMessageMoveTasks.
	select {
	case <-s.ctx.Done():
		s.finalizeMoveTask(taskID, "FAILED", 0, 0, "server shutting down")
		return
	case <-time.After(100 * time.Millisecond):
	}

	var moved, failed int32

	for {
		if s.ctx.Err() != nil {
			s.finalizeMoveTask(taskID, "FAILED", moved, failed, "server shutting down")
			return
		}

		var taskPb pb.MessageMoveTask
		if err := s.tasksStore.GetProto(taskID, &taskPb); err == nil {
			if taskPb.Status == "CANCELLING" {
				s.finalizeMoveTask(taskID, "CANCELLED", moved, failed, "")
				return
			}
		}

		s.msgMutex.Lock()
		batchMoved, batchFailed, done, err := s.moveMessageBatch(sourceURL, destURL)
		s.msgMutex.Unlock()

		moved += batchMoved
		failed += batchFailed

		if err != nil {
			s.finalizeMoveTask(taskID, "FAILED", moved, failed, err.Error())
			return
		}

		s.updateMoveTaskProgress(taskID, moved, failed)

		if done {
			s.finalizeMoveTask(taskID, "COMPLETED", moved, failed, "")
			return
		}

		if maxRate > 0 && maxRate < moveTaskBatchSize {
			sleepDuration := time.Duration(moveTaskBatchSize) * time.Second / time.Duration(maxRate)
			select {
			case <-s.ctx.Done():
				s.finalizeMoveTask(taskID, "FAILED", moved, failed, "server shutting down")
				return
			case <-time.After(sleepDuration):
			}
		}
	}
}

// moveMessageBatch reads up to moveTaskBatchSize messages from sourceURL, copies
// each to destURL, and deletes the original. Returns moved/failed counts and
// whether the source queue is now empty.
//
// FIFO destination queues get sequence number assignment and deduplication
// checking applied to each moved message.
// Each message copy+delete is atomic via storage.Update transaction,
// preventing duplication on partial failure.
func (s *SQSStore) moveMessageBatch(sourceURL, destURL string) (moved, failed int32, done bool, err error) {
	opts := common.ListOptions{Prefix: messagePrefix(sourceURL), MaxItems: moveTaskBatchSize}
	result, listErr := common.ListProto(s.messagesStore, opts,
		func() *pb.Message { return &pb.Message{} }, nil)
	if listErr != nil {
		return 0, 0, false, listErr
	}

	if len(result.Items) == 0 {
		return 0, 0, true, nil
	}

	destQueue, qErr := s.GetQueue(destURL)
	if qErr != nil {
		return 0, 0, false, qErr
	}

	messagesBucket := "sqs-messages-" + s.region
	receiptsBucket := "sqs-receipts-" + s.region
	now := time.Now().UTC()

	for _, msgPb := range result.Items {
		newMsg := proto.Clone(msgPb).(*pb.Message)
		newMsg.Id = uuid.New().String()
		newMsg.SentTimestamp = timestamppb.New(now)
		newMsg.QueueUrl = destURL
		newMsg.QueueArn = destQueue.ARN
		newMsg.ReceiptHandle = ""
		newMsg.ReceivedAt = nil
		newMsg.VisibleAfter = nil
		newMsg.ApproximateReceiveCount = 0
		newMsg.ApproximateFirstReceiveTimestamp = nil

		if newMsg.Attributes == nil {
			newMsg.Attributes = map[string]string{}
		}
		newMsg.Attributes["SentTimestamp"] = fmt.Sprintf("%d", now.UnixMilli())
		newMsg.Attributes["SenderId"] = s.accountID

		// FIFO destination — assign sequence number and check deduplication
		srcKey := messageKey(sourceURL, msgPb.Id)
		var dedupKey string
		if destQueue.FifoQueue {
			s.sequenceMu.Lock()
			counter := s.sequenceCounters[destURL]
			if counter == 0 {
				counter = now.UnixNano()
			}
			counter++
			s.sequenceCounters[destURL] = counter
			newMsg.SequenceNumber = fmt.Sprintf("%d", counter)
			s.sequenceMu.Unlock()

			domainMsg := ProtoToMessage(newMsg)
			dedupKey = s.buildDeduplicationKey(destURL, domainMsg)
			if dedupKey != "" {
				if _, ok := s.getDeduplicationMessageID(dedupKey); ok {
					// Delete source message on dedup hit to prevent
					// infinite re-listing on subsequent batch iterations.
					handle := msgPb.ReceiptHandle
					_ = s.storage.Update(context.Background(), func(txn storage.Transaction) error {
						if err := txn.Bucket(messagesBucket).Delete([]byte(srcKey)); err != nil {
							return err
						}
						if handle != "" {
							_ = txn.Bucket(receiptsBucket).Delete([]byte(handle))
						}
						return nil
					})
					moved++
					continue
				}
			}
		}

		destKey := messageKey(destURL, newMsg.Id)
		data, marshalErr := proto.Marshal(newMsg)
		if marshalErr != nil {
			failed++
			continue
		}

		// Atomic copy + delete in a single transaction
		handle := msgPb.ReceiptHandle
		txErr := s.storage.Update(context.Background(), func(txn storage.Transaction) error {
			if err := txn.Bucket(messagesBucket).Put([]byte(destKey), data); err != nil {
				return err
			}
			if err := txn.Bucket(messagesBucket).Delete([]byte(srcKey)); err != nil {
				return err
			}
			if handle != "" {
				_ = txn.Bucket(receiptsBucket).Delete([]byte(handle))
			}
			return nil
		})
		if txErr != nil {
			failed++
			continue
		}

		if dedupKey != "" {
			s.putDeduplicationEntry(dedupKey, destKey)
		}

		moved++
	}

	return moved, failed, false, nil
}

func (s *SQSStore) updateMoveTaskProgress(taskID string, moved, failed int32) {
	s.taskMu.Lock()
	defer s.taskMu.Unlock()
	var taskPb pb.MessageMoveTask
	if err := s.tasksStore.GetProto(taskID, &taskPb); err != nil {
		return
	}
	// Don't stomp a terminal or cancelling status with stale progress.
	if taskPb.Status != "RUNNING" {
		return
	}
	taskPb.MovedMessages = moved
	taskPb.FailureMessages = failed
	_ = s.tasksStore.PutProto(taskID, &taskPb)
}

func (s *SQSStore) finalizeMoveTask(taskID, status string, moved, failed int32, failureReason string) {
	s.taskMu.Lock()
	defer s.taskMu.Unlock()
	var taskPb pb.MessageMoveTask
	if err := s.tasksStore.GetProto(taskID, &taskPb); err != nil {
		logs.Warn("SQS: failed to load move task for finalisation", logs.String("taskId", taskID), logs.Err(err))
		return
	}
	// Idempotent: already terminal, don't overwrite.
	currentStatus := taskPb.GetStatus()
	if currentStatus == "COMPLETED" || currentStatus == "FAILED" || currentStatus == "CANCELLED" {
		return
	}
	// If the task was marked CANCELLING while the goroutine was processing
	// its final batch, honour the cancel signal rather than COMPLETED/FAILED.
	if currentStatus == "CANCELLING" {
		status = "CANCELLED"
	}
	taskPb.Status = status
	taskPb.MovedMessages = moved
	taskPb.FailureMessages = failed
	taskPb.FailureReason = failureReason
	taskPb.EndTime = timestamppb.Now()
	if err := s.tasksStore.PutProto(taskID, &taskPb); err != nil {
		logs.Warn("SQS: failed to persist finalised move task", logs.String("taskId", taskID), logs.Err(err))
	}
}

// findSourceQueueForDLQ scans all queues for one whose RedrivePolicy points to
// the given DLQ URL's ARN. Used when DestinationArn is omitted.
func (s *SQSStore) findSourceQueueForDLQ(dlqURL string) string {
	dlqQueue, err := s.GetQueue(dlqURL)
	if err != nil {
		return ""
	}

	result, err := s.ListQueues(common.ListOptions{MaxItems: 1000})
	if err != nil {
		return ""
	}

	for _, q := range result.Items {
		if q.RedrivePolicy != nil && q.RedrivePolicy.DeadLetterTargetARN == dlqQueue.ARN {
			return q.URL
		}
	}
	return ""
}

// CancelMessageMoveTask marks a running task as CANCELLING. The background
// worker detects this and transitions to CANCELLED. Returns an error if the
// task has already reached a terminal state (COMPLETED, FAILED, CANCELLED).
func (s *SQSStore) CancelMessageMoveTask(taskId string) (*MessageMoveTask, error) {
	s.taskMu.Lock()
	defer s.taskMu.Unlock()
	var taskPb pb.MessageMoveTask
	if err := s.tasksStore.GetProto(taskId, &taskPb); err != nil {
		return nil, ErrTaskNotFound
	}
	task := ProtoToMessageMoveTask(&taskPb)

	switch task.Status {
	case "COMPLETED", "FAILED", "CANCELLED":
		return nil, ErrTaskAlreadyTerminal
	case "CANCELLING":
		return task, nil
	}

	task.Status = "CANCELLING"
	if err := s.tasksStore.PutProto(taskId, MessageMoveTaskToProto(task)); err != nil {
		return nil, err
	}
	return task, nil
}

// ListMessageMoveTasks lists message move tasks for a source queue, ordered by
// start time descending (most recent first). maxResults limits the count (AWS
// default 1, upper limit 10).
func (s *SQSStore) ListMessageMoveTasks(sourceARN string, maxResults int32) ([]*MessageMoveTask, error) {
	items, err := common.ListMatchingProto[*pb.MessageMoveTask](s.tasksStore, "",
		func() *pb.MessageMoveTask { return &pb.MessageMoveTask{} },
		func(t *pb.MessageMoveTask) bool {
			return t.SourceQueueArn == sourceARN
		})
	if err != nil {
		return nil, err
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].StartTime.AsTime().After(items[j].StartTime.AsTime())
	})

	if maxResults > 0 && int32(len(items)) > maxResults {
		items = items[:maxResults]
	}

	tasks := make([]*MessageMoveTask, len(items))
	for i, t := range items {
		tasks[i] = ProtoToMessageMoveTask(t)
	}
	return tasks, nil
}

// hasActiveMoveTask returns true when there is at least one RUNNING or CANCELLING
// task for the given source ARN.
func (s *SQSStore) hasActiveMoveTask(sourceARN string) bool {
	items, err := common.ListMatchingProto[*pb.MessageMoveTask](s.tasksStore, "",
		func() *pb.MessageMoveTask { return &pb.MessageMoveTask{} },
		func(t *pb.MessageMoveTask) bool {
			return t.SourceQueueArn == sourceARN &&
				(t.Status == "RUNNING" || t.Status == "CANCELLING")
		})
	if err != nil {
		return false
	}
	return len(items) > 0
}

// GetMessageMoveTask retrieves a specific message move task by its ID.
func (s *SQSStore) GetMessageMoveTask(taskId string) (*MessageMoveTask, error) {
	var taskPb pb.MessageMoveTask
	if err := s.tasksStore.GetProto(taskId, &taskPb); err != nil {
		return nil, ErrTaskNotFound
	}
	return ProtoToMessageMoveTask(&taskPb), nil
}

// cleanupTerminalTasks removes old terminal tasks for the given source ARN.
// Called from StartMessageMoveTask to prevent unbounded accumulation.
func (s *SQSStore) cleanupTerminalTasks(sourceARN string) {
	items, err := common.ListMatchingProto[*pb.MessageMoveTask](s.tasksStore, "",
		func() *pb.MessageMoveTask { return &pb.MessageMoveTask{} },
		func(t *pb.MessageMoveTask) bool {
			if t.SourceQueueArn != sourceARN {
				return false
			}
			return t.Status == "COMPLETED" || t.Status == "FAILED" || t.Status == "CANCELLED"
		})
	if err != nil {
		return
	}
	for _, t := range items {
		_ = s.tasksStore.Delete(t.TaskId)
	}
}
