package sqs

import (
	"context"
	"errors"
	"strings"
	"vorpalstacks/internal/common/kmsutil"

	storecommon "vorpalstacks/internal/store/aws/common"
	sqsstore "vorpalstacks/internal/store/aws/sqs"
)

// ---------------------------------------------------------------------------
// Transport-agnostic DTOs shared by HTTP handlers and the admin gRPC handler.
// ---------------------------------------------------------------------------

// CreateQueueInput contains the parameters for creating an SQS queue.
type CreateQueueInput struct {
	QueueName string
	Region    string
	Attrs     map[string]string
	Tags      map[string]string
}

// CreateQueueResult contains the output of a create-queue operation.
type CreateQueueResult struct {
	QueueURL string
}

// DeleteQueueInput contains the parameters for deleting an SQS queue.
type DeleteQueueInput struct {
	QueueURL string
}

// ListQueuesInput contains the parameters for listing SQS queues.
type ListQueuesInput struct {
	QueueNamePrefix string
	MaxResults      int
	NextToken       string
}

// ListQueuesResult contains the output of a list-queues operation.
type ListQueuesResult struct {
	QueueURLs []string
	NextToken string
}

// GetQueueUrlInput contains the parameters for resolving a queue URL.
type GetQueueUrlInput struct {
	QueueName string
}

// GetQueueUrlResult contains the output of a get-queue-url operation.
type GetQueueUrlResult struct {
	QueueURL string
}

// ---------------------------------------------------------------------------
// Core methods — shared business logic used by both the HTTP API handlers
// and the admin gRPC-Web handler. Each method performs validation, delegates
// to the store, and maps store errors to service-level AWS errors.
// ---------------------------------------------------------------------------

// createQueueCore creates a new queue or returns the existing URL when the same
// name and attributes are submitted (idempotent re-creation). It rejects when
// the name is taken with different attributes.
func (s *SQSService) createQueueCore(ctx context.Context, store sqsstore.SQSStoreInterface, in CreateQueueInput) (*CreateQueueResult, error) {
	if in.QueueName == "" {
		return nil, ErrMissingParameter
	}
	if !isValidQueueName(in.QueueName) {
		return nil, ErrInvalidQueueName
	}

	queue := sqsstore.NewQueue(in.QueueName, in.Region, store.GetAccountID())

	if err := applyQueueAttributes(in.Attrs, queue); err != nil {
		return nil, err
	}

	// Validate FIFO queue naming bidirectionally: FifoQueue=true requires
	// ".fifo" suffix and vice-versa. Standard queue names cannot contain dots.
	if queue.FifoQueue && !strings.HasSuffix(in.QueueName, ".fifo") {
		return nil, ErrInvalidParameterValue
	}
	if !queue.FifoQueue && strings.HasSuffix(in.QueueName, ".fifo") {
		return nil, ErrInvalidParameterValue
	}

	// Validate KMS key existence when a KmsMasterKeyId is set.
	if kmsKey, ok := in.Attrs["KmsMasterKeyId"]; ok && kmsKey != "" && s.kmsChecker != nil {
		if err := s.kmsChecker.CheckKey(ctx, in.Region, kmsKey); err != nil {
			return nil, mapKMSError(err)
		}
	}

	if len(in.Tags) > 0 {
		queue.Tags = in.Tags
	}

	created, err := store.CreateQueue(queue)
	if err != nil {
		if errors.Is(err, sqsstore.ErrQueueAlreadyExists) {
			existingQueue, getErr := store.GetQueueByName(in.QueueName)
			if getErr != nil {
				return nil, convertStoreError(getErr)
			}
			if !requestAttrsMatchExisting(in.Attrs, existingQueue) {
				return nil, ErrQueueNameExists
			}
			return &CreateQueueResult{QueueURL: existingQueue.URL}, nil
		}
		return nil, convertStoreError(err)
	}

	return &CreateQueueResult{QueueURL: created.URL}, nil
}

// deleteQueueCore deletes a queue by URL.
func (s *SQSService) deleteQueueCore(store sqsstore.SQSStoreInterface, in DeleteQueueInput) error {
	if in.QueueURL == "" {
		return ErrMissingParameter
	}
	if err := store.DeleteQueue(in.QueueURL); err != nil {
		return convertStoreError(err)
	}
	return nil
}

// listQueuesCore lists queues with optional prefix filtering and pagination.
// The prefix filter runs inside the store listing so MaxResults counts only
// matching queues. NextToken is only surfaced when MaxResults was explicitly
// set ("You must set MaxResults to receive a value for NextToken").
func (s *SQSService) listQueuesCore(store sqsstore.SQSStoreInterface, in ListQueuesInput) (*ListQueuesResult, error) {
	if in.MaxResults < 0 || in.MaxResults > sqsstore.MaxListResults {
		return nil, ErrInvalidParameterValue
	}
	maxResultsSet := in.MaxResults > 0
	maxItems := in.MaxResults
	if !maxResultsSet {
		maxItems = sqsstore.MaxListResults
	}
	result, err := store.ListQueues(storecommon.ListOptions{
		MaxItems: maxItems,
		Marker:   in.NextToken,
	}, in.QueueNamePrefix)
	if err != nil {
		return nil, convertStoreError(err)
	}

	queueURLs := make([]string, 0, len(result.Items))
	for _, queue := range result.Items {
		queueURLs = append(queueURLs, queue.URL)
	}

	nextToken := ""
	if maxResultsSet && result.IsTruncated {
		nextToken = result.NextMarker
	}
	return &ListQueuesResult{
		QueueURLs: queueURLs,
		NextToken: nextToken,
	}, nil
}

// getQueueUrlCore resolves a queue name to its URL.
func (s *SQSService) getQueueUrlCore(store sqsstore.SQSStoreInterface, in GetQueueUrlInput) (*GetQueueUrlResult, error) {
	if in.QueueName == "" {
		return nil, ErrMissingParameter
	}
	queue, err := store.GetQueueByName(in.QueueName)
	if err != nil {
		return nil, convertStoreError(err)
	}
	return &GetQueueUrlResult{QueueURL: queue.URL}, nil
}

// mapKMSError converts a KMS key-checker error into the appropriate SQS Kms*
// error that matches the AWS SQS API error codes.
func mapKMSError(err error) error {
	switch {
	case errors.Is(err, kmsutil.ErrKeyNotFound):
		return ErrKmsNotFound
	case errors.Is(err, kmsutil.ErrKeyDisabled):
		return ErrKmsDisabled
	case errors.Is(err, kmsutil.ErrKeyInvalidState):
		return ErrKmsInvalidState
	case errors.Is(err, kmsutil.ErrKeyInvalidUsage):
		return ErrKmsInvalidKeyUsage
	default:
		return ErrInvalidParameterValue
	}
}
