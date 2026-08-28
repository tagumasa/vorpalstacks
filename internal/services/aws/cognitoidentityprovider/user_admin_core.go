package cognitoidentityprovider

import (
	"context"
	"errors"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	"vorpalstacks/internal/core/logs"
	cognitostore "vorpalstacks/internal/store/aws/cognitoidentityprovider"
)

// AdminCreateUserInput carries the wire parameters of AdminCreateUser.
type AdminCreateUserInput struct {
	UserPoolID             string
	Username               string
	MessageAction          string
	TemporaryPassword      string
	ForceAliasCreation     bool
	DesiredDeliveryMediums []string
	UserAttributes         map[string]string
	ValidationData         map[string]string
	ClientMetadata         map[string]string
}

// AdminUpdateUserAttributesInput carries the wire parameters of
// AdminUpdateUserAttributes.
type AdminUpdateUserAttributesInput struct {
	UserPoolID     string
	Username       string
	UserAttributes map[string]string
}

// AdminDeleteUserAttributesInput carries the wire parameters of
// AdminDeleteUserAttributes.
type AdminDeleteUserAttributesInput struct {
	UserPoolID         string
	Username           string
	UserAttributeNames []string
}

// AdminResetUserPasswordInput carries the wire parameters of
// AdminResetUserPassword.
type AdminResetUserPasswordInput struct {
	UserPoolID     string
	Username       string
	ClientMetadata map[string]string
}

// AdminSetUserPasswordInput carries the wire parameters of
// AdminSetUserPassword.
type AdminSetUserPasswordInput struct {
	UserPoolID string
	Username   string
	Password   string
	Permanent  bool
}

// AdminUserGlobalSignOutInput carries the wire parameters of
// AdminUserGlobalSignOut.
type AdminUserGlobalSignOutInput struct {
	UserPoolID string
	Username   string
}

