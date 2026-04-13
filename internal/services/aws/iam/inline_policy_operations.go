// Package iam provides IAM service operations for vorpalstacks.
package iam

import (
	"context"
	"vorpalstacks/internal/common/request"
)

// PutUserPolicy creates or updates an inline policy for a user.
func (s *IAMService) PutUserPolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	userName := request.GetStringParam(req.Parameters, "UserName")
	policyName := request.GetStringParam(req.Parameters, "PolicyName")
	policyDocument := request.GetStringParam(req.Parameters, "PolicyDocument")

	if userName == "" {
		return nil, ErrNoSuchUser
	}
	if policyName == "" {
		return nil, ErrNoSuchPolicy
	}
	if !validatePolicyDocument(policyDocument) {
		return nil, ErrMalformedPolicyDocument
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if !store.Users().Exists(userName) {
		return nil, NewNoSuchUserError(userName)
	}

	if err := store.InlinePolicies().Put(PrincipalTypeUser, userName, policyName, policyDocument); err != nil {
		return nil, err
	}

	return nil, nil
}

// GetUserPolicy retrieves an inline policy for a user.
func (s *IAMService) GetUserPolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	userName := request.GetStringParam(req.Parameters, "UserName")
	policyName := request.GetStringParam(req.Parameters, "PolicyName")

	if userName == "" {
		return nil, ErrNoSuchUser
	}
	if policyName == "" {
		return nil, ErrNoSuchPolicy
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if !store.Users().Exists(userName) {
		return nil, NewNoSuchUserError(userName)
	}

	policy, err := store.InlinePolicies().Get(PrincipalTypeUser, userName, policyName)
	if err != nil {
		return nil, NewNoSuchPolicyError(policyName)
	}

	return map[string]interface{}{
		"UserName":       userName,
		"PolicyName":     policyName,
		"PolicyDocument": policy.PolicyDocument,
	}, nil
}

// DeleteUserPolicy deletes an inline policy from a user.
func (s *IAMService) DeleteUserPolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	userName := request.GetStringParam(req.Parameters, "UserName")
	policyName := request.GetStringParam(req.Parameters, "PolicyName")

	if userName == "" {
		return nil, ErrNoSuchUser
	}
	if policyName == "" {
		return nil, ErrNoSuchPolicy
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if !store.InlinePolicies().Exists(PrincipalTypeUser, userName, policyName) {
		return nil, NewNoSuchPolicyError(policyName)
	}

	if err := store.InlinePolicies().Delete(PrincipalTypeUser, userName, policyName); err != nil {
		return nil, err
	}

	return nil, nil
}

// ListUserPolicies lists inline policies for a user.
func (s *IAMService) ListUserPolicies(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
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

	policyNames, err := store.InlinePolicies().List(PrincipalTypeUser, userName)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"PolicyNames": policyNames,
		"IsTruncated": false,
	}, nil
}

// PutGroupPolicy creates or updates an inline policy for a group.
func (s *IAMService) PutGroupPolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	groupName := request.GetStringParam(req.Parameters, "GroupName")
	policyName := request.GetStringParam(req.Parameters, "PolicyName")
	policyDocument := request.GetStringParam(req.Parameters, "PolicyDocument")

	if groupName == "" {
		return nil, ErrNoSuchGroup
	}
	if policyName == "" {
		return nil, ErrNoSuchPolicy
	}
	if !validatePolicyDocument(policyDocument) {
		return nil, ErrMalformedPolicyDocument
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if !store.Groups().Exists(groupName) {
		return nil, NewNoSuchGroupError(groupName)
	}

	if err := store.InlinePolicies().Put(PrincipalTypeGroup, groupName, policyName, policyDocument); err != nil {
		return nil, err
	}

	return nil, nil
}

// GetGroupPolicy retrieves an inline policy for a group.
func (s *IAMService) GetGroupPolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	groupName := request.GetStringParam(req.Parameters, "GroupName")
	policyName := request.GetStringParam(req.Parameters, "PolicyName")

	if groupName == "" {
		return nil, ErrNoSuchGroup
	}
	if policyName == "" {
		return nil, ErrNoSuchPolicy
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if !store.Groups().Exists(groupName) {
		return nil, NewNoSuchGroupError(groupName)
	}

	policy, err := store.InlinePolicies().Get(PrincipalTypeGroup, groupName, policyName)
	if err != nil {
		return nil, NewNoSuchPolicyError(policyName)
	}

	return map[string]interface{}{
		"GroupName":      groupName,
		"PolicyName":     policyName,
		"PolicyDocument": policy.PolicyDocument,
	}, nil
}

