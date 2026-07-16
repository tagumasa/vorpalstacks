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
	"encoding/json"

	"github.com/golang-jwt/jwt/v5"
)

// JWTUser defines the interface for user information required for token generation.
// Implement this interface to provide user data when generating JWT tokens.
type JWTUser interface {
	GetID() string
	GetUsername() string
	GetGroups() []string
	GetEmail() string
	GetCustomClaims() map[string]interface{}
}

// CognitoClaims represents the claims structure compatible with AWS Cognito tokens.
// It includes both standard JWT claims and Cognito-specific claims.
type CognitoClaims struct {
	jwt.RegisteredClaims
	Username      string   `json:"cognito:username,omitempty"`
	Email         string   `json:"email,omitempty"`
	EmailVerified bool     `json:"email_verified,omitempty"`
	Groups        []string `json:"cognito:groups,omitempty"`
	ClientID      string   `json:"client_id,omitempty"`
	Scope         string   `json:"scope,omitempty"`
	TokenUse      string   `json:"token_use,omitempty"`
	AuthTime      int64    `json:"auth_time,omitempty"`
	Jti           string   `json:"jti,omitempty"`
	custom        map[string]interface{}
}

// cognitoClaimsAlias is used to avoid infinite recursion in UnmarshalJSON.
type cognitoClaimsAlias CognitoClaims

// UnmarshalJSON implements json.Unmarshaler for custom claims support.
func (c *CognitoClaims) UnmarshalJSON(data []byte) error {
	var alias cognitoClaimsAlias
	if err := json.Unmarshal(data, &alias); err != nil {
		return err
	}

	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	knownFields := map[string]bool{
		"iss": true, "sub": true, "aud": true, "exp": true, "nbf": true,
		"iat": true, "jti": true, "cognito:username": true, "email": true,
		"email_verified": true, "cognito:groups": true, "client_id": true,
		"scope": true, "token_use": true, "auth_time": true,
	}

	custom := make(map[string]interface{})
	for k, v := range raw {
		if !knownFields[k] {
			custom[k] = v
		}
	}

	*c = CognitoClaims(alias)
	if len(custom) > 0 {
		c.custom = custom
	}
	return nil
}

// MarshalJSON implements json.Marshaler for custom claims support.
func (c *CognitoClaims) MarshalJSON() ([]byte, error) {
	type Alias CognitoClaims
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(c),
	}

	data, err := json.Marshal(aux)
	if err != nil {
		return nil, err
	}

	if len(c.custom) == 0 {
		return data, nil
	}

	var merged map[string]interface{}
	if err := json.Unmarshal(data, &merged); err != nil {
		return nil, err
	}

	for k, v := range c.custom {
		merged[k] = v
	}
	return json.Marshal(merged)
}

// GetCustomClaim returns a custom claim value by key.
func (c *CognitoClaims) GetCustomClaim(key string) (interface{}, bool) {
	if c.custom == nil {
		return nil, false
	}
	v, ok := c.custom[key]
	return v, ok
}

// GetCustomClaimString returns a custom claim as string.
func (c *CognitoClaims) GetCustomClaimString(key string) string {
	if v, ok := c.GetCustomClaim(key); ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

// HasAudience checks if the given audience is in the token's audience claim.
func (c *CognitoClaims) HasAudience(aud string) bool {
	for _, a := range c.Audience {
		if a == aud {
			return true
		}
	}
	return false
}

// GetTokenUse returns the token type (access, id, or refresh).
func (c *CognitoClaims) GetTokenUse() string {
	return c.TokenUse
}

// IsAccessToken returns true if the token is an access token.
func (c *CognitoClaims) IsAccessToken() bool {
	return c.TokenUse == "access"
}

// IsIDToken returns true if the token is an ID token.
func (c *CognitoClaims) IsIDToken() bool {
	return c.TokenUse == "id"
}
