// Package iam provides IAM service operations for vorpalstacks.
package iam

import (
	"net/http"

	"vorpalstacks/internal/common/errors"
)

// ErrNoSuchUser is returned when a user with the specified name cannot be found.
var (
	// ErrNoSuchUser is returned when a user with the specified name cannot be found.
	ErrNoSuchUser = errors.NewAWSError("NoSuchEntity", "The user with name {UserName} cannot be found.", http.StatusNotFound)
	// ErrUserAlreadyExists is returned when attempting to create a user that already exists.
	ErrUserAlreadyExists = errors.NewAWSError("EntityAlreadyExists", "User {UserName} already exists.", http.StatusConflict)
	// ErrNoSuchAccessKey is returned when an access key with the specified ID cannot be found.
	ErrNoSuchAccessKey = errors.NewAWSError("NoSuchEntity", "The Access Key with id {AccessKeyId} cannot be found.", http.StatusNotFound)
	// ErrAccessKeyLimitExceeded is returned when the user has reached the maximum number of access keys.
	ErrAccessKeyLimitExceeded = errors.NewAWSError("LimitExceeded", "Cannot exceed quota for AccessKeysPerUser: 2.", http.StatusForbidden)
	// ErrNoSuchLoginProfile is returned when a login profile for the specified user does not exist.
	ErrNoSuchLoginProfile = errors.NewAWSError("NoSuchEntity", "Login profile for user {UserName} does not exist.", http.StatusNotFound)
	// ErrLoginProfileAlreadyExists is returned when a login profile for the user already exists.
	ErrLoginProfileAlreadyExists = errors.NewAWSError("EntityAlreadyExists", "Login profile for user {UserName} already exists.", http.StatusConflict)
	// ErrPasswordPolicyViolation is returned when the password does not meet policy requirements.
	ErrPasswordPolicyViolation = errors.NewAWSError("PasswordPolicyViolation", "The password does not meet the password policy requirements.", http.StatusBadRequest)
	// ErrDeleteConflict is returned when an entity cannot be deleted due to dependencies.
	ErrDeleteConflict = errors.NewAWSError("DeleteConflict", "Cannot delete entity, must delete access keys first.", http.StatusConflict)
	// ErrInvalidUserPath is returned when the specified user path is invalid.
	ErrInvalidUserPath = errors.NewAWSError("InvalidInput", "The specified value for path is invalid.", http.StatusBadRequest)
	// ErrInvalidInput is returned when an input parameter is invalid.
	ErrInvalidInput = errors.NewAWSError("InvalidInput", "The input parameter {Parameter} is invalid.", http.StatusBadRequest)
	// ErrLimitExceeded is returned when a quota limit has been exceeded.
	ErrLimitExceeded = errors.NewAWSError("LimitExceeded", "Cannot exceed quota for Users: 5000.", http.StatusForbidden)
	// ErrNoSuchGroup is returned when a group with the specified name cannot be found.
	ErrNoSuchGroup = errors.NewAWSError("NoSuchEntity", "The group with name {GroupName} cannot be found.", http.StatusNotFound)
	// ErrGroupAlreadyExists is returned when attempting to create a group that already exists.
	ErrGroupAlreadyExists = errors.NewAWSError("EntityAlreadyExists", "Group {GroupName} already exists.", http.StatusConflict)
	// ErrDeleteGroupConflict is returned when a group cannot be deleted due to dependencies.
	ErrDeleteGroupConflict = errors.NewAWSError("DeleteConflict", "Cannot delete entity, must remove users from group first.", http.StatusConflict)
	// ErrUserNotInGroup is returned when the specified user is not a member of the group.
	ErrUserNotInGroup = errors.NewAWSError("NoSuchEntity", "User {UserName} is not in group {GroupName}.", http.StatusNotFound)
	// ErrUserAlreadyInGroup is returned when the user is already a member of the group.
	ErrUserAlreadyInGroup = errors.NewAWSError("EntityAlreadyExists", "User {UserName} is already in group {GroupName}.", http.StatusConflict)
	// ErrNoSuchRole is returned when a role with the specified name cannot be found.
	ErrNoSuchRole = errors.NewAWSError("NoSuchEntity", "The role with name {RoleName} cannot be found.", http.StatusNotFound)
	// ErrRoleAlreadyExists is returned when attempting to create a role that already exists.
	ErrRoleAlreadyExists = errors.NewAWSError("EntityAlreadyExists", "Role {RoleName} already exists.", http.StatusConflict)
	// ErrDeleteRoleConflict is returned when a role cannot be deleted due to dependencies.
	ErrDeleteRoleConflict = errors.NewAWSError("DeleteConflict", "Cannot delete entity, must detach all policies first.", http.StatusConflict)
	// ErrNoSuchInstanceProfile is returned when an instance profile with the specified name cannot be found.
	ErrNoSuchInstanceProfile = errors.NewAWSError("NoSuchEntity", "Instance Profile {InstanceProfileName} does not exist.", http.StatusNotFound)
	// ErrInstanceProfileAlreadyExists is returned when attempting to create an instance profile that already exists.
	ErrInstanceProfileAlreadyExists = errors.NewAWSError("EntityAlreadyExists", "Instance Profile {InstanceProfileName} already exists.", http.StatusConflict)
	// ErrRoleNotInInstanceProfile is returned when the role is not associated with the instance profile.
	ErrRoleNotInInstanceProfile = errors.NewAWSError("NoSuchEntity", "Role {RoleName} not in instance profile {InstanceProfileName}.", http.StatusNotFound)
	// ErrRoleAlreadyInInstanceProfile is returned when the role is already associated with an instance profile.
	ErrRoleAlreadyInInstanceProfile = errors.NewAWSError("LimitExceeded", "Cannot exceed quota for InstanceProfilesPerRole: 1. Already associated with instance profile {InstanceProfileName}.", http.StatusForbidden)
	// ErrMalformedPolicyDocument is returned when a policy document contains invalid JSON.
	ErrMalformedPolicyDocument = errors.NewAWSError("MalformedPolicyDocument", "This policy contains invalid JSON.", http.StatusBadRequest)
	// ErrNoSuchPolicy is returned when a policy with the specified ARN cannot be found.
	ErrNoSuchPolicy = errors.NewAWSError("NoSuchEntity", "The policy with ARN {PolicyArn} does not exist.", http.StatusNotFound)
	// ErrPolicyAlreadyExists is returned when attempting to create a policy that already exists.
	ErrPolicyAlreadyExists = errors.NewAWSError("EntityAlreadyExists", "A policy with the name {PolicyName} already exists.", http.StatusConflict)
	// ErrDeletePolicyConflict is returned when a policy cannot be deleted due to attachments.
	ErrDeletePolicyConflict = errors.NewAWSError("DeleteConflict", "Cannot delete policy {PolicyArn}, there are attachments.", http.StatusConflict)
	// ErrNoSuchPolicyVersion is returned when a policy version with the specified ID cannot be found.
	ErrNoSuchPolicyVersion = errors.NewAWSError("NoSuchEntity", "Policy version {VersionId} does not exist.", http.StatusNotFound)
	// ErrLimitExceededPolicyVersions is returned when the policy has reached the maximum number of versions.
	ErrLimitExceededPolicyVersions = errors.NewAWSError("LimitExceeded", "Cannot exceed quota for PolicyVersions: 5.", http.StatusForbidden)
	// ErrNoSuchMFADevice is returned when an MFA device with the specified serial number cannot be found.
	ErrNoSuchMFADevice = errors.NewAWSError("NoSuchEntity", "MFA Device {SerialNumber} does not exist.", http.StatusNotFound)
	// ErrMFADeviceAlreadyAssigned is returned when the MFA device is already assigned to a user.
	ErrMFADeviceAlreadyAssigned = errors.NewAWSError("EntityAlreadyExists", "MFA Device {SerialNumber} is already assigned to a user.", http.StatusConflict)
	// ErrNoSuchPasswordPolicy is returned when a password policy cannot be found.
	ErrNoSuchPasswordPolicy = errors.NewAWSError("NoSuchEntity", "The Password Policy with domain name {Domain} cannot be found.", http.StatusNotFound)
	// ErrInvalidAuthenticationCode is returned when the MFA authentication code is invalid.
	ErrInvalidAuthenticationCode = errors.NewAWSError("InvalidAuthenticationCode", "Invalid authentication code.", http.StatusBadRequest)
	// ErrPolicyNotAttached is returned when the policy is not attached to the principal.
	ErrPolicyNotAttached = errors.NewAWSError("NoSuchEntity", "The policy with ARN {PolicyArn} is not attached to the principal.", http.StatusNotFound)
)

