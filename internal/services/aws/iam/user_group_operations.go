// Package iam provides IAM service operations for vorpalstacks.
package iam

import (
	"context"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
)

// AddUserToGroup adds a user to a group.
func (s *IAMService) AddUserToGroup(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	userName := request.GetStringParam(req.Parameters, "UserName")
	groupName := request.GetStringParam(req.Parameters, "GroupName")

	if userName == "" {
		return nil, ErrNoSuchUser
	}
	if groupName == "" {
		return nil, ErrNoSuchGroup
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if !store.Users().Exists(userName) {
		return nil, NewNoSuchUserError(userName)
	}

	if !store.Groups().Exists(groupName) {
		return nil, NewNoSuchGroupError(groupName)
	}

	if store.UserGroups().IsUserInGroup(userName, groupName) {
		return nil, NewUserAlreadyInGroupError(userName, groupName)
	}

	if err := store.UserGroups().AddUserToGroup(userName, groupName); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// RemoveUserFromGroup removes a user from a group.
func (s *IAMService) RemoveUserFromGroup(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	userName := request.GetStringParam(req.Parameters, "UserName")
	groupName := request.GetStringParam(req.Parameters, "GroupName")

	if userName == "" {
		return nil, ErrNoSuchUser
	}
	if groupName == "" {
		return nil, ErrNoSuchGroup
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if !store.Users().Exists(userName) {
		return nil, NewNoSuchUserError(userName)
	}

	if !store.Groups().Exists(groupName) {
		return nil, NewNoSuchGroupError(groupName)
	}

	if !store.UserGroups().IsUserInGroup(userName, groupName) {
		return nil, NewUserNotInGroupError(userName, groupName)
	}

	if err := store.UserGroups().RemoveUserFromGroup(userName, groupName); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}
