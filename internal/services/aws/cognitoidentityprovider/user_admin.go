package cognitoidentityprovider

import (
	"context"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
)

// AdminCreateUser creates a new user in the specified user pool with admin privileges.
// This operation bypasses the invitation email and sets the user status to FORCE_CHANGE_PASSWORD
// unless MessageAction is set to SUPPRESS.
func (s *CognitoService) AdminCreateUser(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return s.adminCreateUserCore(ctx, reqCtx, AdminCreateUserInput{
		UserPoolID:             getUserPoolID(req),
		Username:               getUsername(req),
		MessageAction:          req.GetParam("MessageAction"),
		TemporaryPassword:      req.GetParam("TemporaryPassword"),
		ForceAliasCreation:     getBoolParam(req, "ForceAliasCreation"),
		DesiredDeliveryMediums: getStringSliceParam(req, "DesiredDeliveryMediums"),
		UserAttributes:         parseUserAttributes(req),
		ValidationData:         parseValidationData(req),
		ClientMetadata:         parseClientMetadata(req),
	})
}

// AdminDeleteUser permanently deletes the specified user from the user pool.
// The user is removed regardless of their current status.
func (s *CognitoService) AdminDeleteUser(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	if err := s.adminDeleteUserCore(reqCtx.GetRegion(), getUserPoolID(req), getUsername(req)); err != nil {
		return nil, err
	}
	return response.EmptyResponse(), nil
}

// AdminGetUser retrieves the specified user from the user pool with all their attributes.
func (s *CognitoService) AdminGetUser(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	user, err := s.adminGetUserCore(reqCtx.GetRegion(), getUserPoolID(req), getUsername(req))
	if err != nil {
		return nil, err
	}

	result := formatUser(user)
	result["UserAttributes"] = formatUserAttributes(user.Attributes)
	return result, nil
}

// AdminUpdateUserAttributes updates the specified user's attributes in the user pool.
func (s *CognitoService) AdminUpdateUserAttributes(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return s.adminUpdateUserAttributesCore(reqCtx, AdminUpdateUserAttributesInput{
		UserPoolID:     getUserPoolID(req),
		Username:       getUsername(req),
		UserAttributes: parseUserAttributes(req),
	})
}

// ListUsers returns a list of users in the specified user pool.
func (s *CognitoService) ListUsers(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	// Smithy QueryLimitType: range {min: 0, max: 60}
	limit, err := parseListLimit(req.Parameters, "Limit", 60)
	if err != nil {
		return nil, err
	}
	result, err := s.listUsersCore(reqCtx.GetRegion(), ListUsersInput{
		UserPoolID: getUserPoolID(req),
		MaxResults: limit,
		NextToken:  request.GetStringParam(req.Parameters, "PaginationToken"),
		Filter:     req.GetParam("Filter"),
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
		resp["PaginationToken"] = result.NextToken
	}

	return resp, nil
}

// AdminEnableUser enables the specified user in the user pool.
func (s *CognitoService) AdminEnableUser(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	if err := s.adminEnableUserCore(reqCtx.GetRegion(), getUserPoolID(req), getUsername(req)); err != nil {
		return nil, err
	}
	return response.EmptyResponse(), nil
}

// AdminDisableUser disables the specified user in the user pool.
func (s *CognitoService) AdminDisableUser(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	if err := s.adminDisableUserCore(reqCtx.GetRegion(), getUserPoolID(req), getUsername(req)); err != nil {
		return nil, err
	}
	return response.EmptyResponse(), nil
}

// AdminDeleteUserAttributes deletes the specified user attributes from the user.
func (s *CognitoService) AdminDeleteUserAttributes(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return s.adminDeleteUserAttributesCore(reqCtx, AdminDeleteUserAttributesInput{
		UserPoolID:         getUserPoolID(req),
		Username:           getUsername(req),
		UserAttributeNames: getStringSliceParam(req, "UserAttributeNames"),
	})
}

// AdminResetUserPassword forces the specified user to change their password on their next sign-in.
// Sets the user status to RESET_REQUIRED.
func (s *CognitoService) AdminResetUserPassword(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return s.adminResetUserPasswordCore(ctx, reqCtx, AdminResetUserPasswordInput{
		UserPoolID:     getUserPoolID(req),
		Username:       getUsername(req),
		ClientMetadata: parseClientMetadata(req),
	})
}

// AdminSetUserPassword sets the password for the specified user as an administrator.
// If Permanent is true, the password does not expire. Otherwise, the user must change it on next sign-in.
func (s *CognitoService) AdminSetUserPassword(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return s.adminSetUserPasswordCore(reqCtx, AdminSetUserPasswordInput{
		UserPoolID: getUserPoolID(req),
		Username:   getUsername(req),
		Password:   getPassword(req),
		Permanent:  getBoolParam(req, "Permanent"),
	})
}

// AdminUserGlobalSignOut signs out the specified user from all devices by invalidating their refresh tokens.
func (s *CognitoService) AdminUserGlobalSignOut(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return s.adminUserGlobalSignOutCore(reqCtx, AdminUserGlobalSignOutInput{
		UserPoolID: getUserPoolID(req),
		Username:   getUsername(req),
	})
}
