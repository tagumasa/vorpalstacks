// Package iam provides IAM service operations for vorpalstacks.
package iam

import (
	"context"

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
	profile, err := s.getLoginProfileCore(store, userName)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"LoginProfile": loginProfileToResponse(profile),
	}, nil
}

// CreateLoginProfile creates a login profile for a user.
// The password must comply with the account password policy.
// Returns the created login profile details.
func (s *IAMService) CreateLoginProfile(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	input := &CreateLoginProfileInput{
		UserName:              request.GetStringParam(req.Parameters, "UserName"),
		Password:              request.GetStringParam(req.Parameters, "Password"),
		PasswordResetRequired: request.GetBoolParam(req.Parameters, "PasswordResetRequired"),
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	profile, err := s.createLoginProfileCore(store, input)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"LoginProfile": loginProfileToResponse(profile),
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
	if err := s.deleteLoginProfileCore(store, userName); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// UpdateLoginProfile updates the login profile for a user.
// Can update the password and/or password reset requirement.
// The new password must comply with the account password policy.
func (s *IAMService) UpdateLoginProfile(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	input := &UpdateLoginProfileInput{
		UserName: request.GetStringParam(req.Parameters, "UserName"),
		Password: request.GetStringParam(req.Parameters, "Password"),
	}
	if raw, ok := req.Parameters["PasswordResetRequired"]; ok {
		required := false
		switch v := raw.(type) {
		case bool:
			required = v
		case string:
			required = v == "true"
		}
		input.PasswordResetRequired = &required
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	profile, err := s.updateLoginProfileCore(store, input)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"LoginProfile": loginProfileToResponse(profile),
	}, nil
}

// ChangePassword changes the password for a user.
// Requires the old password to be verified before setting the new password.
// The new password must comply with the account password policy.
func (s *IAMService) ChangePassword(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	input := &ChangePasswordInput{
		OldPassword: request.GetStringParam(req.Parameters, "OldPassword"),
		NewPassword: request.GetStringParam(req.Parameters, "NewPassword"),
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := s.changePasswordCore(reqCtx, store, input); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

func loginProfileToResponse(profile *iamstore.LoginProfile) map[string]interface{} {
	return map[string]interface{}{
		"UserName":              profile.UserName,
		"CreateDate":            profile.CreateDate.Format(timeutils.ISO8601SimpleFormat),
		"PasswordResetRequired": profile.PasswordResetRequired,
	}
}
