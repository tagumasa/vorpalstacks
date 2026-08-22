package cognitoidentityprovider

import (
	"context"
	"crypto/rand"
	"crypto/subtle"
	"fmt"
	"time"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"

	"golang.org/x/crypto/bcrypt"
)

func generateConfirmationCode() (string, error) {
	const maxCode = 1000000
	const limit = (1 << 24) / maxCode * maxCode
	for {
		b := make([]byte, 3)
		if _, err := rand.Read(b); err != nil {
			return "", err
		}
		n := int(b[0])<<16 | int(b[1])<<8 | int(b[2])
		if n < limit {
			return fmt.Sprintf("%06d", n%maxCode), nil
		}
	}
}

// refreshAuthToken contains the shared refresh-token flow for both InitiateAuth
// and AdminInitiateAuth. It validates the refresh token, looks up the user, and
// issues new access/ID tokens.
func (s *CognitoService) refreshAuthToken(reqCtx *request.RequestContext, userPoolID string, refreshToken string) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	rt, err := store.GetRefreshTokenByValue(refreshToken)
	if err != nil {
		return nil, ErrNotAuthorized
	}

	if time.Now().After(rt.Expires) {
		return nil, ErrNotAuthorized
	}

	// A refresh token is scoped to the user pool that issued it. When the
	// caller names a pool (the Admin flows carry UserPoolId) it must match
	// the token's pool; otherwise the request could mint tokens signed
	// with another pool's keys and issuer.
	if userPoolID != "" && userPoolID != rt.UserPoolID {
		return nil, ErrNotAuthorized
	}
	poolID := rt.UserPoolID

	user, err := store.GetUserByID(rt.UserID)
	if err != nil {
		return nil, ErrNotAuthorized
	}

	attrs := userAttributesMap(user)
	if err := invokePostAuthentication(reqCtx, s, poolID, user.Username, rt.ClientID, nil, attrs, nil); err != nil {
		return nil, fmt.Errorf("PostAuthentication trigger failed: %w", err)
	}

	accessToken, idToken, _, expiresIn, err := s.CreateTokens(reqCtx, poolID, user.ID, rt.ClientID, TokenGenerationRefreshTokens, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create tokens: %w", err)
	}
	return authResultNoRefresh(accessToken, idToken, expiresIn), nil
}

func (s *CognitoService) handleRefreshTokenAuth(reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	refreshToken, err := parseRefreshTokenParam(req)
	if err != nil {
		return nil, err
	}
	return s.refreshAuthToken(reqCtx, "", refreshToken)
}

func (s *CognitoService) handleAdminRefreshTokenAuth(reqCtx *request.RequestContext, req *request.ParsedRequest, userPoolID string) (interface{}, error) {
	refreshToken, err := parseRefreshTokenParam(req)
	if err != nil {
		return nil, err
	}
	return s.refreshAuthToken(reqCtx, userPoolID, refreshToken)
}

// parseRefreshTokenParam extracts the REFRESH_TOKEN from AuthParameters.
func parseRefreshTokenParam(req *request.ParsedRequest) (string, error) {
	authParams := req.Parameters["AuthParameters"]
	if authParams == nil {
		return "", ErrInvalidParameter
	}
	params, ok := authParams.(map[string]interface{})
	if !ok {
		return "", ErrInvalidParameter
	}
	refreshToken, _ := params["REFRESH_TOKEN"].(string)
	if refreshToken == "" {
		return "", ErrInvalidParameter
	}
	return refreshToken, nil
}

// SignOut signs out a user.
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_SignOut.html
func (s *CognitoService) SignOut(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	accessToken := getAccessToken(req)
	if accessToken == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	// SignOut always returns 200 OK per AWS spec, even for invalid or
	// already-revoked access tokens. A client that calls SignOut after token
	// expiry or a previous sign-out receives an empty success.
	at, err := store.GetAccessTokenByValue(accessToken)
	if err != nil {
		return response.EmptyResponse(), nil
	}

	// Best-effort deletion; the token may have been concurrently revoked.
	_ = store.DeleteAccessToken(at.UserPoolID, at.UserID, accessToken)

	return response.EmptyResponse(), nil
}

