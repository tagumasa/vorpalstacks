package sts

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// webIdentityTokenManager lazily generates signing key pairs for the RS256 and
// ES384 algorithms supported by the GetWebIdentityToken API. Keys are shared
// across all callers in the service instance and regenerated on each process
// restart — acceptable because the tokens are short-lived (default 5 minutes,
// maximum 1 hour).
type webIdentityTokenManager struct {
	once      sync.Once
	rsaKey    *rsa.PrivateKey
	ecKey     *ecdsa.PrivateKey
	rsaKeyID  string
	ecKeyID   string
	issuer    string
	initError error
}

// webIdentityTokenManagerInstance is the package-level singleton. Initialised
// on first GetWebIdentityToken call.
var webIdentityTokenManagerInstance = &webIdentityTokenManager{
	issuer: "sts.amazonaws.com",
}

// init generates the RSA and ECDSA key pairs exactly once.
func (m *webIdentityTokenManager) init() {
	m.once.Do(func() {
		rsaKey, err := rsa.GenerateKey(rand.Reader, 2048)
		if err != nil {
			m.initError = fmt.Errorf("failed to generate RSA key for web identity tokens: %w", err)
			return
		}
		m.rsaKey = rsaKey
		m.rsaKeyID = uuid.New().String()

		ecKey, err := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
		if err != nil {
			m.initError = fmt.Errorf("failed to generate EC P-384 key for web identity tokens: %w", err)
			return
		}
		m.ecKey = ecKey
		m.ecKeyID = uuid.New().String()
	})
}

// generateToken creates a signed JWT representing the calling AWS identity.
// The algorithm must be either RS256 or ES384. The aud claim is populated
// from the caller-supplied audience list. When tags is non-empty, the tags
// are embedded in the JWT under the standard "tags" claim so external
// services consuming the token can read the caller-attached session tags.
func (m *webIdentityTokenManager) generateToken(subject, accountID string, audiences []string, algorithm string, durationSeconds int, tags map[string]string) (string, time.Time, error) {
	m.init()
	if m.initError != nil {
		return "", time.Time{}, m.initError
	}

	now := time.Now().UTC()
	expiration := now.Add(time.Duration(durationSeconds) * time.Second)

	claims := jwt.MapClaims{
		"iss":     m.issuer,
		"sub":     subject,
		"aud":     audiences,
		"iat":     now.Unix(),
		"exp":     expiration.Unix(),
		"jti":     uuid.New().String(),
		"account": accountID,
	}
	if len(tags) > 0 {
		claims["tags"] = tags
	}

	var token *jwt.Token
	var keyID string
	switch algorithm {
	case "RS256":
		token = jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
		keyID = m.rsaKeyID
	case "ES384":
		token = jwt.NewWithClaims(jwt.SigningMethodES384, claims)
		keyID = m.ecKeyID
	default:
		return "", time.Time{}, fmt.Errorf("unsupported signing algorithm: %s", algorithm)
	}
	token.Header["kid"] = keyID

	signed, err := token.SignedString(m.signingKey(algorithm))
	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to sign web identity token: %w", err)
	}
	return signed, expiration, nil
}

// signingKey returns the private key matching the requested algorithm.
func (m *webIdentityTokenManager) signingKey(algorithm string) interface{} {
	switch algorithm {
	case "RS256":
		return m.rsaKey
	case "ES384":
		return m.ecKey
	default:
		return m.rsaKey
	}
}

// extractJWTClaim returns the value of claim from a JSON Web Token's payload
// segment, or the empty string when the token is not a parseable JWT or the
// claim is absent. Used by AssumeRoleWithWebIdentity to extract the subject
// and audience without requiring signature verification — the platform
// cannot reach external IdPs in TEST_MODE and SDK tests pass dummy tokens
// that are not real JWTs.
func extractJWTClaim(token, claim string) string {
	parts := strings.Split(token, ".")
	if len(parts) < 2 {
		return ""
	}
	raw, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		raw, err = base64.StdEncoding.DecodeString(parts[1])
		if err != nil {
			return ""
		}
	}
	var payload map[string]interface{}
	if err := json.Unmarshal(raw, &payload); err != nil {
		return ""
	}
	val, ok := payload[claim]
	if !ok {
		return ""
	}
	switch v := val.(type) {
	case string:
		return v
	case []interface{}:
		// "aud" is the most common array-valued claim. AWS returns a
		// single Audience string (the OIDC client ID that matched the
		// provider configuration). VorpalStacks does not track per-
		// provider client IDs, so return the first array element as
		// the best approximation of AWS behaviour.
		for _, item := range v {
			if s, ok := item.(string); ok {
				return s
			}
		}
		return ""
	default:
		return ""
	}
}

// validateWebIdentityDurationSeconds validates the DurationSeconds parameter
// for GetWebIdentityToken per the Smithy webIdentityTokenDurationSecondsType
// trait (range 60-3600, default 300).
func validateWebIdentityDurationSeconds(durationSeconds int) (int, error) {
	if durationSeconds == 0 {
		return 300, nil
	}
	if durationSeconds < 60 || durationSeconds > 3600 {
		return 0, ErrInvalidWebIdentityDuration
	}
	return durationSeconds, nil
}

// isJWTExpired performs a soft expiry check on a web identity token. It
// extracts the standard "exp" claim from the JWT payload and returns true
// when the token is past its expiry. Tokens that are not parseable JWTs
// (e.g. dummy tokens used in SDK tests) return false — the check is
// intentionally lenient to preserve TEST_MODE compatibility while still
// rejecting real expired OIDC tokens. Signature verification is not
// performed because VorpalStacks cannot reach external IdPs.
//
// This function is only reached in TEST_MODE; in production,
// AssumeRoleWithWebIdentity returns ErrIDPCommunicationError before
// reaching this check.
func isJWTExpired(token string) bool {
	parts := strings.Split(token, ".")
	if len(parts) < 2 {
		return false
	}
	raw, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		raw, err = base64.StdEncoding.DecodeString(parts[1])
		if err != nil {
			return false
		}
	}
	var payload map[string]interface{}
	if err := json.Unmarshal(raw, &payload); err != nil {
		return false
	}
	exp, ok := payload["exp"]
	if !ok {
		// A parseable JWT without an exp claim is malformed — OIDC
		// tokens must declare an expiry. Fail-closed to reject it.
		return true
	}
	// json.Unmarshal into interface{} always returns float64 for JSON
	// numbers, so only that case needs handling.
	if v, ok := exp.(float64); ok {
		return time.Now().UTC().Unix() >= int64(v)
	}
	return false
}
