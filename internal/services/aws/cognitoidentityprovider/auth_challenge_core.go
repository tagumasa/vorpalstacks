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

	awserrors "vorpalstacks/internal/common/errors"
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

	// "You must provide a SECRET_HASH parameter in all challenge responses
	// to an app client that has a client secret."
	if err := s.verifyAuthPlaneSecretHash(reqCtx, in.ClientID, challengeMember(in.Params, "USERNAME"), challengeMember(in.Params, "SECRET_HASH")); err != nil {
		return nil, err
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

	// The admin plane carries the same challenge-response contract: "You
	// must provide a SECRET_HASH parameter in all challenge responses to an
	// app client that has a client secret."
	if err := s.verifyAuthPlaneSecretHash(reqCtx, in.ClientID, challengeMember(in.Params, "USERNAME"), challengeMember(in.Params, "SECRET_HASH")); err != nil {
		return nil, err
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
		s.recordAuthEvent(reqCtx, userPoolID, user.ID, username, clientID, authEventSignIn, authEventResponseFail)
		return nil, ErrNotAuthorized
	}

	if err := validatePassword(newPassword, userPool.PasswordPolicy); err != nil {
		return nil, ErrPasswordPolicyViolation
	}

	// The credential rewrite and the CONFIRMED transition run as one
	// serialised read-modify-write against the fresh record, so a concurrent
	// password writer's history cannot be overwritten by the snapshot this
	// request read.
	err = store.UpdateUserFunc(userPoolID, username, func(u *cognitostore.User) error {
		if err := setNativePasswordCredentials(u, userPool.PasswordPolicy, newPassword); err != nil {
			return err
		}
		u.UserStatus = "CONFIRMED"
		user = u
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

	// Completing NEW_PASSWORD_REQUIRED finishes the primary credentials:
	// the pool's MFA configuration applies before tokens are issued. The
	// new-password session is burned (deferred above); an MFA challenge
	// mints its own session. A remembered device named by DEVICE_KEY
	// replaces the MFA challenge with device SRP authentication.
	challenge, err := s.nextChallengeAfterPrimary(reqCtx, store, userPool, clientID, user, challengeMember(params, "DEVICE_KEY"))
	if err != nil {
		return nil, ErrInternalError
	}
	if challenge != nil {
		s.recordAuthEvent(reqCtx, userPoolID, user.ID, username, clientID, authEventSignIn, authEventResponseInProgress)
		return challenge, nil
	}

	accessToken, idToken, refreshToken, expiresIn, err := s.CreateTokens(reqCtx, userPoolID, user.ID, clientID, TokenGenerationAuthentication, clientMetadata)
	if err != nil {
		return nil, fmt.Errorf("failed to create tokens: %w", err)
	}

	return s.withNewDeviceMetadata(store, userPool, user, challengeMember(params, "DEVICE_KEY"), map[string]interface{}{
		"AuthenticationResult": map[string]interface{}{
			"AccessToken":  accessToken,
			"IdToken":      idToken,
			"RefreshToken": refreshToken,
			"TokenType":    "Bearer",
			"ExpiresIn":    expiresIn,
		},
	}), nil
}

// challengeMember reads a challenge-answer member from a
// RespondToAuthChallenge request. The documented carrier is the
// ChallengeResponses map (per the AWS API reference, e.g. PASSWORD answers
// travel as "ChallengeResponses": {"USERNAME": ..., "PASSWORD": ...});
// every key passed here is one of the documented response names, and
// nothing outside the map is read.
func challengeMember(params map[string]interface{}, keys ...string) string {
	responses, ok := params["ChallengeResponses"].(map[string]interface{})
	if !ok {
		return ""
	}
	for _, k := range keys {
		if s, ok := responses[k].(string); ok && s != "" {
			return s
		}
	}
	return ""
}

// respondToMfaOrCustomChallenge handles all challenge types beyond
// NEW_PASSWORD_REQUIRED and PASSWORD_VERIFIER. It covers MFA challenges
// (SOFTWARE_TOKEN_MFA, SMS_MFA, EMAIL_OTP, the WEB_AUTHN second factor), the
// primary sign-in challenges (PASSWORD, PASSWORD_SRP, WEB_AUTHN), the
// selector challenges (SELECT_MFA_TYPE, SELECT_CHALLENGE, MFA_SETUP), the
// device challenges (DEVICE_SRP_AUTH, DEVICE_PASSWORD_VERIFIER) and
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

	// completesPrimary tracks whether a successful answer completes the
	// user's primary credentials (the pool's MFA configuration applies
	// before tokens are issued); an MFA-typed answer completes the second
	// factor itself and goes straight to tokens.
	completesPrimary := challengeName == "PASSWORD"

	switch challengeName {
	case "SOFTWARE_TOKEN_MFA":
		code := challengeMember(params, "SOFTWARE_TOKEN_MFA_CODE")
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

	case "SMS_MFA", "SMS_OTP", "EMAIL_OTP":
		// Each OTP-typed challenge is answered with its own documented key:
		// SMS_MFA reads SMS_MFA_CODE, SMS_OTP reads SMS_OTP_CODE and
		// EMAIL_OTP reads EMAIL_OTP_CODE.
		var codeKey string
		switch challengeName {
		case "SMS_MFA":
			codeKey = "SMS_MFA_CODE"
		case "SMS_OTP":
			codeKey = "SMS_OTP_CODE"
		default:
			codeKey = "EMAIL_OTP_CODE"
		}
		code := challengeMember(params, codeKey)
		if code == "" {
			return nil, ErrInvalidParameter
		}
		// The answer is compared against the one-time code bound to this
		// challenge session (generated at mint, dying with the session): a
		// code issued for any other purpose — sign-up confirmation or a
		// password reset — never satisfies an OTP challenge. Constant-time
		// comparison, as with the TOTP verification: the attempt limit
		// mitigates but does not eliminate the timing side channel of a
		// naive string comparison.
		if challengeSession.OTPCode == "" || subtle.ConstantTimeCompare([]byte(challengeSession.OTPCode), []byte(code)) != 1 {
			if !testMode {
				recordChallengeFailure(store, challengeSession)
				return nil, ErrCodeMismatch
			}
			// TEST_MODE: any non-empty code answers, as with ConfirmSignUp.
		}

	case "PASSWORD":
		password := challengeMember(params, "PASSWORD")
		if password == "" {
			return nil, ErrInvalidParameter
		}
		if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
			recordChallengeFailure(store, challengeSession)
			return nil, ErrNotAuthorized
		}

	case "PASSWORD_SRP":
		// The initial SRP leg: the client answers with SRP_A and the server
		// responds with the PASSWORD_VERIFIER challenge the client computes
		// its claim against.
		srpAHex := challengeMember(params, "SRP_A")
		if srpAHex == "" {
			return nil, ErrInvalidParameter
		}
		if _, ok := new(big.Int).SetString(srpAHex, 16); !ok {
			return nil, ErrInvalidParameter
		}
		resp, serr := s.issueSrpPasswordVerifier(reqCtx, store, userPool, clientID, user, srpAHex)
		if serr != nil {
			return nil, serr
		}
		_ = store.DeleteChallengeSession(session)
		return resp, nil

	case "WEB_AUTHN":
		credentialJSON := challengeMember(params, "CREDENTIAL")
		if credentialJSON == "" {
			return nil, ErrInvalidParameter
		}
		if verr := verifyWebAuthnAssertion(store, challengeSession, userPool, user, credentialJSON); verr != nil {
			recordChallengeFailure(store, challengeSession)
			s.recordAuthEvent(reqCtx, userPool.ID, user.ID, username, clientID, authEventSignIn, authEventResponseFail)
			if verr == ErrInvalidParameter {
				return nil, ErrInvalidParameter
			}
			return nil, ErrNotAuthorized
		}
		// A client-selected passkey is the primary sign-in method; a
		// session the MFA machinery issued carries the second-factor role
		// and the assertion completes the sign-in.
		completesPrimary = challengeSession.ChallengeRole != challengeRoleMFA

	case "SELECT_MFA_TYPE", "SELECT_CHALLENGE":
		// The documented selection key is ANSWER inside ChallengeResponses;
		// each selector then admits only its own answer set below.
		selected := challengeMember(params, "ANSWER")
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

		if challengeName == "SELECT_CHALLENGE" {
			// PASSWORD, PASSWORD_SRP and WEB_AUTHN complete authentication
			// in the selection response itself: the selected challenge's
			// answer travels alongside ANSWER, so no intermediate challenge
			// round-trip is needed. The OTP selections deliver their code
			// with the next challenge and are answered separately.
			inlineAnswered := false
			switch selected {
			case "PASSWORD":
				if password := challengeMember(params, "PASSWORD"); password != "" {
					if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
						recordChallengeFailure(store, challengeSession)
						return nil, ErrNotAuthorized
					}
					completesPrimary = true
					inlineAnswered = true
				}
			case "PASSWORD_SRP":
				if srpAHex := challengeMember(params, "SRP_A"); srpAHex != "" {
					if _, ok := new(big.Int).SetString(srpAHex, 16); !ok {
						return nil, ErrInvalidParameter
					}
					resp, serr := s.issueSrpPasswordVerifier(reqCtx, store, userPool, clientID, user, srpAHex)
					if serr != nil {
						return nil, serr
					}
					_ = store.DeleteChallengeSession(session)
					return resp, nil
				}
			case "WEB_AUTHN":
				if credentialJSON := challengeMember(params, "CREDENTIAL"); credentialJSON != "" {
					// The selector session carries the assertion request
					// options issued with the AvailableChallenges response;
					// the assertion verifies against them.
					if verr := verifyWebAuthnAssertion(store, challengeSession, userPool, user, credentialJSON); verr != nil {
						recordChallengeFailure(store, challengeSession)
						s.recordAuthEvent(reqCtx, userPool.ID, user.ID, username, clientID, authEventSignIn, authEventResponseFail)
						if verr == ErrInvalidParameter {
							return nil, ErrInvalidParameter
						}
						return nil, ErrNotAuthorized
					}
					completesPrimary = true
					inlineAnswered = true
				}
			}
			if inlineAnswered {
				// The selection answer completed (or began, for SRP) the
				// authentication; an inline PASSWORD or WEB_AUTHN answer
				// continues at the shared completion tail below.
				break
			}
			// A passkey selection without an inline CREDENTIAL answer issues
			// the WEB_AUTHN challenge with fresh assertion options: a bare
			// WEB_AUTHN session carries no challenge bytes and no relying
			// party id, so it could never be answered.
			if selected == "WEB_AUTHN" {
				resp, werr := s.issueWebAuthnChallenge(store, userPool, clientID, user, "")
				if werr != nil {
					return nil, werr
				}
				_ = store.DeleteChallengeSession(session)
				return resp, nil
			}
		}

		// Mint a fresh session of the selected type and burn the selector
		// session so every challenge type is answered by a session issued
		// for exactly that type. The SELECT_MFA_TYPE answer EMAIL_MFA names
		// the email factor; the challenge it selects is EMAIL_OTP — the
		// challenge-name vocabulary carries no EMAIL_MFA member.
		selectedChallenge := selected
		if selected == "EMAIL_MFA" {
			selectedChallenge = "EMAIL_OTP"
		}
		newSession, err := mintChallengeSession(store, userPool.ID, clientID, username, selectedChallenge, challengeSessionTTL)
		if err != nil {
			return nil, ErrInternalError
		}
		s.deliverChallengeCode(reqCtx, userPool.ID, selectedChallenge)
		_ = store.DeleteChallengeSession(session)
		return map[string]interface{}{
			"ChallengeName": selectedChallenge,
			"Session":       newSession,
			"ChallengeParameters": map[string]string{
				"USERNAME": username,
			},
		}, nil

	case "MFA_SETUP":
		// Only an MFA_SETUP-typed session reaches here (validated above).
		// The session VerifySoftwareToken returns "satisfies an MFA_SETUP
		// challenge": once the user holds a verified software token,
		// presenting that session completes the sign-in at the shared
		// completion tail below. Before enrolment the challenge shell is
		// re-issued for the client to enrol through AssociateSoftwareToken
		// and VerifySoftwareToken first.
		if user.SoftwareTokenMfa != nil && user.SoftwareTokenMfa.Enabled && user.SoftwareTokenMfa.Verified {
			break
		}
		return map[string]interface{}{
			"ChallengeName": "MFA_SETUP",
			"Session":       session,
			"ChallengeParameters": map[string]string{
				"USERNAME": username,
			},
		}, nil

	case "DEVICE_SRP_AUTH":
		// The initial leg of device SRP: the client names the device and
		// sends SRP_A; the server generates its ephemeral against the
		// device's stored verifier and responds with the
		// DEVICE_PASSWORD_VERIFIER challenge. The device replaces the
		// sign-in's MFA challenge, so the proof completes the sign-in.
		deviceKey := challengeMember(params, "DEVICE_KEY")
		srpAHex := challengeMember(params, "SRP_A")
		if deviceKey == "" || srpAHex == "" {
			return nil, ErrInvalidParameter
		}
		if _, ok := new(big.Int).SetString(srpAHex, 16); !ok {
			return nil, ErrInvalidParameter
		}
		if challengeSession.DeviceKey != "" && challengeSession.DeviceKey != deviceKey {
			return nil, ErrNotAuthorized
		}
		device, derr := store.GetDevice(userPool.ID, user.ID, deviceKey)
		if derr != nil || device.DeviceRememberedStatus != "remembered" || device.DeviceSecretVerifierB == "" {
			recordChallengeFailure(store, challengeSession)
			return nil, ErrNotAuthorized
		}
		verifier, ok := new(big.Int).SetString(device.DeviceSecretVerifierB, 16)
		if !ok {
			return nil, ErrInternalError
		}
		B, b, gerr := GenerateB(verifier)
		if gerr != nil {
			return nil, ErrInternalError
		}
		secretBlock := make([]byte, 16)
		if _, rerr := rand.Read(secretBlock); rerr != nil {
			return nil, ErrInternalError
		}
		verifierSession := newChallengeSession(userPool.ID, clientID, username, "DEVICE_PASSWORD_VERIFIER", srpChallengeSessionTTL)
		verifierSession.SrpA = srpAHex
		verifierSession.SrpB = B.Text(16)
		verifierSession.SrpPrivateB = b.Text(16)
		verifierSession.SecretBlock = base64.StdEncoding.EncodeToString(secretBlock)
		verifierSession.DeviceKey = deviceKey
		if serr := store.SaveChallengeSession(verifierSession); serr != nil {
			return nil, ErrInternalError
		}
		_ = store.DeleteChallengeSession(session)
		s.recordAuthEvent(reqCtx, userPool.ID, user.ID, username, clientID, authEventSignIn, authEventResponseInProgress)
		return map[string]interface{}{
			"ChallengeName": "DEVICE_PASSWORD_VERIFIER",
			"Session":       verifierSession.SessionID,
			"ChallengeParameters": map[string]interface{}{
				"USERNAME":     username,
				"DEVICE_KEY":   deviceKey,
				"SECRET_BLOCK": verifierSession.SecretBlock,
				"SRP_B":        B.Text(16),
			},
		}, nil

	case "DEVICE_PASSWORD_VERIFIER":
		// The device SRP proof. The claim message hashes the device group
		// key and the device key in place of the pool name and user id of
		// the password claim; the shared secret derives from the device's
		// stored verifier.
		deviceKey := challengeMember(params, "DEVICE_KEY")
		claimSigB64 := challengeMember(params, "PASSWORD_CLAIM_SIGNATURE")
		claimBlockB64 := challengeMember(params, "PASSWORD_CLAIM_SECRET_BLOCK")
		timestamp := challengeMember(params, "TIMESTAMP")
		if deviceKey == "" || claimSigB64 == "" || claimBlockB64 == "" || timestamp == "" {
			return nil, ErrInvalidParameter
		}
		if challengeSession.DeviceKey != deviceKey {
			return nil, ErrNotAuthorized
		}
		if challengeSession.SecretBlock != claimBlockB64 {
			return nil, ErrNotAuthorized
		}
		device, derr := store.GetDevice(userPool.ID, user.ID, deviceKey)
		if derr != nil || device.DeviceSecretVerifierB == "" {
			return nil, ErrNotAuthorized
		}
		if _, verr := s.verifyDeviceSrpClaim(store, challengeSession, userPool, device, claimSigB64, timestamp); verr != nil {
			recordChallengeFailure(store, challengeSession)
			s.recordAuthEvent(reqCtx, userPool.ID, user.ID, username, clientID, authEventSignIn, authEventResponseFail)
			return nil, ErrNotAuthorized
		}
		device.DeviceLastAuthenticatedDate = time.Now().UTC()
		if uerr := store.UpdateDevice(device); uerr != nil {
			return nil, ErrInternalError
		}
		// The device proof replaces the MFA challenge — the sign-in
		// completes with tokens.

	case "CUSTOM_CHALLENGE":
		// Invoke VerifyAuthChallengeResponse Lambda trigger to verify the
		// client's answer.
		if userPool.LambdaConfig != nil && userPool.LambdaConfig.VerifyAuthChallengeResponse != "" {
			answer := challengeMember(params, "ANSWER")
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

	case "ADMIN_NO_SRP_AUTH":
		// The model documents that an ADMIN_NO_SRP_AUTH challenge cannot be
		// responded to with this operation.
		return nil, ErrInvalidParameter

	default:
		return nil, ErrInvalidParameter
	}

	// Challenge passed — burn the session, persist any user state changes
	// and issue tokens.
	_ = store.DeleteChallengeSession(session)
	if err := store.UpdateUser(user); err != nil {
		return nil, ErrInternalError
	}

	// A challenge that completed the primary credentials leaves the pool's
	// MFA configuration to apply before any tokens are issued; an MFA-typed
	// challenge answer completes the second factor itself and goes straight
	// to tokens. A remembered device named by DEVICE_KEY replaces the MFA
	// challenge with device SRP authentication.
	if completesPrimary {
		deviceKey := challengeMember(params, "DEVICE_KEY")
		challenge, cerr := s.nextChallengeAfterPrimary(reqCtx, store, userPool, clientID, user, deviceKey)
		if cerr != nil {
			return nil, ErrInternalError
		}
		if challenge != nil {
			s.recordAuthEvent(reqCtx, userPool.ID, user.ID, username, clientID, authEventSignIn, authEventResponseInProgress)
			return challenge, nil
		}
	}

	accessToken, idToken, refreshToken, expiresIn, err := s.CreateTokens(reqCtx, userPool.ID, user.ID, clientID, TokenGenerationAuthentication, clientMetadata)
	if err != nil {
		return nil, ErrInternalError
	}
	s.recordAuthEvent(reqCtx, userPool.ID, user.ID, username, clientID, authEventSignIn, authEventResponsePass)

	return s.withNewDeviceMetadata(store, userPool, user, challengeMember(params, "DEVICE_KEY"), map[string]interface{}{
		"AuthenticationResult": map[string]interface{}{
			"AccessToken":  accessToken,
			"IdToken":      idToken,
			"RefreshToken": refreshToken,
			"TokenType":    "Bearer",
			"ExpiresIn":    expiresIn,
		},
	}), nil
}

