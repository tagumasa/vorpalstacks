package cognitoidentityprovider

import (
	"context"
	"vorpalstacks/internal/services/aws/common/request"
)

// ListUsersInGroup lists the users in a group.
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_ListUsersInGroup.html
func (s *CognitoService) ListUsersInGroup(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	userPoolID := getUserPoolID(req)
	groupName := getGroupName(req)
	if userPoolID == "" || groupName == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	users, err := store.ListUsersInGroup(userPoolID, groupName)
	if err != nil {
		return nil, ErrGroupNotFound
	}

	userList := make([]map[string]interface{}, 0)
	for _, user := range users {
		userList = append(userList, formatUser(user))
	}

	return map[string]interface{}{
		"Users": userList,
	}, nil
}

// AdminListGroupsForUser lists the groups for a user.
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_AdminListGroupsForUser.html
func (s *CognitoService) AdminListGroupsForUser(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	userPoolID := getUserPoolID(req)
	username := getUsername(req)
	if userPoolID == "" || username == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	groups, err := store.ListGroupsForUser(userPoolID, username)
	if err != nil {
		return nil, ErrUserNotFound
	}

	groupList := make([]map[string]interface{}, 0)
	for _, group := range groups {
		groupList = append(groupList, formatGroup(group))
	}

	return map[string]interface{}{
		"Groups": groupList,
	}, nil
}
