package cognitoidentity

import (
	"errors"
	"net/http"
	"time"

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
	// ErrInvalidIdentityPoolConfig is returned when the identity pool has no roles configured.
	ErrInvalidIdentityPoolConfig = awserrors.NewAWSError("InvalidIdentityPoolConfigurationException", "Invalid identity pool configuration", http.StatusBadRequest)
	// ErrDeveloperUserAlreadyRegistered is returned when the input IdentityId does not
	// match the existing developer identity's IdentityId.
	ErrDeveloperUserAlreadyRegistered = awserrors.NewAWSError("DeveloperUserAlreadyRegisteredException", "This developer user identifier is already registered with Cognito", http.StatusConflict)
	// ErrLimitExceeded is returned when the total number of identity pools
	// has exceeded the per-account limit.
	ErrLimitExceeded = awserrors.NewAWSError("LimitExceededException", "Limit exceeded", http.StatusBadRequest)
)

// mapStoreError maps a store-layer error to the appropriate service-layer
// error. If the error matches the notFoundSentinel via errors.Is, it returns
// ErrResourceNotFound; otherwise it returns ErrInternalError. This replaces
// the previous pattern of unconditionally mapping all store errors to
// ErrResourceNotFound, which masked genuine internal errors.
func mapStoreError(err error, notFoundSentinel error) error {
	if errors.Is(err, notFoundSentinel) {
		return ErrResourceNotFound
	}
	return ErrInternalError
}

// CredentialResult holds temporary credentials issued for a Cognito identity.
type CredentialResult struct {
	AccessKeyID     string
	SecretAccessKey string
	SessionToken    string
	Expiration      time.Time
}

// CredentialIssuer creates temporary STS-backed sessions for the enhanced
// authflow (GetCredentialsForIdentity). Implemented by an adapter wrapping
// the STS SessionStore, injected at server startup via SetCredentialIssuer.
type CredentialIssuer interface {
	IssueSession(roleArn, roleSessionName string, durationSeconds int) (*CredentialResult, error)
}
