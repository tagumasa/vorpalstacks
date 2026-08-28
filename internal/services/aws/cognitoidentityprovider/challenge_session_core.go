package cognitoidentityprovider

import (
	"time"

	cognitostore "vorpalstacks/internal/store/aws/cognitoidentityprovider"
)

// maxChallengeAttempts bounds the number of failed verification attempts
// (OTP codes, TOTP codes, passwords, custom answers) accepted within a
// single challenge session before the session is invalidated. Without this
// bound a six-digit code can be brute-forced by replaying
// RespondToAuthChallenge against one session.
const maxChallengeAttempts = 5

// selectChallengeChoices lists the challenge types a client may choose in a
// SELECT_CHALLENGE response. Per the RespondToAuthChallenge contract the
// selectable sign-in challenges are exactly PASSWORD, PASSWORD_SRP,
// EMAIL_OTP, SMS_OTP and WEB_AUTHN. MFA challenges and MFA_SETUP are issued
// by the server at authentication time and must never be client-selectable:
// accepting MFA_SETUP here would let an unauthenticated caller mint an
// MFA_SETUP-typed session through USER_AUTH (which issues selector sessions
// without verifying credentials) and overwrite a victim's TOTP configuration
// via AssociateSoftwareToken.
var selectChallengeChoices = map[string]bool{
	"PASSWORD":     true,
	"PASSWORD_SRP": true,
	"EMAIL_OTP":    true,
	"SMS_OTP":      true,
	"WEB_AUTHN":    true,
}

// selectMfaTypeChoices lists the MFA types a client may choose in a
// SELECT_MFA_TYPE response: SMS_MFA or SOFTWARE_TOKEN_MFA.
var selectMfaTypeChoices = map[string]bool{
	"SMS_MFA":            true,
	"SOFTWARE_TOKEN_MFA": true,
}

// validateChallengeSession loads a challenge session by ID and enforces the
// binding contract shared by every challenge-response path: the session must
// exist, be unexpired, be within the failed-attempt budget, carry the
// expected challenge type, and belong to the same user pool, app client and
// user as the request. Empty expected values skip the corresponding check
// for callers (such as AssociateSoftwareToken) whose API does not carry that
// parameter. Any mismatch returns NotAuthorizedException without revealing
// which leg failed, mirroring the SRP verifier path.
//
// The session is not consumed here; callers decide when to burn it.
func validateChallengeSession(
	store cognitostore.CognitoStoreInterface,
	sessionID, expectedChallenge, userPoolID, clientID, username string,
) (*cognitostore.ChallengeSession, error) {
	if sessionID == "" {
		return nil, ErrNotAuthorized
	}
	cs, err := store.GetChallengeSession(sessionID)
	if err != nil || cs == nil {
		return nil, ErrNotAuthorized
	}
	if cs.ChallengeName != expectedChallenge {
		return nil, ErrNotAuthorized
	}
	if userPoolID != "" && cs.UserPoolID != userPoolID {
		return nil, ErrNotAuthorized
	}
	if clientID != "" && cs.ClientID != clientID {
		return nil, ErrNotAuthorized
	}
	if username != "" && cs.Username != username {
		return nil, ErrNotAuthorized
	}
	if !cs.ExpiresAt.IsZero() && time.Now().UTC().After(cs.ExpiresAt) {
		_ = store.DeleteChallengeSession(sessionID)
		return nil, ErrNotAuthorized
	}
	if cs.FailedAttempts >= maxChallengeAttempts {
		_ = store.DeleteChallengeSession(sessionID)
		return nil, ErrNotAuthorized
	}
	return cs, nil
}

// recordChallengeFailure increments the session's failed-attempt counter and
// invalidates the session once the budget is exhausted. Callers invoke this
// after every wrong verification answer so retries stay bounded.
func recordChallengeFailure(store cognitostore.CognitoStoreInterface, cs *cognitostore.ChallengeSession) {
	cs.FailedAttempts++
	if cs.FailedAttempts >= maxChallengeAttempts {
		_ = store.DeleteChallengeSession(cs.SessionID)
		return
	}
	_ = store.SaveChallengeSession(cs)
}

// mintChallengeSession creates and persists a fresh challenge session of the
// given type, returning the new session ID.
func mintChallengeSession(store cognitostore.CognitoStoreInterface, userPoolID, clientID, username, challengeName string, ttl time.Duration) (string, error) {
	sessionID := generateSessionID()
	cs := &cognitostore.ChallengeSession{
		SessionID:     sessionID,
		UserPoolID:    userPoolID,
		ClientID:      clientID,
		Username:      username,
		ChallengeName: challengeName,
		CreatedAt:     time.Now().UTC(),
		ExpiresAt:     time.Now().UTC().Add(ttl),
	}
	if err := store.SaveChallengeSession(cs); err != nil {
		return "", err
	}
	return sessionID, nil
}
