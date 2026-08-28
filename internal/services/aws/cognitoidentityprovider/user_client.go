package cognitoidentityprovider

import (
	"context"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
)

// DeleteUser deletes the authenticated user from the user pool.
func (s *CognitoService) DeleteUser(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	if err := s.deleteUserByAccessTokenCore(reqCtx, getAccessToken(req)); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// DeleteUserAttributes deletes the specified user attributes for the authenticated user.
func (s *CognitoService) DeleteUserAttributes(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	if err := s.deleteUserAttributesByAccessTokenCore(reqCtx, getAccessToken(req), getStringSliceParam(req, "UserAttributeNames")); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// GetUser returns the user attributes for the authenticated user.
func (s *CognitoService) GetUser(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	user, err := s.getUserByAccessTokenCore(reqCtx, getAccessToken(req))
	if err != nil {
		return nil, err
	}

	result := formatUser(user)
	result["UserAttributes"] = formatUserAttributes(user.Attributes)
	return result, nil
}

// UpdateUserAttributes updates the user attributes for the authenticated user.
func (s *CognitoService) UpdateUserAttributes(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	if err := s.updateUserAttributesByAccessTokenCore(reqCtx, getAccessToken(req), parseUserAttributes(req)); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}
