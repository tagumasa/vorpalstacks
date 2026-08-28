// Package iam provides IAM service operations for vorpalstacks.
package iam

import (
	"context"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
)

// AddUserToGroup adds a user to a group.
func (s *IAMService) AddUserToGroup(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	input := &UserGroupMembershipInput{
		UserName:  request.GetStringParam(req.Parameters, "UserName"),
		GroupName: request.GetStringParam(req.Parameters, "GroupName"),
	}
	if err := s.addUserToGroupCore(store, input); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// RemoveUserFromGroup removes a user from a group.
func (s *IAMService) RemoveUserFromGroup(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	input := &UserGroupMembershipInput{
		UserName:  request.GetStringParam(req.Parameters, "UserName"),
		GroupName: request.GetStringParam(req.Parameters, "GroupName"),
	}
	if err := s.removeUserFromGroupCore(store, input); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}
