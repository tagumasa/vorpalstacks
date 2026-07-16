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
	"testing"
	"time"
)

type testUser struct {
	id           string
	username     string
	groups       []string
	email        string
	customClaims map[string]interface{}
}

func (u *testUser) GetID() string                           { return u.id }
func (u *testUser) GetUsername() string                     { return u.username }
func (u *testUser) GetGroups() []string                     { return u.groups }
func (u *testUser) GetEmail() string                        { return u.email }
func (u *testUser) GetCustomClaims() map[string]interface{} { return u.customClaims }

func TestGenerateAndValidateAccessToken(t *testing.T) {
	privateKey, err := GenerateRSAKeyPair()
	if err != nil {
		t.Fatalf("failed to generate key: %v", err)
	}

	issuer := "https://cognito-idp.us-east-1.amazonaws.com/us-east-1_test"
	manager := NewManager(privateKey, "test-key-id", issuer)

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
	if !claims.IsAccessToken() {
		t.Error("expected access token")
	}
}

func TestGenerateAndValidateIDToken(t *testing.T) {
	privateKey, err := GenerateRSAKeyPair()
	if err != nil {
		t.Fatalf("failed to generate key: %v", err)
	}

	issuer := "https://cognito-idp.us-east-1.amazonaws.com/us-east-1_test"
	manager := NewManager(privateKey, "test-key-id", issuer)

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
	if !claims.IsIDToken() {
		t.Error("expected ID token")
	}
	if !claims.HasAudience("client-456") {
		t.Error("expected audience client-456")
	}
}

func TestValidateWithWrongAudience(t *testing.T) {
	privateKey, err := GenerateRSAKeyPair()
	if err != nil {
		t.Fatalf("failed to generate key: %v", err)
	}

	issuer := "https://cognito-idp.us-east-1.amazonaws.com/us-east-1_test"
	manager := NewManager(privateKey, "test-key-id", issuer)

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

	manager := NewManager(privateKey, "test-key-id", "https://correct-issuer.com")

	user := &testUser{id: "user-1", username: "test", groups: nil, email: ""}
	token, _ := manager.GenerateAccessToken(user, "client", 3600)

	wrongManager := NewManagerWithPublicKey(&privateKey.PublicKey, "test-key-id", "https://wrong-issuer.com")

	_, err = wrongManager.ValidateToken(token)
	if err != ErrInvalidIssuer {
		t.Errorf("expected ErrInvalidIssuer, got %v", err)
	}
}

func TestCustomClaims(t *testing.T) {
	privateKey, _ := GenerateRSAKeyPair()
	manager := NewManager(privateKey, "key-id", "https://issuer.com")

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
	manager := NewManager(privateKey, "key-id", "https://issuer.com", WithClockSkew(5*time.Minute))

	user := &testUser{id: "user-1", username: "test", groups: nil, email: ""}
	token, _ := manager.GenerateAccessToken(user, "client", 3600)

	_, err := manager.ValidateToken(token)
	if err != nil {
		t.Errorf("token should be valid with clock skew: %v", err)
	}
}

func TestRefreshToken(t *testing.T) {
	privateKey, _ := GenerateRSAKeyPair()
	manager := NewManager(privateKey, "key-id", "https://issuer.com")

	token := manager.GenerateRefreshToken()
	if len(token) < 30 {
		t.Errorf("refresh token too short: %s", token)
	}
}

func TestJWKS(t *testing.T) {
	privateKey, _ := GenerateRSAKeyPair()
	manager := NewManager(privateKey, "my-key-id", "https://issuer.com")

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
	manager := NewManagerWithPublicKey(&privateKey.PublicKey, "key-id", "https://issuer.com")

	user := &testUser{id: "user-1", username: "test", groups: nil, email: ""}
	_, err := manager.GenerateAccessToken(user, "client", 3600)
	if err != ErrNoPrivateKey {
		t.Errorf("expected ErrNoPrivateKey, got %v", err)
	}
}

func TestManagerGetters(t *testing.T) {
	privateKey, _ := GenerateRSAKeyPair()
	expectedIssuer := "https://issuer.com"
	expectedKeyID := "my-key-id"
	expectedClockSkew := 5 * time.Second

	manager := NewManager(privateKey, expectedKeyID, expectedIssuer, WithClockSkew(expectedClockSkew))

	if kid := manager.GetKeyID(); kid != expectedKeyID {
		t.Errorf("expected keyID %s, got %s", expectedKeyID, kid)
	}
	if iss := manager.GetIssuer(); iss != expectedIssuer {
		t.Errorf("expected issuer %s, got %s", expectedIssuer, iss)
	}
	if cs := manager.GetClockSkew(); cs != expectedClockSkew {
		t.Errorf("expected clockSkew %v, got %v", expectedClockSkew, cs)
	}
}

func TestGetTokenUse(t *testing.T) {
	privateKey, _ := GenerateRSAKeyPair()
	manager := NewManager(privateKey, "key-id", "https://issuer.com")

	user := &testUser{id: "user-1", username: "test", groups: nil, email: ""}

	accessToken, _ := manager.GenerateAccessToken(user, "client", 3600)
	accessClaims, _ := manager.ValidateToken(accessToken)
	if use := accessClaims.GetTokenUse(); use != "access" {
		t.Errorf("expected token use 'access', got %s", use)
	}

	idToken, _ := manager.GenerateIDToken(user, "client", 3600)
	idClaims, _ := manager.ValidateToken(idToken)
	if use := idClaims.GetTokenUse(); use != "id" {
		t.Errorf("expected token use 'id', got %s", use)
	}
}

func TestGetJWKSMap(t *testing.T) {
	privateKey, _ := GenerateRSAKeyPair()
	manager := NewManager(privateKey, "my-key-id", "https://issuer.com")

	jwksMap := manager.GetJWKSMap()
	keys, ok := jwksMap["keys"].([]map[string]interface{})
	if !ok || len(keys) != 1 {
		t.Fatalf("expected 1 key in map, got %v", keys)
	}

	key := keys[0]
	if key["kid"] != "my-key-id" {
		t.Errorf("expected kid=my-key-id, got %v", key["kid"])
	}
	if key["alg"] != "RS256" {
		t.Errorf("expected alg=RS256, got %v", key["alg"])
	}
	if key["use"] != "sig" {
		t.Errorf("expected use=sig, got %v", key["use"])
	}
}
