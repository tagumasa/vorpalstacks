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
	if !validateUsernamePattern(username) {
		return nil, ErrInvalidParameter
	}
	if ma := req.GetParam("MessageAction"); ma != "" {
		if !validateMessageAction(ma) {
			return nil, ErrInvalidParameter
		}
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

	validationData := parseValidationData(req)
	forceAliasCreation := getBoolParam(req, "ForceAliasCreation")

	preSignUpResult, err := invokePreSignUp(ctx, s, PreSignUpAdminCreateUser, userPoolID, username, "", userPool.LambdaConfig, userAttrs, validationData, parseClientMetadata(req))
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
		saltHex, verifierHex, verr := computeSrpVerifier(userPoolID, username, tempPassword)
		if verr != nil {
			return nil, ErrInternalError
		}
		user.SrpSalt = saltHex
		user.SrpVerifier = verifierHex
	}

	if preSignUpResult.AutoConfirmUser {
		user.UserStatus = "CONFIRMED"
		markAutoVerifiedAttributes(user, userPool)
	} else if forceAliasCreation {
		markAutoVerifiedAttributes(user, userPool)
	}

	// DesiredDeliveryMediums controls how the invitation is delivered.
	// Valid values are SMS and EMAIL per the Smithy DeliveryMediumType enum.
	for _, dm := range getStringSliceParam(req, "DesiredDeliveryMediums") {
		if !validDeliveryMediums[dm] {
			return nil, ErrInvalidParameter
		}
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
		if _, err := invokeCustomMessage(ctx, s, CustomMessageAdminCreateUser, userPoolID, username, "", userPool.LambdaConfig, "####", attrs, parseClientMetadata(req)); err != nil {
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
	userPool, err := store.GetUserPool(userPoolID)
	if err != nil {
		return nil, ErrResourceNotFound
	}
	user, err := store.GetUser(userPoolID, username)
	if err != nil {
		return nil, ErrUserNotFound
	}

	user.UserStatus = "RESET_REQUIRED"
	if err := store.UpdateUser(user); err != nil {
		return nil, ErrInternalError
	}

	// AdminResetUserPassword always sends a reset code via the CustomMessage
	// trigger. The AWS API does not accept a MessageAction parameter for this
	// operation, so suppression is not an option.
	clientMetadata := parseClientMetadata(req)
	if _, err := invokeCustomMessage(ctx, s, CustomMessageForgotPassword, userPoolID, username, "", userPool.LambdaConfig, "####", userAttributesMap(user), clientMetadata); err != nil {
		logs.Warn("CustomMessage trigger failed for AdminResetUserPassword", logs.Err(err))
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
	saltHex, verifierHex, verr := computeSrpVerifier(userPoolID, username, password)
	if verr != nil {
		return nil, ErrInternalError
	}
	user.SrpSalt = saltHex
	user.SrpVerifier = verifierHex

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
