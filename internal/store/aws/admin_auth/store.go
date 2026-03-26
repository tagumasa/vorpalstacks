// Package admin_auth provides admin authentication data store implementations for vorpalstacks.
package admin_auth

import (
	"time"

	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/store/aws/common"
)

// RefreshToken represents a refresh token for authentication.
type RefreshToken struct {
	Token     string    `json:"token"`
	Username  string    `json:"username"`
	ClientID  string    `json:"clientId"`
	Scope     string    `json:"scope"`
	Expires   time.Time `json:"expires"`
	CreatedAt time.Time `json:"createdAt"`
}

// RefreshTokenStore provides storage for refresh tokens.
type RefreshTokenStore struct {
	*common.BaseStore
}

// NewRefreshTokenStore creates a new refresh token store.
func NewRefreshTokenStore(store storage.BasicStorage) *RefreshTokenStore {
	return &RefreshTokenStore{
		BaseStore: common.NewBaseStore(store.Bucket("admin-refreshtokens"), "admin-refreshtokens"),
	}
}

// Create stores a new refresh token.
func (s *RefreshTokenStore) Create(token *RefreshToken) error {
	return s.Put(token.Token, token)
}

// Get retrieves a refresh token by its token value.
func (s *RefreshTokenStore) Get(token string) (*RefreshToken, error) {
	var rt RefreshToken
	if err := s.BaseStore.Get(token, &rt); err != nil {
		return nil, ErrTokenNotFound
	}
	if time.Now().After(rt.Expires) {
		return nil, ErrTokenExpired
	}
	return &rt, nil
}

// Delete removes a refresh token from the store.
func (s *RefreshTokenStore) Delete(token string) error {
	return s.BaseStore.Delete(token)
}
