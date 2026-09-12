/*
 * Copyright 2026 Vorpalstacks
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package vsjwt

import (
	"crypto/x509"
	"encoding/pem"
	"errors"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type testUser struct {
	id            string
	username      string
	groups        []string
	email         string
	emailVerified bool
	customClaims  map[string]interface{}
}

func (u *testUser) GetID() string                           { return u.id }
func (u *testUser) GetUsername() string                     { return u.username }
func (u *testUser) GetGroups() []string                     { return u.groups }
func (u *testUser) GetEmail() string                        { return u.email }
func (u *testUser) GetEmailVerified() bool                  { return u.emailVerified }
func (u *testUser) GetCustomClaims() map[string]interface{} { return u.customClaims }

func TestGenerateAndValidateAccessToken(t *testing.T) {
	privateKey, err := GenerateRSAKeyPair()
	if err != nil {
		t.Fatalf("failed to generate key: %v", err)
	}

	issuer := "https://cognito-idp.us-east-1.amazonaws.com/us-east-1_test"
	manager := mustNewManager(privateKey, "test-key-id", issuer)

	user := &testUser{
		id:       "user-123",
		username: "alice",
		groups:   []string{"admins"},
		email:    "alice@example.com",
	}

	token, err := manager.GenerateAccessToken(user, "client-123", 3600)
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	claims, err := manager.ValidateToken(token)
	if err != nil {
		t.Fatalf("failed to validate token: %v", err)
	}

	if claims.Subject != user.id {
		t.Errorf("expected subject %s, got %s", user.id, claims.Subject)
	}
	if claims.Username != user.username {
		t.Errorf("expected username %s, got %s", user.username, claims.Username)
	}
	if claims.TokenUse != "access" {
		t.Error("expected access token")
	}
}

func TestGenerateAndValidateIDToken(t *testing.T) {
	privateKey, err := GenerateRSAKeyPair()
	if err != nil {
		t.Fatalf("failed to generate key: %v", err)
	}

	issuer := "https://cognito-idp.us-east-1.amazonaws.com/us-east-1_test"
	manager := mustNewManager(privateKey, "test-key-id", issuer)

	user := &testUser{
		id:       "user-123",
		username: "bob",
		groups:   []string{"users"},
		email:    "bob@example.com",
	}

	token, err := manager.GenerateIDToken(user, "client-456", 3600)
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	claims, err := manager.ValidateTokenWithAudience(token, "client-456")
	if err != nil {
		t.Fatalf("failed to validate token: %v", err)
	}

	if claims.Subject != user.id {
		t.Errorf("expected subject %s, got %s", user.id, claims.Subject)
	}
	if claims.TokenUse != "id" {
		t.Error("expected ID token")
	}
	if !claims.HasAudience("client-456") {
		t.Error("expected audience client-456")
	}
}

// The ID token's email_verified claim follows the principal's email
// state rather than defaulting to true: a verified email yields true,
// an unverified one yields false (serialised as the claim's absence
// under the omitempty encoding).
func TestGenerateIDTokenEmailVerifiedDerivation(t *testing.T) {
	privateKey, err := GenerateRSAKeyPair()
	if err != nil {
		t.Fatalf("failed to generate key: %v", err)
	}
	manager := mustNewManager(privateKey, "test-key-id", "https://cognito-idp.us-east-1.amazonaws.com/us-east-1_test")

	verified := &testUser{id: "u1", username: "verified", email: "v@example.com", emailVerified: true}
	token, err := manager.GenerateIDToken(verified, "client", 3600)
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}
	claims, err := manager.ValidateToken(token)
	if err != nil {
		t.Fatalf("failed to validate token: %v", err)
	}
	if !claims.EmailVerified {
		t.Error("verified principal's ID token lost email_verified")
	}

	unverified := &testUser{id: "u2", username: "unverified", email: "u@example.com", emailVerified: false}
	token, err = manager.GenerateIDToken(unverified, "client", 3600)
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}
	claims, err = manager.ValidateToken(token)
	if err != nil {
		t.Fatalf("failed to validate token: %v", err)
	}
	if claims.EmailVerified {
		t.Error("unverified principal's ID token claimed email_verified")
	}
}

func TestValidateWithWrongAudience(t *testing.T) {
	privateKey, err := GenerateRSAKeyPair()
	if err != nil {
		t.Fatalf("failed to generate key: %v", err)
	}

	issuer := "https://cognito-idp.us-east-1.amazonaws.com/us-east-1_test"
	manager := mustNewManager(privateKey, "test-key-id", issuer)

	user := &testUser{id: "user-1", username: "test", groups: nil, email: ""}

	token, _ := manager.GenerateIDToken(user, "client-correct", 3600)

	_, err = manager.ValidateTokenWithAudience(token, "client-wrong")
	if err != ErrInvalidAudience {
		t.Errorf("expected ErrInvalidAudience, got %v", err)
	}
}

func TestValidateWithWrongIssuer(t *testing.T) {
	privateKey, err := GenerateRSAKeyPair()
	if err != nil {
		t.Fatalf("failed to generate key: %v", err)
	}

	manager := mustNewManager(privateKey, "test-key-id", "https://correct-issuer.com")

	user := &testUser{id: "user-1", username: "test", groups: nil, email: ""}
	token, _ := manager.GenerateAccessToken(user, "client", 3600)

	wrongManager := mustNewManagerWithPublicKey(&privateKey.PublicKey, "test-key-id", "https://wrong-issuer.com")

	_, err = wrongManager.ValidateToken(token)
	if err != ErrInvalidIssuer {
		t.Errorf("expected ErrInvalidIssuer, got %v", err)
	}
}

func TestCustomClaims(t *testing.T) {
	privateKey, _ := GenerateRSAKeyPair()
	manager := mustNewManager(privateKey, "key-id", "https://issuer.com")

	user := &testUser{id: "user-1", username: "test", groups: nil, email: ""}
	custom := map[string]interface{}{
		"custom:role":   "admin",
		"custom:tenant": "acme",
	}

	token, err := manager.GenerateAccessTokenWithClaims(user, "client", 3600, custom)
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	claims, err := manager.ValidateToken(token)
	if err != nil {
		t.Fatalf("failed to validate token: %v", err)
	}

	if role := claims.GetCustomClaimString("custom:role"); role != "admin" {
		t.Errorf("expected custom:role=admin, got %s", role)
	}
	if tenant := claims.GetCustomClaimString("custom:tenant"); tenant != "acme" {
		t.Errorf("expected custom:tenant=acme, got %s", tenant)
	}
}

func TestClockSkew(t *testing.T) {
	privateKey, _ := GenerateRSAKeyPair()
	manager := mustNewManager(privateKey, "key-id", "https://issuer.com", WithClockSkew(5*time.Minute))

	user := &testUser{id: "user-1", username: "test", groups: nil, email: ""}
	token, _ := manager.GenerateAccessToken(user, "client", 3600)

	_, err := manager.ValidateToken(token)
	if err != nil {
		t.Errorf("token should be valid with clock skew: %v", err)
	}
}

func TestRefreshToken(t *testing.T) {
	privateKey, _ := GenerateRSAKeyPair()
	manager := mustNewManager(privateKey, "key-id", "https://issuer.com")

	token := manager.GenerateRefreshToken()
	if len(token) < 30 {
		t.Errorf("refresh token too short: %s", token)
	}
}

func TestJWKS(t *testing.T) {
	privateKey, _ := GenerateRSAKeyPair()
	manager := mustNewManager(privateKey, "my-key-id", "https://issuer.com")

	jwks := manager.GetJWKS()
	if len(jwks.Keys) != 1 {
		t.Fatalf("expected 1 key, got %d", len(jwks.Keys))
	}

	key := jwks.Keys[0]
	if key.Kid != "my-key-id" {
		t.Errorf("expected kid=my-key-id, got %s", key.Kid)
	}
	if key.Alg != "RS256" {
		t.Errorf("expected alg=RS256, got %s", key.Alg)
	}
	if key.Kty != "RSA" {
		t.Errorf("expected kty=RSA, got %s", key.Kty)
	}
}

func TestPEMEncoding(t *testing.T) {
	privateKey, _ := GenerateRSAKeyPair()

	privatePEM := EncodePrivateKeyToPEM(privateKey)
	if privatePEM == "" {
		t.Error("private PEM is empty")
	}

	publicPEM, err := EncodePublicKeyToPEM(&privateKey.PublicKey)
	if err != nil {
		t.Fatalf("failed to encode public key: %v", err)
	}
	if publicPEM == "" {
		t.Error("public PEM is empty")
	}

	decodedPrivate, err := DecodePrivateKeyFromPEM(privatePEM)
	if err != nil {
		t.Fatalf("failed to decode private key: %v", err)
	}
	if decodedPrivate.N.Cmp(privateKey.N) != 0 {
		t.Error("decoded private key doesn't match")
	}

	decodedPublic, err := DecodePublicKeyFromPEM(publicPEM)
	if err != nil {
		t.Fatalf("failed to decode public key: %v", err)
	}
	if decodedPublic.N.Cmp(privateKey.PublicKey.N) != 0 {
		t.Error("decoded public key doesn't match")
	}
}

func TestNoPrivateKey(t *testing.T) {
	privateKey, _ := GenerateRSAKeyPair()
	manager := mustNewManagerWithPublicKey(&privateKey.PublicKey, "key-id", "https://issuer.com")

	user := &testUser{id: "user-1", username: "test", groups: nil, email: ""}
	_, err := manager.GenerateAccessToken(user, "client", 3600)
	if err != ErrNoPrivateKey {
		t.Errorf("expected ErrNoPrivateKey, got %v", err)
	}
}

func TestNewManagerWithPublicKeyNilKey(t *testing.T) {
	_, err := NewManagerWithPublicKey(nil, "key-id", "https://issuer.com")
	if err != ErrNilPublicKey {
		t.Errorf("expected ErrNilPublicKey, got %v", err)
	}
}

// --- Regression: token without exp must be rejected ---

func TestValidateRejectsMissingExp(t *testing.T) {
	privateKey, _ := GenerateRSAKeyPair()
	issuer := "https://issuer.com"
	manager := mustNewManager(privateKey, "key-id", issuer)

	now := time.Now().UTC()
	claims := &CognitoClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:  "user-1",
			Issuer:   issuer,
			IssuedAt: jwt.NewNumericDate(now),
		},
		Username: "test",
		TokenUse: "access",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	token.Header["kid"] = "key-id"
	tokenString, err := token.SignedString(privateKey)
	if err != nil {
		t.Fatalf("failed to sign token: %v", err)
	}

	_, err = manager.ValidateToken(tokenString)
	if err == nil {
		t.Error("expected error for token without exp, got nil")
	}
}

// --- Regression: token signed with RS384 must be rejected ---

func TestValidateRejectsNonRS256(t *testing.T) {
	privateKey, _ := GenerateRSAKeyPair()
	issuer := "https://issuer.com"
	manager := mustNewManager(privateKey, "key-id", issuer)

	user := &testUser{id: "user-1", username: "test", groups: nil, email: ""}
	token, err := manager.GenerateAccessToken(user, "client", 3600)
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	pubManager := mustNewManagerWithPublicKey(&privateKey.PublicKey, "key-id", issuer)

	claims := &CognitoClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   "user-1",
			Issuer:    issuer,
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(time.Hour)),
		},
		Username: "test",
		TokenUse: "access",
	}
	rs384Token := jwt.NewWithClaims(jwt.SigningMethodRS384, claims)
	rs384Token.Header["kid"] = "key-id"
	rs384TokenString, _ := rs384Token.SignedString(privateKey)

	_, err = pubManager.ValidateToken(rs384TokenString)
	if err == nil {
		t.Error("expected error for RS384 token, got nil")
	}

	_ = token
}

// --- Regression: PKCS8 format PEM must be decodable ---

func TestDecodePKCS8PrivateKey(t *testing.T) {
	privateKey, _ := GenerateRSAKeyPair()

	der, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		t.Fatalf("failed to marshal PKCS8: %v", err)
	}

	pemStr := string(pem.EncodeToMemory(&pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: der,
	}))

	decoded, err := DecodePrivateKeyFromPEM(pemStr)
	if err != nil {
		t.Fatalf("failed to decode PKCS8 PEM: %v", err)
	}

	if decoded.N.Cmp(privateKey.N) != 0 {
		t.Error("decoded key does not match original")
	}
}

func TestDecodeInvalidPEM(t *testing.T) {
	_, err := DecodePrivateKeyFromPEM("not a PEM string")
	if err != ErrInvalidPEMBlock {
		t.Errorf("expected ErrInvalidPEMBlock, got %v", err)
	}
}

func TestGenerateAccessTokenWithScopeCarriesScopeClaim(t *testing.T) {
	privateKey, err := GenerateRSAKeyPair()
	if err != nil {
		t.Fatalf("failed to generate key: %v", err)
	}
	manager := mustNewManager(privateKey, "test-key-id", "https://cognito-idp.us-east-1.amazonaws.com/us-east-1_test")

	user := &testUser{id: "user-123", username: "alice"}
	token, err := manager.GenerateAccessTokenWithScope(user, "client-123", "openid profile", 3600)
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	claims, err := manager.ValidateTokenForUse(token, "access", "")
	if err != nil {
		t.Fatalf("failed to validate token: %v", err)
	}
	if claims.Scope != "openid profile" {
		t.Errorf("expected scope %q, got %q", "openid profile", claims.Scope)
	}
	if claims.TokenUse != "access" {
		t.Error("expected access token")
	}
}

func TestValidateTokenForUseRejectsTokenTypeConfusion(t *testing.T) {
	privateKey, err := GenerateRSAKeyPair()
	if err != nil {
		t.Fatalf("failed to generate key: %v", err)
	}
	manager := mustNewManager(privateKey, "test-key-id", "https://cognito-idp.us-east-1.amazonaws.com/us-east-1_test")

	user := &testUser{id: "user-123", username: "alice"}
	idToken, err := manager.GenerateIDToken(user, "client-123", 3600)
	if err != nil {
		t.Fatalf("failed to generate ID token: %v", err)
	}

	if _, err := manager.ValidateTokenForUse(idToken, "access", ""); !errors.Is(err, ErrUnexpectedTokenUse) {
		t.Fatalf("expected ErrUnexpectedTokenUse, got %v", err)
	}

	if _, err := manager.ValidateTokenForUse(idToken, "id", "client-123"); err != nil {
		t.Fatalf("expected ID token to validate for use id with matching audience: %v", err)
	}
	if _, err := manager.ValidateTokenForUse(idToken, "id", "other-client"); !errors.Is(err, ErrInvalidAudience) {
		t.Fatalf("expected ErrInvalidAudience, got %v", err)
	}
}
