package iam

import (
	"context"

	"vorpalstacks/internal/common/pagination"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
)

// inlinePolicyOps encapsulates the principal-type-specific wire parameters
// needed for inline policy operations. This eliminates the copy-paste
// duplication across User/Group/Role variants.
type inlinePolicyOps struct {
	paramName         string
	principalType     string
	emptyErr          error
	responseParamName string
}

var (
	userInlinePolicyOps = inlinePolicyOps{
		paramName:         "UserName",
		principalType:     PrincipalTypeUser,
		emptyErr:          ErrNoSuchUser,
		responseParamName: "UserName",
	}
	groupInlinePolicyOps = inlinePolicyOps{
		paramName:         "GroupName",
		principalType:     PrincipalTypeGroup,
		emptyErr:          ErrNoSuchGroup,
		responseParamName: "GroupName",
	}
	roleInlinePolicyOps = inlinePolicyOps{
		paramName:         "RoleName",
		principalType:     PrincipalTypeRole,
		emptyErr:          ErrNoSuchRole,
		responseParamName: "RoleName",
	}
)

// putInlinePolicy creates or updates an inline policy for a principal.
func putInlinePolicy(ctx context.Context, s *IAMService, reqCtx *request.RequestContext, req *request.ParsedRequest, ops inlinePolicyOps) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	input := &PutInlinePolicyInput{
		PrincipalType:  ops.principalType,
		PrincipalName:  request.GetStringParam(req.Parameters, ops.paramName),
		PolicyName:     request.GetStringParam(req.Parameters, "PolicyName"),
		PolicyDocument: request.GetStringParam(req.Parameters, "PolicyDocument"),
	}
	if err := s.putInlinePolicyCore(store, input); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// getInlinePolicy retrieves an inline policy for a principal.
func getInlinePolicy(ctx context.Context, s *IAMService, reqCtx *request.RequestContext, req *request.ParsedRequest, ops inlinePolicyOps) (interface{}, error) {
	principalName := request.GetStringParam(req.Parameters, ops.paramName)
	policyName := request.GetStringParam(req.Parameters, "PolicyName")

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	policy, err := s.getInlinePolicyCore(store, ops.principalType, principalName, policyName)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		ops.responseParamName: principalName,
		"PolicyName":          policyName,
		"PolicyDocument":      policy.PolicyDocument,
	}, nil
}

// deleteInlinePolicy deletes an inline policy from a principal.
func deleteInlinePolicy(ctx context.Context, s *IAMService, reqCtx *request.RequestContext, req *request.ParsedRequest, ops inlinePolicyOps) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := s.deleteInlinePolicyCore(store,
		ops.principalType,
		request.GetStringParam(req.Parameters, ops.paramName),
		request.GetStringParam(req.Parameters, "PolicyName"),
	); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// listInlinePolicies lists inline policies for a principal.
// Supports pagination via Marker and MaxItems.
func listInlinePolicies(ctx context.Context, s *IAMService, reqCtx *request.RequestContext, req *request.ParsedRequest, ops inlinePolicyOps) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	result, err := s.listInlinePoliciesCore(store,
		ops.principalType,
		request.GetStringParam(req.Parameters, ops.paramName),
		request.GetStringParam(req.Parameters, "Marker"),
		pagination.GetMaxItems(req.Parameters, pagination.DefaultMaxItems),
	)
	if err != nil {
		return nil, err
	}

	resp := map[string]interface{}{
		"PolicyNames": result.Items,
		"IsTruncated": result.IsTruncated,
	}

	if result.IsTruncated && result.NextMarker != "" {
		resp["Marker"] = result.NextMarker
	}

	return resp, nil
}
