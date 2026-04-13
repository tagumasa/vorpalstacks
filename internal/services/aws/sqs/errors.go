package sqs

import (
	awserrors "vorpalstacks/internal/common/errors"
)

var (
	// ErrQueueDoesNotExist is returned when the specified queue does not exist.
	ErrQueueDoesNotExist = awserrors.NewAWSError("QueueDoesNotExist", "The specified queue does not exist.", 400).SetQueryErrorCode("AWS.SimpleQueueService.NonExistentQueue;Sender")
	// ErrQueueNameExists is returned when a queue already exists with the same name and a different configuration.
	ErrQueueNameExists = awserrors.NewAWSError("QueueNameExists", "A queue already exists with the same name and a different configuration.", 400)
	// ErrInvalidQueueName is returned when the queue name is invalid.
	ErrInvalidQueueName = awserrors.NewAWSError("InvalidParameterValue", "Invalid queue name.", 400)
	// ErrInvalidParameterValue is returned when a parameter value is invalid.
	ErrInvalidParameterValue = awserrors.NewAWSError("InvalidParameterValue", "Invalid parameter value.", 400)
	// ErrMissingParameter is returned when a required parameter is missing.
	ErrMissingParameter = awserrors.NewAWSError("MissingParameter", "The request must contain a required parameter.", 400)
	// ErrInvalidBatchEntryId is returned when a batch entry ID is invalid.
	ErrInvalidBatchEntryId = awserrors.NewAWSError("InvalidBatchEntryId", "The batch entry ID is invalid.", 400)
	// ErrBatchEntryIdsNotDistinct is returned when two or more batch entries have the same ID.
	ErrBatchEntryIdsNotDistinct = awserrors.NewAWSError("BatchEntryIdsNotDistinct", "Two or more batch entries have the same ID.", 400)
	// ErrTooManyEntriesInBatch is returned when the batch request contains more than 10 entries.
	ErrTooManyEntriesInBatch = awserrors.NewAWSError("TooManyEntriesInBatchRequest", "Maximum number of entries per request are 10.", 400)
	// ErrEmptyBatchRequest is returned when the batch request contains no entries.
	ErrEmptyBatchRequest = awserrors.NewAWSError("EmptyBatchRequest", "The batch request doesn't contain any entries.", 400)
	// ErrReceiptHandleIsInvalid is returned when the receipt handle provided is not valid.
	ErrReceiptHandleIsInvalid = awserrors.NewAWSError("ReceiptHandleIsInvalid", "The receipt handle provided is not valid.", 400)
	// ErrInvalidAttributeName is returned when the attribute name is invalid.
	ErrInvalidAttributeName = awserrors.NewAWSError("InvalidAttributeName", "The attribute name is invalid.", 400)
	// ErrMessageNotInflight is returned when the message referred to is not in flight.
	ErrMessageNotInflight = awserrors.NewAWSError("MessageNotInflight", "The message referred to is not in flight.", 400)
	// ErrPurgeQueueInProgress is returned when a purge queue operation is already running.
	ErrPurgeQueueInProgress = awserrors.NewAWSError("PurgeQueueInProgress", "There is already a purge queue operation running.", 403)
	// ErrMessageTooLarge is returned when the message contents are invalid.
	ErrMessageTooLarge = awserrors.NewAWSError("InvalidMessageContents", "The message contains invalid message contents.", 400)
	// ErrTooManyTags is returned when too many tags are provided for a queue.
	ErrTooManyTags = awserrors.NewAWSError("InvalidParameterValue", "Too many tags for queue.", 400)
	// ErrInvalidTagKey is returned when a tag key is invalid.
	ErrInvalidTagKey = awserrors.NewAWSError("InvalidParameterValue", "Tag key is invalid.", 400)
	// ErrInvalidTagValue is returned when a tag value is invalid.
	ErrInvalidTagValue = awserrors.NewAWSError("InvalidParameterValue", "Tag value is invalid.", 400)
	// ErrMissingMessageGroupId is returned when the MessageGroupId parameter is required but missing.
	ErrMissingMessageGroupId = awserrors.NewAWSError("InvalidParameterValue", "The request must contain the parameter MessageGroupId.", 400)
	// ErrMissingDeduplicationId is returned when the MessageDeduplicationId parameter is required but missing.
	ErrMissingDeduplicationId = awserrors.NewAWSError("InvalidParameterValue", "The request must contain the parameter MessageDeduplicationId.", 400)
)
