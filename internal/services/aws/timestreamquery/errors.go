package timestreamquery

import (
	"net/http"

	awserrors "vorpalstacks/internal/common/errors"
)

var (
	// ErrValidationException is returned when the request is invalid.
	ErrValidationException = awserrors.NewAWSError("ValidationException", "The request is invalid.", http.StatusBadRequest)
	// ErrConflictException is returned when a resource already exists or is in a conflicting state.
	ErrConflictException = awserrors.NewConflictException("The resource already exists.")
	// ErrInternalServer is returned when an internal server error occurs.
	ErrInternalServer = awserrors.NewAWSError("InternalServerException", "Internal server error.", http.StatusInternalServerError)
	// ErrQueryExecutionError is returned when query execution fails.
	// Smithy: smithy.api#error "client", smithy.api#httpError 400.
	ErrQueryExecutionError = awserrors.NewAWSError("QueryExecutionException", "Query execution failed.", http.StatusBadRequest)
	// ErrResourceNotFound is returned when the specified resource does not exist.
	ErrResourceNotFound = awserrors.NewAWSError("ResourceNotFoundException", "The resource was not found.", http.StatusNotFound)
)
