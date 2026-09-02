package inspection

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// TokenCookieName is the cookie that carries the client token of the
// CAPTCHA and Challenge actions — the Developer Guide's token
// characteristics page names it aws-waf-token.
const TokenCookieName = "aws-waf-token"

// challengeTokenVersion is the payload version; a mismatch rejects the
// token so a format change cannot misread old values.
const challengeTokenVersion = 1

// ChallengeSolutionHexZeros is the number of leading zero hex digits a
// challenge solution's hash must carry. The interstitial's script and
// the exchange endpoint both reference this definition.
const ChallengeSolutionHexZeros = 4

// ChallengeToken is the payload carried by the aws-waf-token cookie.
// AWS encrypts its tokens and publishes only their behavioural state:
// the timestamp of the client's latest successful silent-challenge
// response and that of its latest successful CAPTCHA response. This
// platform's token is HMAC-signed JSON carrying the same state plus the
// domains it was issued for, which is what the token-domain acceptance
// rules need.
type ChallengeToken struct {
	Version           int64    `json:"v"`
	IssuedAt          int64    `json:"iat"`
	ChallengeSolvedAt int64    `json:"chs,omitempty"`
	CaptchaSolvedAt   int64    `json:"cps,omitempty"`
	Domains           []string `json:"d,omitempty"`
}

// TokenValidator verifies token values. The service injects the
// implementation that holds the signing secret; the evaluator stays free
// of key management.
type TokenValidator interface {
	ValidateToken(value string) (ChallengeToken, bool)
}

// SignToken serialises and signs a token with the service's secret. The
// wire form is base64url(payload).base64url(HMAC-SHA256(payload)).
func SignToken(secret []byte, token ChallengeToken) string {
	token.Version = challengeTokenVersion
	payload, err := json.Marshal(token)
	if err != nil {
		return ""
	}
	encoded := base64.RawURLEncoding.EncodeToString(payload)
	mac := hmac.New(sha256.New, secret)
	mac.Write([]byte(encoded))
	return encoded + "." + base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}

// ParseToken verifies a token value's signature and decodes its
// payload. A malformed value, a bad signature or an unknown version
// fails with ok=false.
func ParseToken(secret []byte, value string) (ChallengeToken, bool) {
	var token ChallengeToken
	encoded, signature, found := strings.Cut(value, ".")
	if !found {
		return token, false
	}
	mac := hmac.New(sha256.New, secret)
	mac.Write([]byte(encoded))
	expected := mac.Sum(nil)
	given, err := base64.RawURLEncoding.DecodeString(signature)
	if err != nil || !hmac.Equal(expected, given) {
		return token, false
	}
	payload, err := base64.RawURLEncoding.DecodeString(encoded)
	if err != nil {
		return token, false
	}
	if err := json.Unmarshal(payload, &token); err != nil {
		return token, false
	}
	if token.Version != challengeTokenVersion {
		return ChallengeToken{}, false
	}
	return token, true
}

// CoversHost reports whether the token's domain scope accepts the given
// request host. A token issued for a domain is valid for that host and
// its subdomains — the same scope a cookie set on that domain would
// have. The TokenDomains documentation accepts the protected resource's
// host plus every configured token domain.
func (t ChallengeToken) CoversHost(host string) bool {
	host = strings.ToLower(strings.TrimSpace(host))
	if host == "" {
		return false
	}
	for _, domain := range t.Domains {
		domain = strings.ToLower(strings.TrimSpace(domain))
		if domain == "" {
			continue
		}
		if host == domain || strings.HasSuffix(host, "."+domain) {
			return true
		}
	}
	return false
}

// SolvedWithin reports whether the kind's solve timestamp exists and
// lies within the immunity window ending at now. kind is one of the
// ActionCaptcha and ActionChallenge constants; each kind consults only
// its own timestamp, matching the two independent timestamps the token
// characteristics page describes.
func (t ChallengeToken) SolvedWithin(kind string, now time.Time, immunity time.Duration) bool {
	var solved int64
	switch kind {
	case ActionCaptcha:
		solved = t.CaptchaSolvedAt
	case ActionChallenge:
		solved = t.ChallengeSolvedAt
	default:
		return false
	}
	if solved <= 0 || immunity <= 0 {
		return false
	}
	elapsed := now.Unix() - solved
	return elapsed >= 0 && elapsed < int64(immunity/time.Second)
}

// VerifyChallengeSolution reports whether counter solves the challenge:
// the hex digest of challengeID "." counter must start with
// ChallengeSolutionHexZeros zero digits. The proof is stateless — the
// server only needs the challenge identifier it issued and the counter
// the client returns.
func VerifyChallengeSolution(challengeID, counter string) bool {
	if challengeID == "" || counter == "" {
		return false
	}
	if _, err := strconv.ParseUint(counter, 10, 64); err != nil {
		return false
	}
	digest := fmt.Sprintf("%x", sha256.Sum256([]byte(challengeID+"."+counter)))
	return strings.HasPrefix(digest, strings.Repeat("0", ChallengeSolutionHexZeros))
}
