package sqs

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"time"

	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/core/storage"
	pb "vorpalstacks/internal/pb/storage/storage_sqs"
	"vorpalstacks/internal/store/aws/common"
	svcarn "vorpalstacks/internal/utils/aws/arn"
)

// CreateQueue creates a new SQS queue.
func (s *SQSStore) CreateQueue(queue *Queue) (*Queue, error) {
	if err := ValidateQueueName(queue.Name); err != nil {
		return nil, err
	}
	if err := ValidateFifoQueueName(queue.Name, queue.FifoQueue); err != nil {
		return nil, err
	}
	// Attribute names and value formats run through the same shared
	// validation as SetQueueAttributes, so the store's create path reaches
	// full coverage independently of the caller.
	if err := ValidateQueueAttributes(queue.Attributes); err != nil {
		return nil, err
	}

	if err := validateVisibilityTimeout(queue.VisibilityTimeout); err != nil {
		return nil, err
	}
	if err := validateMaximumMessageSize(queue.MaximumMessageSize); err != nil {
		return nil, err
	}
	if err := validateMessageRetentionPeriod(queue.MessageRetentionPeriod); err != nil {
		return nil, err
	}
	if err := validateDelaySeconds(queue.DelaySeconds); err != nil {
		return nil, err
	}
	if err := validateReceiveMessageWaitTimeSeconds(queue.ReceiveMessageWaitTimeSeconds); err != nil {
		return nil, err
	}
	if len(queue.Tags) > 0 {
		if err := validateTags(queue.Tags); err != nil {
			return nil, err
		}
	}

	queueURL := s.buildQueueURL(queue.Name)

	s.queueMutex.Lock()
	defer s.queueMutex.Unlock()

	if s.Exists(queueURL) {
		return nil, ErrQueueAlreadyExists
	}
	if s.deletedRecently(queueURL) {
		return nil, ErrQueueDeletedRecently
	}

	now := time.Now().UTC()
	queue.URL = queueURL
	queue.ARN = s.buildQueueARN(queue.Name)
	queue.Region = s.region
	queue.AccountID = s.accountID
	queue.CreatedTimestamp = now
	queue.LastModifiedTimestamp = now

	if queue.Attributes == nil {
		queue.Attributes = make(map[string]string)
	}

	queue.Attributes["QueueArn"] = queue.ARN
	queue.Attributes["CreatedTimestamp"] = fmt.Sprintf("%d", now.Unix())
	queue.Attributes["LastModifiedTimestamp"] = fmt.Sprintf("%d", now.Unix())
	queue.Attributes["VisibilityTimeout"] = fmt.Sprintf("%d", queue.VisibilityTimeout)
	queue.Attributes["MaximumMessageSize"] = fmt.Sprintf("%d", queue.MaximumMessageSize)
	queue.Attributes["MessageRetentionPeriod"] = fmt.Sprintf("%d", queue.MessageRetentionPeriod)
	queue.Attributes["DelaySeconds"] = fmt.Sprintf("%d", queue.DelaySeconds)
	queue.Attributes["ReceiveMessageWaitTimeSeconds"] = fmt.Sprintf("%d", queue.ReceiveMessageWaitTimeSeconds)
	queue.Attributes["FifoQueue"] = fmt.Sprintf("%v", queue.FifoQueue)
	queue.Attributes["ContentBasedDeduplication"] = fmt.Sprintf("%v", queue.ContentBasedDeduplication)

	if err := validateSSEExclusion(queue.Attributes); err != nil {
		return nil, err
	}
	if err := validateHighThroughputFifo(queue.Attributes); err != nil {
		return nil, err
	}
	if err := s.validateRedrivePolicyTarget(queue, queue.RedrivePolicy); err != nil {
		return nil, err
	}

	if err := s.BaseStore.PutProto(queueURL, QueueToProto(queue)); err != nil {
		return nil, err
	}

	if len(queue.Tags) > 0 {
		if err := s.TagStore.Tag(queueURL, queue.Tags); err != nil {
			// The create must either happen completely or not at all: the
			// record above already committed, so a failed tag write is
			// compensated by removing the record — otherwise the queue
			// exists tagless forever (the Core's idempotent-retry path
			// returns the existing URL without re-applying the requested
			// tags). No deletion marker is written: this queue was never
			// observable to a client, so an immediate retry must succeed
			// rather than hit the recreate-prohibition window.
			if delErr := s.BaseStore.Delete(queueURL); delErr != nil {
				logs.Error("SQS: queue record survived a failed tag write; the queue exists without its requested tags",
					logs.String("queueUrl", queueURL), logs.Err(delErr))
			}
			return nil, err
		}
	}

	return queue, nil
}

