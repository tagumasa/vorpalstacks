// Package timestreamwrite provides Timestream Write service operations for vorpalstacks.
package timestreamwrite

import (
	"net/http"

	awserrors "vorpalstacks/internal/services/aws/common/errors"
)

// TimestreamWriteError represents an error returned by the Timestream Write service.
type TimestreamWriteError struct {
	*awserrors.AWSError
}

// Unwrap returns the underlying AWS error.
func (e *TimestreamWriteError) Unwrap() error {
	return e.AWSError
}

// GetCode returns the error code.
func (e *TimestreamWriteError) GetCode() string {
	return e.AWSError.GetCode()
}

// GetMessage returns the error message.
func (e *TimestreamWriteError) GetMessage() string {
	return e.AWSError.GetMessage()
}

// GetHTTPStatusCode returns the HTTP status code.
func (e *TimestreamWriteError) GetHTTPStatusCode() int {
	return e.AWSError.GetHTTPStatusCode()
}

// As implements the error interface's As method.
func (e *TimestreamWriteError) As(target interface{}) bool {
	if t, ok := target.(**awserrors.AWSError); ok {
		*t = e.AWSError
		return true
	}
	return false
}

var (
	// ErrResourceNotFound is returned when the specified resource does not exist.
	ErrResourceNotFound = &TimestreamWriteError{awserrors.NewAWSError("ResourceNotFoundException", "The requested resource could not be found.", http.StatusNotFound)}
	// ErrResourceAlreadyExists is returned when attempting to create a resource that already exists.
	ErrResourceAlreadyExists = &TimestreamWriteError{awserrors.NewAWSError("ResourceAlreadyExistsException", "The resource already exists.", http.StatusConflict)}
	// ErrInvalidParameter is returned when a parameter is invalid.
	ErrInvalidParameter = &TimestreamWriteError{awserrors.NewAWSError("ValidationException", "The input fails to satisfy the constraints specified by an AWS service.", http.StatusBadRequest)}
	// ErrValidationException is returned when validation fails.
	ErrValidationException = &TimestreamWriteError{awserrors.NewAWSError("ValidationException", "The input fails to satisfy the constraints specified by an AWS service.", http.StatusBadRequest)}
	// ErrAccessDenied is returned when access is denied.
	ErrAccessDenied = &TimestreamWriteError{awserrors.NewAWSError("AccessDeniedException", "You are not authorised to perform this action.", http.StatusForbidden)}
	// ErrThrottling is returned when the request is throttled.
	ErrThrottling = &TimestreamWriteError{awserrors.NewAWSError("ThrottlingException", "The request was denied due to request throttling.", http.StatusTooManyRequests)}
	// ErrInternalServer is returned when an internal server error occurs.
	ErrInternalServer = &TimestreamWriteError{awserrors.NewAWSError("InternalServerException", "Internal server error.", http.StatusInternalServerError)}
	// ErrServiceQuotaExceeded is returned when a service quota is exceeded.
	ErrServiceQuotaExceeded = &TimestreamWriteError{awserrors.NewAWSError("ServiceQuotaExceededException", "The request exceeded the service quota.", http.StatusTooManyRequests)}
)
