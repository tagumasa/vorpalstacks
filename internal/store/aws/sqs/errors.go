package sqs

import (
	"errors"
)

var (
	// ErrQueueAlreadyExists is returned when attempting to create a queue that
	// already exists.
	ErrQueueAlreadyExists = errors.New("queue already exists")

	// ErrQueueNotFound is returned when the specified queue does not exist.
	ErrQueueNotFound = errors.New("queue does not exist")

	// ErrInvalidQueueName is returned when the queue name does not meet
	// SQS naming requirements.
	ErrInvalidQueueName = errors.New("invalid queue name")

	// ErrMessageNotFound is returned when the specified message cannot be
	// found in the queue.
	ErrMessageNotFound = errors.New("message not found")

	// ErrInvalidReceiptHandle is returned when the receipt handle provided
	// is not valid or has expired.
	ErrInvalidReceiptHandle = errors.New("invalid receipt handle")

	// ErrBatchEntryIdsNotDistinct is returned when two or more batch entries
	// share the same entry ID.
	ErrBatchEntryIdsNotDistinct = errors.New("batch entry ids not distinct")

	// ErrTooManyEntriesInBatch is returned when a batch request contains
	// more than the maximum allowed entries (10).
	ErrTooManyEntriesInBatch = errors.New("too many entries in batch")

	// ErrEmptyBatchRequest is returned when a batch request contains no entries.
	ErrEmptyBatchRequest = errors.New("empty batch request")

	// ErrInvalidAttributeName is returned when the specified attribute name
	// is not valid for the queue.
	ErrInvalidAttributeName = errors.New("invalid attribute name")

	// ErrInvalidMessageContents is returned when the message contents are
	// not valid (e.g., too large or contains invalid characters).
	ErrInvalidMessageContents = errors.New("invalid message contents")

	// ErrUnsupportedOperation is returned when the requested operation
	// is not supported.
	ErrUnsupportedOperation = errors.New("unsupported operation")

	// ErrTaskNotFound is returned when the specified task does not exist.
	ErrTaskNotFound = errors.New("task not found")

	// ErrInvalidBatchEntryId is returned when a batch entry ID is invalid.
	ErrInvalidBatchEntryId = errors.New("invalid batch entry id")

	// ErrInvalidParameterValue is returned when a parameter value is not valid.
	ErrInvalidParameterValue = errors.New("invalid parameter value")

	// ErrMessageTooLarge is returned when the message exceeds the maximum
	// allowed size.
	ErrMessageTooLarge = errors.New("message too large")

	// ErrTooManyTags is returned when too many tags are attached to the queue.
	ErrTooManyTags = errors.New("too many tags")

	// ErrInvalidTagKey is returned when a tag key is invalid.
	ErrInvalidTagKey = errors.New("invalid tag key")

	// ErrInvalidTagValue is returned when a tag value is invalid.
	ErrInvalidTagValue = errors.New("invalid tag value")

	// ErrPurgeQueueInProgress is returned when a purge operation is already
	// in progress for the queue.
	ErrPurgeQueueInProgress = errors.New("purge queue in progress")

	// ErrMissingMessageGroupId is returned when the MessageGroupId is required
	// for a FIFO queue but was not provided.
	ErrMissingMessageGroupId = errors.New("missing message group id for fifo queue")

	// ErrMissingDeduplicationId is returned when the MessageDeduplicationId
	// is required for a FIFO queue but was not provided.
	ErrMissingDeduplicationId = errors.New("missing deduplication id for fifo queue")

	// ErrQueueNameExistsWithDiffConfig is returned when a queue with the same
	// name already exists but with different attributes.
	ErrQueueNameExistsWithDiffConfig = errors.New("queue name exists with different configuration")
)