// GetQueue retrieves a queue by its URL.
func (s *SQSStore) GetQueue(queueURL string) (*Queue, error) {
	var p pb.Queue
	if err := s.BaseStore.GetProto(queueURL, &p); err != nil {
		return nil, ErrQueueNotFound
	}
	return ProtoToQueue(&p), nil
}

// GetQueueByName retrieves a queue by its name.
func (s *SQSStore) GetQueueByName(queueName string) (*Queue, error) {
	queueURL := s.buildQueueURL(queueName)
	return s.GetQueue(queueURL)
}

// UpdateQueue updates an existing queue.
func (s *SQSStore) UpdateQueue(queue *Queue) error {
	if !s.Exists(queue.URL) {
		return ErrQueueNotFound
	}
	// A queue without the raw attribute map would panic on the timestamp
	// write below; CreateQueue always seeds the map, so this guard only
	// covers callers constructing Queue values by hand.
	if queue.Attributes == nil {
		queue.Attributes = make(map[string]string)
	}
	queue.LastModifiedTimestamp = time.Now().UTC()
	queue.Attributes["LastModifiedTimestamp"] = fmt.Sprintf("%d", queue.LastModifiedTimestamp.Unix())
	return s.BaseStore.PutProto(queue.URL, QueueToProto(queue))
}

// DeleteQueue deletes a queue by its URL and removes the move-task records
// that reference it: a task without its source or destination queue can
// neither run nor be meaningfully listed. The task-record cleanup runs after
// the queue lock is released and takes only taskMu, so taskMu is never
// nested inside the queue lock plane.
func (s *SQSStore) DeleteQueue(queueURL string) error {
	queueARN, err := s.deleteQueueRecord(queueURL)
	if err != nil {
		return err
	}
	s.deleteMoveTasksReferencingQueue(queueARN)
	return nil
}

// deleteQueueRecord performs the queue deletion proper and returns the ARN of
// the deleted queue. The queue record, its messages, their receipt-handle
// entries and its deduplication keys are removed in ONE storage.Update
// transaction: a failure on any leg rolls the whole deletion back and the
// queue stays fully intact — messages included — instead of the previous
// unsequenced series that left a half-deleted queue (messages wiped, record
// and/or tags alive) on a later leg's failure. Deletion excludes both lock
// planes: queueMutex blocks every queue mutator for the whole delete, so a
// reader that already resolved the queue cannot resurrect the record through
// UpdateQueue after the delete. Lock order is queueMutex → msgMutex, the
// codebase's queue-then-messages layering.
func (s *SQSStore) deleteQueueRecord(queueURL string) (string, error) {
	s.queueMutex.Lock()
	defer s.queueMutex.Unlock()
	return s.deleteQueueRecordLocked(queueURL)
}

