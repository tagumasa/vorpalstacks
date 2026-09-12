package cognitoidentityprovider

import (
	"context"
	"errors"
	"time"

	awserrors "vorpalstacks/internal/common/errors"
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

	// sub is assigned by the service, never supplied by the client, and the
	// schema validation would otherwise hold the client's placeholder
	// against sub's non-empty minimum length.
	userAttrs := in.UserAttributes
	delete(userAttrs, "sub")
	if err := validateUserAttributesAgainstSchema(userPool, userAttrs, true); err != nil {
		return nil, ErrInvalidParameter
	}
	userAttrs["sub"] = ""

	preSignUpResult, err := invokePreSignUp(ctx, s, PreSignUpAdminCreateUser, in.UserPoolID, in.Username, "", userPool.LambdaConfig, userAttrs, in.ValidationData, in.ClientMetadata)
	if err != nil {
		// A Lambda that raises an error is the documented way to reject the
		// creation; a transport failure is an infrastructure error.
		return nil, classifyTriggerFailure(err)
	}

	delete(userAttrs, "sub")
	if preSignUpResult.UserAttributes != nil {
		// The trigger's overrides are screened before merging: the
		// verified-claim keys are stripped — the auto-verify flags below
		// are the only sanctioned way for a trigger to mark a contact
		// attribute verified, and the request's own verified claims (which
		// pair with ForceAliasCreation) stay untouched — and the merged map
		// re-runs the schema gate.
		delete(preSignUpResult.UserAttributes, "email_verified")
		delete(preSignUpResult.UserAttributes, "phone_number_verified")
		for k, v := range preSignUpResult.UserAttributes {
			userAttrs[k] = v
		}
		if err := validateUserAttributesAgainstSchema(userPool, userAttrs, true); err != nil {
			return nil, ErrInvalidLambdaResponse
		}
	}

	user := cognitostore.NewUser(in.UserPoolID, in.Username)
	user.Attributes = userAttrs
	// Every user's record carries its immutable identifier as the required
	// sub attribute, so projections and filters read it from the attribute
	// map like every other attribute.
	user.Attributes["sub"] = user.ID
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
		if err := setNativePasswordCredentials(user, userPool.PasswordPolicy, tempPassword); err != nil {
			if errors.Is(err, ErrPasswordHistoryViolation) {
				return nil, ErrPasswordHistoryViolation
			}
			return nil, ErrInternalError
		}
	}

	if preSignUpResult.AutoConfirmUser {
		user.UserStatus = "CONFIRMED"
		markAutoVerifiedAttributes(user, userPool)
	}

	// The PreSignUp response's auto-verify flags mark the matching
	// attribute verified at creation, independent of the pool-wide
	// AutoVerifiedAttributes marking.
	if preSignUpResult.AutoVerifyEmail && user.Attributes["email"] != "" {
		user.Attributes["email_verified"] = "true"
	}
	if preSignUpResult.AutoVerifyPhone && user.Attributes["phone_number"] != "" {
		user.Attributes["phone_number_verified"] = "true"
	}

	// DesiredDeliveryMediums controls how the invitation is delivered.
	// Valid values are SMS and EMAIL per the Smithy DeliveryMediumType enum.
	for _, dm := range in.DesiredDeliveryMediums {
		if !validDeliveryMediums[dm] {
			return nil, ErrInvalidParameter
		}
	}

	// ForceAliasCreation applies only when the request also activates a
	// contact alias (email_verified or phone_number_verified true) — the
	// model documents the parameter as "used only if the
	// phone_number_verified or email_verified attribute is set to True.
	// Otherwise, it is ignored" — and then migrates a claimed alias value
	// from its previous holder instead of rejecting the conflict.
	var createErr error
	if in.ForceAliasCreation &&
		(user.Attributes["email_verified"] == "true" || user.Attributes["phone_number_verified"] == "true") {
		createErr = store.CreateUserMigrateAliasClaims(user)
	} else {
		createErr = store.CreateUser(user)
	}
	if createErr != nil {
		if errors.Is(createErr, cognitostore.ErrAliasExists) {
			return nil, ErrAliasExists
		}
		if errors.Is(createErr, cognitostore.ErrUserAlreadyExists) {
			return nil, ErrUserAlreadyExists
		}
		return nil, ErrInternalError
	}

	attrs := userAttributesMap(user)
	if preSignUpResult.AutoConfirmUser || in.MessageAction == "SUPPRESS" {
		invokePostConfirmation(ctx, s, PostConfirmationAdminCreateUser, in.UserPoolID, in.Username, "", userPool.LambdaConfig, attrs)
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
func (s *CognitoService) adminUpdateUserAttributesCore(ctx context.Context, reqCtx *request.RequestContext, in AdminUpdateUserAttributesInput) (interface{}, error) {
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

	if err := validateUserAttributesAgainstSchema(userPool, in.UserAttributes, false); err != nil {
		return nil, ErrInvalidParameter
	}

	if user.Attributes == nil {
		user.Attributes = make(map[string]string)
	}

	for k, v := range in.UserAttributes {
		user.Attributes[k] = v
	}

	// An updated email/phone_number enters the verification round: the flag
	// drops, a VerifyUserAttribute code is minted and the CustomMessage
	// UpdateUserAttribute trigger fires.
	if err := issueAttributeUpdateVerification(ctx, s, user, userPool, in.UserAttributes); err != nil {
		return nil, err
	}

	if err := store.UpdateUser(user); err != nil {
		if errors.Is(err, cognitostore.ErrAliasExists) {
			return nil, ErrAliasExists
		}
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

	// The model documents the operation as sending the user "a password-reset
	// code" that their next sign-in must "complete the reset by confirming";
	// the code is stored on the purpose-bound reset field so a follow-up
	// ConfirmForgotPassword can verify it.
	confirmationCode, err := generateConfirmationCode()
	if err != nil {
		return nil, ErrInternalError
	}
	codeExpiry := time.Now().UTC().Add(verificationCodeTTL)
	// The status transition, the code and the credential deactivation run as
	// one serialised write against the fresh record. The reset deactivates
	// every user's credentials: an import hash dies with the reset so it can
	// never verify again, and the reset marker routes every next sign-in to
	// PasswordResetRequiredException — the hash-less CSV-import class
	// included, whose any-password first-sign-in flow the reset supersedes.
	if err := store.UpdateUserFunc(in.UserPoolID, in.Username, func(u *cognitostore.User) error {
		u.UserStatus = "RESET_REQUIRED"
		if u.PasswordHashAlgo != "" {
			u.PasswordHash = ""
			u.PasswordHashAlgo = ""
		}
		u.MigratedAwaitingReset = true
		u.PasswordResetCode = confirmationCode
		u.PasswordResetExpiry = codeExpiry
		return nil
	}); err != nil {
		if errors.Is(err, cognitostore.ErrUserNotFound) {
			return nil, ErrUserNotFound
		}
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

	// The existence read keeps the UserNotFoundException precedence; the
	// authoritative write below runs against the fresh record.
	if _, err := store.GetUser(in.UserPoolID, in.Username); err != nil {
		return nil, ErrUserNotFound
	}

	if err := validatePassword(in.Password, userPool.PasswordPolicy); err != nil {
		return nil, ErrPasswordPolicyViolation
	}

	// The credential rewrite and the status transition run as one serialised
	// read-modify-write, so a concurrent writer's password history cannot be
	// overwritten by a stale snapshot of the record.
	err = store.UpdateUserFunc(in.UserPoolID, in.Username, func(u *cognitostore.User) error {
		if err := setNativePasswordCredentials(u, userPool.PasswordPolicy, in.Password); err != nil {
			return err
		}
		if in.Permanent {
			u.UserStatus = "CONFIRMED"
		} else {
			u.UserStatus = "FORCE_CHANGE_PASSWORD"
		}
		return nil
	})
	if err != nil {
		if errors.Is(err, ErrPasswordHistoryViolation) {
			return nil, ErrPasswordHistoryViolation
		}
		var awsErr *awserrors.AWSError
		if errors.As(err, &awsErr) {
			return nil, err
		}
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
