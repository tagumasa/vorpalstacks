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

import "errors"

// ErrNoPrivateKey is returned when the token manager attempts to sign a token
// but has not been configured with a private key.
var ErrNoPrivateKey = errors.New("manager does not have a private key configured")

// ErrNilPublicKey is returned when NewManagerWithPublicKey is called with a nil
// public key, which would leave the manager unable to verify tokens.
var ErrNilPublicKey = errors.New("public key must not be nil")

// ErrInvalidToken is returned when a token fails validation due to being
// malformed, expired, or having an invalid signature.
var ErrInvalidToken = errors.New("invalid token")

// ErrInvalidPEMBlock is returned when a PEM-encoded key or certificate
// cannot be parsed correctly.
var ErrInvalidPEMBlock = errors.New("failed to parse PEM block")

// ErrUnexpectedMethod is returned when a token's signing method does not
// match the expected algorithm.
var ErrUnexpectedMethod = errors.New("unexpected signing method")

// ErrInvalidIssuer is returned when the token's issuer claim does not
// match the expected issuer.
var ErrInvalidIssuer = errors.New("invalid token issuer")

// ErrInvalidAudience is returned when the token's audience claim does not
// contain the expected recipient.
var ErrInvalidAudience = errors.New("invalid token audience")

// ErrInvalidSubject is returned when the token's subject claim does not
// match the expected subject.
var ErrInvalidSubject = errors.New("invalid token subject")

// ErrNotRSAPublicKey is returned when a PEM block decodes to a non-RSA key
// (e.g. EC, Ed25519) during public key import.
var ErrNotRSAPublicKey = errors.New("PEM block does not contain an RSA public key")

// ErrInvalidKeySize is returned when an RSA key has fewer than 2048 bits,
// which is below the minimum required for RS256 signing per AWS Cognito spec.
var ErrInvalidKeySize = errors.New("RSA key must be at least 2048 bits")
