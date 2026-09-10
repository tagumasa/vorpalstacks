package apigateway

import (
	"fmt"
	"net/http"

	awserrors "vorpalstacks/internal/common/errors"

	storeerrors "vorpalstacks/internal/store/aws/apigateway"
)

// ApiGatewayError represents an error returned by API Gateway operations.
type ApiGatewayError struct {
	*awserrors.AWSError
}

// Unwrap returns the underlying AWS error.
func (e *ApiGatewayError) Unwrap() error {
	return e.AWSError
}

// NewApiGatewayError creates a new ApiGatewayError with the specified code, message, and HTTP status.
func NewApiGatewayError(code, message string, httpStatus int) *ApiGatewayError {
	return &ApiGatewayError{
		AWSError: awserrors.NewAWSError(code, message, httpStatus),
	}
}

var (
	ErrNotFoundException = NewApiGatewayError("NotFoundException", "The resource specified in the request does not exist.", http.StatusNotFound)
	ErrConflictException = NewApiGatewayError("ConflictException", "The resource already exists.", http.StatusConflict)
	// InternalFailure is the only 500 code the API Gateway contract documents
	// (Common Error Types); it is the service-wide internal-failure taxonomy.
	ErrInternalFailureException = NewApiGatewayError("InternalFailure", "The request can't be processed right now because of an internal server issue.", http.StatusInternalServerError)
)

// NewNotFoundException creates a new not found exception with the specified resource type and name.
func NewNotFoundException(resourceType, resourceName string) *ApiGatewayError {
	return NewApiGatewayError("NotFoundException", fmt.Sprintf("%s '%s' not found.", resourceType, resourceName), http.StatusNotFound)
}

// NewBadRequestException creates a new bad request exception with the specified message.
func NewBadRequestException(message string) *ApiGatewayError {
	return NewApiGatewayError("BadRequestException", message, http.StatusBadRequest)
}

// NewConflictException creates a new conflict exception with the specified message.
func NewConflictException(message string) *ApiGatewayError {
	return NewApiGatewayError("ConflictException", message, http.StatusConflict)
}

// NewInternalFailureException creates a new internal failure exception with
// the specified message.
func NewInternalFailureException(message string) *ApiGatewayError {
	return NewApiGatewayError("InternalFailure", message, http.StatusInternalServerError)
}

// storeErrorMappings maps store-level sentinel errors to API Gateway API
// errors. This follows the data-driven pattern used by CloudTrail, Kinesis,
// SecretsManager and other services.
var storeErrorMappings = []awserrors.StoreErrorMapping{
	// NotFound errors
	{Store: storeerrors.ErrRestApiNotFound, AWS: ErrNotFoundException},
	{Store: storeerrors.ErrResourceNotFound, AWS: ErrNotFoundException},
	{Store: storeerrors.ErrMethodNotFound, AWS: ErrNotFoundException},
	{Store: storeerrors.ErrIntegrationNotFound, AWS: ErrNotFoundException},
	{Store: storeerrors.ErrDeploymentNotFound, AWS: ErrNotFoundException},
	{Store: storeerrors.ErrStageNotFound, AWS: ErrNotFoundException},
	{Store: storeerrors.ErrRequestValidatorNotFound, AWS: ErrNotFoundException},
	{Store: storeerrors.ErrModelNotFound, AWS: ErrNotFoundException},
	{Store: storeerrors.ErrGatewayResponseNotFound, AWS: ErrNotFoundException},
	{Store: storeerrors.ErrApiKeyNotFound, AWS: ErrNotFoundException},
	{Store: storeerrors.ErrUsagePlanNotFound, AWS: ErrNotFoundException},
	{Store: storeerrors.ErrUsagePlanKeyNotFound, AWS: ErrNotFoundException},
	{Store: storeerrors.ErrDomainNameNotFound, AWS: ErrNotFoundException},
	{Store: storeerrors.ErrBasePathMappingNotFound, AWS: ErrNotFoundException},
	{Store: storeerrors.ErrAuthorizerNotFound, AWS: ErrNotFoundException},
	{Store: storeerrors.ErrMethodResponseNotFound, AWS: ErrNotFoundException},
	{Store: storeerrors.ErrIntegrationResponseNotFound, AWS: ErrNotFoundException},
	// AlreadyExists errors
	{Store: storeerrors.ErrRestApiAlreadyExists, AWS: ErrConflictException},
	{Store: storeerrors.ErrResourceAlreadyExists, AWS: ErrConflictException},
	{Store: storeerrors.ErrDeploymentAlreadyExists, AWS: ErrConflictException},
	{Store: storeerrors.ErrStageAlreadyExists, AWS: ErrConflictException},
	{Store: storeerrors.ErrRequestValidatorAlreadyExists, AWS: ErrConflictException},
	{Store: storeerrors.ErrModelAlreadyExists, AWS: ErrConflictException},
	{Store: storeerrors.ErrApiKeyAlreadyExists, AWS: ErrConflictException},
	{Store: storeerrors.ErrUsagePlanAlreadyExists, AWS: ErrConflictException},
	{Store: storeerrors.ErrUsagePlanKeyAlreadyExists, AWS: ErrConflictException},
	{Store: storeerrors.ErrDomainNameAlreadyExists, AWS: ErrConflictException},
	{Store: storeerrors.ErrBasePathMappingAlreadyExists, AWS: ErrConflictException},
	{Store: storeerrors.ErrAuthorizerAlreadyExists, AWS: ErrConflictException},
}

// toApiGatewayError converts a generic error to an ApiGatewayError,
// properly mapping store-level errors to the correct API Gateway error types.
func toApiGatewayError(err error) *ApiGatewayError {
	if err == nil {
		return nil
	}
	if apiErr, ok := err.(*ApiGatewayError); ok {
		return apiErr
	}
	mapped := awserrors.MapStoreError(err, storeErrorMappings)
	if _, ok := mapped.(*ApiGatewayError); ok {
		return mapped.(*ApiGatewayError)
	}
	return ErrInternalFailureException
}
