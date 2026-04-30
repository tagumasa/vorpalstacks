package apigateway

import (
	"errors"
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

// ToJSON returns the error as a JSON string.
func (e *ApiGatewayError) ToJSON() string {
	return e.AWSError.ToJSONWithFormat("rest-json")
}

// GetApiGatewayError converts a generic error to an ApiGatewayError.
var (
	ErrNotFoundException        = NewApiGatewayError("NotFoundException", "The resource specified in the request does not exist.", http.StatusNotFound)
	ErrBadRequestException      = NewApiGatewayError("BadRequestException", "The request is not valid.", http.StatusBadRequest)
	ErrConflictException        = NewApiGatewayError("ConflictException", "The resource already exists.", http.StatusConflict)
	ErrTooManyRequestsException = NewApiGatewayError("TooManyRequestsException", "Too many requests have been made.", http.StatusTooManyRequests)
	ErrServiceException         = NewApiGatewayError("ServiceException", "An internal service error occurred.", http.StatusInternalServerError)
	ErrAccessDeniedException    = NewApiGatewayError("AccessDeniedException", "Access denied.", http.StatusForbidden)
	ErrLimitExceededException   = NewApiGatewayError("LimitExceededException", "The limit has been exceeded.", http.StatusTooManyRequests)
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

// GetApiGatewayError converts a generic error to an ApiGatewayError.
func GetApiGatewayError(err error) *ApiGatewayError {
	if apiErr, ok := err.(*ApiGatewayError); ok {
		return apiErr
	}
	if isNotFoundError(err) {
		return ErrNotFoundException
	}
	return ErrServiceException
}

func isNotFoundError(err error) bool {
	notFoundErrors := []error{
		storeerrors.ErrRestApiNotFound,
		storeerrors.ErrResourceNotFound,
		storeerrors.ErrMethodNotFound,
		storeerrors.ErrIntegrationNotFound,
		storeerrors.ErrDeploymentNotFound,
		storeerrors.ErrStageNotFound,
		storeerrors.ErrRequestValidatorNotFound,
		storeerrors.ErrModelNotFound,
		storeerrors.ErrApiKeyNotFound,
		storeerrors.ErrUsagePlanNotFound,
		storeerrors.ErrUsagePlanKeyNotFound,
		storeerrors.ErrDomainNameNotFound,
		storeerrors.ErrBasePathMappingNotFound,
		storeerrors.ErrAuthorizerNotFound,
		storeerrors.ErrMethodResponseNotFound,
		storeerrors.ErrIntegrationResponseNotFound,
	}
	for _, notFound := range notFoundErrors {
		if errors.Is(err, notFound) {
			return true
		}
	}
	return false
}
