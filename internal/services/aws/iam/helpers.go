// Package iam provides IAM service operations for vorpalstacks.
package iam

import (
	"context"
	"strconv"

	"vorpalstacks/internal/common/pagination"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	"vorpalstacks/internal/common/tags"
	iamstore "vorpalstacks/internal/store/aws/iam"
)

// resolveUserName returns userName if non-empty, otherwise defaults to the
// caller's IAM principal name.  Per AWS spec, several IAM operations allow
// UserName to be omitted, in which case it defaults to the authenticated
// caller.  If the caller is not an IAM user (e.g. anonymous), an error is
// returned.
func resolveUserName(reqCtx *request.RequestContext, userName string) (string, error) {
	if userName != "" {
		return userName, nil
	}
	if reqCtx.PrincipalType == request.PrincipalTypeUser && reqCtx.Principal != "" {
		return reqCtx.Principal, nil
	}
	return "", ErrNoSuchUser
}

type tagOps[T any] struct {
	paramName  string
	emptyErr   error
	notFoundFn func(string) error
	getFn      func(*iamstore.IAMStore, string) (T, error)
	putFn      func(*iamstore.IAMStore, T) error
	tagsFn     func(T) *[]tags.Tag
}

func tagResource[T any](ctx context.Context, s *IAMService, reqCtx *request.RequestContext, req *request.ParsedRequest, ops tagOps[T]) (interface{}, error) {
	name := request.GetStringParam(req.Parameters, ops.paramName)
	if name == "" {
		return nil, ops.emptyErr
	}
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	res, err := ops.getFn(store, name)
	if err != nil {
		return nil, ops.notFoundFn(name)
	}
	currentTags := ops.tagsFn(res)
	newTags := tags.ParseTagsWithQueryFallback(req.Parameters, "Tags")
	if err := validateTagEntries(newTags); err != nil {
		return nil, err
	}
	merged := tags.Apply(*currentTags, newTags)
	if len(merged) > tags.MaxTagsPerResource {
		return nil, NewInvalidInputError("Tags", "exceeds maximum of "+strconv.Itoa(tags.MaxTagsPerResource)+" tags per resource")
	}
	*currentTags = merged
	if err := ops.putFn(store, res); err != nil {
		return nil, err
	}
	return response.EmptyResponse(), nil
}

func untagResource[T any](ctx context.Context, s *IAMService, reqCtx *request.RequestContext, req *request.ParsedRequest, ops tagOps[T]) (interface{}, error) {
	name := request.GetStringParam(req.Parameters, ops.paramName)
	if name == "" {
		return nil, ops.emptyErr
	}
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	res, err := ops.getFn(store, name)
	if err != nil {
		return nil, ops.notFoundFn(name)
	}
	currentTags := ops.tagsFn(res)
	*currentTags = tags.RemoveByTagKeys(*currentTags, tags.ParseTagKeysWithQueryFallback(req.Parameters, "TagKeys"))
	if err := ops.putFn(store, res); err != nil {
		return nil, err
	}
	return response.EmptyResponse(), nil
}

func listResourceTags[T any](ctx context.Context, s *IAMService, reqCtx *request.RequestContext, req *request.ParsedRequest, ops tagOps[T]) (interface{}, error) {
	name := request.GetStringParam(req.Parameters, ops.paramName)
	if name == "" {
		return nil, ops.emptyErr
	}
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	res, err := ops.getFn(store, name)
	if err != nil {
		return nil, ops.notFoundFn(name)
	}
	marker := request.GetStringParam(req.Parameters, "Marker")
	maxItems := pagination.GetMaxItems(req.Parameters, pagination.DefaultMaxItems)
	paged := pagination.PaginateSlice(tags.ToResponse(*ops.tagsFn(res)), marker, maxItems, func(item map[string]interface{}) string {
		if key, ok := item["Key"].(string); ok {
			return key
		}
		return ""
	})
	resp := map[string]interface{}{
		"Tags":        paged.Items,
		"IsTruncated": paged.IsTruncated,
	}
	if paged.NextMarker != "" {
		resp["Marker"] = paged.NextMarker
	}
	return resp, nil
}

