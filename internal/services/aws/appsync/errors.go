package appsync

import (
	"fmt"
	"net/http"

	awserrors "vorpalstacks/internal/common/errors"
)

// AppSyncError wraps a generic AWSError with AppSync-specific JSON serialisation
// using the REST-JSON 1.0 protocol format.
type AppSyncError struct {
	*awserrors.AWSError
}

// Unwrap returns the underlying AWS error for error chain inspection.
func (e *AppSyncError) Unwrap() error {
	return e.AWSError
}

// NewAppSyncError creates a new AppSyncError with the specified error code, message, and HTTP status.
func NewAppSyncError(code, message string, httpStatus int) *AppSyncError {
	return &AppSyncError{
		AWSError: awserrors.NewAWSError(code, message, httpStatus),
	}
}

// ToJSON serialises this error in the REST-JSON 1.0 protocol format.
func (e *AppSyncError) ToJSON() string {
	return e.AWSError.ToJSONWithFormat("rest-json")
}

// Pre-defined sentinel errors matching the AppSync error types.
var (
	ErrNotFoundException                  = NewAppSyncError("NotFoundException", "Resource not found.", http.StatusNotFound)
	ErrBadRequestException                = NewAppSyncError("BadRequestException", "The request is not valid.", http.StatusBadRequest)
	ErrUnauthorizedException              = NewAppSyncError("UnauthorizedException", "You are not authorized to perform this operation.", http.StatusUnauthorized)
	ErrAccessDeniedException              = NewAppSyncError("AccessDeniedException", "Access denied.", http.StatusForbidden)
	ErrConflictException                  = NewAppSyncError("ConflictException", "The resource already exists.", http.StatusConflict)
	ErrConcurrentModificationException    = NewAppSyncError("ConcurrentModificationException", "The resource is being modified by another request.", http.StatusConflict)
	ErrLimitExceededException             = NewAppSyncError("LimitExceededException", "The limit has been exceeded.", http.StatusTooManyRequests)
	ErrApiLimitExceededException          = NewAppSyncError("ApiLimitExceededException", "The API limit has been exceeded.", http.StatusBadRequest)
	ErrApiKeyLimitExceededException       = NewAppSyncError("ApiKeyLimitExceededException", "The API key limit has been exceeded.", http.StatusBadRequest)
	ErrApiKeyValidityOutOfBoundsException = NewAppSyncError("ApiKeyValidityOutOfBoundsException", "The API key validity period is out of bounds.", http.StatusBadRequest)
	ErrGraphQLSchemaException             = NewAppSyncError("GraphQLSchemaException", "The GraphQL schema is not valid.", http.StatusBadRequest)
	ErrInternalFailureException           = NewAppSyncError("InternalFailureException", "An internal failure occurred.", http.StatusInternalServerError)
	ErrServiceQuotaExceededException      = NewAppSyncError("ServiceQuotaExceededException", "The service quota has been exceeded.", http.StatusPaymentRequired)
)

// NewNotFoundException creates a NotFoundException with the specified resource description.
func NewNotFoundException(resource string) *AppSyncError {
	return NewAppSyncError("NotFoundException", fmt.Sprintf("%s not found.", resource), http.StatusNotFound)
}

// NewBadRequestException creates a BadRequestException with a custom message.
func NewBadRequestException(message string) *AppSyncError {
	return NewAppSyncError("BadRequestException", message, http.StatusBadRequest)
}

// NewConflictException creates a ConflictException with a custom message.
func NewConflictException(message string) *AppSyncError {
	return NewAppSyncError("ConflictException", message, http.StatusConflict)
}
