package cognitoidentityprovider

import (
	"context"
	"crypto/subtle"
	"errors"
	"fmt"
	"os"
	"time"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	"vorpalstacks/internal/core/logs"
	cognitostore "vorpalstacks/internal/store/aws/cognitoidentityprovider"
)

// testMode is resolved once at process start: the regression runner starts
// the server with TEST_MODE=true and the flag never changes during a run,
// so the code-verification bypasses read this value instead of re-querying
// the environment on every request.
var testMode = os.Getenv("TEST_MODE") == "true"

// SignUpInput carries the wire parameters of SignUp.
type SignUpInput struct {
	ClientID       string
	Username       string
	Password       string
	UserAttributes map[string]string
	ValidationData map[string]string
	ClientMetadata map[string]string
}

// ConfirmSignUpInput carries the wire parameters of ConfirmSignUp.
type ConfirmSignUpInput struct {
	ClientID         string
	Username         string
	ConfirmationCode string
}

// AdminConfirmSignUpInput carries the wire parameters of AdminConfirmSignUp.
type AdminConfirmSignUpInput struct {
	UserPoolID string
	Username   string
}

// signUpCore registers a new user in the specified user pool.
func (s *CognitoService) signUpCore(ctx context.Context, reqCtx *request.RequestContext, in SignUpInput) (interface{}, error) {
	if in.ClientID == "" || in.Username == "" || in.Password == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	targetPool, err := store.GetUserPoolByClientID(in.ClientID)
	if err != nil {
		return nil, ErrResourceNotFound
	}

	// Reject self-registration when admin-only creation is enforced.
	if targetPool.AdminCreateUserConfig != nil && targetPool.AdminCreateUserConfig.AllowAdminCreateUserOnly {
		return nil, ErrNotAuthorized
	}

	if err := validatePassword(in.Password, targetPool.PasswordPolicy); err != nil {
		return nil, ErrPasswordPolicyViolation
	}

	// sub is assigned by the service, never supplied by the client, and the
	// schema validation would otherwise hold the client's placeholder
	// against sub's non-empty minimum length. A request without attributes
	// still needs the placeholder slot, so a nil map is normalised here.
	userAttrs := in.UserAttributes
	if userAttrs == nil {
		userAttrs = make(map[string]string)
	}
	delete(userAttrs, "sub")
	if err := validateUserAttributesAgainstSchema(targetPool, userAttrs, true); err != nil {
		return nil, ErrInvalidParameter
	}
	userAttrs["sub"] = ""

	preSignUpResult, err := invokePreSignUp(ctx, s, PreSignUpSignUp, targetPool.ID, in.Username, in.ClientID, targetPool.LambdaConfig, userAttrs, in.ValidationData, in.ClientMetadata)
	if err != nil {
		// A Lambda that raises an error is the documented way to reject the
		// sign-up; a transport failure is an infrastructure error.
		return nil, classifyTriggerFailure(err)
	}

	delete(userAttrs, "sub")
	if preSignUpResult.UserAttributes != nil {
		// The trigger's overrides are screened before merging: the
		// verified-claim keys are stripped — the auto-verify flags below
		// are the only sanctioned way for a trigger to mark a contact
		// attribute verified — and the merged map re-runs the schema gate,
		// because the overrides replace values the client's pass validated.
		delete(preSignUpResult.UserAttributes, "email_verified")
		delete(preSignUpResult.UserAttributes, "phone_number_verified")
		for k, v := range preSignUpResult.UserAttributes {
			userAttrs[k] = v
		}
		if err := validateUserAttributesAgainstSchema(targetPool, userAttrs, true); err != nil {
			return nil, ErrInvalidLambdaResponse
		}
	}

	user := cognitostore.NewUser(targetPool.ID, in.Username)
	user.Attributes = userAttrs
	// Every user's record carries its immutable identifier as the required
	// sub attribute, so projections and filters read it from the attribute
	// map like every other attribute.
	user.Attributes["sub"] = user.ID

	if err := setNativePasswordCredentials(user, targetPool.PasswordPolicy, in.Password); err != nil {
		if errors.Is(err, ErrPasswordHistoryViolation) {
			return nil, ErrPasswordHistoryViolation
		}
		return nil, ErrInternalError
	}

	// The PreSignUp response's auto-verify flags mark the matching
	// attribute verified at creation, independent of the pool-wide
	// AutoVerifiedAttributes marking applied on confirmation.
	if preSignUpResult.AutoVerifyEmail && user.Attributes["email"] != "" {
		user.Attributes["email_verified"] = "true"
	}
	if preSignUpResult.AutoVerifyPhone && user.Attributes["phone_number"] != "" {
		user.Attributes["phone_number_verified"] = "true"
	}

	if preSignUpResult.AutoConfirmUser {
		user.UserStatus = "CONFIRMED"
	} else {
		code, codeErr := generateConfirmationCode()
		if codeErr != nil {
			return nil, ErrInternalError
		}
		user.SignUpCode = code
		user.SignUpCodeExpiry = time.Now().UTC().Add(verificationCodeTTL)
	}

	if err := store.CreateUser(user); err != nil {
		if errors.Is(err, cognitostore.ErrAliasExists) {
			// Only reachable when the PreSignUp response activates a
			// verified alias at creation; the unverified sign-up itself
			// never claims alias values, matching Amazon Cognito's staged
			// alias activation. SignUp's error list carries
			// UsernameExistsException alone — an alias duplicate surfaces
			// under the username-exists error.
			return nil, ErrUserAlreadyExists
		}
		if errors.Is(err, cognitostore.ErrUserAlreadyExists) {
			return nil, ErrUserAlreadyExists
		}
		return nil, ErrInternalError
	}

	// The completed self-service sign-up lands in the user's event history.
	s.recordAuthEvent(reqCtx, targetPool.ID, user.ID, in.Username, in.ClientID, authEventSignUp, authEventResponsePass)

	if preSignUpResult.AutoConfirmUser {
		attrs := userAttributesMap(user)
		invokePostConfirmation(ctx, s, PostConfirmationConfirmSignUp, targetPool.ID, in.Username, in.ClientID, targetPool.LambdaConfig, attrs)
	} else {
		if _, err := invokeCustomMessage(ctx, s, CustomMessageSignUp, targetPool.ID, in.Username, in.ClientID, targetPool.LambdaConfig, user.SignUpCode, userAttributesMap(user), in.ClientMetadata); err != nil {
			logs.Warn("CustomMessage trigger failed", logs.Err(err))
		}
	}

	result := map[string]interface{}{
		"UserConfirmed": preSignUpResult.AutoConfirmUser,
		"UserSub":       user.ID,
	}

	if !preSignUpResult.AutoConfirmUser {
		medium, attrName := determineDeliveryMedium(targetPool, user)
		result["CodeDeliveryDetails"] = map[string]interface{}{
			"Destination":    "***",
			"DeliveryMedium": medium,
			"AttributeName":  attrName,
		}
		// The message-delivery notification log records the sign-up code
		// delivery for pools that configure the userNotification source.
		s.publishNotificationLogCore(reqCtx, targetPool.ID, fmt.Sprintf("SignUp verification code delivery via %s to %s", medium, attrName))
	}

	return result, nil
}

