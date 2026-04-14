// Package iam provides IAM service operations for vorpalstacks.
package iam

import (
	"context"
	"errors"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	"vorpalstacks/internal/common/tags"
	"vorpalstacks/internal/core/logs"
	iamstore "vorpalstacks/internal/store/aws/iam"
	"vorpalstacks/internal/utils/timeutils"
)

// CreateGroup creates a new IAM group.
func (s *IAMService) CreateGroup(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	groupName := request.GetStringParam(req.Parameters, "GroupName")
	if groupName == "" {
		return nil, NewInvalidInputError("GroupName", "cannot be empty")
	}

	path := request.GetStringParam(req.Parameters, "Path")
	if path == "" {
		path = "/"
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	group, err := store.Groups().Create(groupName, path, s.accountID)
	if err != nil {
		if errors.Is(err, iamstore.ErrGroupAlreadyExists) {
			return nil, NewGroupAlreadyExistsError(groupName)
		}
		return nil, err
	}

	return map[string]interface{}{
		"Group": s.groupToResponse(reqCtx, group),
	}, nil
}

// GetGroup retrieves information about an IAM group.
func (s *IAMService) GetGroup(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	groupName := request.GetStringParam(req.Parameters, "GroupName")
	if groupName == "" {
		return nil, ErrNoSuchGroup
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	group, err := store.Groups().Get(groupName)
	if err != nil {
		return nil, NewNoSuchGroupError(groupName)
	}

	users, err := store.UserGroups().ListUsersInGroup(groupName)
	if err != nil {
		return nil, err
	}

	userList := make([]interface{}, 0, len(users))
	for _, userName := range users {
		if user, err := store.Users().Get(userName); err == nil {
			userList = append(userList, s.userToResponse(reqCtx, user))
		}
	}

	return map[string]interface{}{
		"Group":       s.groupToResponse(reqCtx, group),
		"Users":       userList,
		"IsTruncated": false,
	}, nil
}

// UpdateGroup updates an IAM group.
func (s *IAMService) UpdateGroup(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	groupName := request.GetStringParam(req.Parameters, "GroupName")
	if groupName == "" {
		return nil, ErrNoSuchGroup
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	group, err := store.Groups().Get(groupName)
	if err != nil {
		return nil, NewNoSuchGroupError(groupName)
	}

	newPath := request.GetStringParam(req.Parameters, "NewPath")
	if newPath != "" {
		group.Path = newPath
		group.Arn = store.ARNBuilder().GroupARN(newPath, group.GroupName)
	}

	newGroupName := request.GetStringParam(req.Parameters, "NewGroupName")
	if newGroupName != "" && newGroupName != groupName {
		if store.Groups().Exists(newGroupName) {
			return nil, NewGroupAlreadyExistsError(newGroupName)
		}

		group.GroupName = newGroupName
		group.Arn = store.ARNBuilder().GroupARN(group.Path, newGroupName)

		if err := store.Groups().Put(group); err != nil {
			return nil, err
		}

		users, err := store.UserGroups().ListUsersInGroup(groupName)
		if err != nil {
			if delErr := store.Groups().Delete(newGroupName); delErr != nil {
				logs.Error("UpdateGroup: failed to rollback group creation", logs.Err(delErr))
			}
			return nil, err
		}
		for _, userName := range users {
			if err := store.UserGroups().RemoveUserFromGroup(userName, groupName); err != nil {
				return nil, err
			}
			if err := store.UserGroups().AddUserToGroup(userName, newGroupName); err != nil {
				return nil, err
			}
		}

		if err := store.InlinePolicies().MigratePrincipal(groupName, newGroupName, PrincipalTypeGroup); err != nil {
			return nil, err
		}

		if err := store.AttachedPolicies().MigratePrincipal(groupName, newGroupName, PrincipalTypeGroup); err != nil {
			return nil, err
		}

		if err := store.Groups().Delete(groupName); err != nil {
			return nil, err
		}
	} else {
		if err := store.Groups().Put(group); err != nil {
			return nil, err
		}
	}

	return map[string]interface{}{
		"Group": s.groupToResponse(reqCtx, group),
	}, nil
}

// DeleteGroup deletes an IAM group.
func (s *IAMService) DeleteGroup(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	groupName := request.GetStringParam(req.Parameters, "GroupName")
	if groupName == "" {
		return nil, ErrNoSuchGroup
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if !store.Groups().Exists(groupName) {
		return nil, NewNoSuchGroupError(groupName)
	}

	userCount := store.UserGroups().CountUsersInGroup(groupName)
	if userCount > 0 {
		return nil, NewDeleteGroupConflictError("Cannot delete entity, must remove users from group first.")
	}

	if err := store.InlinePolicies().DeleteAllForPrincipal(PrincipalTypeGroup, groupName); err != nil {
		return nil, err
	}

	attachedPolicies, err := store.AttachedPolicies().ListAttachedPolicies(PrincipalTypeGroup, groupName)
	if err != nil {
		return nil, err
	}
	for _, policyArn := range attachedPolicies {
		if err := store.AttachedPolicies().Detach(PrincipalTypeGroup, groupName, policyArn); err != nil {
			return nil, err
		}
		if err := store.Policies().DecrementAttachmentCount(policyArn); err != nil {
			return nil, err
		}
	}

	if err := store.Groups().Delete(groupName); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// ListGroups lists IAM groups.
func (s *IAMService) ListGroups(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	pathPrefix := request.GetStringParam(req.Parameters, "PathPrefix")
	marker := request.GetStringParam(req.Parameters, "Marker")
	maxItems := getMaxItems(req.Parameters)

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	result, err := store.Groups().List(pathPrefix, marker, maxItems)
	if err != nil {
		return nil, err
	}

	groups := make([]interface{}, len(result.Groups))
	for i, group := range result.Groups {
		groups[i] = s.groupToResponse(reqCtx, group)
	}

	response := map[string]interface{}{
		"Groups":      groups,
		"IsTruncated": result.IsTruncated,
	}

	if result.Marker != "" {
		response["Marker"] = result.Marker
	}

	return response, nil
}

// ListGroupsForUser lists IAM groups for a user.
func (s *IAMService) ListGroupsForUser(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	userName := request.GetStringParam(req.Parameters, "UserName")
	if userName == "" {
		return nil, ErrNoSuchUser
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if !store.Users().Exists(userName) {
		return nil, NewNoSuchUserError(userName)
	}

	groupNames, err := store.UserGroups().ListGroupsForUser(userName)
	if err != nil {
		return nil, err
	}

	groups := make([]interface{}, 0, len(groupNames))
	for _, groupName := range groupNames {
		if group, err := store.Groups().Get(groupName); err == nil {
			groups = append(groups, s.groupToResponse(reqCtx, group))
		}
	}

	return map[string]interface{}{
		"Groups":      groups,
		"IsTruncated": false,
	}, nil
}

func (s *IAMService) groupToResponse(reqCtx *request.RequestContext, group *iamstore.Group) map[string]interface{} {
	resp := map[string]interface{}{
		"GroupId":    group.ID,
		"Path":       group.Path,
		"GroupName":  group.GroupName,
		"Arn":        group.Arn,
		"CreateDate": group.CreateDate.Format(timeutils.ISO8601SimpleFormat),
	}

	if group.PermissionsBoundary != nil {
		resp["PermissionsBoundary"] = map[string]interface{}{
			"PermissionsBoundaryType": group.PermissionsBoundary.PermissionsBoundaryType,
			"PermissionsBoundaryArn":  group.PermissionsBoundary.PermissionsBoundaryArn,
		}
	}

	if tags := tagsToResponse(group.Tags); tags != nil {
		resp["Tags"] = tags
	}

	return resp
}

// TagGroup adds tags to an IAM group.
func (s *IAMService) TagGroup(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	groupName := request.GetStringParam(req.Parameters, "GroupName")
	if groupName == "" {
		return nil, ErrNoSuchGroup
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	group, err := store.Groups().Get(groupName)
	if err != nil {
		return nil, NewNoSuchGroupError(groupName)
	}

	group.Tags = tags.Apply(group.Tags, tags.ParseTagsWithQueryFallback(req.Parameters, "Tags"))

	if err := store.Groups().Put(group); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// UntagGroup removes tags from an IAM group.
func (s *IAMService) UntagGroup(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	groupName := request.GetStringParam(req.Parameters, "GroupName")
	if groupName == "" {
		return nil, ErrNoSuchGroup
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	group, err := store.Groups().Get(groupName)
	if err != nil {
		return nil, NewNoSuchGroupError(groupName)
	}

	group.Tags = removeTags(group.Tags, parseTagKeys(req.Parameters))

	if err := store.Groups().Put(group); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// ListGroupTags lists tags for an IAM group.
func (s *IAMService) ListGroupTags(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	groupName := request.GetStringParam(req.Parameters, "GroupName")
	if groupName == "" {
		return nil, ErrNoSuchGroup
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	group, err := store.Groups().Get(groupName)
	if err != nil {
		return nil, NewNoSuchGroupError(groupName)
	}

	return map[string]interface{}{
		"Tags":        tagsToResponse(group.Tags),
		"IsTruncated": false,
	}, nil
}
