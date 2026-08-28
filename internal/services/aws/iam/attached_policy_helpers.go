package iam

import (
	"context"

	"vorpalstacks/internal/common/pagination"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
)

// attachOps encapsulates the principal-type-specific wire parameters needed
// for attach, detach, and list-attached-policy operations. This eliminates
// the copy-paste duplication across User/Group/Role variants.
type attachOps struct {
	paramName     string
	principalType string
	emptyErr      error
}

var (
	userAttachOps = attachOps{
		paramName:     "UserName",
		principalType: PrincipalTypeUser,
		emptyErr:      ErrNoSuchUser,
	}
	groupAttachOps = attachOps{
		paramName:     "GroupName",
		principalType: PrincipalTypeGroup,
		emptyErr:      ErrNoSuchGroup,
	}
	roleAttachOps = attachOps{
		paramName:     "RoleName",
		principalType: PrincipalTypeRole,
		emptyErr:      ErrNoSuchRole,
	}
)

// attachPolicy attaches a managed policy to a principal (user, group, or role).
func attachPolicy(ctx context.Context, s *IAMService, reqCtx *request.RequestContext, req *request.ParsedRequest, ops attachOps) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	input := &AttachPolicyInput{
		PrincipalType: ops.principalType,
		PrincipalName: request.GetStringParam(req.Parameters, ops.paramName),
		PolicyArn:     request.GetStringParam(req.Parameters, "PolicyArn"),
	}
	if err := s.attachPolicyCore(store, input); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// detachPolicy detaches a managed policy from a principal.
func detachPolicy(ctx context.Context, s *IAMService, reqCtx *request.RequestContext, req *request.ParsedRequest, ops attachOps) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	input := &AttachPolicyInput{
		PrincipalType: ops.principalType,
		PrincipalName: request.GetStringParam(req.Parameters, ops.paramName),
		PolicyArn:     request.GetStringParam(req.Parameters, "PolicyArn"),
	}
	if err := s.detachPolicyCore(store, input); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// listAttachedPolicies lists the managed policies attached to a principal.
// Supports pagination via Marker and MaxItems, and optional PathPrefix
// filtering.
func listAttachedPolicies(ctx context.Context, s *IAMService, reqCtx *request.RequestContext, req *request.ParsedRequest, ops attachOps) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	result, err := s.listAttachedPoliciesCore(store,
		ops.principalType,
		request.GetStringParam(req.Parameters, ops.paramName),
		request.GetStringParam(req.Parameters, "PathPrefix"),
		request.GetStringParam(req.Parameters, "Marker"),
		pagination.GetMaxItems(req.Parameters, pagination.DefaultMaxItems),
	)
	if err != nil {
		return nil, err
	}

	policies := make([]interface{}, len(result.Entries))
	for i, e := range result.Entries {
		policies[i] = map[string]interface{}{
			"PolicyName": e.PolicyName,
			"PolicyArn":  e.PolicyArn,
		}
	}

	resp := map[string]interface{}{
		"AttachedPolicies": policies,
		"IsTruncated":      result.IsTruncated,
	}

	if result.IsTruncated && result.NextMarker != "" {
		resp["Marker"] = result.NextMarker
	}

	return resp, nil
}
