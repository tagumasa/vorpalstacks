package cognitoidentity

import (
	"crypto/rsa"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"vorpalstacks/pkg/vsjwt"
)

// cognitoIdentityIssuer returns the Cognito Identity OIDC issuer hostname
// for the given AWS partition. The standard partition uses a regionless
// hostname; China and GovCloud partitions include the region.
func cognitoIdentityIssuer(region string) string {
	if strings.HasPrefix(region, "cn-") {
		return "cognito-identity." + region + ".amazonaws.com.cn"
	}
	if strings.HasPrefix(region, "us-gov-") {
		return "cognito-identity." + region + ".amazonaws.com"
	}
	return "cognito-identity.amazonaws.com"
}

// tokenManager lazily generates an RSA key pair for signing Cognito Identity
// OpenID Connect tokens. The key is shared across all identity pools in the
// service instance and is regenerated on each process restart — acceptable
// because the tokens are short-lived (10 minutes for GetOpenIdToken,
// configurable up to 24 hours for GetOpenIdTokenForDeveloperIdentity).
type tokenManager struct {
	once       sync.Once
	privateKey *rsa.PrivateKey
	keyID      string
	issuer     string
	err        error
}

func newTokenManager(region string) *tokenManager {
	return &tokenManager{
		issuer: cognitoIdentityIssuer(region),
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
// with iss, sub, aud, iat, exp, jti claims. When principalTags are provided
// they are embedded as the cognito:principal_tags claim so that downstream
// STS AssumeRoleWithWebIdentity propagates them as session tags.
func (tm *tokenManager) generateOpenIdToken(identityID, poolID string, expiresIn int64, amr []string, principalTags map[string]string) (string, error) {
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
	if len(principalTags) > 0 {
		claims["cognito:principal_tags"] = principalTags
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	token.Header["kid"] = tm.keyID

	return token.SignedString(tm.privateKey)
}
