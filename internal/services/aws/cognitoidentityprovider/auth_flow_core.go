package cognitoidentityprovider

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"math/big"
	"time"

	"vorpalstacks/internal/common/request"
	cognitostore "vorpalstacks/internal/store/aws/cognitoidentityprovider"

	"golang.org/x/crypto/bcrypt"
)

// InitiateAuthInput carries the wire parameters of InitiateAuth. Params holds
// the raw request parameter map; the nested AuthParameters block is read from
// it inside the Core so the extraction semantics stay with the flow logic.
type InitiateAuthInput struct {
	AuthFlow       string
	ClientID       string
	Params         map[string]interface{}
	ValidationData map[string]string
	ClientMetadata map[string]string
}

// AdminInitiateAuthInput carries the wire parameters of AdminInitiateAuth.
type AdminInitiateAuthInput struct {
	UserPoolID     string
	AuthFlow       string
	ClientID       string
	Params         map[string]interface{}
	ValidationData map[string]string
	ClientMetadata map[string]string
}

// initiateAuthCore validates the requested flow and dispatches to the
// flow-specific implementation.
func (s *CognitoService) initiateAuthCore(ctx context.Context, reqCtx *request.RequestContext, in InitiateAuthInput) (interface{}, error) {
	if in.AuthFlow == "" || in.ClientID == "" {
		return nil, ErrInvalidParameter
	}
	if !validateInitiateAuthFlow(in.AuthFlow) {
		return nil, ErrInvalidParameter
	}

	switch in.AuthFlow {
	case "USER_PASSWORD_AUTH":
		return s.userPasswordAuthFlow(ctx, reqCtx, in)
	case "USER_SRP_AUTH":
		return s.userSrpAuthFlow(reqCtx, in.ClientID, in.Params)
	case "REFRESH_TOKEN_AUTH", "REFRESH_TOKEN":
		return s.refreshTokenAuthFlow(reqCtx, in.Params)
	case "CUSTOM_AUTH":
		return s.customAuthFlow(ctx, reqCtx, in.ClientID,
			authMember(in.Params, "Username", "USERNAME"),
			authMember(in.Params, "Password", "PASSWORD"))
	case "USER_AUTH":
		return s.userAuthFlow(reqCtx, in.ClientID, authMember(in.Params, "Username", "USERNAME"))
	default:
		return nil, ErrInvalidParameter
	}
}

