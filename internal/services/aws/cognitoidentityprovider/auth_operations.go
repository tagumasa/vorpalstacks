package cognitoidentityprovider

import (
	"context"
	"crypto/rand"
	"crypto/subtle"
	"errors"
	"fmt"
	"time"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	cognitostore "vorpalstacks/internal/store/aws/cognitoidentityprovider"

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

// InitiateAuth initiates the authentication flow.
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_InitiateAuth.html
func (s *CognitoService) InitiateAuth(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	authFlow := req.GetParam("AuthFlow")
	clientID := getClientId(req)

	if authFlow == "" || clientID == "" {
		return nil, ErrInvalidParameter
	}

	switch authFlow {
	case "USER_PASSWORD_AUTH":
		return s.handleUserPasswordAuth(ctx, reqCtx, req)
	case "USER_SRP_AUTH":
		return nil, ErrInvalidParameter
	case "REFRESH_TOKEN_AUTH", "REFRESH_TOKEN":
		return s.handleRefreshTokenAuth(reqCtx, req)
	default:
		return nil, ErrInvalidParameter
	}
}

// authenticateUser contains the shared authentication logic used by both
// USER_PASSWORD_AUTH (InitiateAuth) and ADMIN_NO_SRP_AUTH (AdminInitiateAuth).
// It handles user lookup, migration, challenge detection, credential
// verification, trigger invocation, and token generation.
func (s *CognitoService) authenticateUser(
	ctx context.Context,
	reqCtx *request.RequestContext,
	userPoolID, clientID, username, password string,
	lambdaConfig *cognitostore.LambdaConfig,
) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	user, err := store.GetUser(userPoolID, username)
	if err != nil {
		migrationResult, migrationErr := invokeUserMigration(ctx, s, userPoolID, username, clientID, password, lambdaConfig)
		if migrationErr != nil || migrationResult == nil {
			return nil, ErrIncorrectPassword
		}

		migratedUser := cognitostore.NewUser(userPoolID, username)
		if migrationResult.UserAttributes != nil {
			migratedUser.Attributes = migrationResult.UserAttributes
		}
		if migrationResult.FinalUserStatus != "" {
			migratedUser.UserStatus = migrationResult.FinalUserStatus
		} else {
			migratedUser.UserStatus = "CONFIRMED"
		}
		if migrationResult.MessageAction != "SUPPRESS" && password != "" {
			hash, hashErr := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
			if hashErr != nil {
				return nil, ErrInternalError
			}
			migratedUser.PasswordHash = string(hash)
		}
		if err := store.CreateUser(migratedUser); err != nil {
			if errors.Is(err, cognitostore.ErrUserAlreadyExists) {
				return nil, ErrIncorrectPassword
			}
			return nil, ErrInternalError
		}
		user = migratedUser
	}

	if !user.Enabled {
		s.recordAuthEvent(reqCtx, userPoolID, user.ID, "SignIn", "Fail")
		return nil, ErrNotAuthorized
	}

	if user.UserStatus == "FORCE_CHANGE_PASSWORD" || user.UserStatus == "RESET_REQUIRED" {
		s.recordAuthEvent(reqCtx, userPoolID, user.ID, "SignIn", "InProgress")
		return s.newPasswordChallenge(reqCtx, userPoolID, clientID, user)
	}

	if user.UserStatus != "CONFIRMED" {
		s.recordAuthEvent(reqCtx, userPoolID, user.ID, "SignIn", "Fail")
		return nil, ErrUserNotConfirmed
	}

	attrs := userAttributesMap(user)
	if err := invokePreAuthentication(ctx, s, userPoolID, username, clientID, lambdaConfig, attrs); err != nil {
		s.recordAuthEvent(reqCtx, userPoolID, user.ID, "SignIn", "Fail")
		return nil, fmt.Errorf("PreAuthentication trigger failed: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		s.recordAuthEvent(reqCtx, userPoolID, user.ID, "SignIn", "Fail")
		return nil, ErrIncorrectPassword
	}

	if err := invokePostAuthentication(ctx, s, userPoolID, username, clientID, lambdaConfig, attrs); err != nil {
		return nil, fmt.Errorf("PostAuthentication trigger failed: %w", err)
	}

	accessToken, idToken, refreshToken, expiresIn, err := s.CreateTokens(reqCtx, userPoolID, user.ID, clientID)
	if err != nil {
		return nil, fmt.Errorf("failed to create tokens: %w", err)
	}
	s.recordAuthEvent(reqCtx, userPoolID, user.ID, "SignIn", "Pass")
	return authResult(accessToken, idToken, refreshToken, expiresIn), nil
}

// newPasswordChallenge creates a NEW_PASSWORD_REQUIRED challenge session and
// returns the challenge response.
func (s *CognitoService) newPasswordChallenge(reqCtx *request.RequestContext, userPoolID, clientID string, user *cognitostore.User) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	session := generateSessionID()
	challengeSession := &cognitostore.ChallengeSession{
		SessionID:     session,
		UserPoolID:    userPoolID,
		ClientID:      clientID,
		Username:      user.Username,
		ChallengeName: "NEW_PASSWORD_REQUIRED",
		CreatedAt:     time.Now().UTC(),
		ExpiresAt:     time.Now().UTC().Add(5 * time.Minute),
	}
	if err := store.SaveChallengeSession(challengeSession); err != nil {
		return nil, ErrInternalError
	}
	return map[string]interface{}{
		"ChallengeName": "NEW_PASSWORD_REQUIRED",
		"Session":       session,
		"ChallengeParameters": map[string]interface{}{
			"USER_ID_FOR_SRP":    user.Username,
			"requiredAttributes": "[]",
		},
	}, nil
}

