# vsjwt

A lightweight JWT implementation with RSA-256 signing, compatible with AWS Cognito token format.

## Features

- JWT token generation (Access, ID, and Refresh tokens)
- Token validation with RSA public key
- **Issuer and Audience validation**
- **Clock skew tolerance** for token validation
- **Custom claims support**
- JWKS endpoint support for OpenID Connect discovery
- RSA key pair generation and PEM encoding/decoding
- Cognito-compatible claims structure
- Zero external storage dependencies

## Installation

```bash
go get github.com/vorpalstacks/vorpalstacks/pkg/vsjwt
```

## Quick Start

### Generating Tokens

```go
package main

import (
    "fmt"
    "time"
    
    "github.com/vorpalstacks/vorpalstacks/pkg/vsjwt"
)

// Implement JWTUser interface
type User struct {
    id       string
    username string
    groups   []string
    email    string
}

func (u *User) GetID() string       { return u.id }
func (u *User) GetUsername() string { return u.username }
func (u *User) GetGroups() []string { return u.groups }
func (u *User) GetEmail() string    { return u.email }

func main() {
    // Generate RSA key pair
    privateKey, err := vsjwt.GenerateRSAKeyPair()
    if err != nil {
        panic(err)
    }
    
    // Create manager with clock skew tolerance
    issuer := "https://cognito-idp.us-east-1.amazonaws.com/us-east-1_abc123"
    manager := vsjwt.NewManager(privateKey, "key-id-123", issuer, vsjwt.WithClockSkew(time.Minute))
    
    // Create user
    user := &User{
        id:       "user-123",
        username: "alice",
        groups:   []string{"admins", "developers"},
        email:    "alice@example.com",
    }
    
    // Generate tokens
    accessToken, _ := manager.GenerateAccessToken(user, "client-id", 3600)
    idToken, _ := manager.GenerateIDToken(user, "client-id", 3600)
    refreshToken := manager.GenerateRefreshToken()
    
    fmt.Println("Access Token:", accessToken)
    fmt.Println("ID Token:", idToken)
    fmt.Println("Refresh Token:", refreshToken)
}
```

### Validating Tokens

```go
// Create validator with public key only
publicKey := &privateKey.PublicKey
validator := vsjwt.NewManagerWithPublicKey(publicKey, "key-id-123", issuer)

// Basic validation (issuer is always checked)
claims, err := validator.ValidateToken(tokenString)
if err != nil {
    // Invalid token
    panic(err)
}

// With audience validation
claims, err = validator.ValidateTokenWithAudience(tokenString, "expected-client-id")
if err != nil {
    // Invalid token or wrong audience
    panic(err)
}

fmt.Println("User:", claims.Username)
fmt.Println("Groups:", claims.Groups)
fmt.Println("Token Use:", claims.TokenUse)
```

### Custom Claims

```go
// Generate token with custom claims
customClaims := map[string]interface{}{
    "custom:role":   "admin",
    "custom:tenant": "acme-corp",
}
token, _ := manager.GenerateAccessTokenWithClaims(user, "client-id", 3600, customClaims)

// Access custom claims after validation
claims, _ := manager.ValidateToken(token)
role := claims.GetCustomClaimString("custom:role")  // "admin"
tenant, ok := claims.GetCustomClaim("custom:tenant")  // "acme-corp", true
```

### JWKS for OpenID Connect

```go
// Get JWKS for /.well-known/jwks.json endpoint
jwks := manager.GetJWKS()

// Serve as JSON
jsonBytes, _ := json.Marshal(jwks)
fmt.Println(string(jsonBytes))
```

### Key Persistence

```go
// Encode keys to PEM for storage
privatePEM := vsjwt.EncodePrivateKeyToPEM(privateKey)
publicPEM, _ := vsjwt.EncodePublicKeyToPEM(&privateKey.PublicKey)

// Store in database or file...

// Decode keys from PEM
decodedPrivate, _ := vsjwt.DecodePrivateKeyFromPEM(privatePEM)
decodedPublic, _ := vsjwt.DecodePublicKeyFromPEM(publicPEM)
```

