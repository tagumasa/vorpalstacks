package iam

import (
	"errors"
	"fmt"
	"net/http"

	awserrors "vorpalstacks/internal/common/errors"
	storecommon "vorpalstacks/internal/store/aws/common"
	iamstore "vorpalstacks/internal/store/aws/iam"
)

// Fixed-message AWS error values. Every parameterised error class is built
// by the New* constructors below, which interpolate the entity name into
// the message; pre-baked sentinels carrying placeholder text are not
// permitted (their unsubstituted "{Placeholder}" text would reach the
// wire verbatim).
var (
	// ErrAccessKeyLimitExceeded is returned when the user has reached the maximum number of access keys.
	ErrAccessKeyLimitExceeded = awserrors.NewAWSError("LimitExceeded", fmt.Sprintf("Cannot exceed quota for AccessKeysPerUser: %d.", iamstore.MaxAccessKeysPerUser), http.StatusConflict)
	// ErrLimitExceededGroupsPerUser is returned when the user is already a
	// member of the maximum number of IAM groups.
	ErrLimitExceededGroupsPerUser = awserrors.NewAWSError("LimitExceeded", fmt.Sprintf("Cannot exceed quota for GroupsPerUser: %d.", iamstore.MaxIAMGroupsPerUser), http.StatusConflict)
	// ErrLimitExceededMFADevicesPerUser is returned when the user already
	// has the maximum number of MFA devices configured.
	ErrLimitExceededMFADevicesPerUser = awserrors.NewAWSError("LimitExceeded", fmt.Sprintf("Cannot exceed quota for MFADevicesPerUser: %d.", iamstore.MaxMFADevicesPerUser), http.StatusConflict)
	// ErrPasswordPolicyViolation is returned when the password does not meet policy requirements.
	ErrPasswordPolicyViolation = awserrors.NewAWSError("PasswordPolicyViolation", "The password does not meet the password policy requirements.", http.StatusBadRequest)
	// ErrInstanceProfileRoleLimit is returned when attempting to add a second
	// role to an instance profile (AWS allows only one role per profile).
	ErrInstanceProfileRoleLimit = awserrors.NewAWSError("LimitExceeded", "Cannot exceed quota for RolesPerInstanceProfile: 1.", http.StatusConflict)
	// ErrMalformedPolicyDocument is returned when a policy document contains invalid JSON.
	ErrMalformedPolicyDocument = awserrors.NewAWSError("MalformedPolicyDocument", "This policy contains invalid JSON.", http.StatusBadRequest)
	// ErrReportExpired reports that the most recent credential report is
	// older than the four-hour validity window; a new one must be
	// generated. The fault and its code are the operation's modelled
	// CredentialReportExpiredException.
	ErrReportExpired = awserrors.NewAWSError("ReportExpired",
		"The request was rejected because the most recent credential report has expired. To generate a new credential report, use GenerateCredentialReport.",
		http.StatusGone)
	// ErrMalformedCertificate is returned when a certificate body cannot be
	// parsed as X.509.
	ErrMalformedCertificate = awserrors.NewAWSError("MalformedCertificate", "The certificate is malformed or invalid.", http.StatusBadRequest)

	// ErrDuplicateCertificate is returned when the certificate is already
	// registered for the user.
	ErrDuplicateCertificate = awserrors.NewAWSError("DuplicateCertificate", "The certificate is already registered for this user.", http.StatusConflict)

	// ErrLimitExceededSigningCertificates is returned when the user already
	// holds the maximum number of signing certificates.
	ErrLimitExceededSigningCertificates = awserrors.NewAWSError("LimitExceeded", fmt.Sprintf("Cannot exceed quota for SigningCertificatesPerUser: %d.", iamstore.MaxSigningCertificatesPerUser), http.StatusConflict)

	// ErrInvalidPublicKey is returned when an SSH public key body cannot be
	// parsed.
	ErrInvalidPublicKey = awserrors.NewAWSError("InvalidPublicKey", "The public key is invalid or unsupported.", http.StatusBadRequest)

	// ErrNotSupportedService is returned when the named service does not
	// support service-specific credentials.
	ErrNotSupportedService = awserrors.NewAWSError("NotSupportedService", "The specified service does not support service-specific credentials.", http.StatusNotFound)

	// ErrFeatureEnabled is returned when enabling outbound web identity
	// federation while it is already enabled.
	ErrFeatureEnabled = awserrors.NewAWSError("FeatureEnabled", "Outbound web identity federation is already enabled for this account.", http.StatusConflict)

	// ErrFeatureDisabled is returned when accessing outbound web identity
	// federation configuration while the feature is disabled.
	ErrFeatureDisabled = awserrors.NewAWSError("FeatureDisabled", "Outbound web identity federation is not enabled for this account.", http.StatusNotFound)

	// ErrKeyPairMismatch is returned when the private key does not match
	// the certificate's public key.
	ErrKeyPairMismatch = awserrors.NewAWSError("KeyPairMismatch", "The private key does not match the public key in the certificate.", http.StatusBadRequest)

	// ErrDuplicateSSHPublicKey is returned when the key material is already
	// registered for the user.
	ErrDuplicateSSHPublicKey = awserrors.NewAWSError("DuplicateSSHPublicKey", "The SSH public key is already associated with this user.", http.StatusBadRequest)

	// ErrLimitExceededSSHPublicKeys is returned when the user already holds
	// the maximum number of SSH public keys.
	ErrLimitExceededSSHPublicKeys = awserrors.NewAWSError("LimitExceeded", fmt.Sprintf("Cannot exceed quota for SSHPublicKeysPerUser: %d.", iamstore.MaxSSHPublicKeysPerUser), http.StatusConflict)

	// ErrLimitExceededPolicyVersions is returned when the policy has reached the maximum number of versions.
	ErrLimitExceededPolicyVersions = awserrors.NewAWSError("LimitExceeded", fmt.Sprintf("Cannot exceed quota for PolicyVersions: %d.", iamstore.MaxPolicyVersions), http.StatusConflict)
	// ErrInvalidAuthenticationCode is returned when the MFA authentication code is invalid.
	ErrInvalidAuthenticationCode = awserrors.NewAWSError("InvalidAuthenticationCode", "Invalid authentication code.", http.StatusBadRequest)
)

