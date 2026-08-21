// Package lambda provides AWS Lambda service operations for vorpalstacks.
package lambda

import (
	"errors"
	"fmt"
	"net/http"

	awserrors "vorpalstacks/internal/common/errors"
	lambdastore "vorpalstacks/internal/store/aws/lambda"
)

// LambdaError represents an error returned by the Lambda service.
type LambdaError struct {
	*awserrors.AWSError
}

// Unwrap returns the underlying error.
func (e *LambdaError) Unwrap() error {
	return e.AWSError
}

// NewLambdaError creates a new LambdaError with the specified code, message and HTTP status.
func NewLambdaError(code, message string, httpStatus int) *LambdaError {
	return &LambdaError{
		AWSError: awserrors.NewAWSError(code, message, httpStatus),
	}
}

// ToJSON returns the LambdaError as a JSON-formatted string.
func (e *LambdaError) ToJSON() string {
	return e.AWSError.ToJSONWithFormat("lambda")
}

var (
	// ErrResourceNotFound is returned when the specified Lambda resource does not exist.
	ErrResourceNotFound = NewLambdaError(
		"ResourceNotFoundException",
		"The resource specified in the request does not exist.",
		http.StatusNotFound,
	)

	// ErrResourceInUse is returned when the Lambda resource is already in use.
	ErrResourceInUse = NewLambdaError(
		"ResourceInUseException",
		"The resource is already in use.",
		http.StatusBadRequest,
	)

	// ErrInvalidParameterValue is returned when a parameter value is invalid.
	ErrInvalidParameterValue = NewLambdaError(
		"InvalidParameterValueException",
		"The value for the parameter is invalid.",
		http.StatusBadRequest,
	)

	// ErrInvalidRuntime is returned when the Lambda runtime is invalid.
	ErrInvalidRuntime = NewLambdaError(
		"InvalidParameterValueException",
		"The runtime parameter is invalid.",
		http.StatusBadRequest,
	)

	// ErrCodeVerificationFailed is returned when code signature verification fails.
	ErrCodeVerificationFailed = NewLambdaError(
		"CodeVerificationFailedException",
		"The code signature failed the signature verification check.",
		http.StatusBadRequest,
	)

	// ErrCodeStorageExceeded is returned when the total code size exceeds the account limit.
	ErrCodeStorageExceeded = NewLambdaError(
		"CodeStorageExceededException",
		"The total code size for the account exceeds the maximum allowed limit.",
		http.StatusBadRequest,
	)

	// ErrTooManyRequests is returned when the request rate limit is exceeded.
	ErrTooManyRequests = NewLambdaError(
		"TooManyRequestsException",
		"Too many requests have been made. Please retry.",
		http.StatusTooManyRequests,
	)

	// ErrServiceException is returned when an internal service error occurs.
	ErrServiceException = NewLambdaError(
		"ServiceException",
		"An internal service error occurred.",
		http.StatusInternalServerError,
	)

	// ErrRequestTooLarge is returned when the request payload exceeds the
	// maximum allowed size (6 MB for synchronous invocation, 256 KB for
	// asynchronous invocation).
	ErrRequestTooLarge = NewLambdaError(
		"RequestTooLargeException",
		"The request payload exceeds the maximum allowed size.",
		http.StatusRequestEntityTooLarge,
	)
)

// NewResourceNotFound creates a new LambdaError for a resource that was not found.
func NewResourceNotFound(resourceType, resourceName string) *LambdaError {
	return &LambdaError{
		AWSError: awserrors.NewResourceNotFoundException(resourceType, resourceName),
	}
}

// NewInvalidParameter creates a new LambdaError for an invalid parameter.
func NewInvalidParameter(paramName, message string) *LambdaError {
	return &LambdaError{
		AWSError: awserrors.NewAWSError("InvalidParameterValueException",
			fmt.Sprintf("Invalid parameter '%s': %s", paramName, message), http.StatusBadRequest),
	}
}

// NewResourceConflict creates a new LambdaError for a resource conflict.
func NewResourceConflict(message string) *LambdaError {
	return &LambdaError{
		AWSError: awserrors.NewAWSError("ResourceConflictException", message, http.StatusConflict),
	}
}

// IsLambdaError checks if the given error is a LambdaError.
func IsLambdaError(err error) bool {
	_, ok := err.(*LambdaError)
	return ok
}

// mapStoreError converts raw store-layer sentinel errors into the AWS error
// contract of the Lambda API. Handlers that surface store errors must route
// through this mapping so unmapped sentinels never escape to the dispatcher,
// which would answer them with an internal-error 500.
func mapStoreError(err error) error {
	switch {
	case err == nil:
		return nil
	case errors.Is(err, lambdastore.ErrEventSourceAlreadyExists):
		return NewResourceConflict("The event source mapping already exists for this event source and function.")
	case errors.Is(err, lambdastore.ErrResourceConflict):
		return NewResourceConflict("The resource conflicts with the current state of the function.")
	case errors.Is(err, lambdastore.ErrFunctionNotFound),
		errors.Is(err, lambdastore.ErrVersionNotFound),
		errors.Is(err, lambdastore.ErrLayerNotFound),
		errors.Is(err, lambdastore.ErrLayerVersionNotFound),
		errors.Is(err, lambdastore.ErrAliasNotFound),
		errors.Is(err, lambdastore.ErrPolicyNotFound):
		return ErrResourceNotFound
	}
	return err
}

// GetLambdaError extracts a LambdaError from the given error, returning ErrServiceException if not found.
func GetLambdaError(err error) *LambdaError {
	if lambdaErr, ok := err.(*LambdaError); ok {
		return lambdaErr
	}
	return ErrServiceException
}

// Response represents a response from a Lambda function invocation.
type Response struct {
	StatusCode int
	Body       interface{}
}

// AWSResponse creates a new Response with the specified status code and body.
func AWSResponse(status int, body interface{}) *Response {
	return &Response{
		StatusCode: status,
		Body:       body,
	}
}
