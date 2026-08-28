package cognitoidentityprovider

import (
	"context"
	"vorpalstacks/internal/common/iam"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
)

// CreateGroup creates a group in a Cognito user pool.
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_CreateGroup.html
func (s *CognitoService) CreateGroup(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	in := CreateGroupInput{
		UserPoolID:  getUserPoolID(req),
		GroupName:   getGroupName(req),
		Description: req.GetParam("Description"),
		RoleArn:     req.GetParam("RoleArn"),
	}
	if precedence, ok := getIntParamOK(req, "Precedence"); ok {
		in.Precedence = &precedence
	}

	group, err := s.createGroupValidatedCore(ctx, reqCtx.GetRegion(), in, reqCtx.GetIAMValidator())
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"Group": formatGroup(group),
	}, nil
}

// GetGroup returns information about a group in a Cognito user pool.
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_GetGroup.html
func (s *CognitoService) GetGroup(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	group, err := s.getGroupCore(reqCtx.GetRegion(), getUserPoolID(req), getGroupName(req))
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{"Group": formatGroup(group)}, nil
}

// DeleteGroup deletes a group from a Cognito user pool.
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_DeleteGroup.html
func (s *CognitoService) DeleteGroup(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	if err := s.deleteGroupCore(reqCtx.GetRegion(), getUserPoolID(req), getGroupName(req)); err != nil {
		return nil, err
	}
	return response.EmptyResponse(), nil
}

// ListGroups lists the groups in a Cognito user pool.
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_ListGroups.html
func (s *CognitoService) ListGroups(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	// Smithy QueryLimitType: range {min: 0, max: 60}
	limit, err := parseListLimit(req.Parameters, "Limit", 60)
	if err != nil {
		return nil, err
	}
	result, err := s.listGroupsCore(reqCtx.GetRegion(), ListGroupsInput{
		UserPoolID: getUserPoolID(req),
		MaxResults: limit,
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

// UpdateGroup updates a group in a Cognito user pool.
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_UpdateGroup.html
func (s *CognitoService) UpdateGroup(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	userPoolID := getUserPoolID(req)
	groupName := getGroupName(req)
	if userPoolID == "" || groupName == "" {
		return nil, ErrInvalidParameter
	}

	group, err := s.getGroupCore(reqCtx.GetRegion(), userPoolID, groupName)
	if err != nil {
		return nil, err
	}

	if description := req.GetParam("Description"); description != "" {
		group.Description = description
	}
	if roleArn := req.GetParam("RoleArn"); roleArn != "" {
		validator := reqCtx.GetIAMValidator()
		if err := validator.ValidateRoleForServiceWithErrors(ctx, roleArn, iam.ServicePrincipalCognito, &iam.RoleErrorFactories{
			RoleNotFoundError:        iam.NewCognitoRoleError,
			RoleCannotBeAssumedError: iam.NewCognitoRoleError,
			InvalidArnError:          iam.NewCognitoRoleError,
		}); err != nil {
			return nil, err
		}
		group.RoleArn = roleArn
	}
	if precedence, ok := getIntParamOK(req, "Precedence"); ok {
		if !validatePrecedence(precedence) {
			return nil, ErrInvalidParameter
		}
		group.Precedence = &precedence
	}

	if err := s.updateGroupCore(reqCtx.GetRegion(), group); err != nil {
		return nil, err
	}

	return map[string]interface{}{"Group": formatGroup(group)}, nil
}

// AdminAddUserToGroup adds a user to a group in a Cognito user pool.
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_AdminAddUserToGroup.html
func (s *CognitoService) AdminAddUserToGroup(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	if err := s.adminAddUserToGroupCore(reqCtx.GetRegion(), getUserPoolID(req), getGroupName(req), getUsername(req)); err != nil {
		return nil, err
	}
	return response.EmptyResponse(), nil
}

// AdminRemoveUserFromGroup removes a user from a group in a Cognito user pool.
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_AdminRemoveUserFromGroup.html
func (s *CognitoService) AdminRemoveUserFromGroup(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	if err := s.adminRemoveUserFromGroupCore(reqCtx.GetRegion(), getUserPoolID(req), getGroupName(req), getUsername(req)); err != nil {
		return nil, err
	}
	return response.EmptyResponse(), nil
}
