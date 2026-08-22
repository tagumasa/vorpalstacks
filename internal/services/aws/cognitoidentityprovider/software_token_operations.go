package cognitoidentityprovider

import (
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1"
	"crypto/subtle"
	"encoding/base32"
	"encoding/binary"
	"fmt"
	"strings"
	"time"

	"vorpalstacks/internal/common/request"
	cognitostore "vorpalstacks/internal/store/aws/cognitoidentityprovider"
)

const totpSecretSize = 20

// totpCodeModulus is 10^totpCodeDigits, used to truncate the HMAC to the
// fixed six-digit code length.
const totpCodeModulus uint32 = 1_000_000

func generateTOTPSecret() (string, error) {
	secret := make([]byte, totpSecretSize)
	if _, err := rand.Read(secret); err != nil {
		return "", fmt.Errorf("failed to generate TOTP secret: %w", err)
	}
	return base32.StdEncoding.EncodeToString(secret), nil
}

// totpCodeAt derives the six-digit TOTP code for the given 30-second step
// from the decoded secret key.
func totpCodeAt(key []byte, step int64) string {
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(step))

	mac := hmac.New(sha1.New, key)
	mac.Write(buf)
	hash := mac.Sum(nil)

	off := hash[len(hash)-1] & 0x0f
	truncated := binary.BigEndian.Uint32(hash[off : off+4])
	truncated &= 0x7fffffff
	return fmt.Sprintf("%0*d", totpCodeDigits, truncated%totpCodeModulus)
}

// validateTOTPCode reports whether code matches the TOTP generated from
// secret for the current 30-second step or one step either side (clock
// drift tolerance). Comparison is constant-time so matching does not leak
// timing information about how many digits were correct.
func validateTOTPCode(secret, code string) bool {
	key, err := base32.StdEncoding.DecodeString(strings.ToUpper(secret))
	if err != nil {
		return false
	}

	step := time.Now().Unix() / totpTimeStepSec
	for offset := int64(-totpAllowedDrift); offset <= totpAllowedDrift; offset++ {
		if subtle.ConstantTimeCompare([]byte(totpCodeAt(key, step+offset)), []byte(code)) == 1 {
			return true
		}
	}

	return false
}

// AssociateSoftwareToken generates a TOTP secret and associates it with the user for MFA setup.
func (s *CognitoService) AssociateSoftwareToken(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	accessToken := req.GetParam("AccessToken")
	session := req.GetParam("Session")

	if accessToken == "" && session == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	var user *cognitostore.User
	if accessToken != "" {
		userID, err := s.ValidateAccessToken(reqCtx, accessToken)
		if err != nil {
			return nil, ErrNotAuthorized
		}
		user, err = store.GetUserByID(userID)
		if err != nil {
			return nil, ErrNotAuthorized
		}
	} else if session != "" {
		// Session-based flow: the Amazon Cognito API accepts a Session in
		// place of an AccessToken for the mid-sign-in MFA enrolment path. The
		// session must carry the MFA_SETUP challenge type so a session minted
		// for any other challenge cannot overwrite an existing (possibly
		// verified) MFA configuration. The service currently issues no
		// MFA_SETUP-typed sessions (the Lambda-facing designation path is
		// closed, see customFlowChallengeNames), so a session reaching this
		// branch cannot validate — the parameter handling remains because it
		// is part of the API contract.
		challengeSession, err := validateChallengeSession(store, session, "MFA_SETUP", "", "", "")
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
	if session != "" {
		result["Session"] = session
	}
	return result, nil
}

// VerifySoftwareToken verifies a TOTP code provided by the user during MFA setup.
func (s *CognitoService) VerifySoftwareToken(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	accessToken := req.GetParam("AccessToken")
	userCode := req.GetParam("UserCode")
	session := req.GetParam("Session")

	if accessToken == "" {
		return nil, ErrInvalidParameter
	}
	if userCode == "" {
		return nil, ErrInvalidParameter
	}
	if !totpCodePattern.MatchString(userCode) {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	userID, err := s.ValidateAccessToken(reqCtx, accessToken)
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

	if !validateTOTPCode(user.SoftwareTokenMfa.SecretKey, userCode) {
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

	if session != "" {
		result["Session"] = session
	}

	return result, nil
}