func authResult(accessToken, idToken, refreshToken string, expiresIn int64) map[string]interface{} {
	return map[string]interface{}{
		"AuthenticationResult": map[string]interface{}{
			"AccessToken":  accessToken,
			"IdToken":      idToken,
			"RefreshToken": refreshToken,
			"TokenType":    "Bearer",
			"ExpiresIn":    expiresIn,
		},
	}
}

func authResultNoRefresh(accessToken, idToken string, expiresIn int64) map[string]interface{} {
	return map[string]interface{}{
		"AuthenticationResult": map[string]interface{}{
			"AccessToken": accessToken,
			"IdToken":     idToken,
			"TokenType":   "Bearer",
			"ExpiresIn":   expiresIn,
		},
	}
}

func (s *CognitoService) handleUserPasswordAuth(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	username, password, err := parseAuthParams(req)
	if err != nil {
		return nil, err
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	userPool, err := store.GetUserPoolByClientID(getClientId(req))
	if err != nil {
		return nil, ErrResourceNotFound
	}

	return s.authenticateUser(ctx, reqCtx, userPool.ID, getClientId(req), username, password, userPool.LambdaConfig)
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

	poolID := userPoolID
	if poolID == "" {
		poolID = rt.UserPoolID
	}

	user, err := store.GetUserByID(rt.UserID)
	if err != nil {
		return nil, ErrNotAuthorized
	}

	attrs := userAttributesMap(user)
	if err := invokePostAuthentication(reqCtx, s, poolID, user.Username, rt.ClientID, nil, attrs); err != nil {
		return nil, fmt.Errorf("PostAuthentication trigger failed: %w", err)
	}

	accessToken, idToken, _, expiresIn, err := s.CreateTokens(reqCtx, poolID, user.ID, rt.ClientID)
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

// RespondToAuthChallenge responds to an authentication challenge.
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_RespondToAuthChallenge.html
func (s *CognitoService) RespondToAuthChallenge(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	clientID := getClientId(req)
	challengeName := req.GetParam("ChallengeName")
	session := req.GetParam("Session")

	if clientID == "" || challengeName == "" {
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

	if challengeName == "NEW_PASSWORD_REQUIRED" {
		return s.respondToNewPasswordChallenge(reqCtx, req, userPool.ID, clientID, session)
	}

	return nil, ErrInvalidParameter
}

// parseAuthParams extracts USERNAME and PASSWORD from the AuthParameters
// block of an InitiateAuth or AdminInitiateAuth request.
func parseAuthParams(req *request.ParsedRequest) (username, password string, err error) {
	authParams := req.Parameters["AuthParameters"]
	if authParams == nil {
		return "", "", ErrInvalidParameter
	}
	params, ok := authParams.(map[string]interface{})
	if !ok {
		return "", "", ErrInvalidParameter
	}
	username, _ = params["USERNAME"].(string)
	password, _ = params["PASSWORD"].(string)
	if username == "" || password == "" {
		return "", "", ErrInvalidParameter
	}
	return username, password, nil
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
	at, err := store.GetAccessTokenByValue(accessToken)
	if err != nil {
		return response.EmptyResponse(), nil
	}

	if err := store.DeleteAccessToken(at.UserPoolID, at.UserID, accessToken); err != nil {
		return nil, err
	}

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

	if err := store.DeleteAllRefreshTokensForUser(user.UserPoolID, user.ID); err != nil {
		return nil, err
	}
	if err := store.DeleteAccessToken(user.UserPoolID, user.ID, accessToken); err != nil {
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
		return nil, ErrUserNotFound
	}

	confirmationCode, err := generateConfirmationCode()
	if err != nil {
		return nil, ErrInternalError
	}
	user.ConfirmationCode = confirmationCode
	user.ConfirmationCodeExpiry = time.Now().Add(24 * time.Hour)
	if err := store.UpdateUser(user); err != nil {
		return nil, ErrInternalError
	}

	attrs := userAttributesMap(user)
	customMsg, _ := invokeCustomMessage(ctx, s, CustomMessageForgotPassword, userPool.ID, username, clientID, userPool.LambdaConfig, "####", attrs)

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
	user.UserStatus = "CONFIRMED"
	user.ConfirmationCode = ""
	user.ConfirmationCodeExpiry = time.Time{}

	if err := store.UpdateUser(user); err != nil {
		return nil, ErrInternalError
	}

	return response.EmptyResponse(), nil
}

// AdminInitiateAuth initiates the admin authentication flow.
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_AdminInitiateAuth.html
func (s *CognitoService) AdminInitiateAuth(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	userPoolID := getUserPoolID(req)
	clientID := getClientId(req)
	authFlow := req.GetParam("AuthFlow")

	if userPoolID == "" || clientID == "" || authFlow == "" {
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

	switch authFlow {
	case "ADMIN_NO_SRP_AUTH", "ADMIN_USER_PASSWORD_AUTH":
		return s.handleAdminNoSrpAuth(ctx, reqCtx, req, userPoolID, clientID, userPool.LambdaConfig)
	case "REFRESH_TOKEN_AUTH", "REFRESH_TOKEN":
		return s.handleAdminRefreshTokenAuth(reqCtx, req, userPoolID)
	default:
		return nil, ErrInvalidParameter
	}
}

func (s *CognitoService) handleAdminNoSrpAuth(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest, userPoolID, clientID string, lambdaConfig *cognitostore.LambdaConfig) (interface{}, error) {
	username, password, err := parseAuthParams(req)
	if err != nil {
		return nil, err
	}

	return s.authenticateUser(ctx, reqCtx, userPoolID, clientID, username, password, lambdaConfig)
}

// AdminRespondToAuthChallenge responds to an admin authentication challenge.
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_AdminRespondToAuthChallenge.html
func (s *CognitoService) AdminRespondToAuthChallenge(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	userPoolID := getUserPoolID(req)
	clientID := getClientId(req)
	challengeName := req.GetParam("ChallengeName")
	session := req.GetParam("Session")

	if userPoolID == "" || clientID == "" || challengeName == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	_, err = store.GetUserPool(userPoolID)
	if err != nil {
		return nil, ErrResourceNotFound
	}

	if challengeName == "NEW_PASSWORD_REQUIRED" {
		return s.respondToNewPasswordChallenge(reqCtx, req, userPoolID, clientID, session)
	}

	return nil, ErrInvalidParameter
}

// respondToNewPasswordChallenge handles the NEW_PASSWORD_REQUIRED challenge
// for both RespondToAuthChallenge and AdminRespondToAuthChallenge.
func (s *CognitoService) respondToNewPasswordChallenge(reqCtx *request.RequestContext, req *request.ParsedRequest, userPoolID, clientID, session string) (interface{}, error) {
	challengeResponses := req.Parameters["ChallengeResponses"]
	if challengeResponses == nil {
		return nil, ErrInvalidParameter
	}

	params, ok := challengeResponses.(map[string]interface{})
	if !ok {
		return nil, ErrInvalidParameter
	}

	username, _ := params["USERNAME"].(string)
	newPassword, _ := params["NEW_PASSWORD"].(string)

	if username == "" || newPassword == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	challengeSession, err := store.GetChallengeSession(session)
	if err != nil || challengeSession == nil {
		return nil, ErrNotAuthorized
	}
	if challengeSession.Username != username {
		return nil, ErrNotAuthorized
	}
	if !challengeSession.ExpiresAt.IsZero() && time.Now().UTC().After(challengeSession.ExpiresAt) {
		_ = store.DeleteChallengeSession(session)
		return nil, ErrNotAuthorized
	}

	userPool, err := store.GetUserPool(userPoolID)
	if err != nil {
		return nil, ErrResourceNotFound
	}

	user, err := store.GetUser(userPoolID, username)
	if err != nil {
		return nil, ErrUserNotFound
	}

	if err := validatePassword(newPassword, userPool.PasswordPolicy); err != nil {
		return nil, ErrPasswordPolicyViolation
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, ErrInternalError
	}
	user.PasswordHash = string(hash)
	user.UserStatus = "CONFIRMED"

	if err := store.UpdateUser(user); err != nil {
		return nil, ErrInternalError
	}

	accessToken, idToken, refreshToken, expiresIn, err := s.CreateTokens(reqCtx, userPoolID, user.ID, clientID)
	if err != nil {
		return nil, fmt.Errorf("failed to create tokens: %w", err)
	}

	return map[string]interface{}{
		"AuthenticationResult": map[string]interface{}{
			"AccessToken":  accessToken,
			"IdToken":      idToken,
			"RefreshToken": refreshToken,
			"TokenType":    "Bearer",
			"ExpiresIn":    expiresIn,
		},
	}, nil
}