// DeleteGroupPolicy deletes an inline policy from a group.
func (s *IAMService) DeleteGroupPolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	groupName := request.GetStringParam(req.Parameters, "GroupName")
	policyName := request.GetStringParam(req.Parameters, "PolicyName")

	if groupName == "" {
		return nil, ErrNoSuchGroup
	}
	if policyName == "" {
		return nil, ErrNoSuchPolicy
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if !store.InlinePolicies().Exists(PrincipalTypeGroup, groupName, policyName) {
		return nil, NewNoSuchPolicyError(policyName)
	}

	if err := store.InlinePolicies().Delete(PrincipalTypeGroup, groupName, policyName); err != nil {
		return nil, err
	}

	return nil, nil
}

// ListGroupPolicies lists inline policies for a group.
func (s *IAMService) ListGroupPolicies(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	groupName := request.GetStringParam(req.Parameters, "GroupName")
	if groupName == "" {
		return nil, ErrNoSuchGroup
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if !store.Groups().Exists(groupName) {
		return nil, NewNoSuchGroupError(groupName)
	}

	policyNames, err := store.InlinePolicies().List(PrincipalTypeGroup, groupName)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"PolicyNames": policyNames,
		"IsTruncated": false,
	}, nil
}

// PutRolePolicy creates or updates an inline policy for a role.
func (s *IAMService) PutRolePolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	roleName := request.GetStringParam(req.Parameters, "RoleName")
	policyName := request.GetStringParam(req.Parameters, "PolicyName")
	policyDocument := request.GetStringParam(req.Parameters, "PolicyDocument")

	if roleName == "" {
		return nil, ErrNoSuchRole
	}
	if policyName == "" {
		return nil, ErrNoSuchPolicy
	}
	if !validatePolicyDocument(policyDocument) {
		return nil, ErrMalformedPolicyDocument
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if !store.Roles().Exists(roleName) {
		return nil, NewNoSuchRoleError(roleName)
	}

	if err := store.InlinePolicies().Put(PrincipalTypeRole, roleName, policyName, policyDocument); err != nil {
		return nil, err
	}

	return nil, nil
}

// GetRolePolicy retrieves an inline policy for a role.
func (s *IAMService) GetRolePolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	roleName := request.GetStringParam(req.Parameters, "RoleName")
	policyName := request.GetStringParam(req.Parameters, "PolicyName")

	if roleName == "" {
		return nil, ErrNoSuchRole
	}
	if policyName == "" {
		return nil, ErrNoSuchPolicy
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if !store.Roles().Exists(roleName) {
		return nil, NewNoSuchRoleError(roleName)
	}

	policy, err := store.InlinePolicies().Get(PrincipalTypeRole, roleName, policyName)
	if err != nil {
		return nil, NewNoSuchPolicyError(policyName)
	}

	return map[string]interface{}{
		"RoleName":       roleName,
		"PolicyName":     policyName,
		"PolicyDocument": policy.PolicyDocument,
	}, nil
}

// DeleteRolePolicy deletes an inline policy from a role.
func (s *IAMService) DeleteRolePolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	roleName := request.GetStringParam(req.Parameters, "RoleName")
	policyName := request.GetStringParam(req.Parameters, "PolicyName")

	if roleName == "" {
		return nil, ErrNoSuchRole
	}
	if policyName == "" {
		return nil, ErrNoSuchPolicy
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if !store.InlinePolicies().Exists(PrincipalTypeRole, roleName, policyName) {
		return nil, NewNoSuchPolicyError(policyName)
	}

	if err := store.InlinePolicies().Delete(PrincipalTypeRole, roleName, policyName); err != nil {
		return nil, err
	}

	return nil, nil
}

// ListRolePolicies lists inline policies for a role.
func (s *IAMService) ListRolePolicies(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	roleName := request.GetStringParam(req.Parameters, "RoleName")
	if roleName == "" {
		return nil, ErrNoSuchRole
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if !store.Roles().Exists(roleName) {
		return nil, NewNoSuchRoleError(roleName)
	}

	policyNames, err := store.InlinePolicies().List(PrincipalTypeRole, roleName)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"PolicyNames": policyNames,
		"IsTruncated": false,
	}, nil
}
