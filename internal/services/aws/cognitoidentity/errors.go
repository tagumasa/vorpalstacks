package cognitoidentity

import (
	"net/http"

	awserrors "vorpalstacks/internal/services/aws/common/errors"
)

// CognitoIdentityError represents an error from Cognito Identity.
type CognitoIdentityError struct {
	*awserrors.AWSError
}

// Unwrap returns the underlying AWS error.
func (e *CognitoIdentityError) Unwrap() error {
	return e.AWSError
}

var (
	// ErrInvalidParameter is returned when an invalid parameter is provided.
	ErrInvalidParameter = &CognitoIdentityError{awserrors.NewAWSError("InvalidParameterException", "Invalid parameter", http.StatusBadRequest)}
	// ErrResourceNotFound is returned when a resource is not found.
	ErrResourceNotFound = &CognitoIdentityError{awserrors.NewAWSError("ResourceNotFoundException", "Resource not found", http.StatusNotFound)}
	// ErrInternalError is returned when an internal error occurs.
	ErrInternalError = &CognitoIdentityError{awserrors.NewAWSError("InternalErrorException", "Internal error", http.StatusInternalServerError)}
	// ErrNotAuthorized is returned when the user is not authorized.
	ErrNotAuthorized = &CognitoIdentityError{awserrors.NewAWSError("NotAuthorizedException", "Not authorized", http.StatusUnauthorized)}
	// ErrResourceInUse is returned when a resource is already in use.
	ErrResourceInUse = &CognitoIdentityError{awserrors.NewAWSError("ResourceInUseException", "Resource already exists", http.StatusConflict)}
)