// adminInitiateAuthCore validates the requested admin flow, resolves the
// named user pool and dispatches to the flow-specific implementation.
func (s *CognitoService) adminInitiateAuthCore(ctx context.Context, reqCtx *request.RequestContext, in AdminInitiateAuthInput) (interface{}, error) {
	if in.UserPoolID == "" || in.ClientID == "" || in.AuthFlow == "" {
		return nil, ErrInvalidParameter
	}
	if !validateAdminInitiateAuthFlow(in.AuthFlow) {
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

	switch in.AuthFlow {
	case "ADMIN_NO_SRP_AUTH", "ADMIN_USER_PASSWORD_AUTH":
		return s.adminNoSrpAuthFlow(ctx, reqCtx, in.UserPoolID, in.ClientID, userPool.LambdaConfig, in.Params, in.ValidationData, in.ClientMetadata)
	case "CUSTOM_AUTH":
		return s.customAuthFlow(ctx, reqCtx, in.ClientID,
			authMember(in.Params, "Username", "USERNAME"),
			authMember(in.Params, "Password", "PASSWORD"))
	case "REFRESH_TOKEN_AUTH", "REFRESH_TOKEN":
		return s.adminRefreshTokenAuthFlow(reqCtx, in.UserPoolID, in.Params)
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

// customAuthFlow implements the CUSTOM_AUTH flow. The client provides a
// USERNAME (and optionally a PASSWORD for the initial verification). The
// server invokes the DefineAuthChallenge and CreateAuthChallenge Lambda
// triggers to produce a custom challenge that the client must answer via
// RespondToAuthChallenge.
func (s *CognitoService) customAuthFlow(ctx context.Context, reqCtx *request.RequestContext, clientID, username, password string) (interface{}, error) {
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
		ExpiresAt:     time.Now().UTC().Add(challengeSessionTTL),
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

// userAuthFlow implements the USER_AUTH flow (the modern recommended auth
// flow). The server inspects the user's configured auth factors and returns
// the set of available challenges for the client to choose from.
func (s *CognitoService) userAuthFlow(reqCtx *request.RequestContext, clientID, username string) (interface{}, error) {
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
	var available []string
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
		ExpiresAt:     time.Now().UTC().Add(challengeSessionTTL),
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
		// A Lambda function that raises an error rejects the migration and
		// fails the sign-in as a wrong password; an invocation-transport
		// failure is an infrastructure error. A missing trigger means no
		// migration path exists and authentication fails as an incorrect
		// password.
		if migrationErr != nil {
			return nil, classifyMigrationFailure(migrationErr)
		}
		if migrationResult == nil {
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
		s.recordAuthEvent(reqCtx, userPoolID, user.ID, username, clientID, "SignIn", "Fail")
		return nil, ErrNotAuthorized
	}

	// The password must verify before any challenge is issued. Amazon
	// Cognito returns NEW_PASSWORD_REQUIRED only for users who signed in
	// successfully with their temporary password; issuing it on user status
	// alone would let anyone reset a FORCE_CHANGE_PASSWORD or RESET_REQUIRED
	// account by username only. Users imported from a CSV are the documented
	// exception: "the first time they sign in, they can enter any password.
	// Amazon Cognito prompts them to enter a new password" — an imported
	// RESET_REQUIRED user with no stored hash goes straight to the
	// NEW_PASSWORD_REQUIRED challenge, and a user holding an imported hash
	// verifies against the import algorithm, migrating to the native
	// bcrypt+SRP credentials on success.
	switch {
	case user.PasswordHashAlgo != "":
		if !verifyImportedPasswordHash(user.PasswordHashAlgo, user.PasswordHash, password) {
			s.recordAuthEvent(reqCtx, userPoolID, user.ID, username, clientID, "SignIn", "Fail")
			return nil, ErrIncorrectPassword
		}
		if err := s.migrateImportedCredentials(store, userPoolID, user, password); err != nil {
			return nil, err
		}
	case user.PasswordHash == "":
		if user.UserStatus != "RESET_REQUIRED" {
			s.recordAuthEvent(reqCtx, userPoolID, user.ID, username, clientID, "SignIn", "Fail")
			return nil, ErrIncorrectPassword
		}
	default:
		if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
			s.recordAuthEvent(reqCtx, userPoolID, user.ID, username, clientID, "SignIn", "Fail")
			return nil, ErrIncorrectPassword
		}
	}

	if user.UserStatus == "FORCE_CHANGE_PASSWORD" || user.UserStatus == "RESET_REQUIRED" {
		s.recordAuthEvent(reqCtx, userPoolID, user.ID, username, clientID, "SignIn", "InProgress")
		return s.newPasswordChallenge(reqCtx, userPoolID, clientID, user)
	}

	if user.UserStatus != "CONFIRMED" {
		s.recordAuthEvent(reqCtx, userPoolID, user.ID, username, clientID, "SignIn", "Fail")
		return nil, ErrUserNotConfirmed
	}

	attrs := userAttributesMap(user)
	if err := invokePreAuthentication(ctx, s, userPoolID, username, clientID, lambdaConfig, attrs, clientMetadata); err != nil {
		s.recordAuthEvent(reqCtx, userPoolID, user.ID, username, clientID, "SignIn", "Fail")
		return nil, fmt.Errorf("PreAuthentication trigger failed: %w", err)
	}

	if err := invokePostAuthentication(ctx, s, userPoolID, username, clientID, lambdaConfig, attrs, clientMetadata); err != nil {
		return nil, fmt.Errorf("PostAuthentication trigger failed: %w", err)
	}

	accessToken, idToken, refreshToken, expiresIn, err := s.CreateTokens(reqCtx, userPoolID, user.ID, clientID, TokenGenerationAuthentication, clientMetadata)
	if err != nil {
		return nil, fmt.Errorf("failed to create tokens: %w", err)
	}
	s.recordAuthEvent(reqCtx, userPoolID, user.ID, username, clientID, "SignIn", "Pass")
	return authResult(accessToken, idToken, refreshToken, expiresIn), nil
}

// migrateImportedCredentials replaces an imported password hash with the
// native bcrypt+SRP pair after successful verification, mirroring AWS's
// transparent credential migration at first sign-in.
func (s *CognitoService) migrateImportedCredentials(store cognitostore.CognitoStoreInterface, userPoolID string, user *cognitostore.User, password string) error {
	if err := setNativePasswordCredentials(user, userPoolID, user.Username, password); err != nil {
		return ErrInternalError
	}
	if err := store.UpdateUser(user); err != nil {
		return ErrInternalError
	}
	return nil
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
		ExpiresAt:     time.Now().UTC().Add(srpChallengeSessionTTL),
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

// userPasswordAuthFlow implements the InitiateAuth USER_PASSWORD_AUTH flow.
func (s *CognitoService) userPasswordAuthFlow(ctx context.Context, reqCtx *request.RequestContext, in InitiateAuthInput) (interface{}, error) {
	username, password, err := authParamsFrom(in.Params)
	if err != nil {
		return nil, err
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	userPool, err := store.GetUserPoolByClientID(in.ClientID)
	if err != nil {
		return nil, ErrResourceNotFound
	}

	return s.authenticateUser(ctx, reqCtx, userPool.ID, in.ClientID, username, password, userPool.LambdaConfig, in.ValidationData, in.ClientMetadata)
}

// adminNoSrpAuthFlow implements the AdminInitiateAuth
// ADMIN_NO_SRP_AUTH / ADMIN_USER_PASSWORD_AUTH flows.
func (s *CognitoService) adminNoSrpAuthFlow(
	ctx context.Context, reqCtx *request.RequestContext,
	userPoolID, clientID string,
	lambdaConfig *cognitostore.LambdaConfig,
	params map[string]interface{},
	validationData, clientMetadata map[string]string,
) (interface{}, error) {
	username, password, err := authParamsFrom(params)
	if err != nil {
		return nil, err
	}

	return s.authenticateUser(ctx, reqCtx, userPoolID, clientID, username, password, lambdaConfig, validationData, clientMetadata)
}

// userSrpAuthFlow handles the InitiateAuth USER_SRP_AUTH flow. It receives
// the client's SRP_A value, looks up the user, generates the server's
// ephemeral B and a fresh SECRET_BLOCK, persists the SRP session state, and
// returns a PASSWORD_VERIFIER challenge. The actual proof verification happens
// in respondToPasswordVerifierCore when the client responds.
//
// Per AWS spec, USER_ID_FOR_SRP returned to the client is the user's username
// (the same value the client must use in the inner hash and claim message).
func (s *CognitoService) userSrpAuthFlow(reqCtx *request.RequestContext, clientID string, params map[string]interface{}) (interface{}, error) {
	username, srpAHex, err := srpAuthParamsFrom(params)
	if err != nil {
		return nil, err
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
		// Avoid revealing whether the username exists; mirror the NotAuthorized
		// semantics used elsewhere in the auth flow.
		return nil, ErrNotAuthorized
	}

	if !user.Enabled {
		s.recordAuthEvent(reqCtx, userPool.ID, user.ID, username, clientID, "SignIn", "Fail")
		return nil, ErrNotAuthorized
	}
	if user.UserStatus != "CONFIRMED" {
		// Users in FORCE_CHANGE_PASSWORD or other transitional states cannot
		// complete SRP; they must first complete the NEW_PASSWORD_REQUIRED
		// challenge via the non-SRP admin flow.
		s.recordAuthEvent(reqCtx, userPool.ID, user.ID, username, clientID, "SignIn", "Fail")
		return nil, ErrUserNotConfirmed
	}
	if user.SrpVerifier == "" || user.SrpSalt == "" {
		// The user was created before SRP support was added. They must reset
		// their password via ForgotPassword/AdminSetUserPassword to obtain a
		// verifier before USER_SRP_AUTH will succeed.
		s.recordAuthEvent(reqCtx, userPool.ID, user.ID, username, clientID, "SignIn", "Fail")
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
		ExpiresAt:     time.Now().UTC().Add(srpChallengeSessionTTL),
		SrpA:          srpAHex,
		SrpB:          B.Text(16),
		SrpPrivateB:   b.Text(16),
		SecretBlock:   base64.StdEncoding.EncodeToString(secretBlock),
	}
	if err := store.SaveChallengeSession(challengeSession); err != nil {
		return nil, ErrInternalError
	}

	s.recordAuthEvent(reqCtx, userPool.ID, user.ID, username, clientID, "SignIn", "InProgress")

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

// refreshTokenAuthFlow implements the InitiateAuth REFRESH_TOKEN_AUTH flow.
func (s *CognitoService) refreshTokenAuthFlow(reqCtx *request.RequestContext, params map[string]interface{}) (interface{}, error) {
	refreshToken, err := refreshTokenFrom(params)
	if err != nil {
		return nil, err
	}
	return s.refreshAuthToken(reqCtx, "", refreshToken)
}

// adminRefreshTokenAuthFlow implements the AdminInitiateAuth
// REFRESH_TOKEN_AUTH flow with an explicit user pool binding.
func (s *CognitoService) adminRefreshTokenAuthFlow(reqCtx *request.RequestContext, userPoolID string, params map[string]interface{}) (interface{}, error) {
	refreshToken, err := refreshTokenFrom(params)
	if err != nil {
		return nil, err
	}
	return s.refreshAuthToken(reqCtx, userPoolID, refreshToken)
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

// authMember reads an InitiateAuth AuthParameters member. The documented
// carrier is the AuthParameters map (per the AWS API reference: USER_AUTH
// takes USERNAME, CUSTOM_AUTH takes USERNAME, USER_PASSWORD_AUTH takes
// USERNAME and PASSWORD); the historical top-level lookups are kept as
// fallbacks for direct-HTTP callers. Keys are tried in order.
func authMember(params map[string]interface{}, keys ...string) string {
	for _, k := range keys {
		if v := request.GetParamLowerFirst(params, k); v != "" {
			return v
		}
	}
	if authParams, ok := params["AuthParameters"].(map[string]interface{}); ok {
		for _, k := range keys {
			if s, ok := authParams[k].(string); ok && s != "" {
				return s
			}
		}
	}
	return ""
}

// authParamsFrom extracts USERNAME and PASSWORD from the AuthParameters
// block of an InitiateAuth or AdminInitiateAuth request.
func authParamsFrom(params map[string]interface{}) (username, password string, err error) {
	authParams := params["AuthParameters"]
	if authParams == nil {
		return "", "", ErrInvalidParameter
	}
	m, ok := authParams.(map[string]interface{})
	if !ok {
		return "", "", ErrInvalidParameter
	}
	username, _ = m["USERNAME"].(string)
	password, _ = m["PASSWORD"].(string)
	if username == "" || password == "" {
		return "", "", ErrInvalidParameter
	}
	return username, password, nil
}

// srpAuthParamsFrom extracts USERNAME and SRP_A from the AuthParameters
// block of an InitiateAuth USER_SRP_AUTH request. SRP_A is a lowercase hex
// string supplied by the client.
func srpAuthParamsFrom(params map[string]interface{}) (username, srpAHex string, err error) {
	authParams := params["AuthParameters"]
	if authParams == nil {
		return "", "", ErrInvalidParameter
	}
	m, ok := authParams.(map[string]interface{})
	if !ok {
		return "", "", ErrInvalidParameter
	}
	username, _ = m["USERNAME"].(string)
	srpAHex, _ = m["SRP_A"].(string)
	if username == "" || srpAHex == "" {
		return "", "", ErrInvalidParameter
	}
	if _, ok := new(big.Int).SetString(srpAHex, 16); !ok {
		return "", "", ErrInvalidParameter
	}
	return username, srpAHex, nil
}

// refreshTokenFrom extracts the REFRESH_TOKEN from AuthParameters.
func refreshTokenFrom(params map[string]interface{}) (string, error) {
	authParams := params["AuthParameters"]
	if authParams == nil {
		return "", ErrInvalidParameter
	}
	m, ok := authParams.(map[string]interface{})
	if !ok {
		return "", ErrInvalidParameter
	}
	refreshToken, _ := m["REFRESH_TOKEN"].(string)
	if refreshToken == "" {
		return "", ErrInvalidParameter
	}
	return refreshToken, nil
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
