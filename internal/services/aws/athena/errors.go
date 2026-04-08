package athena

import (
	"net/http"

	awserrors "vorpalstacks/internal/services/aws/common/errors"
)

var (
	// ErrInvalidRequestException is returned when the input fails constraints.
	ErrInvalidRequestException        = awserrors.NewAWSError("InvalidRequestException", "The input failed to satisfy the constraints specified by an AWS service.", http.StatusBadRequest)
	// ErrInternalServerException is returned when an internal server error occurs.
	ErrInternalServerException        = awserrors.NewInternalErrorException("An internal server error occurred.")
	// ErrInvalidParameterException is returned when an invalid parameter value is provided.
	ErrInvalidParameterException      = awserrors.NewInvalidParameterException("Invalid parameter value.")
	ErrResourceNotFoundException      = awserrors.NewResourceNotFoundException("Resource", "")
	// ErrTooManyRequestsException is returned when too many requests have been received.
	ErrResourceAlreadyExistsException = awserrors.NewResourceAlreadyExistsException("Resource")
	// ErrMetadataException is returned when an error occurs while accessing metadata.
	ErrTooManyRequestsException       = awserrors.NewThrottlingException("Too many requests have been received.")
	// ErrInvalidConfigurationException is returned when the configuration is invalid.
	ErrMetadataException              = awserrors.NewAWSError("MetadataException", "An error occurred while accessing metadata.", http.StatusBadRequest)
	// ErrInvalidConfigurationException is returned when the configuration is invalid.
	ErrInvalidConfigurationException  = awserrors.NewBadRequestException("The configuration is invalid.")
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
