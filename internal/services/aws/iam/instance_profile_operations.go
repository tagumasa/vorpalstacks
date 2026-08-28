// Package iam provides IAM service operations for vorpalstacks.
package iam

import (
	"context"

	"vorpalstacks/internal/common/pagination"
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
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	input := &CreateInstanceProfileInput{
		InstanceProfileName: request.GetStringParam(req.Parameters, "InstanceProfileName"),
		Path:                request.GetStringParam(req.Parameters, "Path"),
		Tags:                tags.ParseTagsWithQueryFallback(req.Parameters, "Tags"),
	}
	profile, err := s.createInstanceProfileCore(store, input)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"InstanceProfile": instanceProfileToResponse(profile),
	}, nil
}

// GetInstanceProfile retrieves information about an instance profile,
// including the roles attached to the instance profile.
func (s *IAMService) GetInstanceProfile(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	result, err := s.getInstanceProfileCore(store, request.GetStringParam(req.Parameters, "InstanceProfileName"))
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"InstanceProfile": s.instanceProfileToResponseWithRoles(result.Profile, result.Roles),
	}, nil
}

// DeleteInstanceProfile deletes an instance profile.
// Returns an error if roles are still attached to the instance profile.
func (s *IAMService) DeleteInstanceProfile(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := s.deleteInstanceProfileCore(store, request.GetStringParam(req.Parameters, "InstanceProfileName")); err != nil {
		return nil, err
	}
	return response.EmptyResponse(), nil
}

// ListInstanceProfiles lists the instance profiles in the account.
// Supports filtering by path prefix and pagination via marker.
func (s *IAMService) ListInstanceProfiles(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	pathPrefix := request.GetStringParam(req.Parameters, "PathPrefix")
	marker := request.GetStringParam(req.Parameters, "Marker")
	maxItems := pagination.GetMaxItems(req.Parameters, pagination.DefaultMaxItems)

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	result, err := s.listInstanceProfilesCore(store, pathPrefix, marker, maxItems)
	if err != nil {
		return nil, err
	}

	profiles := make([]interface{}, len(result.Profiles))
	for i, entry := range result.Profiles {
		profiles[i] = s.instanceProfileToResponseWithRoles(entry.Profile, entry.Roles)
	}

	resp := map[string]interface{}{
		"InstanceProfiles": profiles,
		"IsTruncated":      result.IsTruncated,
	}

	if result.Marker != "" {
		resp["Marker"] = result.Marker
	}

	return resp, nil
}

// AddRoleToInstanceProfile adds a role to an instance profile.
// The role and instance profile must already exist.  AWS enforces a
// maximum of one role per instance profile; the limit is enforced
// atomically inside the store layer.
func (s *IAMService) AddRoleToInstanceProfile(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := s.addRoleToInstanceProfileCore(store,
		request.GetStringParam(req.Parameters, "InstanceProfileName"),
		request.GetStringParam(req.Parameters, "RoleName")); err != nil {
		return nil, err
	}
	return response.EmptyResponse(), nil
}

// RemoveRoleFromInstanceProfile removes a role from an instance profile.
// Returns an error if the role is not attached to the instance profile.
func (s *IAMService) RemoveRoleFromInstanceProfile(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := s.removeRoleFromInstanceProfileCore(store,
		request.GetStringParam(req.Parameters, "InstanceProfileName"),
		request.GetStringParam(req.Parameters, "RoleName")); err != nil {
		return nil, err
	}
	return response.EmptyResponse(), nil
}

var instanceProfileTagOps = tagOps[*iamstore.InstanceProfile]{
	paramName:  "InstanceProfileName",
	emptyErr:   ErrNoSuchInstanceProfile,
	notFoundFn: func(n string) error { return NewNoSuchInstanceProfileError(n) },
	getFn: func(s *iamstore.IAMStore, n string) (*iamstore.InstanceProfile, error) {
		return s.InstanceProfiles().Get(n)
	},
	putFn:  func(s *iamstore.IAMStore, r *iamstore.InstanceProfile) error { return s.InstanceProfiles().Put(r) },
	tagsFn: func(r *iamstore.InstanceProfile) *[]tags.Tag { return &r.Tags },
}

// ListInstanceProfileTags lists the tags attached to an instance profile.
func (s *IAMService) ListInstanceProfileTags(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return listResourceTags(ctx, s, reqCtx, req, instanceProfileTagOps)
}

// TagInstanceProfile adds tags to an instance profile.
func (s *IAMService) TagInstanceProfile(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return tagResource(ctx, s, reqCtx, req, instanceProfileTagOps)
}

// UntagInstanceProfile removes tags from an instance profile.
func (s *IAMService) UntagInstanceProfile(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return untagResource(ctx, s, reqCtx, req, instanceProfileTagOps)
}

func instanceProfileToResponse(profile *iamstore.InstanceProfile) map[string]interface{} {
	resp := map[string]interface{}{
		"InstanceProfileId":   profile.ID,
		"Path":                profile.Path,
		"InstanceProfileName": profile.InstanceProfileName,
		"Arn":                 profile.Arn,
		"CreateDate":          profile.CreateDate.Format(timeutils.ISO8601SimpleFormat),
	}

	if tags := tags.ToResponse(profile.Tags); tags != nil {
		resp["Tags"] = tags
	}

	return resp
}

func (s *IAMService) instanceProfileToResponseWithRoles(profile *iamstore.InstanceProfile, roles []*iamstore.Role) map[string]interface{} {
	resp := instanceProfileToResponse(profile)

	roleList := make([]interface{}, 0, len(roles))
	for _, role := range roles {
		roleList = append(roleList, roleToResponse(role))
	}
	resp["Roles"] = roleList

	return resp
}
