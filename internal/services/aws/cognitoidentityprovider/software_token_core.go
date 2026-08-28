package cognitoidentityprovider

import (
	"context"

	"vorpalstacks/internal/common/request"
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
		// place of an AccessToken for the mid-sign-in MFA enrolment path. The
		// session must carry the MFA_SETUP challenge type so a session minted
		// for any other challenge cannot overwrite an existing (possibly
		// verified) MFA configuration. The service currently issues no
		// MFA_SETUP-typed sessions (the Lambda-facing designation path is
		// closed, see customFlowChallengeNames), so a session reaching this
		// branch cannot validate — the parameter handling remains because it
		// is part of the API contract.
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
		return nil, err
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
	if in.AccessToken == "" {
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

	userID, err := s.ValidateAccessToken(reqCtx, in.AccessToken)
	if err != nil {
		return nil, ErrNotAuthorized
	}

	user, err := store.GetUserByID(userID)
	if err != nil {
		return nil, ErrNotAuthorized
	}

	if user.SoftwareTokenMfa == nil || user.SoftwareTokenMfa.SecretKey == "" {
		return nil, ErrInvalidParameter
	}

	if !validateTOTPCode(user.SoftwareTokenMfa.SecretKey, in.UserCode) {
		return nil, ErrCodeMismatch
	}

	user.SoftwareTokenMfa.Verified = true
	user.SoftwareTokenMfa.Enabled = true
	if err := store.UpdateUser(user); err != nil {
		return nil, err
	}

	result := map[string]interface{}{
		"Status": "SUCCESS",
	}

	if in.Session != "" {
		result["Session"] = in.Session
	}

	return result, nil
}
