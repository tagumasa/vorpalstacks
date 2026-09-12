package cognitoidentityprovider

import (
	"context"
	"crypto/hmac"
	"crypto/sha1"
	"crypto/subtle"
	"encoding/base32"
	"encoding/binary"
	"fmt"
	"math"
	"strings"
	"time"

	"vorpalstacks/internal/common/request"
)

// totpCodeModulus is 10^totpCodeDigits, used to truncate the HMAC to the
// fixed six-digit code length. It is derived from the digit count so the
// two cannot drift apart.
var totpCodeModulus = uint32(math.Pow(10, totpCodeDigits))

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
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_AssociateSoftwareToken.html
func (s *CognitoService) AssociateSoftwareToken(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return s.associateSoftwareTokenCore(ctx, reqCtx, AssociateSoftwareTokenInput{
		AccessToken: req.GetParam("AccessToken"),
		Session:     req.GetParam("Session"),
	})
}

// VerifySoftwareToken verifies a TOTP code provided by the user during MFA setup.
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_VerifySoftwareToken.html
func (s *CognitoService) VerifySoftwareToken(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return s.verifySoftwareTokenCore(ctx, reqCtx, VerifySoftwareTokenInput{
		AccessToken: req.GetParam("AccessToken"),
		UserCode:    req.GetParam("UserCode"),
		Session:     req.GetParam("Session"),
	})
}

// AdminDeleteSoftwareToken deletes a user's registered TOTP MFA factor.
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_AdminDeleteSoftwareToken.html
func (s *CognitoService) AdminDeleteSoftwareToken(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return s.adminDeleteSoftwareTokenCore(reqCtx, AdminDeleteSoftwareTokenInput{
		UserPoolID: req.GetParam("UserPoolId"),
		Username:   req.GetParam("Username"),
	})
}
