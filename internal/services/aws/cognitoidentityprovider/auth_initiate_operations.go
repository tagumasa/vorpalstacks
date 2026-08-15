package cognitoidentityprovider

import (
	"context"
	"errors"
	"fmt"
	"time"

	"vorpalstacks/internal/common/request"
	cognitostore "vorpalstacks/internal/store/aws/cognitoidentityprovider"

	"golang.org/x/crypto/bcrypt"
)

// InitiateAuth initiates the authentication flow.
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_InitiateAuth.html
func (s *CognitoService) InitiateAuth(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	authFlow := req.GetParam("AuthFlow")
	clientID := getClientId(req)

	if authFlow == "" || clientID == "" {
		return nil, ErrInvalidParameter
	}
	if !validateInitiateAuthFlow(authFlow) {
		return nil, ErrInvalidParameter
	}

	switch authFlow {
	case "USER_PASSWORD_AUTH":
		return s.handleUserPasswordAuth(ctx, reqCtx, req)
	case "USER_SRP_AUTH":
		return s.handleUserSrpAuth(reqCtx, req)
	case "REFRESH_TOKEN_AUTH", "REFRESH_TOKEN":
		return s.handleRefreshTokenAuth(reqCtx, req)
	case "CUSTOM_AUTH":
		return s.handleCustomAuth(ctx, reqCtx, req)
	case "USER_AUTH":
		return s.handleUserAuth(ctx, reqCtx, req)
	default:
		return nil, ErrInvalidParameter
	}
}

// customFlowChallengeNames lists the challenge types a DefineAuthChallenge
// Lambda response may nominate as the next challenge in a custom
// authentication flow. Anything outside this set is a malformed response.
// Per the Amazon Cognito developer guide, the DefineAuthChallenge response
// designates only the flow-level challenges (CUSTOM_CHALLENGE, SRP_A,
// PASSWORD_VERIFIER); MFA challenges — including MFA_SETUP — are issued by
// the service itself after credential verification ("You don't need to
// invoke any MFA challenges in your define auth challenge function").
// Accepting MFA_SETUP here would mint an MFA_SETUP-typed session without
// any verified credential, and AssociateSoftwareToken accepts exactly that
// session type, so an unauthenticated caller could overwrite a victim's
// TOTP configuration.
var customFlowChallengeNames = map[string]bool{
	"CUSTOM_CHALLENGE":  true,
	"SRP_A":             true,
	"PASSWORD_VERIFIER": true,
}

// resolveCustomFlowChallenge applies the customFlowChallengeNames contract
// to a DefineAuthChallenge response and returns the challenge the flow
// continues with. A nil or empty response keeps the CUSTOM_CHALLENGE
// default; a name outside the set is a malformed Lambda response.
func resolveCustomFlowChallenge(lambdaResult map[string]interface{}) (string, error) {
	if lambdaResult == nil {
		return "CUSTOM_CHALLENGE", nil
	}
	cn, ok := lambdaResult["challengeName"].(string)
	if !ok || cn == "" {
		return "CUSTOM_CHALLENGE", nil
	}
	if !customFlowChallengeNames[cn] {
		return "", ErrInvalidLambdaResponse
	}
	return cn, nil
}

// handleCustomAuth implements the CUSTOM_AUTH flow. The client provides a
// USERNAME (and optionally a PASSWORD for the initial verification). The
// server invokes the DefineAuthChallenge and CreateAuthChallenge Lambda
// triggers to produce a custom challenge that the client must answer via
// RespondToAuthChallenge.
func (s *CognitoService) handleCustomAuth(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	clientID := getClientId(req)
	username := getUsername(req)
	password := getPassword(req)

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
		// Return the same NotAuthorized error used by other auth flows to
		// prevent user enumeration via distinguishable session or error.
		return nil, ErrNotAuthorized
	}

	// If a password is provided, verify it as part of the initial auth.
	if password != "" {
		if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
			return nil, ErrNotAuthorized
		}
	}

	// Invoke DefineAuthChallenge Lambda trigger to determine the challenge.
	var lambdaResult map[string]interface{}
	if userPool.LambdaConfig != nil && userPool.LambdaConfig.DefineAuthChallenge != "" {
		lambdaResult, _ = s.invokeTrigger(ctx, DefineAuthChallenge, userPool.ID, username, clientID,
			userPool.LambdaConfig.DefineAuthChallenge,
			map[string]interface{}{
				"userAttributes":  userAttributesMap(user),
				"challengeResult": "FAILED",
			},
			map[string]interface{}{
				"challengeName":      "",
				"issueTokens":        false,
				"failAuthentication": false,
			},
			true,
		)
	}

	sessionID := generateSessionID()
	challengeParams := map[string]string{
		"USERNAME": username,
	}

	challengeName, err := resolveCustomFlowChallenge(lambdaResult)
	if err != nil {
		return nil, err
	}

	// Invoke CreateAuthChallenge Lambda trigger to produce challenge parameters.
	if userPool.LambdaConfig != nil && userPool.LambdaConfig.CreateAuthChallenge != "" {
		createResult, _ := s.invokeTrigger(ctx, CreateAuthChallenge, userPool.ID, username, clientID,
			userPool.LambdaConfig.CreateAuthChallenge,
			map[string]interface{}{
				"userAttributes": userAttributesMap(user),
				"challengeName":  challengeName,
			},
			map[string]interface{}{
				"publicChallengeParameters": map[string]string{},
			},
			true,
		)
		if createResult != nil {
			if params, ok := createResult["publicChallengeParameters"].(map[string]interface{}); ok {
				for k, v := range params {
					if vs, ok := v.(string); ok {
						challengeParams[k] = vs
					}
				}
			}
		}
	}

	challengeSession := &cognitostore.ChallengeSession{
		SessionID:     sessionID,
		UserPoolID:    userPool.ID,
		ClientID:      clientID,
		Username:      username,
		ChallengeName: challengeName,
		CreatedAt:     time.Now().UTC(),
		ExpiresAt:     time.Now().UTC().Add(3 * time.Minute),
	}
	if err := store.SaveChallengeSession(challengeSession); err != nil {
		return nil, ErrInternalError
	}

	return map[string]interface{}{
		"ChallengeName":       challengeName,
		"Session":             sessionID,
		"ChallengeParameters": challengeParams,
	}, nil
}

