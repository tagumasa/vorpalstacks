package cognitoidentityprovider

import (
	"context"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	cognitostore "vorpalstacks/internal/store/aws/cognitoidentityprovider"
)

// AssociateSoftwareTokenInput carries the wire parameters of
// AssociateSoftwareToken.
type AssociateSoftwareTokenInput struct {
	AccessToken string
	Session     string
}

// VerifySoftwareTokenInput carries the wire parameters of
// VerifySoftwareToken.
type VerifySoftwareTokenInput struct {
	AccessToken string
	UserCode    string
	Session     string
}

// AdminDeleteSoftwareTokenInput carries the wire parameters of
// AdminDeleteSoftwareToken.
type AdminDeleteSoftwareTokenInput struct {
	UserPoolID string
	Username   string
}

// associateSoftwareTokenCore generates a TOTP secret and associates it with
// the user for MFA setup.
func (s *CognitoService) associateSoftwareTokenCore(ctx context.Context, reqCtx *request.RequestContext, in AssociateSoftwareTokenInput) (interface{}, error) {
	if in.AccessToken == "" && in.Session == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	var user *cognitostore.User
	if in.AccessToken != "" {
		userID, err := s.ValidateAccessToken(reqCtx, in.AccessToken)
		if err != nil {
			return nil, ErrNotAuthorized
		}
		user, err = store.GetUserByID(userID)
		if err != nil {
			return nil, ErrNotAuthorized
		}
	} else if in.Session != "" {
		// Session-based flow: the Amazon Cognito API accepts a Session in
		// place of an AccessToken for the mid-sign-in MFA enrolment path
		// (the MFA_SETUP challenge issued at sign-in). The session must
		// carry the MFA_SETUP challenge type so a session minted for any
		// other challenge cannot overwrite an existing (possibly verified)
		// MFA configuration.
		challengeSession, err := validateChallengeSession(store, in.Session, "MFA_SETUP", "", "", "")
		if err != nil {
			return nil, ErrNotAuthorized
		}
		user, err = store.GetUser(challengeSession.UserPoolID, challengeSession.Username)
		if err != nil {
			return nil, ErrNotAuthorized
		}
	}
	if user == nil {
		return nil, ErrNotAuthorized
	}

	// The access-token path returns the existing secret for an
	// already-associated registration instead of re-arming a fresh one, so
	// a stray call cannot silently destroy a working TOTP registration. The
	// mid-sign-in MFA_SETUP session path is the re-enrolment ceremony and
	// keeps re-arming: the session type check above is what guards it.
	if in.AccessToken != "" && user.SoftwareTokenMfa != nil && user.SoftwareTokenMfa.SecretKey != "" {
		result := map[string]interface{}{
			"SecretCode": user.SoftwareTokenMfa.SecretKey,
		}
		return result, nil
	}

	secret, err := generateTOTPSecret()
	if err != nil {
		return nil, ErrInternalError
	}
	user.SoftwareTokenMfa = &cognitostore.SoftwareTokenMfaSettings{
		Enabled:      false,
		PreferredMfa: false,
		SecretKey:    secret,
		Verified:     false,
	}

	if err := store.UpdateUser(user); err != nil {
		return nil, ErrInternalError
	}

	result := map[string]interface{}{
		"SecretCode": secret,
	}
	if in.Session != "" {
		result["Session"] = in.Session
	}
	return result, nil
}

// verifySoftwareTokenCore verifies a TOTP code provided by the user during
// MFA setup.
func (s *CognitoService) verifySoftwareTokenCore(ctx context.Context, reqCtx *request.RequestContext, in VerifySoftwareTokenInput) (interface{}, error) {
	if in.AccessToken == "" && in.Session == "" {
		return nil, ErrInvalidParameter
	}
	if in.UserCode == "" {
		return nil, ErrInvalidParameter
	}
	if !totpCodePattern.MatchString(in.UserCode) {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	var user *cognitostore.User
	var challengeSession *cognitostore.ChallengeSession
	if in.AccessToken != "" {
		userID, err := s.ValidateAccessToken(reqCtx, in.AccessToken)
		if err != nil {
			return nil, ErrNotAuthorized
		}
		user, err = store.GetUserByID(userID)
		if err != nil {
			return nil, ErrNotAuthorized
		}
	} else {
		// Session-based flow: mid-sign-in MFA enrolment (the MFA_SETUP
		// challenge) carries the session instead of an access token. The
		// same MFA_SETUP-typed session requirement as AssociateSoftwareToken
		// binds the verification to the enrolment it completes.
		var err error
		challengeSession, err = validateChallengeSession(store, in.Session, "MFA_SETUP", "", "", "")
		if err != nil {
			return nil, ErrNotAuthorized
		}
		user, err = store.GetUser(challengeSession.UserPoolID, challengeSession.Username)
		if err != nil {
			return nil, ErrNotAuthorized
		}
	}
	if user == nil {
		return nil, ErrNotAuthorized
	}

	if user.SoftwareTokenMfa == nil || user.SoftwareTokenMfa.SecretKey == "" {
		return nil, ErrInvalidParameter
	}

	// The registration's attempt budget binds the whole access-token
	// verification path, not only its failures: once exhausted, even the
	// correct code is refused until the registration is re-associated. The
	// session path gets the same bound from validateChallengeSession.
	if challengeSession == nil && user.SoftwareTokenMfa.FailedAttempts >= maxChallengeAttempts {
		return nil, ErrNotAuthorized
	}

	if !validateTOTPCode(user.SoftwareTokenMfa.SecretKey, in.UserCode) {
		// Every enrolment verification attempt is budgeted, like every
		// other challenge answer: the session path consumes the challenge
		// session's failure budget, and the access-token path — which
		// carries no session — counts on the registration itself, so the
		// six-digit code cannot be brute-forced at wire speed.
		if challengeSession != nil {
			recordChallengeFailure(store, challengeSession)
		} else {
			user.SoftwareTokenMfa.FailedAttempts++
			if uerr := store.UpdateUser(user); uerr != nil {
				return nil, ErrInternalError
			}
		}
		return nil, ErrCodeMismatch
	}

	user.SoftwareTokenMfa.Verified = true
	user.SoftwareTokenMfa.Enabled = true
	user.SoftwareTokenMfa.FailedAttempts = 0
	if err := store.UpdateUser(user); err != nil {
		return nil, ErrInternalError
	}

	result := map[string]interface{}{
		"Status": "SUCCESS",
	}

	if in.Session != "" {
		result["Session"] = in.Session
	}

	return result, nil
}

// adminDeleteSoftwareTokenCore removes a user's TOTP software token
// registration together with any preference for TOTP MFA; the other MFA
// factors the user has registered stay untouched. A user with no software
// token registration reports ResourceNotFoundException — the deletion
// addresses a factor that is not there, rather than succeeding silently.
func (s *CognitoService) adminDeleteSoftwareTokenCore(reqCtx *request.RequestContext, in AdminDeleteSoftwareTokenInput) (interface{}, error) {
	if in.UserPoolID == "" || in.Username == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	user, err := store.GetUser(in.UserPoolID, in.Username)
	if err != nil {
		return nil, ErrUserNotFound
	}
	if user.SoftwareTokenMfa == nil || user.SoftwareTokenMfa.SecretKey == "" {
		return nil, ErrResourceNotFound
	}

	// Clearing the whole registration also clears its PreferredMfa flag,
	// leaving the remaining factors' preference state as it was.
	user.SoftwareTokenMfa = nil
	if err := store.UpdateUser(user); err != nil {
		return nil, ErrInternalError
	}

	return response.EmptyResponse(), nil
}
