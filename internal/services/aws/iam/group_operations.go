// Package iam provides IAM service operations for vorpalstacks.
package iam

import (
	"context"
	"vorpalstacks/internal/common/pagination"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	iamstore "vorpalstacks/internal/store/aws/iam"
	"vorpalstacks/internal/utils/timeutils"
)

// CreateGroup creates a new IAM group.
func (s *IAMService) CreateGroup(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	input := &CreateGroupInput{
		GroupName: request.GetStringParam(req.Parameters, "GroupName"),
		Path:      request.GetStringParam(req.Parameters, "Path"),
	}
	group, err := s.createGroupCore(store, input)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"Group": s.groupToResponse(reqCtx, group),
	}, nil
}

// GetGroup retrieves information about an IAM group.
func (s *IAMService) GetGroup(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	groupName := request.GetStringParam(req.Parameters, "GroupName")

	marker := request.GetStringParam(req.Parameters, "Marker")
	maxItems := pagination.GetMaxItems(req.Parameters, pagination.DefaultMaxItems)

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	group, err := s.getGroupCore(store, groupName)
	if err != nil {
		return nil, err
	}
	users, err := s.listUsersInGroupCore(store, groupName)
	if err != nil {
		return nil, err
	}

	userList := make([]interface{}, 0, len(users))
	for _, user := range users {
		userList = append(userList, s.userToResponse(reqCtx, user))
	}

	paged := pagination.PaginateSlice(userList, marker, maxItems, func(item interface{}) string {
		if m, ok := item.(map[string]interface{}); ok {
			if name, ok := m["UserName"].(string); ok {
				return name
			}
		}
		return ""
	})

	resp := map[string]interface{}{
		"Group":       s.groupToResponse(reqCtx, group),
		"Users":       paged.Items,
		"IsTruncated": paged.IsTruncated,
	}
	if paged.NextMarker != "" {
		resp["Marker"] = paged.NextMarker
	}
	return resp, nil
}

// UpdateGroup updates an IAM group.
func (s *IAMService) UpdateGroup(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	input := &UpdateGroupInput{
		GroupName:    request.GetStringParam(req.Parameters, "GroupName"),
		NewPath:      request.GetStringParam(req.Parameters, "NewPath"),
		NewGroupName: request.GetStringParam(req.Parameters, "NewGroupName"),
	}
	group, err := s.updateGroupCore(store, input)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"Group": s.groupToResponse(reqCtx, group),
	}, nil
}

// DeleteGroup deletes an IAM group.
func (s *IAMService) DeleteGroup(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	input := &DeleteGroupInput{
		GroupName: request.GetStringParam(req.Parameters, "GroupName"),
	}
	if err := s.deleteGroupCore(store, input); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// ListGroups lists IAM groups.
func (s *IAMService) ListGroups(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	pathPrefix := request.GetStringParam(req.Parameters, "PathPrefix")
	marker := request.GetStringParam(req.Parameters, "Marker")
	maxItems := pagination.GetMaxItems(req.Parameters, pagination.DefaultMaxItems)

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	result, err := s.listGroupsCore(store, pathPrefix, marker, maxItems)
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
		return nil, NewValidationError("UserName")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	groupList, err := s.listGroupsForUserCore(store, userName)
	if err != nil {
		return nil, err
	}

	groups := make([]interface{}, 0, len(groupList))
	for _, group := range groupList {
		groups = append(groups, s.groupToResponse(reqCtx, group))
	}

	marker := request.GetStringParam(req.Parameters, "Marker")
	maxItems := pagination.GetMaxItems(req.Parameters, pagination.DefaultMaxItems)

	paged := pagination.PaginateSlice(groups, marker, maxItems, func(item interface{}) string {
		if m, ok := item.(map[string]interface{}); ok {
			if name, ok := m["GroupName"].(string); ok {
				return name
			}
		}
		return ""
	})

	resp := map[string]interface{}{
		"Groups":      paged.Items,
		"IsTruncated": paged.IsTruncated,
	}
	if paged.NextMarker != "" {
		resp["Marker"] = paged.NextMarker
	}
	return resp, nil
}

func (s *IAMService) groupToResponse(reqCtx *request.RequestContext, group *iamstore.Group) map[string]interface{} {
	return map[string]interface{}{
		"GroupId":    group.ID,
		"Path":       group.Path,
		"GroupName":  group.GroupName,
		"Arn":        group.Arn,
		"CreateDate": group.CreateDate.Format(timeutils.ISO8601SimpleFormat),
	}
}