// adminCreateUserCore creates a new user in the specified user pool with
// admin privileges. This operation bypasses the invitation email and sets the
// user status to FORCE_CHANGE_PASSWORD unless MessageAction is set to
// SUPPRESS.
func (s *CognitoService) adminCreateUserCore(ctx context.Context, reqCtx *request.RequestContext, in AdminCreateUserInput) (interface{}, error) {
	if in.UserPoolID == "" || in.Username == "" {
		return nil, ErrInvalidParameter
	}
	if !validateUsernamePattern(in.Username) {
		return nil, ErrInvalidParameter
	}
	if in.MessageAction != "" {
		if !validateMessageAction(in.MessageAction) {
			return nil, ErrInvalidParameter
		}
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	userPool, err := store.GetUserPool(in.UserPoolID)
	if err != nil {
		return nil, ErrResourceNotFound
	}

	userAttrs := in.UserAttributes
	userAttrs["sub"] = ""

	preSignUpResult, err := invokePreSignUp(ctx, s, PreSignUpAdminCreateUser, in.UserPoolID, in.Username, "", userPool.LambdaConfig, userAttrs, in.ValidationData, in.ClientMetadata)
	if err != nil {
		return nil, ErrInternalError
	}

	if preSignUpResult.UserAttributes != nil {
		userAttrs = preSignUpResult.UserAttributes
	}
	delete(userAttrs, "sub")

	user := cognitostore.NewUser(in.UserPoolID, in.Username)
	user.Attributes = userAttrs
	user.UserStatus = "FORCE_CHANGE_PASSWORD"

	tempPassword := in.TemporaryPassword
	if tempPassword == "" && in.MessageAction != "SUPPRESS" {
		// AWS generates a temporary password when none is supplied and the
		// invitation is not suppressed; the CustomMessage trigger below then
		// receives it in the code parameter.
		generated, gerr := generateTemporaryPassword(userPool.PasswordPolicy)
		if gerr != nil {
			return nil, ErrInternalError
		}
		tempPassword = generated
	}
	if tempPassword != "" {
		if err := validatePassword(tempPassword, userPool.PasswordPolicy); err != nil {
			return nil, ErrPasswordPolicyViolation
		}
		if err := setNativePasswordCredentials(user, in.UserPoolID, in.Username, tempPassword); err != nil {
			return nil, ErrInternalError
		}
	}

	if preSignUpResult.AutoConfirmUser {
		user.UserStatus = "CONFIRMED"
		markAutoVerifiedAttributes(user, userPool)
	} else if in.ForceAliasCreation {
		markAutoVerifiedAttributes(user, userPool)
	}

	// DesiredDeliveryMediums controls how the invitation is delivered.
	// Valid values are SMS and EMAIL per the Smithy DeliveryMediumType enum.
	for _, dm := range in.DesiredDeliveryMediums {
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
	if preSignUpResult.AutoConfirmUser || in.MessageAction == "SUPPRESS" {
		if err := invokePostConfirmation(ctx, s, PostConfirmationAdminCreateUser, in.UserPoolID, in.Username, "", userPool.LambdaConfig, attrs); err != nil {
			logs.Warn("PostConfirmation trigger failed", logs.Err(err))
		}
	} else {
		code := "####"
		if tempPassword != "" {
			code = tempPassword
		}
		if _, err := invokeCustomMessage(ctx, s, CustomMessageAdminCreateUser, in.UserPoolID, in.Username, "", userPool.LambdaConfig, code, attrs, in.ClientMetadata); err != nil {
			logs.Warn("CustomMessage trigger failed", logs.Err(err))
		}
	}

	return map[string]interface{}{
		"User": formatUser(user),
	}, nil
}

// adminUpdateUserAttributesCore updates the specified user's attributes in
// the user pool.
func (s *CognitoService) adminUpdateUserAttributesCore(reqCtx *request.RequestContext, in AdminUpdateUserAttributesInput) (interface{}, error) {
	if in.UserPoolID == "" || in.Username == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	user, err := store.GetUser(in.UserPoolID, in.Username)
	if err != nil {
		return nil, ErrUserNotFound
	}

	if user.Attributes == nil {
		user.Attributes = make(map[string]string)
	}

	for k, v := range in.UserAttributes {
		user.Attributes[k] = v
	}

	if err := store.UpdateUser(user); err != nil {
		return nil, ErrInternalError
	}

	return response.EmptyResponse(), nil
}

// adminDeleteUserAttributesCore deletes the specified user attributes from
// the user.
func (s *CognitoService) adminDeleteUserAttributesCore(reqCtx *request.RequestContext, in AdminDeleteUserAttributesInput) (interface{}, error) {
	if in.UserPoolID == "" || in.Username == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	user, err := store.GetUser(in.UserPoolID, in.Username)
	if err != nil {
		return nil, ErrUserNotFound
	}

	if user.Attributes == nil {
		return response.EmptyResponse(), nil
	}

	for _, name := range in.UserAttributeNames {
		delete(user.Attributes, name)
	}

	if err := store.UpdateUser(user); err != nil {
		return nil, ErrInternalError
	}

	return response.EmptyResponse(), nil
}

// adminResetUserPasswordCore forces the specified user to change their
// password on their next sign-in. Sets the user status to RESET_REQUIRED.
func (s *CognitoService) adminResetUserPasswordCore(ctx context.Context, reqCtx *request.RequestContext, in AdminResetUserPasswordInput) (interface{}, error) {
	if in.UserPoolID == "" || in.Username == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	userPool, err := store.GetUserPool(in.UserPoolID)
	if err != nil {
		return nil, ErrResourceNotFound
	}
	user, err := store.GetUser(in.UserPoolID, in.Username)
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
	if _, err := invokeCustomMessage(ctx, s, CustomMessageForgotPassword, in.UserPoolID, in.Username, "", userPool.LambdaConfig, "####", userAttributesMap(user), in.ClientMetadata); err != nil {
		logs.Warn("CustomMessage trigger failed for AdminResetUserPassword", logs.Err(err))
	}

	return response.EmptyResponse(), nil
}

// adminSetUserPasswordCore sets the password for the specified user as an
// administrator. If Permanent is true, the password does not expire.
// Otherwise, the user must change it on next sign-in.
func (s *CognitoService) adminSetUserPasswordCore(reqCtx *request.RequestContext, in AdminSetUserPasswordInput) (interface{}, error) {
	if in.UserPoolID == "" || in.Username == "" || in.Password == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	userPool, err := store.GetUserPool(in.UserPoolID)
	if err != nil {
		return nil, ErrResourceNotFound
	}

	user, err := store.GetUser(in.UserPoolID, in.Username)
	if err != nil {
		return nil, ErrUserNotFound
	}

	if err := validatePassword(in.Password, userPool.PasswordPolicy); err != nil {
		return nil, ErrPasswordPolicyViolation
	}

	if err := setNativePasswordCredentials(user, in.UserPoolID, in.Username, in.Password); err != nil {
		return nil, ErrInternalError
	}

	if in.Permanent {
		user.UserStatus = "CONFIRMED"
	} else {
		user.UserStatus = "FORCE_CHANGE_PASSWORD"
	}

	if err := store.UpdateUser(user); err != nil {
		return nil, ErrInternalError
	}

	return response.EmptyResponse(), nil
}

// adminUserGlobalSignOutCore signs out the specified user from all devices
// by invalidating their refresh tokens.
func (s *CognitoService) adminUserGlobalSignOutCore(reqCtx *request.RequestContext, in AdminUserGlobalSignOutInput) (interface{}, error) {
	if in.UserPoolID == "" || in.Username == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	user, err := store.GetUser(in.UserPoolID, in.Username)
	if err != nil {
		return nil, ErrUserNotFound
	}

	if err := store.DeleteUserTokens(in.UserPoolID, user.ID); err != nil {
		return nil, ErrInternalError
	}

	return response.EmptyResponse(), nil
}
