package sts

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"os"
	"sync"
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

// delegatedTokenEntry stores the principal ARN associated with a trade-in token.
type delegatedTokenEntry struct {
	PrincipalArn string    `json:"principal_arn"`
	Expires      time.Time `json:"expires"`
}

// SessionStore manages STS session tokens.
type SessionStore struct {
	bucket          storage.Bucket
	accessKeyBucket storage.Bucket
	delegatedBucket storage.Bucket
	ctx             context.Context
	cancel          context.CancelFunc
	wg              sync.WaitGroup
}

// cleanupInterval is how often the background TTL sweeper runs.
const cleanupInterval = 5 * time.Minute

// NewSessionStore creates a new SessionStore instance and starts a
// background goroutine that periodically removes expired sessions,
// their access-key index entries, and expired delegated tokens.
func NewSessionStore(store storage.BasicStorage, region string) *SessionStore {
	bucketName := "sts_sessions-" + region
	akBucketName := "sts_access_keys-" + region
	delegatedBucketName := "sts_delegated_tokens-" + region
	ctx, cancel := context.WithCancel(context.Background())
	s := &SessionStore{
		bucket:          store.Bucket(bucketName),
		accessKeyBucket: store.Bucket(akBucketName),
		delegatedBucket: store.Bucket(delegatedBucketName),
		ctx:             ctx,
		cancel:          cancel,
	}
	s.wg.Add(1)
	go s.startCleanupSweeper(cleanupInterval)
	return s
}

// Close stops the background TTL sweeper goroutine and waits for it
// to exit. Safe to call multiple times.
func (s *SessionStore) Close() {
	if s.cancel != nil {
		s.cancel()
	}
	s.wg.Wait()
}

// startCleanupSweeper periodically removes expired sessions and
// delegated tokens until the store's context is cancelled.
func (s *SessionStore) startCleanupSweeper(interval time.Duration) {
	defer s.wg.Done()
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-s.ctx.Done():
			return
		case <-ticker.C:
			s.sweepExpired()
		}
	}
}

// sweepExpired iterates the session and delegated-token buckets,
// deleting entries whose Expiration is in the past.
func (s *SessionStore) sweepExpired() {
	now := time.Now().UTC()
	_ = s.bucket.ForEach(func(k, v []byte) error {
		var session Session
		if err := json.Unmarshal(v, &session); err != nil {
			return nil // skip corrupt entries
		}
		if session.Expiration.Before(now) {
			_ = s.bucket.Delete(k)
			_ = s.accessKeyBucket.Delete([]byte(session.AccessKeyId))
		}
		return nil
	})
	_ = s.delegatedBucket.ForEach(func(k, v []byte) error {
		var entry delegatedTokenEntry
		if err := json.Unmarshal(v, &entry); err != nil {
			return nil
		}
		if entry.Expires.Before(now) {
			_ = s.delegatedBucket.Delete(k)
		}
		return nil
	})
}

// SeedTestDelegatedTokens populates delegated tokens for test mode.
// This must be called explicitly by the server initialisation when TEST_MODE is enabled.
func (s *SessionStore) SeedTestDelegatedTokens() {
	testTokens := map[string]string{
		"dummy-trade-in-token-verify": "arn:aws:iam::123456789012:root",
	}
	expires := time.Now().UTC().Add(24 * time.Hour)
	for token, arn := range testTokens {
		_ = s.StoreDelegationToken(token, arn, expires)
	}
}

// Create creates a new STS session.
func (s *SessionStore) Create(params CreateSessionParams) (*Session, error) {
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

	durationSeconds := params.DurationSeconds
	if durationSeconds == 0 {
		durationSeconds = 3600
	}

	session := &Session{
		SessionToken:           sessionToken,
		AccessKeyId:            accessKeyId,
		SecretAccessKey:        secretAccessKey,
		Expiration:             time.Now().UTC().Add(time.Duration(durationSeconds) * time.Second),
		PrincipalArn:           params.PrincipalArn,
		PrincipalType:          params.PrincipalType,
		PrincipalName:          params.PrincipalName,
		RoleArn:                params.RoleArn,
		RoleSessionName:        params.RoleSessionName,
		SourceIdentity:         params.SourceIdentity,
		Tags:                   params.Tags,
		MultiFactorAuthPresent: params.MultiFactorAuthPresent,
		TransitiveTagKeys:      params.TransitiveTagKeys,
		Policy:                 params.Policy,
		PolicyArns:             params.PolicyArns,
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
		// Also clean up the access-key index entry to prevent orphaned
		// lookups that dereference a deleted session token.
		_ = s.accessKeyBucket.Delete([]byte(session.AccessKeyId))
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
		// Map the store-level sentinel to the cross-cutting auth
		// sentinel so that the auth middleware can detect expired
		// sessions and return ExpiredTokenException.
		if errors.Is(err, ErrSessionExpired) {
			return nil, auth.ErrSessionExpired
		}
		return nil, err
	}
	return &auth.SessionCredentials{
		AccessKeyID:     session.AccessKeyId,
		SecretAccessKey: session.SecretAccessKey,
		SessionToken:    session.SessionToken,
		PrincipalArn:    session.PrincipalArn,
		PrincipalType:   session.PrincipalType,
		Tags:            session.Tags,
		SourceIdentity:  session.SourceIdentity,
		Policy:          session.Policy,
		PolicyArns:      session.PolicyArns,
	}, nil
}

// StoreDelegationToken stores a trade-in token with its associated principal ARN.
func (s *SessionStore) StoreDelegationToken(token, principalArn string, expires time.Time) error {
	entry := &delegatedTokenEntry{
		PrincipalArn: principalArn,
		Expires:      expires,
	}
	data, err := json.Marshal(entry)
	if err != nil {
		return err
	}
	return s.delegatedBucket.Put([]byte(token), data)
}

// RedeemDelegationToken looks up a trade-in token and returns the
// associated principal ARN. Returns ErrDelegationTokenNotFound if the
// token does not exist, ErrDelegationTokenExpired if it has expired,
// or the underlying storage error for I/O failures.
func (s *SessionStore) RedeemDelegationToken(token string) (string, error) {
	data, err := s.delegatedBucket.Get([]byte(token))
	if err != nil {
		return "", err
	}
	if data == nil {
		return "", ErrDelegationTokenNotFound
	}

	var entry delegatedTokenEntry
	if err := json.Unmarshal(data, &entry); err != nil {
		return "", err
	}

	if entry.Expires.Before(time.Now().UTC()) {
		_ = s.delegatedBucket.Delete([]byte(token))
		return "", ErrDelegationTokenExpired
	}

	_ = s.delegatedBucket.Delete([]byte(token))
	if os.Getenv("TEST_MODE") == "true" {
		_ = s.StoreDelegationToken(token, entry.PrincipalArn, entry.Expires)
	}

	return entry.PrincipalArn, nil
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
