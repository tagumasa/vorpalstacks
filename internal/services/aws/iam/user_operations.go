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

// GetUser retrieves an IAM user by its name.
// UserName is required.
// Returns an error if the user does not exist.
func (s *IAMService) GetUser(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	userName := request.GetStringParam(req.Parameters, "UserName")
	userName, err := resolveUserName(reqCtx, userName)
	if err != nil {
		return nil, err
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	user, err := s.getUserCore(store, userName)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"User": s.userToResponse(reqCtx, user),
	}, nil
}

// CreateUser creates a new IAM user.
// UserName is required and must not be empty.
// Path defaults to "/" if not specified.
// Tags are optional.
func (s *IAMService) CreateUser(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	input := &CreateUserInput{
		UserName:               request.GetStringParam(req.Parameters, "UserName"),
		Path:                   request.GetStringParam(req.Parameters, "Path"),
		PermissionsBoundaryArn: request.GetStringParam(req.Parameters, "PermissionsBoundary"),
		Tags:                   tags.ParseTagsWithQueryFallback(req.Parameters, "Tags"),
	}
	user, err := s.createUserCore(store, input)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"User": s.userToResponse(reqCtx, user),
	}, nil
}

// DeleteUser deletes an IAM user by its name.
// UserName is required.
// Returns an error if the user has MFA devices, access keys, login profile, or attached policies.
// Also removes the user from all groups and deletes inline policies.
func (s *IAMService) DeleteUser(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	input := &DeleteUserInput{
		UserName: request.GetStringParam(req.Parameters, "UserName"),
	}
	if err := s.deleteUserCore(store, input); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// UpdateUser updates the path and/or name of an IAM user.
// UserName is required.
// NewPath and NewUserName are optional parameters to update.
// If NewUserName is provided, migrates all user resources to the new name.
func (s *IAMService) UpdateUser(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	input := &UpdateUserInput{
		UserName:    request.GetStringParam(req.Parameters, "UserName"),
		NewPath:     request.GetStringParam(req.Parameters, "NewPath"),
		NewUserName: request.GetStringParam(req.Parameters, "NewUserName"),
	}
	user, err := s.updateUserCore(store, input)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"User": s.userToResponse(reqCtx, user),
	}, nil
}

// ListUsers lists IAM users.
// PathPrefix filters by path prefix.
// Supports pagination via Marker and MaxItems.
func (s *IAMService) ListUsers(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	pathPrefix := request.GetStringParam(req.Parameters, "PathPrefix")
	marker := request.GetStringParam(req.Parameters, "Marker")
	maxItems := pagination.GetMaxItems(req.Parameters, pagination.DefaultMaxItems)

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	result, err := s.listUsersCore(store, pathPrefix, marker, maxItems)
	if err != nil {
		return nil, err
	}

	users := make([]interface{}, len(result.Users))
	for i, user := range result.Users {
		users[i] = s.userToResponse(reqCtx, user)
	}

	response := map[string]interface{}{
		"Users":       users,
		"IsTruncated": result.IsTruncated,
	}

	if result.Marker != "" {
		response["Marker"] = result.Marker
	}

	return response, nil
}

var userTagOps = tagOps[*iamstore.User]{
	paramName:  "UserName",
	emptyErr:   NewValidationError("UserName"),
	notFoundFn: func(n string) error { return NewNoSuchUserError(n) },
	getFn:      func(s *iamstore.IAMStore, n string) (*iamstore.User, error) { return s.Users().Get(n) },
	putFn:      func(s *iamstore.IAMStore, r *iamstore.User) error { return s.Users().Put(r) },
	tagsFn:     func(r *iamstore.User) *[]types.Tag { return &r.Tags },
}

// TagUser adds tags to an IAM user.
// UserName is required.
// Tags are provided as a list of key-value pairs.
func (s *IAMService) TagUser(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return tagResource(ctx, s, reqCtx, req, userTagOps)
}

// UntagUser removes tags from an IAM user.
// UserName is required.
// TagKeys specifies which tags to remove.
func (s *IAMService) UntagUser(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return untagResource(ctx, s, reqCtx, req, userTagOps)
}

// ListUserTags lists the tags attached to an IAM user.
// UserName is required.
func (s *IAMService) ListUserTags(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return listResourceTags(ctx, s, reqCtx, req, userTagOps)
}

// PutUserPermissionsBoundary sets the permissions boundary for an IAM user.
// UserName is required.
// PermissionsBoundary is the ARN of a managed policy to use as the permissions boundary.
func (s *IAMService) PutUserPermissionsBoundary(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	userName := request.GetStringParam(req.Parameters, "UserName")
	if userName == "" {
		return nil, NewValidationError("UserName")
	}

	permissionsBoundary := request.GetStringParam(req.Parameters, "PermissionsBoundary")
	if permissionsBoundary == "" {
		return nil, NewValidationError("PermissionsBoundary")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	user, err := s.getUserCore(store, userName)
	if err != nil {
		return nil, err
	}
	if err := putUserPermissionsBoundaryCore(store, user, permissionsBoundary); err != nil {
		return nil, err
	}
	return response.EmptyResponse(), nil
}

// DeleteUserPermissionsBoundary removes the permissions boundary from an IAM user.
// UserName is required.
func (s *IAMService) DeleteUserPermissionsBoundary(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	userName := request.GetStringParam(req.Parameters, "UserName")
	if userName == "" {
		return nil, NewValidationError("UserName")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := s.deleteUserPermissionsBoundaryCore(store, userName); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// GetAccountAuthorizationDetails retrieves information about all IAM users, groups,
// roles, and policies in the account, including their relationships.
// Supports pagination via Marker and MaxItems.
func (s *IAMService) GetAccountAuthorizationDetails(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	filterParam := request.GetStringList(req.Parameters, "Filter")
	filters := make(map[string]bool)
	for _, f := range filterParam {
		filters[f] = true
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	input := &AccountAuthorizationDetailsInput{
		Filters:  filters,
		Marker:   pagination.GetMarker(req.Parameters),
		MaxItems: pagination.GetMaxItems(req.Parameters, pagination.DefaultMaxItems),
	}
	return s.getAccountAuthorizationDetailsCore(store, input)
}

func (s *IAMService) userToResponse(reqCtx *request.RequestContext, user *iamstore.User) map[string]interface{} {
	resp := map[string]interface{}{
		"UserId":     user.ID,
		"Path":       user.Path,
		"UserName":   user.UserName,
		"Arn":        user.Arn,
		"CreateDate": user.CreateDate.Format(timeutils.ISO8601SimpleFormat),
	}

	if user.PasswordLastUsed != nil {
		resp["PasswordLastUsed"] = user.PasswordLastUsed.Format(timeutils.ISO8601SimpleFormat)
	}

	if user.PermissionsBoundary != nil {
		resp["PermissionsBoundary"] = map[string]interface{}{
			"PermissionsBoundaryType": user.PermissionsBoundary.PermissionsBoundaryType,
			"PermissionsBoundaryArn":  user.PermissionsBoundary.PermissionsBoundaryArn,
		}
	}

	if tags := tags.ToResponse(user.Tags); tags != nil {
		resp["Tags"] = tags
	}

	return resp
}
