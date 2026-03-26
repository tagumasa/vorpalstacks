package admin_auth

// RefreshTokenStoreInterface defines operations for managing refresh tokens.
type RefreshTokenStoreInterface interface {
	Create(token *RefreshToken) error
	Get(token string) (*RefreshToken, error)
	Delete(token string) error
}

var _ RefreshTokenStoreInterface = (*RefreshTokenStore)(nil)
