package cognitoidentityprovider

import (
	"context"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
)

// DeleteUser deletes the authenticated user from the user pool.
func (s *CognitoService) DeleteUser(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	accessToken := getAccessToken(req)
	if accessToken == "" {
		return nil, ErrInvalidParameter
	}

	userID, err := s.ValidateAccessToken(reqCtx, accessToken)
	if err != nil {
		return nil, ErrNotAuthorized
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	user, err := store.GetUserByID(userID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	if err := store.DeleteUser(user.UserPoolID, user.Username); err != nil {
		return nil, ErrInternalError
	}
	if err := store.DeleteUserTokens(user.UserPoolID, user.ID); err != nil {
		return nil, ErrInternalError
	}

	return response.EmptyResponse(), nil
}

// DeleteUserAttributes deletes the specified user attributes for the authenticated user.
func (s *CognitoService) DeleteUserAttributes(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	accessToken := getAccessToken(req)
	if accessToken == "" {
		return nil, ErrInvalidParameter
	}

	userID, err := s.ValidateAccessToken(reqCtx, accessToken)
	if err != nil {
		return nil, ErrNotAuthorized
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	user, err := store.GetUserByID(userID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	if user.Attributes == nil {
		return response.EmptyResponse(), nil
	}

	attrNames := getStringSliceParam(req, "UserAttributeNames")
	for _, name := range attrNames {
		delete(user.Attributes, name)
	}

	if err := store.UpdateUser(user); err != nil {
		return nil, ErrInternalError
	}

	return response.EmptyResponse(), nil
}

// GetUser returns the user attributes for the authenticated user.
func (s *CognitoService) GetUser(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	accessToken := getAccessToken(req)
	if accessToken == "" {
		return nil, ErrInvalidParameter
	}

	userID, err := s.ValidateAccessToken(reqCtx, accessToken)
	if err != nil {
		return nil, ErrNotAuthorized
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	user, err := store.GetUserByID(userID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	return map[string]interface{}{
		"Username":             user.Username,
		"UserAttributes":       formatUserAttributes(user.Attributes),
		"UserCreateDate":       user.CreatedDate.Unix(),
		"UserLastModifiedDate": user.LastModifiedDate.Unix(),
		"Enabled":              user.Enabled,
		"UserStatus":           user.UserStatus,
	}, nil
}

// UpdateUserAttributes updates the user attributes for the authenticated user.
func (s *CognitoService) UpdateUserAttributes(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	accessToken := getAccessToken(req)
	if accessToken == "" {
		return nil, ErrInvalidParameter
	}

	userID, err := s.ValidateAccessToken(reqCtx, accessToken)
	if err != nil {
		return nil, ErrNotAuthorized
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	user, err := store.GetUserByID(userID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	if user.Attributes == nil {
		user.Attributes = make(map[string]string)
	}

	newAttrs := parseUserAttributes(req)
	for k, v := range newAttrs {
		user.Attributes[k] = v
	}

	if err := store.UpdateUser(user); err != nil {
		return nil, ErrInternalError
	}

	return response.EmptyResponse(), nil
}
