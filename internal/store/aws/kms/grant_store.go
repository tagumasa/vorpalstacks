package kms

// Package kms provides KMS (Key Management Service) data store implementations
// for vorpalstacks.

import (
	"encoding/json"
	"errors"
	"time"

	arnutil "vorpalstacks/internal/utils/aws/arn"

	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/store/aws/common"
)

var errStopIteration = errors.New("stop iteration")

// GrantStore manages KMS grants.
type GrantStore struct {
	*common.BaseStore
	arnBuilder *ARNBuilder
	kl         common.KeyLocker
}

func grantBucketName(region string) string {
	return "kms_grants-" + region
}

// NewGrantStore creates a new GrantStore.
func NewGrantStore(store storage.BasicStorage, accountId, region string) *GrantStore {
	return &GrantStore{
		BaseStore:  common.NewBaseStore(store.Bucket(grantBucketName(region)), "kms"),
		arnBuilder: NewARNBuilder(accountId, region),
	}
}

// Create creates a new grant for a KMS key.
func (s *GrantStore) Create(keyID, granteePrincipal, retiringPrincipal string, operations []string, name string, constraints *GrantConstraints) (*Grant, error) {
	grantID, err := GenerateGrantID()
	if err != nil {
		return nil, err
	}

	grant := &Grant{
		GrantID:           grantID,
		KeyID:             keyID,
		GranteePrincipal:  granteePrincipal,
		RetiringPrincipal: retiringPrincipal,
		Operations:        operations,
		Name:              name,
		IssuingAccount:    arnutil.NewARNBuilder(s.arnBuilder.AccountId(), "").IAM().Root(),
		CreationDate:      time.Now(),
		Constraints:       constraints,
	}

	if err := s.save(grant); err != nil {
		return nil, err
	}

	return grant, nil
}

// CreateWithToken creates a new grant with a pre-generated token.
func (s *GrantStore) CreateWithToken(keyID, granteePrincipal, retiringPrincipal string, operations []string, name string, constraints *GrantConstraints, token string) (*Grant, error) {
	var grant *Grant
	err := s.kl.WithLock(keyID, func() error {
		var err error
		grant, err = s.Create(keyID, granteePrincipal, retiringPrincipal, operations, name, constraints)
		if err != nil {
			return err
		}
		grant.GrantToken = token
		return s.save(grant)
	})
	if err != nil {
		return nil, err
	}
	return grant, nil
}

// GetByToken retrieves a grant by its token.
func (s *GrantStore) GetByToken(token string) (*Grant, error) {
	var found *Grant
	err := s.ForEach(func(key string, value []byte) error {
		var grant Grant
		if err := json.Unmarshal(value, &grant); err != nil {
			return err
		}
		if grant.GrantToken == token {
			found = &grant
			return errStopIteration
		}
		return nil
	})
	if err != nil && err != errStopIteration {
		return nil, NewStoreError("get_grant_by_token", err)
	}
	if found == nil {
		return nil, NewStoreError("get_grant_by_token", ErrGrantNotFound)
	}
	return found, nil
}

// Get retrieves a grant by its ID.
func (s *GrantStore) Get(grantID string) (*Grant, error) {
	var grant Grant
	if err := s.BaseStore.Get(grantID, &grant); err != nil {
		return nil, err
	}
	return &grant, nil
}

func (s *GrantStore) save(grant *Grant) error {
	return s.BaseStore.Put(grant.GrantID, grant)
}

// Delete removes a grant by its ID.
func (s *GrantStore) Delete(grantID string) error {
	return s.kl.WithLock(grantID, func() error {
		return s.BaseStore.Delete(grantID)
	})
}

// List returns grants with optional filtering by key ID and grantee principal.
func (s *GrantStore) List(keyID string, granteePrincipal string, marker string, maxItems int) (*GrantListResult, error) {
	filter := func(g *Grant) bool {
		keyMatch := keyID == "" || g.KeyID == keyID
		principalMatch := granteePrincipal == "" || g.GranteePrincipal == granteePrincipal
		return keyMatch && principalMatch
	}
	result, err := common.List[Grant](s.BaseStore, common.ListOptions{Marker: marker, MaxItems: maxItems}, filter)
	if err != nil {
		return nil, err
	}
	return &GrantListResult{
		Grants:      result.Items,
		IsTruncated: result.IsTruncated,
		NextMarker:  result.NextMarker,
	}, nil
}

// ListByRetiringPrincipal returns grants that can be retired by a specific principal.
func (s *GrantStore) ListByRetiringPrincipal(principal string, marker string, maxItems int) (*GrantListResult, error) {
	filter := func(g *Grant) bool { return g.RetiringPrincipal == principal }
	result, err := common.List[Grant](s.BaseStore, common.ListOptions{Marker: marker, MaxItems: maxItems}, filter)
	if err != nil {
		return nil, err
	}
	return &GrantListResult{
		Grants:      result.Items,
		IsTruncated: result.IsTruncated,
		NextMarker:  result.NextMarker,
	}, nil
}

// DeleteByKeyID removes all grants for a specific key.
func (s *GrantStore) DeleteByKeyID(keyID string) error {
	return s.kl.WithLock(keyID, func() error {
		var grantIDs []string

		err := s.ForEach(func(key string, value []byte) error {
			var grant Grant
			if err := json.Unmarshal(value, &grant); err != nil {
				return err
			}
			if grant.KeyID == keyID {
				grantIDs = append(grantIDs, grant.GrantID)
			}
			return nil
		})

		if err != nil {
			return NewStoreError("delete_grants_by_key", err)
		}

		for _, id := range grantIDs {
			if err := s.BaseStore.Delete(id); err != nil {
				return NewStoreError("delete_grants_by_key", err)
			}
		}

		return nil
	})
}

// Count returns the number of grants for a key.
func (s *GrantStore) Count(keyID string) (int, error) {
	count := 0
	err := s.ForEach(func(key string, value []byte) error {
		var grant Grant
		if err := json.Unmarshal(value, &grant); err != nil {
			return err
		}
		if keyID == "" || grant.KeyID == keyID {
			count++
		}
		return nil
	})
	return count, err
}

// ARNBuilder returns the ARN builder for this store.
func (s *GrantStore) ARNBuilder() *ARNBuilder {
	return s.arnBuilder
}
