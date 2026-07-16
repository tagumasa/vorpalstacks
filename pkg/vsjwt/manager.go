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
	"crypto/rsa"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const defaultClockSkew = 30 * time.Second

// Manager handles JWT token generation and validation.
type Manager struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
	keyID      string
	issuer     string
	clockSkew  time.Duration
}

// ManagerOption configures a Manager.
type ManagerOption func(*Manager)

// WithClockSkew sets the allowed clock skew for token validation.
func WithClockSkew(d time.Duration) ManagerOption {
	return func(m *Manager) {
		m.clockSkew = d
	}
}

// NewManager creates a new JWT Manager with a private key for token generation and validation.
// The issuer is typically in the format: https://cognito-idp.{region}.amazonaws.com/{userPoolID}
func NewManager(privateKey *rsa.PrivateKey, keyID, issuer string, opts ...ManagerOption) *Manager {
	m := &Manager{
		privateKey: privateKey,
		publicKey:  &privateKey.PublicKey,
		keyID:      keyID,
		issuer:     issuer,
		clockSkew:  defaultClockSkew,
	}
	for _, opt := range opts {
		opt(m)
	}
	return m
}

// NewManagerWithPublicKey creates a new JWT Manager with only a public key for token validation.
// This manager cannot generate tokens, only validate them.
func NewManagerWithPublicKey(publicKey *rsa.PublicKey, keyID, issuer string, opts ...ManagerOption) *Manager {
	m := &Manager{
		publicKey: publicKey,
		keyID:     keyID,
		issuer:    issuer,
		clockSkew: defaultClockSkew,
	}
	for _, opt := range opts {
		opt(m)
	}
	return m
}

// GenerateAccessToken generates an access token for the given user.
// Access tokens are used to authorize API requests.
func (m *Manager) GenerateAccessToken(user JWTUser, clientID string, expiresIn int64) (string, error) {
	return m.GenerateAccessTokenWithClaims(user, clientID, expiresIn, nil)
}

// GenerateAccessTokenWithClaims generates an access token with additional custom claims.
func (m *Manager) GenerateAccessTokenWithClaims(user JWTUser, clientID string, expiresIn int64, customClaims map[string]interface{}) (string, error) {
	if m.privateKey == nil {
		return "", ErrNoPrivateKey
	}

	mergedClaims := make(map[string]interface{})
	if userClaims := user.GetCustomClaims(); userClaims != nil {
		for k, v := range userClaims {
			mergedClaims[k] = v
		}
	}
	for k, v := range customClaims {
		mergedClaims[k] = v
	}

	now := time.Now().UTC()
	claims := &CognitoClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   user.GetID(),
			Issuer:    m.issuer,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(expiresIn) * time.Second)),
			ID:        uuid.New().String(),
		},
		Username: user.GetUsername(),
		Groups:   user.GetGroups(),
		ClientID: clientID,
		TokenUse: "access",
		AuthTime: now.Unix(),
		custom:   mergedClaims,
	}

	if email := user.GetEmail(); email != "" {
		claims.Email = email
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	token.Header["kid"] = m.keyID

	return token.SignedString(m.privateKey)
}

// GenerateIDToken generates an ID token for the given user.
// ID tokens contain user identity information.
func (m *Manager) GenerateIDToken(user JWTUser, clientID string, expiresIn int64) (string, error) {
	return m.GenerateIDTokenWithClaims(user, clientID, expiresIn, nil)
}

// GenerateIDTokenWithClaims generates an ID token with additional custom claims.
func (m *Manager) GenerateIDTokenWithClaims(user JWTUser, clientID string, expiresIn int64, customClaims map[string]interface{}) (string, error) {
	if m.privateKey == nil {
		return "", ErrNoPrivateKey
	}

	mergedClaims := make(map[string]interface{})
	if userClaims := user.GetCustomClaims(); userClaims != nil {
		for k, v := range userClaims {
			mergedClaims[k] = v
		}
	}
	for k, v := range customClaims {
		mergedClaims[k] = v
	}

	now := time.Now().UTC()
	claims := &CognitoClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   user.GetID(),
			Issuer:    m.issuer,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(expiresIn) * time.Second)),
			ID:        uuid.New().String(),
			Audience:  jwt.ClaimStrings{clientID},
		},
		Username:      user.GetUsername(),
		Groups:        user.GetGroups(),
		TokenUse:      "id",
		AuthTime:      now.Unix(),
		EmailVerified: true,
		custom:        mergedClaims,
	}

	if email := user.GetEmail(); email != "" {
		claims.Email = email
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	token.Header["kid"] = m.keyID

	return token.SignedString(m.privateKey)
}

// GenerateRefreshToken generates a refresh token.
// Refresh tokens are opaque strings (not JWTs) used to obtain new tokens.
func (m *Manager) GenerateRefreshToken() string {
	return uuid.New().String() + uuid.New().String()
}

// ValidateToken validates a JWT token and returns the claims.
// It verifies the signature, issuer, and expiration with clock skew tolerance.
func (m *Manager) ValidateToken(tokenString string) (*CognitoClaims, error) {
	return m.ValidateTokenWithAudience(tokenString, "")
}

// ValidateTokenWithAudience validates a JWT token with audience checking.
// If expectedAudience is non-empty, the token's aud claim must match.
func (m *Manager) ValidateTokenWithAudience(tokenString, expectedAudience string) (*CognitoClaims, error) {
	claims := &CognitoClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, ErrUnexpectedMethod
		}
		return m.publicKey, nil
	}, jwt.WithLeeway(m.clockSkew))

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, ErrInvalidToken
	}

	if err := m.validateClaims(claims, expectedAudience); err != nil {
		return nil, err
	}

	return claims, nil
}

func (m *Manager) validateClaims(claims *CognitoClaims, expectedAudience string) error {
	if claims.Issuer != m.issuer {
		return ErrInvalidIssuer
	}

	if claims.Subject == "" {
		return ErrInvalidSubject
	}

	if expectedAudience != "" {
		if !claims.HasAudience(expectedAudience) {
			return ErrInvalidAudience
		}
	}

	return nil
}

// GetKeyID returns the key ID (kid) used in the JWT header.
func (m *Manager) GetKeyID() string {
	return m.keyID
}

// GetIssuer returns the issuer (iss) claim value.
func (m *Manager) GetIssuer() string {
	return m.issuer
}

// GetClockSkew returns the configured clock skew duration.
func (m *Manager) GetClockSkew() time.Duration {
	return m.clockSkew
}
