package kms

// Package kms provides KMS (Key Management Service) data store implementations
// for vorpalstacks.

import (
	"encoding/json"
	"time"

	"vorpalstacks/internal/core/storage"
)

// GrantStore manages KMS grants.
type GrantStore struct {
	bucket     storage.Bucket
	arnBuilder *ARNBuilder
}

func grantBucketName(region string) string {
	return "kms_grants-" + region
}

// NewGrantStore creates a new GrantStore.
func NewGrantStore(store storage.BasicStorage, accountId, region string) *GrantStore {
	return &GrantStore{
		bucket:     store.Bucket(grantBucketName(region)),
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
		IssuingAccount:    "arn:aws:iam::" + s.arnBuilder.AccountId() + ":root",
		CreationDate:      time.Now(),
		Constraints:       constraints,
	}

	if err := s.save(grant); err != nil {
		return nil, err
	}

	return grant, nil
}

// Get retrieves a grant by its ID.
func (s *GrantStore) Get(grantID string) (*Grant, error) {
	data, err := s.bucket.Get([]byte(grantID))
	if err != nil {
		return nil, NewStoreError("get_grant", err)
	}
	if data == nil {
		return nil, NewStoreError("get_grant", ErrGrantNotFound)
	}

	var grant Grant
	if err := json.Unmarshal(data, &grant); err != nil {
		return nil, NewStoreError("get_grant", err)
	}
	return &grant, nil
}

func (s *GrantStore) save(grant *Grant) error {
	data, err := json.Marshal(grant)
	if err != nil {
		return NewStoreError("save_grant", err)
	}
	return s.bucket.Put([]byte(grant.GrantID), data)
}

// Delete removes a grant by its ID.
func (s *GrantStore) Delete(grantID string) error {
	return s.bucket.Delete([]byte(grantID))
}

// List returns grants with optional filtering by key ID and grantee principal.
func (s *GrantStore) List(keyID string, granteePrincipal string, marker string, maxItems int) (*GrantListResult, error) {
	if maxItems <= 0 {
		maxItems = 100
	}

	var grants []*Grant
	count := 0
	started := marker == ""
	var lastGrantID string
	hasMore := false

	err := s.bucket.ForEach(func(k, v []byte) error {
		var grant Grant
		if err := json.Unmarshal(v, &grant); err != nil {
			return err
		}

		if !started {
			if grant.GrantID == marker {
				started = true
			}
			return nil
		}

		keyMatch := keyID == "" || grant.KeyID == keyID
		principalMatch := granteePrincipal == "" || grant.GranteePrincipal == granteePrincipal

		if keyMatch && principalMatch {
			if count < maxItems {
				grants = append(grants, &grant)
				count++
				lastGrantID = grant.GrantID
			} else {
				hasMore = true
			}
		}
		return nil
	})

	if err != nil {
		return nil, NewStoreError("list_grants", err)
	}

	result := &GrantListResult{
		Grants:      grants,
		IsTruncated: hasMore,
	}
	if result.IsTruncated {
		result.NextMarker = lastGrantID
	}

	return result, nil
}

// ListByRetiringPrincipal returns grants that can be retired by a specific principal.
func (s *GrantStore) ListByRetiringPrincipal(principal string, marker string, maxItems int) (*GrantListResult, error) {
	if maxItems <= 0 {
		maxItems = 100
	}

	var grants []*Grant
	count := 0
	started := marker == ""
	var lastGrantID string
	hasMore := false

	err := s.bucket.ForEach(func(k, v []byte) error {
		var grant Grant
		if err := json.Unmarshal(v, &grant); err != nil {
			return err
		}

		if !started {
			if grant.GrantID == marker {
				started = true
			}
			return nil
		}

		if grant.RetiringPrincipal == principal {
			if count < maxItems {
				grants = append(grants, &grant)
				count++
				lastGrantID = grant.GrantID
			} else {
				hasMore = true
			}
		}
		return nil
	})

	if err != nil {
		return nil, NewStoreError("list_retirable_grants", err)
	}

	result := &GrantListResult{
		Grants:      grants,
		IsTruncated: hasMore,
	}
	if result.IsTruncated {
		result.NextMarker = lastGrantID
	}

	return result, nil
}

// DeleteByKeyID removes all grants for a specific key.
func (s *GrantStore) DeleteByKeyID(keyID string) error {
	var grantIDs []string

	err := s.bucket.ForEach(func(k, v []byte) error {
		var grant Grant
		if err := json.Unmarshal(v, &grant); err != nil {
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
		if err := s.Delete(id); err != nil {
			return err
		}
	}

	return nil
}

// Count returns the number of grants for a key.
func (s *GrantStore) Count(keyID string) (int, error) {
	count := 0
	err := s.bucket.ForEach(func(k, v []byte) error {
		var grant Grant
		if err := json.Unmarshal(v, &grant); err != nil {
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
