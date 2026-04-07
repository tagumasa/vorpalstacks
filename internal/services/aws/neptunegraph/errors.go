package neptunegraph

import (
	"fmt"
	"net/http"

	awserrors "vorpalstacks/internal/services/aws/common/errors"
)

type httpError struct {
	statusCode int
	code       string
	message    string
}

func (e *httpError) Error() string {
	return e.message
}

func (e *httpError) GetHTTPStatus() int {
	return e.statusCode
}

func (e *httpError) GetAWSError() *awserrors.AWSError {
	return awserrors.NewAWSError(e.code, e.message, e.statusCode)
}

func newHTTPError(status int, code, message string) *httpError {
	return &httpError{statusCode: status, code: code, message: message}
}

func newValidationException(reason, message string) *httpError {
	return newHTTPError(http.StatusBadRequest, "ValidationException", fmt.Sprintf("%s: %s", reason, message))
}

func newResourceNotFoundException(resourceType, id string) *httpError {
	return newHTTPError(http.StatusNotFound, "ResourceNotFoundException", fmt.Sprintf("Resource %s '%s' not found", resourceType, id))
}

func newConflictException(reason string) *httpError {
	return newHTTPError(http.StatusConflict, "ConflictException", reason)
}

func newInternalServerException(err error) *httpError {
	return newHTTPError(http.StatusInternalServerError, "InternalServerException", err.Error())
}

func newThrottlingException() *httpError {
	return newHTTPError(http.StatusTooManyRequests, "ThrottlingException", "Rate exceeded")
}

func newUnprocessableException(reason string) *httpError {
	return newHTTPError(http.StatusUnprocessableEntity, "UnprocessableException", reason)
}