// verifyDeviceSrpClaim verifies the DEVICE_PASSWORD_VERIFIER claim: it
// recomputes the shared secret from the stored server ephemeral state and
// the device's password verifier, derives the HMAC key, and constant-time
// compares the expected claim — HMAC over device group key || device key ||
// secret block || timestamp — against the client's PASSWORD_CLAIM_SIGNATURE.
// As with the password claim, any mismatch surfaces as NotAuthorized to
// avoid leaking which leg of the verification failed.
func (s *CognitoService) verifyDeviceSrpClaim(
	store cognitostore.CognitoStoreInterface,
	challengeSession *cognitostore.ChallengeSession,
	userPool *cognitostore.UserPool,
	device *cognitostore.Device,
	claimSigB64, timestamp string,
) ([]byte, error) {
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
	verifier, ok := new(big.Int).SetString(device.DeviceSecretVerifierB, 16)
	if !ok || verifier.Sign() == 0 {
		return nil, ErrNotAuthorized
	}
	K, err := DeriveServerKey(A, B, b, verifier)
	if err != nil {
		return nil, ErrNotAuthorized
	}
	// The device claim carries the same few-seconds freshness contract as
	// the password claim; a stale or malformed TIMESTAMP fails closed.
	if !freshSrpTimestamp(timestamp, time.Now().UTC()) {
		return nil, ErrNotAuthorized
	}
	deviceGroupKey, ok := poolNameFromID(userPool.ID)
	if !ok {
		return nil, ErrInternalError
	}
	expectedSig := VerifyClaim(K, deviceGroupKey, device.DeviceKey, secretBlock, timestamp)
	if subtle.ConstantTimeCompare(clientSig, expectedSig) != 1 {
		return nil, ErrNotAuthorized
	}
	return expectedSig, nil
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
	// The PASSWORD_VERIFIER response must arrive within a few seconds of the
	// challenge; an older (or malformed) TIMESTAMP claim is a replay of a
	// captured signature and is refused with the same NotAuthorizedException
	// the model documents for the exceeded period.
	if !freshSrpTimestamp(timestamp, time.Now().UTC()) {
		s.recordAuthEvent(reqCtx, userPoolID, user.ID, username, clientID, authEventSignIn, authEventResponseFail)
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
		s.recordAuthEvent(reqCtx, userPoolID, user.ID, username, clientID, authEventSignIn, authEventResponseFail)
		return nil, ErrNotAuthorized
	}

	if !user.Enabled {
		s.recordAuthEvent(reqCtx, userPoolID, user.ID, username, clientID, authEventSignIn, authEventResponseFail)
		return nil, ErrNotAuthorized
	}

	// SRP-verified credentials fall under the pool's MFA configuration: a
	// second-factor challenge replaces the token issuance. The verifier
	// session is burned (deferred above), and the MFA challenge mints its
	// own session. A remembered device named by DEVICE_KEY — which may ride
	// in the PASSWORD_VERIFIER answer's challenge responses — replaces the
	// MFA challenge with device SRP authentication.
	pool, err := store.GetUserPool(userPoolID)
	if err != nil {
		return nil, ErrResourceNotFound
	}
	deviceKey := challengeMember(params, "DEVICE_KEY")
	challenge, err := s.nextChallengeAfterPrimary(reqCtx, store, pool, clientID, user, deviceKey)
	if err != nil {
		return nil, ErrInternalError
	}
	if challenge != nil {
		s.recordAuthEvent(reqCtx, userPoolID, user.ID, username, clientID, authEventSignIn, authEventResponseInProgress)
		return challenge, nil
	}

	accessToken, idToken, refreshToken, expiresIn, err := s.CreateTokens(reqCtx, userPoolID, user.ID, clientID, TokenGenerationAuthentication, clientMetadata)
	if err != nil {
		return nil, fmt.Errorf("failed to create tokens: %w", err)
	}
	s.recordAuthEvent(reqCtx, userPoolID, user.ID, username, clientID, authEventSignIn, authEventResponsePass)

	return s.withNewDeviceMetadata(store, pool, user, deviceKey, map[string]interface{}{
		"AuthenticationResult": map[string]interface{}{
			"AccessToken":  accessToken,
			"IdToken":      idToken,
			"RefreshToken": refreshToken,
			"TokenType":    "Bearer",
			"ExpiresIn":    expiresIn,
		},
	}), nil
}