// deleteQueueRecordLocked is deleteQueueRecord's critical section; the
// caller holds queueMutex. A failed key-scan aborts before the transaction
// with the queue fully intact, never a partial delete.
func (s *SQSStore) deleteQueueRecordLocked(queueURL string) (string, error) {
	if !s.Exists(queueURL) {
		return "", ErrQueueNotFound
	}

	idx := strings.LastIndex(queueURL, "/")
	if idx < 0 || idx == len(queueURL)-1 {
		return "", nil
	}
	queueARN := s.buildQueueARN(queueURL[idx+1:])

	// The key sets are collected under both lock planes so the transaction
	// deletes a stable set; the mutex releases through the closure so a
	// panic inside the transaction cannot unwind past the unlock.
	msgPrefix := []byte(messagePrefix(queueURL))
	dedupPrefix := []byte(queueURL + "#")
	receiptsPrefix := msgPrefix
	var msgKeys, dedupKeys, receiptKeys [][]byte
	collect := func() error {
		messagesBucket := s.storage.Bucket(s.messagesBucketName)
		if err := messagesBucket.ForEach(func(k, _ []byte) error {
			if bytes.HasPrefix(k, msgPrefix) {
				msgKeys = append(msgKeys, append([]byte(nil), k...))
			}
			return nil
		}); err != nil {
			return fmt.Errorf("scanning messages: %w", err)
		}
		dedupBucket := s.storage.Bucket(s.dedupBucket)
		if err := dedupBucket.ForEach(func(k, _ []byte) error {
			if bytes.HasPrefix(k, dedupPrefix) {
				dedupKeys = append(dedupKeys, append([]byte(nil), k...))
			}
			return nil
		}); err != nil {
			return fmt.Errorf("scanning deduplication keys: %w", err)
		}
		receiptsBucket := s.storage.Bucket(s.receiptsBucket)
		if err := receiptsBucket.ForEach(func(k, v []byte) error {
			if bytes.HasPrefix(v, receiptsPrefix) {
				receiptKeys = append(receiptKeys, append([]byte(nil), k...))
			}
			return nil
		}); err != nil {
			return fmt.Errorf("scanning receipt handles: %w", err)
		}
		return nil
	}
	txErr := func() error {
		s.msgMutex.Lock()
		defer s.msgMutex.Unlock()
		// An incomplete key set must never reach the transaction: the commit
		// would delete the queue record while the messages a failed scan
		// missed survive under the deleted queue's prefix — the half-deleted
		// state this single-transaction design exists to prevent.
		if err := collect(); err != nil {
			return err
		}
		return s.storage.Update(context.Background(), func(txn storage.Transaction) error {
			messagesBucket := txn.Bucket(s.messagesBucketName)
			for _, k := range msgKeys {
				if err := messagesBucket.Delete(k); err != nil {
					return err
				}
			}
			receiptsBucket := txn.Bucket(s.receiptsBucket)
			for _, k := range receiptKeys {
				if err := receiptsBucket.Delete(k); err != nil {
					return err
				}
			}
			dedupBucket := txn.Bucket(s.dedupBucket)
			for _, k := range dedupKeys {
				if err := dedupBucket.Delete(k); err != nil {
					return err
				}
			}
			return txn.Bucket(auxiliaryBucketName("sqs-queues", s.region)).Delete([]byte(queueURL))
		})
	}()
	if txErr != nil {
		return "", txErr
	}

	// The in-memory mirrors of the deleted state are reclaimed after the
	// commit: sequence counters (otherwise only ever written), the
	// receive-attempt cache entries (otherwise bounded only by size and age)
	// and any lingering purge timestamp.
	s.deduplicationMu.Lock()
	s.cleanupDeduplicationCacheForQueue(queueURL)
	s.deduplicationMu.Unlock()

	s.sequenceMu.Lock()
	delete(s.sequenceCounters, queueURL)
	s.sequenceMu.Unlock()

	s.receiveAttemptMu.Lock()
	attemptPrefix := queueURL + "#"
	for k := range s.receiveAttemptCache {
		if strings.HasPrefix(k, attemptPrefix) {
			delete(s.receiveAttemptCache, k)
		}
	}
	s.receiveAttemptMu.Unlock()

	s.purgeMutex.Lock()
	delete(s.purgeInProgress, queueURL)
	s.purgeMutex.Unlock()

	// The tag record lives behind the tag framework's own API (a separate
	// bucket the store cannot address through its transaction), so it is
	// removed after the committed deletion. A failure here cannot fail the
	// delete — the queue is already gone — and leaves inert tag residue for
	// a URL no store-level tag operation serves, so it is logged loudly
	// instead of surfacing a "deletion failed" for a deletion that happened.
	if err := s.TagStore.Delete(queueURL); err != nil {
		logs.Error("SQS: queue deleted but its tag record survived (inert residue; no tag operation serves a nonexistent queue)",
			logs.String("queueUrl", queueURL), logs.Err(err))
	}
	s.recordQueueDeletion(queueURL)

	return queueARN, nil
}

