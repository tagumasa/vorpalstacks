package sqs

import (
	awserrors "vorpalstacks/internal/services/aws/common/errors"
)

// SQSError represents an error returned by the SQS service.
type SQSError struct {
	*awserrors.AWSError
}

// Unwrap returns the underlying AWS error.
func (e *SQSError) Unwrap() error {
	return e.AWSError
}

// As implements the error interface's As method.
func (e *SQSError) As(target interface{}) bool {
	if t, ok := target.(**awserrors.AWSError); ok {
		*t = e.AWSError
		return true
	}
	return false
}

var (
	// ErrQueueDoesNotExist is returned when the specified queue does not exist.
	ErrQueueDoesNotExist = &SQSError{awserrors.NewAWSError("QueueDoesNotExist", "The specified queue does not exist.", 400).SetQueryErrorCode("AWS.SimpleQueueService.NonExistentQueue;Sender")}
	// ErrQueueNameExists is returned when a queue already exists with the same name and a different configuration.
	ErrQueueNameExists = &SQSError{awserrors.NewAWSError("QueueNameExists", "A queue already exists with the same name and a different configuration.", 400)}
	// ErrInvalidQueueName is returned when the queue name is invalid.
	ErrInvalidQueueName = &SQSError{awserrors.NewAWSError("InvalidParameterValue", "Invalid queue name.", 400)}
	// ErrInvalidParameterValue is returned when a parameter value is invalid.
	ErrInvalidParameterValue = &SQSError{awserrors.NewAWSError("InvalidParameterValue", "Invalid parameter value.", 400)}
	// ErrMissingParameter is returned when a required parameter is missing.
	ErrMissingParameter = &SQSError{awserrors.NewAWSError("MissingParameter", "The request must contain a required parameter.", 400)}
	// ErrInvalidBatchEntryId is returned when a batch entry ID is invalid.
	ErrInvalidBatchEntryId = &SQSError{awserrors.NewAWSError("InvalidBatchEntryId", "The batch entry ID is invalid.", 400)}
	// ErrBatchEntryIdsNotDistinct is returned when two or more batch entries have the same ID.
	ErrBatchEntryIdsNotDistinct = &SQSError{awserrors.NewAWSError("BatchEntryIdsNotDistinct", "Two or more batch entries have the same ID.", 400)}
	// ErrTooManyEntriesInBatch is returned when the batch request contains more than 10 entries.
	ErrTooManyEntriesInBatch = &SQSError{awserrors.NewAWSError("TooManyEntriesInBatchRequest", "Maximum number of entries per request are 10.", 400)}
	// ErrEmptyBatchRequest is returned when the batch request contains no entries.
	ErrEmptyBatchRequest = &SQSError{awserrors.NewAWSError("EmptyBatchRequest", "The batch request doesn't contain any entries.", 400)}
	// ErrReceiptHandleIsInvalid is returned when the receipt handle provided is not valid.
	ErrReceiptHandleIsInvalid = &SQSError{awserrors.NewAWSError("ReceiptHandleIsInvalid", "The receipt handle provided is not valid.", 400)}
	// ErrInvalidAttributeName is returned when the attribute name is invalid.
	ErrInvalidAttributeName = &SQSError{awserrors.NewAWSError("InvalidAttributeName", "The attribute name is invalid.", 400)}
	// ErrMessageNotInflight is returned when the message referred to is not in flight.
	ErrMessageNotInflight = &SQSError{awserrors.NewAWSError("MessageNotInflight", "The message referred to is not in flight.", 400)}
	// ErrPurgeQueueInProgress is returned when a purge queue operation is already running.
	ErrPurgeQueueInProgress = &SQSError{awserrors.NewAWSError("PurgeQueueInProgress", "There is already a purge queue operation running.", 403)}
	// ErrMessageTooLarge is returned when the message contents are invalid.
	ErrMessageTooLarge = &SQSError{awserrors.NewAWSError("InvalidMessageContents", "The message contains invalid message contents.", 400)}
	// ErrTooManyTags is returned when too many tags are provided for a queue.
	ErrTooManyTags = &SQSError{awserrors.NewAWSError("InvalidParameterValue", "Too many tags for queue.", 400)}
	// ErrInvalidTagKey is returned when a tag key is invalid.
	ErrInvalidTagKey = &SQSError{awserrors.NewAWSError("InvalidParameterValue", "Tag key is invalid.", 400)}
	// ErrInvalidTagValue is returned when a tag value is invalid.
	ErrInvalidTagValue = &SQSError{awserrors.NewAWSError("InvalidParameterValue", "Tag value is invalid.", 400)}
	// ErrMissingMessageGroupId is returned when the MessageGroupId parameter is required but missing.
	ErrMissingMessageGroupId = &SQSError{awserrors.NewAWSError("InvalidParameterValue", "The request must contain the parameter MessageGroupId.", 400)}
	// ErrMissingDeduplicationId is returned when the MessageDeduplicationId parameter is required but missing.
	ErrMissingDeduplicationId = &SQSError{awserrors.NewAWSError("InvalidParameterValue", "The request must contain the parameter MessageDeduplicationId.", 400)}
)
