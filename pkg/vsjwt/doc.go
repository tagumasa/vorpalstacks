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

// Package vsjwt provides a lightweight JWT implementation with RSA-256 signing
// that is compatible with AWS Cognito token format.
//
// The package offers:
//   - JWT token generation (Access, ID, and Refresh tokens)
//   - Token validation with RSA public key
//   - Issuer and Audience validation
//   - Clock skew tolerance for token validation
//   - Custom claims support
//   - JWKS endpoint support for OpenID Connect discovery
//   - RSA key pair generation and PEM encoding/decoding
//   - Cognito-compatible claims structure
//
// # Quick Start
//
// Generating tokens:
//
//	privateKey, _ := vsjwt.GenerateRSAKeyPair()
//	manager := vsjwt.NewManager(privateKey, "key-id-123", "https://my-issuer.com")
//
//	user := &MyUser{id: "user-123", username: "alice"}
//	accessToken, _ := manager.GenerateAccessToken(user, "client-id", 3600)
//	idToken, _ := manager.GenerateIDToken(user, "client-id", 3600)
//
// With custom claims:
//
//	custom := map[string]interface{}{"role": "admin"}
//	token, _ := manager.GenerateAccessTokenWithClaims(user, "client-id", 3600, custom)
//
// Validating tokens:
//
//	publicKey := &privateKey.PublicKey
//	validator := vsjwt.NewManagerWithPublicKey(publicKey, "key-id-123", "https://my-issuer.com")
//	claims, _ := validator.ValidateToken(tokenString)
//
// With audience check:
//
//	claims, _ := validator.ValidateTokenWithAudience(tokenString, "expected-client")
//
// With clock skew tolerance:
//
//	manager := vsjwt.NewManager(privateKey, "key-id", issuer, vsjwt.WithClockSkew(time.Minute))
//
// JWKS for OpenID Connect:
//
//	jwks := manager.GetJWKS()
//	// Serve as JSON at /.well-known/jwks.json
//
// # JWTUser Interface
//
// Implement the JWTUser interface to provide user information for token generation:
//
//	type MyUser struct {
//	    id       string
//	    username string
//	    groups   []string
//	    email    string
//	}
//
//	func (u *MyUser) GetID() string       { return u.id }
//	func (u *MyUser) GetUsername() string { return u.username }
//	func (u *MyUser) GetGroups() []string { return u.groups }
//	func (u *MyUser) GetEmail() string    { return u.email }
package vsjwt
