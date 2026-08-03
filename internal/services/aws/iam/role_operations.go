// Package iam provides IAM service operations for vorpalstacks.
package iam

import (
	"context"

	"vorpalstacks/internal/common/pagination"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	"vorpalstacks/internal/common/tags"
	iamstore "vorpalstacks/internal/store/aws/iam"
	"vorpalstacks/internal/utils/aws/types"
	"vorpalstacks/internal/utils/timeutils"
)

// CreateRole creates a new IAM role.
// RoleName is required and must not be empty.
// Path defaults to "/" if not specified.
// AssumeRolePolicyDocument is the trust policy that controls who can assume the role.
// Description, MaxSessionDuration, and Tags are optional.
func (s *IAMService) CreateRole(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	input := &CreateRoleInput{
		RoleName:                 request.GetStringParam(req.Parameters, "RoleName"),
		Path:                     request.GetStringParam(req.Parameters, "Path"),
		AssumeRolePolicyDocument: request.GetStringParam(req.Parameters, "AssumeRolePolicyDocument"),
		Description:              request.GetStringParam(req.Parameters, "Description"),
		MaxSessionDuration:       request.GetIntParam(req.Parameters, "MaxSessionDuration"),
		PermissionsBoundaryArn:   request.GetStringParam(req.Parameters, "PermissionsBoundary"),
		Tags:                     tags.ParseTagsWithQueryFallback(req.Parameters, "Tags"),
	}
	role, err := s.createRoleCore(store, input)
	if err != nil {
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
		return nil, NewValidationError("RoleName")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	role, err := s.getRoleCore(store, roleName)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"Role": s.roleToResponse(reqCtx, role),
	}, nil
}

// UpdateRole updates the description and maximum session duration of an IAM role.
// RoleName is required.
// Description and MaxSessionDuration are optional parameters to update.
func (s *IAMService) UpdateRole(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	input := &UpdateRoleInput{
		RoleName:           request.GetStringParam(req.Parameters, "RoleName"),
		Description:        request.GetStringParam(req.Parameters, "Description"),
		MaxSessionDuration: request.GetIntParam(req.Parameters, "MaxSessionDuration"),
	}
	role, err := s.updateRoleCore(store, input)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"Role": s.roleToResponse(reqCtx, role),
	}, nil
}

// UpdateRoleDescription updates the description of an IAM role.
// RoleName and Description are required.
func (s *IAMService) UpdateRoleDescription(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	input := &UpdateRoleInput{
		RoleName:    request.GetStringParam(req.Parameters, "RoleName"),
		Description: request.GetStringParam(req.Parameters, "Description"),
	}
	role, err := s.updateRoleCore(store, input)
	if err != nil {
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
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	input := &DeleteRoleInput{
		RoleName: request.GetStringParam(req.Parameters, "RoleName"),
	}
	if err := s.deleteRoleCore(store, input); err != nil {
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
	result, err := s.listRolesCore(store, pathPrefix, marker, maxItems)
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
		return nil, NewValidationError("RoleName")
	}
	policyDocument := request.GetStringParam(req.Parameters, "PolicyDocument")

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := s.updateAssumeRolePolicyCore(store, roleName, policyDocument); err != nil {
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
		return nil, NewValidationError("RoleName")
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
		profiles[i] = s.instanceProfileToResponseWithRoles(reqCtx, profile, store)
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
		return nil, NewValidationError("RoleName")
	}

	permissionsBoundary := request.GetStringParam(req.Parameters, "PermissionsBoundary")
	if permissionsBoundary == "" {
		return nil, NewValidationError("PermissionsBoundary")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	role, err := s.getRoleCore(store, roleName)
	if err != nil {
		return nil, err
	}
	if err := putRolePermissionsBoundaryCore(store, role, permissionsBoundary); err != nil {
		return nil, err
	}
	return response.EmptyResponse(), nil
}

// DeleteRolePermissionsBoundary removes the permissions boundary from an IAM role.
// RoleName is required.
func (s *IAMService) DeleteRolePermissionsBoundary(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	roleName := request.GetStringParam(req.Parameters, "RoleName")
	if roleName == "" {
		return nil, NewValidationError("RoleName")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := s.deleteRolePermissionsBoundaryCore(store, roleName); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}
