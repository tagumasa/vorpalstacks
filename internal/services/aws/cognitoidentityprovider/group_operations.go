package cognitoidentityprovider

import (
	"context"
	"vorpalstacks/internal/services/aws/common/iam"
	"vorpalstacks/internal/services/aws/common/request"
	"vorpalstacks/internal/services/aws/common/response"
	cognitostore "vorpalstacks/internal/store/aws/cognitoidentityprovider"
)

// CreateGroup creates a group in a Cognito user pool.
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_CreateGroup.html
func (s *CognitoService) CreateGroup(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	userPoolID := getUserPoolID(req)
	groupName := getGroupName(req)
	if userPoolID == "" || groupName == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	_, err = store.GetUserPool(userPoolID)
	if err != nil {
		return nil, ErrResourceNotFound
	}

	group := cognitostore.NewGroup(userPoolID, groupName)
	group.Description = getParam(req, "Description")
	group.RoleArn = getParam(req, "RoleArn")
	if precedence, ok := getIntParamOK(req, "Precedence"); ok {
		group.Precedence = precedence
	}

	if group.RoleArn != "" {
		validator := reqCtx.GetIAMValidator()
		if err := validator.ValidateRoleForServiceWithErrors(ctx, group.RoleArn, iam.ServicePrincipalCognito, &iam.RoleErrorFactories{
			RoleNotFoundError:        iam.NewCognitoRoleError,
			RoleCannotBeAssumedError: iam.NewCognitoRoleError,
			InvalidArnError:          iam.NewCognitoRoleError,
		}); err != nil {
			return nil, err
		}
	}

	if err := store.CreateGroup(group); err != nil {
		return nil, ErrGroupAlreadyExists
	}

	return map[string]interface{}{
		"Group": formatGroup(group),
	}, nil
}

// GetGroup returns information about a group in a Cognito user pool.
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_GetGroup.html
func (s *CognitoService) GetGroup(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	userPoolID := getUserPoolID(req)
	groupName := getGroupName(req)
	if userPoolID == "" || groupName == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	group, err := store.GetGroup(userPoolID, groupName)
	if err != nil {
		return nil, ErrGroupNotFound
	}

	return map[string]interface{}{
		"Group": formatGroup(group),
	}, nil
}

// DeleteGroup deletes a group from a Cognito user pool.
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_DeleteGroup.html
func (s *CognitoService) DeleteGroup(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	userPoolID := getUserPoolID(req)
	groupName := getGroupName(req)
	if userPoolID == "" || groupName == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if err := store.DeleteGroup(userPoolID, groupName); err != nil {
		return nil, ErrGroupNotFound
	}

	return response.EmptyResponse(), nil
}

// ListGroups lists the groups in a Cognito user pool.
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_ListGroups.html
func (s *CognitoService) ListGroups(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	userPoolID := getUserPoolID(req)
	if userPoolID == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	groups, err := store.ListGroups(userPoolID)
	if err != nil {
		return nil, ErrInternalError
	}

	groupList := make([]map[string]interface{}, 0)
	for _, group := range groups {
		groupList = append(groupList, formatGroup(group))
	}

	return map[string]interface{}{
		"Groups": groupList,
	}, nil
}

// UpdateGroup updates a group in a Cognito user pool.
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_UpdateGroup.html
func (s *CognitoService) UpdateGroup(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	userPoolID := getUserPoolID(req)
	groupName := getGroupName(req)
	if userPoolID == "" || groupName == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	group, err := store.GetGroup(userPoolID, groupName)
	if err != nil {
		return nil, ErrGroupNotFound
	}

	if description := getParam(req, "Description"); description != "" {
		group.Description = description
	}
	if roleArn := getParam(req, "RoleArn"); roleArn != "" {
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
		group.Precedence = precedence
	}

	if err := store.UpdateGroup(group); err != nil {
		return nil, ErrInternalError
	}

	return map[string]interface{}{
		"Group": formatGroup(group),
	}, nil
}

// AdminAddUserToGroup adds a user to a group in a Cognito user pool.
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_AdminAddUserToGroup.html
func (s *CognitoService) AdminAddUserToGroup(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	userPoolID := getUserPoolID(req)
	groupName := getGroupName(req)
	username := getUsername(req)

	if userPoolID == "" || groupName == "" || username == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if err := store.AddUserToGroup(userPoolID, groupName, username); err != nil {
		if err == cognitostore.ErrGroupNotFound {
			return nil, ErrGroupNotFound
		}
		if err == cognitostore.ErrUserNotFound {
			return nil, ErrUserNotFound
		}
		return nil, ErrInternalError
	}

	return response.EmptyResponse(), nil
}

// AdminRemoveUserFromGroup removes a user from a group in a Cognito user pool.
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_AdminRemoveUserFromGroup.html
func (s *CognitoService) AdminRemoveUserFromGroup(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	userPoolID := getUserPoolID(req)
	groupName := getGroupName(req)
	username := getUsername(req)

	if userPoolID == "" || groupName == "" || username == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if err := store.RemoveUserFromGroup(userPoolID, groupName, username); err != nil {
		if err == cognitostore.ErrGroupNotFound {
			return nil, ErrGroupNotFound
		}
		if err == cognitostore.ErrUserNotFound {
			return nil, ErrUserNotFound
		}
		return nil, ErrInternalError
	}

	return response.EmptyResponse(), nil
}
