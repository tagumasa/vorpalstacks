package acm

import (
	"fmt"
	"net/http"

	awserrors "vorpalstacks/internal/services/aws/common/errors"
)

// APIError represents an error returned by the ACM API.
type APIError struct {
	*awserrors.AWSError
}

// Unwrap returns the underlying AWS error.
func (e *APIError) Unwrap() error {
	return e.AWSError
}

// NewAPIError creates a new ACM API error.
func NewAPIError(code, message string, httpStatus int) *APIError {
	return &APIError{
		AWSError: awserrors.NewAWSError(code, message, httpStatus),
	}
}

// NewResourceNotFoundError creates a new resource not found error.
func NewResourceNotFoundError(resourceType, identifier string) *APIError {
	return &APIError{awserrors.NewResourceNotFoundException(resourceType, identifier)}
}

// NewInvalidArnError creates a new invalid ARN error.
func NewInvalidArnError(arn string) *APIError {
	return NewAPIError("InvalidArnException",
		fmt.Sprintf("Invalid ARN: %s", arn), http.StatusBadRequest)
}

// NewInvalidParameterError creates a new invalid parameter error.
func NewInvalidParameterError(message string) *APIError {
	return &APIError{awserrors.NewInvalidParameterException(message)}
}

// NewResourceInUseError creates a new resource in use error.
func NewResourceInUseError(resourceType, identifier string) *APIError {
	return NewAPIError("ResourceInUseException",
		fmt.Sprintf("%s %s is in use", resourceType, identifier), http.StatusBadRequest)
}

// NewValidationException creates a new validation exception.
func NewValidationException(message string) *APIError {
	return &APIError{awserrors.NewValidationException(message)}
}

// NewInvalidDomainValidationOptionsException creates a new invalid domain validation options exception.
func NewInvalidDomainValidationOptionsException(message string) *APIError {
	return NewAPIError("InvalidDomainValidationOptionsException", message, http.StatusBadRequest)
}

// NewRequestInProgressException creates a new request in progress exception.
func NewRequestInProgressException(message string) *APIError {
	return NewAPIError("RequestInProgressException", message, http.StatusBadRequest)
}

// NewInvalidStateException creates a new invalid state exception.
func NewInvalidStateException(message string) *APIError {
	return NewAPIError("InvalidStateException", message, http.StatusBadRequest)
}

// NewConflictException creates a new conflict exception.
func NewConflictException(message string) *APIError {
	return &APIError{awserrors.NewConflictException(message)}
}

// NewAccessDeniedException creates a new access denied exception.
func NewAccessDeniedException(message string) *APIError {
	return &APIError{awserrors.NewAccessDeniedException(message)}
}

// NewLimitExceededException creates a new limit exceeded exception.
func NewLimitExceededException(message string) *APIError {
	return &APIError{awserrors.NewLimitExceededException(message)}
}
