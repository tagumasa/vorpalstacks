package neptunegraph

import (
	"fmt"
	"net/http"

	awserrors "vorpalstacks/internal/common/errors"
)

// httpError represents a NeptuneGraph service error with HTTP status and AWS error code.
type httpError struct {
	statusCode int
	code       string
	message    string
}

// Error returns the human-readable error message.
func (e *httpError) Error() string {
	return e.message
}

// GetHTTPStatus returns the HTTP status code for this error.
func (e *httpError) GetHTTPStatus() int {
	return e.statusCode
}

// GetAWSError returns the structured AWS error representation.
func (e *httpError) GetAWSError() *awserrors.AWSError {
	return awserrors.NewAWSError(e.code, e.message, e.statusCode)
}

// newHTTPError creates an httpError with the given status, code, and message.
func newHTTPError(status int, code, message string) *httpError {
	return &httpError{statusCode: status, code: code, message: message}
}

// newValidationException creates a 400 ValidationException error.
func newValidationException(reason, message string) *httpError {
	return newHTTPError(http.StatusBadRequest, "ValidationException", fmt.Sprintf("%s: %s", reason, message))
}

// newResourceNotFoundException creates a 404 ResourceNotFoundException error.
func newResourceNotFoundException(resourceType, id string) *httpError {
	return newHTTPError(http.StatusNotFound, "ResourceNotFoundException", fmt.Sprintf("Resource %s '%s' not found", resourceType, id))
}

// newConflictException creates a 409 ConflictException error.
func newConflictException(reason string) *httpError {
	return newHTTPError(http.StatusConflict, "ConflictException", reason)
}

// newInternalServerException creates a 500 InternalServerException from an underlying error.
func newInternalServerException(err error) *httpError {
	return newHTTPError(http.StatusInternalServerError, "InternalServerException", err.Error())
}
