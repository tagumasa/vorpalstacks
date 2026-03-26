package sts

// Package sts provides STS (Security Token Service) data store implementations
// for vorpalstacks.

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"log"
	"time"

	"vorpalstacks/internal/core/storage"
)

// SessionStore manages STS session tokens.
type SessionStore struct {
	bucket storage.Bucket
}

// NewSessionStore creates a new SessionStore instance.
func NewSessionStore(store storage.BasicStorage, region string) *SessionStore {
	bucketName := "sts_sessions-" + region
	return &SessionStore{
		bucket: store.Bucket(bucketName),
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
			log.Printf("Failed to delete expired session: %v", err)
		}
		return nil, ErrSessionExpired
	}

	return &session, nil
}

// Delete removes an STS session.
func (s *SessionStore) Delete(sessionToken string) error {
	return s.bucket.Delete([]byte(sessionToken))
}

func generateSessionToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(bytes), nil
}

func generateAccessKeyId() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return "ASIA" + hex.EncodeToString(bytes)[:16], nil
}

func generateSecretAccessKey() (string, error) {
	bytes := make([]byte, 30)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(bytes), nil
}
