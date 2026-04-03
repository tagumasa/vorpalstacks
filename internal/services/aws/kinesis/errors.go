package kinesis

import (
	"net/http"

	awserrors "vorpalstacks/internal/services/aws/common/errors"
)

var (
	// ErrResourceNotFound is returned when the specified resource does not exist.
	ErrResourceNotFound = awserrors.NewAWSError("ResourceNotFoundException", "Requested resource not found.", http.StatusNotFound)
	// ErrResourceInUse is returned when the resource is in use.
	ErrResourceInUse = awserrors.NewAWSError("ResourceInUseException", "The resource is in use.", http.StatusBadRequest)
	// ErrInvalidArgument is returned when an argument is invalid.
	ErrInvalidArgument = awserrors.NewAWSError("InvalidArgumentException", "Invalid argument.", http.StatusBadRequest)
	// ErrLimitExceeded is returned when the rate limit for the stream is exceeded.
	ErrLimitExceeded = awserrors.NewAWSError("LimitExceededException", "Rate exceeded for this stream.", http.StatusTooManyRequests)
	// ErrProvisionedThroughputExceeded is returned when the provisioned throughput is exceeded.
	ErrProvisionedThroughputExceeded = awserrors.NewAWSError("ProvisionedThroughputExceededException", "Rate exceeded for this shard.", http.StatusTooManyRequests)
	// ErrExpiredIterator is returned when the iterator has expired.
	ErrExpiredIterator = awserrors.NewAWSError("ExpiredIteratorException", "Iterator expired.", http.StatusBadRequest)
	// ErrInvalidIterator is returned when the iterator is invalid.
	ErrInvalidIterator = awserrors.NewAWSError("InvalidArgumentException", "Invalid iterator.", http.StatusBadRequest)
	// ErrShardClosed is returned when the shard is closed.
	ErrShardClosed = awserrors.NewAWSError("InvalidArgumentException", "Shard is closed.", http.StatusBadRequest)
	// ErrAccessDenied is returned when access is denied.
	ErrAccessDenied = awserrors.NewAWSError("AccessDeniedException", "Access denied.", http.StatusForbidden)
	// ErrResourceAlreadyExists is returned when attempting to create a resource that already exists.
	ErrResourceAlreadyExists = awserrors.NewAWSError("ResourceAlreadyExistsException", "Resource already exists.", http.StatusConflict)
	// ErrKMSAccessDenied is returned when KMS access is denied.
	ErrKMSAccessDenied = awserrors.NewAWSError("KMSAccessDeniedException", "KMS access denied.", http.StatusForbidden)
	// ErrKMSDisabled is returned when the KMS key is disabled.
	ErrKMSDisabled = awserrors.NewAWSError("KMSDisabledException", "KMS key is disabled.", http.StatusBadRequest)
	// ErrKMSNotFound is returned when the KMS key is not found.
	ErrKMSNotFound = awserrors.NewAWSError("KMSNotFoundException", "KMS key not found.", http.StatusNotFound)
	// ErrKMSThrottling is returned when KMS requests are throttled.
	ErrKMSThrottling = awserrors.NewAWSError("KMSThrottlingException", "KMS request throttled.", http.StatusTooManyRequests)
	// ErrValidation is returned when validation fails.
	ErrValidation = awserrors.NewAWSError("ValidationException", "Validation error.", http.StatusBadRequest)
)