// NewNoSuchUserError creates a new error indicating that a user with the specified name cannot be found.
func NewNoSuchUserError(userName string) *errors.AWSError {
	return errors.NewNoSuchEntityException("user", userName)
}

// NewUserAlreadyExistsError creates a new error indicating that a user with the specified name already exists.
func NewUserAlreadyExistsError(userName string) *errors.AWSError {
	return errors.NewEntityAlreadyExistsException("User " + userName)
}

// NewNoSuchAccessKeyError creates a new error indicating that an access key with the specified ID cannot be found.
func NewNoSuchAccessKeyError(accessKeyId string) *errors.AWSError {
	return errors.NewNoSuchEntityException("Access Key", accessKeyId)
}

// NewNoSuchLoginProfileError creates a new error indicating that a login profile for the specified user does not exist.
func NewNoSuchLoginProfileError(userName string) *errors.AWSError {
	return errors.NewNoSuchEntityException("Login profile for user", userName)
}

// NewLoginProfileAlreadyExistsError creates a new error indicating that a login profile for the user already exists.
func NewLoginProfileAlreadyExistsError(userName string) *errors.AWSError {
	return errors.NewEntityAlreadyExistsException("Login profile for user " + userName)
}

// NewDeleteConflictError creates a new error indicating that an entity cannot be deleted due to dependencies.
func NewDeleteConflictError(message string) *errors.AWSError {
	return errors.NewDeleteConflictException(message)
}

