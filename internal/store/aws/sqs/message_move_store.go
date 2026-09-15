package sqs

import (
	"context"
	"errors"
	"fmt"
	"sort"
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

// moveTaskCancelPollInterval slices the rate-throttle window so the task
// record is re-read between slices: a cancel marked inside the window is
// observed within one slice instead of after the whole window (which is
// 100 s at the minimum rate of one message per second).
const moveTaskCancelPollInterval = 500 * time.Millisecond

// deadLetterQueueSourceArnAttr is the receive-side system attribute naming the
// queue a message was redriven from (set by moveToDLQ); the move task's
// omitted-destination mode routes each message back to that queue.
const deadLetterQueueSourceArnAttr = "DeadLetterQueueSourceArn"

// Message-move task statuses, persisted as strings in the task records.
// Every status transition and terminal-state check references these
// constants, never a bare literal.
const (
	MoveTaskStatusRunning    = "RUNNING"
	MoveTaskStatusCancelling = "CANCELLING"
	MoveTaskStatusCompleted  = "COMPLETED"
	MoveTaskStatusFailed     = "FAILED"
	MoveTaskStatusCancelled  = "CANCELLED"
)

// StartMessageMoveTask starts a task to move messages from a source queue to a
// destination queue. A background goroutine performs the actual transfer. A
// destination is either the request's explicit DestinationArn or, when that
// member is omitted, each message's own recorded redrive origin.
func (s *SQSStore) StartMessageMoveTask(sourceARN, destARN string, maxRate int32) (*MessageMoveTask, error) {
	sourceURL := s.arnToQueueURL(sourceARN)
	if sourceURL == "" {
		return nil, ErrInvalidAddress
	}
	if !s.Exists(sourceURL) {
		return nil, ErrQueueNotFound
	}

	// The SourceArn contract admits only dead-letter queues: "Currently, only
	// ARNs of dead-letter queues (DLQs) whose sources are other Amazon SQS
	// queues are accepted" (StartMessageMoveTask model). A queue that nothing
	// redrives into is not a move source.
	if s.findSourceQueueForDLQ(sourceURL) == "" {
		return nil, ErrUnsupportedOperation
	}

	destURL := s.arnToQueueURL(destARN)
	// A supplied DestinationArn that does not parse as a queue ARN is a
	// malformed request, not an omitted member: falling back to the redrive
	// homes here would silently replace the caller's chosen destination.
	if destARN != "" && destURL == "" {
		return nil, ErrInvalidAddress
	}
	if destURL != "" && destURL != sourceURL && !s.Exists(destURL) {
		return nil, ErrQueueNotFound
	}
	// destARN is stored as the caller supplied it: the task record reports a
	// destination only when the request named one ("If a DestinationArn has
	// not been specified in the StartMessageMoveTask request, this field
	// value will be NULL", model). An empty destURL is the per-message mode —
	// "the messages will be redriven back to their respective original
	// source queues" — resolved from each message's
	// DeadLetterQueueSourceArn at move time.

	var totalToMove int32
	prefix := messagePrefix(sourceURL)
	if err := common.ForEachAllProto(s.messagesStore, prefix,
		func() *pb.Message { return &pb.Message{} }, nil,
		func(_ *pb.Message) error {
			totalToMove++
			return nil
		},
	); err != nil {
		// The count is advisory ("Approximate"); a failed scan understates
		// it, and the move itself proceeds. Logged, never silent.
		logs.Warn("SQS: move-task source scan failed; ApproximateNumberOfMessagesToMove understated", logs.String("sourceUrl", sourceURL), logs.Err(err))
	}

	task := &MessageMoveTask{
		TaskId:                            uuid.New().String(),
		SourceQueueARN:                    sourceARN,
		DestinationQueueARN:               destARN,
		Status:                            MoveTaskStatusRunning,
		MaxNumberOfMessages:               maxRate,
		StartTime:                         time.Now().UTC(),
		ApproximateNumberOfMessagesToMove: totalToMove,
	}

	// The single-active rule and the RUNNING record form one check-then-act
	// sequence under taskMu: without the lock, two concurrent starts for the
	// same source both pass the active check and both worker goroutines run.
	s.taskMu.Lock()
	defer s.taskMu.Unlock()
	// A start racing Close's wg.Wait would be a positive-delta Add against a
	// Wait that may already observe zero (documented WaitGroup misuse), and
	// the worker would outlive the store into torn-down storage. The closing
	// flag is written under this same mutex, so any Add either happens
	// before Close's critical section (Wait counts it) or is refused here.
	if s.closing {
		return nil, errStoreClosing
	}
	if s.hasActiveMoveTask(sourceARN) {
		return nil, ErrOverLimit
	}

	if err := s.tasksStore.PutProto(task.TaskId, MessageMoveTaskToProto(task)); err != nil {
		return nil, err
	}

	// Clean up old terminal tasks for this source ARN to prevent unbounded
	// accumulation in PebbleDB.
	s.cleanupTerminalTasks(sourceARN)

	s.wg.Add(1)
	go s.runMessageMoveTask(task.TaskId, sourceURL, destURL, maxRate)

	return task, nil
}

// recoverInterruptedMoveTasks finalises move-task records that an unclean
// shutdown left in RUNNING or CANCELLING. Both states are only meaningful
// with a live worker goroutine: after a crash none exists, so a stale
// RUNNING record would hold the single-active slot forever (StartMessageMoveTask
// counts it) and a stale CANCELLING record could never reach a terminal
// state. RUNNING finalises as FAILED — the move did not complete; CANCELLING
// finalises as CANCELLED — the operator's cancel was the last recorded
// intent. Runs synchronously at construction, before any API call can
// observe the stale states.
func (s *SQSStore) recoverInterruptedMoveTasks() {
	s.taskMu.Lock()
	defer s.taskMu.Unlock()
	items, err := common.ListMatchingProto[*pb.MessageMoveTask](s.tasksStore, "",
		func() *pb.MessageMoveTask { return &pb.MessageMoveTask{} },
		func(t *pb.MessageMoveTask) bool {
			return t.Status == MoveTaskStatusRunning || t.Status == MoveTaskStatusCancelling
		})
	if err != nil {
		// Fail-open by design, never silently: stale RUNNING/CANCELLING
		// records survive, holding their sources' single-active slots until
		// the next restart re-runs this recovery.
		logs.Error("SQS: move-task recovery listing failed; interrupted tasks stay un-finalised until the next restart", logs.Err(err))
		return
	}
	now := timestamppb.Now()
	for _, t := range items {
		if t.Status == MoveTaskStatusCancelling {
			t.Status = MoveTaskStatusCancelled
		} else {
			t.Status = MoveTaskStatusFailed
			t.FailureReason = "interrupted by an unclean server shutdown"
		}
		t.EndTime = now
		if err := s.tasksStore.PutProto(t.TaskId, t); err != nil {
			// A surviving stale record holds the single-active slot for its
			// source until the next restart re-runs this recovery — loud,
			// never silent.
			logs.Error("SQS: failed to finalise recovered move task", logs.String("taskId", t.TaskId), logs.Err(err))
		}
	}
}

// runMessageMoveTask is the background worker that iterates source queue
// messages, copies them to the destination, and deletes from source.
func (s *SQSStore) runMessageMoveTask(taskID, sourceURL, destURL string, maxRate int32) {
	defer s.wg.Done()
	defer func() {
		if r := recover(); r != nil {
			logs.Error("SQS: panic in runMessageMoveTask goroutine",
				logs.String("taskId", taskID), logs.Any("panic", r))
			s.finalizeMoveTask(taskID, MoveTaskStatusFailed, 0, 0, fmt.Sprintf("internal error: %v", r))
		}
	}()

	// Yield before the first batch so callers have a window to cancel or
	// list the task while it is still RUNNING. Without this, an empty source
	// queue causes the goroutine to finalise instantly, producing a race with
	// CancelMessageMoveTask and ListMessageMoveTasks.
	select {
	case <-s.ctx.Done():
		s.finalizeMoveTask(taskID, MoveTaskStatusFailed, 0, 0, "server shutting down")
		return
	case <-time.After(100 * time.Millisecond):
	}

	var moved int32
	// failedSeen carries the distinct message IDs whose latest attempt
	// failed: a residual failure is re-listed (and re-attempted) on every
	// iteration, so accumulating per-batch counts would count the same
	// message once per retry. A later success removes its ID again.
	failedSeen := make(map[string]struct{})
	failed := func() int32 { return int32(len(failedSeen)) }

	for {
		if s.ctx.Err() != nil {
			s.finalizeMoveTask(taskID, MoveTaskStatusFailed, moved, failed(), "server shutting down")
			return
		}

		var taskPb pb.MessageMoveTask
		if err := s.tasksStore.GetProto(taskID, &taskPb); err != nil {
			if common.IsNotFound(err) {
				// The task record is gone — DeleteQueue removes the records
				// naming a deleted queue. The worker is then a zombie by
				// definition (unobservable, uncancellable, holding no
				// single-active slot), and its sourceURL may already name a
				// recreated different incarnation: it must stop rather than
				// keep relocating that queue's messages.
				logs.Warn("SQS: move task record disappeared, stopping worker", logs.String("taskId", taskID), logs.Err(err))
				return
			}
			// Any other read failure is transient storage trouble, not
			// disappearance: stopping here would strand the record RUNNING —
			// holding its source's single-active slot and ignoring cancels —
			// until the next restart's recovery sweep. The worker re-reads
			// after one poll interval, keeping the shutdown paths live.
			logs.Warn("SQS: move task record read failed; retrying after one poll interval", logs.String("taskId", taskID), logs.Err(err))
			select {
			case <-s.ctx.Done():
				s.finalizeMoveTask(taskID, MoveTaskStatusFailed, moved, failed(), "server shutting down")
				return
			case <-time.After(receivePollInterval):
			}
			continue
		}
		if taskPb.Status == MoveTaskStatusCancelling {
			s.finalizeMoveTask(taskID, MoveTaskStatusCancelled, moved, failed(), "")
			return
		}

		// The mutex releases through defer so a panic inside the batch (the
		// goroutine's own recover keeps the process alive) cannot unwind
		// past the unlock and poison the whole message plane.
		var batchMoved, batchFailed int32
		var done bool
		var batchErr error
		func() {
			s.msgMutex.Lock()
			defer s.msgMutex.Unlock()
			batchMoved, batchFailed, done, batchErr = s.moveMessageBatch(sourceURL, destURL, failedSeen)
		}()

		moved += batchMoved

		if batchErr != nil {
			s.finalizeMoveTask(taskID, MoveTaskStatusFailed, moved, failed(), batchErr.Error())
			return
		}

		// A batch that moved nothing but failed messages cannot make
		// progress: the failed messages are re-listed unchanged on every
		// iteration. The task finalises FAILED instead — the failed
		// messages stay in the source queue, and the terminal record
		// carries the counts and reason.
		if batchMoved == 0 && batchFailed > 0 {
			s.finalizeMoveTask(taskID, MoveTaskStatusFailed, moved, failed(),
				fmt.Sprintf("no message in the batch could be moved (%d failed)", failed()))
			return
		}

		s.updateMoveTaskProgress(taskID, moved, failed())

		if done {
			// A fully walked source with distinct failures still in it
			// never retries them — terminal status ends the retries — so
			// reporting COMPLETED would hide exactly the residual failures
			// the all-fail rule exists to surface. They stay in the source
			// queue as that rule keeps them, and the task reports FAILED
			// with the moved count it did achieve.
			if failed() > 0 {
				s.finalizeMoveTask(taskID, MoveTaskStatusFailed, moved, failed(),
					fmt.Sprintf("%d messages could not be moved and remain in the source queue", failed()))
				return
			}
			s.finalizeMoveTask(taskID, MoveTaskStatusCompleted, moved, failed(), "")
			return
		}

		if maxRate > 0 {
			// Throttle to the requested rate across the whole documented
			// 1-500 range: after each batch of up to moveTaskBatchSize
			// messages, wait the window that rate implies for the batch
			// (100 messages at 500/s = 200 ms). The batch size is not an
			// upper bound on the rate. The window is slept in slices with
			// the task record re-read between them so a cancel marked
			// inside the window is honoured within one slice: an unsliced
			// sleep would keep a CANCELLING task moving messages for the
			// whole window (100 s at rate 1) with its source's
			// single-active slot held. Read failures inside the window are
			// left to the loop head's own retry handling.
			window := time.Duration(moveTaskBatchSize) * time.Second / time.Duration(maxRate)
			deadline := time.Now().Add(window)
			for time.Now().Before(deadline) {
				slice := time.Until(deadline)
				if slice > moveTaskCancelPollInterval {
					slice = moveTaskCancelPollInterval
				}
				select {
				case <-s.ctx.Done():
					s.finalizeMoveTask(taskID, MoveTaskStatusFailed, moved, failed(), "server shutting down")
					return
				case <-time.After(slice):
				}
				var progress pb.MessageMoveTask
				if err := s.tasksStore.GetProto(taskID, &progress); err == nil && progress.Status == MoveTaskStatusCancelling {
					break // the loop head finalises CANCELLED
				}
			}
		} else if batchMoved == 0 && batchFailed == 0 {
			// The page held only in-flight messages: poll for their
			// visibility to expire instead of re-listing in a tight loop
			// (the receive path's poll interval).
			select {
			case <-s.ctx.Done():
				s.finalizeMoveTask(taskID, MoveTaskStatusFailed, moved, failed(), "server shutting down")
				return
			case <-time.After(receivePollInterval):
			}
		}
	}
}

// moveMessageBatch moves up to moveTaskBatchSize messages from sourceURL,
// copying each to its destination and deleting the original. An empty destURL
// is the omitted-DestinationArn mode: each message's destination is the queue
// its DeadLetterQueueSourceArn names ("the messages will be redriven back to
// their respective original source queues"). A message whose recorded origin
// is absent, unparseable, self-pointing or deleted has no home to return to
// and stays in the source, uncounted. Returns moved/failed counts and whether
// the task is finished. Messages currently in flight under a receive are
// skipped (see moveOneMessage); a walk that reaches the listing's end with
// such messages pending returns done=false, and the worker polls for their
// visibility to expire.
//
// The listing is WALKED past unroutable messages: they stay in the source
// forever, so a front page full of them never drains, and stopping at the
// first page would report COMPLETED while movable messages sit beyond them
// in key order. Pages advance by marker until the walk has moved
// moveTaskBatchSize messages or the listing is exhausted; done means the
// whole source was walked with nothing pending — an all-unroutable source
// still completes, unchanged.
//
// FIFO destination queues get sequence number assignment and deduplication
// window registration applied to each moved message (registration without
// consultation, the same rule the policy redrive follows — see
// moveOneMessage).
// Each message copy+delete is atomic via storage.Update transaction,
// preventing duplication on partial failure.
//
// failedSeen is the worker-owned set of distinct message IDs whose latest
// attempt failed; the batch adds and removes IDs in it (a later success
// retracts an earlier failure) so the task's FailureMessages counts
// messages, not attempts. Error returns carry the partial counts of the
// same call: messages already relocated by earlier pages stay counted.
func (s *SQSStore) moveMessageBatch(sourceURL, destURL string, failedSeen map[string]struct{}) (moved, failed int32, done bool, err error) {
	// An in-place move (destination == source) has nothing to move: the
	// copy+delete transaction would otherwise put and delete the same key
	// and destroy every message, and a per-message skip would re-list the
	// source forever. The task completes with nothing moved instead.
	if destURL != "" && destURL == sourceURL {
		return 0, 0, true, nil
	}

	// Resolved destination queues, cached per target URL: the explicit mode
	// holds one entry, the per-message mode at most one per distinct origin.
	destQueues := make(map[string]*Queue)
	if destURL != "" {
		destQueue, qErr := s.GetQueue(destURL)
		if qErr != nil {
			return 0, 0, false, qErr
		}
		destQueues[destURL] = destQueue
	}

	now := time.Now().UTC()
	var pending int
	exhausted := false
	marker := ""

	for {
		opts := common.ListOptions{Prefix: messagePrefix(sourceURL), MaxItems: moveTaskBatchSize}
		if marker != "" {
			opts.Marker = marker
		}
		result, listErr := common.ListProto(s.messagesStore, opts,
			func() *pb.Message { return &pb.Message{} }, nil)
		if listErr != nil {
			return moved, failed, false, listErr
		}

		if len(result.Items) == 0 {
			exhausted = true
			break
		}

		for _, msgPb := range result.Items {
			targetURL := destURL
			if targetURL == "" {
				targetURL = s.arnToQueueURL(msgPb.Attributes[deadLetterQueueSourceArnAttr])
				if targetURL == "" || targetURL == sourceURL {
					// No home to return to: the message stays in the source.
					continue
				}
				if _, resolved := destQueues[targetURL]; !resolved {
					destQueue, qErr := s.GetQueue(targetURL)
					if qErr != nil {
						// An origin queue deleted after the redrive is the same
						// no-home case as an absent origin: the message stays in
						// the DLQ rather than failing the whole task.
						if errors.Is(qErr, ErrQueueNotFound) {
							continue
						}
						return moved, failed, false, qErr
					}
					destQueues[targetURL] = destQueue
				}
			}
			switch s.moveOneMessage(sourceURL, targetURL, destQueues[targetURL], msgPb, now) {
			case moveOutcomeMoved:
				moved++
				delete(failedSeen, msgPb.Id)
			case moveOutcomeFailed:
				failed++
				failedSeen[msgPb.Id] = struct{}{}
			case moveOutcomeSkipped:
				// Left in the source under its receive.
				pending++
			}
		}

		if !result.IsTruncated || result.NextMarker == "" {
			exhausted = true
			break
		}
		marker = result.NextMarker
		if moved+failed >= moveTaskBatchSize {
			break
		}
	}

	// done only for a fully walked source with nothing in flight: pending
	// messages keep the task running so the worker polls for their
	// visibility to expire (or their consumer to delete them).
	done = exhausted && pending == 0
	return moved, failed, done, nil
}

// moveOutcome classifies one message within a move batch.
type moveOutcome int

const (
	moveOutcomeMoved moveOutcome = iota
	moveOutcomeFailed
	moveOutcomeSkipped
)

// moveOneMessage relocates one source message to the destination queue in a
// single copy+delete transaction.
func (s *SQSStore) moveOneMessage(sourceURL, destURL string, destQueue *Queue, msgPb *pb.Message, now time.Time) moveOutcome {
	// A message under a live receive is left to its consumer: relocating it
	// would invalidate the held receipt handle mid-processing (the consumer's
	// DeleteMessage then fails while the message resurfaces at the
	// destination with its delay lost). The move picks the message up once
	// its visibility timeout expires — or never, if the consumer deletes it.
	if s.isMessageInFlight(msgPb, now) {
		return moveOutcomeSkipped
	}

	srcKey := messageKey(sourceURL, msgPb.Id)

	// The move keeps the identity the policy redrive (moveToDLQ) keeps:
	// the same ID, and the documented enqueue-timestamp policy — "For
	// standard queues ... the enqueue timestamp is unchanged ... For FIFO
	// queues, the enqueue timestamp resets when the message is moved to a
	// dead-letter queue." (AWS SQS Developer Guide). The move task is a
	// redrive-home: its default destination is the original source queue.
	newMsg := proto.Clone(msgPb).(*pb.Message)
	sentTimestamp := now
	if !destQueue.FifoQueue && msgPb.SentTimestamp != nil {
		sentTimestamp = msgPb.SentTimestamp.AsTime()
	}
	newMsg.SentTimestamp = timestamppb.New(sentTimestamp)
	newMsg.QueueUrl = destURL
	newMsg.QueueArn = destQueue.ARN
	newMsg.ReceiptHandle = ""
	newMsg.ReceivedAt = nil
	// A message scheduled for future visibility keeps its remaining delay:
	// the clone already carries VisibleAfter and it is not cleared. The
	// receive count accumulates across queues ("Returns the number of times
	// a message has been received across all queues but not deleted",
	// ReceiveMessage API reference); only the first-receive timestamp
	// restarts, its doc wording being queue-scoped.
	newMsg.ApproximateReceiveCount = msgPb.ApproximateReceiveCount
	newMsg.ApproximateFirstReceiveTimestamp = nil

	if newMsg.Attributes == nil {
		newMsg.Attributes = map[string]string{}
	}
	newMsg.Attributes["SentTimestamp"] = fmt.Sprintf("%d", sentTimestamp.UnixMilli())
	newMsg.Attributes["SenderId"] = s.accountID
	newMsg.Attributes["ApproximateReceiveCount"] = fmt.Sprintf("%d", newMsg.ApproximateReceiveCount)
	// The carried map keeps only message-owned state: the previous hop's
	// DeadLetterQueueSourceArn names the redrive relation of the queue the
	// copy just left (on a move home it would name the destination itself,
	// a self-reference the routing treats as no-home), and a carried
	// SqsManagedSseEnabled would report the source queue's encryption for a
	// message now in another queue ("Only one server-side encryption option
	// is supported per queue", ReceiveMessage API reference). moveToDLQ
	// re-stamps the marker on any future redrive, and the receive path
	// re-derives the SSE entry from the destination's own setting.
	delete(newMsg.Attributes, deadLetterQueueSourceArnAttr)
	delete(newMsg.Attributes, "SqsManagedSseEnabled")

	// Deduplication at a FIFO destination follows SendMessage's gating: only
	// a destination that opted in (ContentBasedDeduplication) or a message
	// carrying an explicit deduplication ID participates. Unconditional
	// body-hash deduplication silently deleted same-body messages into FIFO
	// queues that never enabled it.
	var dedupKey string
	if destQueue.FifoQueue && (destQueue.ContentBasedDeduplication || newMsg.MessageDeduplicationId != "") {
		dedupKey = s.buildDeduplicationKey(destURL, ProtoToMessage(newMsg))
	}
	if dedupKey != "" {
		// The same per-key mutual exclusion as SendMessage holds the copy
		// and the registration as one sequence under the key lock, so a
		// concurrent same-key send cannot interleave between them. The move
		// itself, like the policy redrive (moveToDLQ), registers WITHOUT
		// consulting the window: a move relocates an already-accepted
		// message, and suppressing it on an existing entry — the original
		// send's own window when moving a message back home, or an entry a
		// sibling move just opened — would delete the source copy and
		// deliver nothing, a loss no source licenses. Only the send paths
		// suppress.
		s.dedupKeyLocker.Lock(dedupKey)
		defer s.dedupKeyLocker.Unlock(dedupKey)
	}

	// Sequence numbers are assigned at the actual enqueue: on a FIFO
	// destination the enqueue timestamp is re-stamped atomically with the
	// sequence number (see nextFIFOStamp); on a standard destination the
	// copy's FIFO-only members are cleared — "This parameter applies only
	// to FIFO (first-in-first-out) queues" (model, on SequenceNumber and
	// MessageDeduplicationId both) — from both the typed field and the
	// attribute map an earlier receive may have populated. MessageGroupId
	// stays: the model documents it for standard queues too (fair queues).
	if destQueue.FifoQueue {
		stampTime, seq := s.nextFIFOStamp(destURL)
		newMsg.SentTimestamp = timestamppb.New(stampTime)
		newMsg.SequenceNumber = seq
		newMsg.Attributes["SentTimestamp"] = fmt.Sprintf("%d", stampTime.UnixMilli())
	} else {
		newMsg.SequenceNumber = ""
		delete(newMsg.Attributes, "SequenceNumber")
		newMsg.MessageDeduplicationId = ""
		delete(newMsg.Attributes, "MessageDeduplicationId")
	}

	destKey := messageKey(destURL, newMsg.Id)
	data, marshalErr := proto.Marshal(newMsg)
	if marshalErr != nil {
		return moveOutcomeFailed
	}

	// Atomic copy + delete in a single transaction; a FIFO destination's
	// deduplication entry commits with them, so a storage failure can never
	// move a message whose deduplication window is silently missing.
	handle := msgPb.ReceiptHandle
	dedupValue := deduplicationEntryValue(destKey, newMsg.SequenceNumber)
	var staleReceiptErr error
	txErr := s.storage.Update(context.Background(), func(txn storage.Transaction) error {
		if err := txn.Bucket(s.messagesBucketName).Put([]byte(destKey), data); err != nil {
			return err
		}
		if err := txn.Bucket(s.messagesBucketName).Delete([]byte(srcKey)); err != nil {
			return err
		}
		if handle != "" {
			// Inert on failure and never worth aborting the relocation: the
			// stale handle resolves to a deleted key and the 12 h age sweep
			// reclaims it. Logged after the commit, never silent.
			staleReceiptErr = txn.Bucket(s.receiptsBucket).Delete([]byte(handle))
		}
		if dedupKey != "" {
			return txn.Bucket(s.dedupBucket).Put([]byte(dedupKey), dedupValue)
		}
		return nil
	})
	if txErr != nil {
		return moveOutcomeFailed
	}
	if staleReceiptErr != nil {
		logs.Warn("SQS: move left a stale receipt entry behind; the age sweep reclaims it",
			logs.String("sourceUrl", sourceURL), logs.Err(staleReceiptErr))
	}

	if dedupKey != "" {
		s.registerDeduplicationEntry(dedupKey, destKey, newMsg.SequenceNumber)
	}

	return moveOutcomeMoved
}

func (s *SQSStore) updateMoveTaskProgress(taskID string, moved, failed int32) {
	s.taskMu.Lock()
	defer s.taskMu.Unlock()
	var taskPb pb.MessageMoveTask
	if err := s.tasksStore.GetProto(taskID, &taskPb); err != nil {
		// Fail-open by design, never silently: the missed update leaves the
		// record's counts stale until the next batch (or finalisation).
		logs.Warn("SQS: move-task progress read failed; counts stay stale until the next batch", logs.String("taskId", taskID), logs.Err(err))
		return
	}
	// Don't stomp a terminal or cancelling status with stale progress.
	if taskPb.Status != MoveTaskStatusRunning {
		return
	}
	taskPb.MovedMessages = moved
	taskPb.FailureMessages = failed
	if err := s.tasksStore.PutProto(taskID, &taskPb); err != nil {
		// Progress counts are advisory until finalisation; a failed write
		// leaves the record's counts stale until the next batch updates
		// them. Logged, never silent.
		logs.Warn("SQS: move-task progress write failed; counts stay stale until the next batch", logs.String("taskId", taskID), logs.Err(err))
	}
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
	if currentStatus == MoveTaskStatusCompleted || currentStatus == MoveTaskStatusFailed || currentStatus == MoveTaskStatusCancelled {
		return
	}
	// If the task was marked CANCELLING while the goroutine was processing
	// its final batch, honour the cancel signal rather than COMPLETED/FAILED.
	if currentStatus == MoveTaskStatusCancelling {
		status = MoveTaskStatusCancelled
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
// the given DLQ URL's ARN. Every start calls it as the SourceArn contract's
// gate (only a DLQ some live queue redrives into is a move source), and the
// omitted-destination mode's per-message routing relies on the same
// redrive relation it confirms.
func (s *SQSStore) findSourceQueueForDLQ(dlqURL string) string {
	dlqQueue, err := s.GetQueue(dlqURL)
	if err != nil {
		return ""
	}

	// Every page is examined: a single-page list misses the redrive source
	// once the queue count exceeds one page, and both the omitted-destination
	// resolution and the DLQ-only source rule depend on the answer.
	const pageSize = 100
	var marker string
	for {
		opts := common.ListOptions{MaxItems: pageSize}
		if marker != "" {
			opts.Marker = marker
		}
		result, err := s.ListQueues(opts, "")
		if err != nil {
			return ""
		}

		for _, q := range result.Items {
			if q.RedrivePolicy != nil && q.RedrivePolicy.DeadLetterTargetARN == dlqQueue.ARN {
				return q.URL
			}
		}
		if !result.IsTruncated || result.NextMarker == "" {
			return ""
		}
		marker = result.NextMarker
	}
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
	case MoveTaskStatusCompleted, MoveTaskStatusFailed, MoveTaskStatusCancelled:
		return nil, ErrTaskAlreadyTerminal
	case MoveTaskStatusCancelling:
		return task, nil
	}

	task.Status = MoveTaskStatusCancelling
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
				(t.Status == MoveTaskStatusRunning || t.Status == MoveTaskStatusCancelling)
		})
	if err != nil {
		// Fail-open by design, never silently: the miss admits a second
		// concurrent task for the source.
		logs.Warn("SQS: move-task listing failed; the single-active rule is unenforced for this start", logs.String("sourceArn", sourceARN), logs.Err(err))
		return false
	}
	return len(items) > 0
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
			return t.Status == MoveTaskStatusCompleted || t.Status == MoveTaskStatusFailed || t.Status == MoveTaskStatusCancelled
		})
	if err != nil {
		// Fail-open by design, never silently: the miss leaves terminal
		// records that the next start's cleanup re-attempts.
		logs.Warn("SQS: terminal move-task listing failed; old records stay until the next start", logs.String("sourceArn", sourceARN), logs.Err(err))
		return
	}
	for _, t := range items {
		if err := s.tasksStore.Delete(t.TaskId); err != nil {
			// Fail-open by design, never silently: a surviving terminal
			// record is reclaimed by the next start's cleanup pass.
			logs.Warn("SQS: terminal move-task record delete failed; the next start re-attempts it", logs.String("taskId", t.TaskId), logs.Err(err))
		}
	}
}
