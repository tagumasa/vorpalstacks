// Package iam provides IAM service operations for vorpalstacks.
package iam

import (
	"context"
	"errors"

	"vorpalstacks/internal/common/request"
)

const (
	// PrincipalTypeUser represents an IAM user principal type.
	PrincipalTypeUser = "user"
	// PrincipalTypeGroup represents an IAM group principal type.
	PrincipalTypeGroup = "group"
	// PrincipalTypeRole represents an IAM role principal type.
	PrincipalTypeRole = "role"
)

// AttachUserPolicy attaches a policy to a user.
func (s *IAMService) AttachUserPolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	userName := request.GetStringParam(req.Parameters, "UserName")
	policyArn := request.GetStringParam(req.Parameters, "PolicyArn")

	if userName == "" {
		return nil, ErrNoSuchUser
	}
	if policyArn == "" {
		return nil, ErrNoSuchPolicy
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if !store.Users().Exists(userName) {
		return nil, NewNoSuchUserError(userName)
	}

	if !store.Policies().Exists(policyArn) {
		return nil, NewNoSuchPolicyError(policyArn)
	}

	if store.AttachedPolicies().IsAttached(PrincipalTypeUser, userName, policyArn) {
		return nil, nil
	}

	if err := store.AttachedPolicies().Attach(PrincipalTypeUser, userName, policyArn); err != nil {
		return nil, err
	}

	if err := store.Policies().IncrementAttachmentCount(policyArn); err != nil {
		if rollbackErr := store.AttachedPolicies().Detach(PrincipalTypeUser, userName, policyArn); rollbackErr != nil {
			return nil, errors.Join(err, rollbackErr)
		}
		return nil, err
	}

	return nil, nil
}

// DetachUserPolicy detaches a policy from a user.
func (s *IAMService) DetachUserPolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	userName := request.GetStringParam(req.Parameters, "UserName")
	policyArn := request.GetStringParam(req.Parameters, "PolicyArn")

	if userName == "" {
		return nil, ErrNoSuchUser
	}
	if policyArn == "" {
		return nil, ErrNoSuchPolicy
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if !store.AttachedPolicies().IsAttached(PrincipalTypeUser, userName, policyArn) {
		return nil, NewPolicyNotAttachedError(policyArn)
	}

	if err := store.AttachedPolicies().Detach(PrincipalTypeUser, userName, policyArn); err != nil {
		return nil, err
	}

	if err := store.Policies().DecrementAttachmentCount(policyArn); err != nil {
		if rollbackErr := store.AttachedPolicies().Attach(PrincipalTypeUser, userName, policyArn); rollbackErr != nil {
			return nil, errors.Join(err, rollbackErr)
		}
		return nil, err
	}

	return nil, nil
}

// ListAttachedUserPolicies lists the policies attached to a user.
func (s *IAMService) ListAttachedUserPolicies(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
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

	policyArns, err := store.AttachedPolicies().ListAttachedPolicies(PrincipalTypeUser, userName)
	if err != nil {
		return nil, err
	}

	policies := make([]interface{}, 0, len(policyArns))
	for _, arn := range policyArns {
		if policy, err := store.Policies().Get(arn); err == nil {
			policies = append(policies, map[string]interface{}{
				"PolicyName": policy.PolicyName,
				"PolicyArn":  policy.Arn,
			})
		}
	}

	return map[string]interface{}{
		"AttachedPolicies": policies,
		"IsTruncated":      false,
	}, nil
}

// AttachGroupPolicy attaches a policy to a group.
func (s *IAMService) AttachGroupPolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	groupName := request.GetStringParam(req.Parameters, "GroupName")
	policyArn := request.GetStringParam(req.Parameters, "PolicyArn")

	if groupName == "" {
		return nil, ErrNoSuchGroup
	}
	if policyArn == "" {
		return nil, ErrNoSuchPolicy
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if !store.Groups().Exists(groupName) {
		return nil, NewNoSuchGroupError(groupName)
	}

	if !store.Policies().Exists(policyArn) {
		return nil, NewNoSuchPolicyError(policyArn)
	}

	if store.AttachedPolicies().IsAttached(PrincipalTypeGroup, groupName, policyArn) {
		return nil, nil
	}

	if err := store.AttachedPolicies().Attach(PrincipalTypeGroup, groupName, policyArn); err != nil {
		return nil, err
	}

	if err := store.Policies().IncrementAttachmentCount(policyArn); err != nil {
		if rollbackErr := store.AttachedPolicies().Detach(PrincipalTypeGroup, groupName, policyArn); rollbackErr != nil {
			return nil, errors.Join(err, rollbackErr)
		}
		return nil, err
	}

	return nil, nil
}

