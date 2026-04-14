// Package iam provides IAM service operations for vorpalstacks.
package iam

import (
	"context"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	"vorpalstacks/internal/common/tags"
	iamstore "vorpalstacks/internal/store/aws/iam"
	"vorpalstacks/internal/utils/timeutils"
)

// CreateInstanceProfile creates a new instance profile.
// An instance profile is a container for an IAM role that you can use to
// pass role information to an EC2 instance.
func (s *IAMService) CreateInstanceProfile(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	instanceProfileName := request.GetStringParam(req.Parameters, "InstanceProfileName")
	if instanceProfileName == "" {
		return nil, NewInvalidInputError("InstanceProfileName", "cannot be empty")
	}

	path := request.GetStringParam(req.Parameters, "Path")
	if path == "" {
		path = "/"
	}

	newTags := tags.ParseTagsWithQueryFallback(req.Parameters, "Tags")

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	profile, err := store.InstanceProfiles().Create(instanceProfileName, path, s.accountID, newTags)
	if err != nil {
		return nil, NewInstanceProfileAlreadyExistsError(instanceProfileName)
	}

	return map[string]interface{}{
		"InstanceProfile": s.instanceProfileToResponse(reqCtx, profile),
	}, nil
}

// GetInstanceProfile retrieves information about an instance profile,
// including the roles attached to the instance profile.
func (s *IAMService) GetInstanceProfile(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	instanceProfileName := request.GetStringParam(req.Parameters, "InstanceProfileName")
	if instanceProfileName == "" {
		return nil, ErrNoSuchInstanceProfile
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	profile, err := store.InstanceProfiles().Get(instanceProfileName)
	if err != nil {
		return nil, NewNoSuchInstanceProfileError(instanceProfileName)
	}

	return map[string]interface{}{
		"InstanceProfile": s.instanceProfileToResponseWithRoles(reqCtx, profile),
	}, nil
}

// DeleteInstanceProfile deletes an instance profile.
// Returns an error if roles are still attached to the instance profile.
func (s *IAMService) DeleteInstanceProfile(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	instanceProfileName := request.GetStringParam(req.Parameters, "InstanceProfileName")
	if instanceProfileName == "" {
		return nil, ErrNoSuchInstanceProfile
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	profile, err := store.InstanceProfiles().Get(instanceProfileName)
	if err != nil {
		return nil, NewNoSuchInstanceProfileError(instanceProfileName)
	}

	if len(profile.Roles) > 0 {
		return nil, NewDeleteInstanceProfileConflictError(instanceProfileName)
	}

	if err := store.InstanceProfiles().Delete(instanceProfileName); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// ListInstanceProfiles lists the instance profiles in the account.
// Supports filtering by path prefix and pagination via marker.
func (s *IAMService) ListInstanceProfiles(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	pathPrefix := request.GetStringParam(req.Parameters, "PathPrefix")
	marker := request.GetStringParam(req.Parameters, "Marker")
	maxItems := getMaxItems(req.Parameters)

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	result, err := store.InstanceProfiles().List(pathPrefix, marker, maxItems)
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

// AddRoleToInstanceProfile adds a role to an instance profile.
// The role and instance profile must already exist.
func (s *IAMService) AddRoleToInstanceProfile(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	instanceProfileName := request.GetStringParam(req.Parameters, "InstanceProfileName")
	roleName := request.GetStringParam(req.Parameters, "RoleName")

	if instanceProfileName == "" {
		return nil, ErrNoSuchInstanceProfile
	}
	if roleName == "" {
		return nil, ErrNoSuchRole
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if !store.InstanceProfiles().Exists(instanceProfileName) {
		return nil, NewNoSuchInstanceProfileError(instanceProfileName)
	}

	if !store.Roles().Exists(roleName) {
		return nil, NewNoSuchRoleError(roleName)
	}

	hasRole, err := store.InstanceProfiles().HasRole(instanceProfileName, roleName)
	if err != nil {
		return nil, err
	}
	if hasRole {
		return nil, NewRoleAlreadyInInstanceProfileError(roleName, instanceProfileName)
	}

	if err := store.InstanceProfiles().AddRole(instanceProfileName, roleName); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// RemoveRoleFromInstanceProfile removes a role from an instance profile.
// Returns an error if the role is not attached to the instance profile.
func (s *IAMService) RemoveRoleFromInstanceProfile(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	instanceProfileName := request.GetStringParam(req.Parameters, "InstanceProfileName")
	roleName := request.GetStringParam(req.Parameters, "RoleName")

	if instanceProfileName == "" {
		return nil, ErrNoSuchInstanceProfile
	}
	if roleName == "" {
		return nil, ErrNoSuchRole
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if !store.InstanceProfiles().Exists(instanceProfileName) {
		return nil, NewNoSuchInstanceProfileError(instanceProfileName)
	}

	if !store.Roles().Exists(roleName) {
		return nil, NewNoSuchRoleError(roleName)
	}

	if err := store.InstanceProfiles().RemoveRole(instanceProfileName, roleName); err != nil {
		return nil, NewRoleNotInInstanceProfileError(roleName, instanceProfileName)
	}

	return response.EmptyResponse(), nil
}

// ListInstanceProfileTags lists the tags attached to an instance profile.
func (s *IAMService) ListInstanceProfileTags(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	instanceProfileName := request.GetStringParam(req.Parameters, "InstanceProfileName")
	if instanceProfileName == "" {
		return nil, ErrNoSuchInstanceProfile
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	profile, err := store.InstanceProfiles().Get(instanceProfileName)
	if err != nil {
		return nil, NewNoSuchInstanceProfileError(instanceProfileName)
	}

	return map[string]interface{}{
		"Tags":        tagsToResponse(profile.Tags),
		"IsTruncated": false,
	}, nil
}

// TagInstanceProfile adds tags to an instance profile.
func (s *IAMService) TagInstanceProfile(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	instanceProfileName := request.GetStringParam(req.Parameters, "InstanceProfileName")
	if instanceProfileName == "" {
		return nil, ErrNoSuchInstanceProfile
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	profile, err := store.InstanceProfiles().Get(instanceProfileName)
	if err != nil {
		return nil, NewNoSuchInstanceProfileError(instanceProfileName)
	}

	profile.Tags = tags.Apply(profile.Tags, tags.ParseTagsWithQueryFallback(req.Parameters, "Tags"))

	if err := store.InstanceProfiles().Put(profile); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// UntagInstanceProfile removes tags from an instance profile.
func (s *IAMService) UntagInstanceProfile(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	instanceProfileName := request.GetStringParam(req.Parameters, "InstanceProfileName")
	if instanceProfileName == "" {
		return nil, ErrNoSuchInstanceProfile
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	profile, err := store.InstanceProfiles().Get(instanceProfileName)
	if err != nil {
		return nil, NewNoSuchInstanceProfileError(instanceProfileName)
	}

	profile.Tags = removeTags(profile.Tags, parseTagKeys(req.Parameters))

	if err := store.InstanceProfiles().Put(profile); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

func (s *IAMService) instanceProfileToResponse(reqCtx *request.RequestContext, profile *iamstore.InstanceProfile) map[string]interface{} {
	resp := map[string]interface{}{
		"InstanceProfileId":   profile.ID,
		"Path":                profile.Path,
		"InstanceProfileName": profile.InstanceProfileName,
		"Arn":                 profile.Arn,
		"CreateDate":          profile.CreateDate.Format(timeutils.ISO8601SimpleFormat),
	}

	if tags := tagsToResponse(profile.Tags); tags != nil {
		resp["Tags"] = tags
	}

	return resp
}

func (s *IAMService) instanceProfileToResponseWithRoles(reqCtx *request.RequestContext, profile *iamstore.InstanceProfile) map[string]interface{} {
	resp := s.instanceProfileToResponse(reqCtx, profile)

	store, err := s.store(reqCtx)
	if err != nil {
		return resp
	}
	roles := make([]interface{}, 0, len(profile.Roles))
	for _, roleName := range profile.Roles {
		if role, err := store.Roles().Get(roleName); err == nil {
			roles = append(roles, s.roleToResponse(reqCtx, role))
		}
	}
	resp["Roles"] = roles

	return resp
}
