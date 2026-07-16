package sts

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"fmt"
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
// from the caller-supplied audience list.
func (m *webIdentityTokenManager) generateToken(subject, accountID string, audiences []string, algorithm string, durationSeconds int) (string, time.Time, error) {
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
