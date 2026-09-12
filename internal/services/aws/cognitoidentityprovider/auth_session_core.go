package cognitoidentityprovider

import (
	"context"
	"crypto/subtle"
	"errors"
	"time"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	cognitostore "vorpalstacks/internal/store/aws/cognitoidentityprovider"

	"golang.org/x/crypto/bcrypt"
)

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

	// The token record resolves both the caller and the app client the
	// password change is attributed to in the event history.
	tokenRecord, err := s.validateAccessTokenRecord(reqCtx, in.AccessToken)
	if err != nil {
		return nil, ErrNotAuthorized
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	user, err := store.GetUserByID(tokenRecord.UserID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	userPool, err := store.GetUserPool(user.UserPoolID)
	if err != nil {
		return nil, ErrResourceNotFound
	}

	if err := validatePassword(in.NewPassword, userPool.PasswordPolicy); err != nil {
		return nil, ErrPasswordPolicyViolation
	}

	// The previous-password check and the credential rewrite run as one
	// serialised read-modify-write: with a detached read a concurrent
	// password change could slip between them and have the write overwrite
	// the other writer's history with a stale snapshot.
	err = store.UpdateUserFunc(user.UserPoolID, user.Username, func(u *cognitostore.User) error {
		if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(in.PreviousPassword)); err != nil {
			return ErrIncorrectPassword
		}
		return setNativePasswordCredentials(u, userPool.PasswordPolicy, in.NewPassword)
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

	// The completed password change lands in the user's event history.
	s.recordAuthEvent(reqCtx, user.UserPoolID, user.ID, user.Username, tokenRecord.ClientID, authEventPasswordChange, authEventResponsePass)

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
	user.PasswordResetCode = confirmationCode
	user.PasswordResetExpiry = time.Now().UTC().Add(verificationCodeTTL)
	if err := store.UpdateUser(user); err != nil {
		return nil, ErrInternalError
	}

	// The issued reset code lands in the user's event history.
	s.recordAuthEvent(reqCtx, userPool.ID, user.ID, in.Username, in.ClientID, authEventForgotPassword, authEventResponsePass)

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

	// The existence read keeps the UserNotFoundException precedence ahead of
	// the code and policy checks; the authoritative code, expiry and history
	// checks run against the fresh record below.
	if _, err := store.GetUser(userPool.ID, in.Username); err != nil {
		return nil, ErrUserNotFound
	}

	if err := validatePassword(in.Password, userPool.PasswordPolicy); err != nil {
		return nil, ErrPasswordPolicyViolation
	}

	// The code check and the credential rewrite run as one serialised
	// read-modify-write: the single-use code is consumed, the status becomes
	// CONFIRMED and the password history is written against the freshest
	// record, so a concurrent writer's history cannot be overwritten by a
	// stale snapshot.
	err = store.UpdateUserFunc(userPool.ID, in.Username, func(u *cognitostore.User) error {
		if u.PasswordResetCode == "" || subtle.ConstantTimeCompare([]byte(u.PasswordResetCode), []byte(in.ConfirmationCode)) != 1 {
			return ErrCodeMismatch
		}
		if time.Now().After(u.PasswordResetExpiry) {
			return ErrExpiredCode
		}
		if err := setNativePasswordCredentials(u, userPool.PasswordPolicy, in.Password); err != nil {
			return err
		}
		u.UserStatus = "CONFIRMED"
		u.MigratedAwaitingReset = false
		u.PasswordResetCode = ""
		u.PasswordResetExpiry = time.Time{}
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
