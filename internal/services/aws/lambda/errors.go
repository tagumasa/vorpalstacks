// Package lambda provides AWS Lambda service operations for vorpalstacks.
package lambda

import (
	"fmt"
	"net/http"
	"strings"

	awserrors "vorpalstacks/internal/common/errors"
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

// GetLambdaError extracts a LambdaError from the given error, returning ErrServiceException if not found.
func GetLambdaError(err error) *LambdaError {
	if lambdaErr, ok := err.(*LambdaError); ok {
		return lambdaErr
	}
	return ErrServiceException
}

var validRuntimes = []string{
	"nodejs24.x", "nodejs22.x", "nodejs20.x",
	"python3.14", "python3.13", "python3.12", "python3.11", "python3.10",
	"java25", "java21", "java17", "java11", "java8.al2",
	"dotnet10", "dotnet8",
	"ruby3.4", "ruby3.3", "ruby3.2",
	"provided.al2023", "provided.al2",
}

// ValidateRuntime checks if the provided runtime is a valid Lambda runtime.
func ValidateRuntime(runtime string) bool {
	runtimeLower := strings.ToLower(runtime)
	for _, r := range validRuntimes {
		if r == runtimeLower {
			return true
		}
	}
	return false
}

// ValidateHandler validates the handler string for a Lambda function.
// Checks that the handler is not empty and conforms to runtime-specific format requirements.
func ValidateHandler(runtime, handler string) error {
	if handler == "" {
		return NewInvalidParameter("Handler", "Handler cannot be empty")
	}

	if strings.HasPrefix(runtime, "python") {
		if !strings.Contains(handler, ".") {
			return NewInvalidParameter("Handler", "Python handler must be in the format module.function")
		}
	}

	if strings.HasPrefix(runtime, "nodejs") {
		if !strings.Contains(handler, ".") {
			return NewInvalidParameter("Handler", "Node.js handler must be in the format file.function")
		}
	}

	if strings.HasPrefix(runtime, "java") {
		if !strings.Contains(handler, "::") && !strings.Contains(handler, ".") {
			return NewInvalidParameter("Handler", "Java handler must be in the format package.Class::method")
		}
	}

	return nil
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