// buildAttachedManagedPolicies returns the attached managed policy list for
// a principal, collapsing the copy-paste across User/Group/Role in
// GetAccountAuthorizationDetails.
func buildAttachedManagedPolicies(store *iamstore.IAMStore, principalType, principalName string) []interface{} {
	arns, _ := store.AttachedPolicies().ListAttachedPolicies(principalType, principalName)
	policies := make([]interface{}, 0, len(arns))
	for _, arn := range arns {
		if p, err := store.Policies().Get(arn); err == nil {
			policies = append(policies, map[string]interface{}{
				"PolicyName": p.PolicyName,
				"PolicyArn":  p.Arn,
			})
		}
	}
	return policies
}

// buildInlinePolicyList returns the inline policy name list for a principal.
func buildInlinePolicyList(store *iamstore.IAMStore, principalType, principalName string) []interface{} {
	names, _ := store.InlinePolicies().List(principalType, principalName)
	list := make([]interface{}, 0, len(names))
	for _, pn := range names {
		list = append(list, map[string]interface{}{
			"PolicyName": pn,
		})
	}
	return list
}

// listAllUsers paginates through all users matching pathPrefix.
func listAllUsers(store *iamstore.IAMStore, pathPrefix string) ([]*iamstore.User, error) {
	var all []*iamstore.User
	marker := ""
	for {
		result, err := store.Users().List(pathPrefix, marker, 1000)
		if err != nil {
			return nil, err
		}
		all = append(all, result.Users...)
		if !result.IsTruncated {
			break
		}
		marker = result.Marker
	}
	return all, nil
}

// listAllGroups paginates through all groups matching pathPrefix.
func listAllGroups(store *iamstore.IAMStore, pathPrefix string) ([]*iamstore.Group, error) {
	var all []*iamstore.Group
	marker := ""
	for {
		result, err := store.Groups().List(pathPrefix, marker, 1000)
		if err != nil {
			return nil, err
		}
		all = append(all, result.Groups...)
		if !result.IsTruncated {
			break
		}
		marker = result.Marker
	}
	return all, nil
}

// listAllRoles paginates through all roles matching pathPrefix.
func listAllRoles(store *iamstore.IAMStore, pathPrefix string) ([]*iamstore.Role, error) {
	var all []*iamstore.Role
	marker := ""
	for {
		result, err := store.Roles().List(pathPrefix, marker, 1000)
		if err != nil {
			return nil, err
		}
		all = append(all, result.Roles...)
		if !result.IsTruncated {
			break
		}
		marker = result.Marker
	}
	return all, nil
}

// listAllPolicies paginates through all policies matching the given filters.
func listAllPolicies(store *iamstore.IAMStore, scope, pathPrefix string, onlyAttached bool) ([]*iamstore.Policy, error) {
	var all []*iamstore.Policy
	marker := ""
	for {
		result, err := store.Policies().List(scope, pathPrefix, onlyAttached, marker, 1000)
		if err != nil {
			return nil, err
		}
		all = append(all, result.Policies...)
		if !result.IsTruncated {
			break
		}
		marker = result.Marker
	}
	return all, nil
}

