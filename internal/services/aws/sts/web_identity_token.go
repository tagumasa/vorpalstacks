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
// claim is absent. Used by AssumeRoleWithWebIdentity to honour M5
// (SubjectFromWebIdentityToken) and L4 (Audience) without requiring
// signature verification — the platform cannot reach external IdPs in
// TEST_MODE and SDK tests pass dummy tokens that are not real JWTs.
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
		// "aud" is the most common array-valued claim; join with spaces
		// so the caller can surface it as a single value.
		strs := make([]string, 0, len(v))
		for _, item := range v {
			if s, ok := item.(string); ok {
				strs = append(strs, s)
			}
		}
		return strings.Join(strs, " ")
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