// NewNoSuchUserError creates a new error indicating that a user with the specified name cannot be found.
func NewNoSuchUserError(userName string) *awserrors.AWSError {
	return awserrors.NewNoSuchEntityException("user", userName)
}

// NewUserAlreadyExistsError creates a new error indicating that a user with the specified name already exists.
func NewUserAlreadyExistsError(userName string) *awserrors.AWSError {
	return awserrors.NewEntityAlreadyExistsException("User " + userName)
}

// NewNoSuchAccessKeyError creates a new error indicating that an access key with the specified ID cannot be found.
func NewNoSuchAccessKeyError(accessKeyId string) *awserrors.AWSError {
	return awserrors.NewNoSuchEntityException("Access Key", accessKeyId)
}

// NewNoSuchLoginProfileError creates a new error indicating that a login profile for the specified user does not exist.
func NewNoSuchLoginProfileError(userName string) *awserrors.AWSError {
	return awserrors.NewNoSuchEntityException("Login profile for user", userName)
}

// NewLoginProfileAlreadyExistsError creates a new error indicating that a login profile for the user already exists.
func NewLoginProfileAlreadyExistsError(userName string) *awserrors.AWSError {
	return awserrors.NewEntityAlreadyExistsException("Login profile for user " + userName)
}

// NewDeleteConflictError creates a new error indicating that an entity cannot be deleted due to dependencies.
func NewDeleteConflictError(message string) *awserrors.AWSError {
	return awserrors.NewDeleteConflictException(message)
}

// NewNoSuchGroupError creates a new error indicating that a group with the specified name cannot be found.
func NewNoSuchGroupError(groupName string) *awserrors.AWSError {
	return awserrors.NewNoSuchEntityException("group", groupName)
}

// NewGroupAlreadyExistsError creates a new error indicating that a group with the specified name already exists.
func NewGroupAlreadyExistsError(groupName string) *awserrors.AWSError {
	return awserrors.NewEntityAlreadyExistsException("Group " + groupName)
}

// NewDeleteGroupConflictError creates a new error indicating that a group cannot be deleted due to dependencies.
func NewDeleteGroupConflictError(message string) *awserrors.AWSError {
	return awserrors.NewDeleteConflictException(message)
}

// NewUserNotInGroupError creates a new error indicating that the user is not a member of the specified group.
func NewUserNotInGroupError(userName, groupName string) *awserrors.AWSError {
	return awserrors.NewNoSuchEntityException("User "+userName+" in group", groupName)
}

