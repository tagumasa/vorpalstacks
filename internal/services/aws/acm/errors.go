package acm

import (
	"fmt"
	"net/http"

	awserrors "vorpalstacks/internal/services/aws/common/errors"
)

// NewAPIError creates a new ACM API error.
func NewAPIError(code, message string, httpStatus int) *awserrors.AWSError {
	return awserrors.NewAWSError(code, message, httpStatus)
}

// NewResourceNotFoundError creates a new resource not found error.
func NewResourceNotFoundError(resourceType, identifier string) *awserrors.AWSError {
	return awserrors.NewResourceNotFoundException(resourceType, identifier)
}

// NewInvalidArnError creates a new invalid ARN error.
func NewInvalidArnError(arn string) *awserrors.AWSError {
	return NewAPIError("InvalidArnException",
		fmt.Sprintf("Invalid ARN: %s", arn), http.StatusBadRequest)
}

// NewInvalidParameterError creates a new invalid parameter error.
func NewInvalidParameterError(message string) *awserrors.AWSError {
	return awserrors.NewInvalidParameterException(message)
}

// NewResourceInUseError creates a new resource in use error.
func NewResourceInUseError(resourceType, identifier string) *awserrors.AWSError {
	return NewAPIError("ResourceInUseException",
		fmt.Sprintf("%s %s is in use", resourceType, identifier), http.StatusBadRequest)
}

// NewValidationException creates a new validation exception.
func NewValidationException(message string) *awserrors.AWSError {
	return awserrors.NewValidationException(message)
}

// NewInvalidDomainValidationOptionsException creates a new invalid domain validation options exception.
func NewInvalidDomainValidationOptionsException(message string) *awserrors.AWSError {
	return NewAPIError("InvalidDomainValidationOptionsException", message, http.StatusBadRequest)
}

// NewRequestInProgressException creates a new request in progress exception.
func NewRequestInProgressException(message string) *awserrors.AWSError {
	return NewAPIError("RequestInProgressException", message, http.StatusBadRequest)
}

// NewInvalidStateException creates a new invalid state exception.
func NewInvalidStateException(message string) *awserrors.AWSError {
	return NewAPIError("InvalidStateException", message, http.StatusBadRequest)
}

// NewConflictException creates a new conflict exception.
func NewConflictException(message string) *awserrors.AWSError {
	return awserrors.NewConflictException(message)
}

// NewAccessDeniedException creates a new access denied exception.
func NewAccessDeniedException(message string) *awserrors.AWSError {
	return awserrors.NewAccessDeniedException(message)
}

// NewLimitExceededException creates a new limit exceeded exception.
func NewLimitExceededException(message string) *awserrors.AWSError {
	return awserrors.NewLimitExceededException(message)
}
