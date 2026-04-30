package cognitoidentityprovider

import (
	"context"
	"errors"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	"vorpalstacks/internal/core/logs"
	cognitostore "vorpalstacks/internal/store/aws/cognitoidentityprovider"

	"golang.org/x/crypto/bcrypt"
)

// AdminCreateUser creates a new user in the specified user pool with admin privileges.
// This operation bypasses the invitation email and sets the user status to FORCE_CHANGE_PASSWORD
// unless MessageAction is set to SUPPRESS.
func (s *CognitoService) AdminCreateUser(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	userPoolID := getUserPoolID(req)
	username := getUsername(req)
	if userPoolID == "" || username == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	userPool, err := store.GetUserPool(userPoolID)
	if err != nil {
		return nil, ErrResourceNotFound
	}

	userAttrs := parseUserAttributes(req)
	userAttrs["sub"] = ""

	preSignUpResult, err := invokePreSignUp(ctx, s, PreSignUpAdminCreateUser, userPoolID, username, "", userPool.LambdaConfig, userAttrs)
	if err != nil {
		return nil, ErrInternalError
	}

	if preSignUpResult.UserAttributes != nil {
		userAttrs = preSignUpResult.UserAttributes
	}
	delete(userAttrs, "sub")

	user := cognitostore.NewUser(userPoolID, username)
	user.Attributes = userAttrs
	user.UserStatus = "FORCE_CHANGE_PASSWORD"

	tempPassword := req.GetParam("TemporaryPassword")
	if tempPassword != "" {
		if err := validatePassword(tempPassword, userPool.PasswordPolicy); err != nil {
			return nil, ErrPasswordPolicyViolation
		}
		hash, err := bcrypt.GenerateFromPassword([]byte(tempPassword), bcrypt.DefaultCost)
		if err != nil {
			return nil, ErrInternalError
		}
		user.PasswordHash = string(hash)
	}

	if messageAction := req.GetParam("MessageAction"); messageAction == "SUPPRESS" {
		user.UserStatus = "FORCE_CHANGE_PASSWORD"
	}

	if preSignUpResult.AutoConfirmUser {
		user.UserStatus = "CONFIRMED"
	}

	if err := store.CreateUser(user); err != nil {
		if errors.Is(err, cognitostore.ErrUserAlreadyExists) {
			return nil, ErrUserAlreadyExists
		}
		return nil, ErrInternalError
	}

	attrs := userAttributesMap(user)
	if preSignUpResult.AutoConfirmUser || req.GetParam("MessageAction") == "SUPPRESS" {
		if err := invokePostConfirmation(ctx, s, PostConfirmationAdminCreateUser, userPoolID, username, "", userPool.LambdaConfig, attrs); err != nil {
			logs.Warn("PostConfirmation trigger failed", logs.Err(err))
		}
	} else {
		if _, err := invokeCustomMessage(ctx, s, CustomMessageAdminCreateUser, userPoolID, username, "", userPool.LambdaConfig, "####", attrs); err != nil {
			logs.Warn("CustomMessage trigger failed", logs.Err(err))
		}
	}

	return map[string]interface{}{
		"User": formatUser(user),
	}, nil
}

// AdminDeleteUser permanently deletes the specified user from the user pool.
// The user is removed regardless of their current status.
func (s *CognitoService) AdminDeleteUser(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	userPoolID := getUserPoolID(req)
	username := getUsername(req)
	if userPoolID == "" || username == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := store.DeleteUser(userPoolID, username); err != nil {
		return nil, ErrUserNotFound
	}

	return response.EmptyResponse(), nil
}

