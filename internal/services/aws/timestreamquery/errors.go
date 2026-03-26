package timestreamquery

import (
	"net/http"

	awserrors "vorpalstacks/internal/services/aws/common/errors"
)

// TimestreamQueryError represents an error returned by the Timestream Query service.
type TimestreamQueryError struct {
	*awserrors.AWSError
}

// Unwrap returns the underlying AWS error.
func (e *TimestreamQueryError) Unwrap() error {
	return e.AWSError
}

var (
	// ErrValidationException is returned when the request is invalid.
	ErrValidationException = &TimestreamQueryError{awserrors.NewAWSError("ValidationException", "The request is invalid.", http.StatusBadRequest)}
	// ErrInvalidEndpoint is returned when the endpoint is invalid.
	ErrInvalidEndpoint = &TimestreamQueryError{awserrors.NewAWSError("InvalidEndpoint", "The endpoint is invalid.", http.StatusBadRequest)}
	// ErrAccessDenied is returned when access is denied.
	ErrAccessDenied = &TimestreamQueryError{awserrors.NewAWSError("AccessDeniedException", "Access denied.", http.StatusForbidden)}
	// ErrInternalServer is returned when an internal server error occurs.
	ErrInternalServer = &TimestreamQueryError{awserrors.NewAWSError("InternalServerException", "Internal server error.", http.StatusInternalServerError)}
	// ErrQueryExecutionError is returned when query execution fails.
	ErrQueryExecutionError = &TimestreamQueryError{awserrors.NewAWSError("QueryExecutionException", "Query execution failed.", http.StatusInternalServerError)}
	// ErrResourceNotFound is returned when the specified resource does not exist.
	ErrResourceNotFound = &TimestreamQueryError{awserrors.NewAWSError("ResourceNotFoundException", "The resource was not found.", http.StatusNotFound)}
	// ErrResourceAlreadyExists is returned when attempting to create a resource that already exists.
	ErrResourceAlreadyExists = &TimestreamQueryError{awserrors.NewAWSError("ResourceAlreadyExistsException", "The resource already exists.", http.StatusConflict)}
	// ErrThrottlingException is returned when the request is throttled.
	ErrThrottlingException = &TimestreamQueryError{awserrors.NewAWSError("ThrottlingException", "The request was denied due to request throttling.", http.StatusTooManyRequests)}
	// ErrServiceQuotaExceeded is returned when a service quota is exceeded.
	ErrServiceQuotaExceeded = &TimestreamQueryError{awserrors.NewAWSError("ServiceQuotaExceededException", "The request exceeded the service quota.", http.StatusTooManyRequests)}
)
