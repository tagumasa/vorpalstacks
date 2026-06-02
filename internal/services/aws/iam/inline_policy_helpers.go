package iam

import (
	"context"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	iamstore "vorpalstacks/internal/store/aws/iam"
)

// inlinePolicyOps encapsulates the principal-type-specific parameters needed
// for inline policy operations. This eliminates the copy-paste duplication
// across User/Group/Role variants.
type inlinePolicyOps struct {
	paramName         string
	principalType     string
	emptyErr          error
	notFoundFn        func(string) error
	existsFn          func(*iamstore.IAMStore, string) bool
	responseParamName string
}

var (
	userInlinePolicyOps = inlinePolicyOps{
		paramName:         "UserName",
		principalType:     PrincipalTypeUser,
		emptyErr:          ErrNoSuchUser,
		notFoundFn:        func(n string) error { return NewNoSuchUserError(n) },
		existsFn:          func(st *iamstore.IAMStore, n string) bool { return st.Users().Exists(n) },
		responseParamName: "UserName",
	}
	groupInlinePolicyOps = inlinePolicyOps{
		paramName:         "GroupName",
		principalType:     PrincipalTypeGroup,
		emptyErr:          ErrNoSuchGroup,
		notFoundFn:        func(n string) error { return NewNoSuchGroupError(n) },
		existsFn:          func(st *iamstore.IAMStore, n string) bool { return st.Groups().Exists(n) },
		responseParamName: "GroupName",
	}
	roleInlinePolicyOps = inlinePolicyOps{
		paramName:         "RoleName",
		principalType:     PrincipalTypeRole,
		emptyErr:          ErrNoSuchRole,
		notFoundFn:        func(n string) error { return NewNoSuchRoleError(n) },
		existsFn:          func(st *iamstore.IAMStore, n string) bool { return st.Roles().Exists(n) },
		responseParamName: "RoleName",
	}
)

// putInlinePolicy creates or updates an inline policy for a principal.
func putInlinePolicy(ctx context.Context, s *IAMService, reqCtx *request.RequestContext, req *request.ParsedRequest, ops inlinePolicyOps) (interface{}, error) {
	principalName := request.GetStringParam(req.Parameters, ops.paramName)
	policyName := request.GetStringParam(req.Parameters, "PolicyName")
	policyDocument := request.GetStringParam(req.Parameters, "PolicyDocument")

	if principalName == "" {
		return nil, ops.emptyErr
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
	if !ops.existsFn(store, principalName) {
		return nil, ops.notFoundFn(principalName)
	}

	if err := store.InlinePolicies().Put(ops.principalType, principalName, policyName, policyDocument); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// getInlinePolicy retrieves an inline policy for a principal.
func getInlinePolicy(ctx context.Context, s *IAMService, reqCtx *request.RequestContext, req *request.ParsedRequest, ops inlinePolicyOps) (interface{}, error) {
	principalName := request.GetStringParam(req.Parameters, ops.paramName)
	policyName := request.GetStringParam(req.Parameters, "PolicyName")

	if principalName == "" {
		return nil, ops.emptyErr
	}
	if policyName == "" {
		return nil, ErrNoSuchPolicy
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if !ops.existsFn(store, principalName) {
		return nil, ops.notFoundFn(principalName)
	}

	policy, err := store.InlinePolicies().Get(ops.principalType, principalName, policyName)
	if err != nil {
		return nil, NewNoSuchPolicyError(policyName)
	}

	return map[string]interface{}{
		ops.responseParamName: principalName,
		"PolicyName":          policyName,
		"PolicyDocument":      policy.PolicyDocument,
	}, nil
}

// deleteInlinePolicy deletes an inline policy from a principal.
func deleteInlinePolicy(ctx context.Context, s *IAMService, reqCtx *request.RequestContext, req *request.ParsedRequest, ops inlinePolicyOps) (interface{}, error) {
	principalName := request.GetStringParam(req.Parameters, ops.paramName)
	policyName := request.GetStringParam(req.Parameters, "PolicyName")

	if principalName == "" {
		return nil, ops.emptyErr
	}
	if policyName == "" {
		return nil, ErrNoSuchPolicy
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if !ops.existsFn(store, principalName) {
		return nil, ops.notFoundFn(principalName)
	}
	if !store.InlinePolicies().Exists(ops.principalType, principalName, policyName) {
		return nil, NewNoSuchPolicyError(policyName)
	}

	if err := store.InlinePolicies().Delete(ops.principalType, principalName, policyName); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// listInlinePolicies lists inline policies for a principal.
func listInlinePolicies(ctx context.Context, s *IAMService, reqCtx *request.RequestContext, req *request.ParsedRequest, ops inlinePolicyOps) (interface{}, error) {
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

	policyNames, err := store.InlinePolicies().List(ops.principalType, principalName)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"PolicyNames": policyNames,
		"IsTruncated": false,
	}, nil
}