// deleteMoveTasksReferencingQueue removes every move-task record whose source
// or destination ARN is the deleted queue's, so the tasks bucket holds no
// records for queues that no longer exist.
func (s *SQSStore) deleteMoveTasksReferencingQueue(queueARN string) {
	if queueARN == "" {
		return
	}
	s.taskMu.Lock()
	defer s.taskMu.Unlock()
	items, err := common.ListMatchingProto[*pb.MessageMoveTask](s.tasksStore, "",
		func() *pb.MessageMoveTask { return &pb.MessageMoveTask{} },
		func(t *pb.MessageMoveTask) bool {
			return t.SourceQueueArn == queueARN || t.DestinationQueueArn == queueARN
		})
	if err != nil {
		// Fail-open by design, never silently: the miss leaves task records
		// naming a queue that no longer exists (unobservable — listing is
		// by source ARN — but stored).
		logs.Warn("SQS: move-task residue listing failed; records for the deleted queue stay", logs.String("queueArn", queueARN), logs.Err(err))
		return
	}
	for _, t := range items {
		if err := s.tasksStore.Delete(t.TaskId); err != nil {
			// Fail-open by design, never silently: the surviving record names
			// a queue that no longer exists (unobservable — listing is by
			// source ARN) and stays stored.
			logs.Warn("SQS: move-task residue delete failed; the record stays stored", logs.String("taskId", t.TaskId), logs.Err(err))
		}
	}
}

// ListQueues lists queues with the specified pagination options. When
// queueNamePrefix is non-empty the filter is applied during iteration so
// MaxItems counts only matching queues, matching the AWS behaviour of
// filtering before pagination.
func (s *SQSStore) ListQueues(opts common.ListOptions, queueNamePrefix string) (*common.ListResult[Queue], error) {
	var filter func(*pb.Queue) bool
	if queueNamePrefix != "" {
		filter = func(q *pb.Queue) bool {
			name := q.Name
			if name == "" {
				parts := strings.Split(q.Url, "/")
				name = parts[len(parts)-1]
			}
			return strings.HasPrefix(name, queueNamePrefix)
		}
	}
	result, err := common.ListProto[*pb.Queue](s.BaseStore, opts, func() *pb.Queue { return &pb.Queue{} }, filter)
	if err != nil {
		return nil, err
	}

	queues := make([]*Queue, 0, len(result.Items))
	for _, pbQueue := range result.Items {
		queues = append(queues, ProtoToQueue(pbQueue))
	}

	return &common.ListResult[Queue]{
		Items:       queues,
		NextMarker:  result.NextMarker,
		IsTruncated: result.IsTruncated,
	}, nil
}

