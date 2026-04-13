package cognitoidentityprovider

import (
	"net/http"

	"vorpalstacks/internal/common/errors"
)

// NewCognitoError creates a new CognitoError with the specified code, message and status code.
func NewCognitoError(code, message string, statusCode int) *errors.AWSError {
	return errors.NewAWSError(code, message, statusCode)
}

var (
	// ErrInvalidParameter is returned when a parameter is invalid.
	ErrInvalidParameter = NewCognitoError("InvalidParameterException", "Invalid parameter", http.StatusBadRequest)
	// ErrResourceNotFound is returned when the specified resource does not exist.
	ErrResourceNotFound = NewCognitoError("ResourceNotFoundException", "Resource not found", http.StatusNotFound)
	// ErrResourceExists is returned when attempting to create a resource that already exists.
	ErrResourceExists = NewCognitoError("ResourceExistsException", "Resource already exists", http.StatusConflict)
	// ErrUserNotFound is returned when the specified user does not exist.
	ErrUserNotFound = NewCognitoError("UserNotFoundException", "User not found", http.StatusNotFound)
	// ErrUserAlreadyExists is returned when attempting to create a user that already exists.
	ErrUserAlreadyExists = NewCognitoError("UsernameExistsException", "User already exists", http.StatusConflict)
	// ErrGroupNotFound is returned when the specified group does not exist.
	ErrGroupNotFound = NewCognitoError("ResourceNotFoundException", "Group not found", http.StatusNotFound)
	// ErrGroupAlreadyExists is returned when attempting to create a group that already exists.
	ErrGroupAlreadyExists = NewCognitoError("GroupExistsException", "Group already exists", http.StatusConflict)
	// ErrNotAuthorized is returned when the request is not authorized.
	ErrNotAuthorized = NewCognitoError("NotAuthorizedException", "Not authorized", http.StatusForbidden)
	// ErrIncorrectPassword is returned when the username or password is incorrect.
	ErrIncorrectPassword = NewCognitoError("NotAuthorizedException", "Incorrect username or password", http.StatusForbidden)
	// ErrPasswordPolicyViolation is returned when the password does not meet policy requirements.
	ErrPasswordPolicyViolation = NewCognitoError("InvalidPasswordException", "Password does not conform to policy", http.StatusBadRequest)
	// ErrUserNotConfirmed is returned when the user is not confirmed.
	ErrUserNotConfirmed = NewCognitoError("UserNotConfirmedException", "User is not confirmed", http.StatusBadRequest)
	// ErrUserAlreadyConfirmed is returned when the user is already confirmed.
	ErrUserAlreadyConfirmed = NewCognitoError("UserAlreadyConfirmedException", "User is already confirmed", http.StatusBadRequest)
	// ErrInvalidPassword is returned when the password is invalid.
	ErrInvalidPassword = NewCognitoError("InvalidPasswordException", "Invalid password", http.StatusBadRequest)
	// ErrCodeMismatch is returned when the verification code does not match.
	ErrCodeMismatch = NewCognitoError("CodeMismatchException", "Invalid verification code", http.StatusBadRequest)
	// ErrExpiredCode is returned when the verification code has expired.
	ErrExpiredCode = NewCognitoError("ExpiredCodeException", "Invalid code provided, please request a code again", http.StatusBadRequest)
	// ErrTooManyRequests is returned when the request rate limit is exceeded.
	ErrTooManyRequests = NewCognitoError("TooManyRequestsException", "Too many requests", http.StatusTooManyRequests)
	// ErrInternalError is returned when an internal error occurs.
	ErrInternalError = NewCognitoError("InternalErrorException", "Internal error", http.StatusInternalServerError)
	// ErrClientNotFound is returned when the specified client does not exist.
	ErrClientNotFound = NewCognitoError("ResourceNotFoundException", "Client not found", http.StatusNotFound)
	// ErrClientAlreadyExists is returned when attempting to create a client that already exists.
	ErrClientAlreadyExists = NewCognitoError("ResourceExistsException", "Client already exists", http.StatusConflict)
)
