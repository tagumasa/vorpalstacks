// Package iam provides IAM service operations for vorpalstacks.
package iam

import (
	"context"
	"errors"
	"regexp"
	"unicode"

	"vorpalstacks/internal/services/aws/common/request"
	iamstore "vorpalstacks/internal/store/aws/iam"
	"vorpalstacks/internal/utils/timeutils"
)

func validatePasswordAgainstPolicy(password string, policy *iamstore.AccountPasswordPolicy) bool {
	if len(password) < policy.MinimumPasswordLength {
		return false
	}

	if policy.RequireSymbols {
		hasSymbol := regexp.MustCompile(`[!@#$%^&*()_+\-=\[\]{};':",\.<>?/\\|~]`).MatchString(password)
		if !hasSymbol {
			return false
		}
	}

	if policy.RequireNumbers {
		hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)
		if !hasNumber {
			return false
		}
	}

	if policy.RequireUppercaseCharacters {
		hasUpper := false
		for _, r := range password {
			if unicode.IsUpper(r) {
				hasUpper = true
				break
			}
		}
		if !hasUpper {
			return false
		}
	}

	if policy.RequireLowercaseCharacters {
		hasLower := false
		for _, r := range password {
			if unicode.IsLower(r) {
				hasLower = true
				break
			}
		}
		if !hasLower {
			return false
		}
	}

	return true
}

// GetLoginProfile retrieves the login profile for a user.
// Returns the login profile details including username, creation date,
// and whether a password reset is required.
func (s *IAMService) GetLoginProfile(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	userName := request.GetStringParam(req.Parameters, "UserName")
	if userName == "" {
		return nil, ErrNoSuchUser
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
		return nil, ErrNoSuchUser
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
	if userName == "" {
		return nil, ErrNoSuchUser
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

	return nil, nil
}

// UpdateLoginProfile updates the login profile for a user.
// Can update the password and/or password reset requirement.
// The new password must comply with the account password policy.
func (s *IAMService) UpdateLoginProfile(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	userName := request.GetStringParam(req.Parameters, "UserName")
	if userName == "" {
		return nil, ErrNoSuchUser
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
		return nil, ErrNoSuchUser
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

	if err := store.LoginProfiles().UpdatePassword(userName, newPassword); err != nil {
		return nil, err
	}

	return nil, nil
}
