package cognitoidentityprovider

import (
	"context"

	"vorpalstacks/internal/common/request"
)

// ListUsersInGroup lists the users in a group with store-level pagination.
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_ListUsersInGroup.html
func (s *CognitoService) ListUsersInGroup(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	result, err := s.listUsersInGroupCore(reqCtx.GetRegion(), ListUsersInGroupInput{
		UserPoolID: getUserPoolID(req),
		GroupName:  getGroupName(req),
		MaxResults: request.GetIntParam(req.Parameters, "Limit"),
		NextToken:  request.GetStringParam(req.Parameters, "NextToken"),
	})
	if err != nil {
		return nil, err
	}

	userList := make([]map[string]interface{}, 0, len(result.Users))
	for _, user := range result.Users {
		userList = append(userList, formatUser(user))
	}

	resp := map[string]interface{}{
		"Users": userList,
	}
	if result.NextToken != "" {
		resp["NextToken"] = result.NextToken
	}

	return resp, nil
}

// AdminListGroupsForUser lists the groups for a user with store-level pagination.
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_AdminListGroupsForUser.html
func (s *CognitoService) AdminListGroupsForUser(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	result, err := s.adminListGroupsForUserCore(reqCtx.GetRegion(), AdminListGroupsForUserInput{
		UserPoolID: getUserPoolID(req),
		Username:   getUsername(req),
		MaxResults: request.GetIntParam(req.Parameters, "Limit"),
		NextToken:  request.GetStringParam(req.Parameters, "NextToken"),
	})
	if err != nil {
		return nil, err
	}

	groupList := make([]map[string]interface{}, 0, len(result.Groups))
	for _, group := range result.Groups {
		groupList = append(groupList, formatGroup(group))
	}

	resp := map[string]interface{}{
		"Groups": groupList,
	}
	if result.NextToken != "" {
		resp["NextToken"] = result.NextToken
	}

	return resp, nil
}
