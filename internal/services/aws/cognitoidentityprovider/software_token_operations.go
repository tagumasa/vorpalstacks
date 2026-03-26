package cognitoidentityprovider

import (
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base32"
	"encoding/binary"
	"fmt"
	"math"
	"strings"
	"time"

	"vorpalstacks/internal/services/aws/common/request"
	cognitostore "vorpalstacks/internal/store/aws/cognitoidentityprovider"
)

func generateTOTPSecret() (string, error) {
	secret := make([]byte, 20)
	if _, err := rand.Read(secret); err != nil {
		return "", fmt.Errorf("failed to generate TOTP secret: %w", err)
	}
	return base32.StdEncoding.EncodeToString(secret), nil
}

func validateTOTPCode(secret, code string) bool {
	key, err := base32.StdEncoding.DecodeString(strings.ToUpper(secret))
	if err != nil {
		return false
	}

	timestamp := time.Now().Unix() / 30

	for offset := int64(-1); offset <= 1; offset++ {
		t := timestamp + offset
		buf := make([]byte, 8)
		binary.BigEndian.PutUint64(buf, uint64(t))

		mac := hmac.New(sha1.New, key)
		mac.Write(buf)
		hash := mac.Sum(nil)

		off := hash[len(hash)-1] & 0x0f
		truncated := binary.BigEndian.Uint32(hash[off : off+4])
		truncated &= 0x7fffffff
		digits := truncated % uint32(math.Pow10(6))

		if fmt.Sprintf("%06d", digits) == code {
			return true
		}
	}

	return false
}

// AssociateSoftwareToken generates a TOTP secret and associates it with the user for MFA setup.
func (s *CognitoService) AssociateSoftwareToken(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	accessToken := getParam(req, "AccessToken")
	session := getParam(req, "Session")

	if accessToken == "" && session == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	token, err := store.GetAccessTokenByValue(accessToken)
	if err != nil {
		return nil, ErrNotAuthorized
	}

	user, err := store.GetUserByID(token.UserID)
	if err != nil {
		return nil, ErrNotAuthorized
	}

	secret, err := generateTOTPSecret()
	if err != nil {
		return nil, ErrInvalidParameter
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

	return map[string]interface{}{
		"SecretCode": secret,
	}, nil
}

// VerifySoftwareToken verifies a TOTP code provided by the user during MFA setup.
func (s *CognitoService) VerifySoftwareToken(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	accessToken := getParam(req, "AccessToken")
	userCode := getParam(req, "UserCode")
	session := getParam(req, "Session")

	if accessToken == "" {
		return nil, ErrInvalidParameter
	}
	if userCode == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	token, err := store.GetAccessTokenByValue(accessToken)
	if err != nil {
		return nil, ErrNotAuthorized
	}

	user, err := store.GetUserByID(token.UserID)
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