// NewNoSuchGroupError creates a new error indicating that a group with the specified name cannot be found.
func NewNoSuchGroupError(groupName string) *errors.AWSError {
	return errors.NewNoSuchEntityException("group", groupName)
}

// NewGroupAlreadyExistsError creates a new error indicating that a group with the specified name already exists.
func NewGroupAlreadyExistsError(groupName string) *errors.AWSError {
	return errors.NewEntityAlreadyExistsException("Group " + groupName)
}

// NewDeleteGroupConflictError creates a new error indicating that a group cannot be deleted due to dependencies.
func NewDeleteGroupConflictError(message string) *errors.AWSError {
	return errors.NewDeleteConflictException(message)
}

// NewUserNotInGroupError creates a new error indicating that the user is not a member of the specified group.
func NewUserNotInGroupError(userName, groupName string) *errors.AWSError {
	return errors.NewNoSuchEntityException("User "+userName+" in group", groupName)
}

// NewUserAlreadyInGroupError creates a new error indicating that the user is already a member of the specified group.
func NewUserAlreadyInGroupError(userName, groupName string) *errors.AWSError {
	return errors.NewEntityAlreadyExistsException("User " + userName + " in group " + groupName)
}

// NewNoSuchRoleError creates a new error indicating that a role with the specified name cannot be found.
func NewNoSuchRoleError(roleName string) *errors.AWSError {
	return errors.NewNoSuchEntityException("role", roleName)
}

// NewRoleAlreadyExistsError creates a new error indicating that a role with the specified name already exists.
func NewRoleAlreadyExistsError(roleName string) *errors.AWSError {
	return errors.NewEntityAlreadyExistsException("Role " + roleName)
}

// NewDeleteRoleConflictError creates a new error indicating that a role cannot be deleted due to dependencies.
func NewDeleteRoleConflictError(message string) *errors.AWSError {
	return errors.NewDeleteConflictException(message)
}

// NewNoSuchInstanceProfileError creates a new error indicating that an instance profile with the specified name cannot be found.
func NewNoSuchInstanceProfileError(instanceProfileName string) *errors.AWSError {
	return errors.NewNoSuchEntityException("Instance Profile", instanceProfileName)
}

// NewDeleteInstanceProfileConflictError creates a new error indicating that an instance profile cannot be deleted due to attached roles.
func NewDeleteInstanceProfileConflictError(instanceProfileName string) *errors.AWSError {
	return errors.NewDeleteConflictException("Cannot delete instance profile " + instanceProfileName + ", it still has roles attached.")
}

// NewInstanceProfileAlreadyExistsError creates a new error indicating that an instance profile with the specified name already exists.
func NewInstanceProfileAlreadyExistsError(instanceProfileName string) *errors.AWSError {
	return errors.NewEntityAlreadyExistsException("Instance Profile " + instanceProfileName)
}