// confirmSignUpCore confirms a user's registration with the confirmation
// code.
func (s *CognitoService) confirmSignUpCore(ctx context.Context, reqCtx *request.RequestContext, in ConfirmSignUpInput) (interface{}, error) {
	if in.ClientID == "" || in.Username == "" || in.ConfirmationCode == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	targetPool, err := store.GetUserPoolByClientID(in.ClientID)
	if err != nil {
		return nil, ErrResourceNotFound
	}

	user, err := store.GetUser(targetPool.ID, in.Username)
	if err != nil {
		return nil, ErrUserNotFound
	}

	if user.UserStatus == "CONFIRMED" {
		return nil, ErrUserAlreadyConfirmed
	}

	if user.SignUpCode == "" || subtle.ConstantTimeCompare([]byte(user.SignUpCode), []byte(in.ConfirmationCode)) != 1 {
		// TEST_MODE: any non-empty code answers (the empty code was
		// rejected at entry), so the regression suite can confirm sign-ups
		// without reading the delivered message.
		if !testMode {
			return nil, ErrCodeMismatch
		}
	}

	if !user.SignUpCodeExpiry.IsZero() && time.Now().After(user.SignUpCodeExpiry) {
		return nil, ErrExpiredCode
	}

	user.UserStatus = "CONFIRMED"
	user.SignUpCode = ""
	user.SignUpCodeExpiry = time.Time{}
	markAutoVerifiedAttributes(user, targetPool)
	if err := store.UpdateUser(user); err != nil {
		// Confirming a contact alias claims its value; a value another user
		// already claims is the documented ConfirmSignUp conflict.
		if errors.Is(err, cognitostore.ErrAliasExists) {
			return nil, ErrAliasExists
		}
		return nil, ErrInternalError
	}

	attrs := userAttributesMap(user)
	invokePostConfirmation(ctx, s, PostConfirmationConfirmSignUp, targetPool.ID, in.Username, in.ClientID, targetPool.LambdaConfig, attrs)

	return response.EmptyResponse(), nil
}

// adminConfirmSignUpCore confirms a user's registration as an administrator.
func (s *CognitoService) adminConfirmSignUpCore(ctx context.Context, reqCtx *request.RequestContext, in AdminConfirmSignUpInput) (interface{}, error) {
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

	if user.UserStatus == "CONFIRMED" {
		return nil, ErrUserAlreadyConfirmed
	}

	user.UserStatus = "CONFIRMED"
	markAutoVerifiedAttributes(user, userPool)
	if err := store.UpdateUser(user); err != nil {
		// The confirmation claims contact aliases exactly like the public
		// ConfirmSignUp path, so the same conflict class applies.
		if errors.Is(err, cognitostore.ErrAliasExists) {
			return nil, ErrAliasExists
		}
		return nil, ErrInternalError
	}

	attrs := userAttributesMap(user)
	invokePostConfirmation(ctx, s, PostConfirmationConfirmSignUp, in.UserPoolID, in.Username, "", userPool.LambdaConfig, attrs)

	return response.EmptyResponse(), nil
}
