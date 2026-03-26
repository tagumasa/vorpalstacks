// Package admin_auth provides administrator authentication service for vorpalstacks.
package admin_auth

import (
	iamstore "vorpalstacks/internal/store/aws/iam"
	"vorpalstacks/pkg/vsjwt"
)

// PasswordVerifier defines the interface for verifying user passwords.
// Implementations should return true if the password is correct for the given username.
type PasswordVerifier interface {
	VerifyPassword(username, password string) (bool, error)
}

// UserReader defines the interface for retrieving user information.
type UserReader interface {
	Get(userName string) (*iamstore.User, error)
}

// GroupMembershipReader defines the interface for retrieving group membership information.
type GroupMembershipReader interface {
	ListGroupsForUser(userName string) ([]string, error)
}

// AttachedPolicyReader defines the interface for retrieving attached policy information.
type AttachedPolicyReader interface {
	ListAttachedPolicies(principalType, principalName string) ([]string, error)
}

// TokenGenerator defines the interface for generating and validating authentication tokens.
type TokenGenerator interface {
	GenerateAccessToken(user vsjwt.JWTUser, clientID string, expiresIn int64) (string, error)
	GenerateIDToken(user vsjwt.JWTUser, clientID string, expiresIn int64) (string, error)
	GenerateRefreshToken() string
	ValidateToken(tokenString string) (*vsjwt.CognitoClaims, error)
}
