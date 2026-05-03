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

// LoginProfileChecker defines the interface for checking whether a login profile exists.
// Used to determine if the root user has been initialised.
type LoginProfileChecker interface {
	Exists(userName string) bool
}

// UserReader defines the interface for retrieving user information.
type UserReader interface {
	Get(userName string) (*iamstore.User, error)
}

// PasswordCreator defines the interface for creating login profiles with passwords.
type PasswordCreator interface {
	Create(userName, password string, passwordResetRequired bool) (*iamstore.LoginProfile, error)
}

// GroupMembershipReader defines the interface for retrieving group membership information.
type GroupMembershipReader interface {
	ListGroupsForUser(userName string) ([]string, error)
}

// AttachedPolicyReader defines the interface for retrieving attached policy information.
type AttachedPolicyReader interface {
	ListAttachedPolicies(principalType, principalName string) ([]string, error)
}

// AccessKeyCreator defines the interface for creating IAM access keys.
type AccessKeyCreator interface {
	Create(userName string) (*iamstore.AccessKey, error)
}

// TokenGenerator defines the interface for generating and validating authentication tokens.
type TokenGenerator interface {
	GenerateAccessToken(user vsjwt.JWTUser, clientID string, expiresIn int64) (string, error)
	GenerateIDToken(user vsjwt.JWTUser, clientID string, expiresIn int64) (string, error)
	GenerateRefreshToken() string
	ValidateToken(tokenString string) (*vsjwt.CognitoClaims, error)
}