// NewNoSuchRoleError creates a new error indicating that a role with the specified name cannot be found.
func NewNoSuchRoleError(roleName string) *awserrors.AWSError {
	return awserrors.NewNoSuchEntityException("role", roleName)
}

// NewNoSuchRoleTemplateError creates a new error indicating that a role
// template with the specified ARN cannot be found in the catalogue.
func NewNoSuchRoleTemplateError(templateArn string) *awserrors.AWSError {
	return awserrors.NewNoSuchEntityException("Role template", templateArn)
}

// NewRoleTemplateDisabledError creates a new error indicating that the role
// template (or the addressed template version) is disabled and cannot be
// used to create roles.
func NewRoleTemplateDisabledError(templateArn string) *awserrors.AWSError {
	return awserrors.NewAWSError("RoleTemplateDisabled", "The role template "+templateArn+" is disabled and cannot be used to create new roles.", http.StatusBadRequest)
}

// NewNameConflictError creates a new error indicating that the role name a
// template resolves to conflicts with an existing role in the account.
func NewNameConflictError(roleName string) *awserrors.AWSError {
	return awserrors.NewAWSError("NameConflict", "The request was rejected because the resulting role name "+roleName+" conflicts with an existing role in the account.", http.StatusConflict)
}

// NewRoleAlreadyExistsError creates a new error indicating that a role with the specified name already exists.
func NewRoleAlreadyExistsError(roleName string) *awserrors.AWSError {
	return awserrors.NewEntityAlreadyExistsException("Role " + roleName)
}

// NewDeleteRoleConflictError creates a new error indicating that a role cannot be deleted due to dependencies.
func NewDeleteRoleConflictError(message string) *awserrors.AWSError {
	return awserrors.NewDeleteConflictException(message)
}

// NewNoSuchInstanceProfileError creates a new error indicating that an instance profile with the specified name cannot be found.
func NewNoSuchInstanceProfileError(instanceProfileName string) *awserrors.AWSError {
	return awserrors.NewNoSuchEntityException("Instance Profile", instanceProfileName)
}

// NewDeleteInstanceProfileConflictError creates a new error indicating that an instance profile cannot be deleted due to attached roles.
func NewDeleteInstanceProfileConflictError(instanceProfileName string) *awserrors.AWSError {
	return awserrors.NewDeleteConflictException("Cannot delete instance profile " + instanceProfileName + ", it still has roles attached.")
}

// NewInstanceProfileAlreadyExistsError creates a new error indicating that an instance profile with the specified name already exists.
func NewInstanceProfileAlreadyExistsError(instanceProfileName string) *awserrors.AWSError {
	return awserrors.NewEntityAlreadyExistsException("Instance Profile " + instanceProfileName)
}

// NewRoleNotInInstanceProfileError creates a new error indicating that the role is not associated with the instance profile.
func NewRoleNotInInstanceProfileError(roleName, instanceProfileName string) *awserrors.AWSError {
	return awserrors.NewNoSuchEntityException("Role "+roleName+" in instance profile", instanceProfileName)
}

// NewRoleAlreadyInInstanceProfileError creates a new error indicating that the role is already associated with an instance profile.
func NewRoleAlreadyInInstanceProfileError(roleName, instanceProfileName string) *awserrors.AWSError {
	return awserrors.NewLimitExceededException("Cannot exceed quota for InstanceProfilesPerRole: 1. Already associated with instance profile " + instanceProfileName + ".")
}

// NewNoSuchPolicyError creates a new error indicating that a policy with the specified ARN cannot be found.
func NewNoSuchPolicyError(policyArn string) *awserrors.AWSError {
	return awserrors.NewNoSuchEntityException("policy", policyArn)
}

// NewPolicyAlreadyExistsError creates a new error indicating that a policy with the specified name already exists.
func NewPolicyAlreadyExistsError(policyName string) *awserrors.AWSError {
	return awserrors.NewEntityAlreadyExistsException("Policy " + policyName)
}

// NewDeletePolicyConflictError creates a new error indicating that a policy cannot be deleted due to attachments.
func NewDeletePolicyConflictError(policyArn string) *awserrors.AWSError {
	return awserrors.NewDeleteConflictException("Cannot delete policy " + policyArn + ", there are attachments.")
}

// NewNoSuchPolicyVersionError creates a new error indicating that a policy version with the specified ID cannot be found.
func NewNoSuchPolicyVersionError(versionId string) *awserrors.AWSError {
	return awserrors.NewNoSuchEntityException("Policy version", versionId)
}