// DetachGroupPolicy detaches a policy from a group.
func (s *IAMService) DetachGroupPolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	groupName := request.GetStringParam(req.Parameters, "GroupName")
	policyArn := request.GetStringParam(req.Parameters, "PolicyArn")

	if groupName == "" {
		return nil, ErrNoSuchGroup
	}
	if policyArn == "" {
		return nil, ErrNoSuchPolicy
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if !store.AttachedPolicies().IsAttached(PrincipalTypeGroup, groupName, policyArn) {
		return nil, NewPolicyNotAttachedError(policyArn)
	}

	if err := store.AttachedPolicies().Detach(PrincipalTypeGroup, groupName, policyArn); err != nil {
		return nil, err
	}

	if err := store.Policies().DecrementAttachmentCount(policyArn); err != nil {
		if rollbackErr := store.AttachedPolicies().Attach(PrincipalTypeGroup, groupName, policyArn); rollbackErr != nil {
			return nil, errors.Join(err, rollbackErr)
		}
		return nil, err
	}

	return nil, nil
}

// ListAttachedGroupPolicies lists the policies attached to a group.
func (s *IAMService) ListAttachedGroupPolicies(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
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

	policyArns, err := store.AttachedPolicies().ListAttachedPolicies(PrincipalTypeGroup, groupName)
	if err != nil {
		return nil, err
	}

	policies := make([]interface{}, 0, len(policyArns))
	for _, arn := range policyArns {
		if policy, err := store.Policies().Get(arn); err == nil {
			policies = append(policies, map[string]interface{}{
				"PolicyName": policy.PolicyName,
				"PolicyArn":  policy.Arn,
			})
		}
	}

	return map[string]interface{}{
		"AttachedPolicies": policies,
		"IsTruncated":      false,
	}, nil
}

// AttachRolePolicy attaches a policy to a role.
func (s *IAMService) AttachRolePolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	roleName := request.GetStringParam(req.Parameters, "RoleName")
	policyArn := request.GetStringParam(req.Parameters, "PolicyArn")

	if roleName == "" {
		return nil, ErrNoSuchRole
	}
	if policyArn == "" {
		return nil, ErrNoSuchPolicy
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if !store.Roles().Exists(roleName) {
		return nil, NewNoSuchRoleError(roleName)
	}

	if !store.Policies().Exists(policyArn) {
		return nil, NewNoSuchPolicyError(policyArn)
	}

	if store.AttachedPolicies().IsAttached(PrincipalTypeRole, roleName, policyArn) {
		return nil, nil
	}

	if err := store.AttachedPolicies().Attach(PrincipalTypeRole, roleName, policyArn); err != nil {
		return nil, err
	}

	if err := store.Policies().IncrementAttachmentCount(policyArn); err != nil {
		if rollbackErr := store.AttachedPolicies().Detach(PrincipalTypeRole, roleName, policyArn); rollbackErr != nil {
			return nil, errors.Join(err, rollbackErr)
		}
		return nil, err
	}

	return nil, nil
}

// DetachRolePolicy detaches a policy from a role.
func (s *IAMService) DetachRolePolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	roleName := request.GetStringParam(req.Parameters, "RoleName")
	policyArn := request.GetStringParam(req.Parameters, "PolicyArn")

	if roleName == "" {
		return nil, ErrNoSuchRole
	}
	if policyArn == "" {
		return nil, ErrNoSuchPolicy
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if !store.AttachedPolicies().IsAttached(PrincipalTypeRole, roleName, policyArn) {
		return nil, NewPolicyNotAttachedError(policyArn)
	}

	if err := store.AttachedPolicies().Detach(PrincipalTypeRole, roleName, policyArn); err != nil {
		return nil, err
	}

	if err := store.Policies().DecrementAttachmentCount(policyArn); err != nil {
		if rollbackErr := store.AttachedPolicies().Attach(PrincipalTypeRole, roleName, policyArn); rollbackErr != nil {
			return nil, errors.Join(err, rollbackErr)
		}
		return nil, err
	}

	return nil, nil
}

// ListAttachedRolePolicies lists the policies attached to a role.
func (s *IAMService) ListAttachedRolePolicies(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
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

	policyArns, err := store.AttachedPolicies().ListAttachedPolicies(PrincipalTypeRole, roleName)
	if err != nil {
		return nil, err
	}

	policies := make([]interface{}, 0, len(policyArns))
	for _, arn := range policyArns {
		if policy, err := store.Policies().Get(arn); err == nil {
			policies = append(policies, map[string]interface{}{
				"PolicyName": policy.PolicyName,
				"PolicyArn":  policy.Arn,
			})
		}
	}

	return map[string]interface{}{
		"AttachedPolicies": policies,
		"IsTruncated":      false,
	}, nil
}
