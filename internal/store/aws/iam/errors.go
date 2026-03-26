// Package iam provides AWS IAM store functionality for vorpalstacks.
package iam

import (
	"errors"

	"vorpalstacks/internal/store/aws/common"
)

var (
	// ErrUserNotFound is returned when the requested IAM user does not exist.
	ErrUserNotFound = errors.New("user not found")

	// ErrUserAlreadyExists is returned when attempting to create a user that
	// already exists in IAM.
	ErrUserAlreadyExists = errors.New("user already exists")

	// ErrAccessKeyNotFound is returned when the specified access key ID
	// does not exist for the user.
	ErrAccessKeyNotFound = errors.New("access key not found")

	// ErrLoginProfileNotFound is returned when the user does not have a
	// console login profile configured.
	ErrLoginProfileNotFound = errors.New("login profile not found")

	// ErrLoginProfileExists is returned when attempting to create a login
	// profile for a user who already has one.
	ErrLoginProfileExists = errors.New("login profile already exists")

	// ErrInvalidAccessKeyStatus is returned when the access key status
	// is not valid (e.g., trying to activate an inactive key).
	ErrInvalidAccessKeyStatus = errors.New("invalid access key status")

	// ErrInvalidPassword is returned when the password does not meet
	// IAM password policy requirements.
	ErrInvalidPassword = errors.New("invalid password")

	// ErrUserNotFound is returned when the specified IAM group does not exist.
	ErrGroupNotFound = errors.New("group not found")

	// ErrGroupAlreadyExists is returned when attempting to create a group
	// that already exists.
	ErrGroupAlreadyExists = errors.New("group already exists")

	// ErrUserNotInGroup is returned when attempting to remove a user from
	// a group they are not a member of.
	ErrUserNotInGroup = errors.New("user not in group")

	// ErrUserAlreadyInGroup is returned when attempting to add a user to
	// a group they are already a member of.
	ErrUserAlreadyInGroup = errors.New("user already in group")

	// ErrRoleNotFound is returned when the specified IAM role does not exist.
	ErrRoleNotFound = errors.New("role not found")

	// ErrRoleAlreadyExists is returned when attempting to create a role
	// that already exists.
	ErrRoleAlreadyExists = errors.New("role already exists")

	// ErrInstanceProfileNotFound is returned when the specified instance
	// profile does not exist.
	ErrInstanceProfileNotFound = errors.New("instance profile not found")

	// ErrInstanceProfileAlreadyExists is returned when attempting to create
	// an instance profile that already exists.
	ErrInstanceProfileAlreadyExists = errors.New("instance profile already exists")

	// ErrRoleNotInInstanceProfile is returned when attempting to remove a
	// role from an instance profile it is not attached to.
	ErrRoleNotInInstanceProfile = errors.New("role not in instance profile")

	// ErrRoleAlreadyInInstanceProfile is returned when attempting to attach
	// a role that is already attached to the instance profile.
	ErrRoleAlreadyInInstanceProfile = errors.New("role already in instance profile")

	// ErrPolicyNotFound is returned when the specified IAM policy does not exist.
	ErrPolicyNotFound = errors.New("policy not found")

	// ErrPolicyAlreadyExists is returned when attempting to create a policy
	// that already exists.
	ErrPolicyAlreadyExists = errors.New("policy already exists")

	// ErrMFADeviceNotFound is returned when the specified MFA device
	// does not exist for the user.
	ErrMFADeviceNotFound = errors.New("mfa device not found")

	// ErrPasswordPolicyNotFound is returned when the account does not have
	// a password policy configured.
	ErrPasswordPolicyNotFound = errors.New("password policy not found")

	// ErrServerCertificateNotFound is returned when the specified server
	// certificate does not exist.
	ErrServerCertificateNotFound = errors.New("server certificate not found")

	// ErrServerCertificateAlreadyExists is returned when attempting to create
	// a server certificate that already exists.
	ErrServerCertificateAlreadyExists = errors.New("server certificate already exists")

	// ErrSAMLProviderNotFound is returned when the specified SAML provider
	// does not exist.
	ErrSAMLProviderNotFound = errors.New("saml provider not found")

	// ErrSAMLProviderAlreadyExists is returned when attempting to create a
	// SAML provider that already exists.
	ErrSAMLProviderAlreadyExists = errors.New("saml provider already exists")

	// ErrOpenIDConnectProviderNotFound is returned when the specified OIDC
	// provider does not exist.
	ErrOpenIDConnectProviderNotFound = errors.New("openid connect provider not found")

	// ErrOpenIDConnectProviderAlreadyExists is returned when attempting to
	// create an OIDC provider that already exists.
	ErrOpenIDConnectProviderAlreadyExists = errors.New("openid connect provider already exists")

	// ErrSigningCertificateNotFound is returned when the specified signing
	// certificate does not exist for the user.
	ErrSigningCertificateNotFound = errors.New("signing certificate not found")

	// ErrSSHPublicKeyNotFound is returned when the specified SSH public key
	// does not exist for the user.
	ErrSSHPublicKeyNotFound = errors.New("ssh public key not found")

	// ErrServiceSpecificCredentialNotFound is returned when the specified
	// service-specific credential does not exist.
	ErrServiceSpecificCredentialNotFound = errors.New("service-specific credential not found")
)

// StoreError represents an IAM store operation error, extending the common
// store error type with IAM-specific context.
type StoreError = common.StoreError

// NewStoreError creates a new IAM store error with the given operation and error.
// This is a convenience constructor that wraps the common store error functionality.
func NewStoreError(op string, err error) *StoreError {
	return common.NewStoreError("iam", op, err)
}

// NewStoreErrorWithKey creates a new IAM store error with the given operation,
// key identifier, and error. This is used when the error relates to a specific
// IAM resource identifier.
func NewStoreErrorWithKey(op, key string, err error) *StoreError {
	return common.NewStoreErrorWithKey("iam", op, key, err)
}

// IsNotFound checks if the given error is a "not found" type error for any
// IAM resource type. This includes user, access key, group, role, policy,
// and instance profile not found errors.
func IsNotFound(err error) bool {
	return common.IsNotFound(err) ||
		errors.Is(err, ErrUserNotFound) ||
		errors.Is(err, ErrAccessKeyNotFound) ||
		errors.Is(err, ErrGroupNotFound) ||
		errors.Is(err, ErrRoleNotFound) ||
		errors.Is(err, ErrPolicyNotFound) ||
		errors.Is(err, ErrInstanceProfileNotFound) ||
		errors.Is(err, ErrMFADeviceNotFound) ||
		errors.Is(err, ErrServerCertificateNotFound) ||
		errors.Is(err, ErrSAMLProviderNotFound) ||
		errors.Is(err, ErrOpenIDConnectProviderNotFound) ||
		errors.Is(err, ErrSigningCertificateNotFound) ||
		errors.Is(err, ErrSSHPublicKeyNotFound) ||
		errors.Is(err, ErrServiceSpecificCredentialNotFound)
}
