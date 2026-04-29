// Package iam provides IAM service operations for vorpalstacks.
package iam

import (
	"context"
	"errors"
	"regexp"

	"vorpalstacks/internal/common/pagination"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	"vorpalstacks/internal/common/tags"
	iamstore "vorpalstacks/internal/store/aws/iam"
	"vorpalstacks/internal/utils/aws/types"
	"vorpalstacks/internal/utils/timeutils"
)

var roleNamePattern = regexp.MustCompile(`^[\w+=,.@-]{1,64}$`)

// CreateRole creates a new IAM role.
// RoleName is required and must not be empty.
// Path defaults to "/" if not specified.
// AssumeRolePolicyDocument is the trust policy that controls who can assume the role.
// Description, MaxSessionDuration, and Tags are optional.
func (s *IAMService) CreateRole(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	roleName := request.GetStringParam(req.Parameters, "RoleName")
	if roleName == "" {
		return nil, NewInvalidInputError("RoleName", "cannot be empty")
	}
	if !roleNamePattern.MatchString(roleName) {
		return nil, NewInvalidInputError("RoleName", "must be 1 to 64 alphanumeric characters or any of +=,.@-_")
	}

	path := request.GetStringParam(req.Parameters, "Path")
	if path == "" {
		path = "/"
	}

	assumeRolePolicyDocument := request.GetStringParam(req.Parameters, "AssumeRolePolicyDocument")
	if assumeRolePolicyDocument == "" {
		return nil, ErrMalformedPolicyDocument
	}
	if !validatePolicyDocument(assumeRolePolicyDocument) {
		return nil, ErrMalformedPolicyDocument
	}
	description := request.GetStringParam(req.Parameters, "Description")
	maxSessionDuration := request.GetIntParam(req.Parameters, "MaxSessionDuration")
	if maxSessionDuration == 0 {
		maxSessionDuration = 3600
	}
	if maxSessionDuration < 900 || maxSessionDuration > 43200 {
		return nil, NewInvalidInputError("MaxSessionDuration", "must be between 900 and 43200 seconds")
	}

	newTags := tags.ParseTagsWithQueryFallback(req.Parameters, "Tags")

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	role, err := store.Roles().Create(roleName, path, s.accountID, assumeRolePolicyDocument, description, maxSessionDuration, newTags)
	if err != nil {
		if errors.Is(err, iamstore.ErrRoleAlreadyExists) {
			return nil, NewRoleAlreadyExistsError(roleName)
		}
		return nil, err
	}

	return map[string]interface{}{
		"Role": s.roleToResponse(reqCtx, role),
	}, nil
}