// cascadeDeleteUser removes all resources associated with a user before deleting
// the user record. Used by the admin handler to ensure consistent cleanup.
func cascadeDeleteUser(store *iamstore.IAMStore, userName string) error {
	// Decrement permissions boundary usage count before the cascade removes
	// the user record. AWS allows deleting an entity that still has a
	// permissions boundary attached, so the policy counter must be adjusted
	// here to avoid drift. Best-effort, matching the PutUserPermissionsBoundary
	// / DeleteUserPermissionsBoundary pattern.
	if user, gErr := store.Users().Get(userName); gErr == nil {
		if user.PermissionsBoundary != nil && user.PermissionsBoundary.PermissionsBoundaryArn != "" {
			_ = store.Policies().DecrementPermissionsBoundaryUsageCount(user.PermissionsBoundary.PermissionsBoundaryArn)
		}
	}

	if store.LoginProfiles().Exists(userName) {
		if err := store.LoginProfiles().Delete(userName); err != nil {
			return err
		}
	}

	if err := store.InlinePolicies().DeleteAllForPrincipal(PrincipalTypeUser, userName); err != nil {
		return err
	}

	attachedPolicies, err := store.AttachedPolicies().ListAttachedPolicies(PrincipalTypeUser, userName)
	if err != nil {
		return err
	}
	for _, policyArn := range attachedPolicies {
		if err := store.AttachedPolicies().Detach(PrincipalTypeUser, userName, policyArn); err != nil {
			return err
		}
		if err := store.Policies().DecrementAttachmentCount(policyArn); err != nil {
			return err
		}
	}

	if err := store.UserGroups().RemoveAllGroupsForUser(userName); err != nil {
		return err
	}

	if err := store.SigningCertificates().DeleteAllForUser(userName); err != nil {
		return err
	}

	if err := store.SSHPublicKeys().DeleteAllForUser(userName); err != nil {
		return err
	}

	if err := store.ServiceSpecificCredentials().DeleteAllForUser(userName); err != nil {
		return err
	}

	if err := store.AccessKeys().DeleteByUserName(userName); err != nil {
		return err
	}

	mfaResult, err := store.MFADevices().ListForUser(userName, "", 1000)
	if err != nil {
		return err
	}
	for _, device := range mfaResult.MFADevices {
		if err := store.MFADevices().Deactivate(device.SerialNumber); err != nil {
			return err
		}
	}

	return store.Users().Delete(userName)
}

// cascadeDeleteRole removes all resources associated with a role before deleting
// the role record. Used by the admin handler to ensure consistent cleanup.
func cascadeDeleteRole(store *iamstore.IAMStore, roleName string) error {
	// Decrement permissions boundary usage count before the cascade removes
	// the role record. AWS allows deleting an entity that still has a
	// permissions boundary attached, so the policy counter must be adjusted
	// here to avoid drift. Best-effort, matching the PutRolePermissionsBoundary
	// / DeleteRolePermissionsBoundary pattern.
	if role, gErr := store.Roles().Get(roleName); gErr == nil {
		if role.PermissionsBoundary != nil && role.PermissionsBoundary.PermissionsBoundaryArn != "" {
			_ = store.Policies().DecrementPermissionsBoundaryUsageCount(role.PermissionsBoundary.PermissionsBoundaryArn)
		}
	}

	if err := store.InlinePolicies().DeleteAllForPrincipal(PrincipalTypeRole, roleName); err != nil {
		return err
	}

	attachedPolicies, err := store.AttachedPolicies().ListAttachedPolicies(PrincipalTypeRole, roleName)
	if err != nil {
		return err
	}
	for _, policyArn := range attachedPolicies {
		if err := store.AttachedPolicies().Detach(PrincipalTypeRole, roleName, policyArn); err != nil {
			return err
		}
		if err := store.Policies().DecrementAttachmentCount(policyArn); err != nil {
			return err
		}
	}

	instanceProfiles, err := store.InstanceProfiles().ListForRole(roleName, "", 1000)
	if err != nil {
		return err
	}
	for _, ip := range instanceProfiles.InstanceProfiles {
		if err := store.InstanceProfiles().RemoveRole(ip.InstanceProfileName, roleName); err != nil {
			return err
		}
	}

	return store.Roles().Delete(roleName)
}

// cascadeDeleteGroup removes all resources associated with a group before deleting
// the group record. Used by the admin handler to ensure consistent cleanup.
func cascadeDeleteGroup(store *iamstore.IAMStore, groupName string) error {
	if err := store.InlinePolicies().DeleteAllForPrincipal(PrincipalTypeGroup, groupName); err != nil {
		return err
	}

	attachedPolicies, err := store.AttachedPolicies().ListAttachedPolicies(PrincipalTypeGroup, groupName)
	if err != nil {
		return err
	}
	for _, policyArn := range attachedPolicies {
		if err := store.AttachedPolicies().Detach(PrincipalTypeGroup, groupName, policyArn); err != nil {
			return err
		}
		if err := store.Policies().DecrementAttachmentCount(policyArn); err != nil {
			return err
		}
	}

	if err := store.UserGroups().RemoveAllUsersFromGroup(groupName); err != nil {
		return err
	}

	return store.Groups().Delete(groupName)
}
