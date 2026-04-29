package sts

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"time"

	"vorpalstacks/internal/common/auth"
	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/core/storage"
)

const (
	stsSessionTokenSize    = 32
	stsAccessKeyIDSize     = 16
	stsSecretAccessKeySize = 30
)

// SessionStore manages STS session tokens.
type SessionStore struct {
	bucket          storage.Bucket
	accessKeyBucket storage.Bucket
}

// NewSessionStore creates a new SessionStore instance.
func NewSessionStore(store storage.BasicStorage, region string) *SessionStore {
	bucketName := "sts_sessions-" + region
	akBucketName := "sts_access_keys-" + region
	return &SessionStore{
		bucket:          store.Bucket(bucketName),
		accessKeyBucket: store.Bucket(akBucketName),
	}
}

// Create creates a new STS session.
func (s *SessionStore) Create(principalType, principalName, principalArn, roleArn, roleSessionName string, durationSeconds int) (*Session, error) {
	sessionToken, err := generateSessionToken()
	if err != nil {
		return nil, err
	}

	accessKeyId, err := generateAccessKeyId()
	if err != nil {
		return nil, err
	}

	secretAccessKey, err := generateSecretAccessKey()
	if err != nil {
		return nil, err
	}

	if durationSeconds == 0 {
		durationSeconds = 3600
	}

	session := &Session{
		SessionToken:    sessionToken,
		AccessKeyId:     accessKeyId,
		SecretAccessKey: secretAccessKey,
		Expiration:      time.Now().UTC().Add(time.Duration(durationSeconds) * time.Second),
		PrincipalArn:    principalArn,
		PrincipalType:   principalType,
		PrincipalName:   principalName,
		RoleArn:         roleArn,
		RoleSessionName: roleSessionName,
	}

	data, err := json.Marshal(session)
	if err != nil {
		return nil, err
	}

	if err := s.bucket.Put([]byte(sessionToken), data); err != nil {
		return nil, err
	}

	if err := s.accessKeyBucket.Put([]byte(accessKeyId), []byte(sessionToken)); err != nil {
		_ = s.bucket.Delete([]byte(sessionToken))
		return nil, err
	}

	return session, nil
}

// Get retrieves an STS session by session token.
func (s *SessionStore) Get(sessionToken string) (*Session, error) {
	data, err := s.bucket.Get([]byte(sessionToken))
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, ErrSessionNotFound
	}

	var session Session
	if err := json.Unmarshal(data, &session); err != nil {
		return nil, err
	}

	if session.Expiration.Before(time.Now().UTC()) {
		if err := s.bucket.Delete([]byte(sessionToken)); err != nil {
			logs.Error("Failed to delete expired session", logs.Err(err))
		}
		return nil, ErrSessionExpired
	}

	return &session, nil
}

// Delete removes an STS session.
func (s *SessionStore) Delete(sessionToken string) error {
	data, err := s.bucket.Get([]byte(sessionToken))
	if err == nil && data != nil {
		var session Session
		if json.Unmarshal(data, &session) == nil {
			_ = s.accessKeyBucket.Delete([]byte(session.AccessKeyId))
		}
	}
	return s.bucket.Delete([]byte(sessionToken))
}

// GetByAccessKeyId retrieves an STS session by access key ID.
func (s *SessionStore) GetByAccessKeyId(accessKeyId string) (*Session, error) {
	tokenBytes, err := s.accessKeyBucket.Get([]byte(accessKeyId))
	if err != nil {
		return nil, err
	}
	if tokenBytes == nil {
		return nil, ErrSessionNotFound
	}
	return s.Get(string(tokenBytes))
}

// ResolveSession implements auth.SessionResolver for STS session lookup by access key ID.
func (s *SessionStore) ResolveSession(accessKeyId string) (*auth.SessionCredentials, error) {
	session, err := s.GetByAccessKeyId(accessKeyId)
	if err != nil {
		return nil, err
	}
	return &auth.SessionCredentials{
		AccessKeyID:     session.AccessKeyId,
		SecretAccessKey: session.SecretAccessKey,
		SessionToken:    session.SessionToken,
	}, nil
}

func generateSessionToken() (string, error) {
	bytes := make([]byte, stsSessionTokenSize)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(bytes), nil
}

func generateAccessKeyId() (string, error) {
	bytes := make([]byte, stsAccessKeyIDSize)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return "ASIA" + hex.EncodeToString(bytes)[:16], nil
}

func generateSecretAccessKey() (string, error) {
	bytes := make([]byte, stsSecretAccessKeySize)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(bytes), nil
}
