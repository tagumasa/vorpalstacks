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
	"encoding/base64"
)

// JWKS represents a JSON Web Key Set for OpenID Connect discovery.
type JWKS struct {
	Keys []JWK `json:"keys"`
}

// JWK represents a single JSON Web Key.
type JWK struct {
	Alg string `json:"alg"`
	Kty string `json:"kty"`
	Use string `json:"use"`
	Kid string `json:"kid"`
	N   string `json:"n"`
	E   string `json:"e"`
}

// GetJWKS returns the JSON Web Key Set for this manager's public key.
// This can be served at the /.well-known/jwks.json endpoint for OIDC discovery.
func (m *Manager) GetJWKS() *JWKS {
	n := m.publicKey.N.Bytes()
	e := encodeExponent(m.publicKey.E)

	return &JWKS{
		Keys: []JWK{
			{
				Alg: "RS256",
				Kty: "RSA",
				Use: "sig",
				Kid: m.keyID,
				N:   base64.RawURLEncoding.EncodeToString(n),
				E:   base64.RawURLEncoding.EncodeToString(e),
			},
		},
	}
}

// GetJWKSMap returns the JWKS as a map for compatibility with existing code.
// Deprecated: Use GetJWKS() for typed access.
func (m *Manager) GetJWKSMap() map[string]interface{} {
	jwks := m.GetJWKS()
	return map[string]interface{}{
		"keys": []map[string]interface{}{
			{
				"alg": jwks.Keys[0].Alg,
				"kty": jwks.Keys[0].Kty,
				"use": jwks.Keys[0].Use,
				"kid": jwks.Keys[0].Kid,
				"n":   jwks.Keys[0].N,
				"e":   jwks.Keys[0].E,
			},
		},
	}
}

func encodeExponent(e int) []byte {
	buf := make([]byte, 4)
	for i := 0; i < 4; i++ {
		buf[3-i] = byte(e & 0xff)
		e >>= 8
	}
	for len(buf) > 1 && buf[0] == 0 {
		buf = buf[1:]
	}
	return buf
}
