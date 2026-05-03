// Package admin_auth provides administrative authentication service operations for vorpalstacks.
//
// admin_auth is an IAM sub-service and directly accesses the IAM store
// (internal/store/aws/iam) for login profile and access key management.
// The root user's credentials (login profile and access keys) are stored
// in the IAM store with the special user name iam.RootUserName.
package admin_auth

import (
	"fmt"
	"time"

	"vorpalstacks/internal/core/storage"
	adminauthstore "vorpalstacks/internal/store/aws/admin_auth"
	iamstore "vorpalstacks/internal/store/aws/iam"
	"vorpalstacks/pkg/vsjwt"
)

// Constants for admin authentication token durations.
const (
	// DefaultClientID is the default client ID for admin authentication.
	DefaultClientID = "admin-console"
	// AccessTokenDurationSec is the default access token duration in seconds.
	AccessTokenDurationSec = int64(3600)
	// IDTokenDurationSec is the default ID token duration in seconds.
	IDTokenDurationSec = int64(3600)
	// RefreshTokenDurationSec is the default refresh token duration in seconds.
	RefreshTokenDurationSec = int64(2592000)
)

// Service provides admin authentication operations.
type AdminAuthService struct {
	passwordVerifier  PasswordVerifier
	loginProfileCheck LoginProfileChecker
	userReader        UserReader
	passwordCreator   PasswordCreator
	accessKeyCreator  AccessKeyCreator
	groupReader       GroupMembershipReader
	policyReader      AttachedPolicyReader
	tokenGenerator    TokenGenerator
	refreshTokenStore *adminauthstore.RefreshTokenStore
	accountID         string
}

// NewAdminAuthService creates a new admin auth service.
//
// Parameters:
//   - store: The storage instance
//   - jwtManager: The JWT manager for token generation
//   - accountID: The AWS account ID
//
// Returns:
//   - *AdminAuthService: A new admin auth service instance
func NewAdminAuthService(
	store storage.BasicStorage,
	jwtManager *vsjwt.Manager,
	accountID string,
) *AdminAuthService {
	userStore := iamstore.NewUserStore(store, accountID)
	loginProfileStore := iamstore.NewLoginProfileStore(store)
	accessKeyStore := iamstore.NewAccessKeyStore(store)

	return &AdminAuthService{
		passwordVerifier:  loginProfileStore,
		loginProfileCheck: loginProfileStore,
		userReader:        userStore,
		passwordCreator:   loginProfileStore,
		accessKeyCreator:  accessKeyStore,
		groupReader:       iamstore.NewUserGroupStore(store),
		policyReader:      iamstore.NewAttachedPolicyStore(store),
		tokenGenerator:    jwtManager,
		refreshTokenStore: adminauthstore.NewRefreshTokenStore(store),
		accountID:         accountID,
	}
}

// RootUserAdapter adapts root user data to the JWTUser interface for token generation.
// Root users are not IAM users, so we construct a minimal adapter.
type RootUserAdapter struct {
	accountID string
}

// NewRootUserAdapter creates a new RootUserAdapter.
func NewRootUserAdapter(accountID string) *RootUserAdapter {
	return &RootUserAdapter{accountID: accountID}
}

// GetID returns the root user identifier.
func (u *RootUserAdapter) GetID() string {
	return iamstore.RootUserName
}

// GetUsername returns the root user name.
func (u *RootUserAdapter) GetUsername() string {
	return iamstore.RootUserName
}

// GetEmail returns an empty email for root user.
func (u *RootUserAdapter) GetEmail() string {
	return ""
}

// GetGroups returns empty groups for root user.
func (u *RootUserAdapter) GetGroups() []string {
	return nil
}

// GetCustomClaims returns custom claims for the root user JWT token.
func (u *RootUserAdapter) GetCustomClaims() map[string]interface{} {
	return map[string]interface{}{
		"arn":        fmt.Sprintf("arn:aws:iam::%s:root", u.accountID),
		"user_id":    iamstore.RootUserName,
		"account_id": u.accountID,
		"is_root":    true,
	}
}

var _ vsjwt.JWTUser = (*RootUserAdapter)(nil)

func (s *AdminAuthService) saveRefreshToken(token, username string) error {
	rt := &adminauthstore.RefreshToken{
		Token:     token,
		Username:  username,
		ClientID:  DefaultClientID,
		Expires:   time.Now().Add(time.Duration(RefreshTokenDurationSec) * time.Second),
		CreatedAt: time.Now(),
	}
	return s.refreshTokenStore.Create(rt)
}

func (s *AdminAuthService) getRefreshToken(token string) (*adminauthstore.RefreshToken, error) {
	return s.refreshTokenStore.Get(token)
}

func (s *AdminAuthService) deleteRefreshToken(token string) error {
	return s.refreshTokenStore.Delete(token)
}