// SetQueueAttributes sets attributes for a queue.
func (s *SQSStore) SetQueueAttributes(queueURL string, attributes map[string]string) error {
	s.queueMutex.Lock()
	defer s.queueMutex.Unlock()

	queue, err := s.GetQueue(queueURL)
	if err != nil {
		return fmt.Errorf("getting queue: %w", err)
	}

	// An empty attribute map is a successful no-op, but the queue must still
	// resolve so that QueueDoesNotExist is returned for missing queues.
	if len(attributes) == 0 {
		return nil
	}

	// Names and value formats are validated through the single shared path
	// (the service-layer Core validates the same way before calling).
	if err := ValidateQueueAttributes(attributes); err != nil {
		return err
	}

	if queue.Attributes == nil {
		queue.Attributes = make(map[string]string)
	}

	// Coercion of the validated values onto the typed fields; every value
	// is kept in the raw attribute map.
	for k, v := range attributes {
		switch k {
		case "VisibilityTimeout":
			queue.VisibilityTimeout = ParseInt32Attr(v)
		case "MaximumMessageSize":
			queue.MaximumMessageSize = ParseInt32Attr(v)
		case "MessageRetentionPeriod":
			queue.MessageRetentionPeriod = ParseInt32Attr(v)
		case "DelaySeconds":
			queue.DelaySeconds = ParseInt32Attr(v)
		case "ReceiveMessageWaitTimeSeconds":
			queue.ReceiveMessageWaitTimeSeconds = ParseInt32Attr(v)
		case "Policy":
			queue.Policy = v
		case "RedrivePolicy":
			rdp, _ := ParseRedrivePolicy(v)
			queue.RedrivePolicy = rdp
			if rdp == nil {
				// The empty value clears the association: the raw attribute
				// is removed so the queue reports RedrivePolicy as unset.
				delete(queue.Attributes, k)
				continue
			}
		case "FifoQueue":
			// Queue type is immutable after creation ("You can't change the
			// queue type after you create it").
			if ParseBoolAttr(v) != queue.FifoQueue {
				return ErrInvalidAttributeValue
			}
		case "ContentBasedDeduplication":
			queue.ContentBasedDeduplication = ParseBoolAttr(v)
		}
		queue.Attributes[k] = v
	}

	// Cross-attribute rules are enforced on the merged attribute view so
	// they also catch the case where one option is already set on the queue
	// and the request sets the other.
	if err := validateSSEExclusion(queue.Attributes); err != nil {
		return err
	}
	if err := validateHighThroughputFifo(queue.Attributes); err != nil {
		return err
	}
	// The redrive target is validated only when the request names a
	// RedrivePolicy: an update to unrelated attributes must succeed even
	// while the configured target is missing or misbehaving — revalidating
	// the stored policy here would brick every attribute write on the queue
	// once its dead-letter target is deleted.
	if _, touchesRedrive := attributes["RedrivePolicy"]; touchesRedrive {
		if err := s.validateRedrivePolicyTarget(queue, queue.RedrivePolicy); err != nil {
			return err
		}
	}

	return s.UpdateQueue(queue)
}

// AddPermission adds permission to a queue.
func (s *SQSStore) AddPermission(queueURL, label string, awsAccountIDs []string, actions []string) error {
	if err := validatePermissionLabel(label); err != nil {
		return err
	}
	if err := validateAWSAccountIDs(awsAccountIDs); err != nil {
		return err
	}
	if err := validateSQSActionList(actions); err != nil {
		return err
	}

	s.queueMutex.Lock()
	defer s.queueMutex.Unlock()

	queue, err := s.GetQueue(queueURL)
	if err != nil {
		return fmt.Errorf("getting queue: %w", err)
	}

	if queue.Permissions == nil {
		queue.Permissions = make(map[string]*Permission)
	}

	if _, exists := queue.Permissions[label]; !exists {
		if len(queue.Permissions) >= maxPermissionLabels {
			return ErrOverLimit
		}
	}

	queue.Permissions[label] = &Permission{
		Label:         label,
		AWSAccountIDs: awsAccountIDs,
		Actions:       actions,
	}

	return s.UpdateQueue(queue)
}

// RemovePermission removes permission from a queue.
func (s *SQSStore) RemovePermission(queueURL, label string) error {
	s.queueMutex.Lock()
	defer s.queueMutex.Unlock()

	queue, err := s.GetQueue(queueURL)
	if err != nil {
		return fmt.Errorf("getting queue: %w", err)
	}

	if queue.Permissions != nil {
		delete(queue.Permissions, label)
	}

	return s.UpdateQueue(queue)
}

// ListQueueTags lists all tags for a queue.
// ListQueueTags lists the tags attached to a queue. The queue must exist:
// tag records for nonexistent or deleted queue URLs are not a writable
// surface at the store level.
func (s *SQSStore) ListQueueTags(queueURL string) (map[string]string, error) {
	if !s.Exists(queueURL) {
		return nil, ErrQueueNotFound
	}
	return s.TagStore.List(queueURL)
}

