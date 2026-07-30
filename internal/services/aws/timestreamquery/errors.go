package timestreamquery

import (
	"net/http"

	awserrors "vorpalstacks/internal/common/errors"
)

var (
	// ErrValidationException is returned when the request is invalid.
	ErrValidationException = awserrors.NewAWSError("ValidationException", "The request is invalid.", http.StatusBadRequest)
	// ErrAccessDenied is returned when access is denied.
	ErrAccessDenied = awserrors.NewAWSError("AccessDeniedException", "Access denied.", http.StatusForbidden)
	// ErrConflictException is returned when a resource already exists or is in a conflicting state.
	ErrConflictException = awserrors.NewConflictException("The resource already exists.")
	// ErrInternalServer is returned when an internal server error occurs.
	ErrInternalServer = awserrors.NewAWSError("InternalServerException", "Internal server error.", http.StatusInternalServerError)
	// ErrQueryExecutionError is returned when query execution fails.
	ErrQueryExecutionError = awserrors.NewAWSError("QueryExecutionException", "Query execution failed.", http.StatusInternalServerError)
	// ErrResourceNotFound is returned when the specified resource does not exist.
	ErrResourceNotFound = awserrors.NewAWSError("ResourceNotFoundException", "The resource was not found.", http.StatusNotFound)
	// ErrThrottlingException is returned when the request is throttled.
	ErrThrottlingException = awserrors.NewAWSError("ThrottlingException", "The request was denied due to request throttling.", http.StatusTooManyRequests)
	// ErrServiceQuotaExceeded is returned when a service quota is exceeded.
	// Account-level quota enforcement is deferred to the v0.2.0 quota framework.
	ErrServiceQuotaExceeded = awserrors.NewAWSError("ServiceQuotaExceededException", "The request exceeded the service quota.", http.StatusTooManyRequests)
)
