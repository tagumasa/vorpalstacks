package appsync

import (
	"fmt"
	"net/http"

	awserrors "vorpalstacks/internal/common/errors"
	appsyncstore "vorpalstacks/internal/store/aws/appsync"
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
	ErrServiceQuotaExceededException      = NewAppSyncError("ServiceQuotaExceededException", "The service quota has been exceeded.", http.StatusTooManyRequests)
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

// NewGraphQLSchemaException creates a GraphQLSchemaException with a custom message.
func NewGraphQLSchemaException(message string) *AppSyncError {
	return NewAppSyncError("GraphQLSchemaException", message, http.StatusBadRequest)
}

// storeErrorMappings maps store-level sentinel errors to AppSync service errors.
// Sourced from the 22 sentinel errors in internal/store/aws/appsync/errors.go.
var storeErrorMappings = []awserrors.StoreErrorMapping{
	// NotFound errors
	{Store: appsyncstore.ErrApiNotFound, AWS: NewNotFoundException("API")},
	{Store: appsyncstore.ErrChannelNamespaceNotFound, AWS: NewNotFoundException("Channel namespace")},
	{Store: appsyncstore.ErrGraphqlApiNotFound, AWS: NewNotFoundException("GraphQL API")},
	{Store: appsyncstore.ErrDataSourceNotFound, AWS: NewNotFoundException("Data source")},
	{Store: appsyncstore.ErrResolverNotFound, AWS: NewNotFoundException("Resolver")},
	{Store: appsyncstore.ErrFunctionNotFound, AWS: NewNotFoundException("Function")},
	{Store: appsyncstore.ErrTypeNotFound, AWS: NewNotFoundException("Type")},
	{Store: appsyncstore.ErrSchemaCreationNotFound, AWS: NewNotFoundException("Schema creation status")},
	{Store: appsyncstore.ErrApiKeyNotFound, AWS: NewNotFoundException("API key")},
	{Store: appsyncstore.ErrApiCacheNotFound, AWS: NewNotFoundException("API cache")},
	{Store: appsyncstore.ErrDomainNameNotFound, AWS: NewNotFoundException("Domain name")},
	{Store: appsyncstore.ErrApiAssociationNotFound, AWS: NewNotFoundException("API association")},
	{Store: appsyncstore.ErrMergedApiAssociationNotFound, AWS: NewNotFoundException("Merged API association")},
	// AlreadyExists errors
	{Store: appsyncstore.ErrApiAlreadyExists, AWS: NewConflictException("API already exists")},
	{Store: appsyncstore.ErrChannelNamespaceExists, AWS: NewConflictException("Channel namespace already exists")},
	{Store: appsyncstore.ErrGraphqlApiAlreadyExists, AWS: NewConflictException("GraphQL API already exists")},
	{Store: appsyncstore.ErrDataSourceAlreadyExists, AWS: NewConflictException("Data source already exists")},
	{Store: appsyncstore.ErrResolverAlreadyExists, AWS: NewConflictException("Resolver already exists")},
	{Store: appsyncstore.ErrFunctionAlreadyExists, AWS: NewConflictException("Function already exists")},
	{Store: appsyncstore.ErrTypeAlreadyExists, AWS: NewConflictException("Type already exists")},
	{Store: appsyncstore.ErrApiCacheAlreadyExists, AWS: NewConflictException("API cache already exists")},
	{Store: appsyncstore.ErrDomainNameAlreadyExists, AWS: NewConflictException("Domain name already exists")},
}

// mapStoreError converts a store-level error to the corresponding AppSync
// service error using the data-driven storeErrorMappings table. Unknown
// errors fall back to ErrInternalFailureException.
func mapStoreError(err error) (interface{}, error) {
	if err == nil {
		return nil, nil
	}
	mapped := awserrors.MapStoreError(err, storeErrorMappings)
	if mapped == err {
		return nil, ErrInternalFailureException
	}
	return nil, mapped
}
