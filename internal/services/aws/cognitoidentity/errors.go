package cognitoidentity

import (
	"net/http"

	awserrors "vorpalstacks/internal/common/errors"
)

var (
	// ErrInvalidParameter is returned when an invalid parameter is provided.
	ErrInvalidParameter = awserrors.NewAWSError("InvalidParameterException", "Invalid parameter", http.StatusBadRequest)
	// ErrResourceNotFound is returned when a resource is not found.
	ErrResourceNotFound = awserrors.NewAWSError("ResourceNotFoundException", "Resource not found", http.StatusNotFound)
	// ErrInternalError is returned when an internal error occurs.
	ErrInternalError = awserrors.NewAWSError("InternalErrorException", "Internal error", http.StatusInternalServerError)
	// ErrNotAuthorized is returned when the user is not authorized.
	ErrNotAuthorized = awserrors.NewAWSError("NotAuthorizedException", "Not authorized", http.StatusUnauthorized)
	// ErrResourceInUse is returned when a resource is already in use.
	ErrResourceInUse = awserrors.NewAWSError("ResourceInUseException", "Resource already exists", http.StatusConflict)
)