// NewNoSuchMFADeviceError creates a new error indicating that an MFA device with the specified serial number cannot be found.
func NewNoSuchMFADeviceError(serialNumber string) *awserrors.AWSError {
	return awserrors.NewNoSuchEntityException("MFA Device", serialNumber)
}

// NewMFADeviceAlreadyAssignedError creates a new error indicating that the MFA device is already assigned to a user.
func NewMFADeviceAlreadyAssignedError(serialNumber string) *awserrors.AWSError {
	return awserrors.NewEntityAlreadyExistsException("MFA Device " + serialNumber)
}

// NewMFADeviceStillAssignedError creates a new error indicating that the MFA device is still assigned to a user and must be deactivated first.
func NewMFADeviceStillAssignedError(serialNumber string) *awserrors.AWSError {
	return awserrors.NewDeleteConflictException("MFA Device " + serialNumber + " is still assigned to a user. Deactivate it first.")
}

// NewNoSuchPasswordPolicyError creates a new error indicating that a password policy cannot be found.
func NewNoSuchPasswordPolicyError() *awserrors.AWSError {
	return awserrors.NewNoSuchEntityException("Password Policy", "")
}

// NewPolicyNotAttachedError creates a new error indicating that the policy is not attached to the principal.
func NewPolicyNotAttachedError(policyArn string) *awserrors.AWSError {
	return awserrors.NewNoSuchEntityException("Policy attachment for ARN", policyArn)
}

// NewNoSuchEntityError creates a new error indicating that an entity of the given type cannot be found.
func NewNoSuchEntityError(entityType, identifier string) *awserrors.AWSError {
	return awserrors.NewNoSuchEntityException(entityType, identifier)
}

// NewEntityAlreadyExistsError creates a new error indicating that an entity already exists.
func NewEntityAlreadyExistsError(entity string) *awserrors.AWSError {
	return awserrors.NewEntityAlreadyExistsException(entity)
}

// ErrNotAuthorized is returned when authentication fails.
var ErrNotAuthorized = awserrors.NewAWSError("NotAuthorized", "Not authorized to perform this operation.", http.StatusForbidden)

// NewValidationError creates a new error indicating that a required parameter is missing.
// Uses the InvalidInput wire code per Smithy InvalidInputException.
func NewValidationError(parameter string) *awserrors.AWSError {
	return awserrors.NewAWSError("InvalidInput", "Required parameter "+parameter+" is missing.", http.StatusBadRequest)
}

// NewInvalidInputError creates a new error indicating that an input parameter is invalid.
func NewInvalidInputError(parameter string, message string) *awserrors.AWSError {
	return awserrors.NewInvalidInputException("The input parameter " + parameter + " is invalid: " + message)
}

// storeReadError maps a failed store read to its AWS wire form: notFound
// when the store signalled that the entity is absent (the common
// ErrNotFound sentinel from the generic entity stores, or the IAM family
// sentinel from the specialised stores), and an InternalFailure carrying
// the cause otherwise — an infrastructure fault must not surface as
// NoSuchEntity. notFoundSentinel may be nil for stores that signal
// absence only through the common sentinel.
func storeReadError(err, notFoundSentinel error, notFound *awserrors.AWSError) *awserrors.AWSError {
	if storecommon.IsNotFound(err) || errors.Is(err, notFoundSentinel) {
		return notFound
	}
	return awserrors.NewInternalFailureException("Failed to read from the IAM store: " + err.Error())
}

// storeListError maps the failure of an unkeyed store walk to a 5xx: a
// listing, or the read of an item that the walk itself returned. Unlike a
// keyed read there is no not-found outcome for the caller — every failure is
// an infrastructure fault, or a store inconsistency such as an item that
// vanished between the listing and its read — so the cause is carried in the
// message instead of being mapped to the NoSuchEntity vocabulary.
func storeListError(err error) *awserrors.AWSError {
	return awserrors.NewInternalFailureException("Failed to read from the IAM store: " + err.Error())
}

// NewServiceFailureException builds the ServiceFailure error shape the
// credential-report operations model for their server faults: an unknown
// processing failure on GenerateCredentialReport or GetCredentialReport
// carries this code, not the platform-wide InternalFailure vocabulary.
func NewServiceFailureException(message string) *awserrors.AWSError {
	return awserrors.NewAWSError("ServiceFailure", message, http.StatusInternalServerError)
}
