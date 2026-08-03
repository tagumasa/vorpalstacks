// Package iam provides IAM service operations for vorpalstacks.
package iam

import (
	"context"
	"errors"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	iamstore "vorpalstacks/internal/store/aws/iam"
	"vorpalstacks/internal/utils/timeutils"
)

// GetLoginProfile retrieves the login profile for a user.
// Returns the login profile details including username, creation date,
// and whether a password reset is required.
func (s *IAMService) GetLoginProfile(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	userName := request.GetStringParam(req.Parameters, "UserName")
	userName, err := resolveUserName(reqCtx, userName)
	if err != nil {
		return nil, err
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if !store.Users().Exists(userName) {
		return nil, NewNoSuchUserError(userName)
	}

	profile, err := store.LoginProfiles().Get(userName)
	if err != nil {
		return nil, NewNoSuchLoginProfileError(userName)
	}

	return map[string]interface{}{
		"LoginProfile": map[string]interface{}{
			"UserName":              profile.UserName,
			"CreateDate":            profile.CreateDate.Format(timeutils.ISO8601SimpleFormat),
			"PasswordResetRequired": profile.PasswordResetRequired,
		},
	}, nil
}

// CreateLoginProfile creates a login profile for a user.
// The password must comply with the account password policy.
// Returns the created login profile details.
func (s *IAMService) CreateLoginProfile(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	userName := request.GetStringParam(req.Parameters, "UserName")
	if userName == "" {
		return nil, NewValidationError("UserName")
	}

	password := request.GetStringParam(req.Parameters, "Password")
	if password == "" {
		return nil, ErrPasswordPolicyViolation
	}

	passwordResetRequired := request.GetBoolParam(req.Parameters, "PasswordResetRequired")

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if !store.Users().Exists(userName) {
		return nil, NewNoSuchUserError(userName)
	}

	passwordPolicy := store.PasswordPolicy().GetOrDefault()
	if !validatePasswordAgainstPolicy(password, passwordPolicy) {
		return nil, ErrPasswordPolicyViolation
	}

	profile, err := store.LoginProfiles().Create(userName, password, passwordResetRequired)
	if err != nil {
		if errors.Is(err, iamstore.ErrLoginProfileExists) {
			return nil, NewLoginProfileAlreadyExistsError(userName)
		}
		return nil, err
	}

	return map[string]interface{}{
		"LoginProfile": map[string]interface{}{
			"UserName":              profile.UserName,
			"CreateDate":            profile.CreateDate.Format(timeutils.ISO8601SimpleFormat),
			"PasswordResetRequired": profile.PasswordResetRequired,
		},
	}, nil
}

// DeleteLoginProfile deletes the login profile for a user.
// Returns an error if the user or login profile does not exist.
func (s *IAMService) DeleteLoginProfile(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	userName := request.GetStringParam(req.Parameters, "UserName")
	userName, err := resolveUserName(reqCtx, userName)
	if err != nil {
		return nil, err
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if !store.Users().Exists(userName) {
		return nil, NewNoSuchUserError(userName)
	}

	if !store.LoginProfiles().Exists(userName) {
		return nil, NewNoSuchLoginProfileError(userName)
	}

	if err := store.LoginProfiles().Delete(userName); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// UpdateLoginProfile updates the login profile for a user.
// Can update the password and/or password reset requirement.
// The new password must comply with the account password policy.
func (s *IAMService) UpdateLoginProfile(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	userName := request.GetStringParam(req.Parameters, "UserName")
	if userName == "" {
		return nil, NewValidationError("UserName")
	}

	password := request.GetStringParam(req.Parameters, "Password")
	passwordResetRequired, hasPasswordResetRequired := req.Parameters["PasswordResetRequired"]

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if !store.Users().Exists(userName) {
		return nil, NewNoSuchUserError(userName)
	}

	if !store.LoginProfiles().Exists(userName) {
		return nil, NewNoSuchLoginProfileError(userName)
	}

	if password != "" {
		passwordPolicy := store.PasswordPolicy().GetOrDefault()
		if !validatePasswordAgainstPolicy(password, passwordPolicy) {
			return nil, ErrPasswordPolicyViolation
		}
		if passwordPolicy.PasswordReusePrevention > 0 {
			reused, err := store.LoginProfiles().CheckPasswordReuse(userName, password, passwordPolicy.PasswordReusePrevention)
			if err != nil {
				return nil, err
			}
			if reused {
				return nil, ErrPasswordPolicyViolation
			}
		}
		if err := store.LoginProfiles().UpdatePassword(userName, password); err != nil {
			return nil, err
		}
	}

	if hasPasswordResetRequired {
		required := false
		switch v := passwordResetRequired.(type) {
		case bool:
			required = v
		case string:
			required = v == "true"
		}
		if err := store.LoginProfiles().UpdatePasswordResetRequired(userName, required); err != nil {
			return nil, err
		}
	}

	profile, err := store.LoginProfiles().Get(userName)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"LoginProfile": map[string]interface{}{
			"UserName":              profile.UserName,
			"CreateDate":            profile.CreateDate.Format(timeutils.ISO8601SimpleFormat),
			"PasswordResetRequired": profile.PasswordResetRequired,
		},
	}, nil
}

// ChangePassword changes the password for a user.
// Requires the old password to be verified before setting the new password.
// The new password must comply with the account password policy.
func (s *IAMService) ChangePassword(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	oldPassword := request.GetStringParam(req.Parameters, "OldPassword")
	newPassword := request.GetStringParam(req.Parameters, "NewPassword")

	if oldPassword == "" || newPassword == "" {
		return nil, ErrPasswordPolicyViolation
	}

	userName := request.GetStringParam(req.Parameters, "UserName")
	if userName == "" {
		return nil, NewValidationError("UserName")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if !store.Users().Exists(userName) {
		return nil, NewNoSuchUserError(userName)
	}

	if !store.LoginProfiles().Exists(userName) {
		return nil, NewNoSuchLoginProfileError(userName)
	}

	valid, err := store.LoginProfiles().VerifyPassword(userName, oldPassword)
	if err != nil {
		return nil, err
	}
	if !valid {
		return nil, ErrNotAuthorized
	}

	passwordPolicy := store.PasswordPolicy().GetOrDefault()
	if !validatePasswordAgainstPolicy(newPassword, passwordPolicy) {
		return nil, ErrPasswordPolicyViolation
	}

	if passwordPolicy.PasswordReusePrevention > 0 {
		reused, err := store.LoginProfiles().CheckPasswordReuse(userName, newPassword, passwordPolicy.PasswordReusePrevention)
		if err != nil {
			return nil, err
		}
		if reused {
			return nil, ErrPasswordPolicyViolation
		}
	}

	if err := store.LoginProfiles().UpdatePassword(userName, newPassword); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}
