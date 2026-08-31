package sqs

import (
	"fmt"

	awserrors "vorpalstacks/internal/common/errors"
	sqsstore "vorpalstacks/internal/store/aws/sqs"
)

var (
	// ErrQueueDoesNotExist is returned when the specified queue does not exist.
	ErrQueueDoesNotExist = awserrors.NewAWSError("QueueDoesNotExist", "The specified queue does not exist.", 400).SetQueryErrorCode("AWS.SimpleQueueService.NonExistentQueue")
	// ErrQueueDeletedRecently is returned when a queue is recreated with the
	// same name within 60 seconds of deletion: "You must wait 60 seconds
	// after deleting a queue before you can create another queue with the
	// same name." (AWS SQS API Reference, CreateQueue errors, HTTP 400)
	ErrQueueDeletedRecently = awserrors.NewAWSError("QueueDeletedRecently", "You must wait 60 seconds after deleting a queue before you can create another queue with the same name.", 400).SetQueryErrorCode("AWS.SimpleQueueService.QueueDeletedRecently")
	// ErrQueueNameExists is returned when a queue already exists with the same name and a different configuration.
	ErrQueueNameExists = awserrors.NewAWSError("QueueNameExists", "A queue already exists with the same name and a different configuration.", 400).SetQueryErrorCode("QueueAlreadyExists")
	// ErrInvalidQueueName is returned when the queue name is invalid.
	ErrInvalidQueueName = awserrors.NewAWSError("InvalidParameterValue", "Invalid queue name.", 400)
	// ErrInvalidParameterValue is returned when a parameter value is invalid.
	ErrInvalidParameterValue = awserrors.NewAWSError("InvalidParameterValue", "Invalid parameter value.", 400)
	// ErrSerializationException is returned when a typed request member
	// carries a value of the wrong wire type, e.g. a non-integer string for
	// an Integer member such as MaxNumberOfMessages. The SQS model serves
	// awsJson1_0 (query-compatible); the protocol rejects such payloads
	// during deserialisation before validation, and no service model
	// enumerates this protocol-level error on its operations.
	ErrSerializationException = awserrors.NewAWSError("SerializationException", "The request payload couldn't be deserialized.", 400)
	// ErrMissingParameter is returned when a required parameter is missing.
	ErrMissingParameter = awserrors.NewAWSError("MissingParameter", "The request must contain a required parameter.", 400)
	// ErrInvalidBatchEntryId is returned when a batch entry ID is invalid.
	ErrInvalidBatchEntryId = awserrors.NewAWSError("InvalidBatchEntryId", "The batch entry ID is invalid.", 400).SetQueryErrorCode("AWS.SimpleQueueService.InvalidBatchEntryId")
	// ErrBatchEntryIdsNotDistinct is returned when two or more batch entries have the same ID.
	ErrBatchEntryIdsNotDistinct = awserrors.NewAWSError("BatchEntryIdsNotDistinct", "Two or more batch entries have the same ID.", 400).SetQueryErrorCode("AWS.SimpleQueueService.BatchEntryIdsNotDistinct")
	// ErrTooManyEntriesInBatch is returned when the batch request exceeds the
	// entry limit.
	ErrTooManyEntriesInBatch = awserrors.NewAWSError("TooManyEntriesInBatchRequest", fmt.Sprintf("Maximum number of entries per request are %d.", sqsstore.MaxBatchEntries), 400).SetQueryErrorCode("AWS.SimpleQueueService.TooManyEntriesInBatchRequest")
	// ErrEmptyBatchRequest is returned when the batch request contains no entries.
	ErrEmptyBatchRequest = awserrors.NewAWSError("EmptyBatchRequest", "The batch request doesn't contain any entries.", 400).SetQueryErrorCode("AWS.SimpleQueueService.EmptyBatchRequest")
	// ErrReceiptHandleIsInvalid is returned when the receipt handle provided
	// is not valid. The service model defines this error with HTTP 404
	// (httpError and awsQueryError httpResponseCode both 404).
	ErrReceiptHandleIsInvalid = awserrors.NewAWSError("ReceiptHandleIsInvalid", "The receipt handle provided is not valid.", 404).SetQueryErrorCode("ReceiptHandleIsInvalid")
	// ErrInvalidAttributeName is returned when the attribute name is invalid.
	ErrInvalidAttributeName = awserrors.NewAWSError("InvalidAttributeName", "The attribute name is invalid.", 400)
	// ErrMessageNotInflight is returned when the message referred to is not in flight.
	ErrMessageNotInflight = awserrors.NewAWSError("MessageNotInflight", "The message referred to is not in flight.", 400).SetQueryErrorCode("AWS.SimpleQueueService.MessageNotInflight")
	// ErrPurgeQueueInProgress is returned when a purge queue operation is already running.
	ErrPurgeQueueInProgress = awserrors.NewAWSError("PurgeQueueInProgress", "There is already a purge queue operation running.", 403).SetQueryErrorCode("AWS.SimpleQueueService.PurgeQueueInProgress")
	// ErrMessageTooLarge is returned when the message body exceeds the queue's
	// MaximumMessageSize.
	ErrMessageTooLarge = awserrors.NewAWSError("InvalidMessageContents", "The message contains invalid message contents.", 400)
	// ErrInvalidMessageContents is returned when the message body contains
	// characters outside the allowed set: "The message contains characters
	// outside the allowed set." (AWS SQS API Reference, SendMessage errors,
	// HTTP 400)
	ErrInvalidMessageContents = awserrors.NewAWSError("InvalidMessageContents", "The message contains characters outside the allowed set.", 400)
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
	// ErrInvalidParameterCombination is returned when mutually exclusive parameters are both set.
	ErrInvalidParameterCombination = awserrors.NewAWSError("InvalidParameterCombination", "AttributeNames and MessageSystemAttributeNames are mutually exclusive.", 400)
	// ErrResourceNotFound is returned for message-move operations on an
	// unknown task: "One or more specified resources don't exist." The
	// service model defines this error with HTTP 404 (httpError and
	// awsQueryError httpResponseCode both 404).
	ErrResourceNotFound = awserrors.NewAWSError("ResourceNotFoundException", "One or more specified resources don't exist.", 404).SetQueryErrorCode("ResourceNotFoundException")
	// ErrOverLimit is returned when a resource limit is exceeded.
	ErrOverLimit = awserrors.NewAWSError("OverLimit", "The specified request exceeds the limit.", 403).SetQueryErrorCode("OverLimit")
	// ErrInvalidAttributeValue is returned when an attribute value is invalid or immutable.
	ErrInvalidAttributeValue = awserrors.NewAWSError("InvalidAttributeValue", "The attribute value is invalid.", 400)
	// ErrBatchRequestTooLong is returned when a batch request payload exceeds
	// the maximum allowed total payload size.
	ErrBatchRequestTooLong = awserrors.NewAWSError("BatchRequestTooLong", fmt.Sprintf("Batch requests cannot be longer than %d bytes.", sqsstore.MaxMaximumMessageSize), 400).SetQueryErrorCode("AWS.SimpleQueueService.BatchRequestTooLong")
	// KMS-related errors — granular mapping matching the AWS SQS API. The
	// query-protocol codes use the KMS.* prefix from the service model.
	ErrKmsAccessDenied    = awserrors.NewAWSError("KmsAccessDenied", "Request was denied due to KMS access denied.", 400).SetQueryErrorCode("KMS.AccessDeniedException")
	ErrKmsDisabled        = awserrors.NewAWSError("KmsDisabled", "The KMS key is disabled.", 400).SetQueryErrorCode("KMS.DisabledException")
	ErrKmsInvalidKeyUsage = awserrors.NewAWSError("KmsInvalidKeyUsage", "The KMS key usage is invalid.", 400).SetQueryErrorCode("KMS.InvalidKeyUsageException")
	ErrKmsInvalidState    = awserrors.NewAWSError("KmsInvalidState", "The KMS key is in an invalid state.", 400).SetQueryErrorCode("KMS.InvalidStateException")
	ErrKmsNotFound        = awserrors.NewAWSError("KmsNotFound", "The KMS key was not found.", 400).SetQueryErrorCode("KMS.NotFoundException")
	// ErrKmsOptInRequired is returned when the KMS key requires opt-in
	// before use. The service model defines this error with HTTP 403
	// (httpError and awsQueryError httpResponseCode both 403), unlike the
	// rest of the KMS family which the model defines with 400.
	ErrKmsOptInRequired = awserrors.NewAWSError("KmsOptInRequired", "KMS opt-in is required.", 403).SetQueryErrorCode("KMS.OptInRequired")
	ErrKmsThrottled     = awserrors.NewAWSError("KmsThrottled", "Request was throttled by KMS.", 400).SetQueryErrorCode("KMS.ThrottlingException")
	// ErrInvalidAddress is returned when a request carries an invalid
	// identifier. The service model defines this error with HTTP 404
	// (httpError and awsQueryError httpResponseCode both 404) for every
	// operation.
	ErrInvalidAddress = awserrors.NewAWSError("InvalidAddress", "The specified ID is invalid.", 404).SetQueryErrorCode("InvalidAddress")
	// ErrInvalidSecurity is returned when the request was not made over HTTPS
	// or did not use SigV4 for signing. The service model defines this error
	// with HTTP 403 (httpError and awsQueryError httpResponseCode both 403)
	// for every operation.
	ErrInvalidSecurity = awserrors.NewAWSError("InvalidSecurity", "The request was not made over HTTPS or did not use SigV4 for signing.", 403).SetQueryErrorCode("InvalidSecurity")
	// ErrRequestThrottled is returned when a request exceeds the permitted
	// request rate. The service model defines this error with HTTP 403
	// (httpError and awsQueryError httpResponseCode both 403) for every
	// operation. Per-account request-rate quotas are AWS-account-tied rate
	// limiting, which this platform does not enforce, so there is no return
	// path for it.
	ErrRequestThrottled = awserrors.NewAWSError("RequestThrottled", "The request was denied due to request throttling.", 403).SetQueryErrorCode("RequestThrottled")
	// ErrUnsupportedOperation is returned for request members the service
	// model marks as not implemented, such as message-attribute list values.
	ErrUnsupportedOperation = awserrors.NewAWSError("UnsupportedOperation", "Error code 400. Unsupported operation.", 400).SetQueryErrorCode("AWS.SimpleQueueService.UnsupportedOperation")
	// ErrInvalidIdFormat is a legacy DeleteMessage error the service model
	// marks as deprecated and included in ReceiptHandleIsInvalid, which is
	// what the platform returns for invalid receipt handles. The model gives
	// it no awsQueryError trait, so the query wire falls back to the shape
	// name.
	ErrInvalidIdFormat = awserrors.NewAWSError("InvalidIdFormat", "The specified receipt handle isn't valid for the current version.", 400)
)
