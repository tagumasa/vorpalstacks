package sns

import (
	awserrors "vorpalstacks/internal/common/errors"
)

var (
	// ErrInvalidParameter is returned when a parameter is invalid.
	ErrInvalidParameter = awserrors.NewAWSError("InvalidParameter", "Invalid parameter", 400)
	// ErrNotFound is returned when the requested resource is not found.
	ErrNotFound = awserrors.NewAWSError("NotFound", "Resource not found", 404)
	// ErrAuthorizationError is returned when the request is not authorized.
	ErrAuthorizationError = awserrors.NewAWSError("AuthorizationError", "Authorization error", 403)
	// ErrInternalError is returned when an internal server error occurs.
	ErrInternalError = awserrors.NewAWSError("InternalError", "Internal server error", 500)
	// ErrEndpointDisabled is returned when the endpoint is disabled.
	ErrEndpointDisabled = awserrors.NewAWSError("EndpointDisabled", "Endpoint is disabled", 400)
	// ErrFilterLimitExceeded is returned when the filter limit is exceeded.
	ErrFilterLimitExceeded = awserrors.NewAWSError("FilterLimitExceeded", "Filter limit exceeded", 400)
	// ErrThrottled is returned when the request is throttled.
	ErrThrottled = awserrors.NewAWSError("Throttled", "Request was throttled", 429)
	// ErrValidation is returned when validation fails.
	ErrValidation = awserrors.NewAWSError("ValidationException", "Validation error", 400)
	// ErrTopicNotFound is returned when the topic does not exist.
	ErrTopicNotFound = awserrors.NewAWSError("NotFound", "Topic does not exist", 404)
	// ErrSubscriptionNotFound is returned when the subscription does not exist.
	ErrSubscriptionNotFound = awserrors.NewAWSError("NotFound", "Subscription does not exist", 404)
	// ErrPlatformAppNotFound is returned when the platform application does not exist.
	ErrPlatformAppNotFound = awserrors.NewAWSError("NotFound", "Platform application does not exist", 404)
	// ErrEndpointNotFound is returned when the endpoint does not exist.
	ErrEndpointNotFound = awserrors.NewAWSError("NotFound", "Endpoint does not exist", 404)
	// ErrTagLimitExceeded is returned when the tag limit is exceeded.
	ErrTagLimitExceeded = awserrors.NewAWSError("TagLimitExceeded", "Tag limit exceeded", 400)
	// ErrTagPolicy is returned when there is a tag policy violation.
	ErrTagPolicy = awserrors.NewAWSError("TagPolicy", "Tag policy violation", 400)
	// ErrBatchEntryIdsNotDistinct is returned when two or more batch entries have the same ID.
	ErrBatchEntryIdsNotDistinct = awserrors.NewAWSError("BatchEntryIdsNotDistinct", "Two or more batch entries have the same ID", 400)
	// ErrTooManyEntriesInBatch is returned when the batch request contains more than 10 entries.
	ErrTooManyEntriesInBatch = awserrors.NewAWSError("TooManyEntriesInBatchRequest", "Maximum number of entries per request are 10", 400)
)

// NewTopicNotFoundException creates a new AWSError for non-existent topic errors.
func NewTopicNotFoundException() *awserrors.AWSError {
	return awserrors.NewAWSError("NotFound", "Topic does not exist", 404)
}

// NewNotFoundException creates a new AWSError for not found errors.
func NewNotFoundException(resource string) *awserrors.AWSError {
	return awserrors.NewAWSError("NotFound", resource+" not found", 404)
}


