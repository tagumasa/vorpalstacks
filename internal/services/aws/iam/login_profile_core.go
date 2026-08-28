// Transport-agnostic Core functions for IAM login profiles and password
// changes: validation and store operations shared by the AWS-compatible
// HTTP API handlers and any admin plane paths (the xxxCore pattern).
package iam

import (
	"errors"

	"vorpalstacks/internal/common/request"
	iamstore "vorpalstacks/internal/store/aws/iam"
)

// CreateLoginProfileInput holds the parameters for creating a login
// profile.
type CreateLoginProfileInput struct {
	UserName              string
	Password              string
	PasswordResetRequired bool
}

// UpdateLoginProfileInput holds the parameters for updating a login
// profile.  A nil PasswordResetRequired leaves the flag untouched.
type UpdateLoginProfileInput struct {
	UserName              string
	Password              string
	PasswordResetRequired *bool
}

// ChangePasswordInput holds the parameters for the authenticated caller's
// own password change.
type ChangePasswordInput struct {
	OldPassword string
	NewPassword string
}

// getLoginProfileCore retrieves the login profile for a user.
func (s *IAMService) getLoginProfileCore(store *iamstore.IAMStore, userName string) (*iamstore.LoginProfile, error) {
	if !store.Users().Exists(userName) {
		return nil, NewNoSuchUserError(userName)
	}

	profile, err := store.LoginProfiles().Get(userName)
	if err != nil {
		return nil, NewNoSuchLoginProfileError(userName)
	}
	return profile, nil
}

// createLoginProfileCore validates input and creates a login profile for a
// user.  The password must comply with the account password policy.
func (s *IAMService) createLoginProfileCore(store *iamstore.IAMStore, input *CreateLoginProfileInput) (*iamstore.LoginProfile, error) {
	if input.UserName == "" {
		return nil, NewValidationError("UserName")
	}
	if input.Password == "" {
		return nil, ErrPasswordPolicyViolation
	}

	if !store.Users().Exists(input.UserName) {
		return nil, NewNoSuchUserError(input.UserName)
	}

	passwordPolicy := store.PasswordPolicy().GetOrDefault()
	if !validatePasswordAgainstPolicy(input.Password, passwordPolicy) {
		return nil, ErrPasswordPolicyViolation
	}

	profile, err := store.LoginProfiles().Create(input.UserName, input.Password, input.PasswordResetRequired)
	if err != nil {
		if errors.Is(err, iamstore.ErrLoginProfileExists) {
			return nil, NewLoginProfileAlreadyExistsError(input.UserName)
		}
		return nil, err
	}
	return profile, nil
}

// deleteLoginProfileCore validates input and deletes the login profile for
// a user.
func (s *IAMService) deleteLoginProfileCore(store *iamstore.IAMStore, userName string) error {
	if !store.Users().Exists(userName) {
		return NewNoSuchUserError(userName)
	}

	if !store.LoginProfiles().Exists(userName) {
		return NewNoSuchLoginProfileError(userName)
	}

	return store.LoginProfiles().Delete(userName)
}

// updateLoginProfileCore validates input and updates the login profile for
// a user: the password and/or the password reset requirement.  The new
// password must comply with the account password policy including the
// password reuse prevention window.
func (s *IAMService) updateLoginProfileCore(store *iamstore.IAMStore, input *UpdateLoginProfileInput) (*iamstore.LoginProfile, error) {
	if input.UserName == "" {
		return nil, NewValidationError("UserName")
	}

	if !store.Users().Exists(input.UserName) {
		return nil, NewNoSuchUserError(input.UserName)
	}

	if !store.LoginProfiles().Exists(input.UserName) {
		return nil, NewNoSuchLoginProfileError(input.UserName)
	}

	if input.Password != "" {
		passwordPolicy := store.PasswordPolicy().GetOrDefault()
		if !validatePasswordAgainstPolicy(input.Password, passwordPolicy) {
			return nil, ErrPasswordPolicyViolation
		}
		if passwordPolicy.PasswordReusePrevention > 0 {
			reused, err := store.LoginProfiles().CheckPasswordReuse(input.UserName, input.Password, passwordPolicy.PasswordReusePrevention)
			if err != nil {
				return nil, err
			}
			if reused {
				return nil, ErrPasswordPolicyViolation
			}
		}
		if err := store.LoginProfiles().UpdatePassword(input.UserName, input.Password); err != nil {
			return nil, err
		}
	}

	if input.PasswordResetRequired != nil {
		if err := store.LoginProfiles().UpdatePasswordResetRequired(input.UserName, *input.PasswordResetRequired); err != nil {
			return nil, err
		}
	}

	return store.LoginProfiles().Get(input.UserName)
}

// changePasswordCore changes the password for the authenticated caller.
// The operation targets the caller itself; the request carries no UserName
// member on the wire, so the caller's user name is resolved from the
// request context after the required-member validation.  The old password
// must be verified and the new password must comply with the account
// password policy.
func (s *IAMService) changePasswordCore(reqCtx *request.RequestContext, store *iamstore.IAMStore, input *ChangePasswordInput) error {
	if input.OldPassword == "" || input.NewPassword == "" {
		return ErrPasswordPolicyViolation
	}

	userName, err := resolveUserName(reqCtx, "")
	if err != nil {
		return err
	}

	if !store.Users().Exists(userName) {
		return NewNoSuchUserError(userName)
	}

	if !store.LoginProfiles().Exists(userName) {
		return NewNoSuchLoginProfileError(userName)
	}

	valid, err := store.LoginProfiles().VerifyPassword(userName, input.OldPassword)
	if err != nil {
		return err
	}
	if !valid {
		return ErrNotAuthorized
	}

	passwordPolicy := store.PasswordPolicy().GetOrDefault()
	if !validatePasswordAgainstPolicy(input.NewPassword, passwordPolicy) {
		return ErrPasswordPolicyViolation
	}

	if passwordPolicy.PasswordReusePrevention > 0 {
		reused, err := store.LoginProfiles().CheckPasswordReuse(userName, input.NewPassword, passwordPolicy.PasswordReusePrevention)
		if err != nil {
			return err
		}
		if reused {
			return ErrPasswordPolicyViolation
		}
	}

	return store.LoginProfiles().UpdatePassword(userName, input.NewPassword)
}