// GlobalSignOut signs out a user from all devices.
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_GlobalSignOut.html
func (s *CognitoService) GlobalSignOut(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
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
		return nil, ErrNotAuthorized
	}

	if err := store.DeleteUserTokens(user.UserPoolID, user.ID); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// ChangePassword changes the password for a user.
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_ChangePassword.html
func (s *CognitoService) ChangePassword(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	accessToken := getAccessToken(req)
	previousPassword := getPreviousPassword(req)
	newPassword := getNewPassword(req)

	if accessToken == "" || previousPassword == "" || newPassword == "" {
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

	userPool, err := store.GetUserPool(user.UserPoolID)
	if err != nil {
		return nil, ErrResourceNotFound
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(previousPassword)); err != nil {
		return nil, ErrIncorrectPassword
	}

	if err := validatePassword(newPassword, userPool.PasswordPolicy); err != nil {
		return nil, ErrPasswordPolicyViolation
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, ErrInternalError
	}
	user.PasswordHash = string(hash)
	saltHex, verifierHex, verr := computeSrpVerifier(user.UserPoolID, user.Username, newPassword)
	if verr != nil {
		return nil, ErrInternalError
	}
	user.SrpSalt = saltHex
	user.SrpVerifier = verifierHex

	if err := store.UpdateUser(user); err != nil {
		return nil, ErrInternalError
	}

	return response.EmptyResponse(), nil
}

// ForgotPassword initiates the forgot password flow.
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_ForgotPassword.html
func (s *CognitoService) ForgotPassword(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	clientID := getClientId(req)
	username := getUsername(req)

	if clientID == "" || username == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	userPool, err := store.GetUserPoolByClientID(clientID)
	if err != nil {
		return nil, ErrResourceNotFound
	}

	user, err := store.GetUser(userPool.ID, username)
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

	attrs := userAttributesMap(user)
	customMsg, _ := invokeCustomMessage(ctx, s, CustomMessageForgotPassword, userPool.ID, username, clientID, userPool.LambdaConfig, "####", attrs, nil)

	codeDeliveryDetails := map[string]interface{}{
		"Destination":    "***",
		"DeliveryMedium": "EMAIL",
		"AttributeName":  "email",
	}
	if customMsg != nil && customMsg.EmailSubject != "" {
		codeDeliveryDetails["_customEmailSubject"] = customMsg.EmailSubject
	}

	return map[string]interface{}{
		"CodeDeliveryDetails": codeDeliveryDetails,
	}, nil
}

// ConfirmForgotPassword confirms the forgot password flow.
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_ConfirmForgotPassword.html
func (s *CognitoService) ConfirmForgotPassword(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	clientID := getClientId(req)
	username := getUsername(req)
	password := getPassword(req)
	confirmationCode := getConfirmationCode(req)

	if clientID == "" || username == "" || password == "" || confirmationCode == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	userPool, err := store.GetUserPoolByClientID(clientID)
	if err != nil {
		return nil, ErrResourceNotFound
	}

	user, err := store.GetUser(userPool.ID, username)
	if err != nil {
		return nil, ErrUserNotFound
	}

	if user.ConfirmationCode == "" || subtle.ConstantTimeCompare([]byte(user.ConfirmationCode), []byte(confirmationCode)) != 1 {
		return nil, ErrCodeMismatch
	}

	if time.Now().After(user.ConfirmationCodeExpiry) {
		return nil, ErrExpiredCode
	}

	if err := validatePassword(password, userPool.PasswordPolicy); err != nil {
		return nil, ErrPasswordPolicyViolation
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, ErrInternalError
	}
	user.PasswordHash = string(hash)
	saltHex, verifierHex, verr := computeSrpVerifier(user.UserPoolID, user.Username, password)
	if verr != nil {
		return nil, ErrInternalError
	}
	user.SrpSalt = saltHex
	user.SrpVerifier = verifierHex
	user.UserStatus = "CONFIRMED"
	user.ConfirmationCode = ""
	user.ConfirmationCodeExpiry = time.Time{}

	if err := store.UpdateUser(user); err != nil {
		return nil, ErrInternalError
	}

	return response.EmptyResponse(), nil
}
