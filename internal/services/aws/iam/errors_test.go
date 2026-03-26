package iam

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIAMErrors(t *testing.T) {
	t.Run("predefined errors", func(t *testing.T) {
		assert.Equal(t, "NoSuchEntity: The user with name {UserName} cannot be found.", ErrNoSuchUser.Error())
		assert.Equal(t, http.StatusNotFound, ErrNoSuchUser.GetHTTPStatusCode())

		assert.Equal(t, "EntityAlreadyExists: User {UserName} already exists.", ErrUserAlreadyExists.Error())
		assert.Equal(t, http.StatusConflict, ErrUserAlreadyExists.GetHTTPStatusCode())

		assert.Equal(t, "LimitExceeded: Cannot exceed quota for AccessKeysPerUser: 2.", ErrAccessKeyLimitExceeded.Error())
		assert.Equal(t, http.StatusForbidden, ErrAccessKeyLimitExceeded.GetHTTPStatusCode())

		assert.Equal(t, "PasswordPolicyViolation: The password does not meet the password policy requirements.", ErrPasswordPolicyViolation.Error())
		assert.Equal(t, http.StatusBadRequest, ErrPasswordPolicyViolation.GetHTTPStatusCode())

		assert.Equal(t, "DeleteConflict: Cannot delete entity, must delete access keys first.", ErrDeleteConflict.Error())
		assert.Equal(t, http.StatusConflict, ErrDeleteConflict.GetHTTPStatusCode())

		assert.Equal(t, "LimitExceeded: Cannot exceed quota for Users: 5000.", ErrLimitExceeded.Error())
		assert.Equal(t, http.StatusForbidden, ErrLimitExceeded.GetHTTPStatusCode())

		assert.Equal(t, "MalformedPolicyDocument: This policy contains invalid JSON.", ErrMalformedPolicyDocument.Error())
		assert.Equal(t, http.StatusBadRequest, ErrMalformedPolicyDocument.GetHTTPStatusCode())

		assert.Equal(t, "LimitExceeded: Cannot exceed quota for PolicyVersions: 5.", ErrLimitExceededPolicyVersions.Error())
		assert.Equal(t, http.StatusForbidden, ErrLimitExceededPolicyVersions.GetHTTPStatusCode())

		assert.Equal(t, "InvalidAuthenticationCode: Invalid authentication code.", ErrInvalidAuthenticationCode.Error())
		assert.Equal(t, http.StatusBadRequest, ErrInvalidAuthenticationCode.GetHTTPStatusCode())
	})

	t.Run("NewNoSuchUserError", func(t *testing.T) {
		err := NewNoSuchUserError("testuser")
		assert.Equal(t, "NoSuchEntity: The user with name testuser cannot be found.", err.Error())
		assert.Equal(t, http.StatusNotFound, err.GetHTTPStatusCode())
	})

	t.Run("NewUserAlreadyExistsError", func(t *testing.T) {
		err := NewUserAlreadyExistsError("testuser")
		assert.Equal(t, "EntityAlreadyExists: User testuser already exists", err.Error())
		assert.Equal(t, http.StatusConflict, err.GetHTTPStatusCode())
	})

	t.Run("NewNoSuchAccessKeyError", func(t *testing.T) {
		err := NewNoSuchAccessKeyError("AKIAIOSFODNN7EXAMPLE")
		assert.Equal(t, "NoSuchEntity: The Access Key with name AKIAIOSFODNN7EXAMPLE cannot be found.", err.Error())
		assert.Equal(t, http.StatusNotFound, err.GetHTTPStatusCode())
	})

	t.Run("NewNoSuchLoginProfileError", func(t *testing.T) {
		err := NewNoSuchLoginProfileError("testuser")
		assert.Equal(t, "NoSuchEntity: The Login profile for user with name testuser cannot be found.", err.Error())
		assert.Equal(t, http.StatusNotFound, err.GetHTTPStatusCode())
	})

	t.Run("NewLoginProfileAlreadyExistsError", func(t *testing.T) {
		err := NewLoginProfileAlreadyExistsError("testuser")
		assert.Equal(t, "EntityAlreadyExists: Login profile for user testuser already exists", err.Error())
		assert.Equal(t, http.StatusConflict, err.GetHTTPStatusCode())
	})

	t.Run("NewDeleteConflictError", func(t *testing.T) {
		err := NewDeleteConflictError("must delete access keys first")
		assert.Equal(t, "DeleteConflict: must delete access keys first", err.Error())
		assert.Equal(t, http.StatusConflict, err.GetHTTPStatusCode())
	})

	t.Run("NewNoSuchGroupError", func(t *testing.T) {
		err := NewNoSuchGroupError("testgroup")
		assert.Equal(t, "NoSuchEntity: The group with name testgroup cannot be found.", err.Error())
		assert.Equal(t, http.StatusNotFound, err.GetHTTPStatusCode())
	})

	t.Run("NewGroupAlreadyExistsError", func(t *testing.T) {
		err := NewGroupAlreadyExistsError("testgroup")
		assert.Equal(t, "EntityAlreadyExists: Group testgroup already exists", err.Error())
		assert.Equal(t, http.StatusConflict, err.GetHTTPStatusCode())
	})

	t.Run("NewUserNotInGroupError", func(t *testing.T) {
		err := NewUserNotInGroupError("testuser", "testgroup")
		assert.Equal(t, "NoSuchEntity: The User testuser in group with name testgroup cannot be found.", err.Error())
		assert.Equal(t, http.StatusNotFound, err.GetHTTPStatusCode())
	})

	t.Run("NewUserAlreadyInGroupError", func(t *testing.T) {
		err := NewUserAlreadyInGroupError("testuser", "testgroup")
		assert.Equal(t, "EntityAlreadyExists: User testuser in group testgroup already exists", err.Error())
		assert.Equal(t, http.StatusConflict, err.GetHTTPStatusCode())
	})

	t.Run("NewNoSuchRoleError", func(t *testing.T) {
		err := NewNoSuchRoleError("testrole")
		assert.Equal(t, "NoSuchEntity: The role with name testrole cannot be found.", err.Error())
		assert.Equal(t, http.StatusNotFound, err.GetHTTPStatusCode())
	})

	t.Run("NewRoleAlreadyExistsError", func(t *testing.T) {
		err := NewRoleAlreadyExistsError("testrole")
		assert.Equal(t, "EntityAlreadyExists: Role testrole already exists", err.Error())
		assert.Equal(t, http.StatusConflict, err.GetHTTPStatusCode())
	})

	t.Run("NewDeleteRoleConflictError", func(t *testing.T) {
		err := NewDeleteRoleConflictError("must detach all policies")
		assert.Equal(t, "DeleteConflict: must detach all policies", err.Error())
		assert.Equal(t, http.StatusConflict, err.GetHTTPStatusCode())
	})

	t.Run("NewNoSuchInstanceProfileError", func(t *testing.T) {
		err := NewNoSuchInstanceProfileError("testprofile")
		assert.Equal(t, "NoSuchEntity: The Instance Profile with name testprofile cannot be found.", err.Error())
		assert.Equal(t, http.StatusNotFound, err.GetHTTPStatusCode())
	})

	t.Run("NewInstanceProfileAlreadyExistsError", func(t *testing.T) {
		err := NewInstanceProfileAlreadyExistsError("testprofile")
		assert.Equal(t, "EntityAlreadyExists: Instance Profile testprofile already exists", err.Error())
		assert.Equal(t, http.StatusConflict, err.GetHTTPStatusCode())
	})

	t.Run("NewRoleNotInInstanceProfileError", func(t *testing.T) {
		err := NewRoleNotInInstanceProfileError("testrole", "testprofile")
		assert.Equal(t, "NoSuchEntity: The Role testrole in instance profile with name testprofile cannot be found.", err.Error())
		assert.Equal(t, http.StatusNotFound, err.GetHTTPStatusCode())
	})

	t.Run("NewRoleAlreadyInInstanceProfileError", func(t *testing.T) {
		err := NewRoleAlreadyInInstanceProfileError("testrole", "testprofile")
		assert.Equal(t, "LimitExceededException: Cannot exceed quota for InstanceProfilesPerRole: 1. Already associated with instance profile testprofile.", err.Error())
		assert.Equal(t, http.StatusBadRequest, err.GetHTTPStatusCode())
	})

	t.Run("NewNoSuchPolicyError", func(t *testing.T) {
		err := NewNoSuchPolicyError("arn:aws:iam::123456789012:policy/test")
		assert.Equal(t, "NoSuchEntity: The policy with name arn:aws:iam::123456789012:policy/test cannot be found.", err.Error())
		assert.Equal(t, http.StatusNotFound, err.GetHTTPStatusCode())
	})

	t.Run("NewPolicyAlreadyExistsError", func(t *testing.T) {
		err := NewPolicyAlreadyExistsError("testpolicy")
		assert.Equal(t, "EntityAlreadyExists: Policy testpolicy already exists", err.Error())
		assert.Equal(t, http.StatusConflict, err.GetHTTPStatusCode())
	})

	t.Run("NewDeletePolicyConflictError", func(t *testing.T) {
		err := NewDeletePolicyConflictError("arn:aws:iam::123456789012:policy/test")
		assert.Equal(t, "DeleteConflict: Cannot delete policy arn:aws:iam::123456789012:policy/test, there are attachments.", err.Error())
		assert.Equal(t, http.StatusConflict, err.GetHTTPStatusCode())
	})

	t.Run("NewNoSuchPolicyVersionError", func(t *testing.T) {
		err := NewNoSuchPolicyVersionError("v1")
		assert.Equal(t, "NoSuchEntity: The Policy version with name v1 cannot be found.", err.Error())
		assert.Equal(t, http.StatusNotFound, err.GetHTTPStatusCode())
	})

	t.Run("NewNoSuchMFADeviceError", func(t *testing.T) {
		err := NewNoSuchMFADeviceError("arn:aws:iam::123456789012:mfa/testuser")
		assert.Equal(t, "NoSuchEntity: The MFA Device with name arn:aws:iam::123456789012:mfa/testuser cannot be found.", err.Error())
		assert.Equal(t, http.StatusNotFound, err.GetHTTPStatusCode())
	})

	t.Run("NewMFADeviceAlreadyAssignedError", func(t *testing.T) {
		err := NewMFADeviceAlreadyAssignedError("arn:aws:iam::123456789012:mfa/testuser")
		assert.Equal(t, "EntityAlreadyExists: MFA Device arn:aws:iam::123456789012:mfa/testuser already exists", err.Error())
		assert.Equal(t, http.StatusConflict, err.GetHTTPStatusCode())
	})

	t.Run("NewNoSuchPasswordPolicyError", func(t *testing.T) {
		err := NewNoSuchPasswordPolicyError()
		assert.Equal(t, "NoSuchEntity: The Password Policy with name  cannot be found.", err.Error())
		assert.Equal(t, http.StatusNotFound, err.GetHTTPStatusCode())
	})

	t.Run("NewInvalidInputError", func(t *testing.T) {
		err := NewInvalidInputError("RoleName", "invalid characters")
		assert.Equal(t, "InvalidInput: The input parameter RoleName is invalid: invalid characters", err.Error())
		assert.Equal(t, http.StatusBadRequest, err.GetHTTPStatusCode())
	})
}
