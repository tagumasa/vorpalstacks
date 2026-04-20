package athena

import (
	"net/http"

	awserrors "vorpalstacks/internal/common/errors"
)

var (
	// ErrInvalidRequestException is returned when the input fails constraints.
	ErrInvalidRequestException = awserrors.NewAWSError("InvalidRequestException", "The input failed to satisfy the constraints specified by an AWS service.", http.StatusBadRequest)
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
	ErrMetadataException = awserrors.NewAWSError("MetadataException", "An error occurred while accessing metadata.", http.StatusBadRequest)
	// ErrInvalidConfigurationException is returned when the configuration is invalid.
	ErrInvalidConfigurationException = awserrors.NewBadRequestException("The configuration is invalid.")
)


