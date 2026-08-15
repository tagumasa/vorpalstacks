package cognitoidentityprovider

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"strings"
	"time"

	"github.com/google/uuid"

	"vorpalstacks/internal/store/aws/common"
)

// CreateUserPool creates a new Cognito user pool.
func (s *CognitoStore) CreateUserPool(userPool *UserPool) (*UserPool, error) {
	// Generate RSA key outside createMu to avoid blocking CreateUserPool
	// and CreateGroup operations during the 100-300ms key generation.
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}

	s.createMu.Lock()
	defer s.createMu.Unlock()
	if userPool.Name == "" {
		return nil, ErrInvalidUserPoolName
	}

	if s.Exists(userPool.ID) {
		return nil, ErrUserPoolAlreadyExists
	}

	now := time.Now().UTC()
	userPool.Arn = s.buildUserPoolArn(userPool.ID)
	userPool.CreationDate = now
	userPool.LastModifiedDate = now
	userPool.JwtKeyID = uuid.New().String()[:8]
	userPool.JwtPrivateKey = encodePrivateKeyToPEM(privateKey)
	userPool.JwtPublicKey = encodePublicKeyToPEM(&privateKey.PublicKey)
	if userPool.Status == "" {
		userPool.Status = "ACTIVE"
	}

	if err := s.Put(userPool.ID, userPool); err != nil {
		return nil, err
	}

	return userPool, nil
}

// GetUserPool retrieves a Cognito user pool by ID.
func (s *CognitoStore) GetUserPool(userPoolID string) (*UserPool, error) {
	var userPool UserPool
	if err := s.BaseStore.Get(userPoolID, &userPool); err != nil {
		return nil, ErrUserPoolNotFound
	}
	return &userPool, nil
}

// UpdateUserPool updates an existing Cognito user pool.
func (s *CognitoStore) UpdateUserPool(userPool *UserPool) error {
	if !s.Exists(userPool.ID) {
		return ErrUserPoolNotFound
	}
	userPool.LastModifiedDate = time.Now().UTC()
	return s.Put(userPool.ID, userPool)
}

// DeleteUserPool deletes a Cognito user pool by ID and cascades to users,
// groups, clients, tokens, and challenge sessions.
func (s *CognitoStore) DeleteUserPool(userPoolID string) error {
	if !s.Exists(userPoolID) {
		return ErrUserPoolNotFound
	}

	users, _ := s.ListUsers(userPoolID)
	for _, u := range users {
		prefix := userPoolID + "#" + u.ID + "#"
		_ = s.refreshTokensStore.ScanPrefix(prefix, func(key string, _ []byte) error {
			return s.refreshTokensStore.Delete(key)
		})
		_ = s.idTokensStore.ScanPrefix(prefix, func(key string, _ []byte) error {
			return s.idTokensStore.Delete(key)
		})
		_ = s.accessTokensStore.ScanPrefix(prefix, func(key string, _ []byte) error {
			return s.accessTokensStore.Delete(key)
		})
		_ = s.DeleteUser(userPoolID, u.Username)
	}

	groups, _ := s.ListGroups(userPoolID)
	for _, g := range groups {
		_ = s.DeleteGroup(userPoolID, g.Name)
	}

	clients, _ := s.ListUserPoolClients(userPoolID)
	for _, c := range clients {
		_ = s.DeleteUserPoolClient(userPoolID, c.ClientID)
	}

	_ = s.challengeSessionsStore.ScanPrefix(userPoolID+"#", func(key string, _ []byte) error {
		return s.challengeSessionsStore.Delete(key)
	})

	// Cascade: delete resource servers.
	rsPrefix := resourceServerPrefix(userPoolID)
	_ = s.BaseStore.ScanPrefix(rsPrefix, func(key string, _ []byte) error {
		return s.BaseStore.Delete(key)
	})

	// Cascade: delete identity providers.
	idpPrefix := identityProviderPrefix(userPoolID)
	_ = s.BaseStore.ScanPrefix(idpPrefix, func(key string, _ []byte) error {
		return s.BaseStore.Delete(key)
	})

	// Cascade: delete domains associated with this pool.
	_ = s.BaseStore.ScanPrefix("domain:", func(key string, value []byte) error {
		var entry UserPoolDomain
		if err := json.Unmarshal(value, &entry); err == nil && entry.UserPoolID == userPoolID {
			return s.BaseStore.Delete(key)
		}
		return nil
	})

	if pool, err := s.GetUserPool(userPoolID); err == nil && pool.Arn != "" {
		_ = s.TagStore.Delete(pool.Arn)
	}

	return s.BaseStore.Delete(userPoolID)
}

// ListUserPools lists all Cognito user pools.
func (s *CognitoStore) ListUserPools() ([]*UserPool, error) {
	var userPools []*UserPool
	// The user pool bucket also contains other entities (domains, identity
	// providers, etc.) keyed by prefixed names like "domain:..." or
	// "identityprovider:...". Only keys matching the user pool ID format
	// ({region}_{uuid}) are actual user pool entries.
	err := s.ForEach(func(key string, value []byte) error {
		if !strings.HasPrefix(key, s.region+"_") {
			return nil
		}
		var userPool UserPool
		if err := json.Unmarshal(value, &userPool); err != nil {
			return err
		}
		userPools = append(userPools, &userPool)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return userPools, nil
}

// ListUserPoolsPaginated lists Cognito user pools with server-side pagination.
// The user pool bucket contains other entities keyed by prefixed names; only
// keys matching the "{region}_{uuid}" pattern are actual user pool entries.
func (s *CognitoStore) ListUserPoolsPaginated(opts common.ListOptions) (*common.ListResult[UserPool], error) {
	return common.List[UserPool](s.BaseStore, opts, func(pool *UserPool) bool {
		return strings.HasPrefix(pool.ID, s.region+"_")
	})
}

// SetUserPoolDomain stores a domain configuration for a user pool.
func (s *CognitoStore) SetUserPoolDomain(domain string, entry *UserPoolDomain) error {
	return s.BaseStore.Put(domainKey(domain), entry)
}

// GetUserPoolDomain retrieves the domain configuration for a user pool.
func (s *CognitoStore) GetUserPoolDomain(domain string) (*UserPoolDomain, error) {
	var entry UserPoolDomain
	if err := s.BaseStore.Get(domainKey(domain), &entry); err != nil {
		return nil, ErrUserPoolNotFound
	}
	return &entry, nil
}

// GetUserPoolDomainByPool retrieves the custom domain for a user pool, if any.
func (s *CognitoStore) GetUserPoolDomainByPool(userPoolID string) (*UserPoolDomain, error) {
	var found *UserPoolDomain
	_ = s.BaseStore.ScanPrefix("domain:", func(key string, value []byte) error {
		var entry UserPoolDomain
		if err := json.Unmarshal(value, &entry); err == nil && entry.UserPoolID == userPoolID {
			found = &entry
		}
		return nil
	})
	if found == nil {
		return nil, ErrUserPoolNotFound
	}
	return found, nil
}

// DeleteUserPoolDomain removes the domain configuration for a user pool.
func (s *CognitoStore) DeleteUserPoolDomain(domain string) error {
	return s.BaseStore.Delete(domainKey(domain))
}