// GetRole retrieves an IAM role by its name.
// Returns an error if the role does not exist.
func (s *IAMService) GetRole(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	roleName := request.GetStringParam(req.Parameters, "RoleName")
	if roleName == "" {
		return nil, ErrNoSuchRole
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	role, err := store.Roles().Get(roleName)
	if err != nil {
		return nil, NewNoSuchRoleError(roleName)
	}

	return map[string]interface{}{
		"Role": s.roleToResponse(reqCtx, role),
	}, nil
}

// UpdateRole updates the description and maximum session duration of an IAM role.
// RoleName is required.
// Description and MaxSessionDuration are optional parameters to update.
func (s *IAMService) UpdateRole(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	roleName := request.GetStringParam(req.Parameters, "RoleName")
	if roleName == "" {
		return nil, ErrNoSuchRole
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	role, err := store.Roles().Get(roleName)
	if err != nil {
		return nil, NewNoSuchRoleError(roleName)
	}

	description := request.GetStringParam(req.Parameters, "Description")
	if description != "" {
		role.Description = description
	}

	maxSessionDuration := request.GetIntParam(req.Parameters, "MaxSessionDuration")
	if maxSessionDuration > 0 {
		if maxSessionDuration < 900 || maxSessionDuration > 43200 {
			return nil, NewInvalidInputError("MaxSessionDuration", "must be between 900 and 43200 seconds")
		}
		role.MaxSessionDuration = maxSessionDuration
	}

	if err := store.Roles().Put(role); err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"Role": s.roleToResponse(reqCtx, role),
	}, nil
}

// UpdateRoleDescription updates the description of an IAM role.
// RoleName and Description are required.
func (s *IAMService) UpdateRoleDescription(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	roleName := request.GetStringParam(req.Parameters, "RoleName")
	if roleName == "" {
		return nil, ErrNoSuchRole
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	role, err := store.Roles().Get(roleName)
	if err != nil {
		return nil, NewNoSuchRoleError(roleName)
	}

	description := request.GetStringParam(req.Parameters, "Description")
	role.Description = description

	if err := store.Roles().Put(role); err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"Role": s.roleToResponse(reqCtx, role),
	}, nil
}

// DeleteRole deletes an IAM role by its name.
// RoleName is required.
// Returns an error if the role is attached to instance profiles.
// Also deletes all inline policies and detaches attached policies.
func (s *IAMService) DeleteRole(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
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

	instanceProfiles, err := store.InstanceProfiles().ListForRole(roleName, "", 1)
	if err != nil {
		return nil, err
	}
	if len(instanceProfiles.InstanceProfiles) > 0 {
		return nil, NewDeleteRoleConflictError("Cannot delete entity, must remove role from instance profile first.")
	}

	inlinePolicies, err := store.InlinePolicies().List(PrincipalTypeRole, roleName)
	if err != nil {
		return nil, err
	}
	if len(inlinePolicies) > 0 {
		return nil, NewDeleteRoleConflictError("Cannot delete entity, must delete policies first.")
	}

	attachedPolicies, err := store.AttachedPolicies().ListAttachedPolicies(PrincipalTypeRole, roleName)
	if err != nil {
		return nil, err
	}
	if len(attachedPolicies) > 0 {
		return nil, NewDeleteRoleConflictError("Cannot delete entity, must detach policies first.")
	}

	if err := store.Roles().Delete(roleName); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// ListRoles lists IAM roles.
// PathPrefix filters by path prefix.
// Supports pagination via Marker and MaxItems.
func (s *IAMService) ListRoles(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	pathPrefix := request.GetStringParam(req.Parameters, "PathPrefix")
	marker := request.GetStringParam(req.Parameters, "Marker")
	maxItems := pagination.GetMaxItems(req.Parameters, pagination.DefaultMaxItems)

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	result, err := store.Roles().List(pathPrefix, marker, maxItems)
	if err != nil {
		return nil, err
	}

	roles := make([]interface{}, len(result.Roles))
	for i, role := range result.Roles {
		roles[i] = s.roleToResponse(reqCtx, role)
	}

	response := map[string]interface{}{
		"Roles":       roles,
		"IsTruncated": result.IsTruncated,
	}

	if result.Marker != "" {
		response["Marker"] = result.Marker
	}

	return response, nil
}

// UpdateAssumeRolePolicy updates the trust policy document of an IAM role.
// RoleName is required.
// PolicyDocument must be a valid JSON policy document.
func (s *IAMService) UpdateAssumeRolePolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	roleName := request.GetStringParam(req.Parameters, "RoleName")
	if roleName == "" {
		return nil, ErrNoSuchRole
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	role, err := store.Roles().Get(roleName)
	if err != nil {
		return nil, NewNoSuchRoleError(roleName)
	}

	policyDocument := request.GetStringParam(req.Parameters, "PolicyDocument")
	if !validatePolicyDocument(policyDocument) {
		return nil, ErrMalformedPolicyDocument
	}
	role.AssumeRolePolicyDocument = policyDocument

	if err := store.Roles().Put(role); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

var roleTagOps = tagOps[*iamstore.Role]{
	paramName:  "RoleName",
	emptyErr:   ErrNoSuchRole,
	notFoundFn: func(n string) error { return NewNoSuchRoleError(n) },
	getFn:      func(s *iamstore.IAMStore, n string) (*iamstore.Role, error) { return s.Roles().Get(n) },
	putFn:      func(s *iamstore.IAMStore, r *iamstore.Role) error { return s.Roles().Put(r) },
	tagsFn:     func(r *iamstore.Role) *[]types.Tag { return &r.Tags },
}

// TagRole adds tags to an IAM role.
// RoleName is required.
// Tags are provided as a list of key-value pairs.
func (s *IAMService) TagRole(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return tagResource(ctx, s, reqCtx, req, roleTagOps)
}

// UntagRole removes tags from an IAM role.
// RoleName is required.
// TagKeys specifies which tags to remove.
func (s *IAMService) UntagRole(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return untagResource(ctx, s, reqCtx, req, roleTagOps)
}

// ListRoleTags lists the tags attached to an IAM role.
// RoleName is required.
func (s *IAMService) ListRoleTags(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return listResourceTags(ctx, s, reqCtx, req, roleTagOps)
}

func (s *IAMService) roleToResponse(reqCtx *request.RequestContext, role *iamstore.Role) map[string]interface{} {
	resp := map[string]interface{}{
		"RoleId":             role.ID,
		"Path":               role.Path,
		"RoleName":           role.RoleName,
		"Arn":                role.Arn,
		"CreateDate":         role.CreateDate.Format(timeutils.ISO8601SimpleFormat),
		"MaxSessionDuration": role.MaxSessionDuration,
	}

	if role.AssumeRolePolicyDocument != "" {
		resp["AssumeRolePolicyDocument"] = role.AssumeRolePolicyDocument
	}

	if role.Description != "" {
		resp["Description"] = role.Description
	}

	if role.PermissionsBoundary != nil {
		resp["PermissionsBoundary"] = map[string]interface{}{
			"PermissionsBoundaryType": role.PermissionsBoundary.PermissionsBoundaryType,
			"PermissionsBoundaryArn":  role.PermissionsBoundary.PermissionsBoundaryArn,
		}
	}

	if tags := tags.ToResponse(role.Tags); tags != nil {
		resp["Tags"] = tags
	}

	if role.RoleLastUsed != nil {
		lastUsed := map[string]interface{}{}
		if role.RoleLastUsed.LastUsedDate != nil {
			lastUsed["LastUsedDate"] = role.RoleLastUsed.LastUsedDate.Format(timeutils.ISO8601SimpleFormat)
		}
		if role.RoleLastUsed.Region != "" {
			lastUsed["Region"] = role.RoleLastUsed.Region
		}
		resp["RoleLastUsed"] = lastUsed
	}

	return resp
}

// ListInstanceProfilesForRole lists instance profiles associated with an IAM role.
// RoleName is required.
// Supports pagination via Marker and MaxItems.
func (s *IAMService) ListInstanceProfilesForRole(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
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

	marker := request.GetStringParam(req.Parameters, "Marker")
	maxItems := pagination.GetMaxItems(req.Parameters, pagination.DefaultMaxItems)

	result, err := store.InstanceProfiles().ListForRole(roleName, marker, maxItems)
	if err != nil {
		return nil, err
	}

	profiles := make([]interface{}, len(result.InstanceProfiles))
	for i, profile := range result.InstanceProfiles {
		profiles[i] = s.instanceProfileToResponseWithRoles(reqCtx, profile)
	}

	response := map[string]interface{}{
		"InstanceProfiles": profiles,
		"IsTruncated":      result.IsTruncated,
	}

	if result.Marker != "" {
		response["Marker"] = result.Marker
	}

	return response, nil
}

// PutRolePermissionsBoundary sets the permissions boundary for an IAM role.
// RoleName is required.
// PermissionsBoundary is the ARN of a managed policy to use as the permissions boundary.
func (s *IAMService) PutRolePermissionsBoundary(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	roleName := request.GetStringParam(req.Parameters, "RoleName")
	if roleName == "" {
		return nil, ErrNoSuchRole
	}

	permissionsBoundary := request.GetStringParam(req.Parameters, "PermissionsBoundary")
	if permissionsBoundary == "" {
		return nil, ErrNoSuchPolicy
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	role, err := store.Roles().Get(roleName)
	if err != nil {
		return nil, NewNoSuchRoleError(roleName)
	}

	if !store.Policies().Exists(permissionsBoundary) {
		return nil, NewNoSuchPolicyError(permissionsBoundary)
	}

	role.PermissionsBoundary = &iamstore.PermissionsBoundary{
		PermissionsBoundaryType: "Policy",
		PermissionsBoundaryArn:  permissionsBoundary,
	}

	if err := store.Roles().Put(role); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// DeleteRolePermissionsBoundary removes the permissions boundary from an IAM role.
// RoleName is required.
func (s *IAMService) DeleteRolePermissionsBoundary(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	roleName := request.GetStringParam(req.Parameters, "RoleName")
	if roleName == "" {
		return nil, ErrNoSuchRole
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	role, err := store.Roles().Get(roleName)
	if err != nil {
		return nil, NewNoSuchRoleError(roleName)
	}

	role.PermissionsBoundary = nil

	if err := store.Roles().Put(role); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}
