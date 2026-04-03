package athena

import (
	awserrors "vorpalstacks/internal/services/aws/common/errors"
)

var (
	// ErrInvalidRequestException is returned when the input fails constraints.
	ErrInvalidRequestException = awserrors.NewBadRequestException("The input failed to satisfy the constraints specified by an AWS service.")
	// ErrInternalServerException is returned when an internal server error occurs.
	ErrInternalServerException = awserrors.NewInternalErrorException("An internal server error occurred.")
	// ErrInvalidParameterException is returned when a parameter value is invalid.
	ErrInvalidParameterException = awserrors.NewInvalidParameterException("Invalid parameter value.")
	// ErrResourceNotFoundException is returned when a resource is not found.
	ErrResourceNotFoundException = awserrors.NewResourceNotFoundException("Resource", "")
	// ErrResourceAlreadyExistsException is returned when a resource already exists.
	ErrResourceAlreadyExistsException = awserrors.NewResourceAlreadyExistsException("Resource")
	// ErrTooManyRequestsException is returned when too many requests are received.
	ErrTooManyRequestsException = awserrors.NewThrottlingException("Too many requests have been received.")
	// ErrMetadataException is returned when an error occurs while accessing metadata.
	ErrMetadataException = awserrors.NewInternalErrorException("An error occurred while accessing metadata.")
	// ErrInvalidConfigurationException is returned when the configuration is invalid.
	ErrInvalidConfigurationException = awserrors.NewBadRequestException("The configuration is invalid.")
)

// NewValidationError creates a validation error.
func NewValidationError(message string) *awserrors.AWSError {
	return awserrors.NewBadRequestException(message)
}

// NewResourceNotFoundError creates a resource not found error.
func NewResourceNotFoundError(message string) *awserrors.AWSError {
	return awserrors.NewNotFoundException(message)
}

// NewResourceAlreadyExistsError creates a resource already exists error.
func NewResourceAlreadyExistsError(message string) *awserrors.AWSError {
	return awserrors.NewResourceAlreadyExistsException(message)
}
