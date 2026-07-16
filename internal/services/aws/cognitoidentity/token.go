package cognitoidentity

import (
	"crypto/rsa"
	"fmt"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"vorpalstacks/pkg/vsjwt"
)

// tokenManager lazily generates an RSA key pair for signing Cognito Identity
// OpenID Connect tokens. The key is shared across all identity pools in the
// service instance and is regenerated on each process restart — acceptable
// because the tokens are short-lived (default 15 minutes).
type tokenManager struct {
	once       sync.Once
	privateKey *rsa.PrivateKey
	keyID      string
	issuer     string
	err        error
}

func newTokenManager() *tokenManager {
	return &tokenManager{
		issuer: "cognito-identity.amazonaws.com",
	}
}

func (tm *tokenManager) init() {
	tm.once.Do(func() {
		key, err := vsjwt.GenerateRSAKeyPair()
		if err != nil {
			tm.err = fmt.Errorf("failed to generate RSA key for Cognito Identity tokens: %w", err)
			return
		}
		tm.privateKey = key
		tm.keyID = uuid.New().String()
	})
}

// generateOpenIdToken creates an OpenID Connect JWT for a Cognito identity.
// The token follows the format AWS Cognito Identity uses: signed with RS256,
// with iss, sub, aud, iat, exp claims.
func (tm *tokenManager) generateOpenIdToken(identityID, poolID string, expiresIn int64, amr []string) (string, error) {
	tm.init()
	if tm.err != nil {
		return "", tm.err
	}

	now := time.Now().UTC()
	claims := jwt.MapClaims{
		"iss": tm.issuer,
		"sub": identityID,
		"aud": poolID,
		"iat": now.Unix(),
		"exp": now.Add(time.Duration(expiresIn) * time.Second).Unix(),
		"jti": uuid.New().String(),
	}
	if len(amr) > 0 {
		claims["amr"] = amr
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	token.Header["kid"] = tm.keyID

	return token.SignedString(tm.privateKey)
}