// NewRoleNotInInstanceProfileError creates a new error indicating that the role is not associated with the instance profile.
func NewRoleNotInInstanceProfileError(roleName, instanceProfileName string) *errors.AWSError {
	return errors.NewNoSuchEntityException("Role "+roleName+" in instance profile", instanceProfileName)
}

// NewRoleAlreadyInInstanceProfileError creates a new error indicating that the role is already associated with an instance profile.
func NewRoleAlreadyInInstanceProfileError(roleName, instanceProfileName string) *errors.AWSError {
	return errors.NewLimitExceededException("Cannot exceed quota for InstanceProfilesPerRole: 1. Already associated with instance profile " + instanceProfileName + ".")
}

// NewNoSuchPolicyError creates a new error indicating that a policy with the specified ARN cannot be found.
func NewNoSuchPolicyError(policyArn string) *errors.AWSError {
	return errors.NewNoSuchEntityException("policy", policyArn)
}

// NewPolicyAlreadyExistsError creates a new error indicating that a policy with the specified name already exists.
func NewPolicyAlreadyExistsError(policyName string) *errors.AWSError {
	return errors.NewEntityAlreadyExistsException("Policy " + policyName)
}

// NewDeletePolicyConflictError creates a new error indicating that a policy cannot be deleted due to attachments.
func NewDeletePolicyConflictError(policyArn string) *errors.AWSError {
	return errors.NewDeleteConflictException("Cannot delete policy " + policyArn + ", there are attachments.")
}

// NewNoSuchPolicyVersionError creates a new error indicating that a policy version with the specified ID cannot be found.
func NewNoSuchPolicyVersionError(versionId string) *errors.AWSError {
	return errors.NewNoSuchEntityException("Policy version", versionId)
}

// NewNoSuchMFADeviceError creates a new error indicating that an MFA device with the specified serial number cannot be found.
func NewNoSuchMFADeviceError(serialNumber string) *errors.AWSError {
	return errors.NewNoSuchEntityException("MFA Device", serialNumber)
}

// NewMFADeviceAlreadyAssignedError creates a new error indicating that the MFA device is already assigned to a user.
func NewMFADeviceAlreadyAssignedError(serialNumber string) *errors.AWSError {
	return errors.NewEntityAlreadyExistsException("MFA Device " + serialNumber)
}

// NewMFADeviceStillAssignedError creates a new error indicating that the MFA device is still assigned to a user and must be deactivated first.
func NewMFADeviceStillAssignedError(serialNumber string) *errors.AWSError {
	return errors.NewDeleteConflictException("MFA Device " + serialNumber + " is still assigned to a user. Deactivate it first.")
}

// NewNoSuchPasswordPolicyError creates a new error indicating that a password policy cannot be found.
func NewNoSuchPasswordPolicyError() *errors.AWSError {
	return errors.NewNoSuchEntityException("Password Policy", "")
}

// NewPolicyNotAttachedError creates a new error indicating that the policy is not attached to the principal.
func NewPolicyNotAttachedError(policyArn string) *errors.AWSError {
	return errors.NewNoSuchEntityException("Policy attachment for ARN", policyArn)
}

// NewNoSuchEntityError creates a new error indicating that an entity of the given type cannot be found.
func NewNoSuchEntityError(entityType, identifier string) *errors.AWSError {
	return errors.NewNoSuchEntityException(entityType, identifier)
}

// NewEntityAlreadyExistsError creates a new error indicating that an entity already exists.
func NewEntityAlreadyExistsError(entity string) *errors.AWSError {
	return errors.NewEntityAlreadyExistsException(entity)
}

// ErrValidationRequiredParameter is returned when a required parameter is missing.
var ErrValidationRequiredParameter = errors.NewAWSError("ValidationError", "Required parameter {Parameter} is missing.", http.StatusBadRequest)

// ErrNotAuthorized is returned when authentication fails.
var ErrNotAuthorized = errors.NewAWSError("NotAuthorized", "Not authorized to perform this operation.", http.StatusForbidden)

// NewValidationError creates a new error indicating that a required parameter is missing.
func NewValidationError(parameter string) *errors.AWSError {
	return errors.NewAWSError("ValidationError", "Required parameter "+parameter+" is missing.", http.StatusBadRequest)
}

// NewInvalidInputError creates a new error indicating that an input parameter is invalid.
func NewInvalidInputError(parameter string, message string) *errors.AWSError {
	return errors.NewInvalidInputException("The input parameter " + parameter + " is invalid: " + message)
}