// TagQueue adds tags to a queue, rejecting unknown queue URLs.
func (s *SQSStore) TagQueue(queueURL string, tags map[string]string) error {
	if !s.Exists(queueURL) {
		return ErrQueueNotFound
	}
	if err := validateTags(tags); err != nil {
		return fmt.Errorf("validating tags: %w", err)
	}
	return s.TagStore.Tag(queueURL, tags)
}

// UntagQueue removes tags from a queue, rejecting unknown queue URLs.
func (s *SQSStore) UntagQueue(queueURL string, tagKeys []string) error {
	if !s.Exists(queueURL) {
		return ErrQueueNotFound
	}
	return s.TagStore.Untag(queueURL, tagKeys)
}

// validateRedrivePolicyTarget enforces the documented dead-letter-queue
// constraints on a queue's RedrivePolicy: the target ARN must be well formed
// and in the same account and Region, must resolve to an existing queue that
// is a distinct queue from the source, and "The dead-letter queue of a FIFO
// queue must also be a FIFO queue. Similarly, the dead-letter queue of a
// standard queue must also be a standard queue." (AWS SQS API Reference.)
// The target's own RedriveAllowPolicy governs whether the source may name it.
func (s *SQSStore) validateRedrivePolicyTarget(source *Queue, rdp *RedrivePolicy) error {
	if rdp == nil {
		return nil
	}
	if rdp.DeadLetterTargetARN == "" {
		return ErrInvalidAttributeValue
	}
	_, _, region, accountID, resource := svcarn.SplitARN(rdp.DeadLetterTargetARN)
	if region != s.region || accountID != s.accountID || resource == "" {
		return ErrInvalidAttributeValue
	}
	dlqURL := s.arnToQueueURL(rdp.DeadLetterTargetARN)
	// A redrive moves a message between two distinct queues. A self-targeted
	// policy would put and delete the same message key in one transaction,
	// destroying the message on every over-count receive.
	if dlqURL == source.URL {
		return ErrInvalidAttributeValue
	}
	dlq, err := s.GetQueue(dlqURL)
	if err != nil {
		return ErrInvalidAttributeValue
	}
	if dlq.FifoQueue != source.FifoQueue {
		return ErrInvalidAttributeValue
	}
	if !redriveAllowedByTarget(source, dlq) {
		return ErrInvalidAttributeValue
	}
	return nil
}

// redriveAllowedByTarget applies the destination queue's RedriveAllowPolicy
// to a source queue naming it as its dead-letter target: "allowAll –
// (Default) Any source queues in this AWS account in the same Region can
// specify this queue as the dead-letter queue. / denyAll – No source queues
// can specify this queue as the dead-letter queue. / byQueue – Only queues
// specified by the sourceQueueArns parameter can specify this queue as the
// dead-letter queue." (AWS SQS API Reference.) An unset policy behaves as the
// documented allowAll default.
func redriveAllowedByTarget(source, dlq *Queue) bool {
	policy := dlq.Attributes["RedriveAllowPolicy"]
	if policy == "" {
		return true
	}
	permission, sourceARNs := parseRedriveAllowPolicy(policy)
	switch permission {
	case "denyAll":
		return false
	case "byQueue":
		for _, arn := range sourceARNs {
			if arn == source.ARN {
				return true
			}
		}
		return false
	default:
		return true
	}
}

// ListDeadLetterSourceQueues returns the queues that have the specified dead
// letter queue as their target, honouring the pagination options carried by
// the ListDeadLetterSourceQueues API (MaxResults 1-1000, NextToken).
func (s *SQSStore) ListDeadLetterSourceQueues(dlqARN string, opts common.ListOptions) (*common.ListResult[Queue], error) {
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

	return &common.ListResult[Queue]{
		Items:       queues,
		NextMarker:  result.NextMarker,
		IsTruncated: result.IsTruncated,
	}, nil
}
