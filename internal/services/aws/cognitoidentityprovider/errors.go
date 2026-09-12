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
	// ErrWebAuthnNotEnabled is returned when a WebAuthn registration is
	// started or completed on a user pool whose passkey feature is not
	// enabled.
	ErrWebAuthnNotEnabled = awserrors.NewAWSError("WebAuthnNotEnabledException", "The user pool does not have the passkey feature enabled", http.StatusBadRequest)
	// ErrWebAuthnConfigurationMissing is returned when a WebAuthn
	// registration must run under a relying party id but the user pool
	// carries neither a configured relying party id nor a hosted domain.
	ErrWebAuthnConfigurationMissing = awserrors.NewAWSError("WebAuthnConfigurationMissingException", "The user pool has no configured relying party ID or user pool domain", http.StatusBadRequest)
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
	// ErrUserNotFound is returned when the specified user does not exist.
	ErrUserNotFound = awserrors.NewAWSError("UserNotFoundException", "User not found", http.StatusBadRequest)
	// ErrUserAlreadyExists is returned when attempting to create a user that already exists.
	ErrUserAlreadyExists = awserrors.NewAWSError("UsernameExistsException", "User already exists", http.StatusBadRequest)
	// ErrAliasExists is returned when a write would claim an alias attribute
	// value (a verified email, phone number or preferred_username configured
	// as an alias or username attribute) that another user already claims.
	ErrAliasExists = awserrors.NewAWSError("AliasExistsException", "An account with the given alias already exists", http.StatusBadRequest)
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
	// ErrPasswordHistoryViolation is returned when a new password matches one
	// of the previous passwords the pool's PasswordHistorySize remembers.
	ErrPasswordHistoryViolation = awserrors.NewAWSError("PasswordHistoryPolicyViolationException", "Password cannot be the same as a previous password", http.StatusBadRequest)
	// ErrUserNotConfirmed is returned when the user is not confirmed.
	ErrUserNotConfirmed = awserrors.NewAWSError("UserNotConfirmedException", "User is not confirmed", http.StatusBadRequest)
	// ErrUserAlreadyConfirmed is returned when the user is already confirmed.
	ErrUserAlreadyConfirmed = awserrors.NewAWSError("UserAlreadyConfirmedException", "User already confirmed", http.StatusBadRequest)
	// ErrCodeMismatch is returned when the verification code does not match.
	ErrCodeMismatch = awserrors.NewAWSError("CodeMismatchException", "Invalid verification code", http.StatusBadRequest)
	// ErrDeviceKeyExists is returned when ConfirmDevice names a device key
	// the user has already registered.
	ErrDeviceKeyExists = awserrors.NewAWSError("DeviceKeyExistsException", "Device key already exists", http.StatusBadRequest)
	// ErrExpiredCode is returned when the verification code has expired.
	ErrExpiredCode = awserrors.NewAWSError("ExpiredCodeException", "Invalid code provided, please request a code again", http.StatusBadRequest)
	// ErrInternalError is returned when an internal error occurs.
	ErrInternalError = awserrors.NewAWSError("InternalErrorException", "Internal error", http.StatusInternalServerError)
	// ErrClientNotFound is returned when the specified client does not exist.
	ErrClientNotFound = awserrors.NewAWSError("ResourceNotFoundException", "Client not found", http.StatusBadRequest)
	// ErrOperationNotEnabled is returned when an operation is not available
	// for the current user pool or app client configuration, for example a
	// machine-to-machine token request against an app client that lacks a
	// client secret or the ALLOW_CLIENT_TOKEN_AUTH authentication flow.
	ErrOperationNotEnabled = awserrors.NewAWSError("OperationNotEnabledException", "Operation not enabled", http.StatusBadRequest)
	// ErrInvalidLambdaResponse is returned when a Lambda trigger returns a
	// malformed response, for example an unknown challenge name.
	ErrInvalidLambdaResponse = awserrors.NewAWSError("InvalidLambdaResponseException", "Invalid Lambda response", http.StatusBadRequest)
	// ErrUserLambdaValidation is returned when an authentication-path Lambda
	// trigger (Pre/PostAuthentication, Define/CreateAuthChallenge) raises an
	// error and thereby fails the authentication.
	ErrUserLambdaValidation = awserrors.NewAWSError("UserLambdaValidationException", "User validation with the Lambda service failed", http.StatusBadRequest)
	// ErrPasswordResetRequired is returned when a user whose password a
	// password reset has deactivated tries to sign in: the model documents
	// that "your client app must handle the PasswordResetRequiredException
	// during the authentication flow" and complete the forgot-password flow.
	ErrPasswordResetRequired = awserrors.NewAWSError("PasswordResetRequiredException", "Password reset required for the user", http.StatusBadRequest)
)
