package cognitoidentityprovider

import (
	"net/http"

	awserrors "vorpalstacks/internal/common/errors"
)

var (
	// ErrInvalidParameter is returned when a parameter is invalid.
	ErrInvalidParameter = awserrors.NewAWSError("InvalidParameterException", "Invalid parameter", http.StatusBadRequest)
	// ErrResourceNotFound is returned when the specified resource does not exist.
	ErrResourceNotFound = awserrors.NewAWSError("ResourceNotFoundException", "Resource not found", http.StatusNotFound)
	// ErrResourceExists is returned when attempting to create a resource that already exists.
	ErrResourceExists = awserrors.NewAWSError("ResourceExistsException", "Resource already exists", http.StatusConflict)
	// ErrUserNotFound is returned when the specified user does not exist.
	ErrUserNotFound = awserrors.NewAWSError("UserNotFoundException", "User not found", http.StatusNotFound)
	// ErrUserAlreadyExists is returned when attempting to create a user that already exists.
	ErrUserAlreadyExists = awserrors.NewAWSError("UsernameExistsException", "User already exists", http.StatusConflict)
	// ErrGroupNotFound is returned when the specified group does not exist.
	ErrGroupNotFound = awserrors.NewAWSError("ResourceNotFoundException", "Group not found", http.StatusNotFound)
	// ErrGroupAlreadyExists is returned when attempting to create a group that already exists.
	ErrGroupAlreadyExists = awserrors.NewAWSError("GroupExistsException", "Group already exists", http.StatusConflict)
	// ErrNotAuthorized is returned when the request is not authorized.
	ErrNotAuthorized = awserrors.NewAWSError("NotAuthorizedException", "Not authorized", http.StatusForbidden)
	// ErrIncorrectPassword is returned when the username or password is incorrect.
	ErrIncorrectPassword = awserrors.NewAWSError("NotAuthorizedException", "Incorrect username or password", http.StatusForbidden)
	// ErrPasswordPolicyViolation is returned when the password does not meet policy requirements.
	ErrPasswordPolicyViolation = awserrors.NewAWSError("InvalidPasswordException", "Password does not conform to policy", http.StatusBadRequest)
	// ErrUserNotConfirmed is returned when the user is not confirmed.
	ErrUserNotConfirmed = awserrors.NewAWSError("UserNotConfirmedException", "User is not confirmed", http.StatusBadRequest)
	// ErrUserAlreadyConfirmed is returned when the user is already confirmed.
	ErrUserAlreadyConfirmed = awserrors.NewAWSError("UserAlreadyConfirmedException", "User is already confirmed", http.StatusBadRequest)
	// ErrInvalidPassword is returned when the password is invalid.
	ErrInvalidPassword = awserrors.NewAWSError("InvalidPasswordException", "Invalid password", http.StatusBadRequest)
	// ErrCodeMismatch is returned when the verification code does not match.
	ErrCodeMismatch = awserrors.NewAWSError("CodeMismatchException", "Invalid verification code", http.StatusBadRequest)
	// ErrExpiredCode is returned when the verification code has expired.
	ErrExpiredCode = awserrors.NewAWSError("ExpiredCodeException", "Invalid code provided, please request a code again", http.StatusBadRequest)
	// ErrTooManyRequests is returned when the request rate limit is exceeded.
	ErrTooManyRequests = awserrors.NewAWSError("TooManyRequestsException", "Too many requests", http.StatusTooManyRequests)
	// ErrInternalError is returned when an internal error occurs.
	ErrInternalError = awserrors.NewAWSError("InternalErrorException", "Internal error", http.StatusInternalServerError)
	// ErrClientNotFound is returned when the specified client does not exist.
	ErrClientNotFound = awserrors.NewAWSError("ResourceNotFoundException", "Client not found", http.StatusNotFound)
	// ErrClientAlreadyExists is returned when attempting to create a client that already exists.
	ErrClientAlreadyExists = awserrors.NewAWSError("ResourceExistsException", "Client already exists", http.StatusConflict)
)
