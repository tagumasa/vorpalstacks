package kinesis

import (
	"net/http"

	awserrors "vorpalstacks/internal/services/aws/common/errors"
)

// KinesisError represents an error returned by the Kinesis service.
type KinesisError struct {
	*awserrors.AWSError
}

// Unwrap returns the underlying AWS error.
func (e *KinesisError) Unwrap() error {
	return e.AWSError
}

var (
	// ErrResourceNotFound is returned when the specified resource does not exist.
	ErrResourceNotFound = &KinesisError{awserrors.NewAWSError("ResourceNotFoundException", "Requested resource not found.", http.StatusNotFound)}
	// ErrResourceInUse is returned when the resource is in use.
	ErrResourceInUse = &KinesisError{awserrors.NewAWSError("ResourceInUseException", "The resource is in use.", http.StatusBadRequest)}
	// ErrInvalidArgument is returned when an argument is invalid.
	ErrInvalidArgument = &KinesisError{awserrors.NewAWSError("InvalidArgumentException", "Invalid argument.", http.StatusBadRequest)}
	// ErrLimitExceeded is returned when the rate limit for the stream is exceeded.
	ErrLimitExceeded = &KinesisError{awserrors.NewAWSError("LimitExceededException", "Rate exceeded for this stream.", http.StatusTooManyRequests)}
	// ErrProvisionedThroughputExceeded is returned when the provisioned throughput is exceeded.
	ErrProvisionedThroughputExceeded = &KinesisError{awserrors.NewAWSError("ProvisionedThroughputExceededException", "Rate exceeded for this shard.", http.StatusTooManyRequests)}
	// ErrExpiredIterator is returned when the iterator has expired.
	ErrExpiredIterator = &KinesisError{awserrors.NewAWSError("ExpiredIteratorException", "Iterator expired.", http.StatusBadRequest)}
	// ErrInvalidIterator is returned when the iterator is invalid.
	ErrInvalidIterator = &KinesisError{awserrors.NewAWSError("InvalidArgumentException", "Invalid iterator.", http.StatusBadRequest)}
	// ErrShardClosed is returned when the shard is closed.
	ErrShardClosed = &KinesisError{awserrors.NewAWSError("InvalidArgumentException", "Shard is closed.", http.StatusBadRequest)}
	// ErrAccessDenied is returned when access is denied.
	ErrAccessDenied = &KinesisError{awserrors.NewAWSError("AccessDeniedException", "Access denied.", http.StatusForbidden)}
	// ErrResourceAlreadyExists is returned when attempting to create a resource that already exists.
	ErrResourceAlreadyExists = &KinesisError{awserrors.NewAWSError("ResourceAlreadyExistsException", "Resource already exists.", http.StatusConflict)}
	// ErrKMSAccessDenied is returned when KMS access is denied.
	ErrKMSAccessDenied = &KinesisError{awserrors.NewAWSError("KMSAccessDeniedException", "KMS access denied.", http.StatusForbidden)}
	// ErrKMSDisabled is returned when the KMS key is disabled.
	ErrKMSDisabled = &KinesisError{awserrors.NewAWSError("KMSDisabledException", "KMS key is disabled.", http.StatusBadRequest)}
	// ErrKMSNotFound is returned when the KMS key is not found.
	ErrKMSNotFound = &KinesisError{awserrors.NewAWSError("KMSNotFoundException", "KMS key not found.", http.StatusNotFound)}
	// ErrKMSThrottling is returned when KMS requests are throttled.
	ErrKMSThrottling = &KinesisError{awserrors.NewAWSError("KMSThrottlingException", "KMS request throttled.", http.StatusTooManyRequests)}
	// ErrValidation is returned when validation fails.
	ErrValidation = &KinesisError{awserrors.NewAWSError("ValidationException", "Validation error.", http.StatusBadRequest)}
)
