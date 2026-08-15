package cognitoidentityprovider

import (
	"time"
)

// CreateRefreshToken creates a new Cognito refresh token.
func (s *CognitoStore) CreateRefreshToken(token *RefreshToken) error {
	key := tokenKey(token.UserPoolID, token.UserID, token.Token)
	return s.refreshTokensStore.Put(key, token)
}

// GetRefreshToken retrieves a Cognito refresh token.
func (s *CognitoStore) GetRefreshToken(userPoolID, userID, token string) (*RefreshToken, error) {
	key := tokenKey(userPoolID, userID, token)
	var rt RefreshToken
	if err := s.refreshTokensStore.Get(key, &rt); err != nil {
		return nil, ErrTokenNotFound
	}
	if time.Now().After(rt.Expires) {
		_ = s.refreshTokensStore.Delete(key)
		return nil, ErrTokenExpired
	}
	return &rt, nil
}

// GetRefreshTokenByValue retrieves a Cognito refresh token by its token value.
func (s *CognitoStore) GetRefreshTokenByValue(token string) (*RefreshToken, error) {
	return findTokenByValue(s.refreshTokensStore, token, func(t *RefreshToken) string { return t.Token }, func(t *RefreshToken) time.Time { return t.Expires })
}

// DeleteRefreshToken deletes a Cognito refresh token.
func (s *CognitoStore) DeleteRefreshToken(userPoolID, userID, token string) error {
	key := tokenKey(userPoolID, userID, token)
	return s.refreshTokensStore.Delete(key)
}

// DeleteAllRefreshTokensForUser deletes all refresh tokens for a user.
func (s *CognitoStore) DeleteAllRefreshTokensForUser(userPoolID, userID string) error {
	prefix := userPoolID + "#" + userID + "#"
	return s.refreshTokensStore.ScanPrefix(prefix, func(key string, value []byte) error {
		return s.refreshTokensStore.Delete(key)
	})
}

// CreateIDToken creates a new Cognito ID token.
func (s *CognitoStore) CreateIDToken(token *IDToken) error {
	key := tokenKey(token.UserPoolID, token.UserID, token.Token)
	return s.idTokensStore.Put(key, token)
}

// GetIDToken retrieves a Cognito ID token.
func (s *CognitoStore) GetIDToken(userPoolID, userID, token string) (*IDToken, error) {
	key := tokenKey(userPoolID, userID, token)
	var it IDToken
	if err := s.idTokensStore.Get(key, &it); err != nil {
		return nil, ErrTokenNotFound
	}
	if time.Now().After(it.Expires) {
		_ = s.idTokensStore.Delete(key)
		return nil, ErrTokenExpired
	}
	return &it, nil
}

// GetIDTokenByValue retrieves a Cognito ID token by its token value.
func (s *CognitoStore) GetIDTokenByValue(token string) (*IDToken, error) {
	return findTokenByValue(s.idTokensStore, token, func(t *IDToken) string { return t.Token }, func(t *IDToken) time.Time { return t.Expires })
}

// DeleteIDToken deletes a Cognito ID token.
func (s *CognitoStore) DeleteIDToken(userPoolID, userID, token string) error {
	key := tokenKey(userPoolID, userID, token)
	return s.idTokensStore.Delete(key)
}

// CreateAccessToken creates a new Cognito access token.
func (s *CognitoStore) CreateAccessToken(token *AccessToken) error {
	key := tokenKey(token.UserPoolID, token.UserID, token.Token)
	return s.accessTokensStore.Put(key, token)
}

// GetAccessToken retrieves a Cognito access token.
func (s *CognitoStore) GetAccessToken(userPoolID, userID, token string) (*AccessToken, error) {
	key := tokenKey(userPoolID, userID, token)
	var at AccessToken
	if err := s.accessTokensStore.Get(key, &at); err != nil {
		return nil, ErrTokenNotFound
	}
	if time.Now().After(at.Expires) {
		_ = s.accessTokensStore.Delete(key)
		return nil, ErrTokenExpired
	}
	return &at, nil
}

// GetAccessTokenByValue retrieves a Cognito access token by its token value.
func (s *CognitoStore) GetAccessTokenByValue(token string) (*AccessToken, error) {
	return findTokenByValue(s.accessTokensStore, token, func(t *AccessToken) string { return t.Token }, func(t *AccessToken) time.Time { return t.Expires })
}

// DeleteAccessToken deletes a Cognito access token.
func (s *CognitoStore) DeleteAccessToken(userPoolID, userID, token string) error {
	key := tokenKey(userPoolID, userID, token)
	return s.accessTokensStore.Delete(key)
}

// DeleteUserTokens deletes all tokens for a user.
func (s *CognitoStore) DeleteUserTokens(userPoolID, userID string) error {
	prefix := userPoolID + "#" + userID + "#"
	var firstErr error

	if err := s.refreshTokensStore.ScanPrefix(prefix, func(key string, value []byte) error {
		return s.refreshTokensStore.Delete(key)
	}); err != nil && firstErr == nil {
		firstErr = err
	}

	if err := s.idTokensStore.ScanPrefix(prefix, func(key string, value []byte) error {
		return s.idTokensStore.Delete(key)
	}); err != nil && firstErr == nil {
		firstErr = err
	}

	if err := s.accessTokensStore.ScanPrefix(prefix, func(key string, value []byte) error {
		return s.accessTokensStore.Delete(key)
	}); err != nil && firstErr == nil {
		firstErr = err
	}

	return firstErr
}

// SaveChallengeSession saves a challenge session.
func (s *CognitoStore) SaveChallengeSession(session *ChallengeSession) error {
	return s.challengeSessionsStore.Put(session.SessionID, session)
}

// GetChallengeSession retrieves a challenge session by ID.
func (s *CognitoStore) GetChallengeSession(sessionID string) (*ChallengeSession, error) {
	var session ChallengeSession
	if err := s.challengeSessionsStore.Get(sessionID, &session); err != nil {
		return nil, err
	}
	return &session, nil
}

// DeleteChallengeSession deletes a challenge session.
func (s *CognitoStore) DeleteChallengeSession(sessionID string) error {
	return s.challengeSessionsStore.Delete(sessionID)
}
