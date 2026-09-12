package cognitoidentityprovider

import (
	"errors"

	"vorpalstacks/internal/store/aws/common"
)

var (
	// ErrUserPoolNotFound is returned when the specified Cognito user pool
	// does not exist.
	ErrUserPoolNotFound = errors.New("user pool not found")

	// ErrUserPoolAlreadyExists is returned when attempting to create a user pool
	// that already exists.
	ErrUserPoolAlreadyExists = errors.New("user pool already exists")

	// ErrInvalidUserPoolName is returned when the user pool name is not valid.
	ErrInvalidUserPoolName = errors.New("invalid user pool name")

	// ErrUserNotFound is returned when the specified user does not exist.
	ErrUserNotFound = errors.New("user not found")

	// ErrUserAlreadyExists is returned when attempting to create a user
	// that already exists.
	ErrUserAlreadyExists = errors.New("user already exists")

	// ErrAliasExists is returned when a write would claim an alias
	// attribute value that another user already claims.
	ErrAliasExists = errors.New("alias already exists")

	// ErrInvalidUsername is returned when the username is not valid.
	ErrInvalidUsername = errors.New("invalid username")

	// ErrGroupNotFound is returned when the specified group does not exist.
	ErrGroupNotFound = errors.New("group not found")

	// ErrGroupAlreadyExists is returned when attempting to create a group
	// that already exists.
	ErrGroupAlreadyExists = errors.New("group already exists")

	// ErrInvalidGroupName is returned when the group name is not valid.
	ErrInvalidGroupName = errors.New("invalid group name")

	// ErrUserNotInGroup is returned when attempting to remove a user from
	// a group they are not a member of.
	ErrUserNotInGroup = errors.New("user not in group")

	// ErrUserAlreadyInGroup is returned when attempting to add a user to
	// a group they are already a member of.
	ErrUserAlreadyInGroup = errors.New("user already in group")

	// ErrTokenNotFound is returned when the specified token does not exist.
	ErrTokenNotFound = errors.New("token not found")

	// ErrTokenExpired is returned when the token has expired.
	ErrTokenExpired = errors.New("token expired")

	// ErrUserNotConfirmed is returned when the user has not been confirmed.
	ErrUserNotConfirmed = errors.New("user not confirmed")

	// ErrUserAlreadyConfirmed is returned when the user is already confirmed.
	ErrUserAlreadyConfirmed = errors.New("user already confirmed")

	// ErrIncorrectPassword is returned when the password is incorrect.
	ErrIncorrectPassword = errors.New("incorrect password")

	// ErrPasswordPolicyViolation is returned when the password does not meet
	// the password policy requirements.
	ErrPasswordPolicyViolation = errors.New("password policy violation")

	// ErrClientNotFound is returned when the specified client does not exist.
	ErrClientNotFound = errors.New("client not found")

	// ErrInvalidParameter is returned when a parameter is not valid.
	ErrInvalidParameter = common.ErrInvalidParameter

	// ErrResourceAlreadyExists is returned when a resource already exists.
	ErrResourceAlreadyExists = errors.New("resource already exists")

	// ErrNotFound is returned when a generic resource is not found.
	ErrNotFound = errors.New("not found")

	// ErrManagedLoginBrandingExists is returned when a save would assign an
	// app client a second branding style.
	ErrManagedLoginBrandingExists = errors.New("app client already has an assigned managed login branding style")

	// ErrTermsExists is returned when a save would give an app client a
	// second terms document of the same name.
	ErrTermsExists = errors.New("terms document name already held by the app client")

	// ErrImportJobStatusConflict is returned when a conditional import job
	// status transition finds the job in a different state than expected.
	ErrImportJobStatusConflict = errors.New("user import job status conflict")

	// ErrImportJobActiveExists is returned when starting an import job
	// while another job is still active in the account.
	ErrImportJobActiveExists = errors.New("another user import job is active")

	// ErrUserPoolDomainNotFound is returned when the specified user pool
	// domain does not exist, or when it exists but is bound to a different
	// pool than the one named in the request.
	ErrUserPoolDomainNotFound = errors.New("user pool domain not found")

	// ErrUserPoolDomainInUse is returned when attempting to bind a domain
	// that is already assigned to another user pool.
	ErrUserPoolDomainInUse = errors.New("user pool domain already in use")

	// ErrUserPoolAlreadyHasDomain is returned when attempting to bind a
	// second domain to a user pool that already owns one.
	ErrUserPoolAlreadyHasDomain = errors.New("user pool already has a domain")

	// ErrResourceServerNotFound is returned when the specified resource
	// server does not exist. It is a distinct sentinel from a missing pool
	// so callers can tell the two entities apart.
	ErrResourceServerNotFound = errors.New("resource server not found")

	// ErrIdentityProviderNotFound is returned when the specified identity
	// provider does not exist. It is a distinct sentinel from a missing pool
	// so callers can tell the two entities apart.
	ErrIdentityProviderNotFound = errors.New("identity provider not found")
)