// AdminGetUser retrieves the specified user from the user pool with all their attributes.
func (s *CognitoService) AdminGetUser(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	userPoolID := getUserPoolID(req)
	username := getUsername(req)
	if userPoolID == "" || username == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	user, err := store.GetUser(userPoolID, username)
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

// AdminUpdateUserAttributes updates the specified user's attributes in the user pool.
func (s *CognitoService) AdminUpdateUserAttributes(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	userPoolID := getUserPoolID(req)
	username := getUsername(req)
	if userPoolID == "" || username == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	user, err := store.GetUser(userPoolID, username)
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

// ListUsers returns a list of users in the specified user pool.
func (s *CognitoService) ListUsers(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	userPoolID := getUserPoolID(req)
	if userPoolID == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	users, err := store.ListUsers(userPoolID)
	if err != nil {
		return nil, ErrInternalError
	}

	filter := req.GetParam("Filter")
	userList := make([]map[string]interface{}, 0)
	for _, user := range users {
		if filter != "" && !matchUserFilter(user, filter) {
			continue
		}
		userList = append(userList, formatUser(user))
	}

	return map[string]interface{}{
		"Users": userList,
	}, nil
}

// AdminEnableUser enables the specified user in the user pool.
func (s *CognitoService) AdminEnableUser(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	userPoolID := getUserPoolID(req)
	username := getUsername(req)
	if userPoolID == "" || username == "" {
		return nil, ErrInvalidParameter
	}

	if err := s.setUserEnabled(reqCtx, userPoolID, username, true); err != nil {
		if err == cognitostore.ErrUserNotFound {
			return nil, ErrUserNotFound
		}
		return nil, ErrInternalError
	}

	return response.EmptyResponse(), nil
}

// AdminDisableUser disables the specified user in the user pool.
func (s *CognitoService) AdminDisableUser(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	userPoolID := getUserPoolID(req)
	username := getUsername(req)
	if userPoolID == "" || username == "" {
		return nil, ErrInvalidParameter
	}

	if err := s.setUserEnabled(reqCtx, userPoolID, username, false); err != nil {
		if err == cognitostore.ErrUserNotFound {
			return nil, ErrUserNotFound
		}
		return nil, ErrInternalError
	}

	return response.EmptyResponse(), nil
}

// AdminDeleteUserAttributes deletes the specified user attributes from the user.
func (s *CognitoService) AdminDeleteUserAttributes(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	userPoolID := getUserPoolID(req)
	username := getUsername(req)
	if userPoolID == "" || username == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	user, err := store.GetUser(userPoolID, username)
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

// AdminResetUserPassword forces the specified user to change their password on their next sign-in.
// Sets the user status to RESET_REQUIRED.
func (s *CognitoService) AdminResetUserPassword(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	userPoolID := getUserPoolID(req)
	username := getUsername(req)
	if userPoolID == "" || username == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	user, err := store.GetUser(userPoolID, username)
	if err != nil {
		return nil, ErrUserNotFound
	}

	user.UserStatus = "RESET_REQUIRED"
	if err := store.UpdateUser(user); err != nil {
		return nil, ErrInternalError
	}

	return response.EmptyResponse(), nil
}

// AdminSetUserPassword sets the password for the specified user as an administrator.
// If Permanent is true, the password does not expire. Otherwise, the user must change it on next sign-in.
func (s *CognitoService) AdminSetUserPassword(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	userPoolID := getUserPoolID(req)
	username := getUsername(req)
	password := getPassword(req)
	permanent := getBoolParam(req, "Permanent")

	if userPoolID == "" || username == "" || password == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	userPool, err := store.GetUserPool(userPoolID)
	if err != nil {
		return nil, ErrResourceNotFound
	}

	user, err := store.GetUser(userPoolID, username)
	if err != nil {
		return nil, ErrUserNotFound
	}

	if err := validatePassword(password, userPool.PasswordPolicy); err != nil {
		return nil, ErrPasswordPolicyViolation
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, ErrInternalError
	}
	user.PasswordHash = string(hash)

	if permanent {
		user.UserStatus = "CONFIRMED"
	} else {
		user.UserStatus = "FORCE_CHANGE_PASSWORD"
	}

	if err := store.UpdateUser(user); err != nil {
		return nil, ErrInternalError
	}

	return response.EmptyResponse(), nil
}

// AdminUserGlobalSignOut signs out the specified user from all devices by invalidating their refresh tokens.
func (s *CognitoService) AdminUserGlobalSignOut(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	userPoolID := getUserPoolID(req)
	username := getUsername(req)
	if userPoolID == "" || username == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	user, err := store.GetUser(userPoolID, username)
	if err != nil {
		return nil, ErrUserNotFound
	}

	if err := store.DeleteUserTokens(userPoolID, user.ID); err != nil {
		return nil, ErrInternalError
	}

	return response.EmptyResponse(), nil
}
