package iam

import (
	"context"
	"errors"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	iamstore "vorpalstacks/internal/store/aws/iam"
)

// attachOps encapsulates the principal-type-specific parameters needed for
// attach, detach, and list-attached-policy operations. This eliminates the
// copy-paste duplication across User/Group/Role variants.
type attachOps struct {
	paramName     string
	principalType string
	emptyErr      error
	notFoundFn    func(string) error
	existsFn      func(*iamstore.IAMStore, string) bool
}

var (
	userAttachOps = attachOps{
		paramName:     "UserName",
		principalType: PrincipalTypeUser,
		emptyErr:      ErrNoSuchUser,
		notFoundFn:    func(n string) error { return NewNoSuchUserError(n) },
		existsFn:      func(st *iamstore.IAMStore, n string) bool { return st.Users().Exists(n) },
	}
	groupAttachOps = attachOps{
		paramName:     "GroupName",
		principalType: PrincipalTypeGroup,
		emptyErr:      ErrNoSuchGroup,
		notFoundFn:    func(n string) error { return NewNoSuchGroupError(n) },
		existsFn:      func(st *iamstore.IAMStore, n string) bool { return st.Groups().Exists(n) },
	}
	roleAttachOps = attachOps{
		paramName:     "RoleName",
		principalType: PrincipalTypeRole,
		emptyErr:      ErrNoSuchRole,
		notFoundFn:    func(n string) error { return NewNoSuchRoleError(n) },
		existsFn:      func(st *iamstore.IAMStore, n string) bool { return st.Roles().Exists(n) },
	}
)

// attachPolicy attaches a managed policy to a principal (user, group, or role).
func attachPolicy(ctx context.Context, s *IAMService, reqCtx *request.RequestContext, req *request.ParsedRequest, ops attachOps) (interface{}, error) {
	principalName := request.GetStringParam(req.Parameters, ops.paramName)
	policyArn := request.GetStringParam(req.Parameters, "PolicyArn")

	if principalName == "" {
		return nil, ops.emptyErr
	}
	if policyArn == "" {
		return nil, ErrNoSuchPolicy
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if !ops.existsFn(store, principalName) {
		return nil, ops.notFoundFn(principalName)
	}
	if !store.Policies().Exists(policyArn) {
		return nil, NewNoSuchPolicyError(policyArn)
	}

	if store.AttachedPolicies().IsAttached(ops.principalType, principalName, policyArn) {
		return response.EmptyResponse(), nil
	}

	if err := store.AttachedPolicies().Attach(ops.principalType, principalName, policyArn); err != nil {
		return nil, err
	}
	if err := store.Policies().IncrementAttachmentCount(policyArn); err != nil {
		if rollbackErr := store.AttachedPolicies().Detach(ops.principalType, principalName, policyArn); rollbackErr != nil {
			return nil, errors.Join(err, rollbackErr)
		}
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// detachPolicy detaches a managed policy from a principal.
func detachPolicy(ctx context.Context, s *IAMService, reqCtx *request.RequestContext, req *request.ParsedRequest, ops attachOps) (interface{}, error) {
	principalName := request.GetStringParam(req.Parameters, ops.paramName)
	policyArn := request.GetStringParam(req.Parameters, "PolicyArn")

	if principalName == "" {
		return nil, ops.emptyErr
	}
	if policyArn == "" {
		return nil, ErrNoSuchPolicy
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if !store.AttachedPolicies().IsAttached(ops.principalType, principalName, policyArn) {
		return nil, NewPolicyNotAttachedError(policyArn)
	}

	if err := store.AttachedPolicies().Detach(ops.principalType, principalName, policyArn); err != nil {
		return nil, err
	}
	if err := store.Policies().DecrementAttachmentCount(policyArn); err != nil {
		if rollbackErr := store.AttachedPolicies().Attach(ops.principalType, principalName, policyArn); rollbackErr != nil {
			return nil, errors.Join(err, rollbackErr)
		}
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// listAttachedPolicies lists the managed policies attached to a principal.
func listAttachedPolicies(ctx context.Context, s *IAMService, reqCtx *request.RequestContext, req *request.ParsedRequest, ops attachOps) (interface{}, error) {
	principalName := request.GetStringParam(req.Parameters, ops.paramName)
	if principalName == "" {
		return nil, ops.emptyErr
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if !ops.existsFn(store, principalName) {
		return nil, ops.notFoundFn(principalName)
	}

	policyArns, err := store.AttachedPolicies().ListAttachedPolicies(ops.principalType, principalName)
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
