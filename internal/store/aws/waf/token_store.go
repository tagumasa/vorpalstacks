package waf

import (
	"crypto/rand"
	"sync"

	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/store/aws/common"
)

const tokenSecretBucketName = "waf_token_secret"

// tokenSecretKey is the single storage key of the inspection plane's
// token signing secret.
const tokenSecretKey = "inspection"

// tokenSecretBytes is the length of a generated signing secret — the
// full width of an HMAC-SHA256 key.
const tokenSecretBytes = 32

// TokenStore persists the HMAC signing secret behind the inspection
// plane's aws-waf-token cookies. The secret is generated once and
// survives restarts, so tokens issued before a restart remain verifiable
// for their whole immunity window.
type TokenStore struct {
	*common.BaseStore
	mu     sync.Mutex
	cached []byte
}

// NewTokenStore creates the token-secret store over the given storage.
func NewTokenStore(store storage.BasicStorage) *TokenStore {
	return &TokenStore{
		BaseStore: common.NewBaseStore(store.Bucket(tokenSecretBucketName), "waf"),
	}
}

// SigningKey returns the persistent signing secret, generating and
// storing one on first use.
func (t *TokenStore) SigningKey() ([]byte, error) {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.cached != nil {
		return t.cached, nil
	}
	raw, err := t.GetRaw(tokenSecretKey)
	if err != nil {
		return nil, err
	}
	if len(raw) > 0 {
		t.cached = raw
		return raw, nil
	}
	secret := make([]byte, tokenSecretBytes)
	if _, err := rand.Read(secret); err != nil {
		return nil, err
	}
	if err := t.PutRaw(tokenSecretKey, secret); err != nil {
		return nil, err
	}
	t.cached = secret
	return secret, nil
}
