package cognitoidentityprovider

import (
	"fmt"
	"time"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/core/logs"
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
// SELECT_MFA_TYPE response: the documented ANSWER vocabulary is
// SMS_MFA|EMAIL_MFA|SOFTWARE_TOKEN_MFA.
var selectMfaTypeChoices = map[string]bool{
	"SMS_MFA":            true,
	"EMAIL_MFA":          true,
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
// after every wrong verification answer so retries stay bounded. A session
// whose failure state cannot be persisted is invalidated on the spot: the
// attempt budget must never weaken because the store dropped a write.
func recordChallengeFailure(store cognitostore.CognitoStoreInterface, cs *cognitostore.ChallengeSession) {
	cs.FailedAttempts++
	if cs.FailedAttempts >= maxChallengeAttempts {
		_ = store.DeleteChallengeSession(cs.SessionID)
		return
	}
	if err := store.SaveChallengeSession(cs); err != nil {
		logs.Error("failed to persist challenge-session failure count; invalidating session", logs.Err(err))
		_ = store.DeleteChallengeSession(cs.SessionID)
	}
}

// newChallengeSession builds an unpersisted challenge session of the given
// type. Callers that need extra session state (SRP values, WebAuthn challenge
// bytes, device bindings, the MFA role) set the fields on the returned struct
// and persist it with SaveChallengeSession.
func newChallengeSession(userPoolID, clientID, username, challengeName string, ttl time.Duration) *cognitostore.ChallengeSession {
	return &cognitostore.ChallengeSession{
		SessionID:     generateSessionID(),
		UserPoolID:    userPoolID,
		ClientID:      clientID,
		Username:      username,
		ChallengeName: challengeName,
		CreatedAt:     time.Now().UTC(),
		ExpiresAt:     time.Now().UTC().Add(ttl),
	}
}

// otpDeliveryMedium reports the delivery medium of an OTP-typed challenge —
// the challenge names that carry a one-time code — and whether the name is
// one of them. Membership of this set is the single definition of "carries a
// one-time code" for minting and delivery alike.
func otpDeliveryMedium(challengeName string) (string, bool) {
	switch challengeName {
	case "SMS_MFA", "SMS_OTP":
		return "SMS", true
	case "EMAIL_OTP":
		return "EMAIL", true
	}
	return "", false
}

// mintChallengeSessionState generates the challenge-type-specific session
// state — the one-time code for OTP-typed sessions, generated at mint and
// valid for the session's lifetime alone — and persists the session. This is
// the single mint definition for every issuance path; callers that need extra
// session state (the MFA role, SRP values, WebAuthn bytes) set the fields on
// the session before calling.
func mintChallengeSessionState(store cognitostore.CognitoStoreInterface, cs *cognitostore.ChallengeSession) error {
	if _, otp := otpDeliveryMedium(cs.ChallengeName); otp {
		code, err := generateConfirmationCode()
		if err != nil {
			return err
		}
		cs.OTPCode = code
	}
	return store.SaveChallengeSession(cs)
}

// mintChallengeSession creates and persists a fresh challenge session of the
// given type, returning the new session ID. OTP-typed sessions (SMS_MFA,
// SMS_OTP, EMAIL_OTP) carry their one-time code on the session itself,
// generated here at mint and valid for the session's lifetime alone.
func mintChallengeSession(store cognitostore.CognitoStoreInterface, userPoolID, clientID, username, challengeName string, ttl time.Duration) (string, error) {
	cs := newChallengeSession(userPoolID, clientID, username, challengeName, ttl)
	if err := mintChallengeSessionState(store, cs); err != nil {
		return "", err
	}
	return cs.SessionID, nil
}

// deliverChallengeCode records an OTP challenge's code issuance in the
// message-delivery notification log. The model has the user pool deliver the
// code with the challenge ("Respond with the code that your user pool
// delivered…"); on this platform the notification log is the recorded
// delivery mechanism for one-time codes, so the record is written at mint,
// where AWS dispatches the message. Non-OTP challenges record nothing.
func (s *CognitoService) deliverChallengeCode(reqCtx *request.RequestContext, userPoolID, challengeName string) {
	message, ok := challengeDeliveryMessage(challengeName)
	if !ok {
		return
	}
	s.publishNotificationLogCore(reqCtx, userPoolID, message)
}

// challengeDeliveryMessage renders the message-delivery record for an
// OTP-typed challenge's code; the second return is false for challenges that
// carry no one-time code.
func challengeDeliveryMessage(challengeName string) (string, bool) {
	medium, otp := otpDeliveryMedium(challengeName)
	if !otp {
		return "", false
	}
	return fmt.Sprintf("%s challenge code delivery via %s", challengeName, medium), true
}
