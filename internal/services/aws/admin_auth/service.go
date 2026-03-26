// Package admin_auth provides administrative authentication service operations for vorpalstacks.
package admin_auth

import (
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
type Service struct {
	passwordVerifier  PasswordVerifier
	userReader        UserReader
	groupReader       GroupMembershipReader
	policyReader      AttachedPolicyReader
	tokenGenerator    TokenGenerator
	refreshTokenStore *adminauthstore.RefreshTokenStore
	accountID         string
}

// NewService creates a new admin auth service.
//
// Parameters:
//   - store: The storage instance
//   - jwtManager: The JWT manager for token generation
//   - accountID: The AWS account ID
//
// Returns:
//   - *Service: A new admin auth service instance
func NewService(
	store storage.BasicStorage,
	jwtManager *vsjwt.Manager,
	accountID string,
) *Service {
	return &Service{
		passwordVerifier:  iamstore.NewLoginProfileStore(store),
		userReader:        iamstore.NewUserStore(store, accountID),
		groupReader:       iamstore.NewUserGroupStore(store),
		policyReader:      iamstore.NewAttachedPolicyStore(store),
		tokenGenerator:    jwtManager,
		refreshTokenStore: adminauthstore.NewRefreshTokenStore(store),
		accountID:         accountID,
	}
}

func (s *Service) saveRefreshToken(token, username string) error {
	rt := &adminauthstore.RefreshToken{
		Token:     token,
		Username:  username,
		ClientID:  DefaultClientID,
		Expires:   time.Now().Add(time.Duration(RefreshTokenDurationSec) * time.Second),
		CreatedAt: time.Now(),
	}
	return s.refreshTokenStore.Create(rt)
}

func (s *Service) getRefreshToken(token string) (*adminauthstore.RefreshToken, error) {
	return s.refreshTokenStore.Get(token)
}

func (s *Service) deleteRefreshToken(token string) error {
	return s.refreshTokenStore.Delete(token)
}
