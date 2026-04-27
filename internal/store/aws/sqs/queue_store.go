package sqs

import (
	"bytes"
	"fmt"
	"strconv"
	"time"

	pb "vorpalstacks/internal/pb/storage/storage_sqs"
	"vorpalstacks/internal/store/aws/common"
)

// CreateQueue creates a new SQS queue.
func (s *SQSStore) CreateQueue(queue *Queue) (*Queue, error) {
	if err := validateQueueName(queue.Name); err != nil {
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

	if err := s.BaseStore.PutProto(queueURL, QueueToProto(queue)); err != nil {
		return nil, err
	}

	if len(queue.Tags) > 0 {
		if err := s.TagStore.Tag(queueURL, queue.Tags); err != nil {
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
	queue.LastModifiedTimestamp = time.Now().UTC()
	queue.Attributes["LastModifiedTimestamp"] = fmt.Sprintf("%d", queue.LastModifiedTimestamp.Unix())
	return s.BaseStore.PutProto(queue.URL, QueueToProto(queue))
}

// DeleteQueue deletes a queue by its URL.
func (s *SQSStore) DeleteQueue(queueURL string) error {
	if !s.Exists(queueURL) {
		return ErrQueueNotFound
	}

	if err := s.messagesStore.DeleteByPrefix(messagePrefix(queueURL)); err != nil {
		return fmt.Errorf("failed to delete messages for queue %s: %w", queueURL, err)
	}

	if err := s.TagStore.Delete(queueURL); err != nil {
		return fmt.Errorf("failed to delete tags for queue %s: %w", queueURL, err)
	}

	s.deduplicationMu.Lock()
	s.cleanupDeduplicationCacheForQueue(queueURL)
	s.deduplicationMu.Unlock()

	if s.storage != nil {
		msgPrefix := messagePrefix(queueURL)
		prefixBytes := []byte(msgPrefix)

		receiptsBucket := s.storage.Bucket("sqs-receipts-" + s.region)
		var receiptKeys [][]byte
		_ = receiptsBucket.ForEach(func(k, v []byte) error {
			if bytes.HasPrefix(v, prefixBytes) {
				keyCopy := make([]byte, len(k))
				copy(keyCopy, k)
				receiptKeys = append(receiptKeys, keyCopy)
			}
			return nil
		})
		for _, k := range receiptKeys {
			_ = receiptsBucket.Delete(k)
		}

		dedupBucket := s.storage.Bucket("sqs-dedup-" + s.region)
		dedupPrefix := queueURL + "#"
		var dedupKeys [][]byte
		_ = dedupBucket.ForEach(func(k, v []byte) error {
			if bytes.HasPrefix(k, []byte(dedupPrefix)) {
				keyCopy := make([]byte, len(k))
				copy(keyCopy, k)
				dedupKeys = append(dedupKeys, keyCopy)
			}
			return nil
		})
		for _, k := range dedupKeys {
			_ = dedupBucket.Delete(k)
		}
	}

	return s.BaseStore.Delete(queueURL)
}

// ListQueues lists queues with the specified pagination options.
func (s *SQSStore) ListQueues(opts common.ListOptions) (*common.ListResult[Queue], error) {
	result, err := common.ListProto[*pb.Queue](s.BaseStore, opts, func() *pb.Queue { return &pb.Queue{} }, nil)
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

	if queue.Attributes == nil {
		queue.Attributes = make(map[string]string)
	}

	for k, v := range attributes {
		var valid bool
		switch k {
		case "VisibilityTimeout":
			if val, err := strconv.ParseInt(v, 10, 32); err == nil {
				if err := validateVisibilityTimeout(int32(val)); err != nil {
					return fmt.Errorf("validating visibility timeout: %w", err)
				}
				queue.VisibilityTimeout = int32(val)
				valid = true
			}
		case "MaximumMessageSize":
			if val, err := strconv.ParseInt(v, 10, 32); err == nil {
				if err := validateMaximumMessageSize(int32(val)); err != nil {
					return fmt.Errorf("validating maximum message size: %w", err)
				}
				queue.MaximumMessageSize = int32(val)
				valid = true
			}
		case "MessageRetentionPeriod":
			if val, err := strconv.ParseInt(v, 10, 32); err == nil {
				if err := validateMessageRetentionPeriod(int32(val)); err != nil {
					return fmt.Errorf("validating message retention period: %w", err)
				}
				queue.MessageRetentionPeriod = int32(val)
				valid = true
			}
		case "DelaySeconds":
			if val, err := strconv.ParseInt(v, 10, 32); err == nil {
				if err := validateDelaySeconds(int32(val)); err != nil {
					return fmt.Errorf("validating delay seconds: %w", err)
				}
				queue.DelaySeconds = int32(val)
				valid = true
			}
		case "ReceiveMessageWaitTimeSeconds":
			if val, err := strconv.ParseInt(v, 10, 32); err == nil {
				if err := validateReceiveMessageWaitTimeSeconds(int32(val)); err != nil {
					return fmt.Errorf("validating receive message wait time: %w", err)
				}
				queue.ReceiveMessageWaitTimeSeconds = int32(val)
				valid = true
			}
		case "RedrivePolicy":
			rdp, err := ParseRedrivePolicy(v)
			if err == nil {
				queue.RedrivePolicy = rdp
				valid = true
			}
		default:
			valid = true
		}
		if valid {
			queue.Attributes[k] = v
		}
	}

	return s.UpdateQueue(queue)
}

// AddPermission adds permission to a queue.
func (s *SQSStore) AddPermission(queueURL, label string, awsAccountIDs []string, actions []string) error {
	queue, err := s.GetQueue(queueURL)
	if err != nil {
		return fmt.Errorf("getting queue: %w", err)
	}

	if queue.Permissions == nil {
		queue.Permissions = make(map[string]*Permission)
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
func (s *SQSStore) ListQueueTags(queueURL string) (map[string]string, error) {
	return s.TagStore.List(queueURL)
}

// TagQueue adds tags to a queue.
func (s *SQSStore) TagQueue(queueURL string, tags map[string]string) error {
	if err := validateTags(tags); err != nil {
		return fmt.Errorf("validating tags: %w", err)
	}
	return s.TagStore.Tag(queueURL, tags)
}

// UntagQueue removes tags from a queue.
func (s *SQSStore) UntagQueue(queueURL string, tagKeys []string) error {
	return s.TagStore.Untag(queueURL, tagKeys)
}
