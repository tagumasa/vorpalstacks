// Package waf provides WAF (Web Application Firewall) service operations for vorpalstacks.
package waf

import (
	awserrors "vorpalstacks/internal/services/aws/common/errors"
)

// APIError represents an error returned by the WAF service.
type APIError struct {
	*awserrors.AWSError
}

// Unwrap returns the underlying AWS error.
func (e *APIError) Unwrap() error {
	return e.AWSError
}

// NewAPIError creates a new APIError with the specified code, message and HTTP status.
func NewAPIError(code, message string, httpStatus int) *APIError {
	return &APIError{
		AWSError: awserrors.NewAWSError(code, message, httpStatus),
	}
}

// NewValidationException creates a validation exception error.
func NewValidationException(message string) *APIError {
	return &APIError{
		AWSError: awserrors.NewValidationException(message),
	}
}

// NewResourceNotFoundException creates a resource not found error.
func NewResourceNotFoundException(message string) *APIError {
	return &APIError{
		AWSError: awserrors.NewResourceNotFoundException("Resource", message),
	}
}
