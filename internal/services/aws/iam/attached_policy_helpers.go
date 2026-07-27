package iam

import (
	"context"
	"errors"
	"strings"

	"vorpalstacks/internal/common/pagination"
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
	if !iamPolicyArnPattern.MatchString(policyArn) {
		return nil, NewInvalidInputError("PolicyArn", "must be a valid IAM policy ARN (arn:aws:iam::<account>:policy/<name>)")
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
	if !iamPolicyArnPattern.MatchString(policyArn) {
		return nil, NewInvalidInputError("PolicyArn", "must be a valid IAM policy ARN (arn:aws:iam::<account>:policy/<name>)")
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
// Supports pagination via Marker and MaxItems, and optional PathPrefix
// filtering.
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

	pathPrefix := request.GetStringParam(req.Parameters, "PathPrefix")
	marker := request.GetStringParam(req.Parameters, "Marker")
	maxItems := pagination.GetMaxItems(req.Parameters, pagination.DefaultMaxItems)

	policyArns, err := store.AttachedPolicies().ListAttachedPolicies(ops.principalType, principalName)
	if err != nil {
		return nil, err
	}

	type attachedPolicyEntry struct {
		PolicyName string
		PolicyArn  string
		Path       string
	}
	entries := make([]attachedPolicyEntry, 0, len(policyArns))
	for _, arn := range policyArns {
		policy, err := store.Policies().Get(arn)
		if err != nil {
			continue
		}
		if pathPrefix != "" && !strings.HasPrefix(policy.Path, pathPrefix) {
			continue
		}
		entries = append(entries, attachedPolicyEntry{
			PolicyName: policy.PolicyName,
			PolicyArn:  policy.Arn,
			Path:       policy.Path,
		})
	}

	result := pagination.PaginateSlice(entries, marker, maxItems, func(e attachedPolicyEntry) string {
		return e.PolicyArn
	})

	policies := make([]interface{}, len(result.Items))
	for i, e := range result.Items {
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
