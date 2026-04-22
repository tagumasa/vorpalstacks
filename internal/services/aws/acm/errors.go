package acm

import (
	"fmt"
	"net/http"

	awserrors "vorpalstacks/internal/common/errors"
)

// NewInvalidArnError creates a new invalid ARN error.
func NewInvalidArnError(arn string) *awserrors.AWSError {
	return awserrors.NewAWSError("InvalidArnException",
		fmt.Sprintf("Invalid ARN: %s", arn), http.StatusBadRequest)
}

// NewInvalidParameterError creates a new invalid parameter error.
func NewInvalidParameterError(message string) *awserrors.AWSError {
	return awserrors.NewInvalidParameterException(message)
}

// NewResourceInUseError creates a new resource in use error.
func NewResourceInUseError(resourceType, identifier string) *awserrors.AWSError {
	return awserrors.NewAWSError("ResourceInUseException",
		fmt.Sprintf("%s %s is in use", resourceType, identifier), http.StatusBadRequest)
}

// NewInvalidDomainValidationOptionsException creates a new invalid domain validation options exception.
func NewInvalidDomainValidationOptionsException(message string) *awserrors.AWSError {
	return awserrors.NewAWSError("InvalidDomainValidationOptionsException", message, http.StatusBadRequest)
}

// NewRequestInProgressException creates a new request in progress exception.
func NewRequestInProgressException(message string) *awserrors.AWSError {
	return awserrors.NewAWSError("RequestInProgressException", message, http.StatusBadRequest)
}

// NewInvalidStateException creates a new invalid state exception.
func NewInvalidStateException(message string) *awserrors.AWSError {
	return awserrors.NewAWSError("InvalidStateException", message, http.StatusBadRequest)
}
