package cognitoidentityprovider

import (
	"context"
	"vorpalstacks/internal/common/request"
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

	maxResults := request.GetIntParam(req.Parameters, "Limit")
	if maxResults <= 0 || maxResults > 60 {
		maxResults = 60
	}
	nextToken := request.GetStringParam(req.Parameters, "NextToken")

	started := nextToken == ""
	userList := make([]map[string]interface{}, 0, maxResults)
	for _, user := range users {
		if !started {
			if user.Username == nextToken {
				started = true
			}
			continue
		}
		userList = append(userList, formatUser(user))
		if len(userList) >= maxResults {
			break
		}
	}

	resp := map[string]interface{}{
		"Users": userList,
	}
	if len(userList) >= maxResults && len(userList) > 0 {
		resp["NextToken"] = userList[len(userList)-1]["Username"]
	}

	return resp, nil
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

	maxResults := request.GetIntParam(req.Parameters, "Limit")
	if maxResults <= 0 || maxResults > 60 {
		maxResults = 60
	}
	nextToken := request.GetStringParam(req.Parameters, "NextToken")

	started := nextToken == ""
	groupList := make([]map[string]interface{}, 0, maxResults)
	for _, group := range groups {
		if !started {
			if group.Name == nextToken {
				started = true
			}
			continue
		}
		groupList = append(groupList, formatGroup(group))
		if len(groupList) >= maxResults {
			break
		}
	}

	resp := map[string]interface{}{
		"Groups": groupList,
	}
	if len(groupList) >= maxResults && len(groupList) > 0 {
		resp["NextToken"] = groupList[len(groupList)-1]["GroupName"]
	}

	return resp, nil
}
