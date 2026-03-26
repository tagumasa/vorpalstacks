package athena

import (
	awserrors "vorpalstacks/internal/services/aws/common/errors"
)

// AthenaError represents an error from Athena.
type AthenaError struct {
	*awserrors.AWSError
}

// Unwrap returns the underlying AWS error.
func (e *AthenaError) Unwrap() error {
	return e.AWSError
}

var (
	// ErrInvalidRequestException is returned when the input fails constraints.
	ErrInvalidRequestException = &AthenaError{awserrors.NewBadRequestException("The input failed to satisfy the constraints specified by an AWS service.")}
	// ErrInternalServerException is returned when an internal server error occurs.
	ErrInternalServerException = &AthenaError{awserrors.NewInternalErrorException("An internal server error occurred.")}
	// ErrInvalidParameterException is returned when a parameter value is invalid.
	ErrInvalidParameterException = &AthenaError{awserrors.NewInvalidParameterException("Invalid parameter value.")}
	// ErrResourceNotFoundException is returned when a resource is not found.
	ErrResourceNotFoundException = &AthenaError{awserrors.NewResourceNotFoundException("Resource", "")}
	// ErrResourceAlreadyExistsException is returned when a resource already exists.
	ErrResourceAlreadyExistsException = &AthenaError{awserrors.NewResourceAlreadyExistsException("Resource")}
	// ErrTooManyRequestsException is returned when too many requests are received.
	ErrTooManyRequestsException = &AthenaError{awserrors.NewThrottlingException("Too many requests have been received.")}
	// ErrMetadataException is returned when an error occurs while accessing metadata.
	ErrMetadataException = &AthenaError{awserrors.NewInternalErrorException("An error occurred while accessing metadata.")}
	// ErrInvalidConfigurationException is returned when the configuration is invalid.
	ErrInvalidConfigurationException = &AthenaError{awserrors.NewBadRequestException("The configuration is invalid.")}
)

// NewValidationError creates a validation error.
func NewValidationError(message string) *AthenaError {
	return &AthenaError{awserrors.NewBadRequestException(message)}
}

// NewResourceNotFoundError creates a resource not found error.
func NewResourceNotFoundError(message string) *AthenaError {
	return &AthenaError{awserrors.NewNotFoundException(message)}
}

// NewResourceAlreadyExistsError creates a resource already exists error.
func NewResourceAlreadyExistsError(message string) *AthenaError {
	return &AthenaError{awserrors.NewResourceAlreadyExistsException(message)}
}