## API Reference

### Manager

```go
// Create manager with private key (can generate and validate tokens)
manager := vsjwt.NewManager(privateKey, keyID, issuer, opts...)

// Create manager with public key only (can only validate tokens)
validator := vsjwt.NewManagerWithPublicKey(publicKey, keyID, issuer, opts...)

// Options
vsjwt.WithClockSkew(time.Minute)  // Allow clock skew (default: 30 seconds)

// Generate tokens
accessToken, err := manager.GenerateAccessToken(user, clientID, expiresInSeconds)
idToken, err := manager.GenerateIDToken(user, clientID, expiresInSeconds)
refreshToken := manager.GenerateRefreshToken()

// Generate with custom claims
token, err := manager.GenerateAccessTokenWithClaims(user, clientID, expiresIn, customClaims)
token, err := manager.GenerateIDTokenWithClaims(user, clientID, expiresIn, customClaims)

// Validate token (checks issuer automatically)
claims, err := manager.ValidateToken(tokenString)

// Validate with audience check
claims, err := manager.ValidateTokenWithAudience(tokenString, expectedAudience)

// Get JWKS
jwks := manager.GetJWKS()
```

### Validation Errors

```go
claims, err := manager.ValidateTokenWithAudience(token, audience)
if errors.Is(err, vsjwt.ErrInvalidIssuer) {
    // Token issuer doesn't match
}
if errors.Is(err, vsjwt.ErrInvalidAudience) {
    // Token audience doesn't match
}
if errors.Is(err, vsjwt.ErrInvalidToken) {
    // Token is invalid (signature, expired, etc.)
}
```

### JWTUser Interface

```go
type JWTUser interface {
    GetID() string
    GetUsername() string
    GetGroups() []string
    GetEmail() string
}
```

### CognitoClaims

```go
type CognitoClaims struct {
    // Standard claims
    jwt.RegisteredClaims
    
    // Cognito-specific claims
    Username      string   `json:"cognito:username,omitempty"`
    Email         string   `json:"email,omitempty"`
    EmailVerified bool     `json:"email_verified,omitempty"`
    Groups        []string `json:"cognito:groups,omitempty"`
    ClientID      string   `json:"client_id,omitempty"`
    Scope         string   `json:"scope,omitempty"`
    TokenUse      string   `json:"token_use,omitempty"`  // "access" or "id"
    AuthTime      int64    `json:"auth_time,omitempty"`
    Jti           string   `json:"jti,omitempty"`
}

// Helper methods
claims.IsAccessToken()  // true if token_use == "access"
claims.IsIDToken()      // true if token_use == "id"
claims.HasAudience(aud) // check if audience is present

// Custom claims access
claims.GetCustomClaim(key)        // interface{}, bool
claims.GetCustomClaimString(key)  // string
```

### Key Functions

```go
// Generate RSA key pair
privateKey, err := vsjwt.GenerateRSAKeyPair()

// PEM encoding/decoding
privatePEM := vsjwt.EncodePrivateKeyToPEM(privateKey)
publicPEM, err := vsjwt.EncodePublicKeyToPEM(publicKey)

privateKey, err := vsjwt.DecodePrivateKeyFromPEM(pemString)
publicKey, err := vsjwt.DecodePublicKeyFromPEM(pemString)
```

## Token Types

| Token Type | Purpose | Claims |
|------------|---------|--------|
| Access Token | API authorization | Subject, Username, Groups, ClientID |
| ID Token | User identity | Subject, Username, Email, Groups, Audience |
| Refresh Token | Obtain new tokens | Opaque string (not a JWT) |

## Compatibility

This package is designed to be compatible with AWS Cognito tokens:

- Uses RS256 signing algorithm
- Includes Cognito-specific claims (`cognito:username`, `cognito:groups`)
- Supports JWKS format for OpenID Connect discovery
- Token format matches Cognito's structure

## License

Apache License 2.0