// handleUserAuth implements the USER_AUTH flow (the modern recommended auth
// flow). The server inspects the user's configured auth factors and returns
// the set of available challenges for the client to choose from.
func (s *CognitoService) handleUserAuth(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	clientID := getClientId(req)
	username := getUsername(req)

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	userPool, err := store.GetUserPoolByClientID(clientID)
	if err != nil {
		return nil, ErrResourceNotFound
	}

	user, err := store.GetUser(userPool.ID, username)

	// The selector session is issued regardless of whether the user
	// exists, so InitiateAuth(USER_AUTH) does not leak account existence;
	// every subsequent challenge response resolves the user again and
	// fails uniformly with NotAuthorized for unknown users. The contents
	// of AvailableChallenges reflect the user's enrolled factors, which is
	// the documented behaviour of choice-based authentication.
	available := make([]string, 0, 5)
	available = append(available, "PASSWORD")
	if userPool.WebAuthnConfiguration != nil && userPool.WebAuthnConfiguration.RelyingPartyId != "" {
		available = append(available, "WEB_AUTHN")
	}

	if err == nil && user != nil {
		if user.SoftwareTokenMfa != nil && user.SoftwareTokenMfa.Verified {
			available = append(available, "SOFTWARE_TOKEN_MFA")
		}
		if isAttributeVerified(user.Attributes, "phone_number") {
			available = append(available, "SMS_OTP")
		}
		if isAttributeVerified(user.Attributes, "email") {
			available = append(available, "EMAIL_OTP")
		}
	}

	sessionID := generateSessionID()
	challengeSession := &cognitostore.ChallengeSession{
		SessionID:     sessionID,
		UserPoolID:    userPool.ID,
		ClientID:      clientID,
		Username:      username,
		ChallengeName: "SELECT_CHALLENGE",
		CreatedAt:     time.Now().UTC(),
		ExpiresAt:     time.Now().UTC().Add(3 * time.Minute),
	}
	if err := store.SaveChallengeSession(challengeSession); err != nil {
		return nil, ErrInternalError
	}

	return map[string]interface{}{
		"AvailableChallenges": available,
		"Session":             sessionID,
	}, nil
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
	validationData, clientMetadata map[string]string,
) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	user, err := store.GetUser(userPoolID, username)
	if err != nil {
		migrationResult, migrationErr := invokeUserMigration(ctx, s, userPoolID, username, clientID, password, lambdaConfig, validationData, clientMetadata)
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
			saltHex, verifierHex, verr := computeSrpVerifier(userPoolID, migratedUser.Username, password)
			if verr != nil {
				return nil, ErrInternalError
			}
			migratedUser.SrpSalt = saltHex
			migratedUser.SrpVerifier = verifierHex
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

	// The password must verify before any challenge is issued. Amazon
	// Cognito returns NEW_PASSWORD_REQUIRED only for users who signed in
	// successfully with their temporary password; issuing it on user status
	// alone would let anyone reset a FORCE_CHANGE_PASSWORD or RESET_REQUIRED
	// account by username only.
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		s.recordAuthEvent(reqCtx, userPoolID, user.ID, "SignIn", "Fail")
		return nil, ErrIncorrectPassword
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
	if err := invokePreAuthentication(ctx, s, userPoolID, username, clientID, lambdaConfig, attrs, clientMetadata); err != nil {
		s.recordAuthEvent(reqCtx, userPoolID, user.ID, "SignIn", "Fail")
		return nil, fmt.Errorf("PreAuthentication trigger failed: %w", err)
	}

	if err := invokePostAuthentication(ctx, s, userPoolID, username, clientID, lambdaConfig, attrs, clientMetadata); err != nil {
		return nil, fmt.Errorf("PostAuthentication trigger failed: %w", err)
	}

	accessToken, idToken, refreshToken, expiresIn, err := s.CreateTokens(reqCtx, userPoolID, user.ID, clientID, TokenGenerationAuthentication, clientMetadata)
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

	return s.authenticateUser(ctx, reqCtx, userPool.ID, getClientId(req), username, password, userPool.LambdaConfig, parseValidationData(req), parseClientMetadata(req))
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
	if !validateAdminInitiateAuthFlow(authFlow) {
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
	case "CUSTOM_AUTH":
		return s.handleCustomAuth(ctx, reqCtx, req)
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

	return s.authenticateUser(ctx, reqCtx, userPoolID, clientID, username, password, lambdaConfig, parseValidationData(req), parseClientMetadata(req))
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
