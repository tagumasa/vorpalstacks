package cognitoidentityprovider

import (
	"context"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"math/big"
	"time"

	"vorpalstacks/internal/common/request"
	cognitostore "vorpalstacks/internal/store/aws/cognitoidentityprovider"

	"golang.org/x/crypto/bcrypt"
)

// RespondToAuthChallengeInput carries the wire parameters of
// RespondToAuthChallenge. Params holds the raw request parameter map; the
// nested ChallengeResponses block and the per-challenge code members are
// read from it inside the Core.
type RespondToAuthChallengeInput struct {
	ClientID       string
	ChallengeName  string
	Session        string
	Params         map[string]interface{}
	ClientMetadata map[string]string
}

// AdminRespondToAuthChallengeInput carries the wire parameters of
// AdminRespondToAuthChallenge.
type AdminRespondToAuthChallengeInput struct {
	UserPoolID     string
	ClientID       string
	ChallengeName  string
	Session        string
	Params         map[string]interface{}
	ClientMetadata map[string]string
}

// respondToAuthChallengeCore validates the challenge name, resolves the
// user pool by client ID and dispatches to the challenge-specific
// implementation.
func (s *CognitoService) respondToAuthChallengeCore(ctx context.Context, reqCtx *request.RequestContext, in RespondToAuthChallengeInput) (interface{}, error) {
	if in.ClientID == "" || in.ChallengeName == "" {
		return nil, ErrInvalidParameter
	}
	if !validateChallengeName(in.ChallengeName) {
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

	if in.ChallengeName == "NEW_PASSWORD_REQUIRED" {
		return s.respondToNewPasswordChallenge(reqCtx, userPool.ID, in.ClientID, in.Session, in.Params, in.ClientMetadata)
	}
	if in.ChallengeName == "PASSWORD_VERIFIER" {
		return s.respondToPasswordVerifierCore(reqCtx, userPool.ID, in.ClientID, in.Session, in.Params, in.ClientMetadata)
	}

	return s.respondToMfaOrCustomChallenge(ctx, reqCtx, in.ChallengeName, userPool, in.ClientID, in.Session, in.Params, in.ClientMetadata)
}

// adminRespondToAuthChallengeCore validates the challenge name, checks the
// named user pool and dispatches to the challenge-specific implementation.
func (s *CognitoService) adminRespondToAuthChallengeCore(ctx context.Context, reqCtx *request.RequestContext, in AdminRespondToAuthChallengeInput) (interface{}, error) {
	if in.UserPoolID == "" || in.ClientID == "" || in.ChallengeName == "" {
		return nil, ErrInvalidParameter
	}
	if !validateChallengeName(in.ChallengeName) {
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

	if in.ChallengeName == "NEW_PASSWORD_REQUIRED" {
		return s.respondToNewPasswordChallenge(reqCtx, in.UserPoolID, in.ClientID, in.Session, in.Params, in.ClientMetadata)
	}
	if in.ChallengeName == "PASSWORD_VERIFIER" {
		return s.respondToPasswordVerifierCore(reqCtx, in.UserPoolID, in.ClientID, in.Session, in.Params, in.ClientMetadata)
	}

	return s.respondToMfaOrCustomChallenge(ctx, reqCtx, in.ChallengeName, userPool, in.ClientID, in.Session, in.Params, in.ClientMetadata)
}

// respondToNewPasswordChallenge handles the NEW_PASSWORD_REQUIRED challenge
// for both RespondToAuthChallenge and AdminRespondToAuthChallenge.
func (s *CognitoService) respondToNewPasswordChallenge(reqCtx *request.RequestContext, userPoolID, clientID, session string, params map[string]interface{}, clientMetadata map[string]string) (interface{}, error) {
	challengeResponses := params["ChallengeResponses"]
	if challengeResponses == nil {
		return nil, ErrInvalidParameter
	}

	respParams, ok := challengeResponses.(map[string]interface{})
	if !ok {
		return nil, ErrInvalidParameter
	}

	username, _ := respParams["USERNAME"].(string)
	newPassword, _ := respParams["NEW_PASSWORD"].(string)

	if username == "" || newPassword == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if _, err := validateChallengeSession(store, session, "NEW_PASSWORD_REQUIRED", userPoolID, clientID, username); err != nil {
		return nil, ErrNotAuthorized
	}
	// NEW_PASSWORD_REQUIRED sessions are single-use: burn the session
	// regardless of outcome so a leaked Session identifier cannot be
	// replayed, mirroring the PASSWORD_VERIFIER path.
	defer func() { _ = store.DeleteChallengeSession(session) }()

	userPool, err := store.GetUserPool(userPoolID)
	if err != nil {
		return nil, ErrResourceNotFound
	}

	user, err := store.GetUser(userPoolID, username)
	if err != nil {
		return nil, ErrUserNotFound
	}
	if !user.Enabled {
		s.recordAuthEvent(reqCtx, userPoolID, user.ID, username, clientID, "SignIn", "Fail")
		return nil, ErrNotAuthorized
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

	accessToken, idToken, refreshToken, expiresIn, err := s.CreateTokens(reqCtx, userPoolID, user.ID, clientID, TokenGenerationAuthentication, clientMetadata)
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

// challengeMember reads a challenge-answer member from a
// RespondToAuthChallenge request. The documented carrier is the
// ChallengeResponses map (per the AWS API reference, e.g. PASSWORD answers
// travel as "ChallengeResponses": {"USERNAME": ..., "PASSWORD": ...}); the
// historical top-level lookups are kept as fallbacks for direct-HTTP
// callers. Keys are tried in order.
func challengeMember(params map[string]interface{}, keys ...string) string {
	for _, k := range keys {
		if v := request.GetParamLowerFirst(params, k); v != "" {
			return v
		}
	}
	if responses, ok := params["ChallengeResponses"].(map[string]interface{}); ok {
		for _, k := range keys {
			if s, ok := responses[k].(string); ok && s != "" {
				return s
			}
		}
	}
	return ""
}

// respondToMfaOrCustomChallenge handles all challenge types beyond
// NEW_PASSWORD_REQUIRED and PASSWORD_VERIFIER. It covers MFA challenges
// (SOFTWARE_TOKEN_MFA, SMS_MFA, SMS_OTP, EMAIL_OTP), the PASSWORD
// challenge, SELECT_MFA_TYPE, SELECT_CHALLENGE, MFA_SETUP, and
// CUSTOM_CHALLENGE (via VerifyAuthChallengeResponse Lambda trigger).
// Every response is bound to a live challenge session of the matching type
// issued by InitiateAuth or a selector response; wrong answers increment the
// session's failure counter and exhaust the session after
// maxChallengeAttempts tries. On success it burns the session and issues
// tokens.
func (s *CognitoService) respondToMfaOrCustomChallenge(
	ctx context.Context,
	reqCtx *request.RequestContext,
	challengeName string,
	userPool *cognitostore.UserPool,
	clientID, session string,
	params map[string]interface{},
	clientMetadata map[string]string,
) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	// The type check prevents a session minted for one challenge (e.g.
	// PASSWORD_VERIFIER) from being replayed against another (e.g.
	// SOFTWARE_TOKEN_MFA).
	challengeSession, err := validateChallengeSession(store, session, challengeName, userPool.ID, clientID, "")
	if err != nil {
		return nil, ErrNotAuthorized
	}
	username := challengeSession.Username
	// The client may echo USERNAME; a mismatch is an authorisation failure.
	if echoed := challengeMember(params, "USERNAME"); echoed != "" && echoed != username {
		return nil, ErrNotAuthorized
	}

	user, err := store.GetUser(userPool.ID, username)
	if err != nil {
		return nil, ErrNotAuthorized
	}
	if !user.Enabled {
		return nil, ErrNotAuthorized
	}

	switch challengeName {
	case "SOFTWARE_TOKEN_MFA":
		code := challengeMember(params, "SOFTWARE_TOKEN_MFA_CODE", "SoftwareTokenMfaCode")
		if code == "" {
			return nil, ErrInvalidParameter
		}
		if user.SoftwareTokenMfa == nil || !user.SoftwareTokenMfa.Verified {
			return nil, ErrNotAuthorized
		}
		if !validateTOTPCode(user.SoftwareTokenMfa.SecretKey, code) {
			recordChallengeFailure(store, challengeSession)
			return nil, ErrCodeMismatch
		}

	case "SMS_MFA", "SMS_OTP":
		code := challengeMember(params, "SMS_MFA_CODE", "SmsMfaCode")
		if code == "" {
			return nil, ErrInvalidParameter
		}
		// Constant-time comparison, as with the TOTP verification: the
		// attempt limit mitigates but does not eliminate the timing side
		// channel of a naive string comparison.
		if user.ConfirmationCode == "" || subtle.ConstantTimeCompare([]byte(user.ConfirmationCode), []byte(code)) != 1 {
			recordChallengeFailure(store, challengeSession)
			return nil, ErrCodeMismatch
		}
		user.ConfirmationCode = ""

	case "EMAIL_OTP":
		code := challengeMember(params, "EMAIL_OTP_CODE", "EmailOtpCode")
		if code == "" {
			return nil, ErrInvalidParameter
		}
		if user.ConfirmationCode == "" || subtle.ConstantTimeCompare([]byte(user.ConfirmationCode), []byte(code)) != 1 {
			recordChallengeFailure(store, challengeSession)
			return nil, ErrCodeMismatch
		}
		user.ConfirmationCode = ""

	case "PASSWORD":
		password := challengeMember(params, "PASSWORD", "Password")
		if password == "" {
			return nil, ErrInvalidParameter
		}
		if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
			recordChallengeFailure(store, challengeSession)
			return nil, ErrNotAuthorized
		}

	case "SELECT_MFA_TYPE", "SELECT_CHALLENGE":
		// The documented selection key is ANSWER inside ChallengeResponses;
		// the legacy MFA_TYPE/SELECTED_CHALLENGE reads stay as fallbacks.
		selected := challengeMember(params, "ANSWER", "MFA_TYPE", "MfaType")
		if challengeName == "SELECT_CHALLENGE" {
			selected = challengeMember(params, "ANSWER", "SELECTED_CHALLENGE", "SelectedChallenge")
		}
		if selected == "" {
			return nil, ErrInvalidParameter
		}
		// Each selector accepts only its own answer set: SELECT_CHALLENGE
		// selects a sign-in challenge, SELECT_MFA_TYPE selects an MFA type.
		choices := selectMfaTypeChoices
		if challengeName == "SELECT_CHALLENGE" {
			choices = selectChallengeChoices
		}
		if !choices[selected] {
			return nil, ErrInvalidParameter
		}
		// Mint a fresh session of the selected type and burn the selector
		// session so every challenge type is answered by a session issued
		// for exactly that type.
		newSession, err := mintChallengeSession(store, userPool.ID, clientID, username, selected, 5*time.Minute)
		if err != nil {
			return nil, ErrInternalError
		}
		_ = store.DeleteChallengeSession(session)
		return map[string]interface{}{
			"ChallengeName": selected,
			"Session":       newSession,
			"ChallengeParameters": map[string]string{
				"USERNAME": username,
			},
		}, nil

	case "MFA_SETUP":
		// Only an MFA_SETUP-typed session reaches here (validated above).
		// The client enrols via AssociateSoftwareToken with this session,
		// then retries authentication.
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
			answer := challengeMember(params, "ANSWER", "Answer")
			result, _ := s.invokeTrigger(ctx, VerifyAuthChallengeResponse, userPool.ID, username, clientID,
				userPool.LambdaConfig.VerifyAuthChallengeResponse,
				map[string]interface{}{
					"userAttributes":  userAttributesMap(user),
					"challengeAnswer": answer,
					"clientMetadata":  clientMetadata,
				},
				map[string]interface{}{
					"answerCorrect": false,
				},
				true,
			)
			if result == nil {
				recordChallengeFailure(store, challengeSession)
				return nil, ErrNotAuthorized
			}
			if correct, _ := result["answerCorrect"].(bool); !correct {
				recordChallengeFailure(store, challengeSession)
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

	// Challenge passed — burn the session, persist any user state changes
	// and issue tokens.
	_ = store.DeleteChallengeSession(session)
	if err := store.UpdateUser(user); err != nil {
		return nil, ErrInternalError
	}

	accessToken, idToken, refreshToken, expiresIn, err := s.CreateTokens(reqCtx, userPool.ID, user.ID, clientID, TokenGenerationAuthentication, clientMetadata)
	if err != nil {
		return nil, ErrInternalError
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

// respondToPasswordVerifierCore handles the PASSWORD_VERIFIER challenge for
// both RespondToAuthChallenge and AdminRespondToAuthChallenge. It retrieves
// the SRP session established by userSrpAuthFlow, recomputes the shared
// secret from the stored private scalar b and the user's verifier, derives
// the HMAC key, and constant-time compares the expected claim against the
// client's PASSWORD_CLAIM_SIGNATURE.
//
// AWS Cognito returns a generic NotAuthorizedException on any mismatch
// (wrong password, tampered SRP_A, expired session, etc.) to avoid leaking
// which leg of the verification failed.
func (s *CognitoService) respondToPasswordVerifierCore(reqCtx *request.RequestContext, userPoolID, clientID, session string, params map[string]interface{}, clientMetadata map[string]string) (interface{}, error) {
	challengeResponses := params["ChallengeResponses"]
	if challengeResponses == nil {
		return nil, ErrInvalidParameter
	}
	respParams, ok := challengeResponses.(map[string]interface{})
	if !ok {
		return nil, ErrInvalidParameter
	}
	username, _ := respParams["USERNAME"].(string)
	claimSigB64, _ := respParams["PASSWORD_CLAIM_SIGNATURE"].(string)
	claimBlockB64, _ := respParams["PASSWORD_CLAIM_SECRET_BLOCK"].(string)
	timestamp, _ := respParams["TIMESTAMP"].(string)
	if username == "" || claimSigB64 == "" || claimBlockB64 == "" || timestamp == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	challengeSession, err := validateChallengeSession(store, session, "PASSWORD_VERIFIER", userPoolID, clientID, username)
	if err != nil {
		return nil, ErrNotAuthorized
	}
	// Burn the session regardless of outcome so a leaked Session identifier
	// cannot be replayed.
	defer func() { _ = store.DeleteChallengeSession(session) }()

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
		s.recordAuthEvent(reqCtx, userPoolID, user.ID, username, clientID, "SignIn", "Fail")
		return nil, ErrNotAuthorized
	}

	if !user.Enabled {
		s.recordAuthEvent(reqCtx, userPoolID, user.ID, username, clientID, "SignIn", "Fail")
		return nil, ErrNotAuthorized
	}

	accessToken, idToken, refreshToken, expiresIn, err := s.CreateTokens(reqCtx, userPoolID, user.ID, clientID, TokenGenerationAuthentication, clientMetadata)
	if err != nil {
		return nil, fmt.Errorf("failed to create tokens: %w", err)
	}
	s.recordAuthEvent(reqCtx, userPoolID, user.ID, username, clientID, "SignIn", "Pass")

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
