package cognitoidentityprovider

import (
	"net/http"

	awserrors "vorpalstacks/internal/common/errors"
)

var (
	// Amazon Cognito documents HTTP 400 for its entire client-error family,
	// including ResourceNotFoundException, UserNotFoundException,
	// NotAuthorizedException and TooManyRequestsException (the per-operation
	// Errors sections of the API reference) — unlike the 404/403/409/429
	// conventions most other services follow.
	// ErrInvalidParameter is returned when a parameter is invalid.
	ErrInvalidParameter = awserrors.NewAWSError("InvalidParameterException", "Invalid parameter", http.StatusBadRequest)
	// ErrWebAuthnChallengeNotFound is returned when a WebAuthn registration
	// completion has no matching pending challenge for the user: either
	// StartWebAuthnRegistration was never issued or its challenge expired.
	ErrWebAuthnChallengeNotFound = awserrors.NewAWSError("WebAuthnChallengeNotFoundException", "The challenge from StartWebAuthn registration has expired", http.StatusBadRequest)
	// ErrWebAuthnClientMismatch is returned when a WebAuthn registration
	// completion arrives with an access token issued to a different app
	// client than the one that started the registration.
	ErrWebAuthnClientMismatch = awserrors.NewAWSError("WebAuthnClientMismatchException", "The access token is for a different client than the one in the original StartWebAuthnRegistration request", http.StatusBadRequest)
	// ErrWebAuthnOriginNotAllowed is returned when the passkey credential's
	// registration origin does not align with the user pool relying party id.
	ErrWebAuthnOriginNotAllowed = awserrors.NewAWSError("WebAuthnOriginNotAllowedException", "The passkey credential's registration origin does not align with the user pool relying party id", http.StatusBadRequest)
	// ErrWebAuthnRelyingPartyMismatch is returned when the passkey credential
	// is associated with a different relying party ID than the user pool
	// relying party ID.
	ErrWebAuthnRelyingPartyMismatch = awserrors.NewAWSError("WebAuthnRelyingPartyMismatchException", "The passkey credential is associated with a different relying party ID than the user pool relying party ID", http.StatusBadRequest)
	// ErrWebAuthnCredentialNotSupported is returned when the passkey
	// credential comes from an unsupported device or provider: its response
	// cannot be interpreted or its key algorithm is not one the user pool
	// offered in the credential creation options.
	ErrWebAuthnCredentialNotSupported = awserrors.NewAWSError("WebAuthnCredentialNotSupportedException", "Passkey credentials from an unsupported device or provider", http.StatusBadRequest)
	// ErrManagedLoginBrandingExists is returned when attempting to assign a
	// managed login branding style to an app client that already has one.
	ErrManagedLoginBrandingExists = awserrors.NewAWSError("ManagedLoginBrandingExistsException", "The app client already has an assigned managed login branding style", http.StatusBadRequest)
	// ErrTermsExists is returned when attempting to create terms documents
	// with a TermsName that is already assigned to the app client.
	ErrTermsExists = awserrors.NewAWSError("TermsExistsException", "Terms document names must be unique to the app client", http.StatusBadRequest)
	// ErrPreconditionNotMet is returned when the request preconditions are
	// not met (e.g. starting a user import job when the pool has no
	// auto-verified attribute or another import job is already active).
	ErrPreconditionNotMet = awserrors.NewAWSError("PreconditionNotMetException", "Precondition not met", http.StatusBadRequest)
	// ErrResourceNotFound is returned when the specified resource does not exist.
	ErrResourceNotFound = awserrors.NewAWSError("ResourceNotFoundException", "Resource not found", http.StatusBadRequest)
	// ErrResourceExists is returned when attempting to create a resource that already exists.
	ErrResourceExists = awserrors.NewAWSError("ResourceExistsException", "Resource already exists", http.StatusConflict)
	// ErrUserNotFound is returned when the specified user does not exist.
	ErrUserNotFound = awserrors.NewAWSError("UserNotFoundException", "User not found", http.StatusBadRequest)
	// ErrUserAlreadyExists is returned when attempting to create a user that already exists.
	ErrUserAlreadyExists = awserrors.NewAWSError("UsernameExistsException", "User already exists", http.StatusBadRequest)
	// ErrGroupNotFound is returned when the specified group does not exist.
	ErrGroupNotFound = awserrors.NewAWSError("ResourceNotFoundException", "Group not found", http.StatusBadRequest)
	// ErrGroupAlreadyExists is returned when attempting to create a group that already exists.
	ErrGroupAlreadyExists = awserrors.NewAWSError("GroupExistsException", "Group already exists", http.StatusBadRequest)
	// ErrNotAuthorized is returned when the request is not authorized.
	ErrNotAuthorized = awserrors.NewAWSError("NotAuthorizedException", "Not authorized", http.StatusBadRequest)
	// ErrIncorrectPassword is returned when the username or password is incorrect.
	ErrIncorrectPassword = awserrors.NewAWSError("NotAuthorizedException", "Incorrect username or password", http.StatusBadRequest)
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
	ErrTooManyRequests = awserrors.NewAWSError("TooManyRequestsException", "Too many requests", http.StatusBadRequest)
	// ErrInternalError is returned when an internal error occurs.
	ErrInternalError = awserrors.NewAWSError("InternalErrorException", "Internal error", http.StatusInternalServerError)
	// ErrClientNotFound is returned when the specified client does not exist.
	ErrClientNotFound = awserrors.NewAWSError("ResourceNotFoundException", "Client not found", http.StatusBadRequest)
	// ErrClientAlreadyExists is returned when attempting to create a client that already exists.
	ErrClientAlreadyExists = awserrors.NewAWSError("ResourceExistsException", "Client already exists", http.StatusConflict)
	// ErrInvalidLambdaResponse is returned when a Lambda trigger returns a
	// malformed response, for example an unknown challenge name.
	ErrInvalidLambdaResponse = awserrors.NewAWSError("InvalidLambdaResponseException", "Invalid Lambda response", http.StatusBadRequest)
)
