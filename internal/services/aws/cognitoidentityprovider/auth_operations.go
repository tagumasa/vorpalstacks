package cognitoidentityprovider

import (
	"context"
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"math/big"
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

	challengeName := "CUSTOM_CHALLENGE"
	if lambdaResult != nil {
		if cn, ok := lambdaResult["challengeName"].(string); ok && cn != "" {
			challengeName = cn
		}
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
	if err != nil {
		// Return a minimal response to avoid user enumeration.
		return map[string]interface{}{
			"AvailableChallenges": []interface{}{},
		}, nil
	}

	available := make([]string, 0, 4)
	available = append(available, "PASSWORD")
	if user.SoftwareTokenMfa != nil && user.SoftwareTokenMfa.Verified {
		available = append(available, "SOFTWARE_TOKEN_MFA")
	}
	if isAttributeVerified(user.Attributes, "phone_number") {
		available = append(available, "SMS_OTP")
	}
	if isAttributeVerified(user.Attributes, "email") {
		available = append(available, "EMAIL_OTP")
	}

	return map[string]interface{}{
		"AvailableChallenges": available,
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

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		s.recordAuthEvent(reqCtx, userPoolID, user.ID, "SignIn", "Fail")
		return nil, ErrIncorrectPassword
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

// handleUserSrpAuth handles the InitiateAuth USER_SRP_AUTH flow. It receives
// the client's SRP_A value, looks up the user, generates the server's
// ephemeral B and a fresh SECRET_BLOCK, persists the SRP session state, and
// returns a PASSWORD_VERIFIER challenge. The actual proof verification happens
// in respondToPasswordVerifier when the client responds.
//
// Per AWS spec, USER_ID_FOR_SRP returned to the client is the user's username
// (the same value the client must use in the inner hash and claim message).
func (s *CognitoService) handleUserSrpAuth(reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	username, srpAHex, err := parseSrpAuthParams(req)
	if err != nil {
		return nil, err
	}

	clientID := getClientId(req)
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
		// Avoid revealing whether the username exists; mirror the NotAuthorized
		// semantics used elsewhere in the auth flow.
		return nil, ErrNotAuthorized
	}

	if !user.Enabled {
		s.recordAuthEvent(reqCtx, userPool.ID, user.ID, "SignIn", "Fail")
		return nil, ErrNotAuthorized
	}
	if user.UserStatus != "CONFIRMED" {
		// Users in FORCE_CHANGE_PASSWORD or other transitional states cannot
		// complete SRP; they must first complete the NEW_PASSWORD_REQUIRED
		// challenge via the non-SRP admin flow.
		s.recordAuthEvent(reqCtx, userPool.ID, user.ID, "SignIn", "Fail")
		return nil, ErrUserNotConfirmed
	}
	if user.SrpVerifier == "" || user.SrpSalt == "" {
		// The user was created before SRP support was added. They must reset
		// their password via ForgotPassword/AdminSetUserPassword to obtain a
		// verifier before USER_SRP_AUTH will succeed.
		s.recordAuthEvent(reqCtx, userPool.ID, user.ID, "SignIn", "Fail")
		return nil, ErrNotAuthorized
	}

	verifier, ok := new(big.Int).SetString(user.SrpVerifier, 16)
	if !ok {
		return nil, ErrInternalError
	}
	B, b, err := GenerateB(verifier)
	if err != nil {
		return nil, ErrInternalError
	}

	secretBlock := make([]byte, 16)
	if _, err := rand.Read(secretBlock); err != nil {
		return nil, ErrInternalError
	}

	sessionID := generateSessionID()
	challengeSession := &cognitostore.ChallengeSession{
		SessionID:     sessionID,
		UserPoolID:    userPool.ID,
		ClientID:      clientID,
		Username:      user.Username,
		ChallengeName: "PASSWORD_VERIFIER",
		CreatedAt:     time.Now().UTC(),
		ExpiresAt:     time.Now().UTC().Add(5 * time.Minute),
		SrpA:          srpAHex,
		SrpB:          B.Text(16),
		SrpPrivateB:   b.Text(16),
		SecretBlock:   base64.StdEncoding.EncodeToString(secretBlock),
	}
	if err := store.SaveChallengeSession(challengeSession); err != nil {
		return nil, ErrInternalError
	}

	s.recordAuthEvent(reqCtx, userPool.ID, user.ID, "SignIn", "InProgress")

	return map[string]interface{}{
		"ChallengeName": "PASSWORD_VERIFIER",
		"Session":       sessionID,
		"ChallengeParameters": map[string]interface{}{
			"USERNAME":        user.Username,
			"USER_ID_FOR_SRP": user.Username,
			"SALT":            user.SrpSalt,
			"SECRET_BLOCK":    challengeSession.SecretBlock,
			"SRP_B":           B.Text(16),
		},
	}, nil
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

// RespondToAuthChallenge responds to an authentication challenge.
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_RespondToAuthChallenge.html
func (s *CognitoService) RespondToAuthChallenge(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	clientID := getClientId(req)
	challengeName := req.GetParam("ChallengeName")
	session := req.GetParam("Session")

	if clientID == "" || challengeName == "" {
		return nil, ErrInvalidParameter
	}
	if !validateChallengeName(challengeName) {
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
	if challengeName == "PASSWORD_VERIFIER" {
		return s.respondToPasswordVerifier(reqCtx, req, userPool.ID, clientID, session)
	}

	return s.respondToMfaOrCustomChallenge(ctx, reqCtx, req, challengeName, userPool, clientID, session)
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

// parseSrpAuthParams extracts USERNAME and SRP_A from the AuthParameters
// block of an InitiateAuth USER_SRP_AUTH request. SRP_A is a lowercase hex
// string supplied by the client.
func parseSrpAuthParams(req *request.ParsedRequest) (username, srpAHex string, err error) {
	authParams := req.Parameters["AuthParameters"]
	if authParams == nil {
		return "", "", ErrInvalidParameter
	}
	params, ok := authParams.(map[string]interface{})
	if !ok {
		return "", "", ErrInvalidParameter
	}
	username, _ = params["USERNAME"].(string)
	srpAHex, _ = params["SRP_A"].(string)
	if username == "" || srpAHex == "" {
		return "", "", ErrInvalidParameter
	}
	if _, ok := new(big.Int).SetString(srpAHex, 16); !ok {
		return "", "", ErrInvalidParameter
	}
	return username, srpAHex, nil
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
	user.ConfirmationCodeExpiry = time.Now().Add(24 * time.Hour)
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
	if !validateChallengeName(challengeName) {
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
	if challengeName == "PASSWORD_VERIFIER" {
		return s.respondToPasswordVerifier(reqCtx, req, userPoolID, clientID, session)
	}

	userPool, err := store.GetUserPool(userPoolID)
	if err != nil {
		return nil, ErrResourceNotFound
	}
	return s.respondToMfaOrCustomChallenge(ctx, reqCtx, req, challengeName, userPool, clientID, session)
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
	saltHex, verifierHex, verr := computeSrpVerifier(userPoolID, user.Username, newPassword)
	if verr != nil {
		return nil, ErrInternalError
	}
	user.SrpSalt = saltHex
	user.SrpVerifier = verifierHex
	user.UserStatus = "CONFIRMED"

	if err := store.UpdateUser(user); err != nil {
		return nil, ErrInternalError
	}

	accessToken, idToken, refreshToken, expiresIn, err := s.CreateTokens(reqCtx, userPoolID, user.ID, clientID, TokenGenerationAuthentication, parseClientMetadata(req))
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

// respondToPasswordVerifier handles the PASSWORD_VERIFIER challenge for both
// RespondToAuthChallenge and AdminRespondToAuthChallenge. It retrieves the
// SRP session established by handleUserSrpAuth, recomputes the shared secret
// from the stored private scalar b and the user's verifier, derives the
// HMAC key, and constant-time compares the expected claim against the client's
// PASSWORD_CLAIM_SIGNATURE.
//
// AWS Cognito returns a generic NotAuthorizedException on any mismatch
// (wrong password, tampered SRP_A, expired session, etc.) to avoid leaking
// which leg of the verification failed.
func (s *CognitoService) respondToPasswordVerifier(reqCtx *request.RequestContext, req *request.ParsedRequest, userPoolID, clientID, session string) (interface{}, error) {
	challengeResponses := req.Parameters["ChallengeResponses"]
	if challengeResponses == nil {
		return nil, ErrInvalidParameter
	}
	params, ok := challengeResponses.(map[string]interface{})
	if !ok {
		return nil, ErrInvalidParameter
	}
	username, _ := params["USERNAME"].(string)
	claimSigB64, _ := params["PASSWORD_CLAIM_SIGNATURE"].(string)
	claimBlockB64, _ := params["PASSWORD_CLAIM_SECRET_BLOCK"].(string)
	timestamp, _ := params["TIMESTAMP"].(string)
	if username == "" || claimSigB64 == "" || claimBlockB64 == "" || timestamp == "" {
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
	// Burn the session regardless of outcome so a leaked Session identifier
	// cannot be replayed.
	defer func() { _ = store.DeleteChallengeSession(session) }()

	if challengeSession.ChallengeName != "PASSWORD_VERIFIER" {
		return nil, ErrNotAuthorized
	}
	if challengeSession.Username != username {
		return nil, ErrNotAuthorized
	}
	if !challengeSession.ExpiresAt.IsZero() && time.Now().UTC().After(challengeSession.ExpiresAt) {
		return nil, ErrNotAuthorized
	}
	if challengeSession.SecretBlock != claimBlockB64 {
		// The client must echo back the exact SECRET_BLOCK we issued. A
		// mismatch indicates the client did not use our challenge parameters.
		return nil, ErrNotAuthorized
	}

	A, ok := new(big.Int).SetString(challengeSession.SrpA, 16)
	if !ok {
		return nil, ErrInternalError
	}
	B, ok := new(big.Int).SetString(challengeSession.SrpB, 16)
	if !ok {
		return nil, ErrInternalError
	}
	b, ok := new(big.Int).SetString(challengeSession.SrpPrivateB, 16)
	if !ok {
		return nil, ErrInternalError
	}
	secretBlock, err := base64.StdEncoding.DecodeString(challengeSession.SecretBlock)
	if err != nil {
		return nil, ErrInternalError
	}
	clientSig, err := base64.StdEncoding.DecodeString(claimSigB64)
	if err != nil {
		return nil, ErrNotAuthorized
	}

	user, err := store.GetUser(userPoolID, username)
	if err != nil {
		return nil, ErrNotAuthorized
	}
	verifier, ok := new(big.Int).SetString(user.SrpVerifier, 16)
	if !ok || verifier.Sign() == 0 {
		return nil, ErrNotAuthorized
	}

	K, err := DeriveServerKey(A, B, b, verifier)
	if err != nil {
		// ErrInvalidSrpA indicates a malicious or malformed client value.
		return nil, ErrNotAuthorized
	}

	poolName, ok := poolNameFromID(userPoolID)
	if !ok {
		return nil, ErrInternalError
	}
	// USER_ID_FOR_SRP for Cognito is the username (not the sub).
	expectedSig := VerifyClaim(K, poolName, user.Username, secretBlock, timestamp)

	if subtle.ConstantTimeCompare(clientSig, expectedSig) != 1 {
		s.recordAuthEvent(reqCtx, userPoolID, user.ID, "SignIn", "Fail")
		return nil, ErrNotAuthorized
	}

	if !user.Enabled {
		s.recordAuthEvent(reqCtx, userPoolID, user.ID, "SignIn", "Fail")
		return nil, ErrNotAuthorized
	}

	accessToken, idToken, refreshToken, expiresIn, err := s.CreateTokens(reqCtx, userPoolID, user.ID, clientID, TokenGenerationAuthentication, parseClientMetadata(req))
	if err != nil {
		return nil, fmt.Errorf("failed to create tokens: %w", err)
	}
	s.recordAuthEvent(reqCtx, userPoolID, user.ID, "SignIn", "Pass")

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

// respondToMfaOrCustomChallenge handles all challenge types beyond
// NEW_PASSWORD_REQUIRED and PASSWORD_VERIFIER. It covers MFA challenges
// (SOFTWARE_TOKEN_MFA, SMS_MFA, SMS_OTP, EMAIL_OTP), the PASSWORD
// challenge, SELECT_MFA_TYPE, SELECT_CHALLENGE, MFA_SETUP, and
// CUSTOM_CHALLENGE (via VerifyAuthChallengeResponse Lambda trigger).
// On success it issues tokens and returns an AuthenticationResult.
func (s *CognitoService) respondToMfaOrCustomChallenge(
	ctx context.Context,
	reqCtx *request.RequestContext,
	req *request.ParsedRequest,
	challengeName string,
	userPool *cognitostore.UserPool,
	clientID, session string,
) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	// Resolve the user from the session or AuthParameters.
	username := req.GetParam("USERNAME")
	if username == "" {
		if session != "" {
			cs, err := store.GetChallengeSession(session)
			if err != nil {
				return nil, ErrNotAuthorized
			}
			username = cs.Username
		}
	}
	if username == "" {
		return nil, ErrInvalidParameter
	}

	user, err := store.GetUser(userPool.ID, username)
	if err != nil {
		return nil, ErrNotAuthorized
	}

	switch challengeName {
	case "SOFTWARE_TOKEN_MFA":
		code := req.GetParam("SOFTWARE_TOKEN_MFA_CODE")
		if code == "" {
			code = getStringParam(req.Parameters, "SoftwareTokenMfaCode")
		}
		if code == "" {
			return nil, ErrInvalidParameter
		}
		if user.SoftwareTokenMfa == nil || !user.SoftwareTokenMfa.Verified {
			return nil, ErrNotAuthorized
		}
		if !validateTOTPCode(user.SoftwareTokenMfa.SecretKey, code) {
			return nil, ErrCodeMismatch
		}

	case "SMS_MFA", "SMS_OTP":
		code := req.GetParam("SMS_MFA_CODE")
		if code == "" {
			code = getStringParam(req.Parameters, "SmsMfaCode")
		}
		if code == "" {
			return nil, ErrInvalidParameter
		}
		if user.ConfirmationCode == "" || user.ConfirmationCode != code {
			return nil, ErrCodeMismatch
		}
		user.ConfirmationCode = ""

	case "EMAIL_OTP":
		code := req.GetParam("EMAIL_OTP_CODE")
		if code == "" {
			code = getStringParam(req.Parameters, "EmailOtpCode")
		}
		if code == "" {
			return nil, ErrInvalidParameter
		}
		if user.ConfirmationCode == "" || user.ConfirmationCode != code {
			return nil, ErrCodeMismatch
		}
		user.ConfirmationCode = ""

	case "PASSWORD":
		password := req.GetParam("PASSWORD")
		if password == "" {
			password = getStringParam(req.Parameters, "Password")
		}
		if password == "" {
			return nil, ErrInvalidParameter
		}
		if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
			return nil, ErrNotAuthorized
		}

	case "SELECT_MFA_TYPE":
		mfaType := req.GetParam("MFA_TYPE")
		if mfaType == "" {
			mfaType = getStringParam(req.Parameters, "MfaType")
		}
		if mfaType == "" {
			return nil, ErrInvalidParameter
		}
		// Issue a new challenge of the selected type.
		return map[string]interface{}{
			"ChallengeName": mfaType,
			"Session":       session,
			"ChallengeParameters": map[string]string{
				"USERNAME": username,
			},
		}, nil

	case "SELECT_CHALLENGE":
		selected := req.GetParam("SELECTED_CHALLENGE")
		if selected == "" {
			selected = getStringParam(req.Parameters, "SelectedChallenge")
		}
		if selected == "" {
			return nil, ErrInvalidParameter
		}
		return map[string]interface{}{
			"ChallengeName": selected,
			"Session":       session,
			"ChallengeParameters": map[string]string{
				"USERNAME": username,
			},
		}, nil

	case "MFA_SETUP":
		// The client should enrol an MFA factor via AssociateSoftwareToken
		// or VerifySoftwareToken, then retry authentication.
		return map[string]interface{}{
			"ChallengeName": "MFA_SETUP",
			"Session":       session,
			"ChallengeParameters": map[string]string{
				"USERNAME": username,
			},
		}, nil

	case "CUSTOM_CHALLENGE":
		// Invoke VerifyAuthChallengeResponse Lambda trigger to verify the
		// client's answer.
		if userPool.LambdaConfig != nil && userPool.LambdaConfig.VerifyAuthChallengeResponse != "" {
			answer := req.GetParam("ANSWER")
			if answer == "" {
				answer = getStringParam(req.Parameters, "Answer")
			}
			result, _ := s.invokeTrigger(ctx, VerifyAuthChallengeResponse, userPool.ID, username, clientID,
				userPool.LambdaConfig.VerifyAuthChallengeResponse,
				map[string]interface{}{
					"userAttributes":  userAttributesMap(user),
					"challengeAnswer": answer,
					"clientMetadata":  parseClientMetadata(req),
				},
				map[string]interface{}{
					"answerCorrect": false,
				},
				true,
			)
			if result == nil {
				return nil, ErrNotAuthorized
			}
			if correct, _ := result["answerCorrect"].(bool); !correct {
				return nil, ErrNotAuthorized
			}
		} else {
			return nil, ErrNotAuthorized
		}

	default:
		// DEVICE_SRP_AUTH, DEVICE_PASSWORD_VERIFIER, WEB_AUTHN, PASSWORD_SRP
		// require protocol-specific verification not yet implemented.
		return nil, ErrInvalidParameter
	}

	// Challenge passed — persist any user state changes and issue tokens.
	if err := store.UpdateUser(user); err != nil {
		return nil, ErrInternalError
	}

	accessToken, idToken, _, expiresIn, err := s.CreateTokens(reqCtx, userPool.ID, user.ID, clientID, TokenGenerationAuthentication, parseClientMetadata(req))
	if err != nil {
		return nil, ErrInternalError
	}

	return map[string]interface{}{
		"AuthenticationResult": map[string]interface{}{
			"AccessToken": accessToken,
			"IdToken":     idToken,
			"TokenType":   "Bearer",
			"ExpiresIn":   expiresIn,
		},
	}, nil
}
