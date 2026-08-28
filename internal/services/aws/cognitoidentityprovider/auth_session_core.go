package cognitoidentityprovider

import (
	"context"
	"crypto/subtle"
	"time"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"

	"golang.org/x/crypto/bcrypt"
)

// SignOutInput carries the wire parameters of SignOut.
type SignOutInput struct {
	AccessToken string
}

// GlobalSignOutInput carries the wire parameters of GlobalSignOut.
type GlobalSignOutInput struct {
	AccessToken string
}

// ChangePasswordInput carries the wire parameters of ChangePassword.
type ChangePasswordInput struct {
	AccessToken      string
	PreviousPassword string
	NewPassword      string
}

// ForgotPasswordInput carries the wire parameters of ForgotPassword.
type ForgotPasswordInput struct {
	ClientID string
	Username string
}

// ConfirmForgotPasswordInput carries the wire parameters of
// ConfirmForgotPassword.
type ConfirmForgotPasswordInput struct {
	ClientID         string
	Username         string
	Password         string
	ConfirmationCode string
}

// signOutCore revokes the caller's access token. SignOut always returns
// 200 OK per AWS spec, even for invalid or already-revoked access tokens. A
// client that calls SignOut after token expiry or a previous sign-out
// receives an empty success.
func (s *CognitoService) signOutCore(reqCtx *request.RequestContext, in SignOutInput) (interface{}, error) {
	if in.AccessToken == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	at, err := store.GetAccessTokenByValue(in.AccessToken)
	if err != nil {
		return response.EmptyResponse(), nil
	}

	// Best-effort deletion; the token may have been concurrently revoked.
	_ = store.DeleteAccessToken(at.UserPoolID, at.UserID, in.AccessToken)

	return response.EmptyResponse(), nil
}

// globalSignOutCore revokes every token minted for the caller.
func (s *CognitoService) globalSignOutCore(reqCtx *request.RequestContext, in GlobalSignOutInput) (interface{}, error) {
	if in.AccessToken == "" {
		return nil, ErrInvalidParameter
	}

	userID, err := s.ValidateAccessToken(reqCtx, in.AccessToken)
	if err != nil {
		return nil, ErrNotAuthorized
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	user, err := store.GetUserByID(userID)
	if err != nil {
		return nil, ErrNotAuthorized
	}

	if err := store.DeleteUserTokens(user.UserPoolID, user.ID); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// changePasswordCore verifies the previous password, enforces the pool
// password policy on the new password and rewrites the native credentials.
func (s *CognitoService) changePasswordCore(reqCtx *request.RequestContext, in ChangePasswordInput) (interface{}, error) {
	if in.AccessToken == "" || in.PreviousPassword == "" || in.NewPassword == "" {
		return nil, ErrInvalidParameter
	}

	userID, err := s.ValidateAccessToken(reqCtx, in.AccessToken)
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

	userPool, err := store.GetUserPool(user.UserPoolID)
	if err != nil {
		return nil, ErrResourceNotFound
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(in.PreviousPassword)); err != nil {
		return nil, ErrIncorrectPassword
	}

	if err := validatePassword(in.NewPassword, userPool.PasswordPolicy); err != nil {
		return nil, ErrPasswordPolicyViolation
	}

	if err := setNativePasswordCredentials(user, user.UserPoolID, user.Username, in.NewPassword); err != nil {
		return nil, ErrInternalError
	}

	if err := store.UpdateUser(user); err != nil {
		return nil, ErrInternalError
	}

	return response.EmptyResponse(), nil
}

// forgotPasswordCore issues a confirmation code for the password-reset
// flow. The response is masked so the operation cannot be used to enumerate
// accounts; unknown users receive the same CodeDeliveryDetails shape.
func (s *CognitoService) forgotPasswordCore(ctx context.Context, reqCtx *request.RequestContext, in ForgotPasswordInput) (interface{}, error) {
	if in.ClientID == "" || in.Username == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	userPool, err := store.GetUserPoolByClientID(in.ClientID)
	if err != nil {
		return nil, ErrResourceNotFound
	}

	user, err := store.GetUser(userPool.ID, in.Username)
	if err != nil {
		// Return a masked CodeDeliveryDetails to prevent user enumeration.
		return map[string]interface{}{
			"CodeDeliveryDetails": map[string]interface{}{
				"Destination":    "***",
				"DeliveryMedium": "EMAIL",
				"AttributeName":  "email",
			},
		}, nil
	}

	confirmationCode, err := generateConfirmationCode()
	if err != nil {
		return nil, ErrInternalError
	}
	user.ConfirmationCode = confirmationCode
	user.ConfirmationCodeExpiry = time.Now().UTC().Add(verificationCodeTTL)
	if err := store.UpdateUser(user); err != nil {
		return nil, ErrInternalError
	}

	// The CustomMessage ForgotPassword trigger fires for its side effects
	// (custom email content); the documented response carries only the
	// masked CodeDeliveryDetails shape.
	attrs := userAttributesMap(user)
	_, _ = invokeCustomMessage(ctx, s, CustomMessageForgotPassword, userPool.ID, in.Username, in.ClientID, userPool.LambdaConfig, "####", attrs, nil)

	return map[string]interface{}{
		"CodeDeliveryDetails": map[string]interface{}{
			"Destination":    "***",
			"DeliveryMedium": "EMAIL",
			"AttributeName":  "email",
		},
	}, nil
}

// confirmForgotPasswordCore verifies the confirmation code and replaces the
// user's password.
func (s *CognitoService) confirmForgotPasswordCore(reqCtx *request.RequestContext, in ConfirmForgotPasswordInput) (interface{}, error) {
	if in.ClientID == "" || in.Username == "" || in.Password == "" || in.ConfirmationCode == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	userPool, err := store.GetUserPoolByClientID(in.ClientID)
	if err != nil {
		return nil, ErrResourceNotFound
	}

	user, err := store.GetUser(userPool.ID, in.Username)
	if err != nil {
		return nil, ErrUserNotFound
	}

	if user.ConfirmationCode == "" || subtle.ConstantTimeCompare([]byte(user.ConfirmationCode), []byte(in.ConfirmationCode)) != 1 {
		return nil, ErrCodeMismatch
	}

	if time.Now().After(user.ConfirmationCodeExpiry) {
		return nil, ErrExpiredCode
	}

	if err := validatePassword(in.Password, userPool.PasswordPolicy); err != nil {
		return nil, ErrPasswordPolicyViolation
	}

	if err := setNativePasswordCredentials(user, user.UserPoolID, user.Username, in.Password); err != nil {
		return nil, ErrInternalError
	}
	user.UserStatus = "CONFIRMED"
	user.ConfirmationCode = ""
	user.ConfirmationCodeExpiry = time.Time{}

	if err := store.UpdateUser(user); err != nil {
		return nil, ErrInternalError
	}

	return response.EmptyResponse(), nil
}
