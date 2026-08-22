package cognitoidentityprovider

import (
	"context"
	"crypto/subtle"
	"fmt"
	"time"

	"vorpalstacks/internal/common/request"
	cognitostore "vorpalstacks/internal/store/aws/cognitoidentityprovider"

	"golang.org/x/crypto/bcrypt"
)

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
	req *request.ParsedRequest,
	challengeName string,
	userPool *cognitostore.UserPool,
	clientID, session string,
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
	if echoed := req.GetParam("USERNAME"); echoed != "" && echoed != username {
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
			recordChallengeFailure(store, challengeSession)
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
		// Constant-time comparison, as with the TOTP verification: the
		// attempt limit mitigates but does not eliminate the timing side
		// channel of a naive string comparison.
		if user.ConfirmationCode == "" || subtle.ConstantTimeCompare([]byte(user.ConfirmationCode), []byte(code)) != 1 {
			recordChallengeFailure(store, challengeSession)
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
		if user.ConfirmationCode == "" || subtle.ConstantTimeCompare([]byte(user.ConfirmationCode), []byte(code)) != 1 {
			recordChallengeFailure(store, challengeSession)
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
			recordChallengeFailure(store, challengeSession)
			return nil, ErrNotAuthorized
		}

	case "SELECT_MFA_TYPE", "SELECT_CHALLENGE":
		selected := req.GetParam("MFA_TYPE")
		if challengeName == "SELECT_CHALLENGE" {
			selected = req.GetParam("SELECTED_CHALLENGE")
			if selected == "" {
				selected = getStringParam(req.Parameters, "SelectedChallenge")
			}
		} else if selected == "" {
			selected = getStringParam(req.Parameters, "MfaType")
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
