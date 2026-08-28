package cognitoidentityprovider

import (
	"context"
	"crypto/subtle"
	"os"
	"time"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	"vorpalstacks/internal/core/logs"
	cognitostore "vorpalstacks/internal/store/aws/cognitoidentityprovider"

	"golang.org/x/crypto/bcrypt"
)

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

	userAttrs := in.UserAttributes
	userAttrs["sub"] = ""

	preSignUpResult, err := invokePreSignUp(ctx, s, PreSignUpSignUp, targetPool.ID, in.Username, in.ClientID, targetPool.LambdaConfig, userAttrs, in.ValidationData, in.ClientMetadata)
	if err != nil {
		return nil, ErrInternalError
	}

	if preSignUpResult.UserAttributes != nil {
		userAttrs = preSignUpResult.UserAttributes
	}
	delete(userAttrs, "sub")

	user := cognitostore.NewUser(targetPool.ID, in.Username)
	user.Attributes = userAttrs

	hash, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, ErrInternalError
	}
	user.PasswordHash = string(hash)
	saltHex, verifierHex, verr := computeSrpVerifier(targetPool.ID, in.Username, in.Password)
	if verr != nil {
		return nil, ErrInternalError
	}
	user.SrpSalt = saltHex
	user.SrpVerifier = verifierHex

	if preSignUpResult.AutoConfirmUser {
		user.UserStatus = "CONFIRMED"
	} else {
		code, codeErr := generateConfirmationCode()
		if codeErr != nil {
			return nil, ErrInternalError
		}
		user.ConfirmationCode = code
		user.ConfirmationCodeExpiry = time.Now().UTC().Add(verificationCodeTTL)
	}

	if err := store.CreateUser(user); err != nil {
		return nil, ErrUserAlreadyExists
	}

	if preSignUpResult.AutoConfirmUser {
		attrs := userAttributesMap(user)
		if err := invokePostConfirmation(ctx, s, PostConfirmationConfirmSignUp, targetPool.ID, in.Username, in.ClientID, targetPool.LambdaConfig, attrs); err != nil {
			logs.Warn("PostConfirmation trigger failed", logs.Err(err))
		}
	} else {
		if _, err := invokeCustomMessage(ctx, s, CustomMessageSignUp, targetPool.ID, in.Username, in.ClientID, targetPool.LambdaConfig, user.ConfirmationCode, userAttributesMap(user), in.ClientMetadata); err != nil {
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

	if user.ConfirmationCode == "" || subtle.ConstantTimeCompare([]byte(user.ConfirmationCode), []byte(in.ConfirmationCode)) != 1 {
		if os.Getenv("TEST_MODE") != "true" {
			return nil, ErrCodeMismatch
		}
		if in.ConfirmationCode == "" {
			return nil, ErrCodeMismatch
		}
	}

	if !user.ConfirmationCodeExpiry.IsZero() && time.Now().After(user.ConfirmationCodeExpiry) {
		return nil, ErrExpiredCode
	}

	user.UserStatus = "CONFIRMED"
	user.ConfirmationCode = ""
	user.ConfirmationCodeExpiry = time.Time{}
	markAutoVerifiedAttributes(user, targetPool)
	if err := store.UpdateUser(user); err != nil {
		return nil, ErrInternalError
	}

	attrs := userAttributesMap(user)
	if err := invokePostConfirmation(ctx, s, PostConfirmationConfirmSignUp, targetPool.ID, in.Username, in.ClientID, targetPool.LambdaConfig, attrs); err != nil {
		logs.Warn("PostConfirmation trigger failed", logs.Err(err))
	}

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
		return nil, ErrInternalError
	}

	attrs := userAttributesMap(user)
	if err := invokePostConfirmation(ctx, s, PostConfirmationConfirmSignUp, in.UserPoolID, in.Username, "", userPool.LambdaConfig, attrs); err != nil {
		logs.Warn("PostConfirmation trigger failed", logs.Err(err))
	}

	return response.EmptyResponse(), nil
}
