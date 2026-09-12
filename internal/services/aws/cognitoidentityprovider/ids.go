package cognitoidentityprovider

import (
	"crypto/rand"
	"encoding/base32"
	"encoding/base64"
	"fmt"
	"math/big"
	"time"

	"github.com/google/uuid"
)

// The service's identifier and code generators, one set for the whole
// package. Each value keeps the encoding its wire surface defines: session
// IDs and event IDs are opaque identifiers, confirmation codes are six
// numeric digits, the TOTP secret is the base32 seed the standard
// specifies, and the authorisation code and client-secret members are
// high-entropy base64url strings.

// generateSessionID mints a challenge-session identifier.
func generateSessionID() string {
	return "SESSION_" + uuid.New().String()
}

// generateEventID mints an authentication-event identifier; the timestamp
// fallback keeps event recording alive when the system entropy source
// fails.
func generateEventID() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return fmt.Sprintf("evt-%d", time.Now().UnixNano())
	}
	return fmt.Sprintf("%x", b)
}

// generateImportJobID mints a user-import job identifier.
func generateImportJobID() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("failed to generate random ID: %w", err)
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

// generateAuthCode mints a single-use hosted-UI authorisation code.
func generateAuthCode() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

// generatedClientSecretLength is the length of a service-generated client
// secret, inside the Smithy ClientSecretType length bounds [24, 64]. The
// raw bounds live beside the pattern in client_token_core.go.
const generatedClientSecretLength = 56

// clientSecretAlphabet is the ClientSecretType pattern class [\w+]: upper-
// and lower-case letters, digits, underscore and plus. Base64 output can
// never satisfy the pattern — its padding '=' and alphabet member '/' fall
// outside the class — so values are sampled from this alphabet directly.
const clientSecretAlphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789_+"

// generateSecretValue mints a client-secret value in the Smithy
// ClientSecretType shape: characters drawn uniformly at random from the
// pattern class [\w+], at a length inside the shape's bounds.
func generateSecretValue() (string, error) {
	b := make([]byte, generatedClientSecretLength)
	for i := range b {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(clientSecretAlphabet))))
		if err != nil {
			return "", err
		}
		b[i] = clientSecretAlphabet[n.Int64()]
	}
	return string(b), nil
}

// totpSecretSize is the TOTP shared-secret length in bytes (160 bits, the
// SHA-1 block width the six-digit profile builds on).
const totpSecretSize = 20

// generateTOTPSecret mints the base32-encoded TOTP shared secret.
func generateTOTPSecret() (string, error) {
	secret := make([]byte, totpSecretSize)
	if _, err := rand.Read(secret); err != nil {
		return "", fmt.Errorf("failed to generate TOTP secret: %w", err)
	}
	return base32.StdEncoding.EncodeToString(secret), nil
}

// generateConfirmationCode mints a six-digit numeric confirmation code.
// The rejection loop keeps the modulo bias below the six-digit space.
func generateConfirmationCode() (string, error) {
	const maxCode = 1000000
	const limit = (1 << 24) / maxCode * maxCode
	for {
		b := make([]byte, 3)
		if _, err := rand.Read(b); err != nil {
			return "", err
		}
		n := int(b[0])<<16 | int(b[1])<<8 | int(b[2])
		if n < limit {
			return fmt.Sprintf("%06d", n%maxCode), nil
		}
	}
}
